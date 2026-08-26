import { useState, FormEvent } from 'react';
import { CreateReportRequest } from '../../types';
import { Button } from '../common/Button';
import { Input } from '../common/Input';
import { AnonymityToggle } from './AnonymityToggle';

interface GroupOption {
  id: string;
  name: string;
}

interface CreateReportFormProps {
  groups: GroupOption[];
  onSubmit?: (data: CreateReportRequest) => Promise<unknown>;
}

export function CreateReportForm({ groups, onSubmit }: CreateReportFormProps) {
  const [groupId, setGroupId] = useState('');
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [isAnonymousToGroup, setIsAnonymousToGroup] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [successMessage, setSuccessMessage] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};
    if (!title.trim()) newErrors.title = 'Title is required';
    if (!content.trim()) newErrors.content = 'Content is required';
    if (!groupId) newErrors.groupId = 'Group is required';
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setSuccessMessage('');
    setErrorMessage('');

    if (!validate()) return;

    setIsSubmitting(true);
    try {
      await onSubmit?.({ groupId, title, content, isAnonymousToGroup });
      setSuccessMessage('Report submitted successfully!');
      setTitle('');
      setContent('');
      setIsAnonymousToGroup(false);
    } catch {
      setErrorMessage('Failed to submit report. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-1">
        <label htmlFor="group-select" className="block text-sm font-medium text-gray-700 dark:text-stone-300">
          Group
        </label>
        <select
          id="group-select"
          value={groupId}
          onChange={(e) => setGroupId(e.target.value)}
          className="block w-full rounded-md border border-gray-300 dark:border-stone-600 bg-white dark:bg-stone-700 text-gray-900 dark:text-white px-3 py-2 shadow-sm focus:border-amber-500 focus:ring-1 focus:ring-amber-500"
          aria-invalid={!!errors.groupId}
        >
          <option value="">Select a group</option>
          {groups.map((group) => (
            <option key={group.id} value={group.id}>
              {group.name}
            </option>
          ))}
        </select>
        {errors.groupId && <p className="text-sm text-red-600 dark:text-red-400">{errors.groupId}</p>}
      </div>

      <Input
        label="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        error={errors.title}
        placeholder="Report title"
      />

      <div className="space-y-1">
        <label htmlFor="content" className="block text-sm font-medium text-gray-700 dark:text-stone-300">
          Content
        </label>
        <textarea
          id="content"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          rows={4}
          className={`block w-full rounded-md border ${errors.content ? 'border-red-500' : 'border-gray-300 dark:border-stone-600'} bg-white dark:bg-stone-700 text-gray-900 dark:text-white px-3 py-2 shadow-sm focus:border-amber-500 focus:ring-1 focus:ring-amber-500`}
          placeholder="Share your progress..."
          aria-invalid={!!errors.content}
        />
        {errors.content && <p className="text-sm text-red-600 dark:text-red-400" role="alert">{errors.content}</p>}
      </div>

      <AnonymityToggle value={isAnonymousToGroup} onChange={setIsAnonymousToGroup} />

      {successMessage && <p className="text-sm text-green-600">{successMessage}</p>}
      {errorMessage && <p className="text-sm text-red-600">{errorMessage}</p>}

      <Button type="submit" isLoading={isSubmitting}>
        {isSubmitting ? 'Submitting...' : 'Submit Report'}
      </Button>
    </form>
  );
}
