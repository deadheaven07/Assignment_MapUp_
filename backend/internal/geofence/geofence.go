package geofence

import (
	"encoding/json"
	"errors"

	"geofencing-alerts/backend/internal/models"
)

func ValidatePolygon(points []models.Coordinate) error {
	if len(points) < 4 {
		return errors.New("polygon must contain at least 4 points")
	}

	for _, point := range points {
		if point.Latitude < -90 || point.Latitude > 90 {
			return errors.New("latitude must be between -90 and 90")
		}
		if point.Longitude < -180 || point.Longitude > 180 {
			return errors.New("longitude must be between -180 and 180")
		}
	}

	first := points[0]
	last := points[len(points)-1]
	if first.Latitude != last.Latitude || first.Longitude != last.Longitude {
		return errors.New("polygon must be closed")
	}

	return nil
}

func ParsePolygon(raw []byte) ([]models.Coordinate, error) {
	var points []models.Coordinate
	if err := json.Unmarshal(raw, &points); err != nil {
		return nil, err
	}
	return points, ValidatePolygon(points)
}

func Contains(point models.Coordinate, polygon []models.Coordinate) bool {
	inside := false
	j := len(polygon) - 1

	for i := 0; i < len(polygon); i++ {
		yi := polygon[i].Latitude
		yj := polygon[j].Latitude
		xi := polygon[i].Longitude
		xj := polygon[j].Longitude

		intersects := (yi > point.Latitude) != (yj > point.Latitude)
		if intersects {
			xAtY := (xj-xi)*(point.Latitude-yi)/(yj-yi) + xi
			if point.Longitude < xAtY {
				inside = !inside
			}
		}
		j = i
	}

	return inside
}
