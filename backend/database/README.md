# Database Migration Guide

## Overview
The backend has been updated to use a simplified data model. Run the migration script to update your database schema.

## Migration Script

### Safe Migration (Recommended)
**File:** `update-schema-to-models.sql`

This script preserves existing data while updating the schema to match the current backend models.

```bash
# Using the helper script (recommended):
./run-migration.sh safe

# Or manually connect to PostgreSQL:
psql -h localhost -U postgres -d your_database_name -f update-schema-to-models.sql
```

## What Changes

### Users Table
- ✅ `total_spent` BIGINT (tracks total money spent)
- ✅ `total_orders` INTEGER (tracks number of orders)
- ✅ `last_purchase_at` TIMESTAMP (last purchase timestamp)
- ✅ `is_active` changed from BOOLEAN to SMALLINT (bitwise operations)

### Products Table
- ✅ Simplified: removed `category_id`, `base_price`, `discount_percent`, `is_featured`, `slug`
- ✅ Added: `price` BIGINT, `description` TEXT, `thumbnail` TEXT, `images` JSONB, `tags` JSONB, `stock` INTEGER, `status` SMALLINT
- ✅ `is_active` changed from BOOLEAN to SMALLINT

### Orders Table
- ✅ Simplified: removed `customer_name`, `customer_email`, `customer_address`
- ✅ Added: `customer_phone` TEXT, `shipping_address` JSONB, `items` JSONB, `payment_method` SMALLINT
- ✅ `status` SMALLINT (uses OrderStatus constants - includes payment state)
- ✅ Removed: `is_paid` column (payment state now included in status)

### New Tables
- ✅ `limited_drops` - for flash sales/drops
- ✅ `symbicode` - for anti-counterfeit system

## Verification

After running the migration, restart your backend application. GORM AutoMigrate will handle any remaining schema adjustments automatically.

## Backup First!
⚠️ **ALWAYS BACKUP YOUR DATABASE BEFORE RUNNING MIGRATIONS!**

```bash
# Using the helper script:
./run-migration.sh backup

# Or manually:
pg_dump -h localhost -U postgres your_database_name > backup_$(date +%Y%m%d_%H%M%S).sql
```
