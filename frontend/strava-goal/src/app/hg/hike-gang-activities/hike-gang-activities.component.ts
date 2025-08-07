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
  markerClusterGroup!: L.MarkerClusterGroup;
  hgActivities: Activity[] = [];
  private polylinesById: Record<number, L.Polyline> = {};
  private lastHighlightedPolyline: L.Polyline | null = null;
  private homeBounds!: L.LatLngBounds;
  private isDragging = false;
  private startY = 0;
  private startHeight = 0;

  constructor(private hgService: HgService, private router: Router) {}

  get totalDistanceKm(): number {
    // Sum distances (in meters), then convert to kilometers
    return this.hgActivities.reduce((sum, a) => sum + a.distance, 0) / 1000;
  }

  ngOnInit(): void {
    console.log('HikeGang Activities: Component initializing...');
    this.hgService.loadActivities();
    this.hgService.activities$.subscribe((acts) => {
      console.log(
        'HikeGang Activities: Received activities:',
        acts?.length || 0
      );

      if (acts) {
        // Log all activity names for debugging
        console.log(
          'All activity names:',
          acts.map((a) => a.name)
        );

        // 1) Filter them by "#hg" (case insensitive)
        this.hgActivities = acts.filter((act) =>
          act.name?.toLowerCase().includes('#hg')
        );

        console.log('Filtered HG activities:', this.hgActivities.length);
        console.log(
          'HG activity names:',
          this.hgActivities.map((a) => a.name)
        );
        console.log(
          'Activities with polylines:',
          this.hgActivities.filter((a) => a.map_polyline).length
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
      console.error('Map element not found! Retrying in 500ms...');
      setTimeout(() => this.initMap(), 500);
      return;
    }

    console.log('Initializing map...');

    this.map = L.map('hgMap', {
      center: [-33.83, 18.6], // e.g., near Cape Town
      zoom: 12,
    });

    this.homeBounds = this.map.getBounds();

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
    }).addTo(this.map);

    this.markerClusterGroup = L.markerClusterGroup();
    this.map.addLayer(this.markerClusterGroup);

    console.log('Map initialized successfully');
  }

  resetMap(): void {
    this.map.fitBounds(this.homeBounds);
  }

  displayActivities(): void {
    console.log('Displaying activities:', this.hgActivities.length);
    console.log('Polyline library available:', typeof polyline !== 'undefined');

    for (const act of this.hgActivities) {
      console.log(
        'Processing activity:',
        act.name,
        'Has polyline:',
        !!act.map_polyline
      );

      if (act.map_polyline) {
        try {
          const decodedCoords = polyline.decode(act.map_polyline);
          console.log(
            'Decoded coordinates for',
            act.name,
            ':',
            decodedCoords.length,
            'points'
          );

          // Validate coordinates
          if (!decodedCoords || decodedCoords.length === 0) {
            console.warn('No coordinates decoded for:', act.name);
            continue;
          }

          const latLngs = decodedCoords
            .map((coords) => {
              if (!Array.isArray(coords) || coords.length < 2) {
                console.warn('Invalid coordinate pair:', coords);
                return null;
              }
              return L.latLng(coords[0], coords[1]);
            })
            .filter((coord): coord is L.LatLng => coord !== null); // Type guard filter

          if (latLngs.length === 0) {
            console.warn('No valid coordinates for:', act.name);
            continue;
          }

          const poly = L.polyline(latLngs, {
            color: act.has_summit ? 'green' : 'blue',
            weight: 3,
          }).addTo(this.map);

          console.log('Added polyline for:', act.name);

          // Save reference
          this.polylinesById[act.id] = poly;

          // Add click handler to zoom and highlight
          poly.on('click', () => {
            this.highlightActivity(act);
          });

          // Add popup
          poly.bindPopup(this.buildActivityPopup(act));
        } catch (error) {
          console.error(
            'Error processing polyline for activity:',
            act.name,
            error
          );
          console.error(
            'Polyline data:',
            act.map_polyline?.substring(0, 100) + '...'
          );

          // Log more details about the error
          console.error('Activity details:', {
            id: act.id,
            name: act.name,
            hasPolyline: !!act.map_polyline,
            polylineLength: act.map_polyline?.length,
          });
        }
      } else {
        console.warn('No polyline data for activity:', act.name);
      }
    }

    console.log(
      'Total polylines added to map:',
      Object.keys(this.polylinesById).length
    );
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
    const resizeHandle = document.querySelector(
      '.resize-handle'
    ) as HTMLElement;
    const mapContainer = document.getElementById('hgMap') as HTMLElement;

    if (resizeHandle && mapContainer) {
      resizeHandle.addEventListener('mousedown', this.startResize.bind(this));
      document.addEventListener('mousemove', this.doResize.bind(this));
      document.addEventListener('mouseup', this.stopResize.bind(this));
    }

    // Update the map size on window resize
    window.addEventListener('resize', this.onResize.bind(this));

    // Also, if you have a container that can change size, observe it
    const resizeObserver = new ResizeObserver(() => {
      this.map.invalidateSize();
    });

    if (mapContainer) {
      resizeObserver.observe(mapContainer);
    }
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
