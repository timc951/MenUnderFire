package routes

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"menunderfire/internal/handlers"
	"menunderfire/internal/logger"
	"menunderfire/internal/middleware"

	"github.com/gorilla/mux"
)

// extractEmailFromJWT tries to decode the JWT payload to extract the email for logging.
// This does NOT validate the token - it's purely for logging purposes.
func extractEmailFromJWT(authHeader string) string {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "anonymous"
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "anonymous"
	}
	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "anonymous"
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "anonymous"
	}
	if email, ok := claims["email"].(string); ok && email != "" {
		return email
	}
	return "anonymous"
}

// loggingMiddleware logs all incoming requests with their headers, geo location, and user info
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user email from JWT if present (for logging only)
		userEmail := extractEmailFromJWT(r.Header.Get("Authorization"))

		// Extract Cloudflare geo headers
		cfCountry := r.Header.Get("CF-IPCountry")
		cfRegion := r.Header.Get("CF-Region")

		logEvent := logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Str("user", userEmail)

		if cfCountry != "" {
			logEvent = logEvent.Str("country", cfCountry)
		}
		if cfRegion != "" {
			logEvent = logEvent.Str("region", cfRegion)
		}

		logEvent.Msg("Incoming request")
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware adds CORS headers to allow requests from the configured origin
func corsMiddleware(allowedOrigin string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With, Origin")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")

			// Handle preflight OPTIONS requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Setup(
	userHandler *handlers.UserHandler,
	groupHandler *handlers.GroupHandler,
	groupMessageHandler *handlers.GroupMessageHandler,
	reportHandler *handlers.ReportHandler,
	organizationHandler *handlers.OrganizationHandler,
	invitationHandler *handlers.InvitationHandler,
	formHandler *handlers.FormHandler,
	sitePageHandler *handlers.SitePageHandler,
	dashboardHandler *handlers.DashboardHandler,
	apiKeyHandler *handlers.APIKeyHandler,
	feedbackHandler *handlers.FeedbackHandler,
	pageHitHandler *handlers.PageHitHandler,
	authMiddleware *middleware.AuthMiddleware,
	corsOrigin string,
	trustProxyHeaders bool,
) *mux.Router {

	router := mux.NewRouter()

	// Add logging middleware first to see all requests
	router.Use(loggingMiddleware)
	// Add CORS middleware to all routes
	router.Use(corsMiddleware(corsOrigin))

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Status line is already committed; a failed body write is unrecoverable here.
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// API routes (base)
	api := router.PathPrefix("/api").Subrouter()

	// Global OPTIONS handler for CORS preflight requests
	api.PathPrefix("/").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// ===================
	// PUBLIC ROUTES (No authentication required)
	// ===================

	// Invitation public routes
	api.HandleFunc("/invitations/accept", invitationHandler.Accept).Methods("POST")
	api.HandleFunc("/invitations/validate/{token}", invitationHandler.Validate).Methods("GET")

	// Hit tracking public route (rate limited: 60 requests per minute per IP)
	hitRateLimiter := middleware.NewRateLimiter(60, 1*time.Minute, trustProxyHeaders)
	api.Handle("/hits", hitRateLimiter.Middleware(http.HandlerFunc(pageHitHandler.Record))).Methods("POST")

	// Site page public routes
	api.HandleFunc("/pages", sitePageHandler.List).Methods("GET")
	api.HandleFunc("/pages/preview/{draftId}", sitePageHandler.GetDraft).Methods("GET") // Preview draft (public)
	api.HandleFunc("/pages/{slug}", sitePageHandler.GetBySlug).Methods("GET")

	// ===================
	// PROTECTED ROUTES (Authentication required)
	// ===================

	// Create protected subrouter with auth middleware
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(authMiddleware.Middleware)

	// ===================
	// User Routes
	// ===================
	protected.HandleFunc("/users/me", userHandler.GetMe).Methods("GET")
	protected.HandleFunc("/users/me/permissions", userHandler.GetPermissions).Methods("GET")
	protected.HandleFunc("/users/me", userHandler.UpdateMe).Methods("PUT")
	protected.HandleFunc("/users/me/accept-agreement", userHandler.AcceptAgreement).Methods("POST")

	// ===================
	// Group Routes
	// ===================
	protected.HandleFunc("/groups", groupHandler.Create).Methods("POST")
	protected.HandleFunc("/groups", groupHandler.List).Methods("GET")
	protected.HandleFunc("/groups/join", groupHandler.JoinByCode).Methods("POST")
	protected.HandleFunc("/groups/{id}", groupHandler.Get).Methods("GET")
	protected.HandleFunc("/groups/{id}/join", groupHandler.Join).Methods("POST")
	protected.HandleFunc("/groups/{id}/members/{userId}", groupHandler.RemoveMember).Methods("DELETE")
	protected.HandleFunc("/groups/{id}/members/{userId}/role", groupHandler.UpdateMemberRole).Methods("PUT")
	protected.HandleFunc("/groups/{id}/settings", groupHandler.UpdateSettings).Methods("PUT")

	// ===================
	// Group Message Routes
	// ===================
	protected.HandleFunc("/groups/{groupId}/form-reports", groupMessageHandler.ListGroupForms).Methods("GET")
	protected.HandleFunc("/groups/{groupId}/form-reports/{formId}", groupMessageHandler.ListGroupFormAnswers).Methods("GET")
	protected.HandleFunc("/groups/{groupId}/messages", groupMessageHandler.List).Methods("GET")
	protected.HandleFunc("/groups/{groupId}/messages", groupMessageHandler.Send).Methods("POST")
	protected.HandleFunc("/groups/{groupId}/messages/{messageId}", groupMessageHandler.Delete).Methods("DELETE")

	// ===================
	// Report Routes
	// ===================
	protected.HandleFunc("/reports", reportHandler.Create).Methods("POST")
	protected.HandleFunc("/reports", reportHandler.List).Methods("GET") // Query param: groupId

	// ===================
	// Organization Routes
	// ===================
	protected.HandleFunc("/organizations", organizationHandler.List).Methods("GET")
	protected.HandleFunc("/organizations/all", organizationHandler.ListAll).Methods("GET") // Site Admin only
	protected.HandleFunc("/organizations", organizationHandler.Create).Methods("POST")     // Site Admin only
	protected.HandleFunc("/organizations/{id}", organizationHandler.Get).Methods("GET")
	protected.HandleFunc("/organizations/{id}", organizationHandler.Update).Methods("PUT") // Admin only
	protected.HandleFunc("/organizations/{id}/admins", organizationHandler.ListAdmins).Methods("GET")
	protected.HandleFunc("/organizations/{id}/groups", organizationHandler.ListGroups).Methods("GET")

	// ===================
	// Invitation Routes (Protected)
	// ===================
	protected.HandleFunc("/invitations", invitationHandler.List).Methods("GET")                            // Role-based filtering
	protected.HandleFunc("/invitations/org-admin", invitationHandler.CreateOrgAdmin).Methods("POST")       // Site Admin only
	protected.HandleFunc("/invitations/group-owner", invitationHandler.CreateGroupOwner).Methods("POST")   // Site Admin or Org Admin
	protected.HandleFunc("/invitations/group-member", invitationHandler.CreateGroupMember).Methods("POST") // Site Admin, Org Admin, or Group Owner
	protected.HandleFunc("/invitations/{id}", invitationHandler.Delete).Methods("DELETE")

	// ===================
	// Form Routes
	// ===================
	// Organization forms
	protected.HandleFunc("/organizations/{orgId}/forms", formHandler.ListByOrg).Methods("GET")
	protected.HandleFunc("/organizations/{orgId}/forms", formHandler.Create).Methods("POST") // Org Admin only

	// Form management
	protected.HandleFunc("/forms/{id}", formHandler.Get).Methods("GET")
	protected.HandleFunc("/forms/{id}", formHandler.Update).Methods("PUT")
	protected.HandleFunc("/forms/{id}", formHandler.Delete).Methods("DELETE")

	// Form field management
	protected.HandleFunc("/forms/{id}/fields", formHandler.AddField).Methods("POST")
	protected.HandleFunc("/forms/{id}/fields/reorder", formHandler.ReorderFields).Methods("PUT")
	protected.HandleFunc("/forms/{formId}/fields/{fieldId}", formHandler.UpdateField).Methods("PUT")
	protected.HandleFunc("/forms/{formId}/fields/{fieldId}", formHandler.DeleteField).Methods("DELETE")

	// Form answer management
	protected.HandleFunc("/forms/{id}/answers", formHandler.SubmitAnswer).Methods("POST")
	protected.HandleFunc("/forms/{id}/answers", formHandler.ListAnswers).Methods("GET")
	protected.HandleFunc("/forms/{id}/answers/me", formHandler.GetMyAnswer).Methods("GET")
	protected.HandleFunc("/forms/{id}/answers/history/{userId}", formHandler.GetAnswerHistory).Methods("GET")

	// ===================
	// Site Page Routes (Protected - Site Admin only)
	// ===================
	protected.HandleFunc("/pages", sitePageHandler.Create).Methods("POST")
	protected.HandleFunc("/pages/{id}", sitePageHandler.Update).Methods("PUT")
	protected.HandleFunc("/pages/{id}", sitePageHandler.Delete).Methods("DELETE")
	protected.HandleFunc("/pages/{id}/preview", sitePageHandler.CreateDraft).Methods("POST") // Create preview draft

	// ===================
	// Dashboard Routes
	// ===================
	protected.HandleFunc("/dashboard/stats", dashboardHandler.GetStats).Methods("GET") // Admin only

	// ===================
	// API Key Routes (Currently disabled in Kotlin backend)
	// ===================
	protected.HandleFunc("/api-keys", apiKeyHandler.Create).Methods("POST")
	protected.HandleFunc("/api-keys", apiKeyHandler.List).Methods("GET")
	protected.HandleFunc("/api-keys/{id}", apiKeyHandler.Delete).Methods("DELETE")

	// ===================
	// Hit Tracking Routes (Protected - Site Admin only)
	// ===================
	protected.HandleFunc("/hits/summary", pageHitHandler.GetSummary).Methods("GET")
	protected.HandleFunc("/hits/recent", pageHitHandler.GetRecent).Methods("GET")
	protected.HandleFunc("/hits/by-country", pageHitHandler.GetByCountry).Methods("GET")
	protected.HandleFunc("/hits/by-region", pageHitHandler.GetByRegion).Methods("GET")
	protected.HandleFunc("/hits/hourly", pageHitHandler.GetHourly).Methods("GET")
	protected.HandleFunc("/hits/daily", pageHitHandler.GetDaily).Methods("GET")
	protected.HandleFunc("/hits/range", pageHitHandler.GetHitsInRange).Methods("GET")

	// ===================
	// Feedback Routes
	// ===================
	protected.HandleFunc("/feedback", feedbackHandler.Create).Methods("POST")
	protected.HandleFunc("/feedback/me", feedbackHandler.ListMine).Methods("GET")
	protected.HandleFunc("/feedback", feedbackHandler.ListAll).Methods("GET")                    // Site Admin only
	protected.HandleFunc("/feedback/{id}/status", feedbackHandler.UpdateStatus).Methods("PATCH") // Site Admin only

	return router
}
