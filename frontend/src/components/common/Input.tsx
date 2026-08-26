import { InputHTMLAttributes, forwardRef } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, id, className = '', ...props }, ref) => {
    const inputId = id || label.toLowerCase().replace(/\s+/g, '-');
    return (
      <div className="space-y-1">
        <label htmlFor={inputId} className="block text-sm font-medium text-gray-700 dark:text-stone-300">
          {label}
        </label>
        <input
          ref={ref}
          id={inputId}
          className={`block w-full rounded-md border ${error ? 'border-red-500' : 'border-gray-300 dark:border-stone-600'} bg-white dark:bg-stone-700 text-gray-900 dark:text-white px-3 py-2 shadow-sm focus:border-amber-500 focus:ring-1 focus:ring-amber-500 ${className}`}
          aria-invalid={!!error}
          aria-describedby={error ? `${inputId}-error` : undefined}
          {...props}
        />
        {error && (
          <p id={`${inputId}-error`} className="text-sm text-red-600 dark:text-red-400" role="alert">
            {error}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';
