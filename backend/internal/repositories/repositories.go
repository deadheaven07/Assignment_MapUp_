package repositories

import (
	"time"

	"geofencing-alerts/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateGeofence(geofence *models.Geofence) error {
	return r.db.Create(geofence).Error
}

func (r *Repository) ListGeofences(category string) ([]models.Geofence, error) {
	var geofences []models.Geofence
	query := r.db.Order("created_at desc")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Find(&geofences).Error
	return geofences, err
}

func (r *Repository) CreateVehicle(vehicle *models.Vehicle) error {
	return r.db.Create(vehicle).Error
}

func (r *Repository) ListVehicles() ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	err := r.db.Order("created_at desc").Find(&vehicles).Error
	return vehicles, err
}

func (r *Repository) FindVehicle(id uint) (models.Vehicle, error) {
	var vehicle models.Vehicle
	err := r.db.First(&vehicle, id).Error
	return vehicle, err
}

func (r *Repository) SaveLocation(location *models.VehicleLocation) error {
	return r.db.Create(location).Error
}

func (r *Repository) LatestVehicleLocation(vehicleID uint) (models.VehicleLocation, error) {
	var location models.VehicleLocation
	err := r.db.Where("vehicle_id = ?", vehicleID).Order("timestamp desc").First(&location).Error
	return location, err
}

func (r *Repository) CreateAlertRule(rule *models.AlertRule) error {
	return r.db.Create(rule).Error
}

func (r *Repository) ListAlertRules(vehicleID *uint, geofenceID *uint) ([]models.AlertRule, error) {
	var rules []models.AlertRule
	query := r.db.Order("created_at desc")
	if vehicleID != nil {
		query = query.Where("vehicle_id = ?", *vehicleID)
	}
	if geofenceID != nil {
		query = query.Where("geofence_id = ?", *geofenceID)
	}
	err := query.Find(&rules).Error
	return rules, err
}

func (r *Repository) MatchingRules(vehicleID uint, geofenceID uint, eventType string) ([]models.AlertRule, error) {
	var rules []models.AlertRule
	err := r.db.Where("enabled = ? AND event_type IN ? AND (vehicle_id IS NULL OR vehicle_id = ?) AND geofence_id = ?",
		true, []string{eventType, models.EventBoth}, vehicleID, geofenceID).Find(&rules).Error
	return rules, err
}

func (r *Repository) UpsertState(vehicleID uint, geofenceID uint, inside bool) error {
	state := models.VehicleGeofenceState{VehicleID: vehicleID, GeofenceID: geofenceID, Inside: inside}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "vehicle_id"}, {Name: "geofence_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"inside", "updated_at"}),
	}).Create(&state).Error
}

func (r *Repository) GetState(vehicleID uint, geofenceID uint) (models.VehicleGeofenceState, bool, error) {
	var state models.VehicleGeofenceState
	err := r.db.Where("vehicle_id = ? AND geofence_id = ?", vehicleID, geofenceID).First(&state).Error
	if err == gorm.ErrRecordNotFound {
		return state, false, nil
	}
	return state, true, err
}

func (r *Repository) CreateViolation(violation *models.Violation) error {
	return r.db.Create(violation).Error
}

func (r *Repository) CreateAlertEvent(event *models.AlertEvent) error {
	return r.db.Create(event).Error
}

type ViolationFilter struct {
	VehicleID  *uint
	GeofenceID *uint
	From       *time.Time
	To         *time.Time
	Limit      int
	Offset     int
}

func (r *Repository) ListViolations(filter ViolationFilter) ([]models.Violation, int64, error) {
	query := r.db.Model(&models.Violation{}).Preload("Vehicle").Preload("Geofence")
	if filter.VehicleID != nil {
		query = query.Where("vehicle_id = ?", *filter.VehicleID)
	}
	if filter.GeofenceID != nil {
		query = query.Where("geofence_id = ?", *filter.GeofenceID)
	}
	if filter.From != nil {
		query = query.Where("timestamp >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("timestamp <= ?", *filter.To)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var violations []models.Violation
	err := query.Order("timestamp desc").Offset(filter.Offset).Limit(filter.Limit).Find(&violations).Error
	return violations, total, err
}
