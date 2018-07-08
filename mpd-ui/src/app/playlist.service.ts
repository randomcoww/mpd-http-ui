import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { Observable, of } from 'rxjs';
import { catchError, map, tap } from 'rxjs/operators';

import { PlaylistItem } from './playlist-item';
import { MessageService } from './message.service';

@Injectable({
  providedIn: 'root'
})
export class PlaylistService {

  private playlistUrl = 'http://localhost:3000/playlist';

  constructor(
    private http: HttpClient,
    private messageService: MessageService) { }


  getPlaylistItems(start: number, end: number): Observable<PlaylistItem[]> {
    const url = `${this.playlistUrl}/items?start=${start}&end=${end}`;

    console.log(url);

    return this.http.get<PlaylistItem[]>(url)
      .pipe(
        tap(_ => this.log(`fetched playlist items between start=${start} and end=${end}`)),
        catchError(this.handleError<PlaylistItem[]>('getPlaylistItems', []))
      );
  }

  movePlaylistItems(start: number, end: number, pos: number): Observable<any> {
    const url = `${this.playlistUrl}/items?start=${start}&end=${end}&pos=${pos}`;
    return this.http.put(url, '{}')
      .pipe(
        tap(_ => this.log(`moved playlist items between start=${start} and end=${end} to pos=${pos}`)),
        catchError(this.handleError<any>('movePlaylistItems'))
      );
  }

  deletePlaylistItems(start: number, end: number): Observable<any> {
    const url = `${this.playlistUrl}/items?start=${start}&end=${end}`;
    return this.http.delete(url)
      .pipe(
        tap(_ => this.log(`deleted playlist items between start=${start} and end=${end}`)),
        catchError(this.handleError<any>('deletePlaylistItems'))
      );
  }

  addToPlaylist(path: string, position: number): Observable<any> {
    const url = `${this.playlistUrl}?path=${path}`;
    return this.http.put(url, '{}')
      .pipe(
        tap(_ => this.log(`add playlist items from path=${path}`)),
        catchError(this.handleError<any>('deletePlaylistItems', []))
      );
  }

  clearPlaylist(): Observable<any> {
    return this.http.delete(this.playlistUrl)
      .pipe(
        tap(_ => this.log(`clear playlist`)),
        catchError(this.handleError<any>('deletePlaylistItems', []))
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
