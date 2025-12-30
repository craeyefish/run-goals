# Local Development Guide

## Recommended: Local Dev with Production Database

### 1. Setup `.env`
```bash
cp backend/.env.template backend/.env
```

Edit `backend/.env`:
```
DATABASE_HOST=summitseekers-db-do-user-30273773-0.j.db.ondigitalocean.com
DATABASE_PORT=25060
DATABASE_USER=doadmin
DATABASE_PASSWORD=
DATABASE_DBNAME=summitseekers
DATABASE_SSLMODE=require

STRAVA_CLIENT_ID=your_client_id
STRAVA_CLIENT_SECRET=your_client_secret
JWT_SECRET=same_as_production

SUMMIT_THRESHOLD_METERS=0.0007
DISABLE_SYNC_JOB=true
```

### 2. Run
```bash
docker compose -f docker-compose.prod-db.yaml up --build
```

### 3. Access
- Frontend: http://localhost:4200
- Backend: http://localhost:8080

---

## Alternative: Full Local Stack

```bash
docker compose up --build
```

Uses local PostgreSQL. Need to seed data or login via Strava.

---

## Key Config

| Variable | Purpose |
|----------|---------|
| `DISABLE_SYNC_JOB=true` | Prevents 24hr background sync |
| `DATABASE_SSLMODE` | `disable` (local) / `require` (DO) |

---

## Strava OAuth Setup

1. https://www.strava.com/settings/api
2. Add `http://localhost:4200` to callback domains
3. Callback URL: `http://localhost:4200/auth/callback`

---

## Commands

```bash
# Start with prod DB
docker compose -f docker-compose.prod-db.yaml up --build

# Start full local
docker compose up --build

# View logs
docker compose logs -f backend

# Rebuild specific service
docker compose build backend

# Stop
docker compose down
docker compose down -v  # Also removes volumes
```
