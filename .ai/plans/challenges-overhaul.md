# Challenges System Overhaul - Implementation Plan

**Date:** 2026-01-02
**Goal:** Transform challenges into the core feature of Summit Seekers with intuitive UX, comprehensive goal types, and competitive/collaborative modes.

---

## Executive Summary

This plan restructures the challenges system to support four goal types (distance, elevation, summit_count, specific_summits) with both competitive and collaborative modes. We'll merge the existing Groups Goals functionality into Challenges, add challenge editing, implement a join code system, create three tabs (My Challenges, My Created, Discover), and build a retro RPG-style info/tutorial popup.

---

## Current State Analysis

### What Works Well
✅ **Backend Infrastructure:**
- Solid DAO/Service/Controller architecture
- Auto-detection of summits from activities
- Progress caching system
- Leaderboard calculations

✅ **Groups Goals System:**
- Clean table UI with expandable rows
- Four goal types already implemented
- Progress calculation for all types
- Good UX for creating/editing goals

✅ **Challenges System:**
- Participant tracking
- Summit logging
- Featured/public challenge discovery
- Group adoption feature

### Current Limitations
❌ Only supports `specific_summits` challenge type
❌ No challenge editing capability
❌ No join codes for finding/joining challenges
❌ Public vs. discoverable challenges conflated
❌ UI inconsistency (challenges vs groups pages)
❌ Duplicate functionality (challenges vs group goals)
❌ No clear product narrative

---

## Design Decisions

### 1. Challenge Types & Modes Matrix

| Type | Competitive Mode | Collaborative Mode |
|------|------------------|-------------------|
| **Distance** | Leaderboard by total km | Team total km target |
| **Elevation** | Leaderboard by total m | Team total m target |
| **Summit Count** | Leaderboard by # summits | Team total summits target |
| **Specific Summits** | Leaderboard by completion % | Team completes all peaks |

### 2. Challenge Visibility Model

**OLD (Current):**
- `visibility`: private/friends/public
- `is_featured`: boolean

**NEW (Proposed):**
- **Private Challenges:** Only creator and invited users (requires join code)
- **Discoverable Challenges:** Anyone can find via join code (no password required)
- **Public Challenges:** Admin-curated, featured on Discover tab (requires admin approval)

### 3. Three-Tab Structure

**Tab 1: My Challenges**
- Challenges I've joined (not created by me)
- Shows my progress, rank, and leaderboard
- Empty state: "Join a challenge to get started!"

**Tab 2: My Created**
- Challenges I've created
- Edit/lock/delete capabilities
- Empty state: "Create your first challenge!"

**Tab 3: Discover**
- **Public Challenges Section:** Admin-curated featured challenges
- **Join by Code Section:** Text input to join by code
- **Propose Challenge Section:** Submit challenge idea for admin review

### 4. Join Code System

- **Format:** 6-character alphanumeric (e.g., `CHX7K2`)
- **Generation:** Auto-generated on challenge creation
- **Uniqueness:** Database constraint + retry logic
- **Display:** Shown prominently on challenge detail page for sharing

### 5. Challenge Locking

- **Unlocked (default):** Creator can edit name, description, dates, peaks, target values
- **Locked:** Immutable, shows 🔒 badge
- **Lock Action:** Irreversible, requires confirmation
- **Purpose:** Optional "hardcore mode" for committed participants

### 6. Group Adoption Clarification

**Current Behavior:**
- Groups can "adopt" a challenge via `challenge_groups` table
- Groups can set custom deadline override
- Purpose: Groups can do challenges as a team activity

**Proposed Enhancement:**
- When a group adopts a challenge, all group members auto-join
- Group leaderboard tab shows group members' individual progress
- Group adoption visible on challenge detail page

---

## Database Schema Changes

### 1. Add Join Code to `challenges` Table

```sql
-- Migration: 88_add_challenge_join_code.sql
ALTER TABLE challenges ADD COLUMN IF NOT EXISTS join_code VARCHAR(6) UNIQUE;
ALTER TABLE challenges ADD COLUMN IF NOT EXISTS is_locked BOOLEAN DEFAULT FALSE;

-- Backfill join codes for existing challenges
UPDATE challenges SET join_code = UPPER(SUBSTRING(MD5(RANDOM()::TEXT || id::TEXT) FROM 1 FOR 6)) WHERE join_code IS NULL;

-- Make join_code NOT NULL after backfill
ALTER TABLE challenges ALTER COLUMN join_code SET NOT NULL;

-- Create index for fast lookup
CREATE INDEX IF NOT EXISTS idx_challenges_join_code ON challenges(join_code);
```

### 2. Extend `challenges` Table for Goal Types

```sql
-- Migration: 89_extend_challenges_for_goal_types.sql
-- Already has: challenge_type, competition_mode, target_count

-- Add goal_type to distinguish between distance/elevation/summit_count/specific_summits
ALTER TABLE challenges ADD COLUMN IF NOT EXISTS goal_type VARCHAR(20) DEFAULT 'specific_summits';
ALTER TABLE challenges ADD CONSTRAINT check_goal_type CHECK (goal_type IN ('distance', 'elevation', 'summit_count', 'specific_summits'));

-- Add target_value for distance/elevation goals (in meters)
ALTER TABLE challenges ADD COLUMN IF NOT EXISTS target_value NUMERIC;

-- Rename target_count to target_summit_count for clarity
ALTER TABLE challenges RENAME COLUMN target_count TO target_summit_count;

-- target_summit_count: for summit_count goals
-- target_value: for distance/elevation goals
-- challenge_peaks: for specific_summits goals
```

### 3. Extend `challenge_participants` for Progress Tracking

```sql
-- Migration: 90_extend_participants_progress.sql
-- Add progress fields for different goal types
ALTER TABLE challenge_participants ADD COLUMN IF NOT EXISTS total_distance NUMERIC DEFAULT 0;
ALTER TABLE challenge_participants ADD COLUMN IF NOT EXISTS total_elevation NUMERIC DEFAULT 0;
ALTER TABLE challenge_participants ADD COLUMN IF NOT EXISTS total_summit_count INTEGER DEFAULT 0;

-- Keep existing: peaks_completed (for specific_summits)
-- Keep existing: total_peaks (for specific_summits)
```

### 4. Add Challenge Proposals Table

```sql
-- Migration: 91_create_challenge_proposals.sql
CREATE TABLE IF NOT EXISTS challenge_proposals (
    id BIGSERIAL PRIMARY KEY,
    proposed_by_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    goal_type VARCHAR(20) NOT NULL,
    competition_mode VARCHAR(20) NOT NULL,
    target_value NUMERIC,
    target_summit_count INTEGER,
    region VARCHAR(100),
    difficulty VARCHAR(20),
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected
    admin_notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMPTZ,
    reviewed_by_user_id BIGINT REFERENCES users(id),

    CONSTRAINT check_proposal_status CHECK (status IN ('pending', 'approved', 'rejected'))
);

CREATE INDEX idx_proposals_status ON challenge_proposals(status);
CREATE INDEX idx_proposals_user ON challenge_proposals(proposed_by_user_id);
```

### 5. Migrate Group Goals to Challenges

```sql
-- Migration: 92_migrate_group_goals_to_challenges.sql
-- Create challenges from existing group_goals
INSERT INTO challenges (
    name,
    description,
    challenge_type,
    goal_type,
    competition_mode,
    visibility,
    start_date,
    deadline,
    created_by_group_id,
    target_value,
    target_summit_count,
    region,
    difficulty,
    is_featured,
    join_code
)
SELECT
    gg.name,
    gg.description,
    'custom',
    gg.goal_type,
    'collaborative', -- group goals are collaborative
    'private', -- existing group goals are private to group
    gg.start_date,
    gg.end_date,
    gg.group_id,
    CASE WHEN gg.goal_type IN ('distance', 'elevation') THEN gg.target_value ELSE NULL END,
    CASE WHEN gg.goal_type = 'summit_count' THEN gg.target_value::INTEGER ELSE NULL END,
    NULL, -- no region in group_goals
    'medium', -- default difficulty
    FALSE,
    UPPER(SUBSTRING(MD5(RANDOM()::TEXT || gg.id::TEXT) FROM 1 FOR 6))
FROM group_goals gg;

-- Create challenge_peaks for specific_summits goals
INSERT INTO challenge_peaks (challenge_id, peak_id, sort_order)
SELECT
    c.id,
    unnest(gg.target_summits),
    generate_series(1, array_length(gg.target_summits, 1))
FROM group_goals gg
JOIN challenges c ON c.name = gg.name AND c.created_by_group_id = gg.group_id
WHERE gg.goal_type = 'specific_summits';

-- Auto-join all group members to migrated challenges
INSERT INTO challenge_participants (challenge_id, user_id, total_peaks)
SELECT
    c.id,
    gm.user_id,
    COALESCE((SELECT COUNT(*) FROM challenge_peaks cp WHERE cp.challenge_id = c.id), 0)
FROM challenges c
JOIN group_members gm ON gm.group_id = c.created_by_group_id
WHERE c.created_by_group_id IS NOT NULL
ON CONFLICT (challenge_id, user_id) DO NOTHING;

-- Link challenges to groups via challenge_groups
INSERT INTO challenge_groups (challenge_id, group_id, started_at)
SELECT c.id, c.created_by_group_id, c.created_at
FROM challenges c
WHERE c.created_by_group_id IS NOT NULL
ON CONFLICT (challenge_id, group_id) DO NOTHING;

-- Mark group_goals table for eventual deprecation (keep for rollback safety)
-- ALTER TABLE group_goals RENAME TO group_goals_deprecated;
```

---

## Backend Implementation Tasks

### Phase 1: Database & Models (2-3 hours)

**1.1 Create Database Migrations**
- [ ] `88_add_challenge_join_code.sql` - Add join_code and is_locked
- [ ] `89_extend_challenges_for_goal_types.sql` - Add goal_type, target_value
- [ ] `90_extend_participants_progress.sql` - Add distance/elevation/summit_count tracking
- [ ] `91_create_challenge_proposals.sql` - Challenge proposal system
- [ ] `92_migrate_group_goals_to_challenges.sql` - Migrate existing data

**1.2 Update Go Models**
File: `backend/models/challenge.go`
- [ ] Add `JoinCode` field
- [ ] Add `IsLocked` field
- [ ] Add `GoalType` field (distance/elevation/summit_count/specific_summits)
- [ ] Add `TargetValue` field (for distance/elevation)
- [ ] Rename `TargetCount` → `TargetSummitCount`
- [ ] Update validation constants

**1.3 Update ChallengeParticipant Model**
- [ ] Add `TotalDistance` field
- [ ] Add `TotalElevation` field
- [ ] Add `TotalSummitCount` field

**1.4 Create ChallengeProposal Model**
- [ ] New model for proposal system

### Phase 2: DAO Layer (3-4 hours)

**2.1 Challenge DAO Updates**
File: `backend/daos/challengeDao.go`

- [ ] **CreateChallenge:** Generate unique join_code
- [ ] **UpdateChallenge:** Add is_locked check, allow editing unlocked challenges
- [ ] **LockChallenge:** New method to lock a challenge (irreversible)
- [ ] **GetChallengeByJoinCode:** New method for join-by-code
- [ ] **GetChallengesByGoalType:** Filter challenges by goal_type
- [ ] Update progress calculation queries for all goal types

**2.2 Progress Calculation Refactor**
File: `backend/daos/challengeDao.go`

- [ ] **RefreshParticipantProgressDistance:** Calculate total distance for user in date range
- [ ] **RefreshParticipantProgressElevation:** Calculate total elevation
- [ ] **RefreshParticipantProgressSummitCount:** Calculate total summits
- [ ] **RefreshParticipantProgressSpecificSummits:** Existing logic (keep)
- [ ] Update `UpdateParticipantProgress` to handle all types

**2.3 Leaderboard DAO Updates**
- [ ] Update `GetChallengeLeaderboard` to support all goal types
- [ ] Add sorting logic: distance DESC, elevation DESC, summit_count DESC, peaks_completed DESC
- [ ] Update progress percentage calculation

**2.4 Challenge Proposals DAO**
File: `backend/daos/challengeProposalDao.go` (NEW)

- [ ] `CreateProposal`
- [ ] `GetProposalsByUser`
- [ ] `GetPendingProposals` (admin)
- [ ] `ApproveProposal` (creates challenge, marks approved)
- [ ] `RejectProposal`

### Phase 3: Service Layer (4-5 hours)

**3.1 Challenge Service Updates**
File: `backend/services/challengeService.go`

- [ ] **CreateChallenge:** Add join_code generation with retry logic
- [ ] **UpdateChallenge:** Add is_locked validation
- [ ] **LockChallenge:** New service method
- [ ] **JoinChallengeByCode:** New method using join_code
- [ ] **RefreshParticipantProgress:** Dispatch to correct calculator based on goal_type
- [ ] **ProcessActivityForChallenges:** Update to handle all goal types

**3.2 Activity Processing Hook**
File: `backend/services/challengeService.go`

When an activity is synced:
- [ ] For `distance` challenges: Add distance to participant total_distance
- [ ] For `elevation` challenges: Add elevation to participant total_elevation
- [ ] For `summit_count` challenges: Increment total_summit_count if summit detected
- [ ] For `specific_summits` challenges: Existing peak matching logic
- [ ] Update participant progress and check for completion

**3.3 Challenge Proposal Service**
File: `backend/services/challengeProposalService.go` (NEW)

- [ ] `CreateProposal`
- [ ] `GetMyProposals`
- [ ] `GetPendingProposals` (admin only)
- [ ] `ApproveProposal` (admin only)
- [ ] `RejectProposal` (admin only)

### Phase 4: Controller & API Endpoints (2-3 hours)

**4.1 New Endpoints**
File: `backend/controllers/challengeController.go`

- [ ] `PUT /api/challenge/lock` - Lock a challenge (irreversible)
- [ ] `POST /api/challenge/join-by-code` - Join challenge by join_code
- [ ] `GET /api/my-created-challenges` - Get challenges created by user
- [ ] `POST /api/challenge-proposals` - Create proposal
- [ ] `GET /api/my-proposals` - Get user's proposals
- [ ] `GET /api/admin/proposals` - Get pending proposals (admin)
- [ ] `POST /api/admin/proposals/:id/approve` - Approve proposal (admin)
- [ ] `POST /api/admin/proposals/:id/reject` - Reject proposal (admin)

**4.2 Update Existing Endpoints**
- [ ] `PUT /api/challenge` - Add is_locked validation
- [ ] `GET /api/challenges` - Exclude created challenges (move to /my-created-challenges)
- [ ] `POST /api/challenges` - Add goal_type, target_value handling

---

## Frontend Implementation Tasks

### Phase 5: Info/Tutorial Popup Component (2-3 hours)

**5.1 Create Retro RPG Info Popup Component**
File: `frontend/strava-goal/src/app/components/info-popup/info-popup.component.ts` (NEW)

Features:
- [ ] Full-screen overlay with retro pixel-art style border
- [ ] Animated text appearance (typewriter effect)
- [ ] Multiple pages/sections (swipe/arrow navigation)
- [ ] Close button (X)
- [ ] Trigger from "ℹ️" button in header

**5.2 Write Info Content**
File: `frontend/strava-goal/src/app/components/info-popup/challenges-info-content.ts` (NEW)

Sections:
1. **Welcome to Challenges!** - What are challenges?
2. **Challenge Types** - Distance, Elevation, Summit Count, Specific Summits
3. **Modes** - Competitive vs. Collaborative
4. **Creating & Joining** - How to create/join via code
5. **Progress & Leaderboards** - How tracking works
6. **Groups & Challenges** - Group adoption

**5.3 Styling**
File: `frontend/strava-goal/src/app/components/info-popup/info-popup.component.scss` (NEW)

- [ ] Retro pixel border (thick, colored)
- [ ] Text box background (semi-transparent dark)
- [ ] Typewriter animation keyframes
- [ ] Navigation buttons (< >)
- [ ] Responsive layout

### Phase 6: UI Component Library from Groups (3-4 hours)

**6.1 Extract & Generalize Table Components**

Create reusable table components based on groups tables:

File: `frontend/strava-goal/src/app/components/shared/data-table/data-table.component.ts` (NEW)

- [ ] Generic table with sortable columns
- [ ] Expandable rows support
- [ ] Progress bar integration
- [ ] Action buttons (edit/delete)
- [ ] Empty states

**6.2 Extract Progress Bar Component**
File: `frontend/strava-goal/src/app/components/shared/progress-bar/progress-bar.component.ts` (NEW)

- [ ] Percentage-based progress bar
- [ ] Color theming (green/blue/orange based on progress)
- [ ] Completion badge (✓ when 100%)

**6.3 Extract Form Components**
File: `frontend/strava-goal/src/app/components/shared/form-fields/` (NEW)

- [ ] Date range picker
- [ ] Goal type selector (with icons)
- [ ] Number input with units
- [ ] Peak multi-selector

### Phase 7: Challenge List Redesign (4-5 hours)

**7.1 Three-Tab Structure**
File: `frontend/strava-goal/src/app/pages/challenges/challenge-list/challenge-list.component.html`

Redesign to three tabs:

**Tab 1: My Challenges**
- [ ] Table view of joined challenges (not created by me)
- [ ] Columns: Name, Type, Mode, Progress, Target, Deadline, Actions
- [ ] Sort by: progress, deadline
- [ ] Filter: active/completed
- [ ] Empty state: "Join your first challenge!"

**Tab 2: My Created**
- [ ] Table view of challenges I created
- [ ] Columns: Name, Type, Mode, Participants, Code, Locked, Actions
- [ ] Actions: View, Edit, Lock, Delete
- [ ] Show join code for sharing
- [ ] Empty state: "Create your first challenge!"

**Tab 3: Discover**
- [ ] **Public Challenges Section:**
  - Grid of featured challenge cards
  - Filter by goal_type, region, difficulty
  - Empty state: "No public challenges yet"

- [ ] **Join by Code Section:**
  - Text input for join code
  - "Join Challenge" button
  - Recent joins history

- [ ] **Propose Challenge Section:**
  - "Propose a Public Challenge" button
  - Opens proposal form modal

**7.2 Update Challenge Card Component**
File: `frontend/strava-goal/src/app/components/challenges/challenge-card/challenge-card.component.html`

- [ ] Add goal type icon (🏃 💪 🏔️ 🎯)
- [ ] Show target value with units (e.g., "1000 km", "25000 m", "50 summits")
- [ ] Show participant count
- [ ] Show lock status (🔒 badge)
- [ ] Improve visual hierarchy

### Phase 8: Challenge Create/Edit Forms (5-6 hours)

**8.1 Redesign Create Form**
File: `frontend/strava-goal/src/app/components/challenges/challenge-create-form/challenge-create-form.component.html`

New fields:
- [ ] Challenge Name
- [ ] Description
- [ ] **Goal Type:** Dropdown (Distance/Elevation/Summit Count/Specific Summits)
- [ ] **Competition Mode:** Radio (Competitive/Collaborative)
- [ ] **Target Value:** Number input (shown for distance/elevation/summit_count)
- [ ] **Target Units:** Auto-set based on goal_type (km/m/summits)
- [ ] **Peaks Selector:** Multi-select (shown for specific_summits)
- [ ] **Visibility:** Radio (Private/Discoverable)
- [ ] Start Date
- [ ] Deadline
- [ ] Region (optional)
- [ ] Difficulty (optional)

Dynamic behavior:
- [ ] Show/hide target value vs. peaks selector based on goal_type
- [ ] Validate target value > 0
- [ ] Validate at least one peak for specific_summits

**8.2 Create Edit Form**
File: `frontend/strava-goal/src/app/components/challenges/challenge-edit-form/challenge-edit-form.component.ts` (NEW)

Similar to create form, but:
- [ ] Pre-populate with existing challenge data
- [ ] Disable editing if challenge is locked
- [ ] Show "Lock Challenge" button (with warning)
- [ ] Show participant count and warning if editing will affect them

**8.3 Lock Confirmation Dialog**
File: `frontend/strava-goal/src/app/components/challenges/lock-challenge-dialog/lock-challenge-dialog.component.ts` (NEW)

- [ ] Warning message: "This action is irreversible!"
- [ ] Show what locking means
- [ ] Require checkbox confirmation
- [ ] "Lock It In" button (primary, danger color)
- [ ] Cancel button

### Phase 9: Challenge Detail Page Overhaul (6-7 hours)

**9.1 Header Redesign**
File: `frontend/strava-goal/src/app/pages/challenges/challenge-detail/challenge-detail.component.html`

- [ ] Show join code prominently with copy button
- [ ] Add goal type badge
- [ ] Show target value/units
- [ ] Show lock status (🔒 badge if locked)
- [ ] Edit button (if creator and unlocked)
- [ ] Lock button (if creator and unlocked)
- [ ] Delete button (if creator)

**9.2 Progress Section Updates**

For different goal types:
- [ ] **Distance:** Show km progress bar, pace graph (optional)
- [ ] **Elevation:** Show meters progress bar, elevation profile (optional)
- [ ] **Summit Count:** Show summit count, recent summits list
- [ ] **Specific Summits:** Existing peak checklist

**9.3 Leaderboard Tab Enhancements**

Use table component (like groups page):

**Competitive Mode:**
- [ ] Table columns: Rank, User, Progress Bar, Target Achieved, Completion Date
- [ ] Trophy icons (🥇🥈🥉) for top 3
- [ ] Highlight current user row
- [ ] Show percentile rank

**Collaborative Mode:**
- [ ] Team progress at top (large progress bar)
- [ ] Individual contributions table below
- [ ] Show who's leading contributions

**9.4 Activity/Log Tab**

Different views per goal type:
- [ ] **Distance/Elevation:** Activity feed with distance/elevation contributions
- [ ] **Summit Count:** Summit log with peak names
- [ ] **Specific Summits:** Existing summit log

**9.5 Group Adoption Tab** (if challenge adopted by groups)
- [ ] List of groups that adopted this challenge
- [ ] Group-specific leaderboard for each group
- [ ] Join group button if not a member

### Phase 10: Join by Code Flow (2-3 hours)

**10.1 Join by Code Component**
File: `frontend/strava-goal/src/app/components/challenges/join-by-code/join-by-code.component.ts` (NEW)

- [ ] Text input for 6-character code
- [ ] Auto-uppercase and validate format
- [ ] "Find Challenge" button
- [ ] Show challenge preview before joining
- [ ] "Join Challenge" confirmation
- [ ] Error handling (invalid code, already joined)

**10.2 Integration Points**
- [ ] Discover tab
- [ ] Direct URL: `/challenges/join/:code` (auto-redirect to challenge if found)
- [ ] Share button on challenge detail page (copies link with code)

### Phase 11: Proposal System (3-4 hours)

**11.1 Proposal Form**
File: `frontend/strava-goal/src/app/components/challenges/challenge-proposal-form/challenge-proposal-form.component.ts` (NEW)

Similar to create form, but:
- [ ] Title: "Propose a Public Challenge"
- [ ] Explanation text: "Submit your challenge idea for admin review..."
- [ ] All fields same as create form
- [ ] Submit button: "Submit Proposal"

**11.2 My Proposals View**
File: `frontend/strava-goal/src/app/pages/challenges/my-proposals/my-proposals.component.ts` (NEW)

- [ ] Table of user's proposals
- [ ] Columns: Name, Type, Status (Pending/Approved/Rejected), Submitted Date, Admin Notes
- [ ] Status badges (color-coded)
- [ ] View proposal details

**11.3 Admin Proposal Review** (Admin-only page)
File: `frontend/strava-goal/src/app/pages/admin/proposal-review/proposal-review.component.ts` (NEW)

- [ ] List of pending proposals
- [ ] Preview proposal details
- [ ] Approve/Reject buttons
- [ ] Add admin notes field
- [ ] Approval creates public challenge with is_featured=true

### Phase 12: Services & State Management (3-4 hours)

**12.1 Update Challenge Service**
File: `frontend/strava-goal/src/app/services/challenge.service.ts`

New signals:
- [ ] `myCreatedChallenges` - Challenges I created
- [ ] `myJoinedChallenges` - Challenges I joined (not created)
- [ ] `publicChallenges` - Featured/public challenges
- [ ] `proposals` - My proposals

New methods:
- [ ] `createChallenge(request)` - Support all goal types
- [ ] `updateChallenge(id, request)` - Edit challenge
- [ ] `lockChallenge(id)` - Lock challenge
- [ ] `joinByCode(code)` - Join by code
- [ ] `createProposal(request)` - Submit proposal
- [ ] `getMyProposals()` - Load proposals
- [ ] `approveProposal(id)` (admin)
- [ ] `rejectProposal(id, notes)` (admin)

**12.2 Computed Values**
- [ ] Active vs. completed challenges
- [ ] Progress percentages for all goal types
- [ ] Leaderboard rankings

### Phase 13: Remove Groups Goals UI (1-2 hours)

**13.1 Remove Components**
- [ ] Delete `groups-goals-table` component
- [ ] Delete `goals-create-form` component
- [ ] Delete `goals-edit-form` component
- [ ] Delete `goals-peak-selector` component
- [ ] Delete `goal-delete-confirmation` component

**13.2 Update Groups Page**
File: `frontend/strava-goal/src/app/pages/groups/groups.component.html`

- [ ] Remove Goals tab
- [ ] Add "Challenges" tab that shows group's adopted challenges
- [ ] Link to challenge detail pages

**13.3 Update Groups Service**
- [ ] Remove goal-related signals and methods
- [ ] Keep group-challenge adoption methods

---

## Testing Plan

### Unit Tests
- [ ] Join code generation (uniqueness, format)
- [ ] Progress calculation for all goal types
- [ ] Leaderboard ranking for all goal types
- [ ] Challenge locking validation

### Integration Tests
- [ ] Create challenge → Join → Record activity → Progress updates
- [ ] Group goals migration → Verify challenge creation
- [ ] Proposal → Approve → Public challenge created

### Manual Testing Checklist

**Challenge Creation:**
- [ ] Create distance challenge (competitive)
- [ ] Create elevation challenge (collaborative)
- [ ] Create summit count challenge (competitive)
- [ ] Create specific summits challenge (collaborative)
- [ ] Verify join code generated and unique

**Joining Challenges:**
- [ ] Join by code (valid code)
- [ ] Join by code (invalid code)
- [ ] Join public challenge
- [ ] Verify auto-join for group-adopted challenges

**Progress Tracking:**
- [ ] Record activity → Distance challenge progress updates
- [ ] Record activity → Elevation challenge progress updates
- [ ] Summit detected → Summit count challenge updates
- [ ] Summit detected → Specific summits challenge updates

**Leaderboards:**
- [ ] Competitive distance leaderboard ranks correctly
- [ ] Competitive elevation leaderboard ranks correctly
- [ ] Competitive summit count leaderboard ranks correctly
- [ ] Collaborative mode shows team total

**Challenge Management:**
- [ ] Edit unlocked challenge
- [ ] Attempt to edit locked challenge (should fail)
- [ ] Lock challenge (irreversible warning)
- [ ] Delete challenge (confirmation)

**Proposals:**
- [ ] Submit proposal
- [ ] View my proposals
- [ ] Admin approve proposal → Public challenge created
- [ ] Admin reject proposal → Status updated

**Migration:**
- [ ] Existing group goals appear as challenges
- [ ] Group members auto-joined to migrated challenges
- [ ] Progress preserved

**UI/UX:**
- [ ] Three tabs render correctly
- [ ] Tables from groups page style applied
- [ ] Info popup displays correctly
- [ ] Empty states show for all tabs
- [ ] Responsive on mobile

---

## Rollout Plan

### Phase A: Backend Foundation (Complete database & backend)
1. Run database migrations
2. Update models, DAOs, services, controllers
3. Deploy backend
4. Verify API endpoints with Postman/curl

### Phase B: Frontend Components (Build UI components)
1. Create info popup component
2. Extract table components from groups page
3. Test components in isolation

### Phase C: Challenge Flows (Build main features)
1. Update challenge create/edit forms
2. Implement three-tab structure
3. Build join-by-code flow
4. Update challenge detail page

### Phase D: Proposals & Polish (Add proposal system)
1. Build proposal submission
2. Build admin review UI
3. Polish empty states, error handling
4. Add loading states

### Phase E: Migration & Cleanup (Remove old system)
1. Run group goals → challenges migration
2. Remove groups goals UI
3. Update groups page
4. Deprecate old endpoints

### Phase F: Testing & Launch
1. Full manual testing
2. Beta testing with real users
3. Fix bugs
4. Write user guide
5. Launch! 🚀

---

## Estimated Timeline

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| 1 | Database & Models | 2-3 hours |
| 2 | DAO Layer | 3-4 hours |
| 3 | Service Layer | 4-5 hours |
| 4 | Controllers & API | 2-3 hours |
| 5 | Info Popup | 2-3 hours |
| 6 | UI Components | 3-4 hours |
| 7 | Challenge List | 4-5 hours |
| 8 | Create/Edit Forms | 5-6 hours |
| 9 | Detail Page Overhaul | 6-7 hours |
| 10 | Join by Code | 2-3 hours |
| 11 | Proposal System | 3-4 hours |
| 12 | Services & State | 3-4 hours |
| 13 | Remove Groups Goals | 1-2 hours |
| Testing & Polish | | 4-6 hours |
| **TOTAL** | | **44-59 hours** |

---

## Open Questions

1. **Admin Interface:** Should we build a full admin panel, or just add admin-only routes to the main app?
   - **Recommendation:** Admin-only routes for now (/admin/proposals)

2. **Password Protection:** Do we want optional password protection for private challenges?
   - **Recommendation:** Skip for v1, can add later if needed

3. **Challenge Categories/Tags:** Should challenges be taggable (e.g., "Trail Running", "Peak Bagging")?
   - **Recommendation:** Use `region` and `difficulty` for now, add tags later

4. **Notifications:** Should we notify users when:
   - Someone joins their challenge?
   - A proposal is approved/rejected?
   - They complete a challenge?
   - **Recommendation:** Add basic notifications (can be simple alerts/toasts)

5. **Group-Wide Challenges:** Should group admins be able to create challenges that auto-include all members?
   - **Recommendation:** Use group adoption feature + make join_code shareable in group chat

---

## Success Metrics

After launch, we'll track:
- Number of challenges created (by type)
- Number of challenges joined
- Completion rates (competitive vs. collaborative)
- User engagement (time spent on challenges page)
- Proposal submission rate
- Public challenge participation

---

## Future Enhancements (Post-Launch)

- Challenge templates (quick-create common challenges)
- Challenge badges/achievements
- Challenge chat/comments
- Recurring challenges (e.g., monthly summit count)
- Challenge rewards/prizes
- Social sharing (share challenge card to social media)
- Challenge analytics dashboard for creators

---

## Notes

- Re-use as much existing infrastructure as possible (DAO patterns, progress calculation logic)
- Maintain backward compatibility during migration
- Keep the retro RPG aesthetic consistent
- Focus on intuitive UX - challenges should be the easiest feature to understand
- Make leaderboards fun and competitive
- Ensure collaborative mode feels rewarding (team achievement)

---

**Status:** Ready for Implementation
**Next Step:** Review plan, get approval, start Phase 1 (Database & Models)
