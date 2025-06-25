import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GoalsPeakSelectorComponent } from './goals-peak-selector.component';

describe('GoalsPeakSelectorComponent', () => {
  let component: GoalsPeakSelectorComponent;
  let fixture: ComponentFixture<GoalsPeakSelectorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GoalsPeakSelectorComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(GoalsPeakSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
