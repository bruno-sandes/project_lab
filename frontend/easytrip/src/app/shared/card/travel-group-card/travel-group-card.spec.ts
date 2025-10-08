import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TravelGroupCard } from './travel-group-card';

describe('TravelGroupCard', () => {
  let component: TravelGroupCard;
  let fixture: ComponentFixture<TravelGroupCard>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TravelGroupCard]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TravelGroupCard);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
