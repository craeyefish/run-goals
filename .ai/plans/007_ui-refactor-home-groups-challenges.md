# UI Refactor - Home, Groups, and Challenges Pages

## Status: Planning

## Overview
Major UI refactoring to declutter pages, improve consistency, and create reusable table components across the application.

---

## Phase 1: Reusable Table Component

### 1.1 Create Generic Data Table Component
- [x] Create `data-table.component.ts` in `src/app/components/shared/`
- [x] Define generic table interface that supports:
  - Column definitions (header, field, formatter, sortable)
  - Row data (generic type)
  - Row click handlers
  - Custom cell renderers (for links, badges, etc.)
  - Sorting state
  - Empty state message
- [x] Implement responsive table styles
- [x] Add sorting functionality (click headers to sort)
- [x] Support for action columns (buttons, links)

### 1.2 Table Column Types to Support
- [x] Text columns (plain text)
- [x] Number columns (with formatting: distance, elevation, percentages)
- [x] Date columns (formatted dates)
- [x] Link columns (internal router links and external links)
- [x] Badge columns (status badges, labels)
- [x] Progress columns (progress bars with percentage)
- [x] Action columns (buttons, icons)

---

## Phase 2: Home Page Refactor

### 2.1 Remove Active Challenges Section
- [ ] Remove active challenges display from home page
- [ ] Update layout to use freed space
- [ ] Remove related state/signals if no longer needed

### 2.2 Make Summit Wishlist Collapsible
- [ ] Add collapsible/expandable state (signal)
- [ ] Default to collapsed on page load
- [ ] Add expand/collapse toggle button with icon
- [ ] Animate expansion/collapse
- [ ] Update styles for collapsed vs expanded state
- [ ] Save preference to localStorage (optional)

### 2.3 Home Page Layout Cleanup
- [ ] Reorganize remaining sections
- [ ] Update spacing and styling
- [ ] Ensure responsive design still works
- [ ] Test on mobile/tablet

---

## Phase 3: Groups Page Refactor

### 3.1 Remove Goals Section
- [ ] Remove "Goals" tab/section completely
- [ ] Remove goal-related state management from groups page
- [ ] Remove goal-related API calls
- [ ] Update navigation/tabs to only show relevant sections

### 3.2 Update Members Section
**Decision Point**: Choose one approach:
- Option A: Show total contributions since group creation
- Option B: Show totals for current calendar year (RECOMMENDED)
- Option C: Remove members section entirely (show only in challenges)

**Implementation (Option B - Calendar Year Totals)**:
- [ ] Update backend to calculate member stats for current year
  - [ ] Add `GetGroupMemberYearStats` endpoint
  - [ ] Query activities for each member (year-to-date)
  - [ ] Return: user info, distance, elevation, summit count
- [ ] Create table using new DataTable component
- [ ] Columns: Member Name (link to Strava), Distance, Elevation, Summits
- [ ] Sort by distance by default

### 3.3 Update Challenges Section to Use Table
- [ ] Convert challenges display to use DataTable component
- [ ] Columns: Challenge Name, Type, Goal, Progress, Participants, Deadline
- [ ] Click row to navigate to challenge detail
- [ ] Show appropriate metrics based on goal type
- [ ] Add visual indicators (icons, progress bars)

### 3.4 Clarify "Adopt Challenge" Feature
**Decision Point**: What should adoption mean?
- Option A: Auto-join all group members to challenge
- Option B: Make challenge visible to group (current behavior)
- Option C: Remove adoption, use "Add to Group" instead

**Recommended Approach (Option C)**:
- [ ] Rename "Adopt" to "Add Challenge to Group"
- [ ] Update UI labels and tooltips
- [ ] Add help text explaining what this does
- [ ] Consider adding "Auto-join all members" as separate action

---

## Phase 4: Challenge Detail Page Refactor

### 4.1 Convert Peaks Tab to Table
- [ ] Use DataTable component for peaks list
- [ ] Columns: Peak Name, Elevation, Region, Status (Completed/Pending)
- [ ] Click row to navigate to explore page with peak filter
- [ ] Show completion badges/icons
- [ ] Sort by: name, elevation, status

### 4.2 Convert Participants Tab to Table
- [ ] Use DataTable component for participants
- [ ] Columns vary by challenge type:
  - **Collaborative Distance**: Name (Strava link), Distance Contributed, Joined Date
  - **Collaborative Elevation**: Name (Strava link), Elevation Contributed, Joined Date
  - **Competitive Distance**: Rank, Name (Strava link), Distance, Progress %
  - **Competitive Elevation**: Rank, Name (Strava link), Elevation, Progress %
  - **Summit Count**: Rank, Name (Strava link), Summits, Progress %
  - **Specific Summits**: Rank, Name (Strava link), Peaks Completed, Progress %
- [ ] Show username with fallback to Strava ID
- [ ] Link names to Strava profiles
- [ ] Sort functionality

### 4.3 Convert Activities Tab to Table
- [ ] Use DataTable component for activities
- [ ] Columns vary by challenge type:
  - **Distance/Elevation**: Activity Name (Strava link), User (Strava link), Distance/Elevation, Date
  - **Summit Challenges**: Activity Name (Strava link), User (Strava link), Peaks, Distance, Elevation, Date
- [ ] Click activity name to open Strava activity
- [ ] Click user name to open Strava profile
- [ ] Format dates consistently
- [ ] Show peak names as badges for summit activities

### 4.4 Update Challenge Detail Layout
- [ ] Ensure header/progress section remains prominent
- [ ] Update tab content area to accommodate tables
- [ ] Ensure responsive design for tables
- [ ] Test on mobile (tables should scroll horizontally if needed)

---

## Phase 5: Testing & Polish

### 5.1 Cross-Browser Testing
- [ ] Test on Chrome
- [ ] Test on Firefox
- [ ] Test on Safari
- [ ] Test on mobile browsers

### 5.2 Responsive Design
- [ ] Test all pages on mobile
- [ ] Test all pages on tablet
- [ ] Ensure tables scroll/stack appropriately
- [ ] Test collapsible sections on mobile

### 5.3 Accessibility
- [ ] Ensure keyboard navigation works
- [ ] Add ARIA labels where needed
- [ ] Test with screen reader (basic check)
- [ ] Ensure color contrast is sufficient

### 5.4 Performance
- [ ] Check load times for large tables
- [ ] Implement virtual scrolling if needed (100+ rows)
- [ ] Optimize re-renders

---

## Implementation Order

1. **Start with DataTable Component** (Phase 1)
   - This is the foundation for everything else
   - Can be developed and tested in isolation

2. **Home Page** (Phase 2)
   - Quick wins, simple changes
   - Tests the collapsible pattern

3. **Groups Page** (Phase 3)
   - First real use of DataTable component
   - Tests table component with real data

4. **Challenge Detail Page** (Phase 4)
   - Most complex, uses table in multiple contexts
   - Benefits from lessons learned in groups page

5. **Testing & Polish** (Phase 5)
   - Final pass on everything

---

## Design Decisions Needed

1. **Members Section**: Show year-to-date stats or remove entirely?
2. **Adopt Challenge**: Keep current behavior or change to auto-join?
3. **Table Styling**: Match existing card design or create new table aesthetic?
4. **Mobile Tables**: Horizontal scroll or responsive card layout?

---

## Notes

- DataTable component should be highly reusable
- Keep consistent styling across all tables
- Consider adding export functionality later (CSV)
- May want to add pagination for very large datasets
- Keep accessibility in mind throughout
