import { useState } from 'react';

interface AnonymityToggleProps {
  value?: boolean;
  onChange?: (value: boolean) => void;
}

export function AnonymityToggle({ value, onChange }: AnonymityToggleProps) {
  const [isAnonymous, setIsAnonymous] = useState(value ?? false);
  const [showTooltip, setShowTooltip] = useState(false);

  const currentValue = value !== undefined ? value : isAnonymous;

  const handleToggle = () => {
    const newValue = !currentValue;
    setIsAnonymous(newValue);
    onChange?.(newValue);
  };

  return (
    <div className="space-y-2">
      <div className="flex items-center gap-3">
        <button
          type="button"
          role="switch"
          aria-checked={currentValue}
          aria-label="Anonymous"
          onClick={handleToggle}
          className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
            currentValue ? 'bg-amber-600' : 'bg-gray-200 dark:bg-stone-600'
          }`}
        >
          <span
            className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
              currentValue ? 'translate-x-6' : 'translate-x-1'
            }`}
          />
        </button>
        <span className="text-sm font-medium text-gray-700 dark:text-stone-300">Submit Anonymously</span>
        <div
          className="relative"
          onMouseEnter={() => setShowTooltip(true)}
          onMouseLeave={() => setShowTooltip(false)}
        >
          <span data-testid="info-icon" className="text-gray-400 dark:text-stone-500 cursor-help">ⓘ</span>
          {showTooltip && (
            <div role="tooltip" className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 dark:bg-stone-800 text-white text-xs rounded-md whitespace-nowrap">
              Anonymity hides your name from group members, but leaders can still see it.
            </div>
          )}
        </div>
      </div>
      <p className="text-sm text-gray-500 dark:text-stone-400">
        {currentValue
          ? 'Group members will not see your name. Leader will still see your identity.'
          : 'Your name will be visible to all group members.'}
      </p>
    </div>
  );
}
