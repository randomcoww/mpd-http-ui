<template lang="pug">
v-card
  v-card-title
    v-layout(row wrap style="align-items: center;")
      v-flex.text-xs-left(xs12 sm12 md3 title)
        | {{ currentsong.Artist || 'No Artist' }}
      v-flex.text-xs-right(xs12 sm12 md9 title)
        | {{ currentsong.Album || 'No Album' }}/{{ currentsong.Title || 'No Title' }}
      v-flex.text-xs-left(xs12 sm12 md12)
        | {{ currentsong.file }}
  v-card-text
    v-slider(
      :max="seek_duration"
      :value="seek_elaspsed"
      v-on:mousedown="onmousedown"
      v-on:click="onmouseup"
      v-on:change="onchange")
    v-layout(row wrap style="align-items: center;")
      v-flex(d-flex xs3 sm2 md1)
        v-btn(flat icon color="primary" @click="playprev")
          v-icon fast_rewind
      v-flex(d-flex xs3 sm2 md1)
        v-btn(flat icon color="primary" @click="playid(-1)")
          v-icon play_arrow
      v-flex(d-flex xs3 sm2 md1)
        v-btn(flat icon color="primary" @click="pause")
          v-icon pause
      v-flex(d-flex xs3 sm2 md1)
        v-btn(flat icon color="primary" @click="stop")
          v-icon stop
      v-flex(d-flex xs3 sm2 md1)
        v-btn(flat icon color="primary" @click="playnext")
          v-icon fast_forward
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

  created () {
    this.$socket.sendObj({ mutation: 'currentsong' })
  },

  methods: {
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
