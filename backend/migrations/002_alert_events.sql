CREATE TABLE IF NOT EXISTS alert_events (
    id BIGSERIAL PRIMARY KEY,
    alert_rule_id BIGINT REFERENCES alert_rules(id),
    violation_id BIGINT REFERENCES violations(id),
    vehicle_id BIGINT NOT NULL REFERENCES vehicles(id),
    geofence_id BIGINT NOT NULL REFERENCES geofences(id),
    event_type TEXT NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ
);