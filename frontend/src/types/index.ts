export interface User {
  id: string;
  externalId: string;
  email: string;
  displayName: string;
  agreementAcceptedAt?: string | null;
  createdAt: string;
}

export interface Group {
  id: string;
  name: string;
  description: string;
  inviteCode?: string | null;
  createdAt: string;
  memberCount?: number;
  role?: string;
}

export type MembershipRole = 'OWNER' | 'LEADER' | 'MODERATOR' | 'MEMBER';

export interface GroupMembership {
  id: string;
  userId: string;
  groupId: string;
  role: MembershipRole;
  joinedAt: string;
  displayName: string;
  email?: string;
}

export interface Report {
  id: string;
  title: string;
  content: string;
  groupId: string;
  reporterName: string | null;
  isAnonymous: boolean;
  createdAt: string;
}

export interface ApiKey {
  id: string;
  name: string;
  keyPrefix: string;
  permissions: ApiKeyPermissions;
  expiresAt: string | null;
  createdAt: string;
  lastUsedAt: string | null;
}

export interface ApiKeyPermissions {
  reports: ('read' | 'write')[];
  groups: ('read')[];
}

export interface CreateGroupRequest {
  name: string;
  description?: string;
  organizationId: string;
}

export interface JoinGroupRequest {
  inviteCode: string;
}

export interface CreateReportRequest {
  groupId: string;
  title: string;
  content: string;
  isAnonymousToGroup: boolean;
}

export interface CreateApiKeyRequest {
  name: string;
  permissions: ApiKeyPermissions;
  expiresInDays?: number;
}

export interface CreateApiKeyResponse {
  id: string;
  name: string;
  key: string;
  permissions: ApiKeyPermissions;
  expiresAt: string | null;
}

export interface DashboardStats {
  organizationCount: number | null;
  groupCount: number;
}

// Site Pages
export interface SitePageSummary {
  id: string;
  slug: string;
  title: string;
  isPublished: boolean;
  updatedAt: string;
}

export interface SitePage {
  id: string;
  slug: string;
  title: string;
  content: string;
  isPublished: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateSitePageRequest {
  slug: string;
  title: string;
  content: string;
  isPublished: boolean;
}

export interface UpdateSitePageRequest {
  title: string;
  content: string;
  isPublished: boolean;
}

// Group Messages
export interface GroupMessage {
  id: string;
  groupId: string;
  senderId: string;
  senderName: string;
  content: string;
  notifyMembers: boolean;
  formId?: string;
  formName?: string;
  createdAt: string;
}

export interface SendGroupMessageRequest {
  content: string;
  notifyMembers: boolean;
  formId?: string;
}

// Forms
export type FormFieldType =
  | 'TEXT_DISPLAY'
  | 'TEXT_SMALL'
  | 'TEXT_MEDIUM'
  | 'TEXT_LARGE'
  | 'CHECKBOX'
  | 'RADIO'
  | 'DROPDOWN';

export interface Form {
  id: string;
  organizationId: string;
  name: string;
  description: string | null;
  isActive: boolean;
  fieldCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface FormField {
  id: string;
  fieldType: FormFieldType;
  label: string;
  description: string | null;
  isRequired: boolean;
  fieldOrder: number;
  options: string[] | null;
}

export interface FormDetail {
  id: string;
  organizationId: string;
  name: string;
  description: string | null;
  isActive: boolean;
  fields: FormField[];
  createdAt: string;
  updatedAt: string;
}

export interface FormAnswer {
  id: string;
  formId: string;
  userId: string;
  userName: string | null;
  version: number;
  isCurrent: boolean;
  answers: Record<string, unknown>;
  submittedAt: string;
}

export interface CreateFormRequest {
  name: string;
  description?: string;
}

export interface UpdateFormRequest {
  name: string;
  description?: string;
  isActive: boolean;
}

export interface AddFormFieldRequest {
  fieldType: FormFieldType;
  label: string;
  description?: string;
  isRequired: boolean;
  options?: string[];
}

export interface UpdateFormFieldRequest {
  label: string;
  description?: string;
  isRequired: boolean;
  options?: string[];
}

// Invitations
export interface InvitationListItem {
  id: string;
  email: string;
  type: 'ORG_ADMIN' | 'GROUP_OWNER' | 'GROUP_MEMBER';
  organizationId?: string | null;
  organizationName?: string | null;
  groupId?: string | null;
  groupName?: string | null;
  inviterName?: string | null;
  status: 'PENDING' | 'ACCEPTED' | 'EXPIRED' | 'REVOKED';
  expiresAt: string;
  createdAt: string;
}

// Feedback
export interface Feedback {
  id: string;
  userId: string;
  userDisplayName: string;
  userEmail: string;
  type: 'BUG' | 'FEATURE' | 'OTHER';
  subject: string;
  description: string;
  status: 'OPEN' | 'IN_PROGRESS' | 'RESOLVED' | 'CLOSED' | 'ROADMAP' | 'MORE_INFO_REQUIRED';
  createdAt: string;
  updatedAt: string;
}

export interface CreateFeedbackRequest {
  type: 'BUG' | 'FEATURE' | 'OTHER';
  subject: string;
  description: string;
}
