package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// Helper functions
func strPtr(s string) *string {
	return &s
}

// Mock repositories for form service tests
// Note: Using distinct names to avoid conflicts with mocks in dashboard_service_impl_test.go

type mockFormRepo struct {
	findByIDFn             func(ctx context.Context, formID string) (*models.Form, error)
	findByOrganizationIDFn func(ctx context.Context, orgID string) ([]*models.Form, error)
	createFn               func(ctx context.Context, orgID, name string, description *string, createdBy string) (*models.Form, error)
	updateFn               func(ctx context.Context, formID, name string, description *string, isActive bool) (*models.Form, error)
	deleteFn               func(ctx context.Context, formID string) error
}

func (m *mockFormRepo) FindByID(ctx context.Context, formID string) (*models.Form, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, formID)
	}
	return nil, errors.New("findByIDFn not set")
}

func (m *mockFormRepo) FindByOrganizationID(ctx context.Context, orgID string) ([]*models.Form, error) {
	if m.findByOrganizationIDFn != nil {
		return m.findByOrganizationIDFn(ctx, orgID)
	}
	return nil, errors.New("findByOrganizationIDFn not set")
}

func (m *mockFormRepo) Create(ctx context.Context, orgID, name string, description *string, createdBy string) (*models.Form, error) {
	if m.createFn != nil {
		return m.createFn(ctx, orgID, name, description, createdBy)
	}
	return nil, errors.New("createFn not set")
}

func (m *mockFormRepo) Update(ctx context.Context, formID, name string, description *string, isActive bool) (*models.Form, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, formID, name, description, isActive)
	}
	return nil, errors.New("updateFn not set")
}

func (m *mockFormRepo) Delete(ctx context.Context, formID string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, formID)
	}
	return errors.New("deleteFn not set")
}

type mockFormFieldRepo struct {
	getFieldsFn     func(ctx context.Context, formID string) ([]models.FormField, error)
	addFieldFn      func(ctx context.Context, formID string, fieldType models.FormFieldType, label string, description *string, isRequired bool, fieldOrder int, options []string) (*models.FormField, error)
	updateFieldFn   func(ctx context.Context, fieldID, label string, description *string, isRequired bool, options []string) (*models.FormField, error)
	deleteFieldFn   func(ctx context.Context, fieldID string) error
	reorderFieldsFn func(ctx context.Context, formID string, fieldIDs []string) error
}

func (m *mockFormFieldRepo) GetFields(ctx context.Context, formID string) ([]models.FormField, error) {
	if m.getFieldsFn != nil {
		return m.getFieldsFn(ctx, formID)
	}
	return nil, errors.New("getFieldsFn not set")
}

func (m *mockFormFieldRepo) AddField(ctx context.Context, formID string, fieldType models.FormFieldType, label string, description *string, isRequired bool, fieldOrder int, options []string) (*models.FormField, error) {
	if m.addFieldFn != nil {
		return m.addFieldFn(ctx, formID, fieldType, label, description, isRequired, fieldOrder, options)
	}
	return nil, errors.New("addFieldFn not set")
}

func (m *mockFormFieldRepo) UpdateField(ctx context.Context, fieldID, label string, description *string, isRequired bool, options []string) (*models.FormField, error) {
	if m.updateFieldFn != nil {
		return m.updateFieldFn(ctx, fieldID, label, description, isRequired, options)
	}
	return nil, errors.New("updateFieldFn not set")
}

func (m *mockFormFieldRepo) DeleteField(ctx context.Context, fieldID string) error {
	if m.deleteFieldFn != nil {
		return m.deleteFieldFn(ctx, fieldID)
	}
	return errors.New("deleteFieldFn not set")
}

func (m *mockFormFieldRepo) ReorderFields(ctx context.Context, formID string, fieldIDs []string) error {
	if m.reorderFieldsFn != nil {
		return m.reorderFieldsFn(ctx, formID, fieldIDs)
	}
	return errors.New("reorderFieldsFn not set")
}

type mockFormAnswerRepo struct {
	submitFn              func(ctx context.Context, formID, userID string, answers json.RawMessage) (*models.FormAnswer, error)
	findByFormFn          func(ctx context.Context, formID string) ([]*models.FormAnswer, error)
	findByFormAndUserIDFn func(ctx context.Context, formID string, userIDs []string) ([]*models.FormAnswer, error)
	findCurrentFn         func(ctx context.Context, formID, userID string) (*models.FormAnswer, error)
	findHistoryFn         func(ctx context.Context, formID, userID string) ([]*models.FormAnswer, error)
}

func (m *mockFormAnswerRepo) Submit(ctx context.Context, formID, userID string, answers json.RawMessage) (*models.FormAnswer, error) {
	if m.submitFn != nil {
		return m.submitFn(ctx, formID, userID, answers)
	}
	return nil, errors.New("submitFn not set")
}

func (m *mockFormAnswerRepo) FindByForm(ctx context.Context, formID string) ([]*models.FormAnswer, error) {
	if m.findByFormFn != nil {
		return m.findByFormFn(ctx, formID)
	}
	return nil, errors.New("findByFormFn not set")
}

func (m *mockFormAnswerRepo) FindByFormAndUserIDs(ctx context.Context, formID string, userIDs []string) ([]*models.FormAnswer, error) {
	if m.findByFormAndUserIDFn != nil {
		return m.findByFormAndUserIDFn(ctx, formID, userIDs)
	}
	return nil, errors.New("findByFormAndUserIDFn not set")
}

func (m *mockFormAnswerRepo) FindCurrent(ctx context.Context, formID, userID string) (*models.FormAnswer, error) {
	if m.findCurrentFn != nil {
		return m.findCurrentFn(ctx, formID, userID)
	}
	return nil, errors.New("findCurrentFn not set")
}

func (m *mockFormAnswerRepo) FindHistory(ctx context.Context, formID, userID string) ([]*models.FormAnswer, error) {
	if m.findHistoryFn != nil {
		return m.findHistoryFn(ctx, formID, userID)
	}
	return nil, errors.New("findHistoryFn not set")
}

type mockOrgRepoForForm struct {
	isAdminFn  func(ctx context.Context, orgID, userID string) (bool, error)
	isMemberFn func(ctx context.Context, orgID, userID string) (bool, error)
}

func (m *mockOrgRepoForForm) IsAdmin(ctx context.Context, orgID, userID string) (bool, error) {
	if m.isAdminFn != nil {
		return m.isAdminFn(ctx, orgID, userID)
	}
	return false, errors.New("isAdminFn not set")
}

func (m *mockOrgRepoForForm) FindByID(ctx context.Context, orgID string) (*models.Organization, error) {
	return nil, errors.New("FindByID not implemented")
}

func (m *mockOrgRepoForForm) FindByUserID(ctx context.Context, userID string) ([]*models.Organization, error) {
	return nil, errors.New("FindByUserID not implemented")
}

func (m *mockOrgRepoForForm) FindAll(ctx context.Context) ([]*models.Organization, error) {
	return nil, errors.New("FindAll not implemented")
}

func (m *mockOrgRepoForForm) Create(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error) {
	return nil, errors.New("Create not implemented")
}

func (m *mockOrgRepoForForm) Update(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error) {
	return nil, errors.New("Update not implemented")
}

func (m *mockOrgRepoForForm) FindAdmins(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error) {
	return nil, errors.New("FindAdmins not implemented")
}

func (m *mockOrgRepoForForm) FindAdmin(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
	return nil, errors.New("FindAdmin not implemented")
}

func (m *mockOrgRepoForForm) AddAdmin(ctx context.Context, orgID, userID string) error {
	return errors.New("AddAdmin not implemented")
}

func (m *mockOrgRepoForForm) IsMember(ctx context.Context, orgID, userID string) (bool, error) {
	if m.isMemberFn != nil {
		return m.isMemberFn(ctx, orgID, userID)
	}
	return false, errors.New("isMemberFn not set")
}

func (m *mockOrgRepoForForm) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("Count not implemented")
}

type mockUserRepoForForm struct {
	isSiteAdminFn func(ctx context.Context, userID string) (bool, error)
}

func (m *mockUserRepoForForm) IsSiteAdmin(ctx context.Context, userID string) (bool, error) {
	if m.isSiteAdminFn != nil {
		return m.isSiteAdminFn(ctx, userID)
	}
	return false, errors.New("isSiteAdminFn not set")
}

func (m *mockUserRepoForForm) FindByID(ctx context.Context, userID string) (*models.User, error) {
	return nil, errors.New("FindByID not implemented")
}

func (m *mockUserRepoForForm) FindByExternalID(ctx context.Context, externalID string) (*models.User, error) {
	return nil, errors.New("FindByExternalID not implemented")
}

func (m *mockUserRepoForForm) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, errors.New("FindByEmail not implemented")
}

func (m *mockUserRepoForForm) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("Count not implemented")
}

func (m *mockUserRepoForForm) Create(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
	return nil, errors.New("Create not implemented")
}

func (m *mockUserRepoForForm) CreateAsSiteAdmin(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
	return nil, errors.New("CreateAsSiteAdmin not implemented")
}

func (m *mockUserRepoForForm) Update(ctx context.Context, userID, displayName string) (*models.User, error) {
	return nil, errors.New("Update not implemented")
}

func (m *mockUserRepoForForm) UpdateExternalID(ctx context.Context, userID, externalID string) error {
	return errors.New("UpdateExternalID not implemented")
}

func (m *mockUserRepoForForm) UpdateInvitationInfo(ctx context.Context, userID, invitedByID, invitationID string) error {
	return errors.New("UpdateInvitationInfo not implemented")
}

func (m *mockUserRepoForForm) FindAdminOrganizationIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, errors.New("FindAdminOrganizationIDs not implemented")
}

func (m *mockUserRepoForForm) FindOwnedGroupIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, errors.New("FindOwnedGroupIDs not implemented")
}

func (m *mockUserRepoForForm) FindMemberGroupIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, errors.New("FindMemberGroupIDs not implemented")
}

func (m *mockUserRepoForForm) RecordAgreementAcceptance(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error {
	return errors.New("RecordAgreementAcceptance not implemented")
}

// createTestFormService creates a FormService with mock repositories
func createTestFormService(
	formRepo repositories.FormRepository,
	fieldRepo repositories.FormFieldRepository,
	answerRepo repositories.FormAnswerRepository,
	orgRepo repositories.OrganizationRepository,
	userRepo repositories.UserRepository,
) FormService {
	return NewFormService(formRepo, fieldRepo, answerRepo, orgRepo, userRepo)
}

// ========== ListByOrg Tests ==========

func TestFormServiceImpl_ListByOrg(t *testing.T) {
	tests := []struct {
		name              string
		orgID             string
		userID            string
		isSiteAdminResult bool
		isSiteAdminErr    error
		isOrgAdminResult  bool
		isOrgAdminErr     error
		findFormsResult   []*models.Form
		findFormsErr      error
		want              []*models.Form
		wantErr           error
		wantErrContains   string
	}{
		{
			name:    "empty orgID returns ErrForbidden",
			orgID:   "",
			userID:  "user-1",
			wantErr: ErrForbidden,
		},
		{
			name:            "site admin error returns wrapped error",
			orgID:           "org-1",
			userID:          "user-1",
			isSiteAdminErr:  errors.New("db connection error"),
			wantErrContains: "failed to check site admin status",
		},
		{
			name:              "org admin error returns wrapped error",
			orgID:             "org-1",
			userID:            "user-1",
			isSiteAdminResult: false,
			isOrgAdminErr:     errors.New("db connection error"),
			wantErrContains:   "failed to check org admin status",
		},
		{
			name:              "non-admin user returns ErrForbidden",
			orgID:             "org-1",
			userID:            "user-1",
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			wantErr:           ErrForbidden,
		},
		{
			name:              "site admin can list forms",
			orgID:             "org-1",
			userID:            "site-admin",
			isSiteAdminResult: true,
			isOrgAdminResult:  false,
			findFormsResult:   []*models.Form{{ID: "form-1", Name: "Test Form"}},
			want:              []*models.Form{{ID: "form-1", Name: "Test Form"}},
		},
		{
			name:              "org admin can list forms",
			orgID:             "org-1",
			userID:            "org-admin",
			isSiteAdminResult: false,
			isOrgAdminResult:  true,
			findFormsResult:   []*models.Form{{ID: "form-1", Name: "Test Form"}},
			want:              []*models.Form{{ID: "form-1", Name: "Test Form"}},
		},
		{
			name:              "find forms error returns wrapped error",
			orgID:             "org-1",
			userID:            "site-admin",
			isSiteAdminResult: true,
			findFormsErr:      errors.New("db error"),
			wantErrContains:   "failed to fetch forms for organization",
		},
		{
			name:              "empty forms list returns empty slice",
			orgID:             "org-1",
			userID:            "org-admin",
			isSiteAdminResult: false,
			isOrgAdminResult:  true,
			findFormsResult:   []*models.Form{},
			want:              []*models.Form{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByOrganizationIDFn: func(ctx context.Context, orgID string) ([]*models.Form, error) {
					if tt.findFormsErr != nil {
						return nil, tt.findFormsErr
					}
					return tt.findFormsResult, nil
				},
			}
			orgRepo := &mockOrgRepoForForm{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}
			userRepo := &mockUserRepoForForm{
				isSiteAdminFn: func(ctx context.Context, userID string) (bool, error) {
					if tt.isSiteAdminErr != nil {
						return false, tt.isSiteAdminErr
					}
					return tt.isSiteAdminResult, nil
				},
			}

			svc := createTestFormService(formRepo, &mockFormFieldRepo{}, &mockFormAnswerRepo{}, orgRepo, userRepo)

			got, err := svc.ListByOrg(context.Background(), tt.orgID, tt.userID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ListByOrg() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("ListByOrg() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("ListByOrg() unexpected error = %v", err)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ListByOrg() got %d forms, want %d", len(got), len(tt.want))
			}
		})
	}
}

// ========== Create Tests ==========

func TestFormServiceImpl_Create(t *testing.T) {
	tests := []struct {
		name             string
		orgID            string
		userID           string
		req              *models.CreateFormRequest
		isOrgAdminResult bool
		isOrgAdminErr    error
		createFormResult *models.Form
		createFormErr    error
		want             *models.Form
		wantErr          error
		wantErrContains  string
	}{
		{
			name:            "org admin check error returns wrapped error",
			orgID:           "org-1",
			userID:          "user-1",
			req:             &models.CreateFormRequest{Name: "Test Form"},
			isOrgAdminErr:   errors.New("db error"),
			wantErrContains: "failed to check org admin status",
		},
		{
			name:             "non-org-admin returns ErrForbidden",
			orgID:            "org-1",
			userID:           "user-1",
			req:              &models.CreateFormRequest{Name: "Test Form"},
			isOrgAdminResult: false,
			wantErr:          ErrForbidden,
		},
		{
			name:             "empty name returns ErrValidation",
			orgID:            "org-1",
			userID:           "org-admin",
			req:              &models.CreateFormRequest{Name: ""},
			isOrgAdminResult: true,
			wantErr:          ErrValidation,
		},
		{
			name:             "whitespace only name returns ErrValidation",
			orgID:            "org-1",
			userID:           "org-admin",
			req:              &models.CreateFormRequest{Name: "   "},
			isOrgAdminResult: true,
			wantErr:          ErrValidation,
		},
		{
			name:             "name with leading/trailing whitespace is trimmed",
			orgID:            "org-1",
			userID:           "org-admin",
			req:              &models.CreateFormRequest{Name: "  Test Form  "},
			isOrgAdminResult: true,
			createFormResult: &models.Form{ID: "form-1", Name: "Test Form"},
			want:             &models.Form{ID: "form-1", Name: "Test Form"},
		},
		{
			name:             "nil description passes nil to repository",
			orgID:            "org-1",
			userID:           "org-admin",
			req:              &models.CreateFormRequest{Name: "Test Form", Description: nil},
			isOrgAdminResult: true,
			createFormResult: &models.Form{ID: "form-1", Name: "Test Form"},
			want:             &models.Form{ID: "form-1", Name: "Test Form"},
		},
		{
			name:             "description with whitespace is trimmed",
			orgID:            "org-1",
			userID:           "org-admin",
			req:              &models.CreateFormRequest{Name: "Test Form", Description: strPtr("  Desc  ")},
			isOrgAdminResult: true,
			createFormResult: &models.Form{ID: "form-1", Name: "Test Form", Description: strPtr("Desc")},
			want:             &models.Form{ID: "form-1", Name: "Test Form", Description: strPtr("Desc")},
		},
		{
			name:             "create form error returns wrapped error",
			orgID:            "org-1",
			userID:           "org-admin",
			req:              &models.CreateFormRequest{Name: "Test Form"},
			isOrgAdminResult: true,
			createFormErr:    errors.New("db error"),
			wantErrContains:  "failed to create form",
		},
		{
			name:             "successful creation returns form",
			orgID:            "org-1",
			userID:           "org-admin",
			req:              &models.CreateFormRequest{Name: "Test Form"},
			isOrgAdminResult: true,
			createFormResult: &models.Form{ID: "form-1", Name: "Test Form", OrganizationID: "org-1"},
			want:             &models.Form{ID: "form-1", Name: "Test Form", OrganizationID: "org-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				createFn: func(ctx context.Context, orgID, name string, description *string, createdBy string) (*models.Form, error) {
					if tt.createFormErr != nil {
						return nil, tt.createFormErr
					}
					return tt.createFormResult, nil
				},
			}
			orgRepo := &mockOrgRepoForForm{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}

			svc := createTestFormService(formRepo, &mockFormFieldRepo{}, &mockFormAnswerRepo{}, orgRepo, &mockUserRepoForForm{})

			got, err := svc.Create(context.Background(), tt.orgID, tt.userID, tt.req)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("Create() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("Create() unexpected error = %v", err)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("Create() got ID = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

// ========== GetByID Tests ==========

func TestFormServiceImpl_GetByID(t *testing.T) {
	tests := []struct {
		name            string
		formID          string
		findFormResult  *models.Form
		findFormErr     error
		want            *models.Form
		wantErr         error
		wantErrContains string
	}{
		{
			name:    "empty formID returns ErrFormNotFound",
			formID:  "",
			wantErr: ErrFormNotFound,
		},
		{
			name:            "find form error returns wrapped error",
			formID:          "form-1",
			findFormErr:     errors.New("db error"),
			wantErrContains: "failed to get form by ID",
		},
		{
			name:           "form not found returns ErrFormNotFound",
			formID:         "form-1",
			findFormResult: nil,
			wantErr:        ErrFormNotFound,
		},
		{
			name:           "form found returns form",
			formID:         "form-1",
			findFormResult: &models.Form{ID: "form-1", Name: "Test Form"},
			want:           &models.Form{ID: "form-1", Name: "Test Form"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
					if tt.findFormErr != nil {
						return nil, tt.findFormErr
					}
					return tt.findFormResult, nil
				},
			}

			svc := createTestFormService(formRepo, &mockFormFieldRepo{}, &mockFormAnswerRepo{}, &mockOrgRepoForForm{}, &mockUserRepoForForm{})

			got, err := svc.GetByID(context.Background(), tt.formID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("GetByID() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("GetByID() unexpected error = %v", err)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("GetByID() got ID = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

// ========== GetDetailByID Tests ==========

func TestFormServiceImpl_GetDetailByID(t *testing.T) {
	tests := []struct {
		name              string
		formID            string
		userID            string
		findFormResult    *models.Form
		findFormErr       error
		isSiteAdminResult bool
		isSiteAdminErr    error
		isOrgAdminResult  bool
		isOrgAdminErr     error
		isOrgMemberResult bool
		isOrgMemberErr    error
		getFieldsResult   []models.FormField
		getFieldsErr      error
		want              *models.Form
		wantErr           error
		wantErrContains   string
	}{
		{
			name:            "find form error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			findFormErr:     errors.New("db error"),
			wantErrContains: "failed to find form",
		},
		{
			name:           "form not found returns ErrFormNotFound",
			formID:         "form-1",
			userID:         "user-1",
			findFormResult: nil,
			wantErr:        ErrFormNotFound,
		},
		{
			name:            "site admin check error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			findFormResult:  &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isSiteAdminErr:  errors.New("db error"),
			wantErrContains: "failed to check site admin status",
		},
		{
			name:              "org admin check error returns wrapped error",
			formID:            "form-1",
			userID:            "user-1",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminErr:     errors.New("db error"),
			wantErrContains:   "failed to check org admin status",
		},
		{
			name:              "org member check error returns wrapped error",
			formID:            "form-1",
			userID:            "user-1",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			isOrgMemberErr:    errors.New("db error"),
			wantErrContains:   "failed to check org membership",
		},
		{
			name:              "non-member non-admin returns ErrForbidden",
			formID:            "form-1",
			userID:            "user-1",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			isOrgMemberResult: false,
			wantErr:           ErrForbidden,
		},
		{
			name:              "org member can get detail",
			formID:            "form-1",
			userID:            "member-1",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1", Name: "Test Form"},
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			isOrgMemberResult: true,
			getFieldsResult:   []models.FormField{{ID: "field-1", Label: "Question 1"}},
			want:              &models.Form{ID: "form-1", OrganizationID: "org-1", Name: "Test Form", Fields: []models.FormField{{ID: "field-1", Label: "Question 1"}}},
		},
		{
			name:              "get fields error returns wrapped error",
			formID:            "form-1",
			userID:            "org-admin",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminResult:  true,
			getFieldsErr:      errors.New("db error"),
			wantErrContains:   "failed to get form fields",
		},
		{
			name:              "site admin can get detail with no fields",
			formID:            "form-1",
			userID:            "site-admin",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1", Name: "Test Form"},
			isSiteAdminResult: true,
			getFieldsResult:   []models.FormField{},
			want:              &models.Form{ID: "form-1", OrganizationID: "org-1", Name: "Test Form", Fields: []models.FormField{}},
		},
		{
			name:              "org admin can get detail with fields",
			formID:            "form-1",
			userID:            "org-admin",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1", Name: "Test Form"},
			isSiteAdminResult: false,
			isOrgAdminResult:  true,
			getFieldsResult:   []models.FormField{{ID: "field-1", Label: "Question 1"}},
			want:              &models.Form{ID: "form-1", OrganizationID: "org-1", Name: "Test Form", Fields: []models.FormField{{ID: "field-1", Label: "Question 1"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
					if tt.findFormErr != nil {
						return nil, tt.findFormErr
					}
					return tt.findFormResult, nil
				},
			}
			fieldRepo := &mockFormFieldRepo{
				getFieldsFn: func(ctx context.Context, formID string) ([]models.FormField, error) {
					if tt.getFieldsErr != nil {
						return nil, tt.getFieldsErr
					}
					return tt.getFieldsResult, nil
				},
			}
			orgRepo := &mockOrgRepoForForm{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
				isMemberFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgMemberErr != nil {
						return false, tt.isOrgMemberErr
					}
					return tt.isOrgMemberResult, nil
				},
			}
			userRepo := &mockUserRepoForForm{
				isSiteAdminFn: func(ctx context.Context, userID string) (bool, error) {
					if tt.isSiteAdminErr != nil {
						return false, tt.isSiteAdminErr
					}
					return tt.isSiteAdminResult, nil
				},
			}

			svc := createTestFormService(formRepo, fieldRepo, &mockFormAnswerRepo{}, orgRepo, userRepo)

			got, err := svc.GetDetailByID(context.Background(), tt.formID, tt.userID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetDetailByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("GetDetailByID() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("GetDetailByID() unexpected error = %v", err)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("GetDetailByID() got ID = %v, want %v", got.ID, tt.want.ID)
			}
			if len(got.Fields) != len(tt.want.Fields) {
				t.Errorf("GetDetailByID() got %d fields, want %d", len(got.Fields), len(tt.want.Fields))
			}
		})
	}
}

// ========== Update Tests ==========

func TestFormServiceImpl_Update(t *testing.T) {
	tests := []struct {
		name             string
		formID           string
		userID           string
		req              *models.UpdateFormRequest
		findFormResult   *models.Form
		findFormErr      error
		isOrgAdminResult bool
		isOrgAdminErr    error
		updateFormResult *models.Form
		updateFormErr    error
		want             *models.Form
		wantErr          error
		wantErrContains  string
	}{
		{
			name:            "find form error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			req:             &models.UpdateFormRequest{Name: "Updated Form"},
			findFormErr:     errors.New("db error"),
			wantErrContains: "failed to find form",
		},
		{
			name:           "form not found returns ErrFormNotFound",
			formID:         "form-1",
			userID:         "user-1",
			req:            &models.UpdateFormRequest{Name: "Updated Form"},
			findFormResult: nil,
			wantErr:        ErrFormNotFound,
		},
		{
			name:            "org admin check error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			req:             &models.UpdateFormRequest{Name: "Updated Form"},
			findFormResult:  &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminErr:   errors.New("db error"),
			wantErrContains: "failed to check org admin status",
		},
		{
			name:             "non-org-admin returns ErrForbidden",
			formID:           "form-1",
			userID:           "user-1",
			req:              &models.UpdateFormRequest{Name: "Updated Form"},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: false,
			wantErr:          ErrForbidden,
		},
		{
			name:             "empty name returns ErrValidation",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.UpdateFormRequest{Name: ""},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			wantErr:          ErrValidation,
		},
		{
			name:             "whitespace only name returns ErrValidation",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.UpdateFormRequest{Name: "   "},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			wantErr:          ErrValidation,
		},
		{
			name:             "update form error returns wrapped error",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.UpdateFormRequest{Name: "Updated Form"},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			updateFormErr:    errors.New("db error"),
			wantErrContains:  "failed to update form",
		},
		{
			name:             "successful update returns updated form",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.UpdateFormRequest{Name: "  Updated Form  ", Description: strPtr("  Desc  "), IsActive: true},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			updateFormResult: &models.Form{ID: "form-1", Name: "Updated Form", Description: strPtr("Desc"), IsActive: true},
			want:             &models.Form{ID: "form-1", Name: "Updated Form", Description: strPtr("Desc"), IsActive: true},
		},
		{
			name:             "nil description passes nil to repository",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.UpdateFormRequest{Name: "Updated Form", Description: nil, IsActive: false},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			updateFormResult: &models.Form{ID: "form-1", Name: "Updated Form", IsActive: false},
			want:             &models.Form{ID: "form-1", Name: "Updated Form", IsActive: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
					if tt.findFormErr != nil {
						return nil, tt.findFormErr
					}
					return tt.findFormResult, nil
				},
				updateFn: func(ctx context.Context, formID, name string, description *string, isActive bool) (*models.Form, error) {
					if tt.updateFormErr != nil {
						return nil, tt.updateFormErr
					}
					return tt.updateFormResult, nil
				},
			}
			orgRepo := &mockOrgRepoForForm{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}

			svc := createTestFormService(formRepo, &mockFormFieldRepo{}, &mockFormAnswerRepo{}, orgRepo, &mockUserRepoForForm{})

			got, err := svc.Update(context.Background(), tt.formID, tt.userID, tt.req)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("Update() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("Update() unexpected error = %v", err)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("Update() got ID = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

// ========== Delete Tests ==========

func TestFormServiceImpl_Delete(t *testing.T) {
	tests := []struct {
		name             string
		formID           string
		userID           string
		findFormResult   *models.Form
		findFormErr      error
		isOrgAdminResult bool
		isOrgAdminErr    error
		deleteFormErr    error
		wantErr          error
		wantErrContains  string
	}{
		{
			name:            "find form error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			findFormErr:     errors.New("db error"),
			wantErrContains: "failed to find form",
		},
		{
			name:           "form not found returns ErrFormNotFound",
			formID:         "form-1",
			userID:         "user-1",
			findFormResult: nil,
			wantErr:        ErrFormNotFound,
		},
		{
			name:            "org admin check error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			findFormResult:  &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminErr:   errors.New("db error"),
			wantErrContains: "failed to check org admin status",
		},
		{
			name:             "non-org-admin returns ErrForbidden",
			formID:           "form-1",
			userID:           "user-1",
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: false,
			wantErr:          ErrForbidden,
		},
		{
			name:             "delete form error returns wrapped error",
			formID:           "form-1",
			userID:           "org-admin",
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			deleteFormErr:    errors.New("db error"),
			wantErrContains:  "failed to delete form",
		},
		{
			name:             "successful delete returns nil",
			formID:           "form-1",
			userID:           "org-admin",
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
					if tt.findFormErr != nil {
						return nil, tt.findFormErr
					}
					return tt.findFormResult, nil
				},
				deleteFn: func(ctx context.Context, formID string) error {
					return tt.deleteFormErr
				},
			}
			orgRepo := &mockOrgRepoForForm{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}

			svc := createTestFormService(formRepo, &mockFormFieldRepo{}, &mockFormAnswerRepo{}, orgRepo, &mockUserRepoForForm{})

			err := svc.Delete(context.Background(), tt.formID, tt.userID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("Delete() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("Delete() unexpected error = %v", err)
			}
		})
	}
}

// ========== AddField Tests ==========

func TestFormServiceImpl_AddField(t *testing.T) {
	tests := []struct {
		name             string
		formID           string
		userID           string
		req              *models.AddFormFieldRequest
		findFormResult   *models.Form
		findFormErr      error
		isOrgAdminResult bool
		isOrgAdminErr    error
		getFieldsResult  []models.FormField
		getFieldsErr     error
		addFieldResult   *models.FormField
		addFieldErr      error
		want             *models.FormField
		wantErr          error
		wantErrContains  string
	}{
		{
			name:            "find form error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			req:             &models.AddFormFieldRequest{FieldType: "TEXT_SMALL", Label: "Question"},
			findFormErr:     errors.New("db error"),
			wantErrContains: "failed to find form",
		},
		{
			name:           "form not found returns ErrFormNotFound",
			formID:         "form-1",
			userID:         "user-1",
			req:            &models.AddFormFieldRequest{FieldType: "TEXT_SMALL", Label: "Question"},
			findFormResult: nil,
			wantErr:        ErrFormNotFound,
		},
		{
			name:            "org admin check error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			req:             &models.AddFormFieldRequest{FieldType: "TEXT_SMALL", Label: "Question"},
			findFormResult:  &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminErr:   errors.New("db error"),
			wantErrContains: "failed to check org admin status",
		},
		{
			name:             "non-org-admin returns ErrForbidden",
			formID:           "form-1",
			userID:           "user-1",
			req:              &models.AddFormFieldRequest{FieldType: "TEXT_SMALL", Label: "Question"},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: false,
			wantErr:          ErrForbidden,
		},
		{
			name:             "empty label returns ErrValidation",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "TEXT_SMALL", Label: ""},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			wantErr:          ErrValidation,
		},
		{
			name:             "whitespace only label returns ErrValidation",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "TEXT_SMALL", Label: "   "},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			wantErr:          ErrValidation,
		},
		{
			name:             "CHECKBOX without options returns ErrValidation",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "CHECKBOX", Label: "Question", Options: []string{}},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			wantErr:          ErrValidation,
		},
		{
			name:             "RADIO without options returns ErrValidation",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "RADIO", Label: "Question", Options: nil},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			wantErr:          ErrValidation,
		},
		{
			name:             "DROPDOWN without options returns ErrValidation",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "DROPDOWN", Label: "Question"},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			wantErr:          ErrValidation,
		},
		{
			name:             "get fields error returns wrapped error",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "TEXT_SMALL", Label: "Question"},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			getFieldsErr:     errors.New("db error"),
			wantErrContains:  "failed to get existing fields",
		},
		{
			name:             "add field error returns wrapped error",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "TEXT_SMALL", Label: "Question"},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			getFieldsResult:  []models.FormField{},
			addFieldErr:      errors.New("db error"),
			wantErrContains:  "failed to add field",
		},
		{
			name:             "first field gets order 0",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "TEXT_SMALL", Label: "  Question  ", IsRequired: true},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			getFieldsResult:  []models.FormField{},
			addFieldResult:   &models.FormField{ID: "field-1", Label: "Question", FieldOrder: 0},
			want:             &models.FormField{ID: "field-1", Label: "Question", FieldOrder: 0},
		},
		{
			name:             "new field gets next order after existing",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "TEXT_SMALL", Label: "Question"},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			getFieldsResult:  []models.FormField{{ID: "field-1", FieldOrder: 0}, {ID: "field-2", FieldOrder: 2}},
			addFieldResult:   &models.FormField{ID: "field-3", Label: "Question", FieldOrder: 3},
			want:             &models.FormField{ID: "field-3", Label: "Question", FieldOrder: 3},
		},
		{
			name:             "TEXT field without options is valid",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "TEXT_SMALL", Label: "Question", Options: nil},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			getFieldsResult:  []models.FormField{},
			addFieldResult:   &models.FormField{ID: "field-1", Label: "Question"},
			want:             &models.FormField{ID: "field-1", Label: "Question"},
		},
		{
			name:             "CHECKBOX with options is valid",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "CHECKBOX", Label: "Question", Options: []string{"A", "B"}},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			getFieldsResult:  []models.FormField{},
			addFieldResult:   &models.FormField{ID: "field-1", Label: "Question", Options: []string{"A", "B"}},
			want:             &models.FormField{ID: "field-1", Label: "Question", Options: []string{"A", "B"}},
		},
		{
			name:             "options with whitespace are trimmed and empty removed",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "RADIO", Label: "Question", Options: []string{"  A  ", "", "  B  ", "   "}},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			getFieldsResult:  []models.FormField{},
			addFieldResult:   &models.FormField{ID: "field-1", Label: "Question", Options: []string{"A", "B"}},
			want:             &models.FormField{ID: "field-1", Label: "Question", Options: []string{"A", "B"}},
		},
		{
			name:             "description is trimmed",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.AddFormFieldRequest{FieldType: "TEXT_SMALL", Label: "Question", Description: strPtr("  Desc  ")},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			getFieldsResult:  []models.FormField{},
			addFieldResult:   &models.FormField{ID: "field-1", Label: "Question", Description: strPtr("Desc")},
			want:             &models.FormField{ID: "field-1", Label: "Question", Description: strPtr("Desc")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
					if tt.findFormErr != nil {
						return nil, tt.findFormErr
					}
					return tt.findFormResult, nil
				},
			}
			fieldRepo := &mockFormFieldRepo{
				getFieldsFn: func(ctx context.Context, formID string) ([]models.FormField, error) {
					if tt.getFieldsErr != nil {
						return nil, tt.getFieldsErr
					}
					return tt.getFieldsResult, nil
				},
				addFieldFn: func(ctx context.Context, formID string, fieldType models.FormFieldType, label string, description *string, isRequired bool, fieldOrder int, options []string) (*models.FormField, error) {
					if tt.addFieldErr != nil {
						return nil, tt.addFieldErr
					}
					return tt.addFieldResult, nil
				},
			}
			orgRepo := &mockOrgRepoForForm{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}

			svc := createTestFormService(formRepo, fieldRepo, &mockFormAnswerRepo{}, orgRepo, &mockUserRepoForForm{})

			got, err := svc.AddField(context.Background(), tt.formID, tt.userID, tt.req)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("AddField() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("AddField() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("AddField() unexpected error = %v", err)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("AddField() got ID = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

// ========== ReorderFields Tests ==========

func TestFormServiceImpl_ReorderFields(t *testing.T) {
	tests := []struct {
		name             string
		formID           string
		userID           string
		req              *models.ReorderFieldsRequest
		findFormResult   *models.Form
		findFormErr      error
		isOrgAdminResult bool
		isOrgAdminErr    error
		reorderErr       error
		wantErr          error
		wantErrContains  string
	}{
		{
			name:            "find form error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			req:             &models.ReorderFieldsRequest{FieldIDs: []string{"field-1", "field-2"}},
			findFormErr:     errors.New("db error"),
			wantErrContains: "failed to find form",
		},
		{
			name:           "form not found returns ErrFormNotFound",
			formID:         "form-1",
			userID:         "user-1",
			req:            &models.ReorderFieldsRequest{FieldIDs: []string{"field-1"}},
			findFormResult: nil,
			wantErr:        ErrFormNotFound,
		},
		{
			name:            "org admin check error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			req:             &models.ReorderFieldsRequest{FieldIDs: []string{"field-1"}},
			findFormResult:  &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminErr:   errors.New("db error"),
			wantErrContains: "failed to check org admin status",
		},
		{
			name:             "non-org-admin returns ErrForbidden",
			formID:           "form-1",
			userID:           "user-1",
			req:              &models.ReorderFieldsRequest{FieldIDs: []string{"field-1"}},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: false,
			wantErr:          ErrForbidden,
		},
		{
			name:             "reorder error returns wrapped error",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.ReorderFieldsRequest{FieldIDs: []string{"field-1"}},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			reorderErr:       errors.New("db error"),
			wantErrContains:  "failed to reorder fields",
		},
		{
			name:             "successful reorder returns nil",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.ReorderFieldsRequest{FieldIDs: []string{"field-2", "field-1", "field-3"}},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
		},
		{
			name:             "empty fieldIDs is passed to repository",
			formID:           "form-1",
			userID:           "org-admin",
			req:              &models.ReorderFieldsRequest{FieldIDs: []string{}},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
					if tt.findFormErr != nil {
						return nil, tt.findFormErr
					}
					return tt.findFormResult, nil
				},
			}
			fieldRepo := &mockFormFieldRepo{
				reorderFieldsFn: func(ctx context.Context, formID string, fieldIDs []string) error {
					return tt.reorderErr
				},
			}
			orgRepo := &mockOrgRepoForForm{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}

			svc := createTestFormService(formRepo, fieldRepo, &mockFormAnswerRepo{}, orgRepo, &mockUserRepoForForm{})

			err := svc.ReorderFields(context.Background(), tt.formID, tt.userID, tt.req)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ReorderFields() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("ReorderFields() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("ReorderFields() unexpected error = %v", err)
			}
		})
	}
}

// ========== UpdateField Tests ==========

func TestFormServiceImpl_UpdateField(t *testing.T) {
	tests := []struct {
		name              string
		formID            string
		fieldID           string
		userID            string
		req               *models.UpdateFormFieldRequest
		findFormResult    *models.Form
		findFormErr       error
		isOrgAdminResult  bool
		isOrgAdminErr     error
		updateFieldResult *models.FormField
		updateFieldErr    error
		want              *models.FormField
		wantErr           error
		wantErrContains   string
	}{
		{
			name:            "find form error returns wrapped error",
			formID:          "form-1",
			fieldID:         "field-1",
			userID:          "user-1",
			req:             &models.UpdateFormFieldRequest{Label: "Updated Question"},
			findFormErr:     errors.New("db error"),
			wantErrContains: "failed to find form",
		},
		{
			name:           "form not found returns ErrFormNotFound",
			formID:         "form-1",
			fieldID:        "field-1",
			userID:         "user-1",
			req:            &models.UpdateFormFieldRequest{Label: "Updated Question"},
			findFormResult: nil,
			wantErr:        ErrFormNotFound,
		},
		{
			name:            "org admin check error returns wrapped error",
			formID:          "form-1",
			fieldID:         "field-1",
			userID:          "user-1",
			req:             &models.UpdateFormFieldRequest{Label: "Updated Question"},
			findFormResult:  &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminErr:   errors.New("db error"),
			wantErrContains: "failed to check org admin status",
		},
		{
			name:             "non-org-admin returns ErrForbidden",
			formID:           "form-1",
			fieldID:          "field-1",
			userID:           "user-1",
			req:              &models.UpdateFormFieldRequest{Label: "Updated Question"},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: false,
			wantErr:          ErrForbidden,
		},
		{
			name:             "empty label returns ErrValidation",
			formID:           "form-1",
			fieldID:          "field-1",
			userID:           "org-admin",
			req:              &models.UpdateFormFieldRequest{Label: ""},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			wantErr:          ErrValidation,
		},
		{
			name:             "whitespace only label returns ErrValidation",
			formID:           "form-1",
			fieldID:          "field-1",
			userID:           "org-admin",
			req:              &models.UpdateFormFieldRequest{Label: "   "},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			wantErr:          ErrValidation,
		},
		{
			name:             "update field error returns wrapped error",
			formID:           "form-1",
			fieldID:          "field-1",
			userID:           "org-admin",
			req:              &models.UpdateFormFieldRequest{Label: "Updated Question"},
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			updateFieldErr:   errors.New("db error"),
			wantErrContains:  "failed to update field",
		},
		{
			name:              "successful update returns updated field",
			formID:            "form-1",
			fieldID:           "field-1",
			userID:            "org-admin",
			req:               &models.UpdateFormFieldRequest{Label: "  Updated Question  ", Description: strPtr("  Desc  "), IsRequired: true, Options: []string{"  A  ", "  B  "}},
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult:  true,
			updateFieldResult: &models.FormField{ID: "field-1", Label: "Updated Question", Description: strPtr("Desc"), IsRequired: true, Options: []string{"A", "B"}},
			want:              &models.FormField{ID: "field-1", Label: "Updated Question", Description: strPtr("Desc"), IsRequired: true, Options: []string{"A", "B"}},
		},
		{
			name:              "nil description passes nil",
			formID:            "form-1",
			fieldID:           "field-1",
			userID:            "org-admin",
			req:               &models.UpdateFormFieldRequest{Label: "Question", Description: nil},
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult:  true,
			updateFieldResult: &models.FormField{ID: "field-1", Label: "Question"},
			want:              &models.FormField{ID: "field-1", Label: "Question"},
		},
		{
			name:              "empty options after trimming",
			formID:            "form-1",
			fieldID:           "field-1",
			userID:            "org-admin",
			req:               &models.UpdateFormFieldRequest{Label: "Question", Options: []string{"  ", "", "   "}},
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult:  true,
			updateFieldResult: &models.FormField{ID: "field-1", Label: "Question", Options: []string{}},
			want:              &models.FormField{ID: "field-1", Label: "Question", Options: []string{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
					if tt.findFormErr != nil {
						return nil, tt.findFormErr
					}
					return tt.findFormResult, nil
				},
			}
			fieldRepo := &mockFormFieldRepo{
				updateFieldFn: func(ctx context.Context, fieldID, label string, description *string, isRequired bool, options []string) (*models.FormField, error) {
					if tt.updateFieldErr != nil {
						return nil, tt.updateFieldErr
					}
					return tt.updateFieldResult, nil
				},
			}
			orgRepo := &mockOrgRepoForForm{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}

			svc := createTestFormService(formRepo, fieldRepo, &mockFormAnswerRepo{}, orgRepo, &mockUserRepoForForm{})

			got, err := svc.UpdateField(context.Background(), tt.formID, tt.fieldID, tt.userID, tt.req)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("UpdateField() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("UpdateField() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("UpdateField() unexpected error = %v", err)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("UpdateField() got ID = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

// ========== DeleteField Tests ==========

func TestFormServiceImpl_DeleteField(t *testing.T) {
	tests := []struct {
		name             string
		formID           string
		fieldID          string
		userID           string
		findFormResult   *models.Form
		findFormErr      error
		isOrgAdminResult bool
		isOrgAdminErr    error
		deleteFieldErr   error
		wantErr          error
		wantErrContains  string
	}{
		{
			name:            "find form error returns wrapped error",
			formID:          "form-1",
			fieldID:         "field-1",
			userID:          "user-1",
			findFormErr:     errors.New("db error"),
			wantErrContains: "failed to find form",
		},
		{
			name:           "form not found returns ErrFormNotFound",
			formID:         "form-1",
			fieldID:        "field-1",
			userID:         "user-1",
			findFormResult: nil,
			wantErr:        ErrFormNotFound,
		},
		{
			name:            "org admin check error returns wrapped error",
			formID:          "form-1",
			fieldID:         "field-1",
			userID:          "user-1",
			findFormResult:  &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminErr:   errors.New("db error"),
			wantErrContains: "failed to check org admin status",
		},
		{
			name:             "non-org-admin returns ErrForbidden",
			formID:           "form-1",
			fieldID:          "field-1",
			userID:           "user-1",
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: false,
			wantErr:          ErrForbidden,
		},
		{
			name:             "delete field error returns wrapped error",
			formID:           "form-1",
			fieldID:          "field-1",
			userID:           "org-admin",
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			deleteFieldErr:   errors.New("db error"),
			wantErrContains:  "failed to delete field",
		},
		{
			name:             "successful delete returns nil",
			formID:           "form-1",
			fieldID:          "field-1",
			userID:           "org-admin",
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
					if tt.findFormErr != nil {
						return nil, tt.findFormErr
					}
					return tt.findFormResult, nil
				},
			}
			fieldRepo := &mockFormFieldRepo{
				deleteFieldFn: func(ctx context.Context, fieldID string) error {
					return tt.deleteFieldErr
				},
			}
			orgRepo := &mockOrgRepoForForm{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}

			svc := createTestFormService(formRepo, fieldRepo, &mockFormAnswerRepo{}, orgRepo, &mockUserRepoForForm{})

			err := svc.DeleteField(context.Background(), tt.formID, tt.fieldID, tt.userID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("DeleteField() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("DeleteField() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("DeleteField() unexpected error = %v", err)
			}
		})
	}
}

// ========== SubmitAnswer Tests ==========

func TestFormServiceImpl_SubmitAnswer(t *testing.T) {
	tests := []struct {
		name               string
		formID             string
		userID             string
		req                *models.SubmitFormAnswerRequest
		findFormResult     *models.Form
		findFormErr        error
		submitAnswerResult *models.FormAnswer
		submitAnswerErr    error
		want               *models.FormAnswer
		wantErr            error
		wantErrContains    string
	}{
		{
			name:            "find form error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			req:             &models.SubmitFormAnswerRequest{Answers: json.RawMessage(`{}`)},
			findFormErr:     errors.New("db error"),
			wantErrContains: "failed to find form",
		},
		{
			name:           "form not found returns ErrFormNotFound",
			formID:         "form-1",
			userID:         "user-1",
			req:            &models.SubmitFormAnswerRequest{Answers: json.RawMessage(`{}`)},
			findFormResult: nil,
			wantErr:        ErrFormNotFound,
		},
		{
			name:           "inactive form returns ErrValidation",
			formID:         "form-1",
			userID:         "user-1",
			req:            &models.SubmitFormAnswerRequest{Answers: json.RawMessage(`{}`)},
			findFormResult: &models.Form{ID: "form-1", IsActive: false},
			wantErr:        ErrValidation,
		},
		{
			name:            "submit answer error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			req:             &models.SubmitFormAnswerRequest{Answers: json.RawMessage(`{}`)},
			findFormResult:  &models.Form{ID: "form-1", IsActive: true},
			submitAnswerErr: errors.New("db error"),
			wantErrContains: "failed to submit answer",
		},
		{
			name:               "successful submission returns answer",
			formID:             "form-1",
			userID:             "user-1",
			req:                &models.SubmitFormAnswerRequest{Answers: json.RawMessage(`{"field-1": "answer"}`)},
			findFormResult:     &models.Form{ID: "form-1", IsActive: true},
			submitAnswerResult: &models.FormAnswer{ID: "answer-1", FormID: "form-1", UserID: "user-1"},
			want:               &models.FormAnswer{ID: "answer-1", FormID: "form-1", UserID: "user-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
					if tt.findFormErr != nil {
						return nil, tt.findFormErr
					}
					return tt.findFormResult, nil
				},
			}
			answerRepo := &mockFormAnswerRepo{
				submitFn: func(ctx context.Context, formID, userID string, answers json.RawMessage) (*models.FormAnswer, error) {
					if tt.submitAnswerErr != nil {
						return nil, tt.submitAnswerErr
					}
					return tt.submitAnswerResult, nil
				},
			}

			svc := createTestFormService(formRepo, &mockFormFieldRepo{}, answerRepo, &mockOrgRepoForForm{}, &mockUserRepoForForm{})

			got, err := svc.SubmitAnswer(context.Background(), tt.formID, tt.userID, tt.req)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("SubmitAnswer() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("SubmitAnswer() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("SubmitAnswer() unexpected error = %v", err)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("SubmitAnswer() got ID = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

// ========== ListAnswers Tests ==========

func TestFormServiceImpl_ListAnswers(t *testing.T) {
	tests := []struct {
		name              string
		formID            string
		userID            string
		findFormResult    *models.Form
		findFormErr       error
		isOrgAdminResult  bool
		isOrgAdminErr     error
		findAnswersResult []*models.FormAnswer
		findAnswersErr    error
		want              []*models.FormAnswer
		wantErr           error
		wantErrContains   string
	}{
		{
			name:            "find form error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			findFormErr:     errors.New("db error"),
			wantErrContains: "failed to find form",
		},
		{
			name:           "form not found returns ErrFormNotFound",
			formID:         "form-1",
			userID:         "user-1",
			findFormResult: nil,
			wantErr:        ErrFormNotFound,
		},
		{
			name:            "org admin check error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			findFormResult:  &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminErr:   errors.New("db error"),
			wantErrContains: "failed to check org admin status",
		},
		{
			name:             "non-org-admin returns ErrForbidden",
			formID:           "form-1",
			userID:           "user-1",
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: false,
			wantErr:          ErrForbidden,
		},
		{
			name:             "find answers error returns wrapped error",
			formID:           "form-1",
			userID:           "org-admin",
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: true,
			findAnswersErr:   errors.New("db error"),
			wantErrContains:  "failed to find answers",
		},
		{
			name:              "no answers returns empty slice",
			formID:            "form-1",
			userID:            "org-admin",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult:  true,
			findAnswersResult: []*models.FormAnswer{},
			want:              []*models.FormAnswer{},
		},
		{
			name:              "answers found returns all answers",
			formID:            "form-1",
			userID:            "org-admin",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult:  true,
			findAnswersResult: []*models.FormAnswer{{ID: "answer-1"}, {ID: "answer-2"}},
			want:              []*models.FormAnswer{{ID: "answer-1"}, {ID: "answer-2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
					if tt.findFormErr != nil {
						return nil, tt.findFormErr
					}
					return tt.findFormResult, nil
				},
			}
			answerRepo := &mockFormAnswerRepo{
				findByFormFn: func(ctx context.Context, formID string) ([]*models.FormAnswer, error) {
					if tt.findAnswersErr != nil {
						return nil, tt.findAnswersErr
					}
					return tt.findAnswersResult, nil
				},
			}
			orgRepo := &mockOrgRepoForForm{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}

			svc := createTestFormService(formRepo, &mockFormFieldRepo{}, answerRepo, orgRepo, &mockUserRepoForForm{})

			got, err := svc.ListAnswers(context.Background(), tt.formID, tt.userID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ListAnswers() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("ListAnswers() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("ListAnswers() unexpected error = %v", err)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ListAnswers() got %d answers, want %d", len(got), len(tt.want))
			}
		})
	}
}

// ========== GetMyAnswer Tests ==========

func TestFormServiceImpl_GetMyAnswer(t *testing.T) {
	tests := []struct {
		name              string
		formID            string
		userID            string
		findFormResult    *models.Form
		findFormErr       error
		findCurrentResult *models.FormAnswer
		findCurrentErr    error
		want              *models.FormAnswer
		wantErr           error
		wantErrContains   string
	}{
		{
			name:            "find form error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			findFormErr:     errors.New("db error"),
			wantErrContains: "failed to find form",
		},
		{
			name:           "form not found returns ErrFormNotFound",
			formID:         "form-1",
			userID:         "user-1",
			findFormResult: nil,
			wantErr:        ErrFormNotFound,
		},
		{
			name:            "find current answer error returns wrapped error",
			formID:          "form-1",
			userID:          "user-1",
			findFormResult:  &models.Form{ID: "form-1"},
			findCurrentErr:  errors.New("db error"),
			wantErrContains: "failed to find current answer",
		},
		{
			name:              "no answer returns nil without error",
			formID:            "form-1",
			userID:            "user-1",
			findFormResult:    &models.Form{ID: "form-1"},
			findCurrentResult: nil,
			want:              nil,
		},
		{
			name:              "answer found returns the answer",
			formID:            "form-1",
			userID:            "user-1",
			findFormResult:    &models.Form{ID: "form-1"},
			findCurrentResult: &models.FormAnswer{ID: "answer-1", UserID: "user-1"},
			want:              &models.FormAnswer{ID: "answer-1", UserID: "user-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
					if tt.findFormErr != nil {
						return nil, tt.findFormErr
					}
					return tt.findFormResult, nil
				},
			}
			answerRepo := &mockFormAnswerRepo{
				findCurrentFn: func(ctx context.Context, formID, userID string) (*models.FormAnswer, error) {
					if tt.findCurrentErr != nil {
						return nil, tt.findCurrentErr
					}
					return tt.findCurrentResult, nil
				},
			}

			svc := createTestFormService(formRepo, &mockFormFieldRepo{}, answerRepo, &mockOrgRepoForForm{}, &mockUserRepoForForm{})

			got, err := svc.GetMyAnswer(context.Background(), tt.formID, tt.userID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetMyAnswer() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("GetMyAnswer() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("GetMyAnswer() unexpected error = %v", err)
				return
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("GetMyAnswer() got = %v, want nil", got)
				}
			} else {
				if got == nil || got.ID != tt.want.ID {
					t.Errorf("GetMyAnswer() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// ========== GetAnswerHistory Tests ==========

func TestFormServiceImpl_GetAnswerHistory(t *testing.T) {
	tests := []struct {
		name              string
		formID            string
		targetUserID      string
		requesterID       string
		findFormResult    *models.Form
		findFormErr       error
		isOrgAdminResult  bool
		isOrgAdminErr     error
		findHistoryResult []*models.FormAnswer
		findHistoryErr    error
		want              []*models.FormAnswer
		wantErr           error
		wantErrContains   string
	}{
		{
			name:            "find form error returns wrapped error",
			formID:          "form-1",
			targetUserID:    "user-1",
			requesterID:     "requester-1",
			findFormErr:     errors.New("db error"),
			wantErrContains: "failed to find form",
		},
		{
			name:           "form not found returns ErrFormNotFound",
			formID:         "form-1",
			targetUserID:   "user-1",
			requesterID:    "requester-1",
			findFormResult: nil,
			wantErr:        ErrFormNotFound,
		},
		{
			name:              "viewing own history is allowed",
			formID:            "form-1",
			targetUserID:      "user-1",
			requesterID:       "user-1",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			findHistoryResult: []*models.FormAnswer{{ID: "answer-1"}, {ID: "answer-2"}},
			want:              []*models.FormAnswer{{ID: "answer-1"}, {ID: "answer-2"}},
		},
		{
			name:            "org admin check error returns wrapped error",
			formID:          "form-1",
			targetUserID:    "user-1",
			requesterID:     "admin-1",
			findFormResult:  &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminErr:   errors.New("db error"),
			wantErrContains: "failed to check org admin status",
		},
		{
			name:             "non-admin viewing other's history returns ErrForbidden",
			formID:           "form-1",
			targetUserID:     "user-1",
			requesterID:      "other-user",
			findFormResult:   &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult: false,
			wantErr:          ErrForbidden,
		},
		{
			name:              "org admin can view other's history",
			formID:            "form-1",
			targetUserID:      "user-1",
			requesterID:       "admin-1",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			isOrgAdminResult:  true,
			findHistoryResult: []*models.FormAnswer{{ID: "answer-1"}},
			want:              []*models.FormAnswer{{ID: "answer-1"}},
		},
		{
			name:            "find history error returns wrapped error",
			formID:          "form-1",
			targetUserID:    "user-1",
			requesterID:     "user-1",
			findFormResult:  &models.Form{ID: "form-1", OrganizationID: "org-1"},
			findHistoryErr:  errors.New("db error"),
			wantErrContains: "failed to find answer history",
		},
		{
			name:              "nil history returns empty slice",
			formID:            "form-1",
			targetUserID:      "user-1",
			requesterID:       "user-1",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			findHistoryResult: nil,
			want:              []*models.FormAnswer{},
		},
		{
			name:              "empty history returns empty slice",
			formID:            "form-1",
			targetUserID:      "user-1",
			requesterID:       "user-1",
			findFormResult:    &models.Form{ID: "form-1", OrganizationID: "org-1"},
			findHistoryResult: []*models.FormAnswer{},
			want:              []*models.FormAnswer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formRepo := &mockFormRepo{
				findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
					if tt.findFormErr != nil {
						return nil, tt.findFormErr
					}
					return tt.findFormResult, nil
				},
			}
			answerRepo := &mockFormAnswerRepo{
				findHistoryFn: func(ctx context.Context, formID, userID string) ([]*models.FormAnswer, error) {
					if tt.findHistoryErr != nil {
						return nil, tt.findHistoryErr
					}
					return tt.findHistoryResult, nil
				},
			}
			orgRepo := &mockOrgRepoForForm{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}

			svc := createTestFormService(formRepo, &mockFormFieldRepo{}, answerRepo, orgRepo, &mockUserRepoForForm{})

			got, err := svc.GetAnswerHistory(context.Background(), tt.formID, tt.targetUserID, tt.requesterID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetAnswerHistory() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("GetAnswerHistory() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("GetAnswerHistory() unexpected error = %v", err)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetAnswerHistory() got %d answers, want %d", len(got), len(tt.want))
			}
		})
	}
}

// ========== Additional Edge Case Tests ==========

// TestFormServiceImpl_AddField_FieldOrderCalculation tests the field order calculation logic
func TestFormServiceImpl_AddField_FieldOrderCalculation(t *testing.T) {
	var capturedFieldOrder int

	formRepo := &mockFormRepo{
		findByIDFn: func(ctx context.Context, formID string) (*models.Form, error) {
			return &models.Form{ID: "form-1", OrganizationID: "org-1"}, nil
		},
	}
	fieldRepo := &mockFormFieldRepo{
		getFieldsFn: func(ctx context.Context, formID string) ([]models.FormField, error) {
			// Return fields with orders 0, 5, 3 (non-sequential)
			return []models.FormField{
				{ID: "field-1", FieldOrder: 0},
				{ID: "field-2", FieldOrder: 5},
				{ID: "field-3", FieldOrder: 3},
			}, nil
		},
		addFieldFn: func(ctx context.Context, formID string, fieldType models.FormFieldType, label string, description *string, isRequired bool, fieldOrder int, options []string) (*models.FormField, error) {
			capturedFieldOrder = fieldOrder
			return &models.FormField{ID: "field-4", FieldOrder: fieldOrder}, nil
		},
	}
	orgRepo := &mockOrgRepoForForm{
		isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
			return true, nil
		},
	}

	svc := createTestFormService(formRepo, fieldRepo, &mockFormAnswerRepo{}, orgRepo, &mockUserRepoForForm{})

	_, err := svc.AddField(context.Background(), "form-1", "org-admin", &models.AddFormFieldRequest{
		FieldType: "TEXT_SMALL",
		Label:     "New Question",
	})

	if err != nil {
		t.Fatalf("AddField() unexpected error = %v", err)
	}

	// maxOrder is 5, so next order should be 6
	if capturedFieldOrder != 6 {
		t.Errorf("AddField() field order = %d, want 6", capturedFieldOrder)
	}
}

// TestFormServiceImpl_Create_ContextPropagation tests that context is properly propagated
func TestFormServiceImpl_Create_ContextPropagation(t *testing.T) {
	type contextKey string
	const testKey contextKey = "test-key"
	const testValue = "test-value"

	var capturedCtx context.Context

	formRepo := &mockFormRepo{
		createFn: func(ctx context.Context, orgID, name string, description *string, createdBy string) (*models.Form, error) {
			capturedCtx = ctx
			return &models.Form{ID: "form-1"}, nil
		},
	}
	orgRepo := &mockOrgRepoForForm{
		isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
			return true, nil
		},
	}

	svc := createTestFormService(formRepo, &mockFormFieldRepo{}, &mockFormAnswerRepo{}, orgRepo, &mockUserRepoForForm{})

	ctx := context.WithValue(context.Background(), testKey, testValue)
	_, err := svc.Create(ctx, "org-1", "org-admin", &models.CreateFormRequest{Name: "Test"})

	if err != nil {
		t.Fatalf("Create() unexpected error = %v", err)
	}

	if capturedCtx == nil {
		t.Fatal("context was not propagated")
	}

	if capturedCtx.Value(testKey) != testValue {
		t.Errorf("context value = %v, want %v", capturedCtx.Value(testKey), testValue)
	}
}
