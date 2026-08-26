package repositories

import (
	"context"
	"encoding/json"

	"menunderfire/internal/models"
)

// FormRepository defines the interface for form data access
type FormRepository interface {
	// FindByID retrieves a form by its ID
	FindByID(ctx context.Context, formID string) (*models.Form, error)

	// FindByOrganizationID retrieves all forms for an organization
	FindByOrganizationID(ctx context.Context, orgID string) ([]*models.Form, error)

	// Create creates a new form
	Create(ctx context.Context, orgID, name string, description *string, createdBy string) (*models.Form, error)

	// Update updates a form
	Update(ctx context.Context, formID, name string, description *string, isActive bool) (*models.Form, error)

	// Delete deletes a form
	Delete(ctx context.Context, formID string) error
}

// FormFieldRepository defines the interface for form field data access
type FormFieldRepository interface {
	// GetFields retrieves all fields for a form
	GetFields(ctx context.Context, formID string) ([]models.FormField, error)

	// AddField adds a field to a form
	AddField(ctx context.Context, formID string, fieldType models.FormFieldType, label string, description *string, isRequired bool, fieldOrder int, options []string) (*models.FormField, error)

	// UpdateField updates a form field
	UpdateField(ctx context.Context, fieldID, label string, description *string, isRequired bool, options []string) (*models.FormField, error)

	// DeleteField deletes a form field
	DeleteField(ctx context.Context, fieldID string) error

	// ReorderFields reorders fields in a form
	ReorderFields(ctx context.Context, formID string, fieldIDs []string) error
}

// FormAnswerRepository defines the interface for form answer data access
type FormAnswerRepository interface {
	// Submit submits answers to a form
	Submit(ctx context.Context, formID, userID string, answers json.RawMessage) (*models.FormAnswer, error)

	// FindByForm retrieves all answers for a form
	FindByForm(ctx context.Context, formID string) ([]*models.FormAnswer, error)

	// FindByFormAndUserIDs retrieves current answers for a form filtered to specific user IDs
	FindByFormAndUserIDs(ctx context.Context, formID string, userIDs []string) ([]*models.FormAnswer, error)

	// FindCurrent retrieves the current answer for a user on a form
	FindCurrent(ctx context.Context, formID, userID string) (*models.FormAnswer, error)

	// FindHistory retrieves the answer history for a user on a form
	FindHistory(ctx context.Context, formID, userID string) ([]*models.FormAnswer, error)
}
