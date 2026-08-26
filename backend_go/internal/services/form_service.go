package services

import (
	"context"
	"fmt"
	"strings"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// FormService defines the interface for form-related business logic
type FormService interface {
	// ListByOrg returns all forms for an organization
	ListByOrg(ctx context.Context, orgID string, userID string) ([]*models.Form, error)

	// Create creates a new form (Org Admin only)
	Create(ctx context.Context, orgID string, userID string, req *models.CreateFormRequest) (*models.Form, error)

	// GetByID retrieves a form by its ID
	GetByID(ctx context.Context, formID string) (*models.Form, error)

	// GetDetailByID retrieves a form with all its fields
	GetDetailByID(ctx context.Context, formID string, userID string) (*models.Form, error)

	// Update updates a form
	Update(ctx context.Context, formID string, userID string, req *models.UpdateFormRequest) (*models.Form, error)

	// Delete deletes a form
	Delete(ctx context.Context, formID string, userID string) error

	// AddField adds a field to a form
	AddField(ctx context.Context, formID string, userID string, req *models.AddFormFieldRequest) (*models.FormField, error)

	// ReorderFields reorders fields in a form
	ReorderFields(ctx context.Context, formID string, userID string, req *models.ReorderFieldsRequest) error

	// UpdateField updates a form field
	UpdateField(ctx context.Context, formID string, fieldID string, userID string, req *models.UpdateFormFieldRequest) (*models.FormField, error)

	// DeleteField deletes a form field
	DeleteField(ctx context.Context, formID string, fieldID string, userID string) error

	// SubmitAnswer submits answers to a form
	SubmitAnswer(ctx context.Context, formID string, userID string, req *models.SubmitFormAnswerRequest) (*models.FormAnswer, error)

	// ListAnswers returns all answers for a form
	ListAnswers(ctx context.Context, formID string, userID string) ([]*models.FormAnswer, error)

	// GetMyAnswer returns the current user's latest answer for a form
	GetMyAnswer(ctx context.Context, formID string, userID string) (*models.FormAnswer, error)

	// GetAnswerHistory returns the answer history for a specific user on a form
	GetAnswerHistory(ctx context.Context, formID string, targetUserID string, requesterID string) ([]*models.FormAnswer, error)
}

type formService struct {
	formRepo   repositories.FormRepository
	fieldRepo  repositories.FormFieldRepository
	answerRepo repositories.FormAnswerRepository
	orgRepo    repositories.OrganizationRepository
	userRepo   repositories.UserRepository
}

// NewFormService creates a new FormService implementation with the provided dependencies
func NewFormService(
	formRepo repositories.FormRepository,
	fieldRepo repositories.FormFieldRepository,
	answerRepo repositories.FormAnswerRepository,
	orgRepo repositories.OrganizationRepository,
	userRepo repositories.UserRepository,
) FormService {
	return &formService{
		formRepo:   formRepo,
		fieldRepo:  fieldRepo,
		answerRepo: answerRepo,
		orgRepo:    orgRepo,
		userRepo:   userRepo,
	}
}

// ListByOrg returns all forms for an organization
//
// Edge cases handled:
//   - User is not org admin or site admin: Returns ErrForbidden
//   - Empty orgID: Returns ErrForbidden
//   - User is org admin: Returns all forms
//   - User is site admin: Returns all forms
//   - Database error checking admin status: Returns wrapped error
//   - Database error fetching forms: Returns wrapped error
//   - Organization has no forms: Returns empty slice
func (s *formService) ListByOrg(ctx context.Context, orgID string, userID string) ([]*models.Form, error) {
	if orgID == "" {
		return nil, ErrForbidden
	}

	// Check if user is site admin
	siteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check site admin status: %w", err)
	}

	// Check if user is org admin
	orgAdmin, err := s.orgRepo.IsAdmin(ctx, orgID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org admin status: %w", err)
	}

	// User must be either org admin or site admin
	if !siteAdmin && !orgAdmin {
		return nil, ErrForbidden
	}

	// Fetch all forms for the organization
	forms, err := s.formRepo.FindByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forms for organization: %w", err)
	}

	return forms, nil
}

// Create creates a new form (Org Admin only)
//
// Edge cases handled:
//   - User is not org admin: Returns ErrForbidden
//   - Empty name after trim: Returns ErrValidation
//   - Name is whitespace only: Returns ErrValidation
//   - Valid name with leading/trailing whitespace: Trims and creates
//   - Nil description: Passes nil to repository
//   - Database error checking admin: Returns wrapped error
//   - Database error creating form: Returns wrapped error
func (s *formService) Create(ctx context.Context, orgID string, userID string, req *models.CreateFormRequest) (*models.Form, error) {
	// Check if user is org admin (NOT site admin - only org admin can create)
	isAdmin, err := s.orgRepo.IsAdmin(ctx, orgID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org admin status: %w", err)
	}
	if !isAdmin {
		return nil, ErrForbidden
	}

	// Trim and validate name
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrValidation
	}

	// Trim description if provided
	var description *string
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		description = &trimmed
	}

	// Create form via repository
	form, err := s.formRepo.Create(ctx, orgID, name, description, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create form: %w", err)
	}

	return form, nil
}

// GetByID retrieves a form by its ID
//
// Edge cases handled:
//   - Empty formID: Returns ErrFormNotFound
//   - Form not found: Returns ErrFormNotFound
//   - Database error: Returns wrapped error
//   - Form exists: Returns the form
func (s *formService) GetByID(ctx context.Context, formID string) (*models.Form, error) {
	if formID == "" {
		return nil, ErrFormNotFound
	}

	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to get form by ID: %w", err)
	}
	if form == nil {
		return nil, ErrFormNotFound
	}

	return form, nil
}

// GetDetailByID retrieves a form with all its fields
//
// Edge cases handled:
//   - Form not found: Returns ErrFormNotFound
//   - User not org admin or site admin: Returns ErrForbidden
//   - Form has no fields: Returns form with empty Fields slice
//   - Form has fields: Returns form with Fields populated
//   - Database error finding form: Returns wrapped error
//   - Database error checking admin: Returns wrapped error
//   - Database error getting fields: Returns wrapped error
func (s *formService) GetDetailByID(ctx context.Context, formID string, userID string) (*models.Form, error) {
	// Find form by ID
	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to find form: %w", err)
	}
	if form == nil {
		return nil, ErrFormNotFound
	}

	// Check if user is site admin
	siteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check site admin status: %w", err)
	}

	if !siteAdmin {
		// Check if user is org admin
		orgAdmin, err := s.orgRepo.IsAdmin(ctx, form.OrganizationID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check org admin status: %w", err)
		}

		if !orgAdmin {
			// Check if user is org member (e.g. group member filling out a form)
			orgMember, err := s.orgRepo.IsMember(ctx, form.OrganizationID, userID)
			if err != nil {
				return nil, fmt.Errorf("failed to check org membership: %w", err)
			}
			if !orgMember {
				return nil, ErrForbidden
			}
		}
	}

	// Get form fields
	fields, err := s.fieldRepo.GetFields(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to get form fields: %w", err)
	}

	// Attach fields to form
	form.Fields = fields

	return form, nil
}

// Update updates a form
//
// Edge cases handled:
//   - Form not found: Returns ErrFormNotFound
//   - User not org admin: Returns ErrForbidden
//   - Empty name after trim: Returns ErrValidation
//   - Valid update: Returns updated form
//   - Database error finding form: Returns wrapped error
//   - Database error checking admin: Returns wrapped error
//   - Database error updating: Returns wrapped error
func (s *formService) Update(ctx context.Context, formID string, userID string, req *models.UpdateFormRequest) (*models.Form, error) {
	// Find form by ID
	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to find form: %w", err)
	}
	if form == nil {
		return nil, ErrFormNotFound
	}

	// Check if user is org admin (NOT site admin)
	isAdmin, err := s.orgRepo.IsAdmin(ctx, form.OrganizationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org admin status: %w", err)
	}
	if !isAdmin {
		return nil, ErrForbidden
	}

	// Trim and validate name
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrValidation
	}

	// Trim description if provided
	var description *string
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		description = &trimmed
	}

	// Update form
	updated, err := s.formRepo.Update(ctx, formID, name, description, req.IsActive)
	if err != nil {
		return nil, fmt.Errorf("failed to update form: %w", err)
	}

	return updated, nil
}

// Delete deletes a form
//
// Edge cases handled:
//   - Form not found: Returns ErrFormNotFound
//   - User not org admin: Returns ErrForbidden
//   - Successful delete: Returns nil
//   - Database error finding form: Returns wrapped error
//   - Database error checking admin: Returns wrapped error
//   - Database error deleting: Returns wrapped error
func (s *formService) Delete(ctx context.Context, formID string, userID string) error {
	// Find form by ID
	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		return fmt.Errorf("failed to find form: %w", err)
	}
	if form == nil {
		return ErrFormNotFound
	}

	// Check if user is org admin
	isAdmin, err := s.orgRepo.IsAdmin(ctx, form.OrganizationID, userID)
	if err != nil {
		return fmt.Errorf("failed to check org admin status: %w", err)
	}
	if !isAdmin {
		return ErrForbidden
	}

	// Delete form
	if err := s.formRepo.Delete(ctx, formID); err != nil {
		return fmt.Errorf("failed to delete form: %w", err)
	}

	return nil
}

// AddField adds a field to a form
//
// Edge cases handled:
//   - Form not found: Returns ErrFormNotFound
//   - User not org admin: Returns ErrForbidden
//   - Empty label after trim: Returns ErrValidation
//   - CHECKBOX/RADIO/DROPDOWN without options: Returns ErrValidation
//   - TEXT fields without options: OK (options not required)
//   - First field (no existing fields): fieldOrder = 0
//   - Database error finding form: Returns wrapped error
//   - Database error checking admin: Returns wrapped error
//   - Database error getting fields: Returns wrapped error
//   - Database error adding field: Returns wrapped error
func (s *formService) AddField(ctx context.Context, formID string, userID string, req *models.AddFormFieldRequest) (*models.FormField, error) {
	// Find form by ID
	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to find form: %w", err)
	}
	if form == nil {
		return nil, ErrFormNotFound
	}

	// Check if user is org admin
	isAdmin, err := s.orgRepo.IsAdmin(ctx, form.OrganizationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org admin status: %w", err)
	}
	if !isAdmin {
		return nil, ErrForbidden
	}

	// Trim and validate label
	label := strings.TrimSpace(req.Label)
	if label == "" {
		return nil, ErrValidation
	}

	// Parse field type
	fieldType := models.FormFieldType(req.FieldType)

	// Validate options for field types that require them
	requiresOptions := fieldType == models.FormFieldTypeCheckbox ||
		fieldType == models.FormFieldTypeRadio ||
		fieldType == models.FormFieldTypeDropdown

	if requiresOptions && len(req.Options) == 0 {
		return nil, ErrValidation
	}

	// Get current fields to determine max fieldOrder
	existingFields, err := s.fieldRepo.GetFields(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing fields: %w", err)
	}

	// Calculate next field order
	maxOrder := -1
	for _, field := range existingFields {
		if field.FieldOrder > maxOrder {
			maxOrder = field.FieldOrder
		}
	}
	nextOrder := maxOrder + 1

	// Trim description if provided
	var description *string
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		description = &trimmed
	}

	// Clean options
	var options []string
	for _, opt := range req.Options {
		trimmed := strings.TrimSpace(opt)
		if trimmed != "" {
			options = append(options, trimmed)
		}
	}

	// Add field to database
	field, err := s.fieldRepo.AddField(ctx, formID, fieldType, label, description, req.IsRequired, nextOrder, options)
	if err != nil {
		return nil, fmt.Errorf("failed to add field: %w", err)
	}

	return field, nil
}

// ReorderFields reorders fields in a form
//
// Edge cases handled:
//   - Form not found: Returns ErrFormNotFound
//   - User not org admin: Returns ErrForbidden
//   - Empty fieldIDs: Calls repository (may be valid to clear order)
//   - Successful reorder: Returns nil
//   - Database error finding form: Returns wrapped error
//   - Database error checking admin: Returns wrapped error
//   - Database error reordering: Returns wrapped error
func (s *formService) ReorderFields(ctx context.Context, formID string, userID string, req *models.ReorderFieldsRequest) error {
	// Find form by ID
	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		return fmt.Errorf("failed to find form: %w", err)
	}
	if form == nil {
		return ErrFormNotFound
	}

	// Check if user is org admin
	isAdmin, err := s.orgRepo.IsAdmin(ctx, form.OrganizationID, userID)
	if err != nil {
		return fmt.Errorf("failed to check org admin status: %w", err)
	}
	if !isAdmin {
		return ErrForbidden
	}

	// Reorder fields
	if err := s.fieldRepo.ReorderFields(ctx, formID, req.FieldIDs); err != nil {
		return fmt.Errorf("failed to reorder fields: %w", err)
	}

	return nil
}

// UpdateField updates a form field
//
// Edge cases handled:
//   - Form not found: Returns ErrFormNotFound
//   - User not org admin: Returns ErrForbidden
//   - Empty label after trim: Returns ErrValidation
//   - Field not found: Repository returns error
//   - Successful update: Returns updated field
//   - Database error finding form: Returns wrapped error
//   - Database error checking admin: Returns wrapped error
//   - Database error updating field: Returns wrapped error
func (s *formService) UpdateField(ctx context.Context, formID string, fieldID string, userID string, req *models.UpdateFormFieldRequest) (*models.FormField, error) {
	// Find form by ID
	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to find form: %w", err)
	}
	if form == nil {
		return nil, ErrFormNotFound
	}

	// Check if user is org admin
	isAdmin, err := s.orgRepo.IsAdmin(ctx, form.OrganizationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org admin status: %w", err)
	}
	if !isAdmin {
		return nil, ErrForbidden
	}

	// Trim and validate label
	label := strings.TrimSpace(req.Label)
	if label == "" {
		return nil, ErrValidation
	}

	// Trim description if provided
	var description *string
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		description = &trimmed
	}

	// Clean options
	var options []string
	for _, opt := range req.Options {
		trimmed := strings.TrimSpace(opt)
		if trimmed != "" {
			options = append(options, trimmed)
		}
	}

	// Update field
	field, err := s.fieldRepo.UpdateField(ctx, fieldID, label, description, req.IsRequired, options)
	if err != nil {
		return nil, fmt.Errorf("failed to update field: %w", err)
	}

	return field, nil
}

// DeleteField deletes a form field
//
// Edge cases handled:
//   - Form not found: Returns ErrFormNotFound
//   - User not org admin: Returns ErrForbidden
//   - Field not found: Repository handles (may be idempotent)
//   - Successful delete: Returns nil
//   - Database error finding form: Returns wrapped error
//   - Database error checking admin: Returns wrapped error
//   - Database error deleting field: Returns wrapped error
func (s *formService) DeleteField(ctx context.Context, formID string, fieldID string, userID string) error {
	// Find form by ID
	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		return fmt.Errorf("failed to find form: %w", err)
	}
	if form == nil {
		return ErrFormNotFound
	}

	// Check if user is org admin
	isAdmin, err := s.orgRepo.IsAdmin(ctx, form.OrganizationID, userID)
	if err != nil {
		return fmt.Errorf("failed to check org admin status: %w", err)
	}
	if !isAdmin {
		return ErrForbidden
	}

	// Delete field
	if err := s.fieldRepo.DeleteField(ctx, fieldID); err != nil {
		return fmt.Errorf("failed to delete field: %w", err)
	}

	return nil
}

// SubmitAnswer submits answers to a form
//
// Edge cases handled:
//   - Form not found: Returns ErrFormNotFound
//   - Form is not active: Returns ErrValidation
//   - Form is active: Submits answer successfully
//   - Empty answers: Repository handles validation
//   - Database error finding form: Returns wrapped error
//   - Database error submitting: Returns wrapped error
func (s *formService) SubmitAnswer(ctx context.Context, formID string, userID string, req *models.SubmitFormAnswerRequest) (*models.FormAnswer, error) {
	// Find form by ID
	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to find form: %w", err)
	}
	if form == nil {
		return nil, ErrFormNotFound
	}

	// Check if form is active
	if !form.IsActive {
		return nil, fmt.Errorf("%w: form is not accepting submissions", ErrValidation)
	}

	// Submit answer via repository
	answer, err := s.answerRepo.Submit(ctx, formID, userID, req.Answers)
	if err != nil {
		return nil, fmt.Errorf("failed to submit answer: %w", err)
	}

	return answer, nil
}

// ListAnswers returns all answers for a form
//
// Edge cases handled:
//   - Form not found: Returns ErrFormNotFound
//   - User not org admin: Returns ErrForbidden
//   - Form has no answers: Returns empty slice
//   - Form has answers: Returns all current answers
//   - Database error finding form: Returns wrapped error
//   - Database error checking admin: Returns wrapped error
//   - Database error fetching answers: Returns wrapped error
func (s *formService) ListAnswers(ctx context.Context, formID string, userID string) ([]*models.FormAnswer, error) {
	// Find form by ID
	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to find form: %w", err)
	}
	if form == nil {
		return nil, ErrFormNotFound
	}

	// Check if user is org admin
	isAdmin, err := s.orgRepo.IsAdmin(ctx, form.OrganizationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org admin status: %w", err)
	}
	if !isAdmin {
		return nil, ErrForbidden
	}

	// Return all current answers
	answers, err := s.answerRepo.FindByForm(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to find answers: %w", err)
	}

	return answers, nil
}

// GetMyAnswer returns the current user's latest answer for a form
//
// Edge cases handled:
//   - Form not found: Returns ErrFormNotFound
//   - User has no answer: Returns nil, nil (not an error)
//   - User has answer: Returns the current answer
//   - Database error finding form: Returns wrapped error
//   - Database error finding answer: Returns wrapped error
func (s *formService) GetMyAnswer(ctx context.Context, formID string, userID string) (*models.FormAnswer, error) {
	// Find form by ID
	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to find form: %w", err)
	}
	if form == nil {
		return nil, ErrFormNotFound
	}

	// Get user's current answer
	answer, err := s.answerRepo.FindCurrent(ctx, formID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find current answer: %w", err)
	}

	// If no answer exists, return nil (not an error)
	return answer, nil
}

// GetAnswerHistory returns the answer history for a specific user on a form
//
// Edge cases handled:
//   - Form not found: Returns ErrFormNotFound
//   - User viewing own history: Returns history
//   - Org admin viewing other's history: Returns history
//   - Non-admin viewing other's history: Returns ErrForbidden
//   - User has no history: Returns empty slice
//   - Database error finding form: Returns wrapped error
//   - Database error checking admin: Returns wrapped error
//   - Database error fetching history: Returns wrapped error
func (s *formService) GetAnswerHistory(ctx context.Context, formID string, targetUserID string, requesterID string) ([]*models.FormAnswer, error) {
	// Find form by ID
	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to find form: %w", err)
	}
	if form == nil {
		return nil, ErrFormNotFound
	}

	// Check authorization
	// Users can see their own history
	if requesterID != targetUserID {
		// Org admins can see any user's history
		isAdmin, err := s.orgRepo.IsAdmin(ctx, form.OrganizationID, requesterID)
		if err != nil {
			return nil, fmt.Errorf("failed to check org admin status: %w", err)
		}
		if !isAdmin {
			return nil, ErrForbidden
		}
	}

	// Return all answers via repository
	answers, err := s.answerRepo.FindHistory(ctx, formID, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find answer history: %w", err)
	}

	// Return empty slice if no history
	if answers == nil {
		return []*models.FormAnswer{}, nil
	}

	return answers, nil
}
