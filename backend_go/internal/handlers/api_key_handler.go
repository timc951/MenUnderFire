package handlers

import (
	"net/http"

	"menunderfire/internal/services"
)

// APIKeyHandler handles HTTP requests for API key endpoints
// NOTE: All API key operations are currently disabled (matching Kotlin backend behavior)
type APIKeyHandler struct {
	apiKeyService services.APIKeyService
}

// NewAPIKeyHandler creates a new APIKeyHandler with the given service
func NewAPIKeyHandler(apiKeyService services.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

// Create handles POST /api/api-keys
// Creates a new API key
// NOTE: Currently disabled - returns 403 Forbidden
func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	respondForbidden(w, "API key operations are currently disabled")
}

// List handles GET /api/api-keys
// Returns all API keys for the current user
// NOTE: Currently disabled - returns 403 Forbidden
func (h *APIKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	respondForbidden(w, "API key operations are currently disabled")
}

// Delete handles DELETE /api/api-keys/{id}
// Deletes an API key
// NOTE: Currently disabled - returns 403 Forbidden
func (h *APIKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	respondForbidden(w, "API key operations are currently disabled")
}

// The following commented code shows the full implementation when API keys are enabled:
/*
func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if req.Name == "" {
		respondBadRequest(w, "name is required")
		return
	}

	if req.ExpiresInDays <= 0 {
		req.ExpiresInDays = 90 // default
	}

	apiKey, rawKey, err := h.apiKeyService.Create(ctx, user.ID, &req)
	if err != nil {
		log.Printf("Error creating API key: %v", err)
		respondInternalError(w, "Failed to create API key")
		return
	}

	respondCreated(w, apiKey.ToResponse(&rawKey))
}

func (h *APIKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	keys, err := h.apiKeyService.List(ctx, user.ID)
	if err != nil {
		log.Printf("Error listing API keys for user %s: %v", user.ID, err)
		respondInternalError(w, "Failed to list API keys")
		return
	}

	responses := make([]*models.APIKeyResponse, len(keys))
	for i, key := range keys {
		responses[i] = key.ToResponse(nil) // Never include raw key in list
	}

	respondOK(w, responses)
}

func (h *APIKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	keyID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	err = h.apiKeyService.Delete(ctx, keyID, user.ID)
	if err != nil {
		log.Printf("Error deleting API key %s: %v", keyID, err)
		if errors.Is(err, services.ErrNotFound) {
			respondNotFound(w, "API key not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to delete this API key")
			return
		}
		respondInternalError(w, "Failed to delete API key")
		return
	}

	respondNoContent(w)
}
*/
