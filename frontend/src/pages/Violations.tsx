import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { getGeofences, getVehicles, getViolations } from '../api/queries';
import { StatusBadge } from '../components/StatusBadge';

export function Violations() {
  const [vehicleId, setVehicleId] = useState('');
  const [geofenceId, setGeofenceId] = useState('');
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');
  const [page, setPage] = useState(1);
  const { data: vehicles = [], isLoading: vehiclesLoading } = useQuery({ queryKey: ['vehicles'], queryFn: getVehicles });
  const { data: geofences = [], isLoading: geofencesLoading } = useQuery({ queryKey: ['geofences'], queryFn: getGeofences });
  const { data, isLoading: violationsLoading } = useQuery({ queryKey: ['violations', vehicleId, geofenceId, from, to, page], queryFn: () => getViolations({ vehicle_id: vehicleId || undefined, geofence_id: geofenceId || undefined, start_date: from ? new Date(from).toISOString() : undefined, end_date: to ? new Date(to).toISOString() : undefined, limit: 10, page }) });

  return (
    <section className="space-y-6">
      <div className="grid gap-4 rounded-3xl border border-slate-200 bg-white p-5 shadow-sm md:grid-cols-5">
        <div>
          <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.22em] text-slate-500">Vehicle</label>
          <select className="w-full rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" value={vehicleId} onChange={(e) => setVehicleId(e.target.value)}>
            <option value="">All vehicles</option>
            {vehicles.map((v) => <option key={v.id} value={v.id}>{v.vehicle_number}</option>)}
          </select>
        </div>
        <div>
          <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.22em] text-slate-500">Geofence</label>
          <select className="w-full rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" value={geofenceId} onChange={(e) => setGeofenceId(e.target.value)}>
            <option value="">All geofences</option>
            {geofences.map((g) => <option key={g.id} value={g.id}>{g.name}</option>)}
          </select>
        </div>
        <div>
          <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.22em] text-slate-500">From</label>
          <input className="w-full rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" type="datetime-local" value={from} onChange={(e) => setFrom(e.target.value)} />
        </div>
        <div>
          <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.22em] text-slate-500">To</label>
          <input className="w-full rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" type="datetime-local" value={to} onChange={(e) => setTo(e.target.value)} />
        </div>
        <button className="self-end rounded-3xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white transition hover:bg-slate-800" onClick={() => setPage(1)} type="button">Apply</button>
      </div>
      <div className="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm">
        <div className="grid gap-4 border-b border-slate-200 bg-slate-50 px-4 py-4 text-xs uppercase tracking-[0.22em] text-slate-500 md:grid-cols-[1.2fr_0.9fr_0.9fr_0.9fr_1fr]">
          <span>Type</span>
          <span>Vehicle</span>
          <span>Geofence</span>
          <span>Location</span>
          <span>Timestamp</span>
        </div>
        {violationsLoading ? (
          <div className="space-y-3 p-4">
            {[...Array(5)].map((_, index) => (
              <div key={index} className="h-14 rounded-3xl bg-slate-100" />
            ))}
          </div>
        ) : !data?.violations.length ? (
          <div className="p-10 text-center text-slate-500">No violations found for the selected filters.</div>
        ) : (
          data.violations.map((item, index) => (
            <div key={item.id} className={`grid min-h-[72px] items-center gap-4 px-4 py-4 text-sm ${index % 2 === 0 ? 'bg-white' : 'bg-slate-50'} md:grid-cols-[1.2fr_0.9fr_0.9fr_0.9fr_1fr] transition hover:bg-slate-100`}>
              <StatusBadge value={item.event_type} />
              <span>{item.vehicle_number}</span>
              <span>{item.geofence_name}</span>
              <span>{item.latitude.toFixed(5)}, {item.longitude.toFixed(5)}</span>
              <span>{new Date(item.timestamp).toLocaleString()}</span>
            </div>
          ))
        )}
      </div>
      <div className="flex flex-col gap-3 rounded-3xl border border-slate-200 bg-white p-4 shadow-sm sm:flex-row sm:items-center sm:justify-between">
        <button className="rounded-3xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-700 transition hover:bg-slate-50" disabled={page === 1} onClick={() => setPage((p) => p - 1)}>Previous</button>
        <span className="text-sm text-slate-600">Page {page} / {Math.max(1, Math.ceil((data?.total_count ?? 0) / 10))}</span>
        <button className="rounded-3xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-700 transition hover:bg-slate-50" disabled={page >= Math.ceil((data?.total_count ?? 0) / 10)} onClick={() => setPage((p) => p + 1)}>Next</button>
      </div>
    </section>
  );
}
