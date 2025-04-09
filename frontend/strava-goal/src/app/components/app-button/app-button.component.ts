import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-button',
  standalone: true,
  imports: [],
  templateUrl: './app-button.component.html',
  styleUrl: './app-button.component.scss',
})
export class AppButtonComponent {
  constructor() { }

  @Input() onClick?: () => void;
}
