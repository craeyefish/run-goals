import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GoalDeleteConfirmationComponent } from './goal-delete-confirmation.component';

describe('GoalDeleteConfirmationComponent', () => {
  let component: GoalDeleteConfirmationComponent;
  let fixture: ComponentFixture<GoalDeleteConfirmationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GoalDeleteConfirmationComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(GoalDeleteConfirmationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
