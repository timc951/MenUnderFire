import { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import { useApi } from '../hooks/useApi';
import { GroupMembers } from '../components/groups/GroupMembers';
import { GroupMessages } from '../components/groups/GroupMessages';
import { FormFillModal } from '../components/groups/FormFillModal';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { Button } from '../components/common/Button';
import { Card } from '../components/common/Card';
import { GroupMembership, MembershipRole, GroupMessage, Form, FormAnswer, FormDetail } from '../types';
import { QRCodeSVG } from 'qrcode.react';

// Response from /groups/:id endpoint
interface GroupDetailResponse {
  id: string;
  name: string;
  description: string | null;
  organizationId: string;
  inviteCode?: string | null;
  inviteCodeExpiresAt?: string | null;
  requirePostApproval: boolean;
  allowAnonymousPosts: boolean;
  role: string;
  members: MemberResponse[];
  createdAt: string;
}

interface MemberResponse {
  id: string;
  displayName: string;
  email: string;
  role: string;
  joinedAt: string;
}

interface GroupState {
  id: string;
  name: string;
  description: string;
  organizationId: string;
  inviteCode?: string | null;
  inviteCodeExpiresAt?: string | null;
  requirePostApproval: boolean;
  allowAnonymousPosts: boolean;
  role: string;
  createdAt: string;
}

export function GroupDetail() {
  const { id } = useParams<{ id: string }>();
  const api = useApi();
  const [group, setGroup] = useState<GroupState | null>(null);
  const [members, setMembers] = useState<GroupMembership[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showInviteCode, setShowInviteCode] = useState(false);
  const [activeTab, setActiveTab] = useState<'reports' | 'members' | 'messages' | 'settings'>('messages');
  const [messages, setMessages] = useState<GroupMessage[]>([]);
  const [messagesLoading, setMessagesLoading] = useState(false);
  const [showFormSelector, setShowFormSelector] = useState(false);
  const [orgForms, setOrgForms] = useState<Form[]>([]);
  const [formsLoading, setFormsLoading] = useState(false);
  const [fillingFormId, setFillingFormId] = useState<string | null>(null);

  // Reports tab state
  const [formReports, setFormReports] = useState<GroupMessage[]>([]);
  const [formReportsLoading, setFormReportsLoading] = useState(false);
  const [selectedFormReport, setSelectedFormReport] = useState<GroupMessage | null>(null);
  const [formAnswers, setFormAnswers] = useState<FormAnswer[]>([]);
  const [formAnswersLoading, setFormAnswersLoading] = useState(false);
  const [viewingAnswer, setViewingAnswer] = useState<FormAnswer | null>(null);
  const [formDetail, setFormDetail] = useState<FormDetail | null>(null);

  const fetchGroup = useCallback(async () => {
    if (!id) return;
    setIsLoading(true);
    try {
      const data = await api.get<GroupDetailResponse>(`/groups/${id}`);
      setGroup({
        id: data.id,
        name: data.name,
        description: data.description || '',
        organizationId: data.organizationId,
        inviteCode: data.inviteCode,
        inviteCodeExpiresAt: data.inviteCodeExpiresAt,
        requirePostApproval: data.requirePostApproval,
        allowAnonymousPosts: data.allowAnonymousPosts,
        role: data.role,
        createdAt: data.createdAt,
      });
      // Convert MemberResponse to GroupMembership format
      setMembers(data.members.map((member) => ({
        id: member.id,
        userId: member.id,
        groupId: id,
        role: member.role as MembershipRole,
        joinedAt: member.joinedAt,
        displayName: member.displayName,
        email: member.email,
      })));
    } catch (err) {
      console.error('Failed to load group:', err);
      setGroup(null);
    } finally {
      setIsLoading(false);
    }
  }, [api, id]);

  const fetchMessages = useCallback(async () => {
    if (!id) return;
    setMessagesLoading(true);
    try {
      const data = await api.get<GroupMessage[]>(`/groups/${id}/messages`);
      setMessages(data);
    } catch (err) {
      console.error('Failed to load messages:', err);
    } finally {
      setMessagesLoading(false);
    }
  }, [api, id]);

  const fetchFormReports = useCallback(async () => {
    if (!id) return;
    setFormReportsLoading(true);
    try {
      const data = await api.get<GroupMessage[]>(`/groups/${id}/form-reports`);
      setFormReports(data);
    } catch (err) {
      console.error('Failed to load form reports:', err);
    } finally {
      setFormReportsLoading(false);
    }
  }, [api, id]);

  useEffect(() => {
    fetchGroup();
  }, [fetchGroup]);

  const fetchOrgForms = useCallback(async () => {
    if (!group?.organizationId) return;
    setFormsLoading(true);
    try {
      const data = await api.get<Form[]>(`/organizations/${group.organizationId}/forms`);
      setOrgForms(data);
    } catch (err) {
      console.error('Failed to load forms:', err);
    } finally {
      setFormsLoading(false);
    }
  }, [api, group?.organizationId]);

  useEffect(() => {
    if (activeTab === 'messages') {
      fetchMessages();
    } else if (activeTab === 'reports') {
      fetchFormReports();
    }
  }, [activeTab, fetchMessages, fetchFormReports]);

  // Get active forms for sending
  const activeForms = orgForms.filter(f => f.isActive);

  const handleOpenFormSelector = () => {
    setShowFormSelector(true);
    fetchOrgForms();
  };

  const handleSendForm = async (formId: string, formName: string) => {
    // Send a message with formId attached
    const content = `Please fill out the form: ${formName}`;
    await api.post(`/groups/${id}/messages`, { content, notifyMembers: true, formId });
    await fetchMessages();
    setShowFormSelector(false);
  };

  const handleRemoveMember = async (userId: string) => {
    await api.delete(`/groups/${id}/members/${userId}`);
    await fetchGroup();
  };

  const handleChangeRole = async (userId: string, newRole: string) => {
    await api.put(`/groups/${id}/members/${userId}/role`, { role: newRole });
    await fetchGroup();
  };

  const handleUpdateSettings = async (settings: { requirePostApproval: boolean; allowAnonymousPosts: boolean }) => {
    await api.put(`/groups/${id}/settings`, settings);
    await fetchGroup();
  };

  const handleSendMessage = async (content: string, notifyMembers: boolean) => {
    await api.post(`/groups/${id}/messages`, { content, notifyMembers });
    await fetchMessages();
  };

  const handleDeleteMessage = async (messageId: string) => {
    await api.delete(`/groups/${id}/messages/${messageId}`);
    await fetchMessages();
  };

  const handleViewFormReport = async (formMessage: GroupMessage) => {
    if (!formMessage.formId) return;
    setSelectedFormReport(formMessage);
    setFormAnswersLoading(true);
    try {
      const [answersData, detailData] = await Promise.all([
        api.get<FormAnswer[]>(`/groups/${id}/form-reports/${formMessage.formId}`),
        api.get<FormDetail>(`/forms/${formMessage.formId}`),
      ]);
      setFormAnswers(answersData);
      setFormDetail(detailData);
    } catch (err) {
      console.error('Failed to load form answers:', err);
    } finally {
      setFormAnswersLoading(false);
    }
  };

  const handleBackFromFormReport = () => {
    setSelectedFormReport(null);
    setFormAnswers([]);
    setFormDetail(null);
    setViewingAnswer(null);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const renderAnswerValue = (value: unknown): string => {
    if (Array.isArray(value)) {
      return value.join(', ');
    }
    if (typeof value === 'boolean') {
      return value ? 'Yes' : 'No';
    }
    return String(value ?? '');
  };

  const isLeader = group?.role === 'LEADER' || group?.role === 'OWNER' || group?.role === 'ADMIN';
  const hasAccess = group?.role && group.role !== '';

  if (isLoading) return <LoadingSpinner size="lg" />;
  if (!group) return <p className="text-red-600 dark:text-red-400">Group not found</p>;

  // User has no access to this group - show limited info only
  if (!hasAccess) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{group.name}</h1>
          {group.description && (
            <p className="text-gray-600 dark:text-stone-400">{group.description}</p>
          )}
        </div>
        <div className="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg p-4">
          <p className="text-amber-800 dark:text-amber-200 font-medium">
            You are not a member of this group.
          </p>
          <p className="text-amber-700 dark:text-amber-300 text-sm mt-1">
            Contact a group leader if you would like to join.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{group.name}</h1>
          <p className="text-gray-600 dark:text-stone-400">{group.description}</p>
        </div>
        {activeTab === 'members' && isLeader && (
          <Button onClick={() => setShowInviteCode(!showInviteCode)}>
            {showInviteCode ? 'Cancel' : 'Add Member'}
          </Button>
        )}
        {activeTab === 'messages' && isLeader && (
          <Button variant="secondary" onClick={handleOpenFormSelector}>
            Send Form
          </Button>
        )}
      </div>

      {showInviteCode && activeTab === 'members' && group.inviteCode && (
        <div className="bg-stone-50 dark:bg-stone-700/50 rounded-lg p-4 space-y-3">
          <p className="text-sm text-gray-700 dark:text-stone-300">
            Share this invite link, code, or QR code to join the group:
          </p>
          <div className="flex flex-col sm:flex-row items-start gap-4">
            <div className="space-y-3">
              <div>
                <p className="text-xs font-medium text-gray-500 dark:text-stone-400 mb-1">Direct Link</p>
                <div className="flex items-center gap-2">
                  <code className="text-sm font-mono text-amber-600 dark:text-amber-400 bg-stone-100 dark:bg-stone-800 px-3 py-2 rounded break-all">
                    {`${window.location.origin}/groups/join?code=${group.inviteCode}`}
                  </code>
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => {
                      navigator.clipboard.writeText(`${window.location.origin}/groups/join?code=${group.inviteCode}`);
                    }}
                  >
                    Copy
                  </Button>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <code className="text-lg font-mono font-bold text-amber-600 dark:text-amber-400 bg-stone-100 dark:bg-stone-800 px-3 py-2 rounded">
                  {group.inviteCode}
                </code>
                <Button
                  variant="secondary"
                  onClick={() => {
                    navigator.clipboard.writeText(`${group.inviteCode}`);
                  }}
                >
                  Copy Code
                </Button>
              </div>
              {group.inviteCodeExpiresAt && (
                <p className={`text-xs ${
                  new Date(group.inviteCodeExpiresAt) < new Date()
                    ? 'text-red-600 dark:text-red-400'
                    : 'text-gray-500 dark:text-stone-400'
                }`}>
                  {new Date(group.inviteCodeExpiresAt) < new Date()
                    ? 'Invite code expired on '
                    : 'Expires '}
                  {new Date(group.inviteCodeExpiresAt).toLocaleDateString()}
                </p>
              )}
            </div>
            <div className="bg-white p-3 rounded-lg shadow-sm">
              <QRCodeSVG
                value={`${window.location.origin}/groups/join?code=${group.inviteCode}`}
                size={150}
                level="M"
                marginSize={0}
              />
              <p className="text-xs text-center text-gray-500 mt-1">Scan to join</p>
            </div>
          </div>
        </div>
      )}

      <div className="border-b border-gray-200 dark:border-stone-700">
        <nav className="flex gap-4">
          {isLeader && (
            <button
              onClick={() => setActiveTab('reports')}
              className={`py-2 px-1 border-b-2 text-sm font-medium ${
                activeTab === 'reports'
                  ? 'border-amber-500 text-amber-600 dark:text-amber-400'
                  : 'border-transparent text-gray-500 dark:text-stone-400 hover:text-gray-700 dark:hover:text-stone-300'
              }`}
            >
              Reports
            </button>
          )}
          <button
            onClick={() => setActiveTab('members')}
            className={`py-2 px-1 border-b-2 text-sm font-medium ${
              activeTab === 'members'
                ? 'border-amber-500 text-amber-600 dark:text-amber-400'
                : 'border-transparent text-gray-500 dark:text-stone-400 hover:text-gray-700 dark:hover:text-stone-300'
            }`}
          >
            Members
          </button>
          <button
            onClick={() => setActiveTab('messages')}
            className={`py-2 px-1 border-b-2 text-sm font-medium ${
              activeTab === 'messages'
                ? 'border-amber-500 text-amber-600 dark:text-amber-400'
                : 'border-transparent text-gray-500 dark:text-stone-400 hover:text-gray-700 dark:hover:text-stone-300'
            }`}
          >
            Messages
          </button>
          {isLeader && (
            <button
              onClick={() => setActiveTab('settings')}
              className={`py-2 px-1 border-b-2 text-sm font-medium ${
                activeTab === 'settings'
                  ? 'border-amber-500 text-amber-600 dark:text-amber-400'
                  : 'border-transparent text-gray-500 dark:text-stone-400 hover:text-gray-700 dark:hover:text-stone-300'
              }`}
            >
              Settings
            </button>
          )}
        </nav>
      </div>

      {activeTab === 'reports' && isLeader && (
        <div>
          {formReportsLoading ? (
            <LoadingSpinner size="lg" />
          ) : !selectedFormReport ? (
            // Form reports list
            formReports.length === 0 ? (
              <p className="text-gray-500 dark:text-stone-400 text-center py-8">
                No form reports yet. Send a form from the Messages tab to collect responses.
              </p>
            ) : (
              <div className="space-y-3">
                {formReports.map((msg) => (
                  <Card
                    key={msg.id}
                    className="cursor-pointer hover:bg-stone-50 dark:hover:bg-stone-700/50 transition-colors"
                    onClick={() => handleViewFormReport(msg)}
                  >
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="font-medium text-gray-900 dark:text-white">
                          {msg.formName || 'Form'}
                        </h3>
                        <p className="text-sm text-gray-500 dark:text-stone-400">
                          Sent {formatDate(msg.createdAt)}
                        </p>
                      </div>
                      <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                      </svg>
                    </div>
                  </Card>
                ))}
              </div>
            )
          ) : (
            // Form answers view
            <div className="space-y-4">
              <div className="flex items-center gap-4">
                <button
                  onClick={handleBackFromFormReport}
                  className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                  aria-label="Back"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                  </svg>
                </button>
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                  {selectedFormReport.formName || 'Form'} - Responses
                </h3>
              </div>

              {formAnswersLoading ? (
                <LoadingSpinner size="md" />
              ) : formAnswers.length === 0 ? (
                <p className="text-gray-500 dark:text-stone-400 text-center py-8">
                  No responses yet for this form.
                </p>
              ) : (
                <div className="space-y-3">
                  {formAnswers.map((answer) => (
                    <Card
                      key={answer.id}
                      className="cursor-pointer hover:bg-stone-50 dark:hover:bg-stone-700/50 transition-colors"
                      onClick={() => setViewingAnswer(answer)}
                    >
                      <div className="flex items-center justify-between">
                        <div>
                          <h4 className="font-medium text-gray-900 dark:text-white">
                            {answer.userName || 'Anonymous'}
                          </h4>
                          <p className="text-sm text-gray-500 dark:text-stone-400">
                            Submitted {formatDate(answer.submittedAt)}
                          </p>
                          {answer.version > 1 && (
                            <span className="text-xs text-amber-600 dark:text-amber-400">
                              Version {answer.version}
                            </span>
                          )}
                        </div>
                        <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                      </div>
                    </Card>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {activeTab === 'members' && (
        <GroupMembers
          members={members}
          inviteCode={group.inviteCode}
          currentUserRole={group.role || 'MEMBER'}
          onRemoveMember={handleRemoveMember}
          onChangeRole={handleChangeRole}
        />
      )}

      {activeTab === 'messages' && (
        <GroupMessages
          messages={messages}
          isLoading={messagesLoading}
          isLeader={isLeader}
          onSendMessage={handleSendMessage}
          onDeleteMessage={handleDeleteMessage}
          onFillForm={(formId) => setFillingFormId(formId)}
        />
      )}

      {activeTab === 'settings' && isLeader && group && (
        <Card>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Group Settings</h3>
          <div className="space-y-4">
            <label className="flex items-center justify-between cursor-pointer">
              <div>
                <p className="text-sm font-medium text-gray-900 dark:text-white">Require post approval</p>
                <p className="text-xs text-gray-500 dark:text-stone-400">Leaders must approve messages before they are visible to the group</p>
              </div>
              <button
                type="button"
                role="switch"
                aria-checked={group.requirePostApproval}
                onClick={() => handleUpdateSettings({ requirePostApproval: !group.requirePostApproval, allowAnonymousPosts: group.allowAnonymousPosts })}
                className={`relative inline-flex h-6 w-11 flex-shrink-0 rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2 ${
                  group.requirePostApproval ? 'bg-amber-500' : 'bg-gray-200 dark:bg-stone-600'
                }`}
              >
                <span className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                  group.requirePostApproval ? 'translate-x-5' : 'translate-x-0'
                }`} />
              </button>
            </label>
            <label className="flex items-center justify-between cursor-pointer">
              <div>
                <p className="text-sm font-medium text-gray-900 dark:text-white">Allow anonymous posts</p>
                <p className="text-xs text-gray-500 dark:text-stone-400">Members can post messages without showing their name</p>
              </div>
              <button
                type="button"
                role="switch"
                aria-checked={group.allowAnonymousPosts}
                onClick={() => handleUpdateSettings({ requirePostApproval: group.requirePostApproval, allowAnonymousPosts: !group.allowAnonymousPosts })}
                className={`relative inline-flex h-6 w-11 flex-shrink-0 rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2 ${
                  group.allowAnonymousPosts ? 'bg-amber-500' : 'bg-gray-200 dark:bg-stone-600'
                }`}
              >
                <span className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                  group.allowAnonymousPosts ? 'translate-x-5' : 'translate-x-0'
                }`} />
              </button>
            </label>
          </div>
        </Card>
      )}

      {/* Send Form Modal */}
      {showFormSelector && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-stone-800 rounded-lg p-6 w-full max-w-md mx-4">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Send Form to Group</h2>
            {formsLoading ? (
              <div className="py-8">
                <LoadingSpinner size="md" />
              </div>
            ) : activeForms.length === 0 ? (
              <p className="text-gray-500 dark:text-stone-400 text-center py-8">
                No active forms available. Create a form in your organization first.
              </p>
            ) : (
              <div className="space-y-2 max-h-64 overflow-y-auto">
                {activeForms.map((form) => (
                  <button
                    key={form.id}
                    onClick={() => handleSendForm(form.id, form.name)}
                    className="w-full text-left p-3 rounded-lg border border-gray-200 dark:border-stone-700 hover:bg-stone-50 dark:hover:bg-stone-700/50 transition-colors"
                  >
                    <h3 className="font-medium text-gray-900 dark:text-white">{form.name}</h3>
                    {form.description && (
                      <p className="text-sm text-gray-600 dark:text-stone-400 mt-1">{form.description}</p>
                    )}
                    <p className="text-xs text-gray-500 dark:text-stone-500 mt-1">{form.fieldCount} fields</p>
                  </button>
                ))}
              </div>
            )}
            <div className="flex justify-end mt-6">
              <Button variant="secondary" onClick={() => setShowFormSelector(false)}>
                Cancel
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Form Fill Modal */}
      {fillingFormId && (
        <FormFillModal
          formId={fillingFormId}
          onClose={() => setFillingFormId(null)}
          onSubmitted={() => fetchMessages()}
        />
      )}

      {/* View Answer Modal */}
      {viewingAnswer && formDetail && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 overflow-y-auto">
          <div className="bg-white dark:bg-stone-800 rounded-lg p-6 w-full max-w-2xl mx-4 my-8 max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  {viewingAnswer.userName || 'Anonymous'}
                </h2>
                <p className="text-sm text-gray-500 dark:text-stone-400">
                  Submitted {formatDate(viewingAnswer.submittedAt)}
                </p>
              </div>
              <button
                onClick={() => setViewingAnswer(null)}
                className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div className="space-y-4">
              {formDetail.fields.map((field) => {
                const value = viewingAnswer.answers[field.id];
                if (field.fieldType === 'TEXT_DISPLAY') return null;
                return (
                  <div key={field.id} className="border-b border-gray-200 dark:border-stone-700 pb-3 last:border-0">
                    <p className="text-sm font-medium text-gray-700 dark:text-stone-300">
                      {field.label}
                    </p>
                    <p className="text-gray-900 dark:text-white mt-1">
                      {value !== undefined ? renderAnswerValue(value) : <span className="text-gray-400 dark:text-stone-500 italic">No answer</span>}
                    </p>
                  </div>
                );
              })}
            </div>

            <div className="flex justify-end mt-6">
              <Button variant="secondary" onClick={() => setViewingAnswer(null)}>
                Close
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
