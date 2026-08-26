import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSitePages } from '../hooks/useSitePages';
import { useRbac } from '../rbac';
import { Button } from '../components/common/Button';
import { Card } from '../components/common/Card';
import { RichTextEditor } from '../components/common/RichTextEditor';
import { SitePageSummary, CreateSitePageRequest } from '../types';

interface CreatePageModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: CreateSitePageRequest) => Promise<void>;
}

function CreatePageModal({ isOpen, onClose, onSubmit }: CreatePageModalProps) {
  const [slug, setSlug] = useState('');
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [isPublished, setIsPublished] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!slug.trim() || !title.trim()) return;

    setIsSubmitting(true);
    setError(null);

    try {
      await onSubmit({ slug: slug.trim(), title: title.trim(), content, isPublished });
      setSlug('');
      setTitle('');
      setContent('');
      setIsPublished(true);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create page');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white dark:bg-stone-800 rounded-lg p-6 w-full max-w-3xl max-h-[90vh] overflow-y-auto">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          Create New Page
        </h2>

        {error && (
          <div className="mb-4 p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-md text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="slug" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              URL Slug
            </label>
            <input
              type="text"
              id="slug"
              value={slug}
              onChange={(e) => setSlug(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, '-'))}
              placeholder="about-us"
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500 bg-white dark:bg-stone-700 text-gray-900 dark:text-white"
              required
            />
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              The URL will be: /page/{slug || 'your-slug'}
            </p>
          </div>

          <div>
            <label htmlFor="title" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Title
            </label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="About Us"
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500 bg-white dark:bg-stone-700 text-gray-900 dark:text-white"
              required
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
              Publish immediately
            </label>
          </div>

          <div className="flex justify-end gap-3 pt-4">
            <Button variant="secondary" onClick={onClose} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting || !slug.trim() || !title.trim()}>
              {isSubmitting ? 'Creating...' : 'Create Page'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}

function PageCard({ page, onClick }: { page: SitePageSummary; onClick: () => void }) {
  return (
    <Card className="cursor-pointer hover:shadow-md transition-shadow" onClick={onClick}>
      <div className="flex items-start justify-between">
        <div className="flex-1 min-w-0">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white truncate">
            {page.title}
          </h3>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            /{page.slug}
          </p>
        </div>
        <div className="flex items-center gap-2 ml-4">
          {page.isPublished ? (
            <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400">
              Published
            </span>
          ) : (
            <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400">
              Draft
            </span>
          )}
        </div>
      </div>
      <p className="text-xs text-gray-400 dark:text-gray-500 mt-2">
        Updated {new Date(page.updatedAt).toLocaleDateString()}
      </p>
    </Card>
  );
}

export function SitePages() {
  const navigate = useNavigate();
  const rbac = useRbac();
  const { pages, isLoading, error, createPage } = useSitePages();
  const [showCreateModal, setShowCreateModal] = useState(false);

  const isSiteAdmin = rbac.isSiteAdmin();

  if (!isSiteAdmin) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
            Access Denied
          </h2>
          <p className="text-gray-600 dark:text-gray-400">
            Only Site Admins can manage site pages.
          </p>
        </div>
      </div>
    );
  }

  const handlePageClick = (page: SitePageSummary) => {
    navigate(`/admin/pages/${page.id}`);
  };

  const handleCreatePage = async (data: CreateSitePageRequest) => {
    await createPage(data);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Site Pages</h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
            Manage content pages visible to all users
          </p>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>
          Create Page
        </Button>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500" />
        </div>
      ) : error ? (
        <div className="text-center py-12">
          <p className="text-red-600 dark:text-red-400">{error}</p>
        </div>
      ) : pages.length === 0 ? (
        <div className="text-center py-12">
          <svg
            className="mx-auto h-12 w-12 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            />
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900 dark:text-white">
            No pages yet
          </h3>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Create your first site page to get started.
          </p>
          <div className="mt-6">
            <Button onClick={() => setShowCreateModal(true)}>
              Create Page
            </Button>
          </div>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {pages.map((page) => (
            <PageCard
              key={page.id}
              page={page}
              onClick={() => handlePageClick(page)}
            />
          ))}
        </div>
      )}

      <CreatePageModal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        onSubmit={handleCreatePage}
      />
    </div>
  );
}
