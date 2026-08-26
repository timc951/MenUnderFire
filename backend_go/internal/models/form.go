package models

import (
	"encoding/json"
	"time"
)

// FormFieldType represents the type of a form field
type FormFieldType string

const (
	FormFieldTypeTextDisplay FormFieldType = "TEXT_DISPLAY"
	FormFieldTypeTextSmall   FormFieldType = "TEXT_SMALL"
	FormFieldTypeTextMedium  FormFieldType = "TEXT_MEDIUM"
	FormFieldTypeTextLarge   FormFieldType = "TEXT_LARGE"
	FormFieldTypeCheckbox    FormFieldType = "CHECKBOX"
	FormFieldTypeRadio       FormFieldType = "RADIO"
	FormFieldTypeDropdown    FormFieldType = "DROPDOWN"
)

// Form represents a form entity
type Form struct {
	ID             string      `json:"id"`
	OrganizationID string      `json:"organizationId"`
	Name           string      `json:"name"`
	Description    *string     `json:"description"`
	IsActive       bool        `json:"isActive"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
	Fields         []FormField `json:"fields,omitempty"`
}

// FormField represents a field in a form
type FormField struct {
	ID          string        `json:"id"`
	FormID      string        `json:"formId"`
	FieldType   FormFieldType `json:"fieldType"`
	Label       string        `json:"label"`
	Description *string       `json:"description"`
	IsRequired  bool          `json:"isRequired"`
	FieldOrder  int           `json:"fieldOrder"`
	Options     []string      `json:"options,omitempty"`
}

// FormAnswer represents a user's answer to a form
type FormAnswer struct {
	ID          string          `json:"id"`
	FormID      string          `json:"formId"`
	UserID      string          `json:"userId"`
	Version     int             `json:"version"`
	IsCurrent   bool            `json:"isCurrent"`
	Answers     json.RawMessage `json:"answers"`
	SubmittedAt time.Time       `json:"submittedAt"`
	User        *User           `json:"user,omitempty"`
}

// ===== Request DTOs =====

// CreateFormRequest is the request for creating a form
type CreateFormRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty"`
}

// UpdateFormRequest is the request for updating a form
type UpdateFormRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"isActive"`
}

// AddFormFieldRequest is the request for adding a field to a form
type AddFormFieldRequest struct {
	FieldType   string   `json:"fieldType" validate:"required"`
	Label       string   `json:"label" validate:"required,min=1,max=200"`
	Description *string  `json:"description,omitempty"`
	IsRequired  bool     `json:"isRequired"`
	Options     []string `json:"options,omitempty"`
}

// UpdateFormFieldRequest is the request for updating a form field
type UpdateFormFieldRequest struct {
	Label       string   `json:"label" validate:"required,min=1,max=200"`
	Description *string  `json:"description,omitempty"`
	IsRequired  bool     `json:"isRequired"`
	Options     []string `json:"options,omitempty"`
}

// ReorderFieldsRequest is the request for reordering form fields
type ReorderFieldsRequest struct {
	FieldIDs []string `json:"fieldIds" validate:"required"`
}

// SubmitFormAnswerRequest is the request for submitting form answers
type SubmitFormAnswerRequest struct {
	Answers json.RawMessage `json:"answers" validate:"required"`
}

// ===== Response DTOs =====

// FormResponse is a summary of a form for list views
type FormResponse struct {
	ID             string  `json:"id"`
	OrganizationID string  `json:"organizationId"`
	Name           string  `json:"name"`
	Description    *string `json:"description,omitempty"`
	IsActive       bool    `json:"isActive"`
	FieldCount     int     `json:"fieldCount"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

// FormDetailResponse is detailed form info with all fields
type FormDetailResponse struct {
	ID             string              `json:"id"`
	OrganizationID string              `json:"organizationId"`
	Name           string              `json:"name"`
	Description    *string             `json:"description,omitempty"`
	IsActive       bool                `json:"isActive"`
	Fields         []FormFieldResponse `json:"fields"`
	CreatedAt      string              `json:"createdAt"`
	UpdatedAt      string              `json:"updatedAt"`
}

// FormFieldResponse is the response for a form field
type FormFieldResponse struct {
	ID          string   `json:"id"`
	FieldType   string   `json:"fieldType"`
	Label       string   `json:"label"`
	Description *string  `json:"description,omitempty"`
	IsRequired  bool     `json:"isRequired"`
	FieldOrder  int      `json:"fieldOrder"`
	Options     []string `json:"options,omitempty"`
}

// FormAnswerResponse is the response for a submitted form answer
type FormAnswerResponse struct {
	ID          string          `json:"id"`
	FormID      string          `json:"formId"`
	UserID      string          `json:"userId"`
	UserName    *string         `json:"userName,omitempty"`
	Version     int             `json:"version"`
	IsCurrent   bool            `json:"isCurrent"`
	Answers     json.RawMessage `json:"answers"`
	SubmittedAt string          `json:"submittedAt"`
}

// ===== Conversion Methods =====

// ToResponse converts a Form entity to a FormResponse DTO
func (f *Form) ToResponse() *FormResponse {
	return &FormResponse{
		ID:             f.ID,
		OrganizationID: f.OrganizationID,
		Name:           f.Name,
		Description:    f.Description,
		IsActive:       f.IsActive,
		FieldCount:     len(f.Fields),
		CreatedAt:      f.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      f.UpdatedAt.Format(time.RFC3339),
	}
}

// ToDetailResponse converts a Form entity to a FormDetailResponse DTO
func (f *Form) ToDetailResponse() *FormDetailResponse {
	fields := make([]FormFieldResponse, len(f.Fields))
	for i, field := range f.Fields {
		fields[i] = *field.ToResponse()
	}
	return &FormDetailResponse{
		ID:             f.ID,
		OrganizationID: f.OrganizationID,
		Name:           f.Name,
		Description:    f.Description,
		IsActive:       f.IsActive,
		Fields:         fields,
		CreatedAt:      f.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      f.UpdatedAt.Format(time.RFC3339),
	}
}

// ToResponse converts a FormField entity to a FormFieldResponse DTO
func (f *FormField) ToResponse() *FormFieldResponse {
	return &FormFieldResponse{
		ID:          f.ID,
		FieldType:   string(f.FieldType),
		Label:       f.Label,
		Description: f.Description,
		IsRequired:  f.IsRequired,
		FieldOrder:  f.FieldOrder,
		Options:     f.Options,
	}
}

// ToResponse converts a FormAnswer entity to a FormAnswerResponse DTO
func (a *FormAnswer) ToResponse() *FormAnswerResponse {
	var userName *string
	if a.User != nil {
		userName = &a.User.DisplayName
	}
	return &FormAnswerResponse{
		ID:          a.ID,
		FormID:      a.FormID,
		UserID:      a.UserID,
		UserName:    userName,
		Version:     a.Version,
		IsCurrent:   a.IsCurrent,
		Answers:     a.Answers,
		SubmittedAt: a.SubmittedAt.Format(time.RFC3339),
	}
}
