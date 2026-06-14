export type Coordinate = { latitude: number; longitude: number };
export type CoordinatePair = [number, number];

export type Geofence = {
  id: string;
  name: string;
  description: string;
  coordinates: CoordinatePair[];
  category: string;
  created_at: string;
};

export type Vehicle = {
  id: string;
  vehicle_number: string;
  driver_name: string;
  vehicle_type: string;
  phone: string;
  status: string;
  created_at: string;
};

export type VehicleLocation = {
  id: number;
  vehicle_id: string;
  latitude: number;
  longitude: number;
  timestamp: string;
};

export type AlertRule = {
  alert_id: string;
  vehicle_id?: string | null;
  geofence_id: string;
  geofence_name?: string;
  vehicle_number?: string;
  event_type: 'entry' | 'exit' | 'both';
  status: string;
  created_at: string;
};

export type Violation = {
  id: string;
  vehicle_id: string;
  vehicle_number: string;
  geofence_id: string;
  geofence_name: string;
  event_type: 'entry' | 'exit';
  latitude: number;
  longitude: number;
  timestamp: string;
};

export type AlertEvent = {
  event_id: string;
  event_type: 'entry' | 'exit';
  timestamp: string;
  vehicle: {
    vehicle_id: string;
    vehicle_number: string;
    driver_name: string;
  };
  geofence: {
    geofence_id: string;
    geofence_name: string;
    category: string;
  };
  location: {
    latitude: number;
    longitude: number;
  };
};
