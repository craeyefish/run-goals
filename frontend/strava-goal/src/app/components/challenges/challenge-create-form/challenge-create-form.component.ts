import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Output, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
    CreateChallengeRequest,
    ChallengeType,
    CompetitionMode,
    Visibility,
    GoalType
} from 'src/app/models/challenge.model';
import { PeakPickerComponent, SelectedPeak } from 'src/app/components/peak-picker/peak-picker.component';

@Component({
    selector: 'challenge-create-form',
    standalone: true,
    imports: [CommonModule, FormsModule, PeakPickerComponent],
    templateUrl: './challenge-create-form.component.html',
    styleUrls: ['./challenge-create-form.component.scss'],
})
export class ChallengeCreateFormComponent {
    @Output() onSubmit = new EventEmitter<CreateChallengeRequest>();
    @Output() onCancel = new EventEmitter<void>();

    // Track if mousedown started on overlay
    private mouseDownOnOverlay = false;

    // Form state
    currentStep = signal<number>(1);
    totalSteps = 3;

    // Step 1: Basic Info
    name = '';
    description = '';
    goalType: GoalType = 'specific_summits';
    competitionMode: CompetitionMode = 'collaborative';
    region = '';
    difficulty = '';

    // Step 2: Dates
    hasDeadline = false;
    startDate = '';
    deadline = '';

    // Step 3: Goal-specific fields
    // For specific_summits
    showPeakPicker = signal(false);
    selectedPeaks: SelectedPeak[] = [];

    // For distance (in km, will convert to meters)
    targetDistance = 0;

    // For elevation (in meters)
    targetElevation = 0;

    // For summit_count
    targetSummitCount = 0;

    // Validation
    get step1Valid(): boolean {
        return this.name.trim().length >= 3;
    }

    get step2Valid(): boolean {
        if (!this.hasDeadline) return true;
        if (!this.deadline) return false;
        if (this.startDate && this.deadline) {
            return new Date(this.startDate) <= new Date(this.deadline);
        }
        return true;
    }

    get step3Valid(): boolean {
        switch (this.goalType) {
            case 'specific_summits':
                return this.selectedPeaks.length > 0;
            case 'distance':
                return this.targetDistance > 0;
            case 'elevation':
                return this.targetElevation > 0;
            case 'summit_count':
                return this.targetSummitCount > 0;
            default:
                return false;
        }
    }

    get canSubmit(): boolean {
        return this.step1Valid && this.step2Valid && this.step3Valid;
    }

    get selectedPeakIds(): number[] {
        return this.selectedPeaks.map(p => p.id);
    }

    // Navigation
    nextStep() {
        if (this.currentStep() < this.totalSteps) {
            this.currentStep.update(s => s + 1);
        }
    }

    prevStep() {
        if (this.currentStep() > 1) {
            this.currentStep.update(s => s - 1);
        }
    }

    goToStep(step: number) {
        if (step >= 1 && step <= this.totalSteps) {
            this.currentStep.set(step);
        }
    }

    // Peak selection
    openPeakPicker() {
        this.showPeakPicker.set(true);
    }

    closePeakPicker() {
        this.showPeakPicker.set(false);
    }

    onPeakSelectionChange(peaks: SelectedPeak[]) {
        this.selectedPeaks = peaks;
    }

    onPeakAdded(peak: SelectedPeak) {
        if (!this.selectedPeaks.find(p => p.id === peak.id)) {
            this.selectedPeaks = [...this.selectedPeaks, peak];
        }
    }

    onPeakRemoved(peak: SelectedPeak) {
        this.selectedPeaks = this.selectedPeaks.filter(p => p.id !== peak.id);
    }

    removePeak(peak: SelectedPeak) {
        this.selectedPeaks = this.selectedPeaks.filter(p => p.id !== peak.id);
    }

    // Modal close handling - only close if click started AND ended on overlay
    onOverlayMouseDown(event: MouseEvent) {
        // Only set flag if mousedown is directly on overlay (not bubbled from content)
        if (event.target === event.currentTarget) {
            this.mouseDownOnOverlay = true;
        }
    }

    onOverlayMouseUp(event: MouseEvent) {
        // Only close if both mousedown and mouseup were on overlay
        if (this.mouseDownOnOverlay && event.target === event.currentTarget) {
            this.cancel();
        }
        this.mouseDownOnOverlay = false;
    }

    // Form submission
    submit() {
        if (!this.canSubmit) return;

        // Convert date strings to ISO format for backend
        const formatDateToISO = (dateStr: string | undefined): string | undefined => {
            if (!dateStr) return undefined;
            // HTML date input gives YYYY-MM-DD, convert to ISO timestamp
            return new Date(dateStr + 'T00:00:00Z').toISOString();
        };

        const request: CreateChallengeRequest = {
            name: this.name.trim(),
            description: this.description.trim() || undefined,
            challengeType: 'custom', // Always custom now
            goalType: this.goalType,
            competitionMode: this.competitionMode,
            visibility: 'private', // Always private (only admins can make public)
            startDate: formatDateToISO(this.startDate),
            deadline: this.hasDeadline && this.deadline ? formatDateToISO(this.deadline) : undefined,
            region: this.region.trim() || undefined,
            difficulty: this.difficulty || undefined,
            // Goal-specific fields
            targetValue: this.goalType === 'distance' ? this.targetDistance * 1000 : // Convert km to meters
                         this.goalType === 'elevation' ? this.targetElevation :
                         undefined,
            targetSummitCount: this.goalType === 'summit_count' ? this.targetSummitCount : undefined,
            peakIds: this.goalType === 'specific_summits' ? this.selectedPeaks.map(p => p.id) : [],
        };

        this.onSubmit.emit(request);
    }

    cancel() {
        this.onCancel.emit();
    }

    // Helper methods
    getTotalElevation(): number {
        return this.selectedPeaks.reduce((sum, p) => sum + p.elevation_meters, 0);
    }
}
