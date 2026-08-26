package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"menunderfire/internal/logger"
	"menunderfire/internal/models"
	"menunderfire/internal/services"
)

// parseDateRange extracts optional "from" and "to" query params (YYYY-MM-DD format)
func parseDateRange(r *http.Request) (from, to *time.Time) {
	if f := r.URL.Query().Get("from"); f != "" {
		if t, err := time.Parse("2006-01-02", f); err == nil {
			from = &t
		}
	}
	if t := r.URL.Query().Get("to"); t != "" {
		// Set "to" to end of day (23:59:59)
		if parsed, err := time.Parse("2006-01-02", t); err == nil {
			endOfDay := parsed.Add(24*time.Hour - time.Second)
			to = &endOfDay
		}
	}
	return
}

// PageHitHandler handles HTTP requests for page hit tracking endpoints
type PageHitHandler struct {
	hitService services.PageHitService
	hitToken   string
}

// NewPageHitHandler creates a new PageHitHandler.
// hitToken is a shared secret the frontend sends as X-Hit-Token to prevent casual abuse.
func NewPageHitHandler(hitService services.PageHitService, hitToken string) *PageHitHandler {
	return &PageHitHandler{
		hitService: hitService,
		hitToken:   hitToken,
	}
}

// Record handles POST /api/hits (public - no auth required)
func (h *PageHitHandler) Record(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Validate hit token to prevent casual abuse
	if h.hitToken != "" {
		if r.Header.Get("X-Hit-Token") != h.hitToken {
			respondForbidden(w, "Invalid hit token")
			return
		}
	}

	// Limit request body to 4 KB to prevent payload abuse
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var req models.RecordHitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	// Validate path: must start with / and be reasonable length
	req.Path = strings.TrimSpace(req.Path)
	if req.Path == "" || req.Path[0] != '/' || len(req.Path) > 500 {
		respondBadRequest(w, "Invalid path")
		return
	}

	// Extract IP address from Cloudflare/proxy headers
	ipAddress := r.Header.Get("CF-Connecting-IP")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Real-Ip")
	}
	if ipAddress == "" {
		ipAddress = strings.Split(r.RemoteAddr, ":")[0]
	}
	// Take first IP if multiple in X-Forwarded-For
	if strings.Contains(ipAddress, ",") {
		ipAddress = strings.TrimSpace(strings.Split(ipAddress, ",")[0])
	}

	// Extract geo location from Cloudflare headers
	cfCountry := r.Header.Get("CF-IPCountry")
	cfRegion := r.Header.Get("CF-Region")

	userAgent := r.Header.Get("User-Agent")
	referrer := r.Header.Get("Referer")

	var ip, ua, ref, country, region *string
	if ipAddress != "" {
		ip = &ipAddress
	}
	if userAgent != "" {
		ua = &userAgent
	}
	if referrer != "" {
		ref = &referrer
	}
	if cfCountry != "" {
		country = &cfCountry
	}
	if cfRegion != "" {
		region = &cfRegion
	}

	// Try to get user from context (may not be authenticated on public route)
	var userID *string
	var userEmail *string
	user, err := GetUserFromContext(ctx)
	if err == nil && user != nil {
		userID = &user.ID
		if user.Email != "" {
			userEmail = &user.Email
		}
	}

	// Use email from request body if not from auth context
	if userEmail == nil && req.Email != "" {
		userEmail = &req.Email
	}

	// Determine display name for logging
	displayEmail := "anonymous"
	if userEmail != nil {
		displayEmail = *userEmail
	}

	logger.Info().
		Str("path", req.Path).
		Str("user", displayEmail).
		Str("ip", ipAddress).
		Str("country", cfCountry).
		Str("region", cfRegion).
		Msg("Page hit recorded")

	hit, err := h.hitService.RecordHit(ctx, req.Path, userID, userEmail, ip, ua, ref, country, region)
	if err != nil {
		logger.Error().Err(err).Msg("Error recording page hit")
		if errors.Is(err, services.ErrValidation) {
			respondBadRequest(w, err.Error())
			return
		}
		respondInternalError(w, "Failed to record hit")
		return
	}

	respondCreated(w, map[string]string{"id": hit.ID})
}

// GetSummary handles GET /api/hits/summary (site admin only)
func (h *PageHitHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	summary, err := h.hitService.GetSummary(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting hit summary")
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Only site admins can view hit statistics")
			return
		}
		respondInternalError(w, "Failed to get hit summary")
		return
	}

	respondOK(w, summary)
}

// GetRecent handles GET /api/hits/recent (site admin only)
func (h *PageHitHandler) GetRecent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	hits, err := h.hitService.GetRecentHits(ctx, user.ID, limit)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting recent hits")
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Only site admins can view hit details")
			return
		}
		respondInternalError(w, "Failed to get recent hits")
		return
	}

	respondOK(w, hits)
}

// GetByCountry handles GET /api/hits/by-country (site admin only)
func (h *PageHitHandler) GetByCountry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	from, to := parseDateRange(r)
	data, err := h.hitService.GetHitsByCountry(ctx, user.ID, from, to)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Only site admins can view hit statistics")
			return
		}
		respondInternalError(w, "Failed to get hits by country")
		return
	}
	respondOK(w, data)
}

// GetByRegion handles GET /api/hits/by-region (site admin only)
func (h *PageHitHandler) GetByRegion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	from, to := parseDateRange(r)
	data, err := h.hitService.GetHitsByRegion(ctx, user.ID, from, to)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Only site admins can view hit statistics")
			return
		}
		respondInternalError(w, "Failed to get hits by region")
		return
	}
	respondOK(w, data)
}

// GetHourly handles GET /api/hits/hourly (site admin only)
func (h *PageHitHandler) GetHourly(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	from, to := parseDateRange(r)
	data, err := h.hitService.GetHitsHourly(ctx, user.ID, from, to)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Only site admins can view hit statistics")
			return
		}
		respondInternalError(w, "Failed to get hourly hits")
		return
	}
	respondOK(w, data)
}

// GetDaily handles GET /api/hits/daily (site admin only)
func (h *PageHitHandler) GetDaily(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	from, to := parseDateRange(r)
	data, err := h.hitService.GetHitsDaily(ctx, user.ID, from, to)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Only site admins can view hit statistics")
			return
		}
		respondInternalError(w, "Failed to get daily hits")
		return
	}
	respondOK(w, data)
}

// GetHitsInRange handles GET /api/hits/range (site admin only)
func (h *PageHitHandler) GetHitsInRange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	from, to := parseDateRange(r)

	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	hits, err := h.hitService.GetRecentHitsInRange(ctx, user.ID, from, to, limit)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Only site admins can view hit details")
			return
		}
		respondInternalError(w, "Failed to get hits in range")
		return
	}
	respondOK(w, hits)
}
