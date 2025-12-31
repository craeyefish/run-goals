import { Component, OnInit, OnDestroy, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import * as L from 'leaflet';
import * as polyline from '@mapbox/polyline';
import 'leaflet.markercluster';
import { Subject, combineLatest, forkJoin } from 'rxjs';
import { takeUntil, filter } from 'rxjs/operators';

import { ActivityService, Activity } from '../../services/activity.service';
import { PeakService, Peak } from '../../services/peak.service';
import { PeakSummitService } from '../../services/peak-summit.service';
import { ChallengeService } from '../../services/challenge.service';
import { ChallengeWithProgress, ChallengePeakWithDetails } from '../../models/challenge.model';

// Icons
const defaultPeakIcon = L.icon({
    iconUrl: 'assets/summit-icon.png',
    iconSize: [32, 32],
    iconAnchor: [16, 16],
    popupAnchor: [0, -32],
});

const visitedPeakIcon = L.icon({
    iconUrl: 'assets/summit-icon-green.png',
    iconSize: [32, 32],
    iconAnchor: [16, 16],
    popupAnchor: [0, -32],
});

const challengePeakIcon = L.icon({
    iconUrl: 'assets/summit-icon.png',
    iconSize: [40, 40],
    iconAnchor: [20, 20],
    popupAnchor: [0, -40],
    className: 'challenge-peak-icon',
});

interface PeakWithDetails extends Peak {
    summitCount?: number;
    lastSummitDate?: Date;
    challengeNames?: string[];
}

interface ChallengeWithPeaks extends ChallengeWithProgress {
    peaks?: ChallengePeakWithDetails[];
    isExpanded?: boolean;
}

type TabType = 'peaks' | 'activities' | 'challenges';
type FilterType = 'all' | 'summited' | 'not-summited';

@Component({
    selector: 'app-explore',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './explore.component.html',
    styleUrls: ['./explore.component.scss'],
})
export class ExploreComponent implements OnInit, OnDestroy, AfterViewInit {
    // Map
    map!: L.Map;
    markerClusterGroup!: L.MarkerClusterGroup;
    activityPolylines: L.Polyline[] = [];
    peakMarkers: Map<number, L.Marker> = new Map();

    // Data
    peaks: PeakWithDetails[] = [];
    filteredPeaks: PeakWithDetails[] = [];
    activities: Activity[] = [];
    challenges: ChallengeWithPeaks[] = [];

    // UI State
    activeTab: TabType = 'peaks';
    peakFilter: FilterType = 'all';
    searchQuery = '';
    showActivities = true;
    showPeaks = true;
    selectedChallengeId: number | null = null;
    expandedChallengeId: number | null = null;
    challengeFilterActive = false;

    // Track highlighted activity
    highlightedActivityPolyline: L.Polyline | null = null;
    highlightedActivityOriginalColor: string | null = null;

    // Selected Peak Modal
    selectedPeak: PeakWithDetails | null = null;
    peakChallenges: ChallengeWithProgress[] = [];

    // Loading
    loading = true;

    private destroy$ = new Subject<void>();

    constructor(
        private activityService: ActivityService,
        private peakService: PeakService,
        private peakSummitService: PeakSummitService,
        private challengeService: ChallengeService,
        private router: Router
    ) { }

    ngOnInit(): void {
        this.loadData();
    }

    ngAfterViewInit(): void {
        setTimeout(() => {
            this.initMap();
        }, 100);
    }

    ngOnDestroy(): void {
        this.destroy$.next();
        this.destroy$.complete();
        if (this.map) {
            this.map.remove();
        }
    }

    initMap(): void {
        this.map = L.map('explore-map', {
            center: [-33.9249, 18.4241],
            zoom: 9,
        });

        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; OpenStreetMap contributors',
        }).addTo(this.map);

        this.markerClusterGroup = L.markerClusterGroup({
            maxClusterRadius: 50,
            spiderfyOnMaxZoom: true,
            showCoverageOnHover: false,
        });
        this.map.addLayer(this.markerClusterGroup);

        // Resize handling
        setTimeout(() => this.map.invalidateSize(), 200);
    }

    loadData(): void {
        this.activityService.loadActivities();
        this.peakService.loadPeaks();
        this.challengeService.loadUserChallenges();

        combineLatest([
            this.activityService.activities$.pipe(filter(a => a !== null)),
            this.peakService.peaks$.pipe(filter(p => p !== null)),
            this.peakSummitService.getPeakSummaries(),
        ]).pipe(takeUntil(this.destroy$)).subscribe({
            next: ([activities, peaks, summaries]) => {
                // Sort activities by date (latest first)
                this.activities = (activities || []).sort((a, b) =>
                    new Date(b.start_date).getTime() - new Date(a.start_date).getTime()
                );

                // Enrich peaks with summit data
                this.peaks = (peaks || []).map(peak => {
                    const summary = summaries.find((s: any) => s.peak_id === peak.id);
                    return {
                        ...peak,
                        summitCount: summary?.activities?.length || 0,
                        lastSummitDate: summary?.activities?.[0]?.start_date
                            ? new Date(summary.activities[0].start_date)
                            : undefined,
                    };
                });

                this.applyPeakFilter();
                this.loading = false;

                // Display on map after data loads
                setTimeout(() => {
                    this.displayActivities();
                    this.displayPeaks();
                }, 300);
            },
            error: (err) => {
                console.error('Error loading explore data:', err);
                this.loading = false;
            },
        });

        // Load challenges separately
        this.challengeService.getUserChallenges().pipe(
            takeUntil(this.destroy$)
        ).subscribe({
            next: (response) => {
                this.challenges = response.challenges || [];
            },
            error: (err) => console.error('Error loading challenges:', err),
        });
    }

    // ==================== Map Display ====================

    displayActivities(): void {
        // Clear existing
        this.activityPolylines.forEach(p => this.map.removeLayer(p));
        this.activityPolylines = [];

        if (!this.showActivities) return;

        for (const act of this.activities) {
            if (act.map_polyline) {
                const decodedCoords = polyline.decode(act.map_polyline);
                const latLngs = decodedCoords.map(coords => L.latLng(coords[0], coords[1]));

                const color = act.has_summit
                    ? 'rgba(14, 212, 14, 0.6)'
                    : 'rgba(0, 100, 255, 0.5)';

                const poly = L.polyline(latLngs, { color, weight: 3 }).addTo(this.map);

                const activityUrl = `https://www.strava.com/activities/${act.strava_activity_id}`;
                poly.bindPopup(`
          <strong>${act.name}</strong><br>
          Distance: ${(act.distance / 1000).toFixed(2)} km<br>
          Elevation: ${act.total_elevation_gain?.toFixed(0) || 0} m<br>
          <a href="${activityUrl}" target="_blank">View on Strava</a>
        `);

                poly.on('popupopen', () => poly.setStyle({ color: '#fc4c02', weight: 5 }));
                poly.on('popupclose', () => poly.setStyle({ color, weight: 3 }));

                this.activityPolylines.push(poly);
            }
        }
    }

    displayPeaks(): void {
        this.markerClusterGroup.clearLayers();
        this.peakMarkers.clear();

        if (!this.showPeaks) return;

        this.filteredPeaks.forEach(peak => {
            const iconToUse = peak.is_summited ? visitedPeakIcon : defaultPeakIcon;
            const marker = L.marker([peak.latitude, peak.longitude], { icon: iconToUse });

            marker.bindPopup(this.buildPeakPopup(peak));
            marker.on('click', () => this.openPeakDetail(peak));

            this.markerClusterGroup.addLayer(marker);
            this.peakMarkers.set(peak.id, marker);
        });
    }

    buildPeakPopup(peak: PeakWithDetails): string {
        const status = peak.is_summited ? '✅ Summited' : '⬜ Not summited';
        const summitInfo = peak.summitCount
            ? `<br>Summit count: ${peak.summitCount}`
            : '';

        return `
      <strong>${peak.name}</strong><br>
      ${peak.elevation_meters?.toFixed(0) || '?'} m<br>
      ${status}${summitInfo}
    `;
    }

    // ==================== Filtering ====================

    applyPeakFilter(): void {
        let filtered = [...this.peaks];

        // Apply filter
        if (this.peakFilter === 'summited') {
            filtered = filtered.filter(p => p.is_summited);
        } else if (this.peakFilter === 'not-summited') {
            filtered = filtered.filter(p => !p.is_summited);
        }

        // Apply search
        if (this.searchQuery.trim()) {
            const query = this.searchQuery.toLowerCase();
            filtered = filtered.filter(p =>
                p.name.toLowerCase().includes(query) ||
                p.region?.toLowerCase().includes(query)
            );
        }

        // Sort: summited first, then by name
        filtered.sort((a, b) => {
            if (a.is_summited !== b.is_summited) {
                return a.is_summited ? -1 : 1;
            }
            return a.name.localeCompare(b.name);
        });

        this.filteredPeaks = filtered;

        // Update map if already initialized
        if (this.map) {
            this.displayPeaks();
        }
    }

    onSearchChange(): void {
        this.applyPeakFilter();
    }

    onFilterChange(): void {
        this.applyPeakFilter();
    }

    toggleActivities(): void {
        this.showActivities = !this.showActivities;
        this.displayActivities();
    }

    togglePeaks(): void {
        this.showPeaks = !this.showPeaks;
        this.displayPeaks();
    }

    // ==================== Peak Detail Modal ====================

    openPeakDetail(peak: PeakWithDetails): void {
        this.selectedPeak = peak;

        // Find challenges containing this peak
        // For now, we'll need to check each challenge's peaks
        // This could be optimized with backend support
        this.peakChallenges = [];
    }

    closePeakDetail(): void {
        this.selectedPeak = null;
        this.peakChallenges = [];
    }

    // ==================== Navigation ====================

    setActiveTab(tab: TabType): void {
        this.activeTab = tab;
    }

    focusPeakOnMap(peak: Peak): void {
        this.map.setView([peak.latitude, peak.longitude], 14);

        const marker = this.peakMarkers.get(peak.id);
        if (marker) {
            marker.openPopup();
        }
    }

    focusActivityOnMap(activity: Activity): void {
        // Reset previously highlighted activity
        if (this.highlightedActivityPolyline && this.highlightedActivityOriginalColor) {
            this.highlightedActivityPolyline.setStyle({
                color: this.highlightedActivityOriginalColor,
                weight: 3
            });
        }

        if (activity.map_polyline) {
            const decodedCoords = polyline.decode(activity.map_polyline);
            if (decodedCoords.length > 0) {
                const bounds = L.latLngBounds(decodedCoords.map(c => L.latLng(c[0], c[1])));
                this.map.fitBounds(bounds, { padding: [50, 50] });

                // Find the polyline for this activity and highlight it
                const activityIdx = this.activities.findIndex(a => a.id === activity.id);
                if (activityIdx >= 0 && this.activityPolylines[activityIdx]) {
                    const polylineLayer = this.activityPolylines[activityIdx];
                    this.highlightedActivityOriginalColor = activity.has_summit
                        ? 'rgba(14, 212, 14, 0.6)'
                        : 'rgba(0, 100, 255, 0.5)';
                    this.highlightedActivityPolyline = polylineLayer;

                    // Highlight, bring to front, and open popup
                    polylineLayer.setStyle({ color: '#fc4c02', weight: 5 });
                    polylineLayer.bringToFront();
                    polylineLayer.openPopup();
                }
            }
        }
    }

    navigateToChallenge(challengeId: number): void {
        this.router.navigate(['/challenges', challengeId]);
    }

    openStravaActivity(activity: Activity): void {
        window.open(`https://www.strava.com/activities/${activity.strava_activity_id}`, '_blank');
    }

    // ==================== Challenge Expansion ====================

    toggleChallengeExpand(challenge: ChallengeWithPeaks, event: Event): void {
        event.stopPropagation();

        if (this.expandedChallengeId === challenge.id) {
            // Collapse - reset view
            this.expandedChallengeId = null;
            this.challengeFilterActive = false;
            this.applyPeakFilter(); // Reset to show all peaks
        } else {
            // Expand and load peaks
            this.expandedChallengeId = challenge.id;
            this.challengeFilterActive = true;

            // Load challenge peaks if not already loaded
            if (!challenge.peaks) {
                console.log('Loading peaks for challenge:', challenge.id, challenge.name);
                this.challengeService.getChallengePeaks(challenge.id).pipe(
                    takeUntil(this.destroy$)
                ).subscribe({
                    next: (peaks) => {
                        console.log('Loaded challenge peaks:', peaks);
                        challenge.peaks = peaks;
                        challenge.isExpanded = true;
                        this.filterAndZoomToChallengePeaks(challenge);
                    },
                    error: (err) => console.error('Error loading challenge peaks:', err),
                });
            } else {
                challenge.isExpanded = true;
                this.filterAndZoomToChallengePeaks(challenge);
            }
        }
    }

    filterAndZoomToChallengePeaks(challenge: ChallengeWithPeaks): void {
        if (!challenge.peaks || challenge.peaks.length === 0) return;

        // Filter peaks to only show challenge peaks
        const challengePeakIds = new Set(challenge.peaks.map(p => p.peakId));
        this.filteredPeaks = this.peaks.filter(p => challengePeakIds.has(p.id));

        // Update map
        this.displayPeaks();

        // Zoom to fit all challenge peaks
        const peakCoords = challenge.peaks.map(p =>
            L.latLng(p.latitude, p.longitude)
        );

        if (peakCoords.length > 0) {
            const bounds = L.latLngBounds(peakCoords);
            this.map.fitBounds(bounds, { padding: [80, 80] });
        }
    }

    focusChallengePeakOnMap(challengePeak: ChallengePeakWithDetails, event: Event): void {
        event.stopPropagation();

        this.map.setView([challengePeak.latitude, challengePeak.longitude], 14);

        // Find and open the marker popup
        const marker = this.peakMarkers.get(challengePeak.peakId);
        if (marker) {
            marker.openPopup();
        }
    }

    clearChallengeFilter(): void {
        this.expandedChallengeId = null;
        this.challengeFilterActive = false;
        this.challenges.forEach(c => c.isExpanded = false);
        this.applyPeakFilter();
    }

    // ==================== Stats ====================

    get summitedCount(): number {
        return this.peaks.filter(p => p.is_summited).length;
    }

    get totalPeaks(): number {
        return this.peaks.length;
    }

    get totalActivities(): number {
        return this.activities.length;
    }

    get activeChallengesCount(): number {
        return this.challenges.filter(c => !c.isCompleted).length;
    }
}
