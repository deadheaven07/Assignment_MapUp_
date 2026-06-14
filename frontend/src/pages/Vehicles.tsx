import { FormEvent, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { createVehicle, getVehicles } from '../api/queries';

export function Vehicles() {
  const client = useQueryClient();
  const { data: vehicles = [], isLoading: vehiclesLoading } = useQuery({ queryKey: ['vehicles'], queryFn: getVehicles });
  const [vehicleNumber, setVehicleNumber] = useState('');
  const [driverName, setDriverName] = useState('');
  const [vehicleType, setVehicleType] = useState('');
  const [phone, setPhone] = useState('');
  const mutation = useMutation({ mutationFn: createVehicle, onSuccess: () => client.invalidateQueries({ queryKey: ['vehicles'] }) });
  const submit = (event: FormEvent) => {
    event.preventDefault();
    mutation.mutate({ vehicle_number: vehicleNumber, driver_name: driverName, vehicle_type: vehicleType, phone }, { onSuccess: () => { setVehicleNumber(''); setDriverName(''); setVehicleType(''); setPhone(''); } });
  };

  return (
    <section className="space-y-6">
      <form onSubmit={submit} className="grid gap-4 rounded-3xl border border-slate-200 bg-white p-5 shadow-sm md:grid-cols-[1fr_1fr_1fr_1fr_auto]">
        <input className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" placeholder="Vehicle number" value={vehicleNumber} onChange={(e) => setVehicleNumber(e.target.value)} />
        <input className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" placeholder="Driver name" value={driverName} onChange={(e) => setDriverName(e.target.value)} />
        <input className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" placeholder="Vehicle type" value={vehicleType} onChange={(e) => setVehicleType(e.target.value)} />
        <input className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none" placeholder="Phone" value={phone} onChange={(e) => setPhone(e.target.value)} />
        <button className="rounded-3xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white transition hover:bg-slate-800" type="submit">Create</button>
      </form>
      <div className="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm">
        <div className="grid gap-4 border-b border-slate-200 bg-slate-50 px-4 py-4 text-xs uppercase tracking-[0.22em] text-slate-500 sm:grid-cols-[1fr_1fr_1fr_1fr]">
          <span>Number</span>
          <span>Driver</span>
          <span>Type</span>
          <span>Phone</span>
        </div>
        {vehiclesLoading ? (
          <div className="space-y-3 p-4">
            {[...Array(4)].map((_, index) => (
              <div key={index} className="h-14 rounded-3xl bg-slate-100" />
            ))}
          </div>
        ) : vehicles.length === 0 ? (
          <div className="p-10 text-center text-slate-500">No vehicles registered yet. Add fleet vehicles to start tracking locations.</div>
        ) : (
          vehicles.map((vehicle, index) => (
            <div key={vehicle.id} className={`grid min-h-[72px] items-center gap-4 px-4 py-4 text-sm ${index % 2 === 0 ? 'bg-white' : 'bg-slate-50'} sm:grid-cols-[1fr_1fr_1fr_1fr] transition hover:bg-slate-100`}>
              <span>{vehicle.vehicle_number}</span>
              <span>{vehicle.driver_name}</span>
              <span>{vehicle.vehicle_type}</span>
              <span>{vehicle.phone}</span>
            </div>
          ))
        )}
      </div>
    </section>
  );
}
