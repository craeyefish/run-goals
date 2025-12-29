import {
    Component,
    OnInit,
    OnDestroy,
    AfterViewInit,
    Input,
    Output,
    EventEmitter,
    ElementRef,
    ViewChild,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import * as L from 'leaflet';
import 'leaflet.markercluster';
import { PeakService, Peak } from '../../services/peak.service';
import { Subject } from 'rxjs';
import { takeUntil, filter } from 'rxjs/operators';

// Icons matching the activity-map component
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

// Selected peak uses a larger version with a highlight ring effect
const selectedPeakIcon = L.divIcon({
    className: 'selected-peak-marker',
    html: '<div class="selected-ring"><img src="assets/summit-icon-green.png" /></div>',
    iconSize: [44, 44],
    iconAnchor: [22, 22],
    popupAnchor: [0, -22],
});

export interface SelectedPeak {
    id: number;
    name: string;
    elevation_meters: number;
    is_summited: boolean;
    latitude: number;
    longitude: number;
}

@Component({
    selector: 'app-peak-picker',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './peak-picker.component.html',
    styleUrls: ['./peak-picker.component.scss'],
})
export class PeakPickerComponent implements OnInit, OnDestroy, AfterViewInit {
    @ViewChild('mapContainer') mapContainer!: ElementRef<HTMLDivElement>;

    // Input: IDs of peaks that are already selected (e.g., in wishlist)
    @Input() initialSelectedIds: number[] = [];

    // Output: Emit when selection changes
    @Output() selectionChange = new EventEmitter<SelectedPeak[]>();
    @Output() peakAdded = new EventEmitter<SelectedPeak>();
    @Output() peakRemoved = new EventEmitter<SelectedPeak>();
    @Output() close = new EventEmitter<void>();

    private destroy$ = new Subject<void>();
    private map!: L.Map;
    private markerClusterGroup!: L.MarkerClusterGroup;
    private peakMarkers: Map<number, L.Marker> = new Map();

    allPeaks: Peak[] = [];
    selectedPeaks: SelectedPeak[] = [];

    // Search functionality
    searchQuery = '';
    filteredPeaks: Peak[] = [];
    showSearchResults = false;

    loading = true;

    constructor(private peakService: PeakService) { }

    ngOnInit(): void {
        // Load peaks if not already loaded
        this.peakService.loadPeaks();

        this.peakService.peaks$
            .pipe(
                takeUntil(this.destroy$),
                filter((peaks) => peaks !== null)
            )
            .subscribe((peaks) => {
                this.allPeaks = peaks!;
                this.loading = false;

                // Initialize selected peaks from input IDs
                this.selectedPeaks = this.allPeaks
                    .filter(p => this.initialSelectedIds.includes(p.id))
                    .map(p => ({
                        id: p.id,
                        name: p.name,
                        elevation_meters: p.elevation_meters,
                        is_summited: p.is_summited,
                        latitude: p.latitude,
                        longitude: p.longitude,
                    }));

                // Display peaks on map if map is ready
                if (this.map) {
                    this.displayPeaks();
                }
            });
    }

    ngAfterViewInit(): void {
        // Small delay to ensure container is rendered
        setTimeout(() => {
            this.initMap();
            if (this.allPeaks.length > 0) {
                this.displayPeaks();
            }
        }, 100);
    }

    ngOnDestroy(): void {
        this.destroy$.next();
        this.destroy$.complete();
        if (this.map) {
            this.map.remove();
        }
    }

    private initMap(): void {
        if (!this.mapContainer?.nativeElement) {
            console.error('Map container not found');
            return;
        }

        this.map = L.map(this.mapContainer.nativeElement, {
            center: [-33.9249, 18.4241], // Default: Cape Town area
            zoom: 8,
        });

        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; OpenStreetMap contributors',
        }).addTo(this.map);

        // Initialize marker cluster group
        this.markerClusterGroup = L.markerClusterGroup({
            maxClusterRadius: 50,
            spiderfyOnMaxZoom: true,
            showCoverageOnHover: false,
        });
        this.map.addLayer(this.markerClusterGroup);

        // Invalidate size after a short delay
        setTimeout(() => {
            this.map.invalidateSize();
        }, 200);
    }

    private displayPeaks(): void {
        this.markerClusterGroup.clearLayers();
        this.peakMarkers.clear();

        this.allPeaks.forEach((peak) => {
            const isSelected = this.selectedPeaks.some(sp => sp.id === peak.id);
            const icon = this.getIconForPeak(peak, isSelected);

            const marker = L.marker([peak.latitude, peak.longitude], { icon });

            // Build popup with action button
            marker.bindPopup(this.buildPeakPopup(peak, isSelected));

            // Handle popup open to attach click handlers
            marker.on('popupopen', () => {
                this.attachPopupHandlers(peak, marker);
            });

            this.peakMarkers.set(peak.id, marker);
            this.markerClusterGroup.addLayer(marker);
        });
    }

    private getIconForPeak(peak: Peak, isSelected: boolean): L.Icon | L.DivIcon {
        if (isSelected) {
            return selectedPeakIcon;
        }
        return peak.is_summited ? visitedPeakIcon : defaultPeakIcon;
    }

    private buildPeakPopup(peak: Peak, isSelected: boolean): string {
        const summitStatus = peak.is_summited
            ? '<span class="status-summited">✅ Summited</span>'
            : '<span class="status-not-summited">⭕ Not summited</span>';

        const buttonText = isSelected ? 'Remove from List' : 'Add to List';
        const buttonClass = isSelected ? 'popup-btn-remove' : 'popup-btn-add';

        // Build location info (region or coordinates)
        const locationInfo = peak.region
            ? `<span class="popup-region">📍 ${peak.region}</span><br>`
            : '';

        // Show alternate name if available and different
        const altNameInfo = peak.alt_name && peak.alt_name !== peak.name
            ? `<span class="popup-alt-name">(${peak.alt_name})</span><br>`
            : '';

        return `
      <div class="peak-popup">
        <strong>${peak.name || 'Unnamed Peak'}</strong>${altNameInfo ? '<br>' + altNameInfo : ''}<br>
        ${locationInfo}
        Elev: ${peak.elevation_meters ? `${peak.elevation_meters.toFixed(0)} m` : 'N/A'}<br>
        ${summitStatus}<br>
        <button class="popup-btn ${buttonClass}" data-peak-id="${peak.id}">
          ${buttonText}
        </button>
      </div>
    `;
    }

    private attachPopupHandlers(peak: Peak, marker: L.Marker): void {
        // Find the button in the popup and attach click handler
        setTimeout(() => {
            const btn = document.querySelector(`.popup-btn[data-peak-id="${peak.id}"]`);
            if (btn) {
                btn.addEventListener('click', () => {
                    this.togglePeakSelection(peak);
                    marker.closePopup();
                });
            }
        }, 10);
    }

    togglePeakSelection(peak: Peak): void {
        const index = this.selectedPeaks.findIndex(sp => sp.id === peak.id);
        const selectedPeak: SelectedPeak = {
            id: peak.id,
            name: peak.name,
            elevation_meters: peak.elevation_meters,
            is_summited: peak.is_summited,
            latitude: peak.latitude,
            longitude: peak.longitude,
        };

        if (index >= 0) {
            // Remove from selection
            this.selectedPeaks.splice(index, 1);
            this.peakRemoved.emit(selectedPeak);
        } else {
            // Add to selection
            this.selectedPeaks.push(selectedPeak);
            this.peakAdded.emit(selectedPeak);
        }

        // Update the marker icon
        this.updateMarkerIcon(peak.id);

        // Emit the full selection
        this.selectionChange.emit([...this.selectedPeaks]);
    }

    private updateMarkerIcon(peakId: number): void {
        const marker = this.peakMarkers.get(peakId);
        const peak = this.allPeaks.find(p => p.id === peakId);

        if (marker && peak) {
            const isSelected = this.selectedPeaks.some(sp => sp.id === peakId);
            const icon = this.getIconForPeak(peak, isSelected);
            marker.setIcon(icon);

            // Update popup content
            marker.setPopupContent(this.buildPeakPopup(peak, isSelected));
        }
    }

    removeSelectedPeak(peak: SelectedPeak): void {
        const index = this.selectedPeaks.findIndex(sp => sp.id === peak.id);
        if (index >= 0) {
            this.selectedPeaks.splice(index, 1);
            this.peakRemoved.emit(peak);
            this.updateMarkerIcon(peak.id);
            this.selectionChange.emit([...this.selectedPeaks]);
        }
    }

    // Search functionality
    onSearchInput(): void {
        if (!this.searchQuery.trim()) {
            this.filteredPeaks = [];
            this.showSearchResults = false;
            return;
        }

        const query = this.searchQuery.toLowerCase();
        this.filteredPeaks = this.allPeaks
            .filter(p =>
                p.name?.toLowerCase().includes(query) ||
                p.alt_name?.toLowerCase().includes(query) ||
                p.region?.toLowerCase().includes(query)
            )
            .slice(0, 15); // Limit results

        this.showSearchResults = this.filteredPeaks.length > 0;
    }

    selectFromSearch(peak: Peak): void {
        // Zoom to peak on map at close zoom level to break clusters
        if (this.map) {
            this.map.setView([peak.latitude, peak.longitude], 14);

            // Open the popup for this peak
            const marker = this.peakMarkers.get(peak.id);
            if (marker) {
                marker.openPopup();
            }
        }

        // Clear search
        this.searchQuery = '';
        this.filteredPeaks = [];
        this.showSearchResults = false;
    }

    onSearchBlur(): void {
        // Delay hiding results to allow click to register
        setTimeout(() => {
            this.showSearchResults = false;
        }, 200);
    }

    // Click on selected peak chip to focus on map
    focusOnPeak(peak: SelectedPeak): void {
        if (this.map) {
            this.map.setView([peak.latitude, peak.longitude], 14);

            // Open the popup for this peak
            const marker = this.peakMarkers.get(peak.id);
            if (marker) {
                marker.openPopup();
            }
        }
    }

    onClose(): void {
        this.close.emit();
    }
}
