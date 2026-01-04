import { CommonModule } from '@angular/common';
import { Component, signal, WritableSignal, OnInit, computed } from '@angular/core';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { BreadcrumbService } from 'src/app/services/breadcrumb.service';
import { GroupService } from 'src/app/services/groups.service';
import { ChallengeService } from 'src/app/services/challenge.service';
import { Challenge, ChallengeWithProgress } from 'src/app/models/challenge.model';
import { DataTableComponent, TableColumn } from 'src/app/components/shared/data-table/data-table.component';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'group-details-page',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    DataTableComponent,
  ],
  templateUrl: './group-details.component.html',
  styleUrls: ['./group-details.component.scss'],
})
export class GroupsDetailsPageComponent implements OnInit {
  private destroy$ = new Subject<void>();

  constructor(
    private groupService: GroupService,
    private challengeService: ChallengeService,
    private breadcrumbService: BreadcrumbService,
    private router: Router
  ) { }

  // Challenge-related signals
  groupChallenges = signal<Challenge[]>([]);
  availableChallenges = signal<Challenge[]>([]);
  filteredAvailableChallenges = signal<Challenge[]>([]);
  showAdoptChallengeModal = signal(false);
  loadingChallenges = signal(false);
  challengeSearchQuery = '';

  // Member stats
  memberStats = signal<any[]>([]);
  currentYear = new Date().getFullYear();

  selectedGroup = this.groupService.selectedGroup;

  // Table column definitions
  challengeColumns: TableColumn[] = [
    { header: 'Challenge Name', field: 'name', type: 'text', sortable: true },
    {
      header: 'Type',
      field: 'goalType',
      type: 'text',
      formatter: (value) => this.getGoalTypeLabel(value)
    },
    {
      header: 'Goal',
      field: 'targetValue',
      type: 'text',
      formatter: (value, row) => this.formatChallengeGoal(row)
    },
    {
      header: 'Mode',
      field: 'competitionMode',
      type: 'badge',
      badgeClass: (value) => value === 'collaborative' ? 'badge-success' : 'badge-info',
      formatter: (value) => value === 'collaborative' ? '🤝 Collaborative' : '🏅 Competitive'
    },
    { header: 'Region', field: 'region', type: 'text' },
    { header: 'Deadline', field: 'deadline', type: 'date' },
  ];

  memberColumns: TableColumn[] = [
    {
      header: 'Member',
      field: 'userName',
      type: 'link',
      sortable: true,
      linkFn: (row) => `https://www.strava.com/athletes/${row.stravaAthleteId}`,
      linkExternal: true
    },
    {
      header: 'Distance (km)',
      field: 'totalDistance',
      type: 'number',
      sortable: true,
      formatter: (value) => (value / 1000).toFixed(1),
      align: 'right'
    },
    {
      header: 'Elevation (m)',
      field: 'totalElevation',
      type: 'number',
      sortable: true,
      formatter: (value) => Math.round(value).toLocaleString(),
      align: 'right'
    },
    {
      header: 'Summits',
      field: 'totalSummits',
      type: 'number',
      sortable: true,
      align: 'right'
    },
  ];

  ngOnInit() {
    const selectedGroup = this.groupService.selectedGroup();

    // Check if we have a selected group
    if (!selectedGroup) {
      console.warn('No group selected, redirecting to groups list');
      // Redirect to groups list if no group is selected
      this.router.navigate(['/groups']);
      return;
    }

    // Set breadcrumb with the group name
    this.breadcrumbService.addItem({
      label: selectedGroup.name,
    });

    // Load group challenges
    this.loadGroupChallenges(selectedGroup.id);

    // Load member stats
    this.loadMemberStats(selectedGroup.id);
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  // ==================== Challenge Methods ====================

  loadGroupChallenges(groupId: number): void {
    this.challengeService.getGroupChallenges(groupId).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (response) => {
        this.groupChallenges.set(response.challenges || []);
      },
      error: (err) => console.error('Error loading group challenges:', err)
    });
  }

  openAdoptChallengeModal(): void {
    this.showAdoptChallengeModal.set(true);
    this.loadAvailableChallenges();
  }

  closeAdoptChallengeModal(): void {
    this.showAdoptChallengeModal.set(false);
    this.challengeSearchQuery = '';
  }

  loadAvailableChallenges(): void {
    this.loadingChallenges.set(true);
    this.challengeService.getPublicChallenges().pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (response) => {
        this.availableChallenges.set(response.challenges || []);
        this.filteredAvailableChallenges.set(response.challenges || []);
        this.loadingChallenges.set(false);
      },
      error: (err) => {
        console.error('Error loading available challenges:', err);
        this.loadingChallenges.set(false);
      }
    });
  }

  filterAvailableChallenges(): void {
    const query = this.challengeSearchQuery.toLowerCase().trim();
    if (!query) {
      this.filteredAvailableChallenges.set(this.availableChallenges());
      return;
    }

    const filtered = this.availableChallenges().filter(c =>
      c.name.toLowerCase().includes(query) ||
      c.region?.toLowerCase().includes(query) ||
      c.description?.toLowerCase().includes(query)
    );
    this.filteredAvailableChallenges.set(filtered);
  }

  isChallengeAdopted(challengeId: number): boolean {
    return this.groupChallenges().some(c => c.id === challengeId);
  }

  adoptChallenge(challenge: Challenge): void {
    if (this.isChallengeAdopted(challenge.id)) return;

    const selectedGroup = this.selectedGroup();
    if (!selectedGroup) return;

    this.challengeService.addGroupToChallenge(challenge.id, { groupId: selectedGroup.id }).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: () => {
        console.log('Challenge adopted:', challenge.name);
        this.loadGroupChallenges(selectedGroup.id);
        this.closeAdoptChallengeModal();
      },
      error: (err) => {
        console.error('Error adopting challenge:', err);
        alert('Failed to adopt challenge. Please try again.');
      }
    });
  }

  removeGroupChallenge(challengeId: number, event: Event): void {
    event.stopPropagation();

    const selectedGroup = this.selectedGroup();
    if (!selectedGroup) return;

    if (!confirm('Remove this challenge from the group?')) return;

    this.challengeService.removeGroupFromChallenge(challengeId, selectedGroup.id).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: () => {
        console.log('Challenge removed from group');
        this.loadGroupChallenges(selectedGroup.id);
      },
      error: (err) => {
        console.error('Error removing challenge:', err);
        alert('Failed to remove challenge. Please try again.');
      }
    });
  }

  navigateToChallenge(challengeId: number): void {
    this.router.navigate(['/challenges', challengeId]);
  }

  // ==================== Member Stats Methods ====================

  loadMemberStats(groupId: number): void {
    this.groupService.getGroupMembers(groupId).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (response) => {
        // For now, just set empty stats - will need backend endpoint for actual stats
        // TODO: Replace with actual API call for year-to-date member stats
        // Note: Member interface doesn't have strava_athlete_id, using user_id as placeholder
        const stats = response.members.map(member => ({
          userName: member.username || `User ${member.user_id}`,
          stravaAthleteId: member.user_id, // Using user_id as placeholder until backend provides strava_athlete_id
          totalDistance: 0,
          totalElevation: 0,
          totalSummits: 0
        }));
        this.memberStats.set(stats);
      },
      error: (err) => console.error('Error loading member stats:', err)
    });
  }

  // ==================== Helper Methods ====================

  getGoalTypeLabel(goalType: string): string {
    const labels: Record<string, string> = {
      'distance': 'Distance',
      'elevation': 'Elevation Gain',
      'summit_count': 'Summit Count',
      'specific_summits': 'Specific Summits'
    };
    return labels[goalType] || goalType;
  }

  formatChallengeGoal(challenge: Challenge): string {
    if (challenge.goalType === 'distance' && challenge.targetValue) {
      return `${(challenge.targetValue / 1000).toFixed(0)} km`;
    }
    if (challenge.goalType === 'elevation' && challenge.targetValue) {
      return `${challenge.targetValue.toLocaleString()} m`;
    }
    if (challenge.goalType === 'summit_count' && challenge.targetSummitCount) {
      return `${challenge.targetSummitCount} summits`;
    }
    if (challenge.goalType === 'specific_summits') {
      return 'Specific Peaks';
    }
    return '-';
  }
}
