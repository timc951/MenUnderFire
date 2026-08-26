import { useState, useEffect, useCallback } from 'react';
import { useApi } from '../hooks/useApi';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { Button } from '../components/common/Button';

interface HitSummary {
  totalHits: number;
  todayHits: number;
  pageStats: { path: string; hitCount: number; uniqueHits: number }[];
}

interface HitsByCountry {
  country: string;
  hitCount: number;
}

interface HitsByRegion {
  country: string;
  region: string;
  hitCount: number;
}

interface HitsHourly {
  hour: string;
  hitCount: number;
}

interface HitsDaily {
  date: string;
  hitCount: number;
}

interface RecentHit {
  id: string;
  path: string;
  userEmail: string;
  ipAddress?: string;
  country?: string;
  region?: string;
  city?: string;
  createdAt: string;
}

type TabKey = 'overview' | 'country' | 'region' | 'hourly' | 'daily' | 'recent';

export function PageHitsAdmin() {
  const api = useApi();
  const [activeTab, setActiveTab] = useState<TabKey>('overview');
  const [loading, setLoading] = useState(true);
  const [summary, setSummary] = useState<HitSummary | null>(null);
  const [byCountry, setByCountry] = useState<HitsByCountry[]>([]);
  const [byRegion, setByRegion] = useState<HitsByRegion[]>([]);
  const [hourly, setHourly] = useState<HitsHourly[]>([]);
  const [daily, setDaily] = useState<HitsDaily[]>([]);
  const [recent, setRecent] = useState<RecentHit[]>([]);

  // Date range state
  const [fromDate, setFromDate] = useState('');
  const [toDate, setToDate] = useState('');

  const buildDateParams = useCallback(() => {
    const params = new URLSearchParams();
    if (fromDate) params.set('from', fromDate);
    if (toDate) params.set('to', toDate);
    const qs = params.toString();
    return qs ? `?${qs}` : '';
  }, [fromDate, toDate]);

  const fetchData = useCallback(async (tab: TabKey) => {
    setLoading(true);
    try {
      const dateParams = buildDateParams();
      switch (tab) {
        case 'overview':
          setSummary(await api.get<HitSummary>('/hits/summary'));
          break;
        case 'country':
          setByCountry(await api.get<HitsByCountry[]>(`/hits/by-country${dateParams}`));
          break;
        case 'region':
          setByRegion(await api.get<HitsByRegion[]>(`/hits/by-region${dateParams}`));
          break;
        case 'hourly':
          setHourly(await api.get<HitsHourly[]>(`/hits/hourly${dateParams}`));
          break;
        case 'daily':
          setDaily(await api.get<HitsDaily[]>(`/hits/daily${dateParams}`));
          break;
        case 'recent':
          setRecent(await api.get<RecentHit[]>(`/hits/range${dateParams ? dateParams + '&' : '?'}limit=100`));
          break;
      }
    } catch (err) {
      console.error('Failed to load hit data:', err);
    } finally {
      setLoading(false);
    }
  }, [api, buildDateParams]);

  useEffect(() => {
    fetchData(activeTab);
  }, [activeTab, fetchData]);

  const handleApplyFilter = () => {
    fetchData(activeTab);
  };

  const handleClearFilter = () => {
    setFromDate('');
    setToDate('');
  };

  // Re-fetch when dates are cleared
  useEffect(() => {
    if (fromDate === '' && toDate === '') {
      fetchData(activeTab);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [fromDate, toDate]);

  const tabs: { key: TabKey; label: string }[] = [
    { key: 'overview', label: 'Overview' },
    { key: 'country', label: 'By Country' },
    { key: 'region', label: 'By Region' },
    { key: 'hourly', label: 'Hourly' },
    { key: 'daily', label: 'Daily' },
    { key: 'recent', label: 'Recent Hits' },
  ];

  const maxCount = (counts: number[]) => Math.max(...counts, 1);

  const showDateFilter = activeTab !== 'overview';

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Page Hit Analytics</h1>
        <p className="text-gray-600 dark:text-stone-400">View site traffic and visitor analytics</p>
      </div>

      {/* Tabs */}
      <div className="flex flex-wrap gap-1 border-b border-stone-200 dark:border-stone-700">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`px-4 py-2 text-sm font-medium rounded-t-md transition-colors ${
              activeTab === tab.key
                ? 'bg-amber-500 text-white'
                : 'text-stone-600 dark:text-stone-400 hover:bg-stone-100 dark:hover:bg-stone-700'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Date Range Filter */}
      {showDateFilter && (
        <div className="bg-white dark:bg-stone-800 rounded-lg border border-stone-200 dark:border-stone-700 p-4">
          <div className="flex flex-wrap items-end gap-4">
            <div>
              <label className="block text-xs font-medium text-stone-500 dark:text-stone-400 mb-1">From</label>
              <input
                type="date"
                value={fromDate}
                onChange={(e) => setFromDate(e.target.value)}
                className="rounded-md border border-stone-300 dark:border-stone-600 bg-white dark:bg-stone-700 text-gray-900 dark:text-white px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-amber-500"
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-stone-500 dark:text-stone-400 mb-1">To</label>
              <input
                type="date"
                value={toDate}
                onChange={(e) => setToDate(e.target.value)}
                className="rounded-md border border-stone-300 dark:border-stone-600 bg-white dark:bg-stone-700 text-gray-900 dark:text-white px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-amber-500"
              />
            </div>
            <Button onClick={handleApplyFilter} className="text-sm">
              Apply
            </Button>
            {(fromDate || toDate) && (
              <button
                onClick={handleClearFilter}
                className="text-sm text-stone-500 dark:text-stone-400 hover:text-stone-700 dark:hover:text-stone-200 underline"
              >
                Clear
              </button>
            )}
            {(fromDate || toDate) && (
              <span className="text-xs text-stone-500 dark:text-stone-400">
                Showing: {fromDate || 'start'} &mdash; {toDate || 'now'}
              </span>
            )}
          </div>
        </div>
      )}

      {loading ? (
        <LoadingSpinner size="lg" />
      ) : (
        <div>
          {/* Overview Tab */}
          {activeTab === 'overview' && summary && (
            <div className="space-y-6">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="bg-white dark:bg-stone-800 rounded-lg border border-stone-200 dark:border-stone-700 p-6">
                  <p className="text-sm text-stone-500 dark:text-stone-400">Total Hits</p>
                  <p className="text-3xl font-bold text-gray-900 dark:text-white">{summary.totalHits.toLocaleString()}</p>
                </div>
                <div className="bg-white dark:bg-stone-800 rounded-lg border border-stone-200 dark:border-stone-700 p-6">
                  <p className="text-sm text-stone-500 dark:text-stone-400">Today</p>
                  <p className="text-3xl font-bold text-amber-600 dark:text-amber-400">{summary.todayHits.toLocaleString()}</p>
                </div>
              </div>

              <div className="bg-white dark:bg-stone-800 rounded-lg border border-stone-200 dark:border-stone-700">
                <div className="p-4 border-b border-stone-200 dark:border-stone-700">
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Hits by Page</h2>
                </div>
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="text-left text-xs text-stone-500 dark:text-stone-400 uppercase border-b border-stone-200 dark:border-stone-700">
                        <th className="px-4 py-3">Path</th>
                        <th className="px-4 py-3 text-right">Hits</th>
                        <th className="px-4 py-3 text-right">Unique</th>
                      </tr>
                    </thead>
                    <tbody>
                      {summary.pageStats.map((stat) => (
                        <tr key={stat.path} className="border-b border-stone-100 dark:border-stone-700/50">
                          <td className="px-4 py-3 text-sm text-gray-900 dark:text-white font-mono">{stat.path}</td>
                          <td className="px-4 py-3 text-sm text-right text-gray-900 dark:text-white">{stat.hitCount.toLocaleString()}</td>
                          <td className="px-4 py-3 text-sm text-right text-stone-500 dark:text-stone-400">{stat.uniqueHits.toLocaleString()}</td>
                        </tr>
                      ))}
                      {summary.pageStats.length === 0 && (
                        <tr><td colSpan={3} className="px-4 py-8 text-center text-stone-500 dark:text-stone-400">No data yet</td></tr>
                      )}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          )}

          {/* By Country Tab */}
          {activeTab === 'country' && (
            <div className="bg-white dark:bg-stone-800 rounded-lg border border-stone-200 dark:border-stone-700">
              <div className="p-4 border-b border-stone-200 dark:border-stone-700">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Hits by Country</h2>
              </div>
              <div className="p-4 space-y-3">
                {byCountry.map((item) => {
                  const max = maxCount(byCountry.map((c) => c.hitCount));
                  const pct = (item.hitCount / max) * 100;
                  return (
                    <div key={item.country}>
                      <div className="flex justify-between text-sm mb-1">
                        <span className="text-gray-900 dark:text-white font-medium">{item.country}</span>
                        <span className="text-stone-500 dark:text-stone-400">{item.hitCount.toLocaleString()}</span>
                      </div>
                      <div className="w-full bg-stone-200 dark:bg-stone-700 rounded-full h-2">
                        <div className="bg-amber-500 h-2 rounded-full" style={{ width: `${pct}%` }} />
                      </div>
                    </div>
                  );
                })}
                {byCountry.length === 0 && (
                  <p className="text-center py-8 text-stone-500 dark:text-stone-400">No data yet</p>
                )}
              </div>
            </div>
          )}

          {/* By Region Tab */}
          {activeTab === 'region' && (
            <div className="bg-white dark:bg-stone-800 rounded-lg border border-stone-200 dark:border-stone-700">
              <div className="p-4 border-b border-stone-200 dark:border-stone-700">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Hits by Region</h2>
              </div>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="text-left text-xs text-stone-500 dark:text-stone-400 uppercase border-b border-stone-200 dark:border-stone-700">
                      <th className="px-4 py-3">Country</th>
                      <th className="px-4 py-3">Region</th>
                      <th className="px-4 py-3 text-right">Hits</th>
                    </tr>
                  </thead>
                  <tbody>
                    {byRegion.map((item, i) => (
                      <tr key={`${item.country}-${item.region}-${i}`} className="border-b border-stone-100 dark:border-stone-700/50">
                        <td className="px-4 py-3 text-sm text-gray-900 dark:text-white">{item.country}</td>
                        <td className="px-4 py-3 text-sm text-gray-900 dark:text-white">{item.region}</td>
                        <td className="px-4 py-3 text-sm text-right text-gray-900 dark:text-white">{item.hitCount.toLocaleString()}</td>
                      </tr>
                    ))}
                    {byRegion.length === 0 && (
                      <tr><td colSpan={3} className="px-4 py-8 text-center text-stone-500 dark:text-stone-400">No data yet</td></tr>
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          )}

          {/* Hourly Tab */}
          {activeTab === 'hourly' && (
            <div className="bg-white dark:bg-stone-800 rounded-lg border border-stone-200 dark:border-stone-700">
              <div className="p-4 border-b border-stone-200 dark:border-stone-700">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  Hits per Hour {!fromDate && !toDate ? '(Last 24h)' : ''}
                </h2>
              </div>
              <div className="p-4">
                {hourly.length > 0 ? (
                  <div className="flex items-end gap-1 h-48">
                    {hourly.map((item) => {
                      const max = maxCount(hourly.map((h) => h.hitCount));
                      const pct = (item.hitCount / max) * 100;
                      const hourLabel = item.hour.slice(11, 16);
                      return (
                        <div key={item.hour} className="flex-1 flex flex-col items-center gap-1" title={`${item.hour}: ${item.hitCount} hits`}>
                          <span className="text-xs text-stone-500 dark:text-stone-400">{item.hitCount}</span>
                          <div className="w-full bg-amber-500 rounded-t" style={{ height: `${Math.max(pct, 2)}%` }} />
                          <span className="text-xs text-stone-500 dark:text-stone-400 rotate-[-45deg] origin-center whitespace-nowrap">{hourLabel}</span>
                        </div>
                      );
                    })}
                  </div>
                ) : (
                  <p className="text-center py-8 text-stone-500 dark:text-stone-400">No data for this period</p>
                )}
              </div>
            </div>
          )}

          {/* Daily Tab */}
          {activeTab === 'daily' && (
            <div className="bg-white dark:bg-stone-800 rounded-lg border border-stone-200 dark:border-stone-700">
              <div className="p-4 border-b border-stone-200 dark:border-stone-700">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  Hits per Day {!fromDate && !toDate ? '(Last 30 Days)' : ''}
                </h2>
              </div>
              <div className="p-4">
                {daily.length > 0 ? (
                  <div className="flex items-end gap-1 h-48">
                    {daily.map((item) => {
                      const max = maxCount(daily.map((d) => d.hitCount));
                      const pct = (item.hitCount / max) * 100;
                      const dateLabel = item.date.slice(5);
                      return (
                        <div key={item.date} className="flex-1 flex flex-col items-center gap-1" title={`${item.date}: ${item.hitCount} hits`}>
                          <span className="text-xs text-stone-500 dark:text-stone-400">{item.hitCount}</span>
                          <div className="w-full bg-amber-500 rounded-t" style={{ height: `${Math.max(pct, 2)}%` }} />
                          <span className="text-xs text-stone-500 dark:text-stone-400 rotate-[-45deg] origin-center whitespace-nowrap">{dateLabel}</span>
                        </div>
                      );
                    })}
                  </div>
                ) : (
                  <p className="text-center py-8 text-stone-500 dark:text-stone-400">No data for this period</p>
                )}
              </div>
            </div>
          )}

          {/* Recent Hits Tab */}
          {activeTab === 'recent' && (
            <div className="bg-white dark:bg-stone-800 rounded-lg border border-stone-200 dark:border-stone-700">
              <div className="p-4 border-b border-stone-200 dark:border-stone-700">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Recent Hits</h2>
              </div>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="text-left text-xs text-stone-500 dark:text-stone-400 uppercase border-b border-stone-200 dark:border-stone-700">
                      <th className="px-4 py-3">Time</th>
                      <th className="px-4 py-3">Path</th>
                      <th className="px-4 py-3">User</th>
                      <th className="px-4 py-3">Country</th>
                      <th className="px-4 py-3">Region</th>
                      <th className="px-4 py-3">IP</th>
                    </tr>
                  </thead>
                  <tbody>
                    {recent.map((hit) => (
                      <tr key={hit.id} className="border-b border-stone-100 dark:border-stone-700/50">
                        <td className="px-4 py-3 text-sm text-stone-500 dark:text-stone-400 whitespace-nowrap">
                          {new Date(hit.createdAt).toLocaleString()}
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-900 dark:text-white font-mono">{hit.path}</td>
                        <td className="px-4 py-3 text-sm text-gray-900 dark:text-white">{hit.userEmail}</td>
                        <td className="px-4 py-3 text-sm text-gray-900 dark:text-white">{hit.country || '-'}</td>
                        <td className="px-4 py-3 text-sm text-gray-900 dark:text-white">{hit.region || '-'}</td>
                        <td className="px-4 py-3 text-sm text-stone-500 dark:text-stone-400 font-mono">{hit.ipAddress || '-'}</td>
                      </tr>
                    ))}
                    {recent.length === 0 && (
                      <tr><td colSpan={6} className="px-4 py-8 text-center text-stone-500 dark:text-stone-400">No hits recorded yet</td></tr>
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
