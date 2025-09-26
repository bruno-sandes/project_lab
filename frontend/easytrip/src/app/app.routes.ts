import { Routes } from '@angular/router';

export const routes: Routes = [
    {
        path: '',
        loadComponent: () => import('./features/login/component/login/login').then(m => m.Login)
    },{
        path: 'register',
        loadComponent: () => import('./features/register/component/register/register').then(m => m.RegisterComponent)
    }
];
