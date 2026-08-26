import { useState, KeyboardEvent } from 'react';
import { GroupMessage } from '../../types';
import { LoadingSpinner } from '../common/LoadingSpinner';

interface GroupMessagesProps {
  messages: GroupMessage[];
  isLoading: boolean;
  isLeader: boolean;
  onSendMessage: (content: string, notifyMembers: boolean) => Promise<void>;
  onDeleteMessage: (messageId: string) => Promise<void>;
  onFillForm?: (formId: string) => void;
}

export function GroupMessages({
  messages,
  isLoading,
  isLeader,
  onSendMessage,
  onDeleteMessage,
  onFillForm,
}: GroupMessagesProps) {
  const [content, setContent] = useState('');
  const [isSending, setIsSending] = useState(false);

  const handleSend = async () => {
    if (!content.trim() || isSending) return;
    setIsSending(true);
    try {
      await onSendMessage(content.trim(), false);
      setContent('');
    } finally {
      setIsSending(false);
    }
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
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

  if (isLoading) {
    return <LoadingSpinner size="lg" />;
  }

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <input
          type="text"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          onKeyDown={handleKeyDown}
          className="flex-1 rounded-md border border-gray-300 dark:border-stone-600 bg-white dark:bg-stone-800 px-3 py-2 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
          placeholder="Type a message..."
          disabled={isSending}
        />
        <button
          onClick={handleSend}
          disabled={isSending || !content.trim()}
          className="px-4 py-2 bg-amber-600 hover:bg-amber-700 disabled:bg-amber-400 text-white rounded-md font-medium transition-colors"
        >
          {isSending ? 'Sending...' : 'Send'}
        </button>
      </div>

      {messages.length === 0 ? (
        <p className="text-gray-500 dark:text-stone-400 text-center py-8">
          No messages yet.
        </p>
      ) : (
        <div className="space-y-4">
          {messages.map((message) => (
            <div
              key={message.id}
              className="bg-white dark:bg-stone-800 rounded-lg border border-gray-200 dark:border-stone-700 p-4"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="font-medium text-gray-900 dark:text-white">
                      {message.senderName}
                    </span>
                    <span className="text-sm text-gray-500 dark:text-stone-400">
                      {formatDate(message.createdAt)}
                    </span>
                    {message.notifyMembers && (
                      <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-amber-100 dark:bg-amber-900/50 text-amber-800 dark:text-amber-300">
                        Notified
                      </span>
                    )}
                  </div>
                  {message.formId ? (
                    <div className="mt-2 border border-amber-200 dark:border-amber-800 bg-amber-50 dark:bg-amber-900/20 rounded-lg p-3">
                      <div className="flex items-center gap-2 mb-2">
                        <svg className="w-5 h-5 text-amber-600 dark:text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                        </svg>
                        <span className="font-medium text-amber-800 dark:text-amber-200">
                          {message.formName || 'Form'}
                        </span>
                      </div>
                      {onFillForm && (
                        <button
                          onClick={() => onFillForm(message.formId!)}
                          className="px-3 py-1.5 bg-amber-600 hover:bg-amber-700 text-white text-sm rounded-md font-medium transition-colors"
                        >
                          Fill Out
                        </button>
                      )}
                    </div>
                  ) : (
                    <p className="text-gray-700 dark:text-stone-300 whitespace-pre-wrap">
                      {message.content}
                    </p>
                  )}
                </div>
                {isLeader && (
                  <button
                    onClick={() => onDeleteMessage(message.id)}
                    className="text-gray-400 hover:text-red-500 dark:text-stone-500 dark:hover:text-red-400"
                    title="Delete message"
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                    >
                      <path
                        fillRule="evenodd"
                        d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                        clipRule="evenodd"
                      />
                    </svg>
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
