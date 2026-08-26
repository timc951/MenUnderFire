import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useApi } from '../hooks/useApi';
import { useRbac } from '../rbac';
import { Button } from '../components/common/Button';
import { Card } from '../components/common/Card';
import { RichTextEditor } from '../components/common/RichTextEditor';
import { SitePage, UpdateSitePageRequest } from '../types';

export function SitePageEdit() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const rbac = useRbac();
  const api = useApi();

  const [page, setPage] = useState<SitePage | null>(null);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [isPublished, setIsPublished] = useState(true);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [isPreviewing, setIsPreviewing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [saveError, setSaveError] = useState<string | null>(null);

  const isSiteAdmin = rbac.isSiteAdmin();

  useEffect(() => {
    const fetchPage = async () => {
      if (!id) return;
      setIsLoading(true);
      setError(null);
      try {
        // Fetch by ID - we need to get the page details
        const pages = await api.get<SitePage[]>('/pages');
        const foundPage = pages.find(p => p.id === id);
        if (foundPage) {
          // Now fetch the full page content by slug
          const fullPage = await api.get<SitePage>(`/pages/${foundPage.slug}`);
          setPage(fullPage);
          setTitle(fullPage.title);
          setContent(fullPage.content);
          setIsPublished(fullPage.isPublished);
        } else {
          setError('Page not found');
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load page');
      } finally {
        setIsLoading(false);
      }
    };
    fetchPage();
  }, [api, id]);

  const handleSave = async () => {
    if (!id || !title.trim()) return;

    setIsSaving(true);
    setSaveError(null);

    try {
      const data: UpdateSitePageRequest = {
        title: title.trim(),
        content,
        isPublished,
      };
      await api.put(`/pages/${id}`, data);
      navigate('/admin/pages');
    } catch (err) {
      setSaveError(err instanceof Error ? err.message : 'Failed to save page');
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!id) return;
    if (!confirm('Are you sure you want to delete this page? This action cannot be undone.')) {
      return;
    }

    setIsDeleting(true);
    try {
      await api.delete(`/pages/${id}`);
      navigate('/admin/pages');
    } catch (err) {
      setSaveError(err instanceof Error ? err.message : 'Failed to delete page');
      setIsDeleting(false);
    }
  };

  const handlePreview = async () => {
    if (!title.trim()) {
      setSaveError('Title is required to preview');
      return;
    }

    setIsPreviewing(true);
    setSaveError(null);

    try {
      const response = await api.post<{ id: string }>(`/pages/${id}/preview`, {
        title: title.trim(),
        content,
      });
      // Open preview in new tab
      window.open(`/page/preview/${response.id}`, '_blank');
    } catch (err) {
      setSaveError(err instanceof Error ? err.message : 'Failed to create preview');
    } finally {
      setIsPreviewing(false);
    }
  };

  if (!isSiteAdmin) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
            Access Denied
          </h2>
          <p className="text-gray-600 dark:text-gray-400">
            Only Site Admins can edit site pages.
          </p>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500" />
      </div>
    );
  }

  if (error || !page) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600 dark:text-red-400">{error || 'Page not found'}</p>
        <Button variant="secondary" className="mt-4" onClick={() => navigate('/admin/pages')}>
          Back to Pages
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <button
            onClick={() => navigate('/admin/pages')}
            className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            aria-label="Back"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              Edit Page
            </h1>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              /{page.slug}
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button
            variant="secondary"
            onClick={handlePreview}
            disabled={isPreviewing}
          >
            {isPreviewing ? 'Loading...' : 'Preview'}
          </Button>
          <Button
            variant="secondary"
            onClick={handleDelete}
            disabled={isDeleting}
            className="text-red-600 hover:text-red-700 border-red-300 hover:border-red-400"
          >
            {isDeleting ? 'Deleting...' : 'Delete'}
          </Button>
        </div>
      </div>

      {saveError && (
        <div className="p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-md text-sm">
          {saveError}
        </div>
      )}

      <Card>
        <div className="space-y-4">
          <div>
            <label htmlFor="title" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Title
            </label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500 bg-white dark:bg-stone-700 text-gray-900 dark:text-white"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Content
            </label>
            <RichTextEditor
              content={content}
              onChange={setContent}
              placeholder="Start writing your page content..."
            />
          </div>

          <div className="flex items-center">
            <input
              type="checkbox"
              id="isPublished"
              checked={isPublished}
              onChange={(e) => setIsPublished(e.target.checked)}
              className="h-4 w-4 text-amber-500 focus:ring-amber-500 border-gray-300 rounded accent-amber-500"
            />
            <label htmlFor="isPublished" className="ml-2 text-sm text-gray-700 dark:text-gray-300">
              Published (visible to all users)
            </label>
          </div>

          <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
            <Button variant="secondary" onClick={() => navigate('/admin/pages')} disabled={isSaving}>
              Cancel
            </Button>
            <Button onClick={handleSave} disabled={isSaving || !title.trim()}>
              {isSaving ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
}
