import { useState, useEffect, useCallback } from 'react';
import { useApi } from './useApi';
import { SitePageSummary, SitePage, CreateSitePageRequest, UpdateSitePageRequest } from '../types';

export function useSitePages() {
  const api = useApi();
  const [pages, setPages] = useState<SitePageSummary[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchPages = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await api.get<SitePageSummary[]>('/pages');
      setPages(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load pages');
    } finally {
      setIsLoading(false);
    }
  }, [api]);

  useEffect(() => {
    fetchPages();
  }, [fetchPages]);

  const createPage = async (data: CreateSitePageRequest): Promise<SitePage> => {
    const page = await api.post<SitePage>('/pages', data);
    await fetchPages();
    return page;
  };

  const updatePage = async (id: string, data: UpdateSitePageRequest): Promise<SitePage> => {
    const page = await api.put<SitePage>(`/pages/${id}`, data);
    await fetchPages();
    return page;
  };

  const deletePage = async (id: string): Promise<void> => {
    await api.delete(`/pages/${id}`);
    await fetchPages();
  };

  return { pages, isLoading, error, createPage, updatePage, deletePage, refetch: fetchPages };
}

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

// Public hook for fetching published pages (no auth required)
export function usePublicSitePages() {
  const [pages, setPages] = useState<SitePageSummary[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchPages = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_BASE_URL}/pages`);
      if (!response.ok) {
        throw new Error('Failed to load pages');
      }
      const data = await response.json();
      setPages(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load pages');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchPages();
  }, [fetchPages]);

  return { pages, isLoading, error, refetch: fetchPages };
}

export function useSitePage(slug: string) {
  const [page, setPage] = useState<SitePage | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchPage = useCallback(async () => {
    if (!slug) return;
    setIsLoading(true);
    setError(null);
    try {
      // Public pages don't require authentication
      const response = await fetch(`${API_BASE_URL}/pages/${slug}`);
      if (!response.ok) {
        throw new Error(response.status === 404 ? 'Page not found' : 'Failed to load page');
      }
      const data = await response.json();
      setPage(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load page');
    } finally {
      setIsLoading(false);
    }
  }, [slug]);

  useEffect(() => {
    fetchPage();
  }, [fetchPage]);

  return { page, isLoading, error, refetch: fetchPage };
}
