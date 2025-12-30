# Map Feature Architecture

## Overview

The map feature displays user activities and summited peaks on an interactive Leaflet map. It combines GPS route data from Strava activities with geographical peak data.

## Key Components

### Frontend
- **Activity Map**: `frontend/strava-goal/src/app/components/activity-map/`
- **Peak Picker**: `frontend/strava-goal/src/app/components/peak-picker/` (reusable peak selection)

### Backend
- **PeakService**: Manages peak data from OpenStreetMap
- **SummitService**: Detects peaks crossed during activities

## Data Flow

### Activities Display
1. `GET /api/activities` returns activities with `map_polyline` (encoded GPS route)
2. Frontend decodes polyline using `@mapbox/polyline`
3. Routes displayed as colored polylines (green = has summit, blue = regular)

### Summit Display
1. `GET /api/peaks` returns all peaks with `is_summited` boolean per user
2. Peaks displayed as clustered markers
3. Icons: `summit-icon.png` (unvisited) / `summit-icon-green.png` (visited)

## Peak Data Model

```typescript
interface Peak {
  id: number;
  osm_id: number;
  latitude: number;
  longitude: number;
  name: string;
  elevation_meters: number;
  is_summited: boolean;
  // New metadata fields (Dec 2024)
  alt_name?: string;
  name_en?: string;
  region?: string;
  wikipedia?: string;
  wikidata?: string;
  description?: string;
  prominence?: number;
}
```

## Summit Detection Algorithm

Located in `backend/services/summitService.go`:
1. Decode activity polyline to coordinates
2. Find peaks within route bounding box
3. Calculate minimum distance from route to each peak
4. Mark summit if distance < threshold (0.0007 degrees ≈ 70m)
5. Store in `user_peaks` junction table

## Map Configuration

- **Center**: Cape Town, South Africa `[-33.9249, 18.4241]`
- **Default Zoom**: 7
- **Tile Layer**: OpenStreetMap
- **Clustering**: Enabled via `leaflet.markercluster`
