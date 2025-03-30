import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { HeaderComponent } from 'src/app/components/header/header.component';
import { SidebarComponent } from 'src/app/components/sidebar/sidebar.component';

@Component({
  selector: 'main-layout',
  standalone: true,
  imports: [HeaderComponent, SidebarComponent, RouterOutlet],
  templateUrl: 'main-layout.component.html',
  styleUrl: 'main-layout.component.scss',
})
export class MainLayoutComponent { }
