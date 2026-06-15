# Setup Guide

## Overview

This repository contains a full-stack geofencing and real-time alert application:

- `backend/` - Go REST API and WebSocket alert server
- `frontend/` - React + Vite dashboard using Leaflet maps
- `docker-compose.yml` - local development stack with PostgreSQL, backend, and frontend
- `render.yaml` - Render deployment configuration for backend service and PostgreSQL

The backend auto-migrates schemas on startup and exposes real-time alerts via WebSocket at `/ws/alerts`.

## Prerequisites

- macOS, Linux, or Windows
- Go 1.24+
- Node.js 22+
- Docker Desktop with Docker Compose
- Optional: `curl` for API testing

## Environment variables

### Backend

The backend reads database configuration from `DATABASE_URL`.

Example:

```bash
export DATABASE_URL='postgresql://deadheaven07:password@localhost:5432/mapup_db?sslmode=disable'
export PORT=8080
export CORS_ORIGINS='http://localhost:5173'
```

### Frontend

The frontend uses Vite environment variables:

```bash
export VITE_API_BASE_URL='http://localhost:8080'
export VITE_WS_URL='ws://localhost:8080/ws/alerts'
```

### Docker Compose

The Compose stack expects `POSTGRES_PASSWORD` in the shell before startup.

```bash
export POSTGRES_PASSWORD='password'
```

## Launch locally with Docker Compose

From the repository root:

```bash
docker compose up --build
```

Open the app:

- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080`
- WebSocket alerts: `ws://localhost:8080/ws/alerts`

To stop the stack:

```bash
docker compose down
```

## Backend local development

### Start backend with Docker

```bash
cd backend
docker build -t mapup-backend .
```

### Run backend locally

```bash
cd backend
export DATABASE_URL='postgresql://deadheaven07:password@localhost:5432/mapup_db?sslmode=disable'
export PORT=8080
export CORS_ORIGINS='http://localhost:5173'
go run cmd/server/main.go
```

The backend will:

- connect to PostgreSQL using `DATABASE_URL`
- auto-migrate schema tables
- expose REST endpoints on port `8080`
- serve WebSocket connections at `/ws/alerts`

## Frontend local development

### Install dependencies

```bash
cd frontend
npm install
```

### Run development server

```bash
npm run dev -- --host 0.0.0.0 --port 5173
```

### Production build

```bash
npm run build
```

### Preview production build

```bash
npm run preview -- --host 0.0.0.0 --port 4173
```

## Database and migrations

The backend uses GORM auto-migrations via `backend/internal/database/database.go` and creates these tables on startup:

- `geofences`
- `vehicles`
- `vehicle_locations`
- `alert_rules`
- `violations`
- `vehicle_geofence_states`

The local Compose stack uses PostgreSQL 16 and stores data in the named volume `postgres_data`.

## API endpoints

The primary backend API endpoints are:

- `POST /geofences`
- `GET /geofences`
- `POST /vehicles`
- `GET /vehicles`
- `POST /vehicles/location`
- `GET /vehicles/location/:vehicle_id`
- `POST /alerts/configure`
- `GET /alerts`
- `GET /violations/history`
- `GET /ws/alerts` (WebSocket)

### Example requests

#### Create geofence

```bash
curl -X POST http://localhost:8080/geofences \
  -H 'Content-Type: application/json' \
  -d '{
    "name":"Downtown Delivery Zone",
    "description":"Main delivery area for downtown customers",
    "coordinates":[[37.7749,-122.4194],[37.7849,-122.4194],[37.7849,-122.4094],[37.7749,-122.4094],[37.7749,-122.4194]],
    "category":"delivery_zone"
  }'
```

#### Register vehicle

```bash
curl -X POST http://localhost:8080/vehicles \
  -H 'Content-Type: application/json' \
  -d '{
    "vehicle_number":"KA-01-AB-1234",
    "driver_name":"John Doe",
    "vehicle_type":"truck",
    "phone":"+1234567890"
  }'
```

#### Send location update

```bash
curl -X POST http://localhost:8080/vehicles/location \
  -H 'Content-Type: application/json' \
  -d '{
    "vehicle_id":"veh_1",
    "latitude":37.7849,
    "longitude":-122.4194,
    "timestamp":"2025-01-15T10:35:00Z"
  }'
```

#### Configure alert rule

```bash
curl -X POST http://localhost:8080/alerts/configure \
  -H 'Content-Type: application/json' \
  -d '{
    "geofence_id":"geo_1",
    "vehicle_id":"veh_1",
    "event_type":"entry"
  }'
```

#### Retrieve violation history

```bash
curl 'http://localhost:8080/violations/history?vehicle_id=veh_1&limit=50'
```

## WebSocket alert stream

The frontend consumes a WebSocket stream at:

```text
ws://localhost:8080/ws/alerts
```

The backend broadcasts alert events for geofence entry/exit violations and supports multiple concurrent dashboard clients.

## Deployment notes

### Render backend

The `render.yaml` file configures a Docker-backed web service and PostgreSQL database using Render.

- Backend service root: `backend`
- Database: PostgreSQL 16
- Environment variable: `DATABASE_URL` is injected from the Render-managed database
- CORS origin is set to the deployed frontend URL

### Frontend deployment

The frontend is designed to deploy as a static site. Build output is served by Nginx in `frontend/Dockerfile`.

### Local production validation

To validate in local production mode:

```bash
cd frontend
npm run build
npm run preview -- --host 0.0.0.0 --port 4173
```

Then point the browser to:

```text
http://localhost:4173
```

## Troubleshooting

- If the backend cannot start, confirm `DATABASE_URL` is set and PostgreSQL is reachable.
- If the frontend cannot connect, ensure `VITE_API_BASE_URL` matches backend URL.
- For WebSocket issues, verify `VITE_WS_URL` uses `ws://` and points to `/ws/alerts`.
- Use `docker compose logs -f backend` for backend diagnostics.

## Notes

- The backend app is in `backend/cmd/server/main.go`.
- The frontend dashboard lives under `frontend/src` and uses `react-leaflet` for map rendering.
- API responses include `time_ns` to report request execution timing.
