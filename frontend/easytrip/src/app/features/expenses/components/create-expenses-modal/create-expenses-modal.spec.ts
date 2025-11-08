import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateExpensesModal } from './create-expenses-modal';

describe('CreateExpensesModal', () => {
  let component: CreateExpensesModal;
  let fixture: ComponentFixture<CreateExpensesModal>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreateExpensesModal]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreateExpensesModal);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
