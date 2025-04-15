import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HikeGangComponent } from './hike-gang.component';

describe('HikeGangComponent', () => {
  let component: HikeGangComponent;
  let fixture: ComponentFixture<HikeGangComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HikeGangComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(HikeGangComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
