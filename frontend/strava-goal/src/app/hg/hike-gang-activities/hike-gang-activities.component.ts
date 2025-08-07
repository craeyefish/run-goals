import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import * as L from 'leaflet';
import * as polyline from '@mapbox/polyline';
import 'leaflet.markercluster';
import { Router } from '@angular/router';
import { Activity, HgService } from 'src/app/services/hg.service';

@Component({
  selector: 'app-hike-gang-activities',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './hike-gang-activities.component.html',
  styleUrls: ['./hike-gang-activities.component.scss'],
})
export class HikeGangActivitiesComponent implements OnInit {
  map!: L.Map;
  markerClusterGroup!: L.LayerGroup | L.MarkerClusterGroup; // Support both types
  hgActivities: Activity[] = [];
  private polylinesById: Record<number, L.Polyline> = {};
  private lastHighlightedPolyline: L.Polyline | null = null;
  private homeBounds!: L.LatLngBounds;
  private isDragging = false;
  private startY = 0;
  private startHeight = 0;

  // Bound methods for event listeners
  private boundStartResize = this.startResize.bind(this);
  private boundDoResize = this.doResize.bind(this);
  private boundStopResize = this.stopResize.bind(this);

  constructor(private hgService: HgService, private router: Router) {}

  get totalDistanceKm(): number {
    // Sum distances (in meters), then convert to kilometers
    return this.hgActivities.reduce((sum, a) => sum + a.distance, 0) / 1000;
  }

  ngOnInit(): void {
    this.hgService.loadActivities();
    this.hgService.activities$.subscribe((acts) => {
      if (acts) {
        // 1) Filter them by "#hg" (case insensitive)
        this.hgActivities = acts.filter((act) =>
          act.name?.toLowerCase().includes('#hg')
        );

        // 2) Sort them by date descending (newest first)
        this.hgActivities.sort((a, b) => {
          // If `start_date` is a string, convert both to Date objects
          return (
            new Date(b.start_date).getTime() - new Date(a.start_date).getTime()
          );
        });

        // 3) Initialize the map, etc.
        setTimeout(() => {
          this.initMap();
          this.displayActivities();
          this.setupResizeHandlers();
        }, 100); // Small delay to ensure DOM is ready
      }
    });
  }

  goBack(): void {
    this.router.navigate(['/hg']); // Replace '/hg' with the correct route for your home page
  }

  initMap(): void {
    const mapElement = document.getElementById('hgMap');
    if (!mapElement) {
      setTimeout(() => this.initMap(), 500);
      return;
    }

    try {
      this.map = L.map('hgMap', {
        center: [-33.83, 18.6], // e.g., near Cape Town
        zoom: 12,
      });

      this.homeBounds = this.map.getBounds();

      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; OpenStreetMap contributors',
      }).addTo(this.map);

      // Check if markerClusterGroup is available
      if (typeof (L as any).markerClusterGroup === 'function') {
        this.markerClusterGroup = (L as any).markerClusterGroup();
      } else {
        // Create a simple layer group as fallback
        this.markerClusterGroup = L.layerGroup();
      }

      this.map.addLayer(this.markerClusterGroup);
    } catch (error) {
      setTimeout(() => this.initMap(), 1000);
    }
  }

  resetMap(): void {
    this.map.fitBounds(this.homeBounds);
  }

  displayActivities(): void {
    for (const act of this.hgActivities) {
      if (act.map_polyline) {
        try {
          const decodedCoords = polyline.decode(act.map_polyline);

          // Validate coordinates
          if (!decodedCoords || decodedCoords.length === 0) {
            continue;
          }

          const latLngs = decodedCoords
            .map((coords) => {
              if (!Array.isArray(coords) || coords.length < 2) {
                return null;
              }
              return L.latLng(coords[0], coords[1]);
            })
            .filter((coord): coord is L.LatLng => coord !== null); // Type guard filter

          if (latLngs.length === 0) {
            continue;
          }

          const poly = L.polyline(latLngs, {
            color: act.has_summit ? 'green' : 'blue',
            weight: 3,
          }).addTo(this.map);

          // Save reference
          this.polylinesById[act.id] = poly;

          // Add click handler to zoom and highlight
          poly.on('click', () => {
            this.highlightActivity(act);
          });

          // Add popup
          poly.bindPopup(this.buildActivityPopup(act));
        } catch (error) {
          // Skip activities with invalid polylines
        }
      } else {
        // Skip activities without polyline data
      }
    }
  }

  highlightActivity(act: Activity): void {
    // Find the polyline for this activity
    const poly = this.polylinesById[act.id];
    if (!poly) return;

    // Optionally reset the previously highlighted polyline
    if (this.lastHighlightedPolyline) {
      this.lastHighlightedPolyline.setStyle({ color: 'blue', weight: 3 });
    }

    // Highlight this one
    poly.bringToFront(); // ensure it’s on top
    poly.setStyle({ color: 'red', weight: 5 });

    // Remember it
    this.lastHighlightedPolyline = poly;

    // Optionally zoom/fit the map to the polyline
    const bounds = poly.getBounds();
    this.map.fitBounds(bounds, {
      padding: [50, 50],
    });
  }

  buildActivityPopup(act: Activity): string {
    const badges = this.extractBadges(act);

    let popupHtml = `
      <strong>${act.name}</strong><br>
      Distance: ${(act.distance / 1000).toFixed(2)} km<br>
      Start Date: ${new Date(act.start_date).toLocaleString()}<br>
      Badges: ${badges.join(', ')}<br>
    `;

    // If there's a photoUrl, add an <img> with a max width to keep it from getting huge
    if (act.photo_url) {
      popupHtml += `
        <img src="${act.photo_url}" alt="Hike Photo" style="max-width:200px; margin-top: 8px;"/>
      `;
    }

    return popupHtml;
  }

  /**
   * Very simple "badge" logic:
   * If the activity description has #peak -> 'Peak Badge'
   * If the activity description has #night -> 'Night Hike Badge'
   * etc
   */
  extractBadges(act: Activity): string[] {
    const badges: string[] = [];
    const desc = ''; //act.description || '';

    if (desc.includes('#peak')) badges.push('Peak Badge');
    if (desc.includes('#night')) badges.push('Night Hike Badge');
    // You can add more logic or do a regex scan for #whatever

    return badges;
  }

  setupResizeHandlers(): void {
    // Use a longer delay to ensure DOM is ready
    setTimeout(() => {
      const resizeHandle = document.querySelector(
        '.resize-handle'
      ) as HTMLElement;
      const mapContainer = document.getElementById('hgMap') as HTMLElement;

      if (resizeHandle && mapContainer) {
        // Remove any existing listeners first
        resizeHandle.removeEventListener('mousedown', this.boundStartResize);
        document.removeEventListener('mousemove', this.boundDoResize);
        document.removeEventListener('mouseup', this.boundStopResize);

        // Add the listeners
        resizeHandle.addEventListener('mousedown', this.boundStartResize);
        document.addEventListener('mousemove', this.boundDoResize);
        document.addEventListener('mouseup', this.boundStopResize);

        // Add touch support for mobile/tablet devices
        resizeHandle.addEventListener('touchstart', (e) => {
          const touch = e.touches[0];
          this.boundStartResize({
            clientY: touch.clientY,
            preventDefault: () => e.preventDefault(),
          } as MouseEvent);
        });

        document.addEventListener('touchmove', (e) => {
          if (this.isDragging) {
            const touch = e.touches[0];
            this.boundDoResize({
              clientY: touch.clientY,
            } as MouseEvent);
          }
        });

        document.addEventListener('touchend', () => {
          if (this.isDragging) {
            this.boundStopResize();
          }
        });
      } else {
        // Retry after another delay
        setTimeout(() => this.setupResizeHandlers(), 1000);
      }

      // Update the map size on window resize
      window.addEventListener('resize', this.onResize.bind(this));

      // Also, if you have a container that can change size, observe it
      const resizeObserver = new ResizeObserver(() => {
        if (this.map) {
          this.map.invalidateSize();
        }
      });

      if (mapContainer) {
        resizeObserver.observe(mapContainer);
      }
    }, 500); // Increased delay for production
  }

  startResize(event: MouseEvent): void {
    this.isDragging = true;
    this.startY = event.clientY;

    const mapContainer = document.getElementById('hgMap') as HTMLElement;
    const resizeHandle = document.querySelector(
      '.resize-handle'
    ) as HTMLElement;

    if (mapContainer) {
      this.startHeight = mapContainer.offsetHeight;
    }

    if (resizeHandle) {
      resizeHandle.classList.add('dragging');
    }

    document.body.classList.add('resizing');
    event.preventDefault();
  }

  doResize(event: MouseEvent): void {
    if (!this.isDragging) return;

    const mapContainer = document.getElementById('hgMap') as HTMLElement;
    if (!mapContainer) return;

    const deltaY = event.clientY - this.startY;
    const newHeight = Math.max(
      300,
      Math.min(window.innerHeight * 0.8, this.startHeight + deltaY)
    );

    mapContainer.style.height = `${newHeight}px`;
    mapContainer.style.flex = `0 0 ${newHeight}px`;

    // Invalidate map size after a short delay to ensure proper rendering
    setTimeout(() => {
      if (this.map) {
        this.map.invalidateSize();
      }
    }, 10);
  }

  stopResize(): void {
    this.isDragging = false;

    const resizeHandle = document.querySelector(
      '.resize-handle'
    ) as HTMLElement;
    if (resizeHandle) {
      resizeHandle.classList.remove('dragging');
    }

    document.body.classList.remove('resizing');
  }

  onResize(): void {
    // Debounced resize handler
    if (this.map) {
      this.map.invalidateSize();
    }
  }
}
