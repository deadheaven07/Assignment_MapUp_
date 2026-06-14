package handlers

import (
	"net/http"
	"strconv"
	"time"

	"geofencing-alerts/backend/internal/geofence"
	"geofencing-alerts/backend/internal/models"
	"geofencing-alerts/backend/internal/repositories"
	"geofencing-alerts/backend/internal/services"
	"geofencing-alerts/backend/internal/websocket"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
	ws      *websocket.Manager
}

func New(service *services.Service, ws *websocket.Manager) *Handler {
	return &Handler{service: service, ws: ws}
}

func (h *Handler) Register(router *gin.Engine) {
	router.POST("/geofences", h.createGeofence)
	router.GET("/geofences", h.listGeofences)
	router.POST("/vehicles", h.createVehicle)
	router.GET("/vehicles", h.listVehicles)
	router.POST("/vehicles/location", h.updateLocation)
	router.GET("/vehicles/location/:vehicle_id", h.latestLocation)
	router.POST("/alerts/configure", h.createAlertRule)
	router.GET("/alerts", h.listAlertRules)
	router.GET("/violations/history", h.listViolations)
	router.GET("/ws/alerts", func(c *gin.Context) { h.ws.Serve(c.Writer, c.Request) })
}

func (h *Handler) createGeofence(c *gin.Context) {
	started := time.Now()
	var input services.CreateGeofenceInput
	if !bindJSON(c, started, &input) {
		return
	}
	geofence, err := h.service.CreateGeofence(input)
	if err != nil {
		fail(c, started, http.StatusBadRequest, err.Error())
		return
	}
	respond(c, started, http.StatusCreated, gin.H{"id": publicID("geo", geofence.ID), "name": geofence.Name, "status": "active"})
}

func (h *Handler) listGeofences(c *gin.Context) {
	started := time.Now()
	geofences, err := h.service.ListGeofences(c.Query("category"))
	if err != nil {
		fail(c, started, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]ginMap, 0, len(geofences))
	for _, geofence := range geofences {
		items = append(items, geofenceDTO(geofence))
	}
	respond(c, started, http.StatusOK, gin.H{"geofences": items})
}

func (h *Handler) createVehicle(c *gin.Context) {
	started := time.Now()
	var input services.CreateVehicleInput
	if !bindJSON(c, started, &input) {
		return
	}
	vehicle, err := h.service.CreateVehicle(input)
	if err != nil {
		fail(c, started, http.StatusBadRequest, err.Error())
		return
	}
	respond(c, started, http.StatusCreated, gin.H{"id": publicID("veh", vehicle.ID), "vehicle_number": vehicle.VehicleNumber, "status": "active"})
}

func (h *Handler) listVehicles(c *gin.Context) {
	started := time.Now()
	vehicles, err := h.service.ListVehicles()
	if err != nil {
		fail(c, started, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]ginMap, 0, len(vehicles))
	for _, vehicle := range vehicles {
		items = append(items, vehicleDTO(vehicle))
	}
	respond(c, started, http.StatusOK, gin.H{"vehicles": items})
}

func (h *Handler) updateLocation(c *gin.Context) {
	started := time.Now()
	var request struct {
		VehicleID string    `json:"vehicle_id" binding:"required"`
		Latitude  float64   `json:"latitude" binding:"required"`
		Longitude float64   `json:"longitude" binding:"required"`
		Timestamp time.Time `json:"timestamp" binding:"required"`
	}
	if !bindJSON(c, started, &request) {
		return
	}
	vehicleID, err := parsePublicID(request.VehicleID, "veh")
	if err != nil {
		fail(c, started, http.StatusBadRequest, err.Error())
		return
	}
	input := services.LocationInput{VehicleID: vehicleID, Latitude: request.Latitude, Longitude: request.Longitude, Timestamp: request.Timestamp}
	result, err := h.service.UpdateLocation(input)
	if err != nil {
		fail(c, started, http.StatusBadRequest, err.Error())
		return
	}
	current := make([]ginMap, 0, len(result.ActiveGeofences))
	for _, geofence := range result.ActiveGeofences {
		current = append(current, ginMap{"geofence_id": publicID("geo", geofence.ID), "geofence_name": geofence.Name, "status": "inside"})
	}
	respond(c, started, http.StatusCreated, gin.H{"vehicle_id": request.VehicleID, "location_updated": true, "current_geofences": current})
}

func (h *Handler) latestLocation(c *gin.Context) {
	started := time.Now()
	vehicleID, err := parsePublicID(c.Param("vehicle_id"), "veh")
	if err != nil {
		fail(c, started, http.StatusBadRequest, "invalid vehicle_id")
		return
	}
	location, err := h.service.LatestLocation(vehicleID)
	if err != nil {
		fail(c, started, http.StatusNotFound, err.Error())
		return
	}
	vehicles, err := h.service.ListVehicles()
	if err != nil {
		fail(c, started, http.StatusInternalServerError, err.Error())
		return
	}
	var vehicleNumber string
	for _, vehicle := range vehicles {
		if vehicle.ID == vehicleID {
			vehicleNumber = vehicle.VehicleNumber
			break
		}
	}
	geofences, err := h.service.ListGeofences("")
	if err != nil {
		fail(c, started, http.StatusInternalServerError, err.Error())
		return
	}
	current := currentGeofencesForLocation(location, geofences)
	respond(c, started, http.StatusOK, gin.H{"vehicle_id": c.Param("vehicle_id"), "vehicle_number": vehicleNumber, "current_location": locationDTO(location), "current_geofences": current})
}

func (h *Handler) createAlertRule(c *gin.Context) {
	started := time.Now()
	var request struct {
		GeofenceID string `json:"geofence_id" binding:"required"`
		VehicleID  string `json:"vehicle_id"`
		EventType  string `json:"event_type" binding:"required"`
	}
	if !bindJSON(c, started, &request) {
		return
	}
	geofenceID, err := parsePublicID(request.GeofenceID, "geo")
	if err != nil {
		fail(c, started, http.StatusBadRequest, err.Error())
		return
	}
	var vehicleID *uint
	if request.VehicleID != "" {
		parsed, err := parsePublicID(request.VehicleID, "veh")
		if err != nil {
			fail(c, started, http.StatusBadRequest, err.Error())
			return
		}
		vehicleID = &parsed
	}
	input := services.AlertRuleInput{VehicleID: vehicleID, GeofenceID: &geofenceID, EventType: request.EventType}
	rule, err := h.service.CreateAlertRule(input)
	if err != nil {
		fail(c, started, http.StatusBadRequest, err.Error())
		return
	}
	payload := gin.H{"alert_id": publicID("alert", rule.ID), "geofence_id": request.GeofenceID, "vehicle_id": nil, "event_type": rule.EventType, "status": "active"}
	if request.VehicleID != "" {
		payload["vehicle_id"] = request.VehicleID
	}
	respond(c, started, http.StatusCreated, payload)
}

func (h *Handler) listAlertRules(c *gin.Context) {
	started := time.Now()
	vehicleID, geofenceID, ok := parseOptionalAlertFilters(c, started)
	if !ok {
		return
	}
	rules, err := h.service.ListAlertRules(vehicleID, geofenceID)
	if err != nil {
		fail(c, started, http.StatusInternalServerError, err.Error())
		return
	}
	geofences, _ := h.service.ListGeofences("")
	vehicles, _ := h.service.ListVehicles()
	items := make([]ginMap, 0, len(rules))
	for _, rule := range rules {
		items = append(items, alertDTO(rule, findGeofence(rule.GeofenceID, geofences), findVehicle(rule.VehicleID, vehicles)))
	}
	respond(c, started, http.StatusOK, gin.H{"alerts": items})
}

func (h *Handler) listViolations(c *gin.Context) {
	started := time.Now()
	filter, err := parseViolationFilter(c)
	if err != nil {
		fail(c, started, http.StatusBadRequest, err.Error())
		return
	}
	violations, total, err := h.service.ListViolations(filter)
	if err != nil {
		fail(c, started, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]ginMap, 0, len(violations))
	for _, violation := range violations {
		items = append(items, violationDTO(violation))
	}
	respond(c, started, http.StatusOK, gin.H{"violations": items, "total_count": total})
}

func parseViolationFilter(c *gin.Context) (repositories.ViolationFilter, error) {
	filter := repositories.ViolationFilter{Limit: 50}
	if value := c.Query("limit"); value != "" {
		size, err := strconv.Atoi(value)
		if err != nil || size < 1 || size > 500 {
			return filter, strconv.ErrSyntax
		}
		filter.Limit = size
	}
	page := 1
	if value := c.Query("page"); value != "" {
		p, err := strconv.Atoi(value)
		if err == nil && p > 0 {
			page = p
		}
	}
	if filter.Limit > 0 {
		filter.Offset = (page - 1) * filter.Limit
	}
	if value := c.Query("vehicle_id"); value != "" {
		id, err := parsePublicID(value, "veh")
		if err != nil {
			return filter, err
		}
		filter.VehicleID = &id
	}
	if value := c.Query("geofence_id"); value != "" {
		id, err := parsePublicID(value, "geo")
		if err != nil {
			return filter, err
		}
		filter.GeofenceID = &id
	}
	if value := c.Query("start_date"); value != "" {
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return filter, err
		}
		filter.From = &t
	}
	if value := c.Query("end_date"); value != "" {
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return filter, err
		}
		filter.To = &t
	}
	return filter, nil
}

func parseOptionalAlertFilters(c *gin.Context, started time.Time) (*uint, *uint, bool) {
	var vehicleID *uint
	var geofenceID *uint
	if value := c.Query("vehicle_id"); value != "" {
		id, err := parsePublicID(value, "veh")
		if err != nil {
			fail(c, started, http.StatusBadRequest, err.Error())
			return nil, nil, false
		}
		vehicleID = &id
	}
	if value := c.Query("geofence_id"); value != "" {
		id, err := parsePublicID(value, "geo")
		if err != nil {
			fail(c, started, http.StatusBadRequest, err.Error())
			return nil, nil, false
		}
		geofenceID = &id
	}
	return vehicleID, geofenceID, true
}

func currentGeofencesForLocation(location models.VehicleLocation, geofences []models.Geofence) []ginMap {
	current := make([]ginMap, 0)
	point := models.Coordinate{Latitude: location.Latitude, Longitude: location.Longitude}
	for _, fence := range geofences {
		polygon, err := geofence.ParsePolygon(fence.Polygon)
		if err != nil {
			continue
		}
		if geofence.Contains(point, polygon) {
			current = append(current, geofenceSummaryDTO(fence))
		}
	}
	return current
}

func findGeofence(id *uint, geofences []models.Geofence) *models.Geofence {
	if id == nil {
		return nil
	}
	for i := range geofences {
		if geofences[i].ID == *id {
			return &geofences[i]
		}
	}
	return nil
}

func findVehicle(id *uint, vehicles []models.Vehicle) *models.Vehicle {
	if id == nil {
		return nil
	}
	for i := range vehicles {
		if vehicles[i].ID == *id {
			return &vehicles[i]
		}
	}
	return nil
}
