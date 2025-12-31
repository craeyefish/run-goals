# Plan 004: Challenge-Centric Redesign

## Vision Statement

Transform Summit Seekers from a "Strava activity tracker with summit detection" into a **challenge-driven summit hunting platform**. Challenges become the core unit of engagement - users discover, join, create, and complete challenges either solo or with groups.

---

## Target User Experience

### The Solo Hiker
> "I want to complete the Cape Town 13 Peaks challenge this year. I can see my progress, which peaks I've done, and plan my next hike."

### The Social Group
> "My hiking group is doing a winter challenge - first to summit 10 peaks wins. We can see each other's progress and it keeps us motivated."

### The Goal Setter
> "I just want to summit 20 peaks this year, doesn't matter which ones. The app tracks my progress and celebrates when I hit milestones."

---

## Information Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         NAVIGATION                               │
│  [Home]  [Challenges]  [Explore]  [Groups]  [Profile]           │
└─────────────────────────────────────────────────────────────────┘

HOME
├── Stats Banner (distance, elevation, summits this year)
├── Active Challenges (cards with progress bars)
│   ├── Personal challenges
│   ├── Group challenges I'm part of
│   └── Quick actions: "View All", "Join Challenge"
├── Recent Summits (last 3-5 with quick links)
└── 2025 Summit Goal progress (if set)

CHALLENGES (★ New Primary Section)
├── My Challenges
│   ├── Active (in progress)
│   ├── Completed (archive)
│   └── Created by me
├── Discover
│   ├── Featured/Popular challenges
│   ├── By region (Cape Town, Garden Route, etc.)
│   └── Search/filter
├── Create Challenge
│   ├── Pick peaks (using peak-picker component)
│   ├── Set deadline (optional)
│   ├── Set type: Collaborative vs Competitive
│   └── Visibility: Private / Friends / Public
└── Challenge Detail Page
    ├── Header: Name, description, deadline, progress
    ├── Peak List: Checkmarks for completed, click to see on map
    ├── Map View: All challenge peaks highlighted
    ├── Participants: Who else is doing this
    ├── Leaderboard (if competitive)
    └── Activity Feed: Recent summits in this challenge

EXPLORE (Combines Map + Discovery)
├── Interactive Map
│   ├── All peaks (color-coded by completion status)
│   ├── My activity routes
│   ├── Filter by: Region, Elevation, Completion
│   └── Challenge overlay toggle
├── Peak Directory
│   ├── Searchable list of all peaks
│   ├── Click → Peak detail (elevation, region, challenges containing it)
│   └── "Add to Challenge" quick action
└── Tabs or toggle: Map View / List View

GROUPS (Social Layer - De-emphasized)
├── My Groups
├── Group Detail
│   ├── Members
│   ├── Group Challenges (shared challenges)
│   ├── Leaderboard
│   └── Activity Feed
└── Create/Join Group

PROFILE
├── Summit Log (all summits, sortable/filterable)
├── Activity History
├── Stats & Achievements
├── 2025 Goal Settings
└── Account Settings
```

---

## Data Model

### New Tables

```sql
-- Core challenge definition
CREATE TABLE challenges (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Type determines behavior
    challenge_type VARCHAR(50) NOT NULL, 
    -- 'predefined'    = Admin-created, discoverable (e.g., "Cape Town 13 Peaks")
    -- 'custom'        = User-created, can be private or public
    -- 'yearly_goal'   = Special type, deadline = end of year, tracks count not specific peaks
    
    -- Competition mode
    competition_mode VARCHAR(50) NOT NULL DEFAULT 'collaborative',
    -- 'collaborative' = Group works together to complete all peaks
    -- 'competitive'   = Leaderboard, first to complete wins / most summits wins
    
    -- Visibility
    visibility VARCHAR(50) NOT NULL DEFAULT 'private',
    -- 'private'  = Only creator and invited participants
    -- 'friends'  = Friends of creator can see
    -- 'public'   = Discoverable by anyone
    
    -- Dates
    start_date DATE,
    deadline DATE,              -- NULL = no deadline
    
    -- Ownership
    created_by_user_id BIGINT REFERENCES users(id),
    created_by_group_id BIGINT REFERENCES groups(id), -- If group-owned
    
    -- For yearly goals specifically
    target_count INT,           -- "Summit 20 peaks" = 20
    
    -- Metadata
    region VARCHAR(255),        -- For filtering: "Cape Town", "Western Cape"
    difficulty VARCHAR(50),     -- 'easy', 'moderate', 'hard', 'extreme'
    is_featured BOOLEAN DEFAULT FALSE,  -- Admin can feature challenges
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Peaks in a challenge (not used for yearly_goal type)
CREATE TABLE challenge_peaks (
    id BIGSERIAL PRIMARY KEY,
    challenge_id BIGINT REFERENCES challenges(id) ON DELETE CASCADE,
    peak_id BIGINT REFERENCES peaks(id) ON DELETE CASCADE,
    sort_order INT,             -- For ordered challenges (e.g., traverse routes)
    
    UNIQUE(challenge_id, peak_id)
);

-- User participation in challenges
CREATE TABLE challenge_participants (
    id BIGSERIAL PRIMARY KEY,
    challenge_id BIGINT REFERENCES challenges(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    
    joined_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,     -- When user completed the challenge
    
    -- Cached progress (updated on summit detection)
    peaks_completed INT DEFAULT 0,
    total_peaks INT,            -- Denormalized for quick queries
    
    UNIQUE(challenge_id, user_id)
);

-- Group participation (group adopts a challenge)
CREATE TABLE challenge_groups (
    id BIGSERIAL PRIMARY KEY,
    challenge_id BIGINT REFERENCES challenges(id) ON DELETE CASCADE,
    group_id BIGINT REFERENCES groups(id) ON DELETE CASCADE,
    
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    
    -- Group can override deadline
    deadline_override DATE,
    
    UNIQUE(challenge_id, group_id)
);

-- Track which summits count toward which challenges
-- (Derived from user_peaks but cached for performance)
CREATE TABLE challenge_summit_log (
    id BIGSERIAL PRIMARY KEY,
    challenge_id BIGINT REFERENCES challenges(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    peak_id BIGINT REFERENCES peaks(id),
    activity_id BIGINT REFERENCES activity(id),
    summited_at TIMESTAMP,
    
    UNIQUE(challenge_id, user_id, peak_id)  -- One credit per peak per challenge
);
```

### Relationship Diagram

```
┌──────────────┐     ┌─────────────────────┐     ┌──────────────┐
│    users     │────▶│ challenge_participants│◀────│  challenges  │
└──────────────┘     └─────────────────────┘     └──────────────┘
       │                                                │
       │              ┌─────────────────────┐           │
       └─────────────▶│ challenge_summit_log │◀─────────┘
                      └─────────────────────┘
                                │
                      ┌─────────────────────┐
                      │    challenge_peaks   │
                      └─────────────────────┘
                                │
                      ┌─────────────────────┐
                      │       peaks          │
                      └─────────────────────┘

┌──────────────┐     ┌─────────────────────┐     ┌──────────────┐
│    groups    │────▶│  challenge_groups    │◀────│  challenges  │
└──────────────┘     └─────────────────────┘     └──────────────┘
```

---

## Implementation Phases

### Phase 1: Foundation (Backend + Data Model)
**Goal**: Build the challenge infrastructure without changing existing UI

- [ ] Create database migrations for new tables
- [ ] Create Challenge model and DAO
- [ ] Create ChallengeService with CRUD operations
- [ ] Create ChallengeController with REST endpoints
- [ ] Integrate with existing summit detection (when summit detected → check challenges)
- [ ] Seed 2-3 predefined challenges (Cape Town 13 Peaks, etc.)

**Endpoints**:
```
GET    /api/challenges                    # List challenges (with filters)
GET    /api/challenges/:id                # Challenge detail
POST   /api/challenges                    # Create challenge
PUT    /api/challenges/:id                # Update challenge
DELETE /api/challenges/:id                # Delete challenge

GET    /api/challenges/:id/peaks          # Peaks in challenge
POST   /api/challenges/:id/peaks          # Add peak to challenge
DELETE /api/challenges/:id/peaks/:peakId  # Remove peak

POST   /api/challenges/:id/join           # Join challenge
DELETE /api/challenges/:id/leave          # Leave challenge
GET    /api/challenges/:id/participants   # List participants
GET    /api/challenges/:id/leaderboard    # Leaderboard (competitive mode)

GET    /api/my/challenges                 # User's challenges
GET    /api/groups/:id/challenges         # Group's challenges
POST   /api/groups/:id/challenges/:challengeId/adopt  # Group adopts challenge
```

**Deliverable**: API fully functional, testable via curl/Postman

---

### Phase 2: Challenge Detail Page (Frontend)
**Goal**: Build the challenge experience in isolation

- [ ] Challenge detail page component
  - [ ] Header with progress ring/bar
  - [ ] Peak list with completion indicators
  - [ ] Map view showing challenge peaks
  - [ ] Participant list
  - [ ] Leaderboard (if competitive)
- [ ] Join/Leave challenge actions
- [ ] Share challenge functionality

**Deliverable**: Can view and interact with a single challenge

---

### Phase 3: Challenges List & Discovery (Frontend)
**Goal**: Users can browse and find challenges

- [ ] My Challenges page (active, completed, created)
- [ ] Discover Challenges page
  - [ ] Featured section
  - [ ] Filter by region, difficulty
  - [ ] Search
- [ ] Challenge cards component (reusable)

**Deliverable**: Full challenges section functional

---

### Phase 4: Create Challenge Flow (Frontend)
**Goal**: Users can create their own challenges

- [ ] Create challenge form
  - [ ] Name, description
  - [ ] Peak picker integration (reuse existing component!)
  - [ ] Deadline picker (optional)
  - [ ] Competition mode toggle
  - [ ] Visibility settings
- [ ] Edit challenge (for owners)
- [ ] Delete challenge (with confirmation)

**Deliverable**: Full CRUD for challenges from UI

---

### Phase 5: Home Page Integration ✅
**Goal**: Challenges become central to the dashboard

- [x] Active Challenges section on home
  - [x] Cards with progress indicators
  - [x] Quick actions (View All, navigate to challenge detail)
- [x] Keep old "wishlist" for now (to be migrated later)
- [x] Keep stats banner and charts
- [x] Recent summits section (last 5)

**Deliverable**: Home page showcases challenge progress

---

### Phase 6: Explore Page Consolidation
**Goal**: Combine map, summits, activities into unified exploration

- [ ] New Explore page with map as primary view
- [ ] Sidebar/panel with:
  - [ ] Peak list (filterable)
  - [ ] Activity list
  - [ ] Challenge overlays
- [ ] Peak detail modal (shows challenges containing this peak)
- [ ] Retire separate Summits and Activities pages (or make them tabs)

**Deliverable**: Single exploration experience

---

### Phase 7: Group Challenges
**Goal**: Groups can adopt and compete in challenges

- [ ] Group challenge section in group detail page
- [ ] "Adopt Challenge" flow for group admins
- [ ] Group leaderboard for competitive challenges
- [ ] Group progress for collaborative challenges

**Deliverable**: Social challenge experience complete

---

### Phase 8: Yearly Goals (Special Challenge Type)
**Goal**: Personal yearly summit goals as a challenge variant

- [ ] "Set 2025 Goal" flow → creates yearly_goal challenge
- [ ] Special UI for count-based goals (not peak-specific)
- [ ] Milestone celebrations (25%, 50%, 75%, 100%)
- [ ] Year-end summary

**Deliverable**: Yearly goals integrated into challenge system

---

### Phase 9: Polish & Quality of Life
- [ ] Challenge completion celebrations (confetti, badges)
- [ ] Push notifications for challenge milestones
- [ ] Challenge activity feed
- [ ] Performance optimization (caching challenge progress)
- [ ] Mobile responsiveness audit

---

## Future Phases (Post-MVP)

### Community Features
- [ ] User-submitted challenges (with admin review)
- [ ] Challenge ratings and reviews
- [ ] "Fork" a challenge (copy and customize)

### Peak Data Improvements
- [ ] Submit peak information corrections
- [ ] Submit peak photos
- [ ] Submit GPX routes to peaks
- [ ] Summit verification (photo proof)

### Gamification
- [ ] Achievement badges
- [ ] Summit streak tracking
- [ ] Seasonal challenges (admin-created, time-limited)

---

## Migration Strategy

### Existing Data
1. **User wishlist peaks** → Convert to personal "My Wishlist" custom challenge
2. **Group goals** → Keep as-is initially, consider deprecating later
3. **user_peaks** → Remains source of truth, challenge_summit_log is derived

### Feature Flags
Use feature flags to gradually roll out:
- `CHALLENGES_ENABLED` - Show challenges nav item
- `NEW_EXPLORE_PAGE` - Use consolidated explore vs old separate pages
- `YEARLY_GOALS_V2` - Use challenge-based yearly goals

---

## Success Metrics

1. **Engagement**: Users joining challenges (target: 50% of active users in a challenge)
2. **Completion**: Challenge completion rate (target: 30% of started challenges completed)
3. **Creation**: User-created challenges (target: 10% of users create a challenge)
4. **Social**: Group challenge adoption (target: 50% of groups have active challenge)

---

## Open Questions

1. **Should yearly goals allow BOTH specific peaks AND count-based?**
   - Decision: No, pick one per user per year. Keeps it simple.

2. **Can a peak be in multiple challenges?**
   - Decision: Yes, absolutely. "Table Mountain" might be in 5 different challenges.

3. **What happens when a user summits a peak in a challenge they haven't joined?**
   - Decision: Nothing automatically. They need to join first. But we could show "You've already completed 3/13 peaks in this challenge!" as encouragement.

4. **How to handle challenge deadlines?**
   - Decision: Soft deadline - challenge doesn't "fail", just shows "completed after deadline" or stays in progress.

---

## Technical Notes

### Summit Detection Integration
When a summit is detected (in `summitService.go`):
1. Record in `user_peaks` (existing)
2. Find all challenges user is participating in
3. Check if peak is in any of those challenges
4. Update `challenge_summit_log` and `challenge_participants.peaks_completed`
5. Check for challenge completion → trigger celebration

### Performance Considerations
- Cache challenge progress in `challenge_participants`
- Leaderboard queries could be expensive → consider materialized views
- Challenge discovery queries need good indexing on region, visibility, is_featured

### Existing Code Reuse
- **PeakPickerComponent** → Perfect for challenge creation
- **ActivityMapComponent** → Adapt for challenge map view
- **StatsService patterns** → Use for challenge progress calculations
