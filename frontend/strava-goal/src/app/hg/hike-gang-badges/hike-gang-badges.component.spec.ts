import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HikeGangBadgesComponent } from './hike-gang-badges.component';

describe('HikeGangBadgesComponent', () => {
  let component: HikeGangBadgesComponent;
  let fixture: ComponentFixture<HikeGangBadgesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HikeGangBadgesComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(HikeGangBadgesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
