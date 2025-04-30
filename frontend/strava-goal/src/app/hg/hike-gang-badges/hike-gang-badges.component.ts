import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-hike-gang-badges',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './hike-gang-badges.component.html',
  styleUrls: ['./hike-gang-badges.component.scss'],
})
export class HikeGangBadgesComponent {
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
      tier: 'gold',
    },
    {
      name: 'Dawn Patrol',
      description: 'Start a hike before sunrise',
      image: '/assets/badges/dawn-patrol-gold.png',
      tier: 'gold',
    },
    // {
    //   name: 'Peak Seeker',
    //   description: 'Reach the summit of a mountain',
    //   image: '/assets/badges/peak-seeker-gold.png',
    //   tier: 'gold',
    // },
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
    // {
    //   name: 'Night Walker',
    //   description: 'Hike under the stars',
    //   image: '/assets/badges/night-walker.png',
    //   tier: 'silver',
    // },
    {
      name: 'Storm Braver',
      description: 'Complete a hike in the rain',
      image: '/assets/badges/storm-braver-gold.png',
      tier: 'unachieved',
    },
    {
      name: 'Tiny Steps',
      description: 'Join your first hike!',
      image: '/assets/badges/tiny-steps-gold.png',
      tier: 'unachieved',
    },
    {
      name: 'Pajama Party',
      description: 'Do a hike in pajamas',
      image: '/assets/badges/pajama-party-gold.png',
      tier: 'unachieved',
    },
    {
      name: 'Bird Spotter',
      description: 'Spot a certain bird on your travvels',
      image: '/assets/badges/bird-spotter-gold.png',
      tier: 'unachieved',
    },
  ];

  showInfo(badge: any) {
    this.selectedBadge = badge;
  }

  getTierFilter(tier: string): string {
    switch (tier) {
      case 'silver':
        return 'grayscale(1) brightness(1.2) contrast(1.1)';
      case 'bronze':
        return 'sepia(1) hue-rotate(-25deg) saturate(3) brightness(0.6)';
      case 'unachieved':
        return 'grayscale(1) brightness(0.4) contrast(0.8)';
      default:
        return 'none';
    }
  }
}
