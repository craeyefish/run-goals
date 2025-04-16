import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HikeGangHomeComponent } from './hike-gang-home.component';

describe('HikeGangHomeComponent', () => {
  let component: HikeGangHomeComponent;
  let fixture: ComponentFixture<HikeGangHomeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HikeGangHomeComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(HikeGangHomeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
