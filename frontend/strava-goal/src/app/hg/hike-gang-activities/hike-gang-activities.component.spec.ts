import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HikeGangActivitiesComponent } from './hike-gang-activities.component';

describe('HikeGangActivitiesComponent', () => {
  let component: HikeGangActivitiesComponent;
  let fixture: ComponentFixture<HikeGangActivitiesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HikeGangActivitiesComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(HikeGangActivitiesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
