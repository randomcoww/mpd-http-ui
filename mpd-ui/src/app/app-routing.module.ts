import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PlaylistComponent } from './playlist/playlist.component';

const routes: Routes = [
  { path: '', redirectTo: '/playlist', pathMatch: 'full' },
  { path: 'playlist', component: PlaylistComponent }
];

@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRoutingModule { }
