import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { Observable, of } from 'rxjs';
import { catchError, map, tap } from 'rxjs/operators';

import { DatabaseItem } from './database-item';
import { MessageService } from './message.service';

@Injectable({
  providedIn: 'root'
})
export class DatabaseService {

  private databaseUrl = 'http://localhost:3000/database';

  constructor(
    private http: HttpClient,
    private messageService: MessageService) { }

  searchDatabaseItems(query: string): Observable<DatabaseItem[]> {
    const url = `${this.databaseUrl}/search?q=${query}`;
    return this.http.get<DatabaseItem[]>(url)
      .pipe(
        tap(_ => this.log(`search database items by q=${query}`)),
        catchError(this.handleError<DatabaseItem[]>('searchDatabaseItems', []))
      );
  }


  private handleError<T> (operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {

      // TODO: send the error to remote logging infrastructure
      console.error(error); // log to console instead

      // TODO: better job of transforming error for user consumption
      this.log(`${operation} failed: ${error.message}`);

      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }

  private log(message: string) {
    this.messageService.add('PlaylistService: ' + message);
  }
}
