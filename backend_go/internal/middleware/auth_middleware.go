package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"menunderfire/internal/config"
	"menunderfire/internal/handlers"
	"menunderfire/internal/logger"
	"menunderfire/internal/models"
	"menunderfire/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

// AuthMiddleware handles JWT validation and user provisioning
type AuthMiddleware struct {
	config         *config.Config
	userRepo       repositories.UserRepository
	jwksCache      map[string]*rsa.PublicKey
	jwksCacheMutex sync.RWMutex
	jwksCacheTime  time.Time
	log            zerolog.Logger
}

// JWKS represents the JSON Web Key Set structure
type JWKS struct {
	Keys []JSONWebKey `json:"keys"`
}

// JSONWebKey represents a single key in the JWKS
type JSONWebKey struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(cfg *config.Config, userRepo repositories.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		config:    cfg,
		userRepo:  userRepo,
		jwksCache: make(map[string]*rsa.PublicKey),
		log:       logger.WithComponent("auth"),
	}
}

// Middleware returns the HTTP middleware function
func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		am.log.Debug().Str("path", r.URL.Path).Msg("Processing request")

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			am.log.Debug().Msg("No Authorization header present")
			respondUnauthorized(w, "Missing authorization header")
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			am.log.Debug().Msg("Authorization header doesn't start with 'Bearer '")
			respondUnauthorized(w, "Invalid authorization header format")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		am.log.Debug().Int("token_length", len(tokenString)).Msg("Token extracted")

		// Validate JWT token
		token, err := am.validateToken(tokenString)
		if err != nil {
			am.log.Warn().Err(err).Msg("Token validation failed")
			respondUnauthorized(w, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			am.log.Warn().Msg("Failed to parse token claims")
			respondUnauthorized(w, "Invalid token claims")
			return
		}

		// Extract user information from token
		externalID, ok := claims["sub"].(string)
		if !ok || externalID == "" {
			am.log.Warn().Msg("No subject (sub) claim in token")
			respondUnauthorized(w, "Invalid token: missing subject")
			return
		}

		email := am.extractEmail(claims)
		displayName := am.extractName(claims)

		am.log.Debug().
			Str("sub", externalID).
			Str("email", email).
			Str("name", displayName).
			Msg("Token validated")

		// Get or create user
		user, err := am.getOrCreateUser(r.Context(), externalID, email, displayName)
		if err != nil {
			am.log.Error().Err(err).Msg("Failed to get/create user")
			respondInternalError(w, "Failed to process user")
			return
		}

		am.log.Debug().
			Str("user_id", user.ID).
			Str("email", user.Email).
			Msg("User loaded")

		// Set user in context
		ctx := handlers.SetUserInContext(r.Context(), user)

		// Continue with the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateToken validates a JWT token and returns the parsed token
func (am *AuthMiddleware) validateToken(tokenString string) (*jwt.Token, error) {
	// Parse token to get the kid (key ID) from header
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get kid from token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("no kid in token header")
		}

		am.log.Debug().Str("kid", kid).Msg("Token kid")

		// Get public key for this kid
		publicKey, err := am.getPublicKey(kid)
		if err != nil {
			return nil, fmt.Errorf("failed to get public key: %w", err)
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	// Verify issuer
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	issuer, ok := claims["iss"].(string)
	if !ok || issuer != am.config.AuthIssuer {
		return nil, fmt.Errorf("invalid issuer: got %s, expected %s", issuer, am.config.AuthIssuer)
	}

	// Verify audience if configured
	if am.config.AuthAudience != "" && am.config.AuthAudience != "none" {
		if !am.verifyAudience(claims, am.config.AuthAudience) {
			return nil, errors.New("invalid audience")
		}
	}

	return token, nil
}

// verifyAudience checks if the token's audience matches the expected audience
func (am *AuthMiddleware) verifyAudience(claims jwt.MapClaims, expectedAudience string) bool {
	aud, ok := claims["aud"]
	if !ok {
		return false
	}

	switch v := aud.(type) {
	case string:
		return v == expectedAudience
	case []interface{}:
		for _, a := range v {
			if audStr, ok := a.(string); ok && audStr == expectedAudience {
				return true
			}
		}
	}

	return false
}

// getPublicKey retrieves the public key for the given kid from JWKS
func (am *AuthMiddleware) getPublicKey(kid string) (*rsa.PublicKey, error) {
	// Check cache first (with read lock)
	am.jwksCacheMutex.RLock()
	if key, exists := am.jwksCache[kid]; exists && time.Since(am.jwksCacheTime) < 24*time.Hour {
		am.jwksCacheMutex.RUnlock()
		am.log.Debug().Str("kid", kid).Msg("Using cached public key")
		return key, nil
	}
	am.jwksCacheMutex.RUnlock()

	// Fetch JWKS
	var jwksURL string
	if strings.Contains(am.config.AuthIssuer, "/realms/") {
		// Keycloak pattern: use AUTH_DOMAIN for JWKS endpoint
		protocol := "http"
		if strings.HasPrefix(am.config.AuthDomain, "https://") {
			protocol = "https"
		}
		domain := strings.TrimPrefix(strings.TrimPrefix(am.config.AuthDomain, "https://"), "http://")
		parts := strings.Split(am.config.AuthIssuer, "/realms/")
		realm := ""
		if len(parts) > 1 {
			realm = parts[1]
		}
		jwksURL = fmt.Sprintf("%s://%s/realms/%s/protocol/openid-connect/certs", protocol, domain, realm)
	} else {
		// Auth0/generic OIDC pattern
		protocol := "https"
		if strings.HasPrefix(am.config.AuthDomain, "localhost") || strings.HasPrefix(am.config.AuthDomain, "127.0.0.1") {
			protocol = "http"
		}
		jwksURL = fmt.Sprintf("%s://%s/.well-known/jwks.json", protocol, am.config.AuthDomain)
	}

	am.log.Debug().Str("url", jwksURL).Msg("Fetching JWKS")

	// Bounded timeout: an unresponsive IdP must not pin request goroutines.
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(jwksURL) // #nosec G107 -- URL is derived from operator-supplied AUTH_* config, not user input
	if err != nil {
		am.log.Error().Err(err).Str("url", jwksURL).Msg("JWKS fetch failed")
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			am.log.Debug().Err(cerr).Msg("JWKS body close failed")
		}
	}()

	am.log.Debug().Int("status", resp.StatusCode).Msg("JWKS response")
	if resp.StatusCode != http.StatusOK {
		am.log.Error().Int("status", resp.StatusCode).Msg("JWKS fetch returned non-200 status")
		return nil, fmt.Errorf("JWKS endpoint returned status %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		am.log.Error().Err(err).Msg("JWKS decode failed")
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	am.log.Debug().Int("key_count", len(jwks.Keys)).Msg("JWKS loaded")

	// Find the key with matching kid
	for _, key := range jwks.Keys {
		if key.Kid == kid && key.Kty == "RSA" {
			am.log.Debug().Str("kid", kid).Msg("Found matching key")
			publicKey, err := am.parseRSAPublicKey(key)
			if err != nil {
				am.log.Error().Err(err).Str("kid", kid).Msg("RSA key parse failed")
				return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
			}

			// Cache the key (with write lock)
			am.jwksCacheMutex.Lock()
			am.jwksCache[kid] = publicKey
			am.jwksCacheTime = time.Now()
			am.jwksCacheMutex.Unlock()

			am.log.Info().Str("kid", kid).Msg("Public key cached")
			return publicKey, nil
		}
	}

	am.log.Error().Str("kid", kid).Int("searched", len(jwks.Keys)).Msg("No matching key found")
	return nil, fmt.Errorf("no matching key found for kid: %s", kid)
}

// parseRSAPublicKey converts a JWK to an RSA public key
func (am *AuthMiddleware) parseRSAPublicKey(key JSONWebKey) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)

	var e int
	for _, b := range eBytes {
		e = e*256 + int(b)
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

// extractEmail extracts the email from JWT claims
func (am *AuthMiddleware) extractEmail(claims jwt.MapClaims) string {
	if email, ok := claims["email"].(string); ok && email != "" {
		return email
	}
	if email, ok := claims["https://menunderfire.com/email"].(string); ok && email != "" {
		return email
	}
	return ""
}

// extractName extracts the name from JWT claims
func (am *AuthMiddleware) extractName(claims jwt.MapClaims) string {
	if name, ok := claims["name"].(string); ok && name != "" {
		return name
	}
	if name, ok := claims["https://menunderfire.com/name"].(string); ok && name != "" {
		return name
	}
	return am.extractEmail(claims)
}

// getOrCreateUser gets an existing user or creates a new one (auto-provisioning)
func (am *AuthMiddleware) getOrCreateUser(ctx context.Context, externalID, email, displayName string) (*models.User, error) {
	// Try finding by external ID first (most common case for returning users)
	user, err := am.userRepo.FindByExternalID(ctx, externalID)
	if err == nil {
		am.log.Info().
			Str("user_id", user.ID).
			Str("db_email", user.Email).
			Str("db_external_id", user.ExternalID).
			Str("jwt_email", email).
			Str("jwt_external_id", externalID).
			Bool("is_site_admin", user.IsSiteAdmin).
			Msg("User found by external ID")
		return user, nil
	}

	// Fallback: try finding by email (handles invitation-created users whose
	// external_id may not match the JWT sub claim)
	if email != "" {
		user, err = am.userRepo.FindByEmail(ctx, email)
		if err == nil {
			am.log.Info().
				Str("user_id", user.ID).
				Str("email", email).
				Str("external_id", externalID).
				Msg("Found user by email, linking external ID")

			if err := am.userRepo.UpdateExternalID(ctx, user.ID, externalID); err != nil {
				am.log.Error().Err(err).Msg("Failed to update external ID")
				return nil, fmt.Errorf("failed to link external ID: %w", err)
			}
			user.ExternalID = externalID
			return user, nil
		}
	}

	am.log.Info().
		Str("external_id", externalID).
		Str("email", email).
		Str("display_name", displayName).
		Msg("Auto-provisioning new user")

	if displayName == "" {
		displayName = email
	}
	if displayName == "" {
		displayName = externalID // Fallback if no email/name in token
	}

	user, err = am.userRepo.Create(ctx, email, displayName, externalID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	am.log.Info().Str("user_id", user.ID).Msg("User auto-provisioned")

	return user, nil
}

// respondUnauthorized writes a 401 Unauthorized response
func respondUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	// Status line is already committed; a failed body write is unrecoverable here.
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   "unauthorized",
		"message": message,
	})
}

// respondInternalError writes a 500 Internal Server Error response
func respondInternalError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	// Status line is already committed; a failed body write is unrecoverable here.
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   "internal_error",
		"message": message,
	})
}
