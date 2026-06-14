package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"geofencing-alerts/backend/internal/models"
)

type coordinatePair [2]float64

func publicID(prefix string, id uint) string {
	return fmt.Sprintf("%s_%d", prefix, id)
}

func parsePublicID(value string, prefix string) (uint, error) {
	raw, ok := strings.CutPrefix(value, prefix+"_")
	if !ok {
		return 0, fmt.Errorf("invalid %s id", prefix)
	}
	id, err := strconv.ParseUint(raw, 10, 64)
	return uint(id), err
}

func coordinatesFromPolygon(raw []byte) []coordinatePair {
	var points []models.Coordinate
	if err := json.Unmarshal(raw, &points); err != nil {
		return nil
	}
	coordinates := make([]coordinatePair, 0, len(points))
	for _, point := range points {
		coordinates = append(coordinates, coordinatePair{point.Latitude, point.Longitude})
	}
	return coordinates
}

func geofenceSummaryDTO(geofence models.Geofence) ginMap {
	return ginMap{
		"geofence_id":   publicID("geo", geofence.ID),
		"geofence_name": geofence.Name,
		"category":      geofence.Category,
	}
}

func geofenceDTO(geofence models.Geofence) ginMap {
	return ginMap{
		"id":          publicID("geo", geofence.ID),
		"name":        geofence.Name,
		"description": geofence.Description,
		"coordinates": coordinatesFromPolygon(geofence.Polygon),
		"category":    geofence.Category,
		"created_at":  geofence.CreatedAt,
	}
}

func vehicleDTO(vehicle models.Vehicle) ginMap {
	return ginMap{
		"id":             publicID("veh", vehicle.ID),
		"vehicle_number": vehicle.VehicleNumber,
		"driver_name":    vehicle.DriverName,
		"vehicle_type":   vehicle.VehicleType,
		"phone":          vehicle.Phone,
		"status":         "active",
		"created_at":     vehicle.CreatedAt,
	}
}

func alertDTO(alert models.AlertRule, geofence *models.Geofence, vehicle *models.Vehicle) ginMap {
	item := ginMap{
		"alert_id":    publicID("alert", alert.ID),
		"geofence_id": nil,
		"vehicle_id":  nil,
		"event_type":  alert.EventType,
		"status":      "active",
		"created_at":  alert.CreatedAt,
	}
	if alert.GeofenceID != nil {
		item["geofence_id"] = publicID("geo", *alert.GeofenceID)
	}
	if alert.VehicleID != nil {
		item["vehicle_id"] = publicID("veh", *alert.VehicleID)
	}
	if geofence != nil {
		item["geofence_name"] = geofence.Name
	}
	if vehicle != nil {
		item["vehicle_number"] = vehicle.VehicleNumber
	}
	return item
}

func violationDTO(violation models.Violation) ginMap {
	return ginMap{
		"id":             publicID("viol", violation.ID),
		"vehicle_id":     publicID("veh", violation.VehicleID),
		"vehicle_number": violation.Vehicle.VehicleNumber,
		"geofence_id":    publicID("geo", violation.GeofenceID),
		"geofence_name":  violation.Geofence.Name,
		"event_type":     violation.EventType,
		"latitude":       violation.Latitude,
		"longitude":      violation.Longitude,
		"timestamp":      violation.Timestamp,
	}
}

func locationDTO(location models.VehicleLocation) ginMap {
	return ginMap{
		"latitude":  location.Latitude,
		"longitude": location.Longitude,
		"timestamp": location.Timestamp,
	}
}

func formatTimeNS(started time.Time) string {
	return strconv.FormatInt(time.Since(started).Nanoseconds(), 10)
}

type ginMap map[string]any
