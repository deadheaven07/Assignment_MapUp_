import { FormEvent, useMemo, useState } from 'react';
import { useMutation, useQueries, useQuery, useQueryClient } from '@tanstack/react-query';
import { getGeofences, getLatestLocation, getVehicles, postLocation } from '../api/queries';
import { MapPanel } from '../components/MapPanel';
import type { VehicleLocation } from '../types';

export function LocationUpdates() {
  const client = useQueryClient();
  const { data: vehicles = [], isLoading: vehiclesLoading } = useQuery({ queryKey: ['vehicles'], queryFn: getVehicles });
  const { data: geofences = [], isLoading: geofencesLoading } = useQuery({ queryKey: ['geofences'], queryFn: getGeofences });
  const [vehicleId, setVehicleId] = useState('');
  const [latitude, setLatitude] = useState('');
  const [longitude, setLongitude] = useState('');
  const latest = useQueries({ queries: vehicles.map((vehicle) => ({ queryKey: ['latest-location', vehicle.id], queryFn: () => getLatestLocation(vehicle.id), retry: false })) });
  const locations = latest.map((query) => query.data).filter((location): location is VehicleLocation => Boolean(location));
  const mutation = useMutation({ mutationFn: postLocation, onSuccess: () => client.invalidateQueries({ queryKey: ['latest-location'] }) });
  const activeNames = useMemo(() => mutation.data?.current_geofences.map((fence) => fence.geofence_name).join(', ') || 'None', [mutation.data]);

  const submit = (event: FormEvent) => {
    event.preventDefault();
    mutation.mutate({ vehicle_id: vehicleId, latitude: Number(latitude), longitude: Number(longitude), timestamp: new Date().toISOString() });
  };

  return (
    <section className="space-y-6">
      <MapPanel geofences={geofences} locations={locations} onPick={(point) => { setLatitude(String(point.latitude)); setLongitude(String(point.longitude)); }} />
      <form onSubmit={submit} className="grid gap-4 rounded-3xl border border-slate-200 bg-white p-5 shadow-sm md:grid-cols-[1.5fr_1fr_1fr_auto]">
        <select className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" value={vehicleId} onChange={(e) => setVehicleId(e.target.value)}>
          <option value="">Vehicle</option>
          {vehicles.map((vehicle) => <option key={vehicle.id} value={vehicle.id}>{vehicle.vehicle_number}</option>)}
        </select>
        <input className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" placeholder="Latitude" value={latitude} onChange={(e) => setLatitude(e.target.value)} />
        <input className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" placeholder="Longitude" value={longitude} onChange={(e) => setLongitude(e.target.value)} />
        <button className="rounded-3xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white transition hover:bg-slate-800" type="submit">Send</button>
      </form>
      <div className="grid gap-4 lg:grid-cols-[1fr_1fr]">
        <div className="rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-sm uppercase tracking-[0.22em] text-slate-500">Live positions</p>
          <p className="mt-3 text-3xl font-semibold text-slate-900">{locations.length}</p>
          <p className="mt-2 text-sm text-slate-500">Active vehicles with current location data on the map.</p>
        </div>
        <div className="rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-sm uppercase tracking-[0.22em] text-slate-500">Geofence status</p>
          <p className="mt-3 text-base font-semibold text-slate-900">{activeNames}</p>
          <p className="mt-2 text-sm text-slate-500">Last location update will show which geofence the vehicle is currently inside.</p>
        </div>
      </div>
    </section>
  );
}
