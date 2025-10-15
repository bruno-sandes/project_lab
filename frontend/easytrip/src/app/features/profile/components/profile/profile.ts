import { Component, signal } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';

@Component({
  selector: 'app-profile',
  imports: [ReactiveFormsModule],
  templateUrl: './profile.html',
  styleUrl: './profile.css'
})
export class Profile {
  
  userEmail = signal<string>('');
  
  profileForm = new FormGroup({
    name: new FormControl('', { 
      nonNullable: true, 
      validators: [Validators.required, Validators.minLength(3)] 
    }),
   
  });



  ngOnInit(): void {
    this.loadProfile();
  }

  loadProfile() {
    const initialName = 'Sophia Clark';
    const initialEmail = 'sophia.clark@email.com';


    this.userEmail.set(initialEmail);
    
  
    this.profileForm.controls.name.setValue(initialName);
    
   
    this.profileForm.markAsPristine();
  }

  onSubmit() {
    if (this.profileForm.invalid) {
      console.error('Formulário inválido. Verifique o nome.');
      return;
    }

    const newName = this.profileForm.controls.name.value;
    console.log('Tentando salvar novo nome:', newName);

  
  }

  // Getter de conveniência
  get nameControl() {
    return this.profileForm.get('name') as FormControl;
  }
}
