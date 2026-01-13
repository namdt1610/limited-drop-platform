#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

DB_FILE="backend/database.db"

echo -e "${YELLOW}=== ULTIMATE INTEGRITY TEST ===${NC}"

# 1. Reset Database
echo -e "\n1. Resetting Database..."
sqlite3 $DB_FILE "DELETE FROM orders; DELETE FROM symbicodes; UPDATE limited_drops SET sold = 0, total_stock = 5000 WHERE id = 1;"
echo -e "${GREEN}Database cleaned. Stock set to 5000.${NC}"

# 2. Run Stress Test
echo -e "\n2. Running Stress Test (10s @ 10,000 RPS)..."
echo "Expected Total Requests: ~100,000"
echo "Expected Sold: 5,000 (Sold Out)"
# We use a custom run to simulate massive overselling
K6_MODE=stress RATE=10000 DURATION=10s MAX_VUS=5000 PREALLOC_VUS=2000 k6 run tests/load/k6-runner.js > /dev/null
echo -e "${GREEN}Stress test completed.${NC}"

# 3. Verify Integrity
echo -e "\n3. Verifying Data Integrity..."

SOLD=$(sqlite3 $DB_FILE "SELECT sold FROM limited_drops WHERE id = 1;")
ORDERS=$(sqlite3 $DB_FILE "SELECT count(*) FROM orders WHERE items LIKE '%\"product_id\":1%';")

echo "----------------------------------------"
echo "Stock Limit:     5000"
echo "DB Sold Count:   $SOLD"
echo "Actual Orders:   $ORDERS"
echo "----------------------------------------"

if [ "$SOLD" -eq 5000 ] && [ "$ORDERS" -eq 5000 ]; then
    echo -e "${GREEN}SUCCESS: Data is perfectly consistent!${NC}"
    echo "System processed ~100k requests and sold exactly 5000 items without error."
else
    echo -e "${RED}FAILURE: Data mismatch detected!${NC}"
    exit 1
fi
