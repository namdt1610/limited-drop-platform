-- Safe migration script to update database schema to match current backend models
-- This preserves existing data while adding missing columns

-- ===== USERS TABLE =====
-- Add missing columns for User struct
ALTER TABLE users ADD COLUMN IF NOT EXISTS total_spent BIGINT DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS total_orders INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_purchase_at TIMESTAMP;
-- Change is_active from boolean to SMALLINT (uint8 in Go)
ALTER TABLE users ALTER COLUMN is_active TYPE SMALLINT USING CASE WHEN is_active THEN 1 ELSE 0 END;
ALTER TABLE users ALTER COLUMN is_active SET DEFAULT 1;
-- Note: SMALLINT can hold uint8 values (0-255)

-- ===== PRODUCTS TABLE =====
-- Remove complex columns no longer needed
ALTER TABLE products DROP COLUMN IF EXISTS slug CASCADE;
ALTER TABLE products DROP COLUMN IF EXISTS category_id CASCADE;
ALTER TABLE products DROP COLUMN IF EXISTS base_price CASCADE;
ALTER TABLE products DROP COLUMN IF EXISTS discount_percent CASCADE;
ALTER TABLE products DROP COLUMN IF EXISTS is_featured CASCADE;

-- Add simplified columns
ALTER TABLE products ADD COLUMN IF NOT EXISTS price BIGINT DEFAULT 0;
ALTER TABLE products ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS thumbnail TEXT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS images JSONB DEFAULT '[]'::jsonb;
ALTER TABLE products ADD COLUMN IF NOT EXISTS tags JSONB DEFAULT '[]'::jsonb;
ALTER TABLE products ADD COLUMN IF NOT EXISTS stock INTEGER DEFAULT 0;
ALTER TABLE products ADD COLUMN IF NOT EXISTS status SMALLINT DEFAULT 0;
-- Change is_active from boolean to SMALLINT
ALTER TABLE products ALTER COLUMN is_active TYPE SMALLINT USING CASE WHEN is_active THEN 1 ELSE 0 END;
ALTER TABLE products ALTER COLUMN is_active SET DEFAULT 1;
-- Note: SMALLINT can hold uint8 values (0-255)

-- ===== ORDERS TABLE =====
-- Simplify order structure
ALTER TABLE orders DROP COLUMN IF EXISTS customer_name CASCADE;
ALTER TABLE orders DROP COLUMN IF EXISTS customer_email CASCADE;
ALTER TABLE orders DROP COLUMN IF EXISTS customer_address CASCADE;
ALTER TABLE orders DROP COLUMN IF EXISTS order_number CASCADE; -- Remove order_number column, now using Base32(ID)

-- Add simplified columns
ALTER TABLE orders ADD COLUMN IF NOT EXISTS customer_phone TEXT;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS shipping_address JSONB;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS items JSONB DEFAULT '[]'::jsonb;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_method SMALLINT DEFAULT 0;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS status SMALLINT DEFAULT 0;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS payos_order_code BIGINT;

-- Drop old is_paid column after migrating data
ALTER TABLE orders DROP COLUMN IF EXISTS is_paid CASCADE;

-- ===== CREATE NEW TABLES =====

-- Limited Drops table
CREATE TABLE IF NOT EXISTS limited_drops (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    name TEXT NOT NULL,
    total_stock INTEGER NOT NULL DEFAULT 0 CHECK (total_stock >= 0),
    drop_size INTEGER NOT NULL DEFAULT 1 CHECK (drop_size > 0),
    sold INTEGER NOT NULL DEFAULT 0 CHECK (sold >= 0),
    is_active SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Symbicode table
CREATE TABLE IF NOT EXISTS symbicode (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT,
    product_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    activated_at TIMESTAMP,
    code UUID NOT NULL UNIQUE, -- UUID v7 binary
    secret_key TEXT NOT NULL,
    activated_ip TEXT,
    is_activated SMALLINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ===== CREATE INDEXES =====

-- Limited drops indexes
CREATE INDEX IF NOT EXISTS idx_limited_drops_product_id ON limited_drops(product_id);
CREATE INDEX IF NOT EXISTS idx_limited_drops_start_time ON limited_drops(start_time);
CREATE INDEX IF NOT EXISTS idx_limited_drops_end_time ON limited_drops(end_time);
CREATE INDEX IF NOT EXISTS idx_limited_drops_name ON limited_drops(name);
CREATE INDEX IF NOT EXISTS idx_limited_drops_is_active ON limited_drops(is_active);

-- Symbicode indexes
CREATE INDEX IF NOT EXISTS idx_symbicode_order_id ON symbicode(order_id);
CREATE INDEX IF NOT EXISTS idx_symbicode_product_id ON symbicode(product_id);
CREATE INDEX IF NOT EXISTS idx_symbicode_created_at ON symbicode(created_at);
CREATE INDEX IF NOT EXISTS idx_symbicode_activated_at ON symbicode(activated_at);
CREATE INDEX IF NOT EXISTS idx_symbicode_code ON symbicode(code);
CREATE INDEX IF NOT EXISTS idx_symbicode_activated_ip ON symbicode(activated_ip);
CREATE INDEX IF NOT EXISTS idx_symbicode_is_activated ON symbicode(is_activated);

-- Orders indexes
CREATE INDEX IF NOT EXISTS idx_orders_customer_phone ON orders(customer_phone);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_payment_method ON orders(payment_method);
CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_payos_order_code ON orders(payos_order_code);

-- ===== MIGRATE EXISTING DATA =====

-- Migrate product prices from base_price if available
UPDATE products SET price = COALESCE(base_price, 0) WHERE price = 0 AND base_price IS NOT NULL;

-- Migrate order status to include payment state
-- If order was paid (is_paid = true), set status to OrderPaid (2)
-- Otherwise keep existing status or set to OrderPending (0)
UPDATE orders SET status = CASE
    WHEN is_paid = true THEN 2  -- OrderPaid
    WHEN status IS NULL OR status = 0 THEN 0  -- OrderPending
    ELSE status
END WHERE status IS NOT NULL OR is_paid IS NOT NULL;

-- Set default payment method for orders
UPDATE orders SET payment_method = 0 WHERE payment_method IS NULL;

-- Set default status for orders that don't have one
UPDATE orders SET status = 0 WHERE status IS NULL;

-- ===== PERFORMANCE OPTIMIZATION =====
ANALYZE users;
ANALYZE products;
ANALYZE orders;
ANALYZE limited_drops;
ANALYZE symbicode;
