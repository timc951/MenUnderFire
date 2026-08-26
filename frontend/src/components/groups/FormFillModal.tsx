import { useState, useEffect } from 'react';
import { useApi } from '../../hooks/useApi';
import { FormDetail, FormAnswer } from '../../types';
import { LoadingSpinner } from '../common/LoadingSpinner';
import { Button } from '../common/Button';

interface FormFillModalProps {
  formId: string;
  onClose: () => void;
  onSubmitted: () => void;
}

export function FormFillModal({ formId, onClose, onSubmitted }: FormFillModalProps) {
  const api = useApi();
  const [form, setForm] = useState<FormDetail | null>(null);
  const [existingAnswer, setExistingAnswer] = useState<FormAnswer | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [answers, setAnswers] = useState<Record<string, unknown>>({});
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true);
      try {
        const [formData, myAnswer] = await Promise.all([
          api.get<FormDetail>(`/forms/${formId}`),
          api.get<FormAnswer>(`/forms/${formId}/answers/me`).catch(() => null),
        ]);
        setForm(formData);
        setExistingAnswer(myAnswer);
        if (myAnswer?.answers) {
          setAnswers(myAnswer.answers);
        }
      } catch (err) {
        console.error('Failed to load form:', err);
        setError('Failed to load form');
      } finally {
        setIsLoading(false);
      }
    };
    fetchData();
  }, [api, formId]);

  const handleFieldChange = (fieldId: string, value: unknown) => {
    setAnswers((prev) => ({ ...prev, [fieldId]: value }));
  };

  const handleCheckboxChange = (fieldId: string, option: string, checked: boolean) => {
    const current = (answers[fieldId] as string[]) || [];
    const updated = checked
      ? [...current, option]
      : current.filter((v) => v !== option);
    handleFieldChange(fieldId, updated);
  };

  const handleSubmit = async () => {
    if (!form) return;
    setIsSubmitting(true);
    setError(null);
    try {
      await api.post(`/forms/${formId}/answers`, { answers });
      onSubmitted();
      onClose();
    } catch (err) {
      console.error('Failed to submit form:', err);
      setError('Failed to submit form. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 overflow-y-auto">
      <div className="bg-white dark:bg-stone-800 rounded-lg p-6 w-full max-w-2xl mx-4 my-8 max-h-[90vh] overflow-y-auto">
        {isLoading ? (
          <div className="py-8">
            <LoadingSpinner size="md" />
          </div>
        ) : error && !form ? (
          <div className="text-center py-8">
            <p className="text-red-600 dark:text-red-400">{error}</p>
            <Button variant="secondary" onClick={onClose} className="mt-4">
              Close
            </Button>
          </div>
        ) : form ? (
          <>
            <div className="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg p-3 mb-4">
              <p className="text-sm text-amber-800 dark:text-amber-200">
                <span className="font-semibold">Please note:</span> Group leaders are mandated reporters. If your responses indicate a risk of harm to yourself or others, your leader may be required to report it.
              </p>
            </div>

            <div className="flex items-center justify-between mb-4">
              <div>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  {form.name}
                </h2>
                {form.description && (
                  <p className="text-sm text-gray-600 dark:text-stone-400 mt-1">
                    {form.description}
                  </p>
                )}
              </div>
              <button
                onClick={onClose}
                className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div className="space-y-5">
              {form.fields.map((field) => {
                if (field.fieldType === 'TEXT_DISPLAY') {
                  return (
                    <div key={field.id} className="text-gray-700 dark:text-stone-300">
                      <p className="font-medium">{field.label}</p>
                      {field.description && (
                        <p className="text-sm text-gray-500 dark:text-stone-400 mt-1">{field.description}</p>
                      )}
                    </div>
                  );
                }

                return (
                  <div key={field.id}>
                    <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">
                      {field.label}
                      {field.isRequired && <span className="text-red-500 ml-1">*</span>}
                    </label>
                    {field.description && (
                      <p className="text-xs text-gray-500 dark:text-stone-400 mb-2">{field.description}</p>
                    )}

                    {field.fieldType === 'TEXT_SMALL' && (
                      <input
                        type="text"
                        value={(answers[field.id] as string) || ''}
                        onChange={(e) => handleFieldChange(field.id, e.target.value)}
                        className="w-full rounded-md border border-gray-300 dark:border-stone-600 bg-white dark:bg-stone-700 px-3 py-2 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                      />
                    )}

                    {field.fieldType === 'TEXT_MEDIUM' && (
                      <textarea
                        value={(answers[field.id] as string) || ''}
                        onChange={(e) => handleFieldChange(field.id, e.target.value)}
                        rows={3}
                        className="w-full rounded-md border border-gray-300 dark:border-stone-600 bg-white dark:bg-stone-700 px-3 py-2 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                      />
                    )}

                    {field.fieldType === 'TEXT_LARGE' && (
                      <textarea
                        value={(answers[field.id] as string) || ''}
                        onChange={(e) => handleFieldChange(field.id, e.target.value)}
                        rows={6}
                        className="w-full rounded-md border border-gray-300 dark:border-stone-600 bg-white dark:bg-stone-700 px-3 py-2 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                      />
                    )}

                    {field.fieldType === 'CHECKBOX' && field.options && (
                      <div className="space-y-2">
                        {field.options.map((option) => (
                          <label key={option} className="flex items-center gap-2">
                            <input
                              type="checkbox"
                              checked={((answers[field.id] as string[]) || []).includes(option)}
                              onChange={(e) => handleCheckboxChange(field.id, option, e.target.checked)}
                              className="rounded border-gray-300 dark:border-stone-600 text-amber-600 focus:ring-amber-500"
                            />
                            <span className="text-gray-700 dark:text-stone-300">{option}</span>
                          </label>
                        ))}
                      </div>
                    )}

                    {field.fieldType === 'RADIO' && field.options && (
                      <div className="space-y-2">
                        {field.options.map((option) => (
                          <label key={option} className="flex items-center gap-2">
                            <input
                              type="radio"
                              name={field.id}
                              checked={answers[field.id] === option}
                              onChange={() => handleFieldChange(field.id, option)}
                              className="border-gray-300 dark:border-stone-600 text-amber-600 focus:ring-amber-500"
                            />
                            <span className="text-gray-700 dark:text-stone-300">{option}</span>
                          </label>
                        ))}
                      </div>
                    )}

                    {field.fieldType === 'DROPDOWN' && field.options && (
                      <select
                        value={(answers[field.id] as string) || ''}
                        onChange={(e) => handleFieldChange(field.id, e.target.value)}
                        className="w-full rounded-md border border-gray-300 dark:border-stone-600 bg-white dark:bg-stone-700 px-3 py-2 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                      >
                        <option value="">Select...</option>
                        {field.options.map((option) => (
                          <option key={option} value={option}>
                            {option}
                          </option>
                        ))}
                      </select>
                    )}
                  </div>
                );
              })}
            </div>

            {error && (
              <p className="text-red-600 dark:text-red-400 text-sm mt-4">{error}</p>
            )}

            <div className="flex justify-end gap-3 mt-6">
              <Button variant="secondary" onClick={onClose}>
                Cancel
              </Button>
              <Button onClick={handleSubmit} disabled={isSubmitting}>
                {isSubmitting ? 'Submitting...' : existingAnswer ? 'Update' : 'Submit'}
              </Button>
            </div>
          </>
        ) : null}
      </div>
    </div>
  );
}
