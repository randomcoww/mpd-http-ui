#!/bin/sh

mkdir -p /mpd/cache /mpd/playlists /mpd/log

touch \
  /mpd/cache/tag_cache \
  /mpd/cache/state \
  /mpd/cache/sticker.sql

chown -R mpd /mpd/cache /mpd/playlists /mpd/log

## start
exec mpd \
  --no-daemon \
  /etc/mpd.conf "$@"
