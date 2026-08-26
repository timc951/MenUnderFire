package models

import "time"

// APIKey represents an API key entity
type APIKey struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	Name        string     `json:"name"`
	KeyHash     string     `json:"-"` // Never expose the hash
	Permissions []string   `json:"permissions"`
	ExpiresAt   *time.Time `json:"expiresAt"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// ===== Request DTOs =====

// CreateAPIKeyRequest is the request for creating an API key
type CreateAPIKeyRequest struct {
	Name          string   `json:"name" validate:"required,min=1,max=100"`
	Permissions   []string `json:"permissions"`
	ExpiresInDays int      `json:"expiresInDays"`
}

// ===== Response DTOs =====

// APIKeyResponse is the response for an API key
// The raw key is only included on initial creation
type APIKeyResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Key         *string  `json:"key,omitempty"` // Only returned once at creation time
	Permissions []string `json:"permissions"`
	ExpiresAt   *string  `json:"expiresAt,omitempty"`
	CreatedAt   string   `json:"createdAt"`
}

// ===== Conversion Methods =====

// ToResponse converts an APIKey entity to an APIKeyResponse DTO
// rawKey should only be provided during creation
func (k *APIKey) ToResponse(rawKey *string) *APIKeyResponse {
	var expiresAt *string
	if k.ExpiresAt != nil {
		exp := k.ExpiresAt.Format(time.RFC3339)
		expiresAt = &exp
	}
	return &APIKeyResponse{
		ID:          k.ID,
		Name:        k.Name,
		Key:         rawKey,
		Permissions: k.Permissions,
		ExpiresAt:   expiresAt,
		CreatedAt:   k.CreatedAt.Format(time.RFC3339),
	}
}
