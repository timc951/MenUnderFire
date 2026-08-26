import { http, HttpResponse } from 'msw';

const API_BASE = 'http://localhost:7001/api';

const mockGroups = [
  { id: '1', name: 'Fitness Group', description: 'Stay fit together', inviteCode: 'ABC123', createdBy: 'user-1', createdAt: '2024-01-01T00:00:00Z', memberCount: 5, role: 'leader' },
  { id: '2', name: 'Study Group', description: 'Study accountability', inviteCode: 'DEF456', createdBy: 'user-2', createdAt: '2024-01-02T00:00:00Z', memberCount: 3, role: 'member' },
];

const mockMembers = [
  { id: 'm1', userId: 'user-1', groupId: '1', role: 'leader', joinedAt: '2024-01-01T00:00:00Z', displayName: 'Leader User' },
  { id: 'm2', userId: 'user-2', groupId: '1', role: 'member', joinedAt: '2024-01-02T00:00:00Z', displayName: 'Member One' },
  { id: 'm3', userId: 'user-3', groupId: '1', role: 'member', joinedAt: '2024-01-03T00:00:00Z', displayName: 'Member Two' },
];

const mockReports = [
  { id: 'r1', title: 'Weekly Check-in', content: 'I exercised 4 times this week.', groupId: '1', reporterName: 'John Doe', isAnonymous: false, createdAt: '2024-01-15T10:00:00Z' },
  { id: 'r2', title: 'Progress Update', content: 'Completed all my tasks.', groupId: '1', reporterName: null, isAnonymous: true, createdAt: '2024-01-14T09:00:00Z' },
];

const mockApiKeys = [
  { id: 'key1', name: 'My Integration Key', keyPrefix: 'ak_abc1', permissions: { reports: ['read', 'write'], groups: ['read'] }, expiresAt: '2025-01-01T00:00:00Z', createdAt: '2024-01-01T00:00:00Z', lastUsedAt: '2024-01-10T00:00:00Z' },
];

const mockUser = {
  id: 'user-1',
  externalId: 'oidc|user-1',
  email: 'test@example.com',
  displayName: 'Test User',
  createdAt: '2024-01-01T00:00:00Z',
};

// Default permissions (can be overridden in tests)
const mockPermissions = {
  isSiteAdmin: false,
  adminOfOrganizationIds: [],
  ownedGroupIds: ['1'],
  memberGroupIds: ['2'],
};

const mockOrganizations = [
  { id: 'org-1', name: 'Test Organization', createdAt: '2024-01-01T00:00:00Z' },
  { id: 'org-2', name: 'Another Organization', createdAt: '2024-01-02T00:00:00Z' },
];

const mockInvitations = [
  { id: 'inv-1', email: 'new@example.com', type: 'group_member', groupId: '1', status: 'pending', expiresAt: '2025-01-01T00:00:00Z', createdAt: '2024-01-01T00:00:00Z' },
];

export const handlers = [
  http.get(`${API_BASE}/users/me`, () => {
    return HttpResponse.json(mockUser);
  }),

  http.get(`${API_BASE}/groups`, () => {
    return HttpResponse.json(mockGroups);
  }),

  http.get(`${API_BASE}/groups/:id`, ({ params }) => {
    const group = mockGroups.find(g => g.id === params.id);
    if (!group) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json(group);
  }),

  http.post(`${API_BASE}/groups`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>;
    return HttpResponse.json(
      { id: 'new-group-1', ...body, inviteCode: 'NEW123', createdAt: new Date().toISOString(), memberCount: 1, role: 'leader' },
      { status: 201 }
    );
  }),

  http.post(`${API_BASE}/groups/join`, () => {
    return HttpResponse.json(mockGroups[0]);
  }),

  http.get(`${API_BASE}/groups/:id/members`, () => {
    return HttpResponse.json(mockMembers);
  }),

  http.delete(`${API_BASE}/groups/:groupId/members/:userId`, () => {
    return new HttpResponse(null, { status: 204 });
  }),

  http.get(`${API_BASE}/groups/:groupId/reports`, () => {
    return HttpResponse.json(mockReports);
  }),

  http.post(`${API_BASE}/reports`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>;
    return HttpResponse.json(
      { id: 'new-report-1', ...body, createdAt: new Date().toISOString() },
      { status: 201 }
    );
  }),

  http.get(`${API_BASE}/api-keys`, () => {
    return HttpResponse.json(mockApiKeys);
  }),

  http.post(`${API_BASE}/api-keys`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>;
    return HttpResponse.json({
      id: 'new-key-1',
      name: body.name,
      key: 'ak_test1234567890abcdef',
      permissions: body.permissions,
      expiresAt: '2025-06-01T00:00:00Z',
    }, { status: 201 });
  }),

  http.delete(`${API_BASE}/api-keys/:id`, () => {
    return new HttpResponse(null, { status: 204 });
  }),

  // RBAC endpoints
  http.get(`${API_BASE}/users/me/permissions`, () => {
    return HttpResponse.json(mockPermissions);
  }),

  http.get(`${API_BASE}/organizations`, () => {
    return HttpResponse.json(mockOrganizations);
  }),

  // Invitation endpoints
  http.get(`${API_BASE}/invitations/validate/:token`, ({ params }) => {
    if (params.token === 'invalid-token') {
      return HttpResponse.json({ valid: false });
    }
    return HttpResponse.json({
      valid: true,
      email: 'test@example.com',
      type: 'group_member',
      organizationName: null,
      groupName: 'Fitness Group',
      inviterName: 'Admin User',
      expiresAt: '2025-01-01T00:00:00Z',
      capabilities: [
        'Submit accountability reports',
        'View group reports',
        'Participate in group activities',
      ],
    });
  }),

  http.post(`${API_BASE}/invitations/accept`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>;
    return HttpResponse.json({
      invitation: { id: 'inv-1', email: body.email, status: 'accepted' },
      user: { id: 'user-new', email: body.email, displayName: body.displayName, createdAt: new Date().toISOString() },
    });
  }),

  http.post(`${API_BASE}/invitations/org-admin`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>;
    return HttpResponse.json({
      id: 'inv-new',
      email: body.email,
      type: 'org_admin',
      organizationId: body.organizationId,
      status: 'pending',
      expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
      createdAt: new Date().toISOString(),
    }, { status: 201 });
  }),

  http.post(`${API_BASE}/invitations/group-owner`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>;
    return HttpResponse.json({
      id: 'inv-new',
      email: body.email,
      type: 'group_owner',
      groupId: body.groupId,
      status: 'pending',
      expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
      createdAt: new Date().toISOString(),
    }, { status: 201 });
  }),

  http.post(`${API_BASE}/invitations/group-member`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>;
    return HttpResponse.json({
      id: 'inv-new',
      email: body.email,
      type: 'group_member',
      groupId: body.groupId,
      status: 'pending',
      expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
      createdAt: new Date().toISOString(),
    }, { status: 201 });
  }),

  http.delete(`${API_BASE}/invitations/:id`, () => {
    return new HttpResponse(null, { status: 204 });
  }),
];
