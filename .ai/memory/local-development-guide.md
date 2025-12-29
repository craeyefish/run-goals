# Local Development Guide

## Quick Start: Local Dev with Production Database

This is the recommended approach for most development work.

### 1. Set up your `.env` file

```bash
cp backend/.env.template backend/.env
# Edit backend/.env with your production DB credentials and Strava API keys
```

Your `backend/.env` should contain:
```
DATABASE_HOST=your-db-cluster.db.ondigitalocean.com
DATABASE_PORT=25060
DATABASE_USER=doadmin
DATABASE_PASSWORD=your_password
DATABASE_DBNAME=run_goals
DATABASE_SSLMODE=require

STRAVA_CLIENT_ID=your_client_id
STRAVA_CLIENT_SECRET=your_client_secret
JWT_SECRET=same_as_production

SUMMIT_THRESHOLD_METERS=0.0007
DISTANCE_CACHE_TTL=1
DISABLE_SYNC_JOB=true
```

### 2. Run with production DB

```bash
docker compose -f docker-compose.prod-db.yaml up --build
```

This runs only frontend + backend (no local database container).

### 3. Access

- Frontend: http://localhost:4200
- Backend API: http://localhost:8080

---

## Alternative: Full Local Stack (with local database)

Use the standard docker-compose for a completely isolated local environment:

```bash
docker compose up --build
```

This starts all 3 services (db, backend, frontend) with a local PostgreSQL.

**Note:** You'll need to seed the local database with your data or login via Strava OAuth to create your user.

---

## Key Configuration

### DISABLE_SYNC_JOB

Set `DISABLE_SYNC_JOB=true` in your `.env` to prevent the background job that syncs all users' activities every 24 hours. This avoids unnecessary Strava API calls during development.

### DATABASE_SSLMODE

- Local DB: `disable`
- Production DO DB: `require`

---

## Strava OAuth for Local Development

For Strava OAuth to work locally, you need to configure your Strava API app:

1. Go to https://www.strava.com/settings/api
2. Add `http://localhost:4200` to authorized callback domains
3. Your callback URL should be `http://localhost:4200/auth/callback`

---

## Quick Commands

```bash
# Run with production DB
docker compose -f docker-compose.prod-db.yaml up --build

# Run full local stack
docker compose up --build

# Stop and clean up
docker compose down
docker compose down -v  # Also removes volumes (resets local DB)

# View logs
docker compose logs -f backend
docker compose logs -f frontend

# Rebuild a specific service
docker compose build backend
```
