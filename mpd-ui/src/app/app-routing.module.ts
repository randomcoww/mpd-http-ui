import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PlaylistComponent } from './playlist/playlist.component';
import { DatabaseSearchComponent } from './database-search/database-search.component';

const routes: Routes = [
  // { path: '', redirectTo: '/playlist', pathMatch: 'full' },
  { path: 'playlist', component: PlaylistComponent },
  { path: 'database-search', component: DatabaseSearchComponent }
];

@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRoutingModule { }
