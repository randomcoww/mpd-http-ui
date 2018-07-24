<template lang="pug">
v-card
  audio(
    src="http://localhost:8000/mpd"
    autoplay="autoplay"
    ref="mpdplayer"
    preload="none"
    @end="reloadmpd"
    @error="reloadmpd"
    @ratechange="reloadmpd")

  v-toolbar(dark)
    v-toolbar-side-icon
    v-toolbar-title
      | {{ currentsong.Artist || 'No Artist' }}/{{ currentsong.Title || 'No Title' }}
    v-spacer
    v-btn(icon ripple @click="playprev")
      v-icon fast_rewind
    v-btn(icon ripple @click="playid(-1)")
      v-icon play_arrow
    v-btn(icon ripple @click="pause")
      v-icon pause
    v-btn(icon ripple @click="stop")
      v-icon stop
    v-btn(icon ripple @click="playnext")
      v-icon fast_forward

  v-list(two-line subheader)
    v-list-tile(@click="")
      v-list-tile-avatar
        v-icon(color="primary lighten-1") play_arrow
      v-list-tile-content
        v-list-tile-title
          | {{ currentsong.Artist || 'No Artist' }}/{{ currentsong.Title || 'No Title' }}
        v-list-tile-sub-title
          | {{ currentsong.Album || 'No Album' }}
    v-subheader(inset)
      | {{ currentsong.file }}

    v-list-tile(@click="")
      v-slider(
        :max="seek_duration"
        :value="seek_elaspsed"
        v-on:mousedown="onmousedown"
        v-on:click="onmouseup"
        v-on:change="onchange")

</template>

<script>
export default {
  data () {
    return {
      dragStartValue: null
    }
  },

  computed: {
    currentsong () {
      return this.$store.state.websocket.socket.currentsong
    },
    seek_elaspsed () {
      if (this.dragStartValue != null) {
        return this.dragStartValue
      } else {
        return this.$store.state.websocket.socket.elapsed
      }
    },
    seek_duration () {
      return this.$store.state.websocket.socket.duration
    }
  },

  watch: {
    currentsong: function () {
      this.reloadmpd()
    }
  },

  created () {
    this.$socket.sendObj({ mutation: 'currentsong' })
  },

  methods: {
    reloadmpd () {
      console.info('reload mpd')
      this.$refs.mpdplayer.load()
    },

    playid (id) {
      this.$socket.sendObj({ mutation: 'playid', value: parseInt(id) })
    },

    stop () {
      this.$socket.sendObj({ mutation: 'stop' })
    },

    pause () {
      this.$socket.sendObj({ mutation: 'pause' })
    },

    playnext () {
      this.$socket.sendObj({ mutation: 'playnext' })
    },

    playprev () {
      this.$socket.sendObj({ mutation: 'playprev' })
    },

    onchange (value) {
      this.dragStartValue = null
      this.$socket.sendObj({ mutation: 'seek', value: value })
      this.$store.commit('elapsed', { value: value })
    },
    onmousedown () {
      this.dragStartValue = this.$store.state.websocket.socket.elapsed
    },
    onmouseup () {
      this.dragStartValue = null
    }
  }
}
</script>
