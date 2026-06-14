import { MapContainer, Marker, Polygon, Popup, TileLayer, useMapEvents } from 'react-leaflet';
import type { Coordinate, Geofence, VehicleLocation } from '../types';

type Props = {
  geofences: Geofence[];
  locations?: VehicleLocation[];
  draft?: Coordinate[];
  onPick?: (point: Coordinate) => void;
};

function ClickPicker({ onPick }: { onPick?: (point: Coordinate) => void }) {
  useMapEvents({
    click(event) {
      onPick?.({ latitude: event.latlng.lat, longitude: event.latlng.lng });
    },
  });
  return null;
}

export function MapPanel({ geofences, locations = [], draft = [], onPick }: Props) {
  return (
    <MapContainer center={[28.6139, 77.209]} zoom={11} className="h-[420px] w-full overflow-hidden rounded border border-slate-200">
      <TileLayer attribution="&copy; OpenStreetMap contributors" url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
      <ClickPicker onPick={onPick} />
      {geofences.map((fence) => (
        <Polygon key={fence.id} positions={fence.coordinates.map((point) => [point[0], point[1]])} pathOptions={{ color: '#0f766e' }}>
          <Popup>{fence.name}</Popup>
        </Polygon>
      ))}
      {draft.length > 1 && <Polygon positions={draft.map((point) => [point.latitude, point.longitude])} pathOptions={{ color: '#dc2626' }} />}
      {locations.map((location) => (
        <Marker key={`${location.vehicle_id}-${location.id}`} position={[location.latitude, location.longitude]}>
          <Popup>Vehicle #{location.vehicle_id}</Popup>
        </Marker>
      ))}
    </MapContainer>
  );
}
