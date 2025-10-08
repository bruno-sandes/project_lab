import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GroupsDashboard } from './groups-dashboard';

describe('GroupsDashboard', () => {
  let component: GroupsDashboard;
  let fixture: ComponentFixture<GroupsDashboard>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GroupsDashboard]
    })
    .compileComponents();

    fixture = TestBed.createComponent(GroupsDashboard);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
