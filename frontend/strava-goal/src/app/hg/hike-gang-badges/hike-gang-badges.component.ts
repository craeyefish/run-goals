import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { Activity, HgService } from 'src/app/services/hg.service';

@Component({
  selector: 'app-hike-gang-badges',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './hike-gang-badges.component.html',
  styleUrls: ['./hike-gang-badges.component.scss'],
})
export class HikeGangBadgesComponent implements OnInit {
  selectedBadge: {
    name: string;
    description: string;
    image: string;
    tier: 'bronze' | 'silver' | 'gold';
  } | null = null;

  diagnosticsData: any = null;
  syncStatus: string = '';

  badges = [
    {
      name: 'Long Trek',
      description: 'Complete a hike over 15km',
      image: '/assets/badges/long-trek-gold.png',
      tier: 'unachieved',
    },
    {
      name: 'Arbor Day',
      description: 'Plant a tree during your hike',
      image: '/assets/badges/arbor-day-gold.png',
      tier: 'unachieved',
    },
    {
      name: 'Hobbit Feet',
      description: 'Complete a barefoot hike',
      image: '/assets/badges/hobbit-feet-gold.png',
      tier: 'unachieved',
    },
    {
      name: 'Werewalker',
      description: 'Hike under the full moon',
      image: '/assets/badges/werewalker-gold.png',
      tier: 'unachieved',
    },
    {
      name: 'Storm Braver',
      description: 'Complete a hike in a storm',
      image: '/assets/badges/storm-braver-gold.png',
      tier: 'unachieved',
    },
    {
      name: 'Pajama Party',
      description: 'Do a hike in pajamas',
      image: '/assets/badges/pajama-party-gold.png',
      tier: 'unachieved',
    },
    {
      name: 'Trusted Companion',
      description: 'Take a shelter dog on a hike',
      image: '/assets/badges/trusted-companion-gold.png',
      tier: 'unachieved',
    },
    {
      name: 'Connected Dwellings',
      description: 'Use the hike routes to connect your homes on the map',
      image: '/assets/badges/connected-dwellings-gold.png',
      tier: 'unachieved',
    },
  ];

  constructor(private hgService: HgService, private router: Router) {}

  ngOnInit(): void {
    console.log('HikeGang Badges: Component initializing...');
    this.hgService.loadActivities();
    this.hgService.activities$.subscribe((activities) => {
      console.log(
        'HikeGang Badges: Received activities:',
        activities?.length || 0
      );

      if (activities) {
        // Filter activities tagged with #hg (case insensitive)
        const hgActivities = activities.filter((act) =>
          act.name?.toLowerCase().includes('#hg')
        );

        console.log('Filtered HG activities for badges:', hgActivities.length);
        console.log(
          'Activities with descriptions:',
          hgActivities.filter((a) => a.description).length
        );

        // Update badge tiers based on activity descriptions
        this.updateBadgeTiers(hgActivities);
      }
    });
  }

  goBack(): void {
    this.router.navigate(['/hg']); // Replace '/hg' with the correct route for your home page
  }

  updateBadgeTiers(activities: Activity[]): void {
    console.log('Updating badge tiers for', activities.length, 'activities');

    // Log the structure of the first activity to see what fields are available
    if (activities.length > 0) {
      console.log('First activity structure:', activities[0]);
      console.log('Available activity properties:', Object.keys(activities[0]));

      // Check if there might be badge info in activity names or other fields
      console.log('Sample activity for badge analysis:', {
        name: activities[0].name,
        description: activities[0].description,
        // Check if there are any other fields that might contain badge info
        allFields: activities[0],
      });
    }

    // For testing: let's manually create a test badge
    console.log('=== BADGE TEST ===');
    console.log(
      'Available badges:',
      this.badges.map((b) => ({ name: b.name, tier: b.tier }))
    );

    // Temporary test: manually award a badge to see if the system works
    if (activities.length > 0) {
      const testBadge = this.badges.find((b) => b.name === 'Long Trek');
      if (testBadge) {
        testBadge.tier = 'gold';
        console.log(
          'TEST: Manually awarded Long Trek badge to test badge system'
        );
      }
    }

    activities.forEach((activity, index) => {
      // Log the first few activities in detail
      if (index < 3) {
        console.log(`Activity ${index + 1} details:`, {
          name: activity.name,
          description: activity.description,
          hasDescription: !!activity.description,
          descriptionType: typeof activity.description,
          descriptionLength: activity.description?.length || 0,
        });
      }

      if (!activity.description) {
        if (index < 5) {
          // Only log first 5 to avoid spam
          console.log('Activity without description:', activity.name);
        }
        return;
      }

      console.log(
        'Processing activity description:',
        activity.name,
        activity.description
      );

      // Extract medal tags from the description (e.g., #hobbit_feet-gold)
      const medalRegex = /#(\w+)-(\w+)/g;
      let match;
      let foundBadges = 0;

      while ((match = medalRegex.exec(activity.description)) !== null) {
        foundBadges++;
        const [_, badgeKey, tier] = match;

        console.log('Found badge:', badgeKey, 'tier:', tier);

        // Find the badge and update its tier
        const badge = this.badges.find((b) =>
          b.name.toLowerCase().replace(/\s+/g, '_').includes(badgeKey)
        );
        if (badge) {
          badge.tier = tier as 'bronze' | 'silver' | 'gold';
          console.log('Updated badge:', badge.name, 'to tier:', tier);
        } else {
          console.warn('Badge not found for key:', badgeKey);
        }
      }

      if (foundBadges === 0) {
        console.log('No badge tags found in description for:', activity.name);
      }
    });

    const achievedBadges = this.badges.filter(
      (b) => b.tier !== 'unachieved'
    ).length;
    console.log(
      'Total achieved badges:',
      achievedBadges,
      'out of',
      this.badges.length
    );

    // Debug: Show final badge states
    console.log(
      'Final badge states:',
      this.badges.map((b) => ({
        name: b.name,
        tier: b.tier,
      }))
    );

    console.log('=== BADGE SYSTEM DIAGNOSIS ===');
    console.log('Issue: All activities have empty descriptions');
    console.log(
      'Badge tags should be in activity descriptions like: "#long_trek-gold"'
    );
    console.log('Possible solutions:');
    console.log('1. Check if Strava activities actually have descriptions');
    console.log('2. Verify backend is correctly fetching/storing descriptions');
    console.log('3. Add descriptions to activities manually for testing');
    console.log('===============================');
  }

  showInfo(badge: any) {
    this.selectedBadge = badge;
  }

  getTierFilter(tier: string): string {
    switch (tier) {
      case 'silver':
        return 'grayscale(1) brightness(1.0) contrast(1.3) saturate(1.2)';
      case 'bronze':
        return 'sepia(0.8) hue-rotate(-18deg) saturate(1.5) brightness(0.8)';
      case 'unachieved':
        return 'grayscale(1) brightness(0.4) contrast(0.8)';
      default:
        return 'none';
    }
  }

  runDiagnostics(): void {
    console.log('Running activity diagnostics...');
    this.hgService.getDiagnostics().subscribe({
      next: (data) => {
        console.log('Diagnostics data received:', data);
        this.diagnosticsData = data;
      },
      error: (error) => {
        console.error('Error getting diagnostics:', error);
        this.diagnosticsData = { error: 'Failed to get diagnostics' };
      }
    });
  }

  triggerSync(): void {
    console.log('Triggering activity sync...');
    this.syncStatus = 'Triggering sync...';
    this.hgService.triggerSync().subscribe({
      next: (response) => {
        console.log('Sync response:', response);
        this.syncStatus = response.message || 'Sync triggered successfully';
        // Clear status after 10 seconds
        setTimeout(() => {
          this.syncStatus = '';
        }, 10000);
      },
      error: (error) => {
        console.error('Error triggering sync:', error);
        this.syncStatus = 'Error: Failed to trigger sync';
        setTimeout(() => {
          this.syncStatus = '';
        }, 10000);
      }
    });
  }
}

// badges = [
//   {
//     name: 'Peak Seeker',
//     description: 'Reach the summit of a mountain',
//     image: '/assets/badges/peak-seeker-gold.png',
//     tier: 'gold',
//   },
//   {
//     name: 'Night Walker',
//     description: 'Hike under the stars',
//     image: '/assets/badges/night-walker.png',
//     tier: 'silver',
//   },
//   {
//     name: 'Tiny Steps',
//     description: 'Join your first hike!',
//     image: '/assets/badges/tiny-steps-gold.png',
//     tier: 'unachieved',
//   },
//   {
//     name: 'Bird Spotter',
//     description: 'Spot a certain bird on your travels',
//     image: '/assets/badges/bird-spotter-gold.png',
//     tier: 'unachieved',
//   },
// {
//   name: 'Dawn Patrol',
//   description: 'Start a hike before sunrise',
//   image: '/assets/badges/dawn-patrol-gold.png',
//   tier: 'unachieved',
// },
// ];
