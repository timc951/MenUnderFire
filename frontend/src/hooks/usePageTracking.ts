import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { useAuth } from './useAuth';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';
const HIT_TOKEN = import.meta.env.VITE_HIT_TOKEN || '';

export function usePageTracking() {
  const location = useLocation();
  const { user } = useAuth();

  useEffect(() => {
    // Fire-and-forget hit tracking
    fetch(`${API_BASE_URL}/hits`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Hit-Token': HIT_TOKEN,
      },
      body: JSON.stringify({
        path: location.pathname,
        email: user?.email || '',
      }),
    }).catch(() => {
      // Silently ignore tracking failures
    });
  }, [location.pathname, user?.email]);
}
