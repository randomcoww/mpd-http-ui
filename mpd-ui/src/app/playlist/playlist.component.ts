import { Component, OnInit } from '@angular/core';

import { PlaylistItem } from '../playlist-item';
import { PlaylistService } from '../playlist.service';

@Component({
  selector: 'app-playlist',
  templateUrl: './playlist.component.html',
  styleUrls: ['./playlist.component.css']
})
export class PlaylistComponent implements OnInit {
  playlistItems: PlaylistItem[];

  constructor(private playlistService: PlaylistService) { }

  ngOnInit() {
    this.getPlaylistItems();
  }

  getPlaylistItems(): void {
    this.playlistService.getPlaylistItems(-1, -1)
    .subscribe(playlistItems => this.playlistItems = playlistItems);
  }
}
