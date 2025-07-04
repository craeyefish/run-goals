# Map Feature Architecture

## Description

Comprehensive overview of the map functionality in the Run Goals application, including activities and summit/peak display, data flow, and technical implementation details.

## Content

### High-Level Overview

The map feature is a core component that displays user activities and summited peaks on an interactive Leaflet map. It combines GPS route data from Strava activities with geographical peak data to provide visual tracking of fitness activities and summit achievements.

### Frontend Architecture

#### Main Component

- **Location**: `frontend/strava-goal/src/app/components/activity-map/activity-map.component.ts`
- **Page Integration**: Used in `pages/map-page/map-page.component.ts`
- **Technology**: Leaflet.js with marker clustering support

#### Key Dependencies

- **Leaflet**: Core mapping library (`import * as L from 'leaflet'`)
- **Polyline**: Route decoding (`import * as polyline from '@mapbox/polyline'`)
- **Clustering**: `import 'leaflet.markercluster'` for performance with many markers

#### Data Services

1. **ActivityService**: Loads user activities with GPS routes
2. **PeakService**: Loads peak/summit data with user summiting status
3. **APIs Used**:
   - `GET /api/activities` - User activities with GPS polylines
   - `GET /api/peaks` - Peak data with summiting status

### Data Flow

#### Activities Display

1. Component calls `activityService.loadActivities()`
2. Activities with `map_polyline` field are decoded using polyline library
3. Routes displayed as colored polylines (green if has summit, blue otherwise)
4. Activity popups show: name, distance, date, Strava link

#### Summit Display (Current Issue)

1. Component calls `peakService.loadPeaks()`
2. Peak data loaded but `displayPeaks()` method never called
3. Should display peaks as clustered markers with different icons:
   - `assets/summit-icon.png` - Unvisited peaks
   - `assets/summit-icon-green.png` - Visited/summited peaks

### Backend Data Models

#### Peak Data Structure

```go
// models/peak.go
type Peak struct {
    ID              int64   `json:"id"`
    OsmID           int64   `json:"osm_id"`
    Latitude        float64 `json:"latitude"`
    Longitude       float64 `json:"longitude"`
    Name            string  `json:"name"`
    ElevationMeters float64 `json:"elevation_meters"`
}

// models/peakSummited.go - API response format
type PeakSummited struct {
    Peak
    IsSummited bool `json:"is_summited"`
}
```

#### Activity Data Structure

```go
type Activity struct {
    // ... other fields
    MapPolyline string `json:"map_polyline"` // Encoded GPS route
    HasSummit   bool   `json:"has_summit"`   // Summit detected flag
}
```

### Summit Detection Logic

#### Backend Process

1. **Service**: `backend/services/summitService.go`
2. **Algorithm**:
   - Decodes activity polyline coordinates
   - Finds candidate peaks within route bounding box
   - Calculates minimum distance from route to each peak
   - Marks summit if distance < threshold (configurable)
3. **Storage**: Results stored in `user_peaks` table linking users, peaks, and activities

#### Database Relationships

- `peaks` table: All geographical peak data (from OpenStreetMap)
- `user_peaks` table: Junction table tracking which users summited which peaks
- `activities` table: GPS routes and summit detection flags

### Map Configuration

#### Default Settings

- **Center**: Cape Town, South Africa `[-33.9249, 18.4241]`
- **Zoom**: Level 7
- **Tile Layer**: OpenStreetMap
- **Clustering**: Enabled for performance

#### Toggle Controls

- **"Show Summits" checkbox**: Controls peak marker visibility
- **Current Issue**: Toggle exists but `displayPeaks()` never called

### Key Technical Patterns

#### Reactive Data Loading

```typescript
// Activities
this.activityService.activities$
  .pipe(filter((acts) => acts !== null))
  .subscribe((acts) => {
    this.activities = acts!;
    this.displayActivities();
  });

// Peaks (missing displayPeaks() call)
this.peakService.peaks$
  .pipe(filter((peaks) => peaks !== null))
  .subscribe((peaks) => {
    this.peaks = peaks!;
    // MISSING: this.displayPeaks();
  });
```

#### Marker Clustering

- Single cluster group for both activities and peaks
- Performance optimization for many markers
- Markers added to cluster group, not directly to map

#### Polyline Styling

```typescript
const color = act.has_summit
  ? 'rgba(14, 212, 14, 0.61)' // Green for summit activities
  : 'rgba(0, 0, 255, 0.6)'; // Blue for regular activities
```

### Frontend Interface Definitions

#### Peak Interface

```typescript
// services/peak.service.ts
export interface Peak {
  id: number;
  osm_id: number;
  latitude: number;
  longitude: number;
  name: string;
  elevation_meters: number;
  is_summited: boolean; // User context flag
}
```

### Performance Considerations

- **Clustering**: Prevents map overload with many markers
- **Lazy Loading**: Data loaded on component initialization
- **Caching**: Services use BehaviorSubject to cache loaded data
- **Route Filtering**: Only activities with polylines are displayed

### Current Known Issues

1. **Summit Display**: `displayPeaks()` method exists but never called
2. **Toggle Functionality**: Peak display commented out in toggle handler
3. **User Context**: Need to verify `is_summited` field properly set by backend

### Integration Points

- **Groups Feature**: Summit goals can target specific peaks
- **Progress Tracking**: Summit achievements count toward group and personal goals
- **Activity Sync**: Strava activities automatically analyzed for summit detection
- **Profile Stats**: Summit counts displayed in user profiles

### Future Enhancement Areas

- Real-time summit detection during activity import
- Advanced filtering (by elevation, region, completion status)
- Route planning with target peaks
- Social features (shared summit achievements)
- Offline map support for remote areas
