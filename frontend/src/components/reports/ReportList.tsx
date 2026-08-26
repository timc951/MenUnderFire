import { Report } from '../../types';
import { ReportCard } from './ReportCard';
import { LoadingSpinner } from '../common/LoadingSpinner';

interface ReportListProps {
  reports: Report[];
  isLoading?: boolean;
  error?: string | null;
}

export function ReportList({ reports, isLoading, error }: ReportListProps) {
  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <p className="text-red-600">{error}</p>
      </div>
    );
  }

  if (reports.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-500">No reports yet. Be the first to submit one!</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {reports.map((report) => (
        <ReportCard key={report.id} report={report} />
      ))}
    </div>
  );
}
