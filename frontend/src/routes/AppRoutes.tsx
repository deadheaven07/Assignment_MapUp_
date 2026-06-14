import { createBrowserRouter } from 'react-router-dom';
import { AppLayout } from '../layouts/AppLayout';
import { Alerts } from '../pages/Alerts';
import { Dashboard } from '../pages/Dashboard';
import { Geofences } from '../pages/Geofences';
import { LocationUpdates } from '../pages/LocationUpdates';
import { Vehicles } from '../pages/Vehicles';
import { Violations } from '../pages/Violations';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <AppLayout />,
    children: [
      { index: true, element: <Dashboard /> },
      { path: 'geofences', element: <Geofences /> },
      { path: 'vehicles', element: <Vehicles /> },
      { path: 'locations', element: <LocationUpdates /> },
      { path: 'alerts', element: <Alerts /> },
      { path: 'violations', element: <Violations /> },
    ],
  },
]);
