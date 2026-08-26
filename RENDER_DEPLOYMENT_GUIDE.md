# Men Under Fire - Render.com Deployment Guide

This guide covers deploying the Men Under Fire application to Render.com using Docker, including the frontend, backend, PostgreSQL database, and Keycloak authentication server.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Prerequisites](#prerequisites)
3. [Cost Estimation](#cost-estimation)
4. [Step 1: Create Dockerfiles](#step-1-create-dockerfiles)
5. [Step 2: Set Up PostgreSQL Database](#step-2-set-up-postgresql-database)
6. [Step 3: Deploy Keycloak](#step-3-deploy-keycloak)
7. [Step 4: Deploy Backend](#step-4-deploy-backend)
8. [Step 5: Deploy Frontend](#step-5-deploy-frontend)
9. [Step 6: Configure Keycloak for Production](#step-6-configure-keycloak-for-production)
10. [Step 7: Connect Everything](#step-7-connect-everything)
11. [Troubleshooting](#troubleshooting)

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                          Render.com                                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────────┐   │
│  │   Frontend   │    │   Backend    │    │      Keycloak        │   │
│  │   (Static)   │───▶│   (Docker)   │◀───│      (Docker)        │   │
│  │              │    │   Port 7001  │    │      Port 8080       │   │
│  └──────────────┘    └──────┬───────┘    └──────────┬───────────┘   │
│                             │                       │               │
│                             ▼                       ▼               │
│                     ┌──────────────┐        ┌──────────────┐        │
│                     │  PostgreSQL  │        │  PostgreSQL  │        │
│                     │  (App DB)    │        │ (Keycloak DB)│        │
│                     └──────────────┘        └──────────────┘        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Prerequisites

- [Render.com account](https://render.com)
- GitHub/GitLab repository with your code
- Docker installed locally (for testing)
- Domain name (optional but recommended for production)

---

## Cost Estimation

### Render.com Pricing (as of 2025)

| Service | Plan | Monthly Cost | Notes |
|---------|------|--------------|-------|
| **Frontend** | Static Site (Free) | $0 | 100 GB bandwidth/month |
| **Backend** | Starter | $7 | 512 MB RAM, 0.5 CPU |
| **Backend** | Standard | $25 | 2 GB RAM, 1 CPU (recommended) |
| **Keycloak** | Starter | $7 | 512 MB RAM (minimum viable) |
| **Keycloak** | Standard | $25 | 2 GB RAM (recommended) |
| **PostgreSQL (App)** | Starter | $7 | 256 MB RAM, 1 GB storage |
| **PostgreSQL (App)** | Standard | $20 | 1 GB RAM, 16 GB storage |
| **PostgreSQL (Keycloak)** | Starter | $7 | 256 MB RAM, 1 GB storage |

### Estimated Monthly Costs

| Configuration | Monthly Cost | Description |
|---------------|--------------|-------------|
| **Minimum (Dev/Testing)** | ~$28 | Frontend (free) + Backend (Starter $7) + Keycloak (Starter $7) + 2x PostgreSQL (Starter $14) |
| **Recommended (Production)** | ~$77 | Frontend (free) + Backend (Standard $25) + Keycloak (Standard $25) + 2x PostgreSQL (Starter $14) + SSL |
| **High Performance** | ~$130+ | Standard/Pro instances for all services |

### Free Tier Limitations

- Free static sites: 100 GB/month bandwidth
- Free web services: Spin down after 15 min inactivity (not suitable for production)
- Free databases: 90-day expiration

---

## Step 1: Create Dockerfiles

### Backend Dockerfile

Create `backend/Dockerfile`:

```dockerfile
# Build stage
FROM gradle:8.5-jdk21 AS build
WORKDIR /app

# Copy gradle files first for caching
COPY build.gradle.kts settings.gradle.kts ./
COPY gradle ./gradle

# Download dependencies
RUN gradle dependencies --no-daemon

# Copy source code
COPY src ./src

# Build the application
RUN gradle buildFatJar --no-daemon

# Runtime stage
FROM eclipse-temurin:21-jre-alpine
WORKDIR /app

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the built jar
COPY --from=build /app/build/libs/*-all.jar app.jar

# Change ownership
RUN chown -R appuser:appgroup /app
USER appuser

# Expose port
EXPOSE 7001

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:7001/health || exit 1

# Run the application
CMD ["java", "-jar", "app.jar"]
```

### Frontend Dockerfile

Create `frontend/Dockerfile`:

```dockerfile
# Build stage
FROM node:20-alpine AS build
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build arguments for environment variables
ARG VITE_API_URL
ARG VITE_KEYCLOAK_URL
ARG VITE_KEYCLOAK_REALM
ARG VITE_KEYCLOAK_CLIENT_ID

# Set environment variables for build
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_KEYCLOAK_URL=$VITE_KEYCLOAK_URL
ENV VITE_KEYCLOAK_REALM=$VITE_KEYCLOAK_REALM
ENV VITE_KEYCLOAK_CLIENT_ID=$VITE_KEYCLOAK_CLIENT_ID

# Build the application
RUN npm run build

# Production stage - using nginx
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### Frontend nginx.conf

Create `frontend/nginx.conf`:

```nginx
server {
    listen 80;
    server_name _;
    root /usr/share/nginx/html;
    index index.html;

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;

    # Handle SPA routing
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Cache static assets
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
}
```

### Keycloak Dockerfile (Custom)

Create `keycloak/Dockerfile`:

```dockerfile
FROM quay.io/keycloak/keycloak:24.0 AS builder

# Enable health and metrics support
ENV KC_HEALTH_ENABLED=true
ENV KC_METRICS_ENABLED=true

# Configure database vendor
ENV KC_DB=postgres

WORKDIR /opt/keycloak

# Build optimized Keycloak
RUN /opt/keycloak/bin/kc.sh build

FROM quay.io/keycloak/keycloak:24.0
COPY --from=builder /opt/keycloak/ /opt/keycloak/

# Environment variables (to be overridden at runtime)
ENV KC_DB=postgres
ENV KC_HEALTH_ENABLED=true
ENV KC_METRICS_ENABLED=true

ENTRYPOINT ["/opt/keycloak/bin/kc.sh"]
```

---

## Step 2: Set Up PostgreSQL Database

### 2.1 Create Application Database

1. Go to [Render Dashboard](https://dashboard.render.com)
2. Click **New** → **PostgreSQL**
3. Configure:
   - **Name**: `menunderfire-db`
   - **Database**: `menunderfire`
   - **User**: `menunderfire_user`
   - **Region**: Choose closest to your users
   - **Plan**: Starter ($7/month) or Standard ($20/month)
4. Click **Create Database**
5. Wait for provisioning and note the connection details:
   - **Internal Database URL** (for services in same region)
   - **External Database URL** (for local development)

### 2.2 Create Keycloak Database

1. Click **New** → **PostgreSQL**
2. Configure:
   - **Name**: `keycloak-db`
   - **Database**: `keycloak`
   - **User**: `keycloak_user`
   - **Region**: Same as application database
   - **Plan**: Starter ($7/month)
3. Click **Create Database**
4. Note the connection details

---

## Step 3: Deploy Keycloak

### 3.1 Create Keycloak Web Service

1. Click **New** → **Web Service**
2. Connect your repository or use **Deploy from a Docker Registry**
3. Configure:
   - **Name**: `menunderfire-keycloak`
   - **Region**: Same as databases
   - **Instance Type**: Starter ($7) or Standard ($25)
   - **Docker Command**: `start --optimized`

### 3.2 Environment Variables for Keycloak

Add these environment variables in Render:

| Variable | Value | Notes |
|----------|-------|-------|
| `KC_DB` | `postgres` | Database type |
| `KC_DB_URL` | `jdbc:postgresql://[INTERNAL_HOST]:5432/keycloak` | Use internal hostname from Keycloak DB |
| `KC_DB_USERNAME` | `keycloak_user` | From database creation |
| `KC_DB_PASSWORD` | `[PASSWORD]` | From database creation |
| `KC_HOSTNAME` | `menunderfire-keycloak.onrender.com` | Your Render hostname |
| `KC_HOSTNAME_STRICT` | `false` | Allow reverse proxy |
| `KC_PROXY` | `edge` | Render handles TLS termination |
| `KEYCLOAK_ADMIN` | `admin` | Initial admin username |
| `KEYCLOAK_ADMIN_PASSWORD` | `[STRONG_PASSWORD]` | Generate a strong password |
| `KC_HTTP_ENABLED` | `true` | Allow HTTP (Render handles HTTPS) |
| `KC_HEALTH_ENABLED` | `true` | Enable health checks |

### 3.3 Alternative: Use render.yaml Blueprint

Create `keycloak/render.yaml`:

```yaml
services:
  - type: web
    name: menunderfire-keycloak
    env: docker
    dockerfilePath: ./keycloak/Dockerfile
    dockerCommand: start --optimized
    healthCheckPath: /health/ready
    envVars:
      - key: KC_DB
        value: postgres
      - key: KC_DB_URL
        fromDatabase:
          name: keycloak-db
          property: connectionString
      - key: KC_DB_USERNAME
        fromDatabase:
          name: keycloak-db
          property: user
      - key: KC_DB_PASSWORD
        fromDatabase:
          name: keycloak-db
          property: password
      - key: KC_HOSTNAME
        value: menunderfire-keycloak.onrender.com
      - key: KC_PROXY
        value: edge
      - key: KEYCLOAK_ADMIN
        sync: false
      - key: KEYCLOAK_ADMIN_PASSWORD
        generateValue: true
```

---

## Step 4: Deploy Backend

### 4.1 Create Backend Web Service

1. Click **New** → **Web Service**
2. Connect your GitHub/GitLab repository
3. Configure:
   - **Name**: `menunderfire-api`
   - **Root Directory**: `backend`
   - **Environment**: Docker
   - **Dockerfile Path**: `./Dockerfile`
   - **Region**: Same as databases
   - **Instance Type**: Starter ($7) or Standard ($25)

### 4.2 Environment Variables for Backend

| Variable | Value | Notes |
|----------|-------|-------|
| `PORT` | `7001` | Application port |
| `DATABASE_URL` | `jdbc:postgresql://[INTERNAL_HOST]:5432/menunderfire` | Use internal hostname |
| `DATABASE_USER` | `menunderfire_user` | From database creation |
| `DATABASE_PASSWORD` | `[PASSWORD]` | From database creation |
| `AUTH_JWKS_URL` | `https://menunderfire-keycloak.onrender.com/realms/menunderfire/protocol/openid-connect/certs` | Keycloak JWKS endpoint |
| `AUTH0_DOMAIN` | `menunderfire-keycloak.onrender.com/realms/menunderfire` | Keycloak realm URL |
| `AUTH0_ISSUER` | `https://menunderfire-keycloak.onrender.com/realms/menunderfire` | Token issuer |
| `AUTH0_AUDIENCE` | `menunderfire-api` | API audience (configure in Keycloak) |

### 4.3 Add Health Endpoint to Backend

The backend needs a `/health` endpoint for Render's health checks. Add this to `Routing.kt`:

```kotlin
// In configureRouting function, add before the /api route block:
routing {
    // Health check endpoint (no auth required)
    get("/health") {
        call.respond(HttpStatusCode.OK, mapOf("status" to "healthy"))
    }

    route("/api") {
        // ... existing routes
    }
}
```

You'll also need to add the import:
```kotlin
import io.ktor.server.response.*
```

### 4.4 Health Check Configuration

Set in Render Dashboard:
- **Health Check Path**: `/health`
- **Health Check Timeout**: 30 seconds

---

## Step 5: Deploy Frontend

### Option A: Static Site (Recommended - Free)

1. Click **New** → **Static Site**
2. Connect your repository
3. Configure:
   - **Name**: `menunderfire-web`
   - **Root Directory**: `frontend`
   - **Build Command**: `npm ci && npm run build`
   - **Publish Directory**: `dist`

4. Add Environment Variables:
   | Variable | Value |
   |----------|-------|
   | `VITE_API_URL` | `https://menunderfire-api.onrender.com` |
   | `VITE_KEYCLOAK_URL` | `https://menunderfire-keycloak.onrender.com` |
   | `VITE_KEYCLOAK_REALM` | `menunderfire` |
   | `VITE_KEYCLOAK_CLIENT_ID` | `menunderfire-frontend` |

5. Add Redirect/Rewrite Rules for SPA:
   - In **Redirects/Rewrites**, add:
     - Source: `/*`
     - Destination: `/index.html`
     - Action: Rewrite

### Option B: Docker Web Service

1. Click **New** → **Web Service**
2. Configure:
   - **Name**: `menunderfire-web`
   - **Root Directory**: `frontend`
   - **Environment**: Docker
   - **Instance Type**: Starter ($7)

3. Add the same environment variables as build args

---

## Step 6: Configure Keycloak for Production

After Keycloak is deployed and running:

### 6.1 Access Admin Console

1. Navigate to `https://menunderfire-keycloak.onrender.com`
2. Login with admin credentials

### 6.2 Create Realm

1. Click dropdown (top-left showing "master")
2. Click **Create realm**
3. **Realm name**: `menunderfire`
4. Click **Create**

### 6.3 Create Frontend Client

1. Go to **Clients** → **Create client**
2. **Client type**: OpenID Connect
3. **Client ID**: `menunderfire-frontend`
4. Click **Next**
5. **Client authentication**: OFF (public SPA)
6. Click **Next**
7. **Login settings**:
   - **Root URL**: `https://menunderfire-web.onrender.com`
   - **Valid redirect URIs**: `https://menunderfire-web.onrender.com/*`
   - **Valid post logout redirect URIs**: `https://menunderfire-web.onrender.com/*`
   - **Web origins**: `https://menunderfire-web.onrender.com`
8. Click **Save**

### 6.4 Create API Client (Optional - for backend validation)

1. Go to **Clients** → **Create client**
2. **Client ID**: `menunderfire-api`
3. **Client authentication**: ON (confidential)
4. This creates an audience that can be used for API validation

### 6.5 Create Roles

1. Go to **Realm roles** → **Create role**
2. Create roles: `user`, `org_admin`, `site_admin`

### 6.6 Create Test User

1. Go to **Users** → **Add user**
2. Fill in details and click **Create**
3. Go to **Credentials** tab → **Set password**
4. Go to **Role mapping** → **Assign role**

---

## Step 7: Connect Everything

### 7.1 Complete render.yaml Blueprint

Create `render.yaml` in your repository root:

```yaml
services:
  # Backend API
  - type: web
    name: menunderfire-api
    env: docker
    rootDir: backend
    dockerfilePath: ./Dockerfile
    healthCheckPath: /health
    envVars:
      - key: PORT
        value: 7001
      - key: DATABASE_URL
        fromDatabase:
          name: menunderfire-db
          property: connectionString
      - key: DATABASE_USER
        fromDatabase:
          name: menunderfire-db
          property: user
      - key: DATABASE_PASSWORD
        fromDatabase:
          name: menunderfire-db
          property: password
      - key: AUTH_JWKS_URL
        sync: false
      - key: AUTH0_DOMAIN
        sync: false
      - key: AUTH0_ISSUER
        sync: false

  # Keycloak
  - type: web
    name: menunderfire-keycloak
    env: docker
    rootDir: keycloak
    dockerfilePath: ./Dockerfile
    dockerCommand: start --optimized
    healthCheckPath: /health/ready
    envVars:
      - key: KC_DB
        value: postgres
      - key: KC_DB_URL
        fromDatabase:
          name: keycloak-db
          property: connectionString
      - key: KC_DB_USERNAME
        fromDatabase:
          name: keycloak-db
          property: user
      - key: KC_DB_PASSWORD
        fromDatabase:
          name: keycloak-db
          property: password
      - key: KC_HOSTNAME
        sync: false
      - key: KC_PROXY
        value: edge
      - key: KC_HTTP_ENABLED
        value: "true"
      - key: KEYCLOAK_ADMIN
        sync: false
      - key: KEYCLOAK_ADMIN_PASSWORD
        generateValue: true

  # Frontend (Static)
  - type: web
    name: menunderfire-web
    env: static
    rootDir: frontend
    buildCommand: npm ci && npm run build
    staticPublishPath: dist
    headers:
      - path: /*
        name: X-Frame-Options
        value: SAMEORIGIN
    routes:
      - type: rewrite
        source: /*
        destination: /index.html
    envVars:
      - key: VITE_API_URL
        sync: false
      - key: VITE_KEYCLOAK_URL
        sync: false
      - key: VITE_KEYCLOAK_REALM
        value: menunderfire
      - key: VITE_KEYCLOAK_CLIENT_ID
        value: menunderfire-frontend

databases:
  - name: menunderfire-db
    databaseName: menunderfire
    user: menunderfire_user

  - name: keycloak-db
    databaseName: keycloak
    user: keycloak_user
```

### 7.2 Deploy Using Blueprint

1. Push `render.yaml` to your repository
2. In Render Dashboard, click **Blueprints**
3. Click **New Blueprint Instance**
4. Select your repository
5. Review and configure sync settings
6. Click **Apply**

---

## Troubleshooting

### Common Issues

#### Keycloak fails to start

- Check database connectivity
- Ensure `KC_PROXY=edge` is set
- Verify `KC_HOSTNAME` matches your Render URL
- Check logs for database migration errors

#### Backend can't connect to database

- Use **internal** database URL, not external
- Verify the JDBC URL format: `jdbc:postgresql://host:5432/database`
- Check if database is in the same region

#### Frontend auth issues

- Verify Keycloak client redirect URIs match exactly
- Check CORS settings in Keycloak client
- Ensure all VITE_* variables are set during build

#### CORS errors

- Add frontend URL to Keycloak client's Web Origins
- Verify backend CORS configuration includes frontend URL

### Checking Logs

1. Go to your service in Render Dashboard
2. Click **Logs** tab
3. Use filters to search for errors

### Health Check Failures

If health checks fail:
1. Verify the health endpoint exists and works locally
2. Increase health check timeout (30-60 seconds for Keycloak)
3. Check if the service needs more memory (upgrade instance)

---

## Security Checklist

Before going to production:

- [ ] Change default Keycloak admin password
- [ ] Use strong database passwords (auto-generated by Render)
- [ ] Enable HTTPS only (Render does this by default)
- [ ] Configure proper CORS origins
- [ ] Set up database backups in Render
- [ ] Configure rate limiting in backend
- [ ] Remove or secure any debug endpoints
- [ ] Set proper Content Security Policy headers
- [ ] Configure Keycloak brute force detection
- [ ] Set up monitoring/alerting in Render

---

## Custom Domain Setup

1. Go to your service settings in Render
2. Click **Custom Domains**
3. Add your domain (e.g., `app.menunderfire.com`)
4. Add the provided DNS records to your domain registrar
5. Wait for SSL certificate provisioning
6. Update Keycloak `KC_HOSTNAME` to your custom domain
7. Update all redirect URIs in Keycloak clients

---

## Maintenance

### Database Backups

Render automatically backs up PostgreSQL databases:
- Starter: Daily backups, 7-day retention
- Standard: Daily backups, 14-day retention

### Updating Services

1. Push changes to your repository
2. Render auto-deploys on push (if auto-deploy enabled)
3. Or manually trigger deploy from Dashboard

### Scaling

To handle more traffic:
1. Upgrade to higher instance types
2. Enable auto-scaling (Pro plan)
3. Consider adding a CDN for frontend assets
