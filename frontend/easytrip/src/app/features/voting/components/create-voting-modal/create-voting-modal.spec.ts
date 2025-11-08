import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateVotingModal } from './create-voting-modal';

describe('CreateVotingModal', () => {
  let component: CreateVotingModal;
  let fixture: ComponentFixture<CreateVotingModal>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreateVotingModal]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreateVotingModal);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
