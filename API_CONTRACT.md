# Men Under Fire - API Contract

Base URL: `/api`

All routes require JWT authentication unless marked as **Public**.

---

## Table of Contents

1. [User Routes](#user-routes)
2. [Group Routes](#group-routes)
3. [Group Message Routes](#group-message-routes)
4. [Report Routes](#report-routes)
5. [Organization Routes](#organization-routes)
6. [Invitation Routes](#invitation-routes)
7. [Form Routes](#form-routes)
8. [Site Page Routes](#site-page-routes)
9. [Dashboard Routes](#dashboard-routes)
10. [API Key Routes](#api-key-routes)

---

## User Routes

### GET /api/users/me
Get the current authenticated user's profile.

**Response:** `UserResponse`
```json
{
  "id": "string (UUID)",
  "email": "string",
  "displayName": "string",
  "createdAt": "string (ISO 8601)"
}
```

---

### GET /api/users/me/permissions
Get the current user's permissions and roles.

**Response:** `UserPermissionsResponse`
```json
{
  "isSiteAdmin": "boolean",
  "adminOfOrganizationIds": ["string (UUID)"],
  "ownedGroupIds": ["string (UUID)"],
  "memberGroupIds": ["string (UUID)"]
}
```

---

### PUT /api/users/me
Update the current user's profile.

**Request Body:** `UpdateUserRequest`
```json
{
  "displayName": "string"
}
```

**Response:** `UserResponse`

---

## Group Routes

### POST /api/groups
Create a new accountability group (caller becomes owner).

**Request Body:** `CreateGroupRequest`
```json
{
  "name": "string",
  "description": "string | null",
  "organizationId": "string (UUID)"
}
```

**Response:** `GroupResponse` (201 Created)
```json
{
  "id": "string (UUID)",
  "name": "string",
  "description": "string | null",
  "inviteCode": "string | null",
  "role": "string | null",
  "memberCount": "number | null",
  "createdAt": "string (ISO 8601)"
}
```

---

### GET /api/groups
List all groups the current user belongs to.

**Response:** `GroupResponse[]`

---

### GET /api/groups/{id}
Get detailed information about a specific group.

**Path Parameters:**
- `id` - Group UUID

**Response:** `GroupDetailResponse`
```json
{
  "id": "string (UUID)",
  "name": "string",
  "description": "string | null",
  "organizationId": "string (UUID)",
  "inviteCode": "string | null",
  "role": "string",
  "members": [
    {
      "id": "string (UUID)",
      "displayName": "string",
      "email": "string",
      "role": "string",
      "joinedAt": "string (ISO 8601)"
    }
  ],
  "createdAt": "string (ISO 8601)"
}
```

**Notes:**
- `inviteCode` is only visible to group leaders/owners

---

### POST /api/groups/{id}/join
Join a group using an invite code.

**Path Parameters:**
- `id` - Group UUID

**Request Body:** `JoinGroupRequest`
```json
{
  "inviteCode": "string"
}
```

**Response:** `MessageResponse`
```json
{
  "message": "string"
}
```

---

### DELETE /api/groups/{id}/members/{userId}
Remove a member from a group.

**Path Parameters:**
- `id` - Group UUID
- `userId` - User UUID to remove

**Response:** 204 No Content

---

## Group Message Routes

### GET /api/groups/{groupId}/messages
List all messages for a group.

**Path Parameters:**
- `groupId` - Group UUID

**Response:** `GroupMessageResponse[]`
```json
[
  {
    "id": "string (UUID)",
    "groupId": "string (UUID)",
    "senderId": "string (UUID)",
    "senderName": "string",
    "content": "string",
    "notifyMembers": "boolean",
    "createdAt": "string (ISO 8601)"
  }
]
```

---

### POST /api/groups/{groupId}/messages
Send a message to all group members.

**Path Parameters:**
- `groupId` - Group UUID

**Request Body:** `SendGroupMessageRequest`
```json
{
  "content": "string",
  "notifyMembers": "boolean (default: false)"
}
```

**Response:** `GroupMessageResponse` (201 Created)

**Authorization:** Group Owner/Leader or Admin only

---

### DELETE /api/groups/{groupId}/messages/{messageId}
Delete a group message.

**Path Parameters:**
- `groupId` - Group UUID
- `messageId` - Message UUID

**Response:** 204 No Content

---

## Report Routes

### POST /api/reports
Submit an accountability report to a group.

**Request Body:** `CreateReportRequest`
```json
{
  "groupId": "string (UUID)",
  "title": "string",
  "content": "string",
  "isAnonymousToGroup": "boolean (default: false)"
}
```

**Response:** `ReportResponse` (201 Created)
```json
{
  "id": "string (UUID)",
  "groupId": "string (UUID)",
  "title": "string",
  "content": "string",
  "isAnonymousToGroup": "boolean",
  "reporterName": "string | null",
  "reporterId": "string | null",
  "isOwnReport": "boolean",
  "createdAt": "string (ISO 8601)"
}
```

**Notes:**
- `reporterName` and `reporterId` are null for anonymous reports when viewed by non-leaders

---

### GET /api/reports
List reports for a group.

**Query Parameters:**
- `groupId` (required) - Group UUID

**Response:** `ReportResponse[]`

---

## Organization Routes

### GET /api/organizations
List organizations the current user belongs to.

**Response:** `OrganizationResponse[]`
```json
[
  {
    "id": "string (UUID)",
    "name": "string",
    "description": "string | null",
    "createdAt": "string (ISO 8601)"
  }
]
```

---

### GET /api/organizations/all
List all organizations in the system.

**Authorization:** Site Admin only

**Response:** `OrganizationResponse[]`

---

### POST /api/organizations
Create a new organization.

**Authorization:** Site Admin only

**Request Body:** `CreateOrganizationRequest`
```json
{
  "name": "string",
  "description": "string | null"
}
```

**Response:** `OrganizationResponse` (201 Created)

---

### GET /api/organizations/{id}
Get detailed information about an organization.

**Path Parameters:**
- `id` - Organization UUID

**Response:** `OrganizationDetailResponse`
```json
{
  "id": "string (UUID)",
  "name": "string",
  "description": "string | null",
  "createdById": "string (UUID)",
  "canEdit": "boolean",
  "isSiteAdmin": "boolean",
  "createdAt": "string (ISO 8601)",
  "updatedAt": "string (ISO 8601)"
}
```

---

### PUT /api/organizations/{id}
Update an organization.

**Authorization:** Admin only (Site Admin or Org Admin)

**Path Parameters:**
- `id` - Organization UUID

**Request Body:** `UpdateOrganizationRequest`
```json
{
  "name": "string",
  "description": "string | null"
}
```

**Response:** `OrganizationResponse`

---

### GET /api/organizations/{id}/admins
List administrators of an organization.

**Path Parameters:**
- `id` - Organization UUID

**Response:** `OrganizationAdminResponse[]`
```json
[
  {
    "id": "string (UUID)",
    "userId": "string (UUID)",
    "organizationId": "string (UUID)",
    "displayName": "string",
    "email": "string",
    "joinedAt": "string (ISO 8601)"
  }
]
```

---

### GET /api/organizations/{id}/groups
List all groups within an organization.

**Path Parameters:**
- `id` - Organization UUID

**Response:** `GroupResponse[]`

---

## Invitation Routes

### POST /api/invitations/org-admin
Invite a new organization admin.

**Authorization:** Site Admin only

**Request Body:** `CreateOrgAdminInvitationRequest`
```json
{
  "email": "string",
  "organizationId": "string (UUID)"
}
```

**Response:** `InvitationResponse` (201 Created)
```json
{
  "id": "string (UUID)",
  "token": "string",
  "email": "string",
  "type": "string",
  "organizationId": "string | null",
  "groupId": "string | null",
  "status": "string",
  "expiresAt": "string (ISO 8601)",
  "createdAt": "string (ISO 8601)"
}
```

---

### POST /api/invitations/group-owner
Invite a new group owner.

**Authorization:** Site Admin or Org Admin

**Request Body:** `CreateGroupOwnerInvitationRequest`
```json
{
  "email": "string",
  "groupId": "string (UUID)"
}
```

**Response:** `InvitationResponse` (201 Created)

---

### POST /api/invitations/group-member
Invite a new group member.

**Authorization:** Site Admin, Org Admin, or Group Owner

**Request Body:** `CreateGroupMemberInvitationRequest`
```json
{
  "email": "string",
  "groupId": "string (UUID)"
}
```

**Response:** `InvitationResponse` (201 Created)

---

### DELETE /api/invitations/{id}
Cancel/delete an invitation.

**Path Parameters:**
- `id` - Invitation UUID

**Response:** 204 No Content

---

### POST /api/invitations/accept
Accept an invitation. **Public** - No authentication required.

**Request Body:** `AcceptInvitationRequest`
```json
{
  "token": "string",
  "externalId": "string",
  "displayName": "string"
}
```

**Response:** `AcceptInvitationResponse`
```json
{
  "invitation": { /* InvitationResponse */ },
  "user": { /* UserResponse */ }
}
```

---

### GET /api/invitations/validate/{token}
Validate an invitation token. **Public** - No authentication required.

**Path Parameters:**
- `token` - Invitation token string

**Response:** `ValidateInvitationResponse`
```json
{
  "valid": "boolean",
  "email": "string | null",
  "type": "string | null",
  "organizationName": "string | null",
  "groupName": "string | null",
  "inviterName": "string | null",
  "expiresAt": "string | null",
  "capabilities": ["string"] | null
}
```

---

## Form Routes

### GET /api/organizations/{orgId}/forms
List all forms for an organization.

**Path Parameters:**
- `orgId` - Organization UUID

**Response:** `FormResponse[]`
```json
[
  {
    "id": "string (UUID)",
    "organizationId": "string (UUID)",
    "name": "string",
    "description": "string | null",
    "isActive": "boolean",
    "fieldCount": "number",
    "createdAt": "string (ISO 8601)",
    "updatedAt": "string (ISO 8601)"
  }
]
```

---

### POST /api/organizations/{orgId}/forms
Create a new form.

**Authorization:** Org Admin only

**Path Parameters:**
- `orgId` - Organization UUID

**Request Body:** `CreateFormRequest`
```json
{
  "name": "string",
  "description": "string | null"
}
```

**Response:** `FormResponse` (201 Created)

---

### GET /api/forms/{id}
Get detailed form information with all fields.

**Path Parameters:**
- `id` - Form UUID

**Response:** `FormDetailResponse`
```json
{
  "id": "string (UUID)",
  "organizationId": "string (UUID)",
  "name": "string",
  "description": "string | null",
  "isActive": "boolean",
  "fields": [
    {
      "id": "string (UUID)",
      "fieldType": "string",
      "label": "string",
      "description": "string | null",
      "isRequired": "boolean",
      "fieldOrder": "number",
      "options": ["string"] | null
    }
  ],
  "createdAt": "string (ISO 8601)",
  "updatedAt": "string (ISO 8601)"
}
```

**Field Types:**
- `TEXT_DISPLAY` - Display-only text
- `TEXT_SMALL` - Single-line text input
- `TEXT_MEDIUM` - Multi-line text input
- `TEXT_LARGE` - Large text area
- `CHECKBOX` - Multiple choice checkboxes
- `RADIO` - Single choice radio buttons
- `DROPDOWN` - Single choice dropdown

---

### PUT /api/forms/{id}
Update a form.

**Path Parameters:**
- `id` - Form UUID

**Request Body:** `UpdateFormRequest`
```json
{
  "name": "string",
  "description": "string | null",
  "isActive": "boolean (default: true)"
}
```

**Response:** `FormResponse`

---

### DELETE /api/forms/{id}
Delete a form.

**Path Parameters:**
- `id` - Form UUID

**Response:** 204 No Content

---

### POST /api/forms/{id}/fields
Add a field to a form.

**Path Parameters:**
- `id` - Form UUID

**Request Body:** `AddFormFieldRequest`
```json
{
  "fieldType": "string",
  "label": "string",
  "description": "string | null",
  "isRequired": "boolean (default: false)",
  "options": ["string"] | null
}
```

**Response:** `FormFieldResponse` (201 Created)

---

### PUT /api/forms/{id}/fields/reorder
Reorder fields in a form.

**Path Parameters:**
- `id` - Form UUID

**Request Body:** `ReorderFieldsRequest`
```json
{
  "fieldIds": ["string (UUID)"]
}
```

**Response:** `MessageResponse`

---

### PUT /api/forms/{formId}/fields/{fieldId}
Update a form field.

**Path Parameters:**
- `formId` - Form UUID
- `fieldId` - Field UUID

**Request Body:** `UpdateFormFieldRequest`
```json
{
  "label": "string",
  "description": "string | null",
  "isRequired": "boolean (default: false)",
  "options": ["string"] | null
}
```

**Response:** `FormFieldResponse`

---

### DELETE /api/forms/{formId}/fields/{fieldId}
Delete a form field.

**Path Parameters:**
- `formId` - Form UUID
- `fieldId` - Field UUID

**Response:** 204 No Content

---

### POST /api/forms/{id}/answers
Submit answers to a form.

**Path Parameters:**
- `id` - Form UUID

**Request Body:** `SubmitFormAnswerRequest`
```json
{
  "answers": {
    "field_id": "value (JSON element)"
  }
}
```

**Response:** `FormAnswerResponse` (201 Created)
```json
{
  "id": "string (UUID)",
  "formId": "string (UUID)",
  "userId": "string (UUID)",
  "userName": "string | null",
  "version": "number",
  "isCurrent": "boolean",
  "answers": { "field_id": "value" },
  "submittedAt": "string (ISO 8601)"
}
```

---

### GET /api/forms/{id}/answers
List all answers for a form.

**Path Parameters:**
- `id` - Form UUID

**Response:** `FormAnswerResponse[]`

---

### GET /api/forms/{id}/answers/me
Get the current user's latest answer for a form.

**Path Parameters:**
- `id` - Form UUID

**Response:** `FormAnswerResponse` or 404 Not Found

---

### GET /api/forms/{id}/answers/history/{userId}
Get answer history for a specific user on a form.

**Path Parameters:**
- `id` - Form UUID
- `userId` - User UUID

**Response:** `FormAnswerResponse[]`

---

## Site Page Routes

### GET /api/pages
List all site pages. **Public** - No authentication required.

**Response:** `SitePageSummaryResponse[]`
```json
[
  {
    "id": "string (UUID)",
    "slug": "string",
    "title": "string",
    "isPublished": "boolean",
    "updatedAt": "string (ISO 8601)"
  }
]
```

---

### GET /api/pages/{slug}
Get a site page by slug. **Public** - No authentication required.

**Path Parameters:**
- `slug` - Page slug string

**Response:** `SitePageResponse`
```json
{
  "id": "string (UUID)",
  "slug": "string",
  "title": "string",
  "content": "string (markdown)",
  "isPublished": "boolean",
  "createdAt": "string (ISO 8601)",
  "updatedAt": "string (ISO 8601)"
}
```

---

### POST /api/pages
Create a new site page.

**Authorization:** Site Admin only

**Request Body:** `CreateSitePageRequest`
```json
{
  "slug": "string",
  "title": "string",
  "content": "string (default: '')",
  "isPublished": "boolean (default: true)"
}
```

**Response:** `SitePageResponse` (201 Created)

---

### PUT /api/pages/{id}
Update a site page.

**Authorization:** Site Admin only

**Path Parameters:**
- `id` - Page UUID

**Request Body:** `UpdateSitePageRequest`
```json
{
  "title": "string",
  "content": "string",
  "isPublished": "boolean (default: true)"
}
```

**Response:** `SitePageResponse`

---

### DELETE /api/pages/{id}
Delete a site page.

**Authorization:** Site Admin only

**Path Parameters:**
- `id` - Page UUID

**Response:** 204 No Content

---

## Dashboard Routes

### GET /api/dashboard/stats
Get dashboard statistics.

**Authorization:** Admin only (Site Admin or Org Admin)

**Response:** `DashboardStatsResponse`
```json
{
  "organizationCount": "number | null",
  "groupCount": "number"
}
```

**Notes:**
- `organizationCount` is only visible to Site Admins
- Org Admins see group count within their organizations
- Regular users see minimal stats

---

## API Key Routes

> **Note:** All API key operations are currently disabled (403 Forbidden).

### POST /api/api-keys
Create a new API key. **Currently disabled.**

**Request Body:** `CreateApiKeyRequest`
```json
{
  "name": "string",
  "permissions": ["string"] (default: []),
  "expiresInDays": "number (default: 90)"
}
```

---

### GET /api/api-keys
List all API keys. **Currently disabled.**

---

### DELETE /api/api-keys/{id}
Delete an API key. **Currently disabled.**

**Path Parameters:**
- `id` - API Key UUID

---

## Common Response Types

### MessageResponse
```json
{
  "message": "string"
}
```

### Error Response
```json
{
  "error": "string",
  "message": "string"
}
```

---

## Authentication

All protected routes require a valid JWT token in the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

---

## CORS Configuration

- Allowed Methods: `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`
- Rate Limiting: Enabled (configurable)
