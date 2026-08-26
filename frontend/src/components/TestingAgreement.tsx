import { useState } from 'react';

interface TestingAgreementProps {
  onAccepted: () => void;
  onDecline: () => void;
  getToken: () => Promise<string>;
}

const AGREEMENT_VERSION = '1.0';

export function TestingAgreement({ onAccepted, onDecline, getToken }: TestingAgreementProps) {
  const [isChecked, setIsChecked] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAccept = async () => {
    if (!isChecked) return;

    setIsSubmitting(true);
    setError(null);

    try {
      const token = await getToken();
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

      const response = await fetch(`${apiUrl}/users/me/accept-agreement`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ agreementVersion: AGREEMENT_VERSION }),
      });

      if (!response.ok) {
        const errorBody = await response.json().catch(() => ({ message: 'Request failed' }));
        throw new Error(errorBody.message || `HTTP ${response.status}`);
      }

      onAccepted();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to record acceptance');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-stone-900 overflow-y-auto p-4">
      <div className="max-w-2xl w-full bg-stone-800 rounded-xl p-5 sm:p-8 border border-stone-700 shadow-2xl my-auto">
        <div className="flex items-center space-x-3 mb-4 sm:mb-6">
          <img
            src="https://content.menunderfire.com/logo.webp"
            alt=""
            width={40}
            height={40}
            className="w-10 h-10"
          />
          <h1 className="text-xl sm:text-2xl font-bold text-amber-500">Testing Agreement</h1>
        </div>

        <div className="text-stone-300 space-y-3 sm:space-y-4 mb-6 sm:mb-8">
          <p className="text-base sm:text-lg font-medium text-white">
            Welcome! Before you continue, please read and acknowledge the following:
          </p>

          <div className="bg-stone-900/50 rounded-lg p-3 sm:p-4 border border-stone-600 space-y-3 text-sm sm:text-base">
            <div className="flex items-start space-x-3">
              <span className="text-amber-500 font-bold mt-0.5">1.</span>
              <p>
                <strong className="text-white">Testing Environment:</strong> This site is currently
                in a testing phase. Features may be incomplete, change without notice, or behave
                unexpectedly.
              </p>
            </div>
            <div className="flex items-start space-x-3">
              <span className="text-amber-500 font-bold mt-0.5">2.</span>
              <p>
                <strong className="text-white">No Private Information:</strong> Do not enter any
                sensitive personal information, passwords, or private data. Use placeholder or test
                data only.
              </p>
            </div>
            <div className="flex items-start space-x-3">
              <span className="text-amber-500 font-bold mt-0.5">3.</span>
              <p>
                <strong className="text-white">No Data Guarantees:</strong> There are no backups.
                Data may be lost, deleted, or reset at any time without prior notice.
              </p>
            </div>
            <div className="flex items-start space-x-3">
              <span className="text-amber-500 font-bold mt-0.5">4.</span>
              <p>
                <strong className="text-white">Bug & Feature Reporting:</strong> If you encounter
                issues or have feature suggestions, please use the{' '}
                <strong className="text-amber-400">Feedback</strong> form available within the
                application.
              </p>
            </div>
          </div>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-900/30 border border-red-700 text-red-400 rounded-md text-sm">
            {error}
          </div>
        )}

        <div className="space-y-4">
          <label className="flex items-start space-x-3 cursor-pointer select-none">
            <input
              type="checkbox"
              checked={isChecked}
              onChange={(e) => setIsChecked(e.target.checked)}
              className="mt-1 h-5 w-5 rounded border-stone-500 text-amber-600 focus:ring-amber-500"
            />
            <span className="text-stone-200">
              I have read and understand that this is a testing environment. I agree not to enter
              private information and acknowledge that data may be lost.
            </span>
          </label>

          <button
            onClick={handleAccept}
            disabled={!isChecked || isSubmitting}
            className="w-full py-3 px-4 bg-amber-600 text-white font-semibold rounded-lg hover:bg-amber-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {isSubmitting ? 'Recording acceptance...' : 'I Understand & Accept'}
          </button>

          <button
            onClick={onDecline}
            disabled={isSubmitting}
            className="w-full py-2 px-4 text-stone-400 hover:text-stone-200 text-sm transition-colors"
          >
            I do not accept, log out
          </button>
        </div>
      </div>
    </div>
  );
}
