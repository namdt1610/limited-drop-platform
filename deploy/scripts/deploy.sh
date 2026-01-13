#!/bin/bash

# Production Deployment Script
# Sets up VPS with Docker (DB + NocoDB only) + Go binary backend

set -e

echo "Starting Production Deployment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running as root
if [[ $EUID -eq 0 ]]; then
   echo -e "${RED}This script should not be run as root${NC}"
   exit 1
fi

# Update system
echo -e "${YELLOW}Updating system...${NC}"
sudo apt update && sudo apt upgrade -y

# Install Docker and Docker Compose
echo -e "${YELLOW}Installing Docker...${NC}"
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

echo -e "${YELLOW}Installing Docker Compose...${NC}"
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create project directory
echo -e "${YELLOW}Setting up project directory...${NC}"
mkdir -p ~/donald
cd ~/donald

# Create production .env file
echo -e "${YELLOW}Creating production environment file...${NC}"
cat > .env << EOF
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=$(openssl rand -hex 16)
DB_NAME=donald_prod
DB_SSLMODE=require

# Server Configuration
PORT=3030

# Performance Tuning
GOGC=200

# CORS (update with your domain)
CORS_ORIGIN=https://donaldvibe.xyz

# Rate Limiting
RATE_LIMIT_DISABLED=false

# Cloudinary (add your credentials)
CLOUDINARY_CLOUD_NAME=your-cloud-name
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret
CLOUDINARY_UPLOAD_PRESET=production

# PayOS (add your credentials)
PAYOS_CLIENT_ID=your-client-id
PAYOS_API_KEY=your-api-key
PAYOS_CHECKSUM_KEY=your-checksum-key

# Email (add your credentials)
RESEND_API_KEY=your-resend-key

# NocoDB Configuration
NOCODB_PORT=8080
NOCODB_PUBLIC_URL=https://nocodb.donaldvibe.xyz
NOCODB_DISABLE_TELEMETRY=true
NOCODB_JWT_SECRET=$(openssl rand -hex 32)
NOCODB_ADMIN_EMAIL=admin@donaldvibe.xyz
NOCODB_ADMIN_PASSWORD=$(openssl rand -hex 8)
EOF

echo -e "${GREEN}Environment file created. Please edit .env with your actual credentials!${NC}"

# Download docker-compose.yml (you'll need to upload this)
echo -e "${YELLOW}Please upload your docker-compose.yml file to this directory${NC}"
echo -e "${YELLOW}Then run: docker-compose up -d${NC}"

# Set up systemd service for Go backend
echo -e "${YELLOW}Setting up systemd service for Go backend...${NC}"

sudo tee /etc/systemd/system/donald-backend.service > /dev/null <<EOF
[Unit]
Description=E-commerce Backend Service
After=network.target postgresql.service
Requires=postgresql.service

[Service]
Type=simple
User=$USER
WorkingDirectory=/home/$USER/donald
ExecStart=/home/$USER/donald/server-linux
Restart=always
RestartSec=5
EnvironmentFile=/home/$USER/donald/.env

# Security
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/home/$USER/donald

# Limits
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload

echo -e "${GREEN}Production deployment setup complete!${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Edit .env with your actual credentials"
echo "2. Upload docker-compose.yml"
echo "3. Run: docker-compose up -d"
echo "4. Upload server-linux binary"
echo "5. Run: sudo systemctl enable donald-backend && sudo systemctl start donald-backend"
echo ""
echo -e "${GREEN}Services will be available at:${NC}"
echo "- Database: localhost:5432"
echo "- NocoDB: donaldvibe.xyz:8080"
echo "- Backend API: donaldvibe.xyz:3030"
echo "- Frontend: Deploy separately (Vercel/Netlify)"
