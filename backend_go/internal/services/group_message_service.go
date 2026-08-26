package services

import (
	"context"
	"fmt"
	"strings"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// GroupMessageService defines the interface for group message business logic
type GroupMessageService interface {
	// List returns all messages for a group
	List(ctx context.Context, groupID string, userID string) ([]*models.GroupMessage, error)

	// Send sends a message to a group
	Send(ctx context.Context, groupID string, userID string, req *models.SendGroupMessageRequest) (*models.GroupMessage, error)

	// Delete deletes a message from a group
	Delete(ctx context.Context, groupID string, messageID string, userID string) error

	// GetByID retrieves a message by its ID
	GetByID(ctx context.Context, messageID string) (*models.GroupMessage, error)

	// ListGroupForms returns all form messages sent to a group
	ListGroupForms(ctx context.Context, groupID string, userID string) ([]*models.GroupMessage, error)

	// ListGroupFormAnswers returns form answers for a group, filtered to group members
	ListGroupFormAnswers(ctx context.Context, groupID string, formID string, userID string) ([]*models.FormAnswer, error)
}

type groupMessageService struct {
	messageRepo    repositories.GroupMessageRepository
	groupRepo      repositories.GroupRepository
	groupService   GroupService
	orgRepo        repositories.OrganizationRepository
	userRepo       repositories.UserRepository
	formRepo       repositories.FormRepository
	formAnswerRepo repositories.FormAnswerRepository
}

// NewGroupMessageService creates a new GroupMessageService implementation with the provided dependencies
func NewGroupMessageService(
	messageRepo repositories.GroupMessageRepository,
	groupRepo repositories.GroupRepository,
	groupService GroupService,
	orgRepo repositories.OrganizationRepository,
	userRepo repositories.UserRepository,
	formRepo repositories.FormRepository,
	formAnswerRepo repositories.FormAnswerRepository,
) GroupMessageService {
	return &groupMessageService{
		messageRepo:    messageRepo,
		groupRepo:      groupRepo,
		groupService:   groupService,
		orgRepo:        orgRepo,
		userRepo:       userRepo,
		formRepo:       formRepo,
		formAnswerRepo: formAnswerRepo,
	}
}

// List returns all messages for a group
func (s *groupMessageService) List(ctx context.Context, groupID string, userID string) ([]*models.GroupMessage, error) {
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to find group: %w", err)
	}
	if group == nil {
		return nil, ErrGroupNotFound
	}

	role, err := s.groupService.GetUserRole(ctx, groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	if role != "" {
		messages, err := s.messageRepo.FindByGroupID(ctx, groupID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch messages: %w", err)
		}
		return messages, nil
	}

	siteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check site admin status: %w", err)
	}
	if siteAdmin {
		messages, err := s.messageRepo.FindByGroupID(ctx, groupID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch messages: %w", err)
		}
		return messages, nil
	}

	orgAdmin, err := s.orgRepo.IsAdmin(ctx, group.OrganizationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org admin status: %w", err)
	}
	if orgAdmin {
		messages, err := s.messageRepo.FindByGroupID(ctx, groupID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch messages: %w", err)
		}
		return messages, nil
	}

	return nil, ErrForbidden
}

// Send sends a message to a group
func (s *groupMessageService) Send(ctx context.Context, groupID string, userID string, req *models.SendGroupMessageRequest) (*models.GroupMessage, error) {
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to find group: %w", err)
	}
	if group == nil {
		return nil, ErrGroupNotFound
	}

	role, err := s.groupService.GetUserRole(ctx, groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}
	if role == "" {
		return nil, ErrForbidden
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, ErrValidation
	}

	message, err := s.messageRepo.Create(ctx, groupID, userID, content, req.NotifyMembers, req.FormID)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return message, nil
}

// Delete deletes a message from a group
func (s *groupMessageService) Delete(ctx context.Context, groupID string, messageID string, userID string) error {
	message, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to find message: %w", err)
	}
	if message == nil {
		return ErrMessageNotFound
	}

	if message.GroupID != groupID {
		return ErrMessageNotFound
	}

	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to find group: %w", err)
	}
	if group == nil {
		return ErrGroupNotFound
	}

	if message.SenderID == userID {
		if err := s.messageRepo.Delete(ctx, messageID); err != nil {
			return fmt.Errorf("failed to delete message: %w", err)
		}
		return nil
	}

	role, err := s.groupService.GetUserRole(ctx, groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}
	if role == "OWNER" || role == "LEADER" {
		if err := s.messageRepo.Delete(ctx, messageID); err != nil {
			return fmt.Errorf("failed to delete message: %w", err)
		}
		return nil
	}

	siteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check site admin status: %w", err)
	}
	if siteAdmin {
		if err := s.messageRepo.Delete(ctx, messageID); err != nil {
			return fmt.Errorf("failed to delete message: %w", err)
		}
		return nil
	}

	orgAdmin, err := s.orgRepo.IsAdmin(ctx, group.OrganizationID, userID)
	if err != nil {
		return fmt.Errorf("failed to check org admin status: %w", err)
	}
	if orgAdmin {
		if err := s.messageRepo.Delete(ctx, messageID); err != nil {
			return fmt.Errorf("failed to delete message: %w", err)
		}
		return nil
	}

	return ErrForbidden
}

// GetByID retrieves a message by its ID
func (s *groupMessageService) GetByID(ctx context.Context, messageID string) (*models.GroupMessage, error) {
	if messageID == "" {
		return nil, ErrMessageNotFound
	}

	message, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message by ID: %w", err)
	}
	if message == nil {
		return nil, ErrMessageNotFound
	}

	return message, nil
}

// ListGroupForms returns all form messages sent to a group (for the Reports tab)
func (s *groupMessageService) ListGroupForms(ctx context.Context, groupID string, userID string) ([]*models.GroupMessage, error) {
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to find group: %w", err)
	}
	if group == nil {
		return nil, ErrGroupNotFound
	}

	// Check if user is a leader/owner
	role, err := s.groupService.GetUserRole(ctx, groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}
	if role == "OWNER" || role == "LEADER" {
		return s.messageRepo.FindFormMessagesByGroupID(ctx, groupID)
	}

	// Check site admin
	siteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check site admin status: %w", err)
	}
	if siteAdmin {
		return s.messageRepo.FindFormMessagesByGroupID(ctx, groupID)
	}

	// Check org admin
	orgAdmin, err := s.orgRepo.IsAdmin(ctx, group.OrganizationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org admin status: %w", err)
	}
	if orgAdmin {
		return s.messageRepo.FindFormMessagesByGroupID(ctx, groupID)
	}

	return nil, ErrForbidden
}

// ListGroupFormAnswers returns form answers for a group, filtered to group members
func (s *groupMessageService) ListGroupFormAnswers(ctx context.Context, groupID string, formID string, userID string) ([]*models.FormAnswer, error) {
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to find group: %w", err)
	}
	if group == nil {
		return nil, ErrGroupNotFound
	}

	// Check if user is a leader/owner
	role, err := s.groupService.GetUserRole(ctx, groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	hasAccess := role == "OWNER" || role == "LEADER"

	if !hasAccess {
		siteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check site admin status: %w", err)
		}
		hasAccess = siteAdmin
	}

	if !hasAccess {
		orgAdmin, err := s.orgRepo.IsAdmin(ctx, group.OrganizationID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check org admin status: %w", err)
		}
		hasAccess = orgAdmin
	}

	if !hasAccess {
		return nil, ErrForbidden
	}

	// Get group members to filter answers
	members, err := s.groupRepo.FindMembers(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}

	memberIDs := make([]string, len(members))
	for i, m := range members {
		memberIDs[i] = m.UserID
	}

	if len(memberIDs) == 0 {
		return []*models.FormAnswer{}, nil
	}

	answers, err := s.formAnswerRepo.FindByFormAndUserIDs(ctx, formID, memberIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get form answers: %w", err)
	}

	return answers, nil
}
