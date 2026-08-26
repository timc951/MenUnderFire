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

// GroupMessageHandler handles HTTP requests for group message endpoints
type GroupMessageHandler struct {
	messageService services.GroupMessageService
}

// NewGroupMessageHandler creates a new GroupMessageHandler with the given service
func NewGroupMessageHandler(messageService services.GroupMessageService) *GroupMessageHandler {
	return &GroupMessageHandler{
		messageService: messageService,
	}
}

// List handles GET /api/groups/{groupId}/messages
// Returns all messages for a group
func (h *GroupMessageHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	groupID := vars["groupId"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	messages, err := h.messageService.List(ctx, groupID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("group_id", groupID).Msg("Error listing messages for group")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You are not a member of this group")
			return
		}
		respondInternalError(w, "Failed to list messages")
		return
	}

	responses := make([]*models.GroupMessageResponse, len(messages))
	for i, msg := range messages {
		responses[i] = msg.ToResponse()
	}

	respondOK(w, responses)
}

// Send handles POST /api/groups/{groupId}/messages
// Sends a message to all group members
func (h *GroupMessageHandler) Send(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	groupID := vars["groupId"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.SendGroupMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateSendGroupMessageRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	message, err := h.messageService.Send(ctx, groupID, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Str("group_id", groupID).Msg("Error sending message to group")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to send messages to this group")
			return
		}
		respondInternalError(w, "Failed to send message")
		return
	}

	respondCreated(w, message.ToResponse())
}

// Delete handles DELETE /api/groups/{groupId}/messages/{messageId}
// Deletes a message from a group
func (h *GroupMessageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	groupID := vars["groupId"]
	messageID := vars["messageId"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	err = h.messageService.Delete(ctx, groupID, messageID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("message_id", messageID).Str("group_id", groupID).Msg("Error deleting message from group")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			respondNotFound(w, "Message not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to delete this message")
			return
		}
		respondInternalError(w, "Failed to delete message")
		return
	}

	respondNoContent(w)
}

// ListGroupForms handles GET /api/groups/{groupId}/form-reports
// Returns all form messages sent to a group
func (h *GroupMessageHandler) ListGroupForms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	groupID := vars["groupId"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	messages, err := h.messageService.ListGroupForms(ctx, groupID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("group_id", groupID).Msg("Error listing group forms")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to view form reports")
			return
		}
		respondInternalError(w, "Failed to list group forms")
		return
	}

	responses := make([]*models.GroupMessageResponse, len(messages))
	for i, msg := range messages {
		responses[i] = msg.ToResponse()
	}

	respondOK(w, responses)
}

// ListGroupFormAnswers handles GET /api/groups/{groupId}/form-reports/{formId}
// Returns form answers for a group, filtered to group members
func (h *GroupMessageHandler) ListGroupFormAnswers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	groupID := vars["groupId"]
	formID := vars["formId"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	answers, err := h.messageService.ListGroupFormAnswers(ctx, groupID, formID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("group_id", groupID).Str("form_id", formID).Msg("Error listing group form answers")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to view form answers")
			return
		}
		respondInternalError(w, "Failed to list form answers")
		return
	}

	responses := make([]*models.FormAnswerResponse, len(answers))
	for i, a := range answers {
		responses[i] = a.ToResponse()
	}

	respondOK(w, responses)
}

// validateSendGroupMessageRequest validates the SendGroupMessageRequest
func validateSendGroupMessageRequest(req *models.SendGroupMessageRequest) error {
	if req.Content == "" {
		return errors.New("content is required")
	}
	return nil
}
