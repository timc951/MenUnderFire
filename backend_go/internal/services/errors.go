package services

import "errors"

// Common service errors
var (
	// ErrNotFound is returned when a requested resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrGroupNotFound is returned when a group is not found
	ErrGroupNotFound = errors.New("group not found")

	// ErrOrganizationNotFound is returned when an organization is not found
	ErrOrganizationNotFound = errors.New("organization not found")

	// ErrInvitationNotFound is returned when an invitation is not found
	ErrInvitationNotFound = errors.New("invitation not found")

	// ErrFormNotFound is returned when a form is not found
	ErrFormNotFound = errors.New("form not found")

	// ErrPageNotFound is returned when a site page is not found
	ErrPageNotFound = errors.New("page not found")

	// ErrDraftNotFound is returned when a page draft is not found or expired
	ErrDraftNotFound = errors.New("draft not found or expired")

	// ErrMessageNotFound is returned when a message is not found
	ErrMessageNotFound = errors.New("message not found")

	// ErrReportNotFound is returned when a report is not found
	ErrReportNotFound = errors.New("report not found")

	// ErrValidation is returned when validation fails
	ErrValidation = errors.New("validation error")

	// ErrUnauthorized is returned when the user is not authenticated
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when the user lacks permission for an action
	ErrForbidden = errors.New("forbidden")

	// ErrConflict is returned when there is a conflict (e.g., duplicate entry)
	ErrConflict = errors.New("conflict")

	// ErrInvalidToken is returned when a token is invalid or expired
	ErrInvalidToken = errors.New("invalid or expired token")

	// ErrInvalidInviteCode is returned when an invite code is invalid
	ErrInvalidInviteCode = errors.New("invalid invite code")
)
