import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HikeGangCoverComponent } from './hike-gang-cover.component';

describe('HikeGangCoverComponent', () => {
  let component: HikeGangCoverComponent;
  let fixture: ComponentFixture<HikeGangCoverComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HikeGangCoverComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(HikeGangCoverComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
