# Nginx Configuration Guide

## Overview

This nginx configuration provides:

- ✅ Reverse proxy for Alpine frontend
- ✅ Reverse proxy for Fiber backend
- ✅ SSL/TLS termination (HTTPS)
- ✅ Gzip compression
- ✅ Rate limiting (API: 100 r/s, General: 50 r/s)
- ✅ Static asset caching
- ✅ Security headers
- ✅ SPA hash routing support

## Server Ports

| Port | Purpose                   | Environment |
| ---- | ------------------------- | ----------- |
| 80   | HTTP (redirects to HTTPS) | Production  |
| 443  | HTTPS (SSL)               | Production  |
| 8080 | Development (no SSL)      | Development |

## Development Setup

### Using Docker

```bash
# Build nginx image
docker build -t donald-nginx .

# Run with local services
docker run -d \
  --name donald-nginx \
  -p 8080:8080 \
  -v $(pwd)/nginx.conf:/etc/nginx/nginx.conf:ro \
  donald-nginx
```

Visit: http://localhost:8080

### Using docker-compose

Add to your docker-compose.yml:

```yaml
nginx:
  build: ./nginx
  container_name: donald-nginx
  ports:
    - "8080:8080"
    - "80:80"
    - "443:443"
  volumes:
    - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
    - ./certs:/etc/letsencrypt:ro # For production SSL
  depends_on:
    - backend
    - alpine
  networks:
    - donald-network
```

## Production Setup with SSL

### 1. Generate SSL Certificates with Let's Encrypt

```bash
# Using certbot
sudo apt-get install certbot python3-certbot-nginx
sudo certbot certonly --standalone -d yourdomain.com -d www.yourdomain.com

# Certificates location:
# - /etc/letsencrypt/live/yourdomain.com/fullchain.pem
# - /etc/letsencrypt/live/yourdomain.com/privkey.pem
```

### 2. Update nginx.conf

Uncomment and update these lines in the production server block:

```nginx
ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
ssl_trusted_certificate /etc/letsencrypt/live/yourdomain.com/chain.pem;

ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers HIGH:!aNULL:!MD5;
ssl_prefer_server_ciphers on;
ssl_session_cache shared:SSL:10m;
ssl_session_timeout 10m;
ssl_stapling on;
ssl_stapling_verify on;
```

### 3. Auto-Renew SSL Certificates

```bash
# Create renewal script
sudo nano /etc/letsencrypt/renewal-hook/nginx-reload.sh
```

```bash
#!/bin/bash
systemctl reload nginx
```

```bash
chmod +x /etc/letsencrypt/renewal-hook/nginx-reload.sh

# Certbot will auto-run this on renewal
```

## Environment Variables Support

### Update FRONTEND_URL in Backend

Make sure your backend is running with the correct frontend URL:

```bash
# Development
export FRONTEND_URL=http://localhost:8080

# Production
export FRONTEND_URL=https://yourdomain.com
```

### Update CORS_ORIGINS in Backend

```bash
# Development
export CORS_ORIGINS=http://localhost:8080

# Production
export CORS_ORIGINS=https://yourdomain.com
```

## Rate Limiting

Current limits (per IP address):

- **API routes** (`/api/*`): 100 requests/second, burst of 200
- **General routes**: 50 requests/second

Adjust in nginx.conf:

```nginx
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=100r/s;
limit_req_zone $binary_remote_addr zone=general_limit:10m rate=50r/s;
```

## Cache Strategy

### Static Assets (30 days)

- Extensions: `.js`, `.css`, `.png`, `.jpg`, `.gif`, `.ico`, `.svg`, `.woff`, `.woff2`, `.ttf`, `.eot`
- Cache-Control: `public, immutable, max-age=31536000`

### API Endpoints (no cache)

- No caching for `/api/*` routes
- Cache-Control: `no-cache, no-store, must-revalidate`

### HTML/SPA (no cache)

- HTML files served fresh on each request
- Allows 404 fallback to index.html for hash routing

## Testing

### Check nginx configuration syntax

```bash
# Inside Docker
docker exec donald-nginx nginx -t

# Locally
sudo nginx -t
```

### View logs

```bash
# Inside Docker
docker logs donald-nginx
docker exec donald-nginx tail -f /var/log/nginx/access.log

# Locally
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### Test API rate limiting

```bash
# Should work (under limit)
for i in {1..50}; do curl -s http://localhost:8080/api/products | head -c 20; done

# Will be rate limited
for i in {1..300}; do curl -s http://localhost:8080/api/products; done
```

### Test SPA hash routing

```bash
curl http://localhost:8080/#drop
curl http://localhost:8080/#payment-success
curl http://localhost:8080/#verify
# All should return index.html (200)
```

## Monitoring

### Health Check Endpoint

Nginx automatically checks backend health via `/health` endpoint:

```bash
curl http://localhost:8080/health
```

If backend is down, nginx will retry other upstream servers (if configured).

## Troubleshooting

### 502 Bad Gateway

- Check if backend (port 3030) is running
- Check if alpine frontend (port 3000) is running
- View nginx error log: `docker logs donald-nginx`

### 404 on API endpoints

- Verify backend is serving `/api/` routes
- Check CORS configuration in backend
- Use `docker exec` to curl backend directly: `curl http://backend:3030/api/products`

### Frontend not loading

- Check if Alpine dev server is running
- Verify docker network connectivity
- Try accessing frontend directly: `http://localhost:3000`

### SSL certificate errors

- Check certificate paths in nginx.conf
- Verify certificates exist: `ls -la /etc/letsencrypt/live/yourdomain.com/`
- Certbot validation: `certbot certificates`

### Rate limiting too strict/loose

- Adjust `limit_req_zone` and `limit_req` directives
- Test with: `ab -n 1000 -c 100 http://localhost:8080/api/products`
