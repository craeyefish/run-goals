import { ComponentFixture, TestBed } from "@angular/core/testing";

import { SummitsComponent } from "./summits.component";

describe("GroupsComponent", () => {
  let component: SummitsComponent;
  let fixture: ComponentFixture<SummitsComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [SummitsComponent],
    });
    fixture = TestBed.createComponent(SummitsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
