import { FormEvent, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { createGeofence, getGeofences } from '../api/queries';
import { MapPanel } from '../components/MapPanel';
import type { Coordinate } from '../types';

export function Geofences() {
  const client = useQueryClient();
  const { data: geofences = [], isLoading: geofencesLoading } = useQuery({ queryKey: ['geofences'], queryFn: getGeofences });
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [category, setCategory] = useState('delivery_zone');
  const [draft, setDraft] = useState<Coordinate[]>([]);
  const mutation = useMutation({ mutationFn: createGeofence, onSuccess: () => client.invalidateQueries({ queryKey: ['geofences'] }) });

  const closePolygon = () => setDraft((points) => (points.length ? [...points, points[0]] : points));
  const submit = (event: FormEvent) => {
    event.preventDefault();
    mutation.mutate({ name, description, category, coordinates: draft.map((point) => [point.latitude, point.longitude]) }, { onSuccess: () => { setName(''); setDescription(''); setDraft([]); } });
  };

  return (
    <section className="space-y-6">
      <MapPanel geofences={geofences} draft={draft} onPick={(point) => setDraft((points) => [...points, point])} />
      <form onSubmit={submit} className="grid gap-4 rounded-3xl border border-slate-200 bg-white p-5 shadow-sm md:grid-cols-[1.2fr_1fr_1fr_auto_auto]">
        <input className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" placeholder="Name" value={name} onChange={(e) => setName(e.target.value)} />
        <input className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" placeholder="Description" value={description} onChange={(e) => setDescription(e.target.value)} />
        <select className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" value={category} onChange={(e) => setCategory(e.target.value)}>
          <option value="delivery_zone">delivery_zone</option>
          <option value="restricted_zone">restricted_zone</option>
          <option value="toll_zone">toll_zone</option>
          <option value="customer_area">customer_area</option>
        </select>
        <button type="button" className="rounded-3xl border border-slate-200 bg-slate-50 px-5 py-3 text-sm font-semibold text-slate-700 transition hover:border-slate-300" onClick={closePolygon}>Close polygon</button>
        <button className="rounded-3xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white transition hover:bg-slate-800">Create</button>
      </form>
      <div className="grid gap-4 lg:grid-cols-2">
        {geofencesLoading ? (
          [...Array(4)].map((_, index) => (
            <div key={index} className="h-40 rounded-3xl bg-slate-100" />
          ))
        ) : geofences.length === 0 ? (
          <div className="rounded-3xl border border-dashed border-slate-200 bg-slate-50 p-10 text-center text-slate-500">
            No geofences yet. Use the map to draw a polygon and create your first zone.
          </div>
        ) : (
          geofences.map((fence) => (
            <div key={fence.id} className="rounded-3xl border border-slate-200 bg-white p-5 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
              <div className="flex items-center justify-between gap-3">
                <h3 className="text-lg font-semibold text-slate-900">{fence.name}</h3>
                <span className="text-sm uppercase tracking-[0.18em] text-slate-500">{fence.category}</span>
              </div>
              <p className="mt-3 text-sm leading-6 text-slate-600">{fence.description || 'No description provided.'}</p>
              <div className="mt-4 flex flex-wrap gap-2 text-sm text-slate-500">
                <span>{fence.coordinates.length} points</span>
                <span>Created on {new Date(fence.created_at).toLocaleDateString()}</span>
              </div>
            </div>
          ))
        )}
      </div>
    </section>
  );
}
