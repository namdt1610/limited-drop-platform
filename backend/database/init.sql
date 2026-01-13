-- Initial database setup for Donald Watch Ecommerce
-- SQLite version - creates all necessary tables and indexes

-- ===== USERS TABLE =====
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT,
    phone TEXT UNIQUE,
    total_spent INTEGER DEFAULT 0,
    total_orders INTEGER DEFAULT 0,
    last_purchase_at DATETIME,
    is_active INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ===== PRODUCTS TABLE =====
CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    thumbnail TEXT,
    images TEXT DEFAULT '[]',
    tags TEXT DEFAULT '[]',
    price INTEGER DEFAULT 0,
    stock INTEGER DEFAULT 0,
    is_active INTEGER DEFAULT 1,
    status INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- ===== ORDERS TABLE =====
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    total_amount INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    customer_phone TEXT,
    shipping_address TEXT,
    items TEXT DEFAULT '[]',
    payment_method INTEGER DEFAULT 0,
    status INTEGER DEFAULT 0,
    payos_order_code INTEGER UNIQUE
);

-- ===== LIMITED DROPS TABLE =====
CREATE TABLE IF NOT EXISTS limited_drops (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    name TEXT NOT NULL,
    total_stock INTEGER NOT NULL DEFAULT 0 CHECK (total_stock >= 0),
    drop_size INTEGER NOT NULL DEFAULT 1 CHECK (drop_size > 0),
    sold INTEGER NOT NULL DEFAULT 0 CHECK (sold >= 0),
    is_active INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ===== SYMBICODE TABLE =====
CREATE TABLE IF NOT EXISTS symbicode (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER,
    product_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    activated_at DATETIME,
    code TEXT NOT NULL UNIQUE,
    secret_key TEXT NOT NULL,
    activated_ip TEXT,
    is_activated INTEGER NOT NULL DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ===== CREATE INDEXES =====

-- Users indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_total_spent ON users(total_spent);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- Products indexes
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_price ON products(price);
CREATE INDEX IF NOT EXISTS idx_products_stock ON products(stock);
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);

-- Orders indexes
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_orders_customer_phone ON orders(customer_phone);
CREATE INDEX IF NOT EXISTS idx_orders_payment_method ON orders(payment_method);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);

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

-- ===== PERFORMANCE OPTIMIZATION =====
ANALYZE users;
ANALYZE products;
ANALYZE orders;
ANALYZE limited_drops;
ANALYZE symbicode;
