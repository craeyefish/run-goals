import { Component } from '@angular/core';
import { AuthComponent } from 'src/app/auth/auth.component';

@Component({
  selector: 'auth-layout',
  standalone: true,
  imports: [AuthComponent],
  templateUrl: 'auth-layout.component.html',
  styleUrl: 'auth-layout.component.scss',
})
export class AuthLayoutComponent { }
