import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StravaCallbackComponent } from './strava-callback.component';

describe('StravaCallbackComponent', () => {
  let component: StravaCallbackComponent;
  let fixture: ComponentFixture<StravaCallbackComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [StravaCallbackComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(StravaCallbackComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
