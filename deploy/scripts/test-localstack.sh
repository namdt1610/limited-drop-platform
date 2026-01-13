#!/bin/bash
# =============================================================================
# LocalStack Local Testing Script
# Run this to test the application with LocalStack locally
# =============================================================================

set -e

echo "=== LocalStack Local Testing ==="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check dependencies
check_dependency() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}Error: $1 is not installed${NC}"
        exit 1
    fi
}

check_dependency docker
check_dependency curl

# Start LocalStack
echo -e "${YELLOW}Starting LocalStack...${NC}"
docker compose -f docker-compose.localstack.yml up -d localstack

# Wait for LocalStack to be healthy
echo -e "${YELLOW}Waiting for LocalStack to be healthy...${NC}"
timeout 60 bash -c '
    until curl -sf http://localhost:4566/_localstack/health > /dev/null 2>&1; do 
        echo "Waiting..."
        sleep 2
    done
'
echo -e "${GREEN}LocalStack is ready!${NC}"

# Initialize resources
echo -e "${YELLOW}Initializing LocalStack resources...${NC}"

# Install awslocal if not present
if ! command -v awslocal &> /dev/null; then
    echo "Installing awscli-local..."
    pip install awscli-local awscli --quiet
fi

# Create S3 bucket
echo "Creating S3 bucket..."
awslocal s3 mb s3://donald-vibe --endpoint-url=http://localhost:4566 2>/dev/null || echo "Bucket already exists"

# Create SQS queues
echo "Creating SQS queues..."
awslocal sqs create-queue --queue-name donald-orders --endpoint-url=http://localhost:4566 2>/dev/null || echo "Queue already exists"
awslocal sqs create-queue --queue-name donald-emails --endpoint-url=http://localhost:4566 2>/dev/null || echo "Queue already exists"
awslocal sqs create-queue --queue-name donald-notifications --endpoint-url=http://localhost:4566 2>/dev/null || echo "Queue already exists"

# Create Secrets
echo "Creating secrets..."
awslocal secretsmanager create-secret --name donald/jwt-secret --secret-string "test-jwt-secret" --endpoint-url=http://localhost:4566 2>/dev/null || echo "Secret already exists"
awslocal secretsmanager create-secret --name donald/db-password --secret-string "test-db-password" --endpoint-url=http://localhost:4566 2>/dev/null || echo "Secret already exists"

echo -e "${GREEN}LocalStack resources initialized!${NC}"

# Verify resources
echo ""
echo -e "${YELLOW}=== Verifying LocalStack Resources ===${NC}"

echo ""
echo "S3 Buckets:"
awslocal s3 ls --endpoint-url=http://localhost:4566

echo ""
echo "SQS Queues:"
awslocal sqs list-queues --endpoint-url=http://localhost:4566 2>/dev/null | grep -o 'donald-[a-z]*' || echo "No queues found"

echo ""
echo "Secrets:"
awslocal secretsmanager list-secrets --endpoint-url=http://localhost:4566 2>/dev/null | grep -o '"Name": "[^"]*"' || echo "No secrets found"

# Run backend tests
echo ""
echo -e "${YELLOW}=== Running Backend Tests ===${NC}"
cd backend

export AWS_ENDPOINT_URL=http://localhost:4566
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
export USE_S3=true
export USE_SQS=true
export USE_SECRETS_MANAGER=true
export S3_BUCKET=donald-vibe

echo "Running unit tests..."
go test ./tests/unit/... -v 2>&1 | tail -20

echo ""
echo "Running integration tests..."
go test ./tests/integration/... -v 2>&1 | tail -20

cd ..

# Run frontend tests
echo ""
echo -e "${YELLOW}=== Running Frontend Tests ===${NC}"
cd alpine
bun run test:run 2>&1 | tail -20
cd ..

# Test S3 upload
echo ""
echo -e "${YELLOW}=== Testing S3 Upload ===${NC}"
echo "test content $(date)" > /tmp/test-upload.txt
awslocal s3 cp /tmp/test-upload.txt s3://donald-vibe/test-upload.txt --endpoint-url=http://localhost:4566
awslocal s3 ls s3://donald-vibe/ --endpoint-url=http://localhost:4566

# Test SQS messaging
echo ""
echo -e "${YELLOW}=== Testing SQS Messaging ===${NC}"
awslocal sqs send-message \
    --queue-url http://localhost:4566/000000000000/donald-orders \
    --message-body '{"orderId": "test-'$(date +%s)'", "action": "test"}' \
    --endpoint-url=http://localhost:4566 | grep -o '"MessageId": "[^"]*"'

echo ""
echo -e "${GREEN}=== All LocalStack Tests Completed ===${NC}"

# Summary
echo ""
echo "=== Summary ==="
echo "LocalStack Dashboard: http://localhost:8055"
echo "LocalStack Endpoint:  http://localhost:4566"
echo ""
echo "To stop LocalStack:"
echo "  docker compose -f docker-compose.localstack.yml down"
echo ""
echo "To view logs:"
echo "  docker logs localstack-main -f"
