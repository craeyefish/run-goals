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
    this.hgService.loadActivities();
    this.hgService.activities$.subscribe((activities) => {
      if (activities) {
        // Filter activities tagged with #hg (case insensitive)
        const hgActivities = activities.filter((act) =>
          act.name?.toLowerCase().includes('#hg')
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
    activities.forEach((activity) => {
      if (!activity.description) {
        return;
      }

      // Extract medal tags from the description (e.g., #hobbit_feet-gold)
      const medalRegex = /#(\w+)-(\w+)/g;
      let match;

      while ((match = medalRegex.exec(activity.description)) !== null) {
        const [_, badgeKey, tier] = match;

        // Find the badge and update its tier
        const badge = this.badges.find((b) =>
          b.name.toLowerCase().replace(/\s+/g, '_').includes(badgeKey)
        );
        if (badge) {
          badge.tier = tier as 'bronze' | 'silver' | 'gold';
        }
      }
    });
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

  triggerSync(): void {
    this.syncStatus = 'Triggering sync...';
    this.hgService.triggerSync().subscribe({
      next: (response) => {
        this.syncStatus = response.message || 'Sync triggered successfully';
        // Clear status after 10 seconds
        setTimeout(() => {
          this.syncStatus = '';
        }, 10000);
      },
      error: (error) => {
        this.syncStatus = 'Error: Failed to trigger sync';
        setTimeout(() => {
          this.syncStatus = '';
        }, 10000);
      },
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
