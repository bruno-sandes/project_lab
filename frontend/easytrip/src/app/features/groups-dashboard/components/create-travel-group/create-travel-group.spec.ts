import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateTravelGroup } from './create-travel-group';

describe('CreateTravelGroup', () => {
  let component: CreateTravelGroup;
  let fixture: ComponentFixture<CreateTravelGroup>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreateTravelGroup]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreateTravelGroup);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
