import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VoteModal } from './vote-modal';

describe('VoteModal', () => {
  let component: VoteModal;
  let fixture: ComponentFixture<VoteModal>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VoteModal]
    })
    .compileComponents();

    fixture = TestBed.createComponent(VoteModal);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
