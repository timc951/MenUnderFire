import { useState, useEffect, useCallback } from 'react';
import { useApi } from '../hooks/useApi';
import { useRbac } from '../rbac';
import { Button } from '../components/common/Button';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { Feedback as FeedbackType, CreateFeedbackRequest } from '../types';

const STATUS_OPTIONS = [
  { value: 'OPEN', label: 'Open' },
  { value: 'IN_PROGRESS', label: 'In Progress' },
  { value: 'ROADMAP', label: 'Added to Roadmap' },
  { value: 'MORE_INFO_REQUIRED', label: 'More Info Required' },
  { value: 'RESOLVED', label: 'Resolved' },
  { value: 'CLOSED', label: 'Closed' },
];

export function Feedback() {
  const api = useApi();
  const rbac = useRbac();
  const isSiteAdmin = rbac.isSiteAdmin();

  const [feedbackList, setFeedbackList] = useState<FeedbackType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [type, setType] = useState<CreateFeedbackRequest['type']>('BUG');
  const [subject, setSubject] = useState('');
  const [description, setDescription] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [updatingStatus, setUpdatingStatus] = useState<string | null>(null);

  const fetchFeedback = useCallback(async () => {
    setIsLoading(true);
    try {
      const endpoint = isSiteAdmin ? '/feedback' : '/feedback/me';
      const data = await api.get<FeedbackType[]>(endpoint);
      setFeedbackList(data);
    } catch (err) {
      console.error('Failed to load feedback:', err);
    } finally {
      setIsLoading(false);
    }
  }, [api, isSiteAdmin]);

  useEffect(() => {
    fetchFeedback();
  }, [fetchFeedback]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      await api.post('/feedback', { type, subject, description });
      setShowForm(false);
      setSubject('');
      setDescription('');
      setType('BUG');
      await fetchFeedback();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to submit feedback');
    } finally {
      setSubmitting(false);
    }
  };

  const handleStatusChange = async (feedbackId: string, newStatus: string) => {
    setUpdatingStatus(feedbackId);
    try {
      await api.patch(`/feedback/${feedbackId}/status`, { status: newStatus });
      setFeedbackList((prev) =>
        prev.map((item) =>
          item.id === feedbackId ? { ...item, status: newStatus as FeedbackType['status'] } : item
        )
      );
    } catch (err) {
      console.error('Failed to update status:', err);
    } finally {
      setUpdatingStatus(null);
    }
  };

  const statusColor = (status: string) => {
    switch (status) {
      case 'OPEN': return 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400';
      case 'IN_PROGRESS': return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400';
      case 'ROADMAP': return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400';
      case 'MORE_INFO_REQUIRED': return 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400';
      case 'RESOLVED': return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400';
      case 'CLOSED': return 'bg-stone-100 text-stone-600 dark:bg-stone-700 dark:text-stone-400';
      default: return 'bg-stone-100 text-stone-600 dark:bg-stone-700 dark:text-stone-400';
    }
  };

  const statusLabel = (status: string) => {
    const opt = STATUS_OPTIONS.find((o) => o.value === status);
    return opt ? opt.label : status.replace(/_/g, ' ');
  };

  const typeLabel = (t: string) => {
    switch (t) {
      case 'BUG': return 'Bug Report';
      case 'FEATURE': return 'Feature Request';
      case 'OTHER': return 'Other';
      default: return t;
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Feedback</h1>
          <p className="text-gray-600 dark:text-stone-400">
            {isSiteAdmin ? 'Manage all user feedback' : 'Report issues or suggest features'}
          </p>
        </div>
        <Button onClick={() => setShowForm(!showForm)}>
          {showForm ? 'Cancel' : 'New Feedback'}
        </Button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="bg-white dark:bg-stone-800 rounded-lg border border-stone-200 dark:border-stone-700 p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">Type</label>
            <select
              value={type}
              onChange={(e) => setType(e.target.value as CreateFeedbackRequest['type'])}
              className="w-full rounded-md border border-stone-300 dark:border-stone-600 bg-white dark:bg-stone-700 text-gray-900 dark:text-white px-3 py-2 focus:outline-none focus:ring-2 focus:ring-amber-500"
            >
              <option value="BUG">Bug Report</option>
              <option value="FEATURE">Feature Request</option>
              <option value="OTHER">Other</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">Subject</label>
            <input
              type="text"
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
              required
              maxLength={200}
              placeholder="Brief summary of your feedback"
              className="w-full rounded-md border border-stone-300 dark:border-stone-600 bg-white dark:bg-stone-700 text-gray-900 dark:text-white px-3 py-2 focus:outline-none focus:ring-2 focus:ring-amber-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              required
              rows={5}
              placeholder="Describe the issue or feature in detail..."
              className="w-full rounded-md border border-stone-300 dark:border-stone-600 bg-white dark:bg-stone-700 text-gray-900 dark:text-white px-3 py-2 focus:outline-none focus:ring-2 focus:ring-amber-500"
            />
          </div>
          {error && <p className="text-red-600 dark:text-red-400 text-sm">{error}</p>}
          <div className="flex justify-end">
            <Button type="submit" disabled={submitting}>
              {submitting ? 'Submitting...' : 'Submit Feedback'}
            </Button>
          </div>
        </form>
      )}

      {isLoading ? (
        <LoadingSpinner size="lg" />
      ) : feedbackList.length === 0 ? (
        <div className="text-center py-12 text-gray-500 dark:text-stone-400">
          <p>No feedback submitted yet.</p>
          <p className="text-sm mt-1">Click "New Feedback" to report an issue or suggest a feature.</p>
        </div>
      ) : (
        <div className="space-y-3">
          {feedbackList.map((item) => (
            <div key={item.id} className="bg-white dark:bg-stone-800 rounded-lg border border-stone-200 dark:border-stone-700 p-4">
              <div className="flex items-start justify-between gap-3">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 flex-wrap">
                    <h3 className="font-medium text-gray-900 dark:text-white">{item.subject}</h3>
                    {isSiteAdmin ? (
                      <select
                        value={item.status}
                        onChange={(e) => handleStatusChange(item.id, e.target.value)}
                        disabled={updatingStatus === item.id}
                        className={`text-xs font-medium rounded-full px-2 py-0.5 border-0 cursor-pointer focus:outline-none focus:ring-2 focus:ring-amber-500 ${statusColor(item.status)} ${updatingStatus === item.id ? 'opacity-50' : ''}`}
                      >
                        {STATUS_OPTIONS.map((opt) => (
                          <option key={opt.value} value={opt.value} className="bg-white dark:bg-stone-800 text-gray-900 dark:text-white">{opt.label}</option>
                        ))}
                      </select>
                    ) : (
                      <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusColor(item.status)}`}>
                        {statusLabel(item.status)}
                      </span>
                    )}
                  </div>
                  <p className="text-xs text-gray-500 dark:text-stone-500 mt-1">
                    {typeLabel(item.type)}
                    {isSiteAdmin && item.userDisplayName && (
                      <> &middot; {item.userDisplayName} ({item.userEmail})</>
                    )}
                    {' '}&middot; {new Date(item.createdAt).toLocaleDateString()}
                  </p>
                  <p className="text-sm text-gray-600 dark:text-stone-400 mt-2 whitespace-pre-wrap">{item.description}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
