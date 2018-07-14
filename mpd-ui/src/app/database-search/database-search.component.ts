import { Component, OnInit } from '@angular/core';
import { Observable, Subject } from 'rxjs';

import {
   debounceTime, distinctUntilChanged, switchMap
 } from 'rxjs/operators';

import { DatabaseItem } from '../database-item';
import { DatabaseService } from '../database.service';

@Component({
  selector: 'app-database-search',
  templateUrl: './database-search.component.html',
  styleUrls: ['./database-search.component.css']
})
export class DatabaseSearchComponent implements OnInit {
  databaseItems$: Observable<DatabaseItem[]>;
  private searchTerms = new Subject<string>();

  constructor(private databaseService: DatabaseService) { }

  // Push a search term into the observable stream.
  search(term: string): void {
    this.searchTerms.next(term);
  }

  ngOnInit(): void {
    this.databaseItems$ = this.searchTerms.pipe(
      // wait 300ms after each keystroke before considering the term
      debounceTime(300),

      // ignore new term if same as previous term
      distinctUntilChanged(),

      // switch to new search observable each time the term changes
      switchMap((term: string) => this.databaseService.searchDatabaseItems(term)),
    );
  }
}
