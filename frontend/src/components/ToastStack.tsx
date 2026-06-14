import { Bell, MapPin, Truck } from 'lucide-react';
import type { AlertEvent } from '../types';

export function ToastStack({ alerts }: { alerts: AlertEvent[] }) {
  return (
    <div className="pointer-events-none fixed right-4 top-4 z-[1000] w-72 space-y-3 sm:w-80">
      {alerts.slice(0, 3).map((alert) => (
        <div key={`${alert.event_id}-${alert.timestamp}`} className="pointer-events-auto overflow-hidden rounded-3xl border border-slate-200 bg-white/95 p-4 shadow-2xl shadow-slate-900/5 backdrop-blur transition duration-200 hover:shadow-slate-900/10">
          <div className="flex items-center justify-between gap-3">
            <div className="flex items-center gap-2 text-sm font-semibold text-slate-900">
              <Bell size={16} />
              {alert.event_type.toUpperCase()} ALERT
            </div>
            <span className="text-xs text-slate-500">{new Date(alert.timestamp).toLocaleTimeString()}</span>
          </div>
          <div className="mt-3 flex items-center gap-2 text-sm text-slate-600">
            <Truck size={16} />
            <span>{alert.vehicle.vehicle_number}</span>
          </div>
          <div className="mt-2 flex items-center gap-2 text-sm text-slate-600">
            <MapPin size={16} />
            <span>{alert.geofence.geofence_name}</span>
          </div>
        </div>
      ))}
    </div>
  );
}
