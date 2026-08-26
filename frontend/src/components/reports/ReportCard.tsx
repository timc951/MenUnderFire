import { useState } from 'react';
import { Report } from '../../types';
import { formatDate, truncateText } from '../../utils/helpers';
import { Card } from '../common/Card';

interface ReportCardProps {
  report: Report;
}

const TRUNCATE_LENGTH = 200;

export function ReportCard({ report }: ReportCardProps) {
  const [expanded, setExpanded] = useState(false);
  const isLong = report.content.length > TRUNCATE_LENGTH;

  return (
    <Card className="space-y-3">
      <div className="flex items-center justify-between">
        <h3 className="font-semibold text-gray-900">{report.title}</h3>
        <span className="text-sm text-gray-500">{formatDate(report.createdAt)}</span>
      </div>
      <div className="flex items-center gap-2">
        {report.isAnonymous ? (
          <>
            <span className="text-sm text-gray-600">Anonymous</span>
            <span data-testid="anonymous-badge" className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-600">
              Anonymous
            </span>
          </>
        ) : (
          <span className="text-sm text-gray-600">{report.reporterName}</span>
        )}
      </div>
      <p className="text-gray-700">
        {isLong && !expanded ? truncateText(report.content, TRUNCATE_LENGTH) : report.content}
      </p>
      {isLong && (
        <button
          onClick={() => setExpanded(!expanded)}
          className="text-sm text-primary-600 hover:text-primary-700 font-medium"
        >
          {expanded ? 'Show less' : 'Show more'}
        </button>
      )}
    </Card>
  );
}
