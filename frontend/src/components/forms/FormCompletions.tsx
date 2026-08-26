import { useState } from 'react';
import { Form, FormAnswer, FormDetail } from '../../types';
import { Card } from '../common/Card';
import { LoadingSpinner } from '../common/LoadingSpinner';
import { Button } from '../common/Button';
import { useApi } from '../../hooks/useApi';

interface FormCompletionsProps {
  forms: Form[];
  isLoading: boolean;
}

export function FormCompletions({ forms, isLoading }: FormCompletionsProps) {
  const api = useApi();
  const [selectedForm, setSelectedForm] = useState<Form | null>(null);
  const [answers, setAnswers] = useState<FormAnswer[]>([]);
  const [answersLoading, setAnswersLoading] = useState(false);
  const [viewingAnswer, setViewingAnswer] = useState<FormAnswer | null>(null);
  const [formDetail, setFormDetail] = useState<FormDetail | null>(null);

  const handleViewForm = async (form: Form) => {
    setSelectedForm(form);
    setAnswersLoading(true);
    try {
      const [answersData, detailData] = await Promise.all([
        api.get<FormAnswer[]>(`/forms/${form.id}/answers`),
        api.get<FormDetail>(`/forms/${form.id}`),
      ]);
      setAnswers(answersData.filter(a => a.isCurrent));
      setFormDetail(detailData);
    } catch (err) {
      console.error('Failed to load form answers:', err);
    } finally {
      setAnswersLoading(false);
    }
  };

  const handleBack = () => {
    setSelectedForm(null);
    setAnswers([]);
    setFormDetail(null);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const renderAnswerValue = (value: unknown): string => {
    if (Array.isArray(value)) {
      return value.join(', ');
    }
    if (typeof value === 'boolean') {
      return value ? 'Yes' : 'No';
    }
    return String(value ?? '');
  };

  if (isLoading) {
    return <LoadingSpinner size="lg" />;
  }

  if (forms.length === 0) {
    return (
      <p className="text-gray-500 dark:text-stone-400 text-center py-8">
        No forms available. Create forms in your organization to collect responses.
      </p>
    );
  }

  // Show list of forms
  if (!selectedForm) {
    return (
      <div className="space-y-3">
        {forms.map((form) => (
          <Card
            key={form.id}
            className="cursor-pointer hover:bg-stone-50 dark:hover:bg-stone-700/50 transition-colors"
            onClick={() => handleViewForm(form)}
          >
            <div className="flex items-center justify-between">
              <div>
                <h3 className="font-medium text-gray-900 dark:text-white">{form.name}</h3>
                {form.description && (
                  <p className="text-sm text-gray-600 dark:text-stone-400 mt-1">{form.description}</p>
                )}
                <p className="text-xs text-gray-500 dark:text-stone-500 mt-1">
                  {form.fieldCount} fields
                </p>
              </div>
              <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </div>
          </Card>
        ))}
      </div>
    );
  }

  // Show form answers
  return (
    <div className="space-y-4">
      <div className="flex items-center gap-4">
        <button
          onClick={handleBack}
          className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
          aria-label="Back"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
          {selectedForm.name} - Responses
        </h3>
      </div>

      {answersLoading ? (
        <LoadingSpinner size="md" />
      ) : answers.length === 0 ? (
        <p className="text-gray-500 dark:text-stone-400 text-center py-8">
          No responses yet for this form.
        </p>
      ) : (
        <div className="space-y-3">
          {answers.map((answer) => (
            <Card
              key={answer.id}
              className="cursor-pointer hover:bg-stone-50 dark:hover:bg-stone-700/50 transition-colors"
              onClick={() => setViewingAnswer(answer)}
            >
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="font-medium text-gray-900 dark:text-white">
                    {answer.userName || 'Anonymous'}
                  </h4>
                  <p className="text-sm text-gray-500 dark:text-stone-400">
                    Submitted {formatDate(answer.submittedAt)}
                  </p>
                  {answer.version > 1 && (
                    <span className="text-xs text-amber-600 dark:text-amber-400">
                      Version {answer.version}
                    </span>
                  )}
                </div>
                <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* View Answer Modal */}
      {viewingAnswer && formDetail && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 overflow-y-auto">
          <div className="bg-white dark:bg-stone-800 rounded-lg p-6 w-full max-w-2xl mx-4 my-8 max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  {viewingAnswer.userName || 'Anonymous'}
                </h2>
                <p className="text-sm text-gray-500 dark:text-stone-400">
                  Submitted {formatDate(viewingAnswer.submittedAt)}
                </p>
              </div>
              <button
                onClick={() => setViewingAnswer(null)}
                className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div className="space-y-4">
              {formDetail.fields.map((field) => {
                const value = viewingAnswer.answers[field.id];
                if (field.fieldType === 'TEXT_DISPLAY') return null;
                return (
                  <div key={field.id} className="border-b border-gray-200 dark:border-stone-700 pb-3 last:border-0">
                    <p className="text-sm font-medium text-gray-700 dark:text-stone-300">
                      {field.label}
                    </p>
                    <p className="text-gray-900 dark:text-white mt-1">
                      {value !== undefined ? renderAnswerValue(value) : <span className="text-gray-400 dark:text-stone-500 italic">No answer</span>}
                    </p>
                  </div>
                );
              })}
            </div>

            <div className="flex justify-end mt-6">
              <Button variant="secondary" onClick={() => setViewingAnswer(null)}>
                Close
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
