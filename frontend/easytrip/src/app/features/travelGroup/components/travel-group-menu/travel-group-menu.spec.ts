import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TravelGroupMenu } from './travel-group-menu';

describe('TravelGroupMenu', () => {
  let component: TravelGroupMenu;
  let fixture: ComponentFixture<TravelGroupMenu>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TravelGroupMenu]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TravelGroupMenu);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
