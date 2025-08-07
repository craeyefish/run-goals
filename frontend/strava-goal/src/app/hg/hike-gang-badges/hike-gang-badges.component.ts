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

    activities.forEach((activity) => {
      if (!activity.description) {
        console.log('Activity without description:', activity.name);
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
