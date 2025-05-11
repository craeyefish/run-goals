import { ComponentFixture, TestBed } from "@angular/core/testing";

import { SummitsPageComponent } from "./summits-page.component";

describe("GroupsComponent", () => {
  let component: SummitsPageComponent;
  let fixture: ComponentFixture<SummitsPageComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [SummitsPageComponent],
    });
    fixture = TestBed.createComponent(SummitsPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
