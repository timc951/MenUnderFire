package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"menunderfire/internal/logger"
	"menunderfire/internal/models"
	"menunderfire/internal/services"

	"github.com/gorilla/mux"
)

// FormHandler handles HTTP requests for form-related endpoints
type FormHandler struct {
	formService services.FormService
}

// NewFormHandler creates a new FormHandler with the given service
func NewFormHandler(formService services.FormService) *FormHandler {
	return &FormHandler{
		formService: formService,
	}
}

// ListByOrg handles GET /api/organizations/{orgId}/forms
// Returns all forms for an organization
func (h *FormHandler) ListByOrg(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	orgID := vars["orgId"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	forms, err := h.formService.ListByOrg(ctx, orgID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("org_id", orgID).Msg("Error listing forms")
		if errors.Is(err, services.ErrOrganizationNotFound) {
			respondNotFound(w, "Organization not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have access to this organization")
			return
		}
		respondInternalError(w, "Failed to list forms")
		return
	}

	responses := make([]*models.FormResponse, len(forms))
	for i, form := range forms {
		responses[i] = form.ToResponse()
	}

	respondOK(w, responses)
}

// Create handles POST /api/organizations/{orgId}/forms
// Creates a new form (Org Admin only)
func (h *FormHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	orgID := vars["orgId"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.CreateFormRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateCreateFormRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	form, err := h.formService.Create(ctx, orgID, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating form")
		if errors.Is(err, services.ErrOrganizationNotFound) {
			respondNotFound(w, "Organization not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Organization admin access required")
			return
		}
		respondInternalError(w, "Failed to create form")
		return
	}

	respondCreated(w, form.ToResponse())
}

// Get handles GET /api/forms/{id}
// Returns detailed form information with all fields
func (h *FormHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	formID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	form, err := h.formService.GetDetailByID(ctx, formID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("form_id", formID).Msg("Error getting form")
		if errors.Is(err, services.ErrFormNotFound) {
			respondNotFound(w, "Form not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have access to this form")
			return
		}
		respondInternalError(w, "Failed to get form")
		return
	}

	respondOK(w, form.ToDetailResponse())
}

// Update handles PUT /api/forms/{id}
// Updates a form
func (h *FormHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	formID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.UpdateFormRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateUpdateFormRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	form, err := h.formService.Update(ctx, formID, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Str("form_id", formID).Msg("Error updating form")
		if errors.Is(err, services.ErrFormNotFound) {
			respondNotFound(w, "Form not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to update this form")
			return
		}
		respondInternalError(w, "Failed to update form")
		return
	}

	respondOK(w, form.ToResponse())
}

// Delete handles DELETE /api/forms/{id}
// Deletes a form
func (h *FormHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	formID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	err = h.formService.Delete(ctx, formID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("form_id", formID).Msg("Error deleting form")
		if errors.Is(err, services.ErrFormNotFound) {
			respondNotFound(w, "Form not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to delete this form")
			return
		}
		respondInternalError(w, "Failed to delete form")
		return
	}

	respondNoContent(w)
}

// AddField handles POST /api/forms/{id}/fields
// Adds a field to a form
func (h *FormHandler) AddField(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	formID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.AddFormFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateAddFormFieldRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	field, err := h.formService.AddField(ctx, formID, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Str("form_id", formID).Msg("Error adding field to form")
		if errors.Is(err, services.ErrFormNotFound) {
			respondNotFound(w, "Form not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to modify this form")
			return
		}
		if errors.Is(err, services.ErrValidation) {
			respondBadRequest(w, err.Error())
			return
		}
		respondInternalError(w, "Failed to add field")
		return
	}

	respondCreated(w, field.ToResponse())
}

// ReorderFields handles PUT /api/forms/{id}/fields/reorder
// Reorders fields in a form
func (h *FormHandler) ReorderFields(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	formID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.ReorderFieldsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if len(req.FieldIDs) == 0 {
		respondBadRequest(w, "fieldIds is required")
		return
	}

	err = h.formService.ReorderFields(ctx, formID, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Str("form_id", formID).Msg("Error reordering fields for form")
		if errors.Is(err, services.ErrFormNotFound) {
			respondNotFound(w, "Form not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to modify this form")
			return
		}
		if errors.Is(err, services.ErrValidation) {
			respondBadRequest(w, err.Error())
			return
		}
		respondInternalError(w, "Failed to reorder fields")
		return
	}

	respondOK(w, MessageResponse{Message: "Fields reordered successfully"})
}

// UpdateField handles PUT /api/forms/{formId}/fields/{fieldId}
// Updates a form field
func (h *FormHandler) UpdateField(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	formID := vars["formId"]
	fieldID := vars["fieldId"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.UpdateFormFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateUpdateFormFieldRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	field, err := h.formService.UpdateField(ctx, formID, fieldID, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Str("field_id", fieldID).Str("form_id", formID).Msg("Error updating field")
		if errors.Is(err, services.ErrFormNotFound) {
			respondNotFound(w, "Form not found")
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			respondNotFound(w, "Field not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to modify this form")
			return
		}
		respondInternalError(w, "Failed to update field")
		return
	}

	respondOK(w, field.ToResponse())
}

// DeleteField handles DELETE /api/forms/{formId}/fields/{fieldId}
// Deletes a form field
func (h *FormHandler) DeleteField(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	formID := vars["formId"]
	fieldID := vars["fieldId"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	err = h.formService.DeleteField(ctx, formID, fieldID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("field_id", fieldID).Str("form_id", formID).Msg("Error deleting field")
		if errors.Is(err, services.ErrFormNotFound) {
			respondNotFound(w, "Form not found")
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			respondNotFound(w, "Field not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to modify this form")
			return
		}
		respondInternalError(w, "Failed to delete field")
		return
	}

	respondNoContent(w)
}

// SubmitAnswer handles POST /api/forms/{id}/answers
// Submits answers to a form
func (h *FormHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	formID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.SubmitFormAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if len(req.Answers) == 0 {
		respondBadRequest(w, "answers is required")
		return
	}

	answer, err := h.formService.SubmitAnswer(ctx, formID, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Str("form_id", formID).Msg("Error submitting answer to form")
		if errors.Is(err, services.ErrFormNotFound) {
			respondNotFound(w, "Form not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have access to this form")
			return
		}
		if errors.Is(err, services.ErrValidation) {
			respondBadRequest(w, err.Error())
			return
		}
		respondInternalError(w, "Failed to submit answer")
		return
	}

	respondCreated(w, answer.ToResponse())
}

// ListAnswers handles GET /api/forms/{id}/answers
// Returns all answers for a form
func (h *FormHandler) ListAnswers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	formID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	answers, err := h.formService.ListAnswers(ctx, formID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("form_id", formID).Msg("Error listing answers for form")
		if errors.Is(err, services.ErrFormNotFound) {
			respondNotFound(w, "Form not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to view answers")
			return
		}
		respondInternalError(w, "Failed to list answers")
		return
	}

	responses := make([]*models.FormAnswerResponse, len(answers))
	for i, answer := range answers {
		responses[i] = answer.ToResponse()
	}

	respondOK(w, responses)
}

// GetMyAnswer handles GET /api/forms/{id}/answers/me
// Returns the current user's latest answer for a form
func (h *FormHandler) GetMyAnswer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	formID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	answer, err := h.formService.GetMyAnswer(ctx, formID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("form_id", formID).Msg("Error getting my answer for form")
		if errors.Is(err, services.ErrFormNotFound) {
			respondNotFound(w, "Form not found")
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			respondNotFound(w, "No answer found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have access to this form")
			return
		}
		respondInternalError(w, "Failed to get answer")
		return
	}

	if answer == nil {
		respondNotFound(w, "No answer found")
		return
	}

	respondOK(w, answer.ToResponse())
}

// GetAnswerHistory handles GET /api/forms/{id}/answers/history/{userId}
// Returns the answer history for a specific user on a form
func (h *FormHandler) GetAnswerHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	formID := vars["id"]
	targetUserID := vars["userId"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	answers, err := h.formService.GetAnswerHistory(ctx, formID, targetUserID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", targetUserID).Str("form_id", formID).Msg("Error getting answer history")
		if errors.Is(err, services.ErrFormNotFound) {
			respondNotFound(w, "Form not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to view this user's answers")
			return
		}
		respondInternalError(w, "Failed to get answer history")
		return
	}

	responses := make([]*models.FormAnswerResponse, len(answers))
	for i, answer := range answers {
		responses[i] = answer.ToResponse()
	}

	respondOK(w, responses)
}

// Validation helpers

func validateCreateFormRequest(req *models.CreateFormRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Name) > 100 {
		return errors.New("name must be 100 characters or less")
	}
	return nil
}

func validateUpdateFormRequest(req *models.UpdateFormRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Name) > 100 {
		return errors.New("name must be 100 characters or less")
	}
	return nil
}

func validateAddFormFieldRequest(req *models.AddFormFieldRequest) error {
	if req.FieldType == "" {
		return errors.New("fieldType is required")
	}
	validTypes := map[string]bool{
		"TEXT_DISPLAY": true,
		"TEXT_SMALL":   true,
		"TEXT_MEDIUM":  true,
		"TEXT_LARGE":   true,
		"CHECKBOX":     true,
		"RADIO":        true,
		"DROPDOWN":     true,
	}
	if !validTypes[req.FieldType] {
		return errors.New("invalid fieldType")
	}
	if req.Label == "" {
		return errors.New("label is required")
	}
	if len(req.Label) > 200 {
		return errors.New("label must be 200 characters or less")
	}
	// Options required for CHECKBOX, RADIO, DROPDOWN
	if (req.FieldType == "CHECKBOX" || req.FieldType == "RADIO" || req.FieldType == "DROPDOWN") && len(req.Options) == 0 {
		return errors.New("options are required for this field type")
	}
	return nil
}

func validateUpdateFormFieldRequest(req *models.UpdateFormFieldRequest) error {
	if req.Label == "" {
		return errors.New("label is required")
	}
	if len(req.Label) > 200 {
		return errors.New("label must be 200 characters or less")
	}
	return nil
}
