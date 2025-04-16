import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Router, RouterLink } from '@angular/router';

@Component({
  selector: 'hike-gang-cover',
  imports: [CommonModule, RouterLink],
  templateUrl: './hike-gang-cover.component.html',
  styleUrl: './hike-gang-cover.component.css',
})
export class HikeGangCoverComponent {
  constructor(private router: Router) {}
}
