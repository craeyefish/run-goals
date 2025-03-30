import { Component } from '@angular/core';
import { HeaderComponent } from 'src/app/components/header/header.component';
import { SidebarComponent } from 'src/app/components/sidebar/sidebar.component';
import { HomeComponent } from 'src/app/pages/home/home.component';

@Component({
  selector: 'main-layout',
  standalone: true,
  imports: [HeaderComponent, SidebarComponent, HomeComponent],
  templateUrl: 'main-layout.component.html',
  styleUrl: 'main-layout.component.scss',
})
export class MainLayoutComponent { }
