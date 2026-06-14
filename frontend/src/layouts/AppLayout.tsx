import { AlertTriangle, Bell, Home, Map, MapPin, Truck } from 'lucide-react';
import { NavLink, Outlet } from 'react-router-dom';
import { LiveAlerts } from '../components/LiveAlerts';
import { ToastStack } from '../components/ToastStack';
import { useLiveAlerts } from '../hooks/useLiveAlerts';

const links = [
  { to: '/', label: 'Dashboard', Icon: Home },
  { to: '/geofences', label: 'Geofences', Icon: MapPin },
  { to: '/vehicles', label: 'Vehicles', Icon: Truck },
  { to: '/locations', label: 'Location Updates', Icon: Map },
  { to: '/alerts', label: 'Alerts', Icon: Bell },
  { to: '/violations', label: 'Violations', Icon: AlertTriangle },
];

export function AppLayout() {
  const { alerts, connected } = useLiveAlerts();

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900">
      <ToastStack alerts={alerts} />
      <header className="border-b border-slate-200 bg-white/95 backdrop-blur">
        <div className="mx-auto flex max-w-7xl flex-col gap-4 px-5 py-5 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <h1 className="text-xl font-semibold tracking-tight text-slate-900">Geofencing Alerts</h1>
          </div>
          <nav className="flex flex-wrap gap-2">
            {links.map(({ to, label, Icon }) => (
              <NavLink
                key={to}
                to={to}
                className={({ isActive }) =>
                  `group inline-flex items-center gap-2 rounded-full px-3 py-2 text-sm font-medium transition ${
                    isActive ? 'bg-slate-900 text-white shadow-sm' : 'text-slate-600 hover:bg-slate-100'
                  }`
                }
              >
                <Icon size={16} className="shrink-0" />
                {label}
              </NavLink>
            ))}
          </nav>
        </div>
      </header>
      <main className="mx-auto grid max-w-7xl gap-6 px-5 py-6 lg:grid-cols-[1fr_320px]">
        <Outlet />
        <LiveAlerts alerts={alerts} connected={connected} />
      </main>
    </div>
  );
}
