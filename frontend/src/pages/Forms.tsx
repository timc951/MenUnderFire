import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useApi } from '../hooks/useApi';
import { Form } from '../types';
import { Button } from '../components/common/Button';
import { Card } from '../components/common/Card';
import { LoadingSpinner } from '../components/common/LoadingSpinner';

export function Forms() {
  const { orgId } = useParams<{ orgId: string }>();
  const navigate = useNavigate();
  const api = useApi();
  const [forms, setForms] = useState<Form[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newFormName, setNewFormName] = useState('');
  const [newFormDescription, setNewFormDescription] = useState('');
  const [isCreating, setIsCreating] = useState(false);

  const fetchForms = useCallback(async () => {
    if (!orgId) return;
    setIsLoading(true);
    try {
      const data = await api.get<Form[]>(`/organizations/${orgId}/forms`);
      setForms(data);
    } catch (err) {
      console.error('Failed to load forms:', err);
    } finally {
      setIsLoading(false);
    }
  }, [api, orgId]);

  useEffect(() => {
    fetchForms();
  }, [fetchForms]);

  const handleCreateForm = async () => {
    if (!newFormName.trim() || !orgId) return;
    setIsCreating(true);
    try {
      const form = await api.post<Form>(`/organizations/${orgId}/forms`, {
        name: newFormName.trim(),
        description: newFormDescription.trim() || undefined,
      });
      setShowCreateModal(false);
      setNewFormName('');
      setNewFormDescription('');
      navigate(`/forms/${form.id}`);
    } catch (err) {
      console.error('Failed to create form:', err);
    } finally {
      setIsCreating(false);
    }
  };

  if (isLoading) {
    return <LoadingSpinner size="lg" />;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <button
            onClick={() => navigate(`/organizations/${orgId}`)}
            className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            aria-label="Back"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Forms</h1>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>Create Form</Button>
      </div>

      {forms.length === 0 ? (
        <Card>
          <p className="text-gray-500 dark:text-stone-400 text-center py-8">
            No forms yet. Create your first form to get started.
          </p>
        </Card>
      ) : (
        <div className="grid gap-4">
          {forms.map((form) => (
            <Card
              key={form.id}
              className="cursor-pointer hover:bg-stone-50 dark:hover:bg-stone-700/50 transition-colors"
              onClick={() => navigate(`/forms/${form.id}`)}
            >
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium text-gray-900 dark:text-white">{form.name}</h3>
                  {form.description && (
                    <p className="text-sm text-gray-600 dark:text-stone-400 mt-1">{form.description}</p>
                  )}
                  <div className="flex items-center gap-4 mt-2 text-xs text-gray-500 dark:text-stone-500">
                    <span>{form.fieldCount} fields</span>
                    <span>{form.isActive ? 'Active' : 'Inactive'}</span>
                  </div>
                </div>
                <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* Create Form Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-stone-800 rounded-lg p-6 w-full max-w-md mx-4">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Create New Form</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">
                  Form Name
                </label>
                <input
                  type="text"
                  value={newFormName}
                  onChange={(e) => setNewFormName(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md bg-white dark:bg-stone-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  placeholder="Enter form name"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">
                  Description (optional)
                </label>
                <textarea
                  value={newFormDescription}
                  onChange={(e) => setNewFormDescription(e.target.value)}
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md bg-white dark:bg-stone-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  placeholder="Enter description"
                />
              </div>
            </div>
            <div className="flex justify-end gap-3 mt-6">
              <Button variant="secondary" onClick={() => setShowCreateModal(false)}>
                Cancel
              </Button>
              <Button onClick={handleCreateForm} disabled={isCreating || !newFormName.trim()}>
                {isCreating ? 'Creating...' : 'Create'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
