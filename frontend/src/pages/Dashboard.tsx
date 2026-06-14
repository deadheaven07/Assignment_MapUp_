import { useMemo } from 'react';
import { useQueries, useQuery } from '@tanstack/react-query';
import { MapPanel } from '../components/MapPanel';
import { getAlerts, getGeofences, getLatestLocation, getVehicles, getViolations } from '../api/queries';
import type { VehicleLocation } from '../types';

export function Dashboard() {
  const { data: geofences = [], isLoading: geofencesLoading } = useQuery({ queryKey: ['geofences'], queryFn: getGeofences });
  const { data: vehicles = [], isLoading: vehiclesLoading } = useQuery({ queryKey: ['vehicles'], queryFn: getVehicles });
  const { data: alerts = [], isLoading: alertsLoading } = useQuery({ queryKey: ['alerts'], queryFn: getAlerts });
  const { data: violationsData, isLoading: violationsLoading } = useQuery({ queryKey: ['violations', 'dashboard'], queryFn: () => getViolations({ limit: 100, page: 1 }) });
  const locationQueries = useQueries({
    queries: vehicles.map((vehicle) => ({ queryKey: ['latest-location', vehicle.id], queryFn: () => getLatestLocation(vehicle.id), retry: false })),
  });
  const locations = locationQueries.map((query) => query.data).filter((location): location is VehicleLocation => Boolean(location));
  const isLoading = geofencesLoading || vehiclesLoading || alertsLoading || violationsLoading;

  const alertCounts = useMemo(
    () => ({
      entry: alerts.filter((alert) => alert.event_type === 'entry').length,
      exit: alerts.filter((alert) => alert.event_type === 'exit').length,
    }),
    [alerts],
  );

  const violations = violationsData?.violations ?? [];
  const violationCounts = useMemo(
    () => ({
      entry: violations.filter((item) => item.event_type === 'entry').length,
      exit: violations.filter((item) => item.event_type === 'exit').length,
    }),
    [violations],
  );

  return (
    <section className="space-y-6">
      <div className="grid gap-4 lg:grid-cols-4">
        <DashboardStat label="Total geofences" value={geofences.length} loading={geofencesLoading} />
        <DashboardStat label="Total vehicles" value={vehicles.length} loading={vehiclesLoading} />
        <DashboardStat label="Total alerts" value={alerts.length} loading={alertsLoading} />
        <DashboardStat label="Total violations" value={violationsData?.total_count ?? 0} loading={violationsLoading} />
      </div>

      <div className="grid gap-4 xl:grid-cols-[1.5fr_1fr]">
        <div className="grid gap-4">
          <div className="rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
              <div>
                <p className="text-sm uppercase tracking-[0.22em] text-slate-500">Event summary</p>
                <h2 className="mt-2 text-xl font-semibold text-slate-900">Entry vs exit alerts</h2>
              </div>
              <p className="text-sm text-slate-500">{alerts.length} total events</p>
            </div>
            <div className="mt-6 space-y-4">
              {['entry', 'exit'].map((type) => {
                const count = type === 'entry' ? alertCounts.entry : alertCounts.exit;
                const total = Math.max(alertCounts.entry + alertCounts.exit, 1);
                const ratio = Math.round((count / total) * 100);
                return (
                  <div key={type}>
                    <div className="flex items-center justify-between text-sm font-medium text-slate-700">
                      <span>{type}</span>
                      <span>{count}</span>
                    </div>
                    <div className="mt-2 h-2 rounded-full bg-slate-200">
                      <div className={`h-2 rounded-full ${type === 'entry' ? 'bg-emerald-500' : 'bg-rose-500'}`} style={{ width: `${ratio}%` }} />
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          <div className="rounded-3xl border border-slate-200 bg-white p-5 shadow-sm transition-transform duration-200 ease-out hover:-translate-y-1 hover:shadow-lg">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
              <div>
                <p className="text-sm uppercase tracking-[0.22em] text-slate-500">Violations overview</p>
                <h2 className="mt-2 text-xl font-semibold text-slate-900">Recent violation types</h2>
              </div>
              <p className="text-sm text-slate-500">Showing latest {violations.length} records</p>
            </div>
            <div className="mt-6 space-y-4">
              {['entry', 'exit'].map((type) => {
                const count = type === 'entry' ? violationCounts.entry : violationCounts.exit;
                const total = Math.max(violationCounts.entry + violationCounts.exit, 1);
                const ratio = Math.round((count / total) * 100);
                return (
                  <div key={type}>
                    <div className="flex items-center justify-between text-sm font-medium text-slate-700">
                      <span>{type}</span>
                      <span>{count}</span>
                    </div>
                    <div className="mt-2 h-2 rounded-full bg-slate-200">
                      <div className={`h-2 rounded-full ${type === 'entry' ? 'bg-sky-500' : 'bg-orange-500'}`} style={{ width: `${ratio}%` }} />
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </div>

        <div className="rounded-3xl border border-slate-200 bg-white p-5 shadow-sm transition-transform duration-200 ease-out hover:-translate-y-1 hover:shadow-lg">
          <div className="flex items-center justify-between gap-4">
            <div>
              <p className="text-sm uppercase tracking-[0.22em] text-slate-500">Fleet snapshot</p>
              <h2 className="mt-2 text-xl font-semibold text-slate-900">Map overview</h2>
            </div>
            <p className="text-sm text-slate-500">{locations.length} active positions</p>
          </div>
          <div className="mt-5 rounded-3xl border border-slate-200 bg-slate-50 p-4 text-sm text-slate-600">
            {isLoading ? (
              <div className="space-y-3">
                <div className="h-3 w-32 animate-pulse rounded-full bg-slate-200" />
                <div className="h-3 w-56 animate-pulse rounded-full bg-slate-200" />
              </div>
            ) : (
              <div className="space-y-3">
                <div className="rounded-2xl bg-slate-100 p-4">
                  <p className="text-sm text-slate-500">Geofences</p>
                  <p className="mt-1 text-2xl font-semibold text-slate-900">{geofences.length}</p>
                </div>
                <div className="rounded-2xl bg-slate-100 p-4">
                  <p className="text-sm text-slate-500">Vehicles reporting</p>
                  <p className="mt-1 text-2xl font-semibold text-slate-900">{vehicles.length}</p>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      <MapPanel geofences={geofences} locations={locations} />
    </section>
  );
}

function DashboardStat({ label, value, loading }: { label: string; value: number; loading: boolean }) {
  return (
    <div className="overflow-hidden rounded-3xl border border-slate-200 bg-white p-5 shadow-sm transition-transform duration-200 ease-out hover:-translate-y-1 hover:shadow-lg">
      <div className="text-sm font-medium uppercase tracking-[0.24em] text-slate-500">{label}</div>
      <div className="mt-4 text-4xl font-semibold text-slate-900">{loading ? <span className="inline-block h-10 w-24 animate-pulse rounded-full bg-slate-200" /> : value}</div>
    </div>
  );
}
