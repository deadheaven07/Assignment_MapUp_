package geofence

import (
	"testing"

	"geofencing-alerts/backend/internal/models"
)

func TestValidatePolygon(t *testing.T) {
	valid := []models.Coordinate{{Latitude: 0, Longitude: 0}, {Latitude: 0, Longitude: 1}, {Latitude: 1, Longitude: 1}, {Latitude: 0, Longitude: 0}}
	if err := ValidatePolygon(valid); err != nil {
		t.Fatalf("expected valid polygon, got %v", err)
	}

	if err := ValidatePolygon([]models.Coordinate{{Latitude: 0, Longitude: 0}, {Latitude: 0, Longitude: 1}, {Latitude: 1, Longitude: 1}}); err == nil {
		t.Fatal("expected invalid polygon with fewer than 4 points")
	}

	if err := ValidatePolygon([]models.Coordinate{{Latitude: 0, Longitude: 0}, {Latitude: 0, Longitude: 1}, {Latitude: 1, Longitude: 1}, {Latitude: 1, Longitude: 2}}); err == nil {
		t.Fatal("expected invalid polygon because it is not closed")
	}

	if err := ValidatePolygon([]models.Coordinate{{Latitude: 100, Longitude: 0}, {Latitude: 0, Longitude: 1}, {Latitude: 1, Longitude: 1}, {Latitude: 100, Longitude: 0}}); err == nil {
		t.Fatal("expected invalid latitude out of range")
	}
}

func TestContains(t *testing.T) {
	polygon := []models.Coordinate{{Latitude: 0, Longitude: 0}, {Latitude: 0, Longitude: 10}, {Latitude: 10, Longitude: 10}, {Latitude: 0, Longitude: 0}}

	inside := Contains(models.Coordinate{Latitude: 5, Longitude: 5}, polygon)
	if !inside {
		t.Fatal("expected point inside polygon")
	}

	outside := Contains(models.Coordinate{Latitude: 15, Longitude: 5}, polygon)
	if outside {
		t.Fatal("expected point outside polygon")
	}

	nearBoundary := Contains(models.Coordinate{Latitude: 0.0001, Longitude: 5}, polygon)
	if !nearBoundary {
		t.Fatal("expected point just inside boundary to be inside polygon")
	}
}
