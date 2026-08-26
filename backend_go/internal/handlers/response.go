package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error response body
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// respondJSON writes a JSON response with the given status code
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		// Status line is already committed; a failed body write is unrecoverable here.
		_ = json.NewEncoder(w).Encode(data)
	}
}

// respondError writes an error response with the given status code
func respondError(w http.ResponseWriter, status int, err string, message string) {
	respondJSON(w, status, ErrorResponse{
		Error:   err,
		Message: message,
	})
}

// respondBadRequest writes a 400 Bad Request response
func respondBadRequest(w http.ResponseWriter, message string) {
	respondError(w, http.StatusBadRequest, "bad_request", message)
}

// respondUnauthorized writes a 401 Unauthorized response
func respondUnauthorized(w http.ResponseWriter, message string) {
	respondError(w, http.StatusUnauthorized, "unauthorized", message)
}

// respondForbidden writes a 403 Forbidden response
func respondForbidden(w http.ResponseWriter, message string) {
	respondError(w, http.StatusForbidden, "forbidden", message)
}

// respondNotFound writes a 404 Not Found response
func respondNotFound(w http.ResponseWriter, message string) {
	respondError(w, http.StatusNotFound, "not_found", message)
}

// respondInternalError writes a 500 Internal Server Error response
func respondInternalError(w http.ResponseWriter, message string) {
	respondError(w, http.StatusInternalServerError, "internal_error", message)
}

// respondCreated writes a 201 Created response with data
func respondCreated(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusCreated, data)
}

// respondOK writes a 200 OK response with data
func respondOK(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusOK, data)
}

// respondNoContent writes a 204 No Content response
func respondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
