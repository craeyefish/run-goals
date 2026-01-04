# UI Refactor Implementation Status

## ✅ Completed

### Phase 1: DataTable Component
- ✅ Created reusable `DataTableComponent` at `frontend/strava-goal/src/app/components/shared/data-table/`
- ✅ Supports all column types: text, number, date, link, badge, progress, custom
- ✅ Sortable columns with visual indicators
- ✅ Row click handlers
- ✅ Responsive design
- ✅ Empty states

### Phase 2: Home Page Refactor
- ✅ Removed Active Challenges section
- ✅ Made Summit Wishlist collapsible (starts collapsed, saves to localStorage)
- ✅ Updated layout and styling

### Phase 3: Groups Page Refactor (IN PROGRESS)
- ✅ Removed Goals section HTML
- ✅ Updated Challenges section to use DataTable in HTML
- ✅ Updated Members section to use DataTable in HTML
- ✅ Updated modal text ("Add Challenge to Group" instead of "Adopt")
- ⏳ **NEEDS**: TypeScript updates for group-details.component.ts:
  - Import DataTableComponent
  - Define `challengeColumns` for table
  - Define `memberColumns` for table
  - Add `memberStats` signal
  - Add `currentYear` property
  - Remove goals-related imports and methods
  - Load member stats data

### Phase 4: Challenge Detail Page (NOT STARTED)
- ⏳ Convert Peaks tab to DataTable
- ⏳ Convert Participants tab to DataTable
- ⏳ Convert Activities tab to DataTable

---

## Next Steps for Groups Page TypeScript

The HTML is updated but TypeScript needs these changes in `group-details.component.ts`:

```typescript
// 1. Add imports
import { DataTableComponent, TableColumn } from 'src/app/components/shared/data-table/data-table.component';

// 2. Update imports array
imports: [
  CommonModule,
  FormsModule,
  DataTableComponent, // ADD THIS
  // Remove: GoalsCreateFormComponent, GoalsEditFormComponent, GoalDeleteConfirmationComponent, GroupsGoalsTableComponent, GroupsMembersTableComponent
],

// 3. Remove goal-related properties and methods (openCreateGoalForm, openEditGoalForm, etc.)

// 4. Add new properties
currentYear = new Date().getFullYear();
memberStats = signal<any[]>([]);

// 5. Define challenge table columns
challengeColumns: TableColumn[] = [
  { header: 'Challenge Name', field: 'name', type: 'text', sortable: true },
  { header: 'Type', field: 'goalType', type: 'text',
    formatter: (value) => this.getGoalTypeLabel(value) },
  { header: 'Mode', field: 'competitionMode', type: 'badge',
    badgeClass: (value) => value === 'collaborative' ? 'badge-success' : 'badge-info',
    formatter: (value) => value === 'collaborative' ? '🤝 Collaborative' : '🏅 Competitive' },
  { header: 'Region', field: 'region', type: 'text' },
  { header: 'Deadline', field: 'deadline', type: 'date' },
];

// 6. Define member table columns
memberColumns: TableColumn[] = [
  { header: 'Member', field: 'userName', type: 'link', sortable: true,
    linkFn: (row) => `https://www.strava.com/athletes/${row.stravaAthleteId}`,
    linkExternal: true },
  { header: 'Distance (km)', field: 'totalDistance', type: 'number', sortable: true,
    formatter: (value) => (value / 1000).toFixed(1), align: 'right' },
  { header: 'Elevation (m)', field: 'totalElevation', type: 'number', sortable: true,
    formatter: (value) => Math.round(value).toLocaleString(), align: 'right' },
  { header: 'Summits', field: 'totalSummits', type: 'number', sortable: true, align: 'right' },
];

// 7. Load member stats in ngOnInit
// Need backend endpoint: GET /api/group-member-stats?groupId={id}&year={year}
```

---

## Backend Needed for Groups Page

Create endpoint to get member stats:
```go
// In backend/handlers/ApiHandler.go
func (h *ApiHandler) GetGroupMemberStats(w http.ResponseWriter, r *http.Request) {
    groupID := parseGroupID(r)
    year := parseYear(r) // default to current year

    // Get all group members
    // For each member, calculate YTD stats:
    //   - Total distance
    //   - Total elevation
    //   - Total summits

    // Return array of member stats
}
```

---

## Challenge Detail Page - Table Conversions Needed

### Peaks Tab
```typescript
peakColumns: TableColumn[] = [
  { header: 'Peak Name', field: 'name', type: 'link', sortable: true,
    linkFn: (row) => `/explore?peakId=${row.peakId}` },
  { header: 'Elevation', field: 'elevation', type: 'number', sortable: true,
    formatter: (value) => `${value}m`, align: 'right' },
  { header: 'Region', field: 'region', type: 'text', sortable: true },
  { header: 'Status', field: 'isSummited', type: 'badge',
    badgeClass: (value) => value ? 'badge-success' : 'badge-default',
    formatter: (value) => value ? '✓ Summited' : 'Pending' },
];
```

### Participants Tab (Collaborative)
```typescript
collaborativeParticipantsColumns: TableColumn[] = [
  { header: 'Member', field: 'userName', type: 'link', sortable: true,
    linkFn: (row) => `https://www.strava.com/athletes/${row.stravaAthleteId}`,
    linkExternal: true,
    formatter: (value, row) => value || row.stravaAthleteId },
  { header: 'Contributed', field: 'totalDistance', type: 'number', sortable: true,
    formatter: (value, row) => {
      if (challenge.goalType === 'distance') return `${(value/1000).toFixed(1)} km`;
      if (challenge.goalType === 'elevation') return `${Math.round(row.totalElevation)} m`;
      if (challenge.goalType === 'summit_count') return `${row.totalSummitCount} summits`;
      return `${row.peaksCompleted} peaks`;
    }},
  { header: 'Joined', field: 'joinedAt', type: 'date' },
];
```

### Activities Tab
```typescript
activityColumns: TableColumn[] = [
  { header: 'Activity', field: 'name', type: 'link', sortable: false,
    linkFn: (row) => `https://www.strava.com/activities/${row.strava_activity_id}`,
    linkExternal: true },
  { header: 'User', field: 'userName', type: 'link',
    linkFn: (row) => `https://www.strava.com/athletes/${row.stravaAthleteId}`,
    linkExternal: true,
    formatter: (value, row) => value || row.stravaAthleteId },
  // Conditional columns based on goalType
  // + Peaks column for summit challenges (shows peakNames)
  // + Distance for distance challenges
  // + Elevation for elevation challenges
  { header: 'Date', field: 'start_date', type: 'date', sortable: true },
];
```

---

## Files Modified So Far

### Frontend
- ✅ `frontend/strava-goal/src/app/components/shared/data-table/` (NEW)
- ✅ `frontend/strava-goal/src/app/pages/home/home.component.{html,ts,scss}`
- ✅ `frontend/strava-goal/src/app/pages/groups/group-details/group-details.component.html`
- ⏳ `frontend/strava-goal/src/app/pages/groups/group-details/group-details.component.ts` (NEEDS UPDATE)
- ⏳ `frontend/strava-goal/src/app/pages/challenges/challenge-detail/challenge-detail.component.{html,ts}` (NOT STARTED)

### Backend
- No changes needed yet for Phase 1-2
- ⏳ Need group member stats endpoint for Phase 3

---

## Testing Checklist

### Home Page
- [ ] Summit wishlist starts collapsed
- [ ] Clicking header expands/collapses wishlist
- [ ] Preference saves to localStorage
- [ ] Action buttons don't trigger collapse
- [ ] Active challenges section is gone

### Groups Page
- [ ] Goals section completely removed
- [ ] Challenges display in table format
- [ ] Challenge table is sortable
- [ ] Clicking challenge row navigates to detail
- [ ] "Add challenge to group" button works
- [ ] Modal shows updated text
- [ ] Members show year-to-date stats
- [ ] Member table is sortable
- [ ] Strava profile links work

### Challenge Detail Page (Not Yet Implemented)
- [ ] Peaks tab uses table
- [ ] Participants tab uses table
- [ ] Activities tab uses table
- [ ] All tables are sortable
- [ ] All Strava links work
- [ ] Peak names show in activities for summit challenges
