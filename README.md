# 🚗 Assignment MapUp

A real-time vehicle tracking and geofencing system engineered for high-throughput location updates, automatic boundary violation detection, and live alert streaming.

---

## 🚀 Live Deployments

The application is fully containerized and deployed across a distributed cloud infrastructure.

### Frontend Dashboard

https://assignment-map-up.vercel.app/

### Production API Backend

https://mapup-backend-hou5.onrender.com

### Database Engine

Managed PostgreSQL 16 Cluster hosted on Render.

---

## ✨ Features

### Real-Time Telemetry Streaming

Live vehicle position tracking using persistent WebSocket connections for low-latency updates.

### Geofencing & Polygon Evaluation

Automatic backend validation to determine whether incoming vehicle coordinates are inside or outside configured geofence boundaries.

### Instant Violation Alerting

Real-time push notifications broadcast to all connected dashboards whenever a vehicle breaches a geofence.

### Comprehensive Vehicle Management

Full CRUD support for registering, retrieving, and managing tracked vehicles.

### Production-Ready Deployment

Containerized architecture with Docker, Docker Compose, Render Blueprints, and Vercel deployment pipelines.

---

## 🛠️ Tech Stack

### Backend

* Go (Golang)
* Gin Gonic
* GORM
* PostgreSQL 16
* Gorilla WebSocket

### Frontend

* React 18
* TypeScript
* Vite
* Axios
* Native WebSocket APIs

### Infrastructure & DevOps

* Docker
* Docker Compose
* Render Blueprints (`render.yaml`)
* Vercel

---

## 🏗️ System Architecture

```text
                    ┌──────────────────┐
                    │ React Dashboard  │
                    │  (Vercel Hosted) │
                    └─────────┬────────┘
                              │
                     REST APIs│
                     WebSocket│
                              ▼
                    ┌──────────────────┐
                    │   Go Backend     │
                    │ Gin + GORM + WS  │
                    └─────────┬────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │ PostgreSQL 16    │
                    │  (Render Hosted) │
                    └──────────────────┘

Incoming Vehicle Coordinates
            │
            ▼
      Geofence Engine
            │
            ▼
    Violation Detection
            │
            ▼
 Real-Time Alert Broadcast
```

---

## 📁 Project Structure

```text
Assignment_MapUp_/
├── backend/
│   ├── cmd/
│   ├── config/
│   ├── controllers/
│   ├── models/
│   ├── repository/
│   ├── services/
│   ├── Dockerfile
│   └── go.mod
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── pages/
│   │   └── App.tsx
│   ├── Dockerfile
│   └── package.json
│
├── docker-compose.yml
└── render.yaml
```

---

## 💻 Local Development Setup

### Prerequisites

* Go 1.21+
* Node.js 18+
* Docker Desktop

---

### 1. Clone Repository

```bash
git clone https://github.com/deadheaven07/Assignment_MapUp_.git
cd Assignment_MapUp_
```

---

### 2. Run Using Docker Compose

Launch Backend, Frontend, and PostgreSQL together:

```bash
docker compose up --build
```

Services:

| Service    | URL                   |
| ---------- | --------------------- |
| Frontend   | http://localhost:5173 |
| Backend    | http://localhost:8080 |
| PostgreSQL | localhost:5432        |

---

### 3. Manual Backend Setup

```bash
cd backend

cp .env.example .env

# Configure environment variables

go run cmd/main.go
```

---

### 4. Manual Frontend Setup

```bash
cd frontend

cp .env.example .env

npm install

npm run dev
```

---

## 🔐 Environment Variables

### Backend

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=mapup
DB_SSLMODE=disable
```

### Frontend

```env
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080/ws/alerts
```

---

## 📡 API Reference

### Vehicle Management

| Method | Endpoint          | Description                |
| ------ | ----------------- | -------------------------- |
| GET    | /api/vehicles     | Retrieve all vehicles      |
| POST   | /api/vehicles     | Register a new vehicle     |
| GET    | /api/vehicles/:id | Retrieve vehicle details   |
| PUT    | /api/vehicles/:id | Update vehicle information |
| DELETE | /api/vehicles/:id | Remove a vehicle           |

---

### Geofence Management

| Method | Endpoint           | Description               |
| ------ | ------------------ | ------------------------- |
| GET    | /api/geofences     | Retrieve all geofences    |
| POST   | /api/geofences     | Create a geofence         |
| GET    | /api/geofences/:id | Retrieve geofence details |
| DELETE | /api/geofences/:id | Remove a geofence         |

---

### Vehicle Location Tracking

| Method | Endpoint                   | Description                      |
| ------ | -------------------------- | -------------------------------- |
| POST   | /api/vehicles/location     | Submit vehicle coordinates       |
| GET    | /api/vehicles/location/:id | Retrieve latest vehicle location |

---

## 🔄 WebSocket Streaming

### Alert Stream

```text
ws://localhost:8080/ws/alerts
```

The WebSocket server broadcasts:

* Vehicle location updates
* Geofence breach events
* Boundary entry notifications
* Boundary exit notifications

### Sample Alert Payload

```json
{
  "vehicle_id": 1,
  "geofence_id": 3,
  "event": "outside_geofence",
  "timestamp": "2026-06-15T12:30:00Z"
}
```

---

## 🐳 Docker Deployment

Build and run locally:

```bash
docker compose up --build
```

Stop services:

```bash
docker compose down
```

Rebuild after changes:

```bash
docker compose up --build --force-recreate
```

---

## 🎯 Assignment Requirements Covered

* Vehicle CRUD operations
* Geofence creation and management
* Polygon containment checks
* Real-time WebSocket updates
* PostgreSQL persistence
* Dockerized deployment
* Production cloud deployment
* Responsive frontend dashboard
* Separation of concerns architecture

---

## 🎥 Demo Videos

### Application Demo

https://www.loom.com/share/3dfc05b0cdaf42309731613d42d2b1d6

### Technical Walkthrough

https://www.loom.com/share/673c183c7a4a43a8ad3bc1eee9ce7d19


---

## 🔮 Future Improvements

* Authentication & Authorization
* JWT-based secure APIs
* Historical route playback
* Multi-geofence support
* Map clustering for large fleets
* Redis pub/sub for horizontal scaling
* Kafka-based telemetry ingestion
* Role-based access control
* Advanced analytics dashboard

---

## 👨‍💻 Author

**Harsh Raghuwanshi**

Built as part of the MapUp Full-Stack Developer Assessment.
