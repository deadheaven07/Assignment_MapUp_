import { Bell, CircleDashed, MapPin, Truck } from 'lucide-react';
import type { AlertEvent } from '../types';
import { StatusBadge } from './StatusBadge';

export function LiveAlerts({ alerts, connected }: { alerts: AlertEvent[]; connected: boolean }) {
  return (
    <aside className="space-y-4 rounded-3xl border border-slate-200 bg-white p-5 shadow-sm lg:sticky lg:top-5">
      <div className="flex items-center justify-between gap-3">
        <div className="flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.18em] text-slate-700">
          <Bell size={18} />
          Live feed
        </div>
        <div className="flex items-center gap-2 text-sm text-slate-500">
          <span className={`h-2.5 w-2.5 rounded-full ${connected ? 'bg-emerald-500' : 'bg-slate-300'}`} />
          {connected ? 'Connected' : 'Offline'}
        </div>
      </div>
      <div className="space-y-3 overflow-auto max-h-[680px]">
        {alerts.length === 0 ? (
          <div className="rounded-3xl border border-dashed border-slate-200 bg-slate-50 p-6 text-sm text-slate-500">
            No live alerts yet. Incoming events appear here in real time.
          </div>
        ) : (
          alerts.map((alert) => (
            <article key={`${alert.event_id}-${alert.timestamp}`} className="overflow-hidden rounded-3xl border border-slate-200 bg-slate-50 p-4 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
              <div className="flex items-start justify-between gap-3">
                <div className="flex items-center gap-2 text-sm font-semibold text-slate-900">
                  <Truck size={18} />
                  {alert.vehicle.vehicle_number}
                </div>
                <StatusBadge value={alert.event_type} />
              </div>
              <div className="mt-3 flex items-center gap-2 text-sm text-slate-600">
                <MapPin size={16} />
                <span>{alert.geofence.geofence_name}</span>
              </div>
              <div className="mt-3 flex items-center gap-2 text-xs uppercase tracking-[0.18em] text-slate-500">
                <CircleDashed size={14} />
                <span>{new Date(alert.timestamp).toLocaleString()}</span>
              </div>
            </article>
          ))
        )}
      </div>
    </aside>
  );
}
