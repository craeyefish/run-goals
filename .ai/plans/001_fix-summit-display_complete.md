# Plan: Fix Summit Display on Map

## Description

Summits/peaks are not showing up on the activity map despite activities being displayed correctly. The main issue is that the `displayPeaks()` method exists but is never called, and there are some potential issues with the backend API response format.

## Scope

### In Scope

- Fix the primary issue where `displayPeaks()` is never called in the activity map component
- Ensure peaks are properly loaded and displayed with correct icons (visited vs unvisited)
- Verify backend API returns peaks with proper user context (`is_summited` field)
- Ensure proper error handling and loading states for peak display
- Test that the toggle functionality works correctly
- Verify map clustering works for both activities and peaks

### Out of Scope

- Major refactoring of summit detection algorithms
- Changes to the database schema or peak data import process
- UI/UX improvements beyond fixing the core functionality
- Performance optimizations (clustering is already implemented)

## Steps

[*] **1. Fix Frontend Peak Display Logic**

- Update `activity-map.component.ts` to call `displayPeaks()` after peaks are loaded
- Ensure peaks are displayed when the component initializes and when peaks data is received
- Fix the toggle functionality to properly call `displayPeaks()` when showing peaks

[*] **2. Verify Backend API Response Format**

- Check that `/api/peaks` endpoint returns the `PeakSummited` model with `is_summited` field
- Ensure the backend properly identifies which peaks have been summited by the current user
- Verify the API response matches the frontend `Peak` interface expectations

[*] **3. Test Peak Icon Display**

- Verify that `defaultPeakIcon` and `visitedPeakIcon` are properly imported and used
- Test that visited peaks show green icons and unvisited peaks show default icons
- Ensure summit icons are properly positioned and sized on the map

[*] **4. Fix Peak Loading and Error Handling**

- Add proper loading states for peak data
- Implement error handling for peak loading failures
- Add console logging to debug peak loading issues

[*] **5. Verify Peak Popup Information**

- Test that peak popups show correct information (name, elevation)
- Ensure popup content is properly formatted and displays correctly

[*] **6. Test Integration with Map Clustering**

- Verify that peaks are added to the correct marker cluster group
- Test that peak markers cluster properly with activities
- Ensure toggle functionality properly adds/removes peak markers from clusters

[*] **7. Update Memory Documentation**

- Document the summit display issue and resolution in `.ai/memory/summit-display-fix.md`
- Include key technical details about the peak loading flow
- Note any potential future improvements or edge cases discovered

[*] **8. Test End-to-End Functionality**

- Test that summits appear on the map when "Show Summits" is checked
- Verify that summits disappear when "Show Summits" is unchecked
- Test that both activities and summits can be displayed simultaneously
- Test peak popup functionality by clicking on summit markers

[*] **9. Verify User Context in Peak Display**

- Ensure that only the current user's summited peaks show as visited (green icons)
- Test that peaks summited by other users don't show as visited for the current user
- Verify that peak summiting status is correctly determined

[*] **10. Performance and Cleanup**

- Ensure peak markers are properly cleaned up when toggling display
- Verify that memory usage is reasonable when displaying many peaks
- Test map performance with both activities and peaks displayed

## Technical Notes

### Root Cause Analysis

1. **Primary Issue**: `displayPeaks()` method exists but is never called in the component lifecycle
2. **Secondary Issue**: The toggle functionality has `displayPeaks()` commented out
3. **Potential Issue**: Backend API might not be returning the `is_summited` field properly

### Key Files to Modify

- `frontend/strava-goal/src/app/components/activity-map/activity-map.component.ts`
- Potentially `backend/services/peakService.go` if API response format needs fixing
- `frontend/strava-goal/src/app/services/peak.service.ts` for any interface mismatches

### API Endpoints Involved

- `GET /api/peaks` - Returns list of peaks with summiting status for current user
- The response should match the `Peak` interface with `is_summited: boolean` field

### Dependencies

- Leaflet clustering (`this.markerClusterGroup`)
- Summit icons: `assets/summit-icon.png` and `assets/summit-icon-green.png`
- Peak service for data loading
