import { useState, FormEvent } from 'react';
import { Modal } from '../common/Modal';
import { Input } from '../common/Input';
import { Button } from '../common/Button';

interface JoinGroupModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (inviteCode: string) => Promise<void>;
}

export function JoinGroupModal({ isOpen, onClose, onSubmit }: JoinGroupModalProps) {
  const [inviteCode, setInviteCode] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!inviteCode.trim()) {
      setError('Invite code is required');
      return;
    }
    setIsSubmitting(true);
    setError('');
    try {
      await onSubmit(inviteCode.trim());
      setInviteCode('');
      onClose();
    } catch {
      setError('Invalid invite code or group not found');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Join Group">
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Invite Code"
          value={inviteCode}
          onChange={(e) => setInviteCode(e.target.value)}
          error={error}
          placeholder="Enter invite code"
        />
        <div className="flex gap-3 justify-end">
          <Button variant="secondary" type="button" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" isLoading={isSubmitting}>
            Join Group
          </Button>
        </div>
      </form>
    </Modal>
  );
}
