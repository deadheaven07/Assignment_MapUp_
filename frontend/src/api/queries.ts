import { api } from './client';
import type { AlertRule, CoordinatePair, Geofence, Vehicle, VehicleLocation, Violation } from '../types';

export const getGeofences = async () => (await api.get<{ geofences: Geofence[] }>('/geofences')).data.geofences;
export const createGeofence = async (payload: { name: string; description: string; coordinates: CoordinatePair[]; category: string }) =>
  (await api.post('/geofences', payload)).data;

export const getVehicles = async () => (await api.get<{ vehicles: Vehicle[] }>('/vehicles')).data.vehicles;
export const createVehicle = async (payload: { vehicle_number: string; driver_name: string; vehicle_type: string; phone: string }) => (await api.post('/vehicles', payload)).data;

export const postLocation = async (payload: { vehicle_id: string; latitude: number; longitude: number; timestamp: string }) =>
  (await api.post<{ vehicle_id: string; location_updated: boolean; current_geofences: { geofence_id: string; geofence_name: string; status: string }[] }>('/vehicles/location', payload)).data;

export const getLatestLocation = async (vehicleId: string) =>
  (await api.get<{ current_location: VehicleLocation }>(`/vehicles/location/${vehicleId}`)).data.current_location;

export const getAlerts = async () => (await api.get<{ alerts: AlertRule[] }>('/alerts')).data.alerts;
export const createAlert = async (payload: { vehicle_id?: string; geofence_id: string; event_type: 'entry' | 'exit' | 'both' }) =>
  (await api.post('/alerts/configure', payload)).data;

export const getViolations = async (params: Record<string, string | number | undefined>) =>
  (await api.get<{ violations: Violation[]; total_count: number }>('/violations/history', { params })).data;
