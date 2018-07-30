<template lang="pug">
v-card
  audio(
    src="/stream"
    autoplay="autoplay"
    ref="mpdplayer"
    preload="none"
    @canplay="onMpdEvent"
    @play="onMpdEvent"
    @playing="onMpdEvent"
    @emptied="onMpdEvent"
    @error="reloadAudio"
    @ratechange="onMpdEvent"
    @ended="reloadAudio"
    @stalled="reloadAudio"
    @suspended="onMpdEvent"
    @waiting="onMpdEvent")

  v-list(two-line subheader)
    v-list-tile(@click="")
      v-list-tile-avatar
        template(v-if="this.playerState >= 4")
          v-icon(color="primary lighten-1") play_arrow
        template(v-else)
          v-icon(color="primary lighten-1") pause
      v-list-tile-content
        v-list-tile-title
          | {{ currentSong.Artist || 'No Artist' }}/{{ currentSong.Title || 'No Title' }}
        v-list-tile-sub-title
          | {{ currentSong.Album || 'No Album' }}
        v-list-tile-sub-title
          | {{ currentSong.file }}

    v-list-tile(@click="")
      v-list-tile-avatar
      v-list-tile-content
        v-list-tile-title
          v-slider(
            :max="seekDuration"
            :value="seekElaspsed"
            v-on:mousedown="onSeekMouseDown"
            v-on:click.stop="onSeekMouseUp"
            v-on:change="onSeekChanged")
        v-list-tile-sub-title
          | {{ seekElaspsed | round }}/{{ seekDuration | round }}

</template>

<script>
import _ from 'lodash'

export default {
  filters: {
    round (v) {
      return Math.round(v)
    }
  },

  data () {
    return {
      dragStartValue: null,
      playerState: null
    }
  },

  computed: {
    currentSong () {
      return this.$store.state.websocket.socket.currentsong
    },
    seekElaspsed () {
      if (this.dragStartValue != null) {
        return this.dragStartValue
      } else {
        return this.$store.state.websocket.socket.elapsed
      }
    },
    seekDuration () {
      return this.$store.state.websocket.socket.duration
    }
  },

  watch: {
    currentSong: function () {
      this.reloadAudio()
    }
  },

  created () {
    this.$socket.sendObj({ mutation: 'currentsong' })
  },

  methods: {
    onMpdEvent: _.debounce(function (event) {
      this.playerState = this.$refs.mpdplayer.readyState
      // console.info('player state', event.type, this.playerState)
    }, 1000),

    reloadAudio: _.debounce(function () {
      // console.info('player reload')
      this.$refs.mpdplayer.load()
    }, 1000),

    onSeekChanged (value) {
      this.dragStartValue = null
      this.$socket.sendObj({ mutation: 'seek', value: value })
      this.$store.commit('elapsed', { value: value })
    },
    onSeekMouseDown () {
      this.dragStartValue = this.$store.state.websocket.socket.elapsed
    },
    onSeekMouseUp () {
      this.dragStartValue = null
    }
  }
}
</script>
