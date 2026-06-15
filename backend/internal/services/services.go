package services

import (
	"encoding/json"
	"errors"
	"time"

	"geofencing-alerts/backend/internal/geofence"
	"geofencing-alerts/backend/internal/models"
	"geofencing-alerts/backend/internal/repositories"
	"geofencing-alerts/backend/internal/websocket"

	"gorm.io/datatypes"
)

type Service struct {
	repo      *repositories.Repository
	publisher websocket.Publisher
}

func New(repo *repositories.Repository, publisher websocket.Publisher) *Service {
	return &Service{repo: repo, publisher: publisher}
}

type CreateGeofenceInput struct {
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description"`
	Coordinates [][]float64 `json:"coordinates" binding:"required"`
	Category    string      `json:"category" binding:"required"`
}

func (s *Service) CreateGeofence(input CreateGeofenceInput) (models.Geofence, error) {
	points := make([]models.Coordinate, 0, len(input.Coordinates))
	for _, pair := range input.Coordinates {
		if len(pair) != 2 {
			return models.Geofence{}, errors.New("coordinates must contain [latitude, longitude] pairs")
		}
		points = append(points, models.Coordinate{Latitude: pair[0], Longitude: pair[1]})
	}
	if err := validateCategory(input.Category); err != nil {
		return models.Geofence{}, err
	}
	if err := geofence.ValidatePolygon(points); err != nil {
		return models.Geofence{}, err
	}
	raw, _ := json.Marshal(points)
	item := models.Geofence{Name: input.Name, Description: input.Description, Polygon: datatypes.JSON(raw), Category: input.Category}
	return item, s.repo.CreateGeofence(&item)
}

func (s *Service) ListGeofences(category string) ([]models.Geofence, error) {
	return s.repo.ListGeofences(category)
}

type CreateVehicleInput struct {
	VehicleNumber string `json:"vehicle_number" binding:"required"`
	DriverName    string `json:"driver_name" binding:"required"`
	VehicleType   string `json:"vehicle_type" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
}

func (s *Service) CreateVehicle(input CreateVehicleInput) (models.Vehicle, error) {
	vehicle := models.Vehicle{VehicleNumber: input.VehicleNumber, DriverName: input.DriverName, VehicleType: input.VehicleType, Phone: input.Phone}
	return vehicle, s.repo.CreateVehicle(&vehicle)
}

func (s *Service) ListVehicles() ([]models.Vehicle, error) {
	return s.repo.ListVehicles()
}

type LocationInput struct {
	VehicleID uint      `json:"vehicle_id" binding:"required"`
	Latitude  float64   `json:"latitude" binding:"required"`
	Longitude float64   `json:"longitude" binding:"required"`
	Timestamp time.Time `json:"timestamp"`
}

type LocationResult struct {
	Location        models.VehicleLocation `json:"location"`
	ActiveGeofences []models.Geofence      `json:"active_geofences"`
}

func (s *Service) UpdateLocation(input LocationInput) (LocationResult, error) {
	if input.Latitude < -90 || input.Latitude > 90 || input.Longitude < -180 || input.Longitude > 180 {
		return LocationResult{}, errors.New("invalid latitude or longitude")
	}
	if input.Timestamp.IsZero() {
		input.Timestamp = time.Now().UTC()
	}

	vehicle, err := s.repo.FindVehicle(input.VehicleID)
	if err != nil {
		return LocationResult{}, err
	}

	location := models.VehicleLocation{
		VehicleID: input.VehicleID,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Timestamp: input.Timestamp,
	}
	if err := s.repo.SaveLocation(&location); err != nil {
		return LocationResult{}, err
	}

	geofences, err := s.repo.ListGeofences("")
	if err != nil {
		return LocationResult{}, err
	}

	point := models.Coordinate{Latitude: input.Latitude, Longitude: input.Longitude}
	active := make([]models.Geofence, 0)
	for _, fence := range geofences {
		polygon, err := geofence.ParsePolygon(fence.Polygon)
		if err != nil {
			continue
		}
		inside := geofence.Contains(point, polygon)
		if inside {
			active = append(active, fence)
		}

		prev, found, err := s.repo.GetState(input.VehicleID, fence.ID)
		if err != nil {
			return LocationResult{}, err
		}
		if !found {
			if inside {
				eventType := models.EventEntry
				violation := models.Violation{
					VehicleID:  input.VehicleID,
					GeofenceID: fence.ID,
					EventType:  eventType,
					Latitude:   input.Latitude,
					Longitude:  input.Longitude,
					Timestamp:  input.Timestamp,
				}
				if err := s.repo.CreateViolation(&violation); err != nil {
					return LocationResult{}, err
				}
				rules, err := s.repo.MatchingRules(input.VehicleID, fence.ID, eventType)
				if err != nil {
					return LocationResult{}, err
				}
				for _, rule := range rules {
					event := models.AlertEvent{
						AlertRuleID: &rule.ID,
						ViolationID: &violation.ID,
						VehicleID:   input.VehicleID,
						GeofenceID:  fence.ID,
						EventType:   eventType,
						Latitude:    input.Latitude,
						Longitude:   input.Longitude,
						Timestamp:   input.Timestamp,
					}
					if err := s.repo.CreateAlertEvent(&event); err != nil {
						return LocationResult{}, err
					}
					s.publisher.Publish(websocket.NewAlertEvent(event.ID, eventType, input.Timestamp, vehicle, fence, location))
				}
			}
			if err := s.repo.UpsertState(input.VehicleID, fence.ID, inside); err != nil {
				return LocationResult{}, err
			}
			continue
		}
		if prev.Inside == inside {
			continue
		}

		eventType := models.EventExit
		if inside {
			eventType = models.EventEntry
		}
		violation := models.Violation{
			VehicleID:  input.VehicleID,
			GeofenceID: fence.ID,
			EventType:  eventType,
			Latitude:   input.Latitude,
			Longitude:  input.Longitude,
			Timestamp:  input.Timestamp,
		}
		if err := s.repo.CreateViolation(&violation); err != nil {
			return LocationResult{}, err
		}
		if err := s.repo.UpsertState(input.VehicleID, fence.ID, inside); err != nil {
			return LocationResult{}, err
		}
		rules, err := s.repo.MatchingRules(input.VehicleID, fence.ID, eventType)
		if err != nil {
			return LocationResult{}, err
		}
		for _, rule := range rules {
			event := models.AlertEvent{
				AlertRuleID: &rule.ID,
				ViolationID: &violation.ID,
				VehicleID:   input.VehicleID,
				GeofenceID:  fence.ID,
				EventType:   eventType,
				Latitude:    input.Latitude,
				Longitude:   input.Longitude,
				Timestamp:   input.Timestamp,
			}
			if err := s.repo.CreateAlertEvent(&event); err != nil {
				return LocationResult{}, err
			}
			s.publisher.Publish(websocket.NewAlertEvent(event.ID, eventType, input.Timestamp, vehicle, fence, location))
		}
	}

	return LocationResult{Location: location, ActiveGeofences: active}, nil
}

func (s *Service) LatestLocation(vehicleID uint) (models.VehicleLocation, error) {
	return s.repo.LatestVehicleLocation(vehicleID)
}

type AlertRuleInput struct {
	VehicleID   *uint  `json:"vehicle_id"`
	GeofenceID  *uint  `json:"geofence_id"`
	EventType   string `json:"event_type" binding:"required"`
	Enabled     *bool  `json:"enabled"`
	Description string `json:"description"`
}

func (s *Service) CreateAlertRule(input AlertRuleInput) (models.AlertRule, error) {
	if input.GeofenceID == nil {
		return models.AlertRule{}, errors.New("geofence_id is required")
	}
	if input.EventType != models.EventEntry && input.EventType != models.EventExit && input.EventType != models.EventBoth {
		return models.AlertRule{}, errors.New("event_type must be entry, exit, or both")
	}
	enabled := true
	if input.Enabled != nil {
		enabled = *input.Enabled
	}
	rule := models.AlertRule{VehicleID: input.VehicleID, GeofenceID: input.GeofenceID, EventType: input.EventType, Enabled: enabled, Description: input.Description}
	return rule, s.repo.CreateAlertRule(&rule)
}

func (s *Service) ListAlertRules(vehicleID *uint, geofenceID *uint) ([]models.AlertRule, error) {
	return s.repo.ListAlertRules(vehicleID, geofenceID)
}

func (s *Service) ListViolations(filter repositories.ViolationFilter) ([]models.Violation, int64, error) {
	return s.repo.ListViolations(filter)
}

func validateCategory(category string) error {
	switch category {
	case "delivery_zone", "restricted_zone", "toll_zone", "customer_area":
		return nil
	default:
		return errors.New("category must be delivery_zone, restricted_zone, toll_zone, or customer_area")
	}
}
