package models

import "time"

// GroupMessage represents a message sent to a group
type GroupMessage struct {
	ID            string    `json:"id"`
	GroupID       string    `json:"groupId"`
	SenderID      string    `json:"senderId"`
	Content       string    `json:"content"`
	NotifyMembers bool      `json:"notifyMembers"`
	FormID        *string   `json:"formId,omitempty"`
	FormName      *string   `json:"formName,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	Sender        *User     `json:"sender,omitempty"`
}

// ===== Request DTOs =====

// SendGroupMessageRequest is the request body for sending a group message
type SendGroupMessageRequest struct {
	Content       string  `json:"content" validate:"required,min=1"`
	NotifyMembers bool    `json:"notifyMembers"`
	FormID        *string `json:"formId,omitempty"`
}

// ===== Response DTOs =====

// GroupMessageResponse is the response for a group message
type GroupMessageResponse struct {
	ID            string  `json:"id"`
	GroupID       string  `json:"groupId"`
	SenderID      string  `json:"senderId"`
	SenderName    string  `json:"senderName"`
	Content       string  `json:"content"`
	NotifyMembers bool    `json:"notifyMembers"`
	FormID        *string `json:"formId,omitempty"`
	FormName      *string `json:"formName,omitempty"`
	CreatedAt     string  `json:"createdAt"`
}

// ===== Conversion Methods =====

// ToResponse converts a GroupMessage entity to a GroupMessageResponse DTO
func (m *GroupMessage) ToResponse() *GroupMessageResponse {
	senderName := ""
	if m.Sender != nil {
		senderName = m.Sender.DisplayName
	}
	return &GroupMessageResponse{
		ID:            m.ID,
		GroupID:       m.GroupID,
		SenderID:      m.SenderID,
		SenderName:    senderName,
		Content:       m.Content,
		NotifyMembers: m.NotifyMembers,
		FormID:        m.FormID,
		FormName:      m.FormName,
		CreatedAt:     m.CreatedAt.Format(time.RFC3339),
	}
}
