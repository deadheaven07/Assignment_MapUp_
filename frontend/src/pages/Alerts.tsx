import { FormEvent, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { createAlert, getAlerts, getGeofences, getVehicles } from '../api/queries';
import { StatusBadge } from '../components/StatusBadge';

export function Alerts() {
  const client = useQueryClient();
  const { data: alerts = [], isLoading: alertsLoading } = useQuery({ queryKey: ['alerts'], queryFn: getAlerts });
  const { data: vehicles = [], isLoading: vehiclesLoading } = useQuery({ queryKey: ['vehicles'], queryFn: getVehicles });
  const { data: geofences = [], isLoading: geofencesLoading } = useQuery({ queryKey: ['geofences'], queryFn: getGeofences });
  const [vehicleId, setVehicleId] = useState('');
  const [geofenceId, setGeofenceId] = useState('');
  const [eventType, setEventType] = useState<'entry' | 'exit' | 'both'>('entry');
  const mutation = useMutation({ mutationFn: createAlert, onSuccess: () => client.invalidateQueries({ queryKey: ['alerts'] }) });
  const submit = (event: FormEvent) => {
    event.preventDefault();
    if (!geofenceId) return;
    mutation.mutate({ vehicle_id: vehicleId || undefined, geofence_id: geofenceId, event_type: eventType });
  };

  return (
    <section className="space-y-6">
      <form onSubmit={submit} className="grid gap-4 rounded-3xl border border-slate-200 bg-white p-5 shadow-sm md:grid-cols-[1.4fr_1fr_1fr_auto]">
        <select className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" value={vehicleId} onChange={(e) => setVehicleId(e.target.value)}>
          <option value="">Any vehicle</option>
          {vehicles.map((v) => <option key={v.id} value={v.id}>{v.vehicle_number}</option>)}
        </select>
        <select className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" value={geofenceId} onChange={(e) => setGeofenceId(e.target.value)}>
          <option value="">Geofence</option>
          {geofences.map((g) => <option key={g.id} value={g.id}>{g.name}</option>)}
        </select>
        <select className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" value={eventType} onChange={(e) => setEventType(e.target.value as 'entry' | 'exit' | 'both')}>
          <option value="entry">entry</option>
          <option value="exit">exit</option>
          <option value="both">both</option>
        </select>
        <button className="rounded-3xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white transition hover:bg-slate-800" type="submit">Configure</button>
      </form>

      <div className="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm">
        <div className="grid gap-4 border-b border-slate-200 bg-slate-50 px-4 py-4 text-xs uppercase tracking-[0.22em] text-slate-500 sm:grid-cols-[1.5fr_1fr_1fr_1fr]">
          <span>Event</span>
          <span>Vehicle</span>
          <span>Geofence</span>
          <span>Status</span>
        </div>
        {alertsLoading ? (
          <div className="space-y-3 p-4">
            {[...Array(4)].map((_, index) => (
              <div key={index} className="h-14 rounded-3xl bg-slate-100" />
            ))}
          </div>
        ) : alerts.length === 0 ? (
          <div className="p-10 text-center text-slate-500">
            No alert configurations yet. Create a rule to receive live entry and exit notifications.
          </div>
        ) : (
          alerts.map((alert, index) => (
            <div key={alert.alert_id} className={`grid min-h-[72px] items-center gap-4 px-4 py-4 text-sm sm:grid-cols-[1.5fr_1fr_1fr_1fr] ${index % 2 === 0 ? 'bg-white' : 'bg-slate-50'} transition hover:bg-slate-100`}>
              <div className="flex items-center gap-2">
                <StatusBadge value={alert.event_type} />
              </div>
              <span>{alert.vehicle_number ?? 'Any vehicle'}</span>
              <span>{alert.geofence_name ?? alert.geofence_id}</span>
              <StatusBadge value={alert.status ?? 'unknown'} />
            </div>
          ))
        )}
      </div>
    </section>
  );
}
