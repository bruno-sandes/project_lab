import { Routes } from '@angular/router';
import { AppMainLayout } from './shared/layout/app-main-layout/app-main-layout';
import { NotFoundComponent } from './shared/not-found-component/not-found-component';

export const routes: Routes = [
    {
        path: 'login',
        loadComponent: () => import('./features/login/component/login/login').then(m => m.Login)

    }, {
        path: 'register',
        loadComponent: () => import('./features/register/component/register/register').then(m => m.RegisterComponent)
    },
    {
        path: '',
        redirectTo: 'login',
        pathMatch: 'full'

    },
    {
        path: 'inicio',
        component: AppMainLayout,
        //to do colocar um authguard aqui que intercepte qualquer unauthorized da apI
        children: [
            {
                path: '',
                redirectTo: 'groups-dashboard',
                pathMatch: 'full'
            },
            {
                path: 'groups-dashboard',
                loadComponent: () => import('./features/groups-dashboard/components/groups-dashboard/groups-dashboard').then(m => m.GroupsDashboard)
            },
            {
                path: 'create-group',
                loadComponent: () => import('./features/groups-dashboard/components/create-travel-group/create-travel-group').then(m => m.CreateTravelGroup)
            },
            {
                path: 'profile',
                loadComponent:() => import('./features/profile/components/profile/profile').then(m => m.Profile)
            }
        ]
    }
];
