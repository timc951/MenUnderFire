import { useState, useEffect, useCallback } from 'react';
import { useApi } from '../hooks/useApi';
import { Button } from '../components/common/Button';
import { Modal } from '../components/common/Modal';
import { Input } from '../components/common/Input';
import { Card } from '../components/common/Card';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ApiKey, CreateApiKeyResponse } from '../types';
import { formatDate } from '../utils/helpers';

export function ApiKeys() {
  const api = useApi();
  const [keys, setKeys] = useState<ApiKey[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showConfirmDelete, setShowConfirmDelete] = useState<string | null>(null);
  const [newKeyName, setNewKeyName] = useState('');
  const [newKeyValue, setNewKeyValue] = useState<string | null>(null);
  const [toast, setToast] = useState<string | null>(null);

  const fetchKeys = useCallback(async () => {
    setIsLoading(true);
    try {
      const data = await api.get<ApiKey[]>('/api-keys');
      setKeys(data);
    } catch {
      // handled by loading state
    } finally {
      setIsLoading(false);
    }
  }, [api]);

  useEffect(() => {
    fetchKeys();
  }, [fetchKeys]);

  const handleCreate = async () => {
    const result = await api.post<CreateApiKeyResponse>('/api-keys', {
      name: newKeyName,
      permissions: { reports: ['read', 'write'], groups: ['read'] },
    });
    setNewKeyValue(result.key);
    setNewKeyName('');
    await fetchKeys();
  };

  const handleRevoke = async (keyId: string) => {
    await api.delete(`/api-keys/${keyId}`);
    setShowConfirmDelete(null);
    setToast('Key revoked successfully');
    await fetchKeys();
    setTimeout(() => setToast(null), 3000);
  };

  if (isLoading) return <LoadingSpinner />;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">API Keys</h1>
        <Button onClick={() => setShowCreateModal(true)}>
          Generate New Key
        </Button>
      </div>

      {keys.length === 0 ? (
        <p className="text-gray-500">No API keys yet. Generate one to get started.</p>
      ) : (
        <div className="space-y-3">
          {keys.map((key) => (
            <Card key={key.id} className="flex items-center justify-between">
              <div>
                <p className="font-medium">{key.name}</p>
                <p className="text-sm text-gray-500">
                  {key.keyPrefix}... • Created {formatDate(key.createdAt)}
                  {key.expiresAt && ` • Expires ${formatDate(key.expiresAt)}`}
                </p>
              </div>
              <Button variant="danger" size="sm" onClick={() => setShowConfirmDelete(key.id)}>
                Revoke
              </Button>
            </Card>
          ))}
        </div>
      )}

      <Modal isOpen={showCreateModal} onClose={() => { setShowCreateModal(false); setNewKeyValue(null); }} title="Generate API Key">
        {newKeyValue ? (
          <div className="space-y-4">
            <p className="text-sm text-green-600 font-medium">Copy this key now - you won't be able to see it again!</p>
            <code className="block bg-gray-100 p-3 rounded text-sm break-all">{newKeyValue}</code>
            <Button onClick={() => { navigator.clipboard.writeText(newKeyValue); setToast('Copied!'); setTimeout(() => setToast(null), 2000); }}>
              Copy to Clipboard
            </Button>
          </div>
        ) : (
          <div className="space-y-4">
            <Input label="Key Name" value={newKeyName} onChange={(e) => setNewKeyName(e.target.value)} placeholder="e.g., My Integration Key" />
            <Button onClick={handleCreate} disabled={!newKeyName.trim()}>
              Create Key
            </Button>
          </div>
        )}
      </Modal>

      {showConfirmDelete && (
        <Modal isOpen={true} onClose={() => setShowConfirmDelete(null)} title="Revoke API Key">
          <div className="space-y-4">
            <p>Are you sure you want to revoke this key? This cannot be undone.</p>
            <div className="flex gap-3 justify-end">
              <Button variant="secondary" onClick={() => setShowConfirmDelete(null)}>Cancel</Button>
              <Button variant="danger" onClick={() => handleRevoke(showConfirmDelete)}>Confirm Revoke</Button>
            </div>
          </div>
        </Modal>
      )}

      {toast && (
        <div className="fixed bottom-4 right-4 bg-green-50 border-l-4 border-green-400 p-4 rounded shadow-lg" role="alert">
          <p className="text-green-800">{toast}</p>
        </div>
      )}
    </div>
  );
}
