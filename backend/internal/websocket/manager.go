package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"geofencing-alerts/backend/internal/models"

	gws "github.com/gorilla/websocket"
)

type AlertEvent struct {
	EventID   string          `json:"event_id"`
	EventType string          `json:"event_type"`
	Timestamp time.Time       `json:"timestamp"`
	Vehicle   alertVehicle   `json:"vehicle"`
	Geofence  alertGeofence  `json:"geofence"`
	Location  alertLocation  `json:"location"`
}

type alertVehicle struct {
	VehicleID     string `json:"vehicle_id"`
	VehicleNumber string `json:"vehicle_number"`
	DriverName    string `json:"driver_name"`
}

type alertGeofence struct {
	GeofenceID   string `json:"geofence_id"`
	GeofenceName string `json:"geofence_name"`
	Category     string `json:"category"`
}

type alertLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewAlertEvent(eventID uint, eventType string, timestamp time.Time, vehicle models.Vehicle, geofence models.Geofence, location models.VehicleLocation) AlertEvent {
	return AlertEvent{
		EventID:   publicID("evt", eventID),
		EventType: eventType,
		Timestamp: timestamp,
		Vehicle: alertVehicle{
			VehicleID:     publicID("veh", vehicle.ID),
			VehicleNumber: vehicle.VehicleNumber,
			DriverName:    vehicle.DriverName,
		},
		Geofence: alertGeofence{
			GeofenceID:   publicID("geo", geofence.ID),
			GeofenceName: geofence.Name,
			Category:     geofence.Category,
		},
		Location: alertLocation{Latitude: location.Latitude, Longitude: location.Longitude},
	}
}

func publicID(prefix string, id uint) string {
	return fmt.Sprintf("%s_%d", prefix, id)
}

type Manager struct {
	clients    map[*gws.Conn]bool
	register   chan *gws.Conn
	unregister chan *gws.Conn
	events     chan AlertEvent
	mu         sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*gws.Conn]bool),
		register:   make(chan *gws.Conn),
		unregister: make(chan *gws.Conn),
		events:     make(chan AlertEvent, 256),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case conn := <-m.register:
			m.mu.Lock()
			m.clients[conn] = true
			m.mu.Unlock()
		case conn := <-m.unregister:
			m.remove(conn)
		case event := <-m.events:
			m.broadcast(event)
		}
	}
}

func (m *Manager) Publish(event AlertEvent) {
	select {
	case m.events <- event:
	default:
	}
}

func (m *Manager) Serve(w http.ResponseWriter, r *http.Request) {
	upgrader := gws.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	m.register <- conn

	go func() {
		defer func() { m.unregister <- conn }()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}

func (m *Manager) broadcast(event AlertEvent) {
	m.mu.RLock()
	conns := make([]*gws.Conn, 0, len(m.clients))
	for conn := range m.clients {
		conns = append(conns, conn)
	}
	m.mu.RUnlock()

	for _, conn := range conns {
		if err := conn.WriteJSON(event); err != nil {
			m.remove(conn)
		}
	}
}

func (m *Manager) remove(conn *gws.Conn) {
	m.mu.Lock()
	if m.clients[conn] {
		delete(m.clients, conn)
		_ = conn.Close()
	}
	m.mu.Unlock()
}
