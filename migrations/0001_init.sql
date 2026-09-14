-- 0001_init.sql
-- Core schema: drivers, orders, stock. Extend as the matching engine's needs become clearer.

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE drivers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'offline', -- offline | available | assigned
    location GEOGRAPHY(POINT, 4326),        -- lat/lng via PostGIS
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Spatial index for nearest-driver queries
CREATE INDEX idx_drivers_location ON drivers USING GIST (location);
CREATE INDEX idx_drivers_status ON drivers (status);

CREATE TABLE merchants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    location GEOGRAPHY(POINT, 4326)
);

CREATE TABLE stock_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES merchants(id),
    name TEXT NOT NULL,
    quantity INT NOT NULL CHECK (quantity >= 0)
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES merchants(id),
    driver_id UUID REFERENCES drivers(id),
    status TEXT NOT NULL DEFAULT 'pending', -- pending | matched | picked_up | delivered | cancelled
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_status ON orders (status);
