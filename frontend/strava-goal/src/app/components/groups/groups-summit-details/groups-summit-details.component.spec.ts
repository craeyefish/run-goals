import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GroupsSummitDetailsComponent } from './groups-summit-details.component';

describe('GroupsSummitDetailsComponent', () => {
  let component: GroupsSummitDetailsComponent;
  let fixture: ComponentFixture<GroupsSummitDetailsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GroupsSummitDetailsComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(GroupsSummitDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
