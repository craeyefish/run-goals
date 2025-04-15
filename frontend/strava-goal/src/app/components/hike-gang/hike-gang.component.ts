import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import * as L from 'leaflet';
import * as polyline from '@mapbox/polyline';
import 'leaflet.markercluster';
import { ActivityService, Activity } from 'src/app/services/activity.service';

@Component({
  selector: 'app-hike-gang',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './hike-gang.component.html',
  styleUrls: ['./hike-gang.component.css'],
})
export class HikeGangComponent implements OnInit {
  map!: L.Map;
  markerClusterGroup!: L.MarkerClusterGroup;
  hgActivities: Activity[] = [];
  private polylinesById: Record<number, L.Polyline> = {};
  private lastHighlightedPolyline: L.Polyline | null = null;

  constructor(private activityService: ActivityService) {}

  ngOnInit(): void {
    this.activityService.loadActivities();
    this.activityService.activities$.subscribe((acts) => {
      if (acts) {
        // 1) Filter them by "#hg"
        this.hgActivities = acts.filter((act) => act.name?.includes('#hg'));

        // 2) Sort them by date descending (newest first)
        this.hgActivities.sort((a, b) => {
          // If `start_date` is a string, convert both to Date objects
          return (
            new Date(b.start_date).getTime() - new Date(a.start_date).getTime()
          );
        });

        // 3) Initialize the map, etc.
        this.initMap();
        this.displayActivities();
      }
    });
  }

  initMap(): void {
    this.map = L.map('hgMap', {
      center: [-33.83, 18.6], // e.g., near Cape Town
      zoom: 12,
    });

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
    }).addTo(this.map);

    this.markerClusterGroup = L.markerClusterGroup();
    this.map.addLayer(this.markerClusterGroup);
  }

  displayActivities(): void {
    for (const act of this.hgActivities) {
      if (act.map_polyline) {
        const decodedCoords = polyline.decode(act.map_polyline);
        const latLngs = decodedCoords.map((coords) =>
          L.latLng(coords[0], coords[1])
        );

        const poly = L.polyline(latLngs, {
          color: act.has_summit ? 'green' : 'blue',
          weight: 3,
        }).addTo(this.map);

        // Save reference so we can style it later
        this.polylinesById[act.id] = poly;

        // Existing popup logic
        poly.bindPopup(this.buildActivityPopup(act));
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
}
