import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreatDestinationModal } from './creat-destination-modal';

describe('CreatDestinationModal', () => {
  let component: CreatDestinationModal;
  let fixture: ComponentFixture<CreatDestinationModal>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreatDestinationModal]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreatDestinationModal);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
