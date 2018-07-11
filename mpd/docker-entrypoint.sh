#!/bin/sh

LOGPATH=${LOGPATH:-"/mpd/logs/log"}
[[ -p $LOGPATH ]] || exit 1

mkdir -p \
  /mpd/cache \
  /mpd/playlists

touch \
  /mpd/cache/tag_cache \
  /mpd/cache/state \
  /mpd/cache/sticker.sql

chown -R mpd \
  /mpd/cache \
  /mpd/playlists
  # $LOGPATH

## start
exec mpd \
  --no-daemon \
  /etc/mpd.conf "$@"
