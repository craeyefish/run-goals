import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PeakSummitTableComponent } from './peak-summit-table.component';

describe('PeakSummitTableComponent', () => {
  let component: PeakSummitTableComponent;
  let fixture: ComponentFixture<PeakSummitTableComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [PeakSummitTableComponent]
    });
    fixture = TestBed.createComponent(PeakSummitTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
