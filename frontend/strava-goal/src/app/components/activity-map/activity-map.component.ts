import { AfterViewInit, Component, OnInit } from '@angular/core';
import * as L from 'leaflet';
import * as polyline from '@mapbox/polyline';
import 'leaflet.markercluster';
import { ActivityService, Activity } from 'src/app/services/activity.service';
import { PeakService, Peak } from 'src/app/services/peak.service';
import { filter } from 'rxjs';

export const defaultPeakIcon = L.icon({
  iconUrl: 'assets/summit-icon.png', // your default summit icon
  iconSize: [32, 32],
  iconAnchor: [16, 16],
  popupAnchor: [0, -32],
});

export const visitedPeakIcon = L.icon({
  iconUrl: 'assets/summit-icon-green.png', // a green version for visited peaks
  iconSize: [32, 32],
  iconAnchor: [16, 16],
  popupAnchor: [0, -32],
});

@Component({
  selector: 'app-activity-map',
  standalone: true,
  templateUrl: './activity-map.component.html',
  styleUrls: ['./activity-map.component.scss'],
})
export class ActivityMapComponent implements OnInit, AfterViewInit {
  showPeaks = true;
  map!: L.Map;

  // Store the retrieved data
  activities: Activity[] = [];
  markerClusterGroup!: L.MarkerClusterGroup; // We'll store the cluster group here

  peaks: Peak[] = []; // or a separate array if you have multiple data sets

  constructor(
    private activityService: ActivityService,
    private peakService: PeakService
  ) {}

  ngOnInit(): void {
    this.initMap();

    this.activityService.loadActivities();
    this.activityService.activities$
      .pipe(filter((acts) => acts !== null))
      .subscribe((acts) => {
        this.activities = acts!;
        this.displayActivities();
      });

    this.peakService.loadPeaks();
    this.peakService.peaks$.pipe(filter((peaks) => peaks !== null)).subscribe({
      next: (peaks) => {
        this.peaks = peaks!;
        this.displayPeaks();
      },
      error: (error) => {
        console.error('Failed to load peaks for map:', error);
      },
    });
  }

  ngAfterViewInit() {
    setTimeout(() => {
      this.map.invalidateSize();
    }, 200);
  }

  initMap(): void {
    this.map = L.map('map', {
      center: [-33.9249, 18.4241], // e.g., near Cape Town
      zoom: 7,
    });

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
    }).addTo(this.map);

    // Initialize the cluster group
    this.markerClusterGroup = L.markerClusterGroup();
    // We'll add markers to this group once we fetch them
    this.map.addLayer(this.markerClusterGroup);
  }

  displayActivities(): void {
    for (const act of this.activities) {
      if (act.map_polyline) {
        const decodedCoords = polyline.decode(act.map_polyline);
        const latLngs = decodedCoords.map((coords) =>
          L.latLng(coords[0], coords[1])
        );

        const color = act.has_summit
          ? 'rgba(14, 212, 14, 0.61)'
          : 'rgba(0, 0, 255, 0.6)';

        // Default polyline style
        const poly = L.polyline(latLngs, { color: color, weight: 3 }).addTo(
          this.map
        );

        // Construct HTML for your popup
        const activityUrl = `https://www.strava.com/activities/${act.strava_activity_id}`;
        const infoHtml = `
          <strong>${act.name}</strong><br>
          Distance: ${(act.distance / 1000).toFixed(2)} km<br>
          Start Date: ${new Date(act.start_date).toLocaleString()}<br>
          <a href="${activityUrl}" target="_blank">View on Strava</a>
        `;

        // Bind the popup
        poly.bindPopup(infoHtml);

        // Highlight on popup open
        poly.on('popupopen', () => {
          poly.setStyle({ color: 'red', weight: 5 });
        });

        // Revert style on popup close
        poly.on('popupclose', () => {
          poly.setStyle({ color: 'blue', weight: 3 });
        });
      }
    }
  }

  displayPeaks(): void {
    // Only display peaks if the toggle is enabled
    if (!this.showPeaks) {
      return;
    }

    this.peaks.forEach((peak) => {
      // Choose the icon based on whether it's summited
      const iconToUse = peak.is_summited ? visitedPeakIcon : defaultPeakIcon;

      const marker = L.marker([peak.latitude, peak.longitude], {
        icon: iconToUse,
      });
      marker.bindPopup(this.buildPeakPopup(peak));

      // Add marker to the cluster group instead of the map directly
      this.markerClusterGroup.addLayer(marker);
    });
  }

  buildPeakPopup(peak: Peak): string {
    const summitStatus = peak.is_summited
      ? '<span style="color: green;">✅ Summited</span>'
      : '<span style="color: orange;">⭕ Not summited</span>';

    return `
      <strong>${peak.name || 'Unnamed Peak'}</strong><br>
      Elev: ${peak.elevation_meters ? `${peak.elevation_meters} m` : 'N/A'}<br>
      Status: ${summitStatus}
    `;
  }

  onTogglePeaks(event: Event): void {
    const inputElement = event.target as HTMLInputElement;
    this.showPeaks = inputElement.checked;

    if (this.showPeaks) {
      // Show peaks by calling displayPeaks to add them to cluster
      this.displayPeaks();
    } else {
      // Hide peaks by clearing and re-displaying only activities
      this.markerClusterGroup.clearLayers();
      this.displayActivities();
    }
  }
}
