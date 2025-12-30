# Summit Seekers (Run Goals) - Project Overview

## Quick Summary

A fitness goal tracking application integrated with Strava that enables users to:
- Track running/hiking activities synced from Strava
- Log summit achievements by detecting peaks crossed during activities
- Join groups and contribute to shared distance/elevation/summit goals
- Visualize activities and summits on an interactive map
- **Personal summit wishlist** - users can select peaks they want to summit

**Live URL**: summitseekers.co.za

---

## Recent Session Context (December 2024)

### What Was Done Recently

1. **Peak Data Enhancement** - Extended the `peaks` table with additional OSM metadata:
   - New fields: `alt_name`, `name_en`, `region`, `wikipedia`, `wikidata`, `description`, `prominence`
   - Admin endpoint `POST /admin/refresh-peaks?admin_key=dev-admin-key` to refresh peak data
   - Now fetches ALL peaks from Western Cape (including unnamed) - **1,629 peaks total**

2. **Peak Picker Component** - New reusable component for selecting peaks:
   - Map-based selection with marker clusters
   - Search by name, region, or alt_name
   - Shows region/elevation for peak differentiation
   - Located at `frontend/strava-goal/src/app/components/peak-picker/`

3. **Home Page Redesign** - Summit-focused dashboard:
   - Slim stats banner (distance, elevation, activities, summits)
   - Distance & elevation charts
   - Personal summit wishlist with add/remove functionality

4. **Sync Optimizations**:
   - Increased page size 30→200 for Strava API calls
   - Split into daily (30 days) and weekly (full) sync jobs
   - Added `summits_calculated` flag for incremental summit detection

5. **#hg Activities Fix** (Just Done):
   - Fixed `FetchAndStoreDetailedActivity` to include `MovingTime` and `Elevation`
   - Updated workflow to re-fetch #hg activities missing these fields
   - Triggered sync: `POST /hikegang/sync`

### Activity Types Filtered

Only these Strava activity types are synced:
- Run, Walk, Hike, VirtualRun, TrailRun

---

## Architecture Diagram

```
┌────────────────────────────────────────────────────────────────────┐
│                        EXTERNAL SERVICES                           │
│  ┌──────────────┐  ┌─────────────────┐  ┌──────────────────────┐  │
│  │  Strava API  │  │  OpenStreetMap  │  │  Cloudflare Tunnel   │  │
│  │  (OAuth +    │  │  (Peak data via │  │  (DNS + ingress)     │  │
│  │   webhooks)  │  │   Overpass API) │  │                      │  │
│  └──────┬───────┘  └────────┬────────┘  └──────────┬───────────┘  │
└─────────┼───────────────────┼──────────────────────┼──────────────┘
          │                   │                      │
          ▼                   ▼                      ▼
┌────────────────────────────────────────────────────────────────────┐
│              DigitalOcean Kubernetes Cluster                       │
│                                                                    │
│  ┌────────────────────────────────────────────────────────────┐   │
│  │                     Frontend (Angular)                      │   │
│  │  - Interactive Leaflet map with activity routes             │   │
│  │  - Group management UI                                      │   │
│  │  - Goal progress dashboards                                 │   │
│  │  - Strava OAuth login flow                                  │   │
│  └────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│                              ▼                                     │
│  ┌────────────────────────────────────────────────────────────┐   │
│  │                     Backend (Go REST API)                   │   │
│  │  Endpoints: /api/* (authenticated), /auth/*, /webhook/*     │   │
│  │                                                             │   │
│  │  Core Services:                                             │   │
│  │  - StravaService: OAuth, activity sync, webhook handling    │   │
│  │  - SummitService: Detects peaks crossed in activity routes  │   │
│  │  - GroupsService: Group & goal CRUD, member management      │   │
│  │  - PeakService: Peak data from OpenStreetMap                │   │
│  │  - GoalProgressService: Calculates group goal completion    │   │
│  │                                                             │   │
│  │  Background Jobs:                                           │   │
│  │  - Daily activity sync for all users                        │   │
│  │  - Summit detection on new activities                       │   │
│  └────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│  ┌──────────────────┐        │        ┌──────────────────────┐    │
│  │  Flux CD GitOps  │        │        │    Cloudflared Pod   │    │
│  │  (auto-deploy    │        │        │    (tunnel ingress)  │    │
│  │   from GHCR)     │        │        │                      │    │
│  └──────────────────┘        │        └──────────────────────┘    │
└──────────────────────────────┼────────────────────────────────────┘
                               │
                               ▼
┌────────────────────────────────────────────────────────────────────┐
│           DigitalOcean Managed PostgreSQL Database                 │
│                                                                    │
│  Tables:                                                           │
│  - users          (Strava auth tokens, user metadata)              │
│  - activity       (Synced Strava activities with GPS polylines)    │
│  - peaks          (Geographic peak data from OpenStreetMap)        │
│  - user_peaks     (Junction: which users summited which peaks)     │
│  - groups         (Group definitions)                              │
│  - group_members  (User-group membership with roles)               │
│  - group_goals    (Distance/elevation/summit targets)              │
└────────────────────────────────────────────────────────────────────┘
```

---

## Tech Stack

| Layer      | Technology                     |
|------------|--------------------------------|
| Frontend   | Angular 17+, TypeScript, Leaflet.js, SCSS |
| Backend    | Go 1.21+, net/http (no framework) |
| Database   | PostgreSQL 16 (DO Managed)     |
| Auth       | JWT + Strava OAuth 2.0         |
| CI/CD      | GitHub Actions → GHCR → Flux CD |
| Infra      | DigitalOcean Kubernetes (DOKS) |
| DNS/Tunnel | Cloudflare Tunnel              |

---

## Key Data Flows

### 1. User Login (Strava OAuth)
```
User → Frontend → /auth/strava/callback → Backend exchanges code → 
Stores tokens in `users` table → Returns JWT → Frontend stores JWT
```

### 2. Activity Sync
```
Strava Webhook → /webhook/strava → Backend fetches activity details →
Stores in `activity` table → Runs summit detection → 
Updates `user_peaks` if summit found
```

### 3. Group Goal Progress
```
Frontend requests /api/groups/{id}/goals → Backend queries:
- group_goals for targets
- activity table for member distances/elevation
- user_peaks for summit counts
→ Calculates percentage complete → Returns progress
```

---

## Project Structure

```
run-goals/
├── backend/              # Go REST API
│   ├── config/           # Environment config struct
│   ├── controllers/      # HTTP request handlers
│   ├── daos/             # Database access objects
│   ├── database/         # PostgreSQL connection
│   ├── dto/              # Request/response DTOs
│   ├── handlers/         # Route multiplexing
│   ├── middleware/       # JWT auth middleware
│   ├── models/           # Domain models
│   ├── services/         # Business logic
│   └── workflows/        # Background jobs
├── frontend/strava-goal/ # Angular SPA
│   └── src/app/
│       ├── components/   # Reusable UI components
│       ├── pages/        # Route pages
│       ├── services/     # API clients
│       └── guards/       # Route guards
├── database/             # PostgreSQL init scripts
│   └── sql/tables/       # Schema DDL files
├── k8s/                  # Kubernetes manifests
│   ├── app/              # App deployments & secrets
│   └── flux-system/      # GitOps configuration
└── docker-compose.yaml   # Local development stack
```

---

## Environment Variables (Backend)

| Variable               | Description                           |
|------------------------|---------------------------------------|
| `DATABASE_HOST`        | PostgreSQL host                       |
| `DATABASE_PORT`        | PostgreSQL port (5432 local, 25060 DO)|
| `DATABASE_USER`        | DB username                           |
| `DATABASE_PASSWORD`    | DB password                           |
| `DATABASE_DBNAME`      | Database name (`run_goals`)           |
| `DATABASE_SSLMODE`     | `disable` locally, `require` for DO   |
| `STRAVA_CLIENT_ID`     | Strava API app client ID              |
| `STRAVA_CLIENT_SECRET` | Strava API app client secret          |
| `JWT_SECRET`           | HMAC secret for JWT signing           |
| `SUMMIT_THRESHOLD_METERS` | Distance to consider summit (0.0007) |

---

## Previous Infrastructure

Originally hosted on 3x Raspberry Pi 4 nodes running k3s. Migrated to DigitalOcean in late 2024 for reliability and managed database benefits.

---

## Quick Commands

```bash
# Local dev (docker-compose)
docker compose up --build

# Check container logs
docker compose logs -f backend

# Access local DB
docker exec -it run-goals-db psql -U postgres -d run_goals

# Production (kubectl)
kubectl get pods -n default
kubectl logs deployment/summitseekers-backend
```

---

## Key Gotchas

1. **Strava Rate Limits**: Be careful with activity fetching during development
2. **Summit Detection**: Uses 0.0007 degree threshold (~70m) to match route to peaks
3. **Peak Data**: Fetched from OpenStreetMap Overpass API on backend startup (Western Cape region)
4. **Background Job**: Daily (recent 30 days) + Weekly (full) sync - see `workflows/useractivities.go`
5. **Managed DB SSL**: Production requires `sslmode=require`
6. **#hg Activities**: These are "HikeGang" activities fetched separately via detailed API (not list API) to get full data
7. **Admin Endpoints**: Use `/admin/refresh-peaks?admin_key=dev-admin-key` for admin operations (no JWT)

---

## Key Files to Know

| File | Purpose |
|------|---------|
| `backend/services/StravaService.go` | Strava OAuth, activity sync, webhook handling |
| `backend/services/summitService.go` | Detects peaks crossed in activity GPS routes |
| `backend/workflows/useractivities.go` | Background sync jobs (daily/weekly) |
| `backend/controllers/supportController.go` | Admin endpoints (delete account, refresh peaks) |
| `frontend/.../components/peak-picker/` | Reusable peak selection component |
| `frontend/.../pages/home-page/` | Main dashboard with stats, charts, wishlist |
| `database/sql/migrations/` | Database migration scripts |
