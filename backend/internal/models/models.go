package models

import (
	"time"

	"gorm.io/datatypes"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Geofence struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Polygon     datatypes.JSON `gorm:"type:jsonb;not null" json:"-"`
	Category    string         `gorm:"not null" json:"category"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type Vehicle struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	VehicleNumber string    `gorm:"uniqueIndex;not null" json:"vehicle_number"`
	DriverName    string    `gorm:"not null" json:"driver_name"`
	VehicleType   string    `gorm:"not null" json:"vehicle_type"`
	Phone         string    `gorm:"not null" json:"phone"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type VehicleLocation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VehicleID uint      `gorm:"not null;index" json:"vehicle_id"`
	Latitude  float64   `gorm:"not null" json:"latitude"`
	Longitude float64   `gorm:"not null" json:"longitude"`
	Timestamp time.Time `gorm:"not null;index" json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

type AlertRule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VehicleID   *uint     `gorm:"index" json:"vehicle_id"`
	GeofenceID  *uint     `gorm:"index" json:"geofence_id"`
	EventType   string    `gorm:"not null" json:"event_type"`
	Enabled     bool      `gorm:"not null;default:true" json:"enabled"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Violation struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	VehicleID  uint      `gorm:"not null;index" json:"vehicle_id"`
	GeofenceID uint      `gorm:"not null;index" json:"geofence_id"`
	EventType  string    `gorm:"not null;index" json:"event_type"`
	Latitude   float64   `gorm:"not null" json:"latitude"`
	Longitude  float64   `gorm:"not null" json:"longitude"`
	Timestamp  time.Time `gorm:"not null;index" json:"timestamp"`
	CreatedAt  time.Time `json:"created_at"`
	Vehicle    Vehicle   `gorm:"foreignKey:VehicleID" json:"vehicle"`
	Geofence   Geofence  `gorm:"foreignKey:GeofenceID" json:"geofence"`
}

type VehicleGeofenceState struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	VehicleID  uint      `gorm:"not null;uniqueIndex:idx_vehicle_geofence" json:"vehicle_id"`
	GeofenceID uint      `gorm:"not null;uniqueIndex:idx_vehicle_geofence" json:"geofence_id"`
	Inside     bool      `gorm:"not null" json:"inside"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AlertEvent struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AlertRuleID *uint     `gorm:"index" json:"alert_rule_id"`
	ViolationID *uint     `gorm:"index" json:"violation_id"`
	VehicleID   uint      `gorm:"not null;index" json:"vehicle_id"`
	GeofenceID  uint      `gorm:"not null;index" json:"geofence_id"`
	EventType   string    `gorm:"not null" json:"event_type"`
	Latitude    float64   `gorm:"not null" json:"latitude"`
	Longitude   float64   `gorm:"not null" json:"longitude"`
	Timestamp   time.Time `gorm:"not null;index" json:"timestamp"`
	CreatedAt   time.Time `json:"created_at"`
}

const (
	EventEntry = "entry"
	EventExit  = "exit"
	EventBoth  = "both"
)
