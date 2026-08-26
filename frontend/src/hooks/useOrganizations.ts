import { useState, useEffect, useCallback } from 'react';
import { useApi } from './useApi';

export interface Organization {
  id: string;
  name: string;
  description: string | null;
  createdAt: string;
}

export interface OrganizationDetail {
  id: string;
  name: string;
  description: string | null;
  createdById: string;
  canEdit: boolean;
  isSiteAdmin: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateOrganizationData {
  name: string;
  description?: string;
}

export interface UpdateOrganizationData {
  name: string;
  description?: string;
}

export function useOrganizations() {
  const api = useApi();
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchOrganizations = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await api.get<Organization[]>('/organizations');
      setOrganizations(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load organizations');
    } finally {
      setIsLoading(false);
    }
  }, [api]);

  useEffect(() => {
    fetchOrganizations();
  }, [fetchOrganizations]);

  const createOrganization = useCallback(
    async (data: CreateOrganizationData) => {
      const newOrg = await api.post<Organization>('/organizations', data);
      setOrganizations((prev) => [...prev, newOrg]);
      return newOrg;
    },
    [api]
  );

  const refetch = useCallback(() => {
    fetchOrganizations();
  }, [fetchOrganizations]);

  return {
    organizations,
    isLoading,
    error,
    createOrganization,
    refetch,
  };
}

export function useOrganization(id: string) {
  const api = useApi();
  const [organization, setOrganization] = useState<OrganizationDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchOrganization = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await api.get<OrganizationDetail>(`/organizations/${id}`);
      setOrganization(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load organization');
    } finally {
      setIsLoading(false);
    }
  }, [api, id]);

  useEffect(() => {
    if (id) {
      fetchOrganization();
    }
  }, [id, fetchOrganization]);

  const updateOrganization = useCallback(
    async (data: UpdateOrganizationData) => {
      const updated = await api.put<Organization>(`/organizations/${id}`, data);
      setOrganization((prev) =>
        prev ? { ...prev, name: updated.name, description: updated.description } : null
      );
      return updated;
    },
    [api, id]
  );

  const refetch = useCallback(() => {
    fetchOrganization();
  }, [fetchOrganization]);

  return {
    organization,
    isLoading,
    error,
    updateOrganization,
    refetch,
  };
}
