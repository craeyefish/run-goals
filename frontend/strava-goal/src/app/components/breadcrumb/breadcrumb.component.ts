import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';// Adjust path if needed
import { RouterLink } from '@angular/router';
import { BreadcrumbItem, BreadcrumbService } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'app-breadcrumb',
  imports: [
    RouterLink,
    CommonModule
  ],
  templateUrl: './breadcrumb.component.html',
  styleUrls: ['./breadcrumb.component.scss'],
})
export class BreadcrumbComponent {
  constructor(public breadcrumbService: BreadcrumbService) { }

  breadcrumbItems = this.breadcrumbService.items;

  separator = '->';
}
