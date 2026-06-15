package services

import (
	"encoding/json"
	"testing"
	"time"

	"geofencing-alerts/backend/internal/geofence"
	"geofencing-alerts/backend/internal/models"
	"geofencing-alerts/backend/internal/repositories"
	"geofencing-alerts/backend/internal/websocket"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testPublisher struct {
	published []websocket.AlertEvent
}

func (p *testPublisher) Publish(event websocket.AlertEvent) {
	p.published = append(p.published, event)
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Geofence{}, &models.Vehicle{}, &models.VehicleLocation{}, &models.AlertRule{}, &models.Violation{}, &models.VehicleGeofenceState{}, &models.AlertEvent{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func TestUpdateLocationTriggersEntryAndExit(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.New(db)
	publisher := &testPublisher{}
	svc := New(repo, publisher)

	var err error
	vehicle := models.Vehicle{VehicleNumber: "V1", DriverName: "Test", VehicleType: "truck", Phone: "123"}
	if err = repo.CreateVehicle(&vehicle); err != nil {
		t.Fatal(err)
	}
	geofenceData := []models.Coordinate{{Latitude: 0, Longitude: 0}, {Latitude: 0, Longitude: 10}, {Latitude: 10, Longitude: 10}, {Latitude: 0, Longitude: 0}}
	if err = geofence.ValidatePolygon(geofenceData); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(geofenceData)
	fence := models.Geofence{Name: "Zone", Description: "Test zone", Polygon: datatypes.JSON(raw), Category: "delivery_zone"}
	if err = repo.CreateGeofence(&fence); err != nil {
		t.Fatal(err)
	}
	rule := models.AlertRule{GeofenceID: &fence.ID, EventType: models.EventEntry, Enabled: true}
	if err = repo.CreateAlertRule(&rule); err != nil {
		t.Fatal(err)
	}

	_, err = svc.UpdateLocation(LocationInput{VehicleID: vehicle.ID, Latitude: 5, Longitude: 5, Timestamp: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}

	if len(publisher.published) != 1 {
		t.Fatalf("expected 1 published alert, got %d", len(publisher.published))
	}
	if publisher.published[0].EventType != models.EventEntry {
		t.Fatalf("expected entry event, got %s", publisher.published[0].EventType)
	}

	_, err = svc.UpdateLocation(LocationInput{VehicleID: vehicle.ID, Latitude: 15, Longitude: 15, Timestamp: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}

	if len(publisher.published) != 1 {
		t.Fatalf("expected only 1 published alert because exit rule not configured, got %d", len(publisher.published))
	}
}
