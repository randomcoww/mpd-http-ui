<template lang="pug">
v-card
  v-card-title
    v-layout(row wrap style="align-items: center;")
      v-flex.text-xs-left(xs12 sm12 md4 title)
        | {{ currentsong.Artist }}
      v-flex.text-xs-right(xs12 sm12 md8 title)
        | {{ currentsong.Album }}/{{ currentsong.Title }}
  v-card-text
    v-slider(
      :max="seek_duration"
      :value="seek_elaspsed"
      v-on:mousedown="onmousedown"
      v-on:click="onmouseup"
      v-on:change="onchange")
    v-layout(row wrap style="align-items: center;")
      v-flex(d-flex xs3 sm2 md1)
        v-btn(flat icon color="primary")
          v-icon fast_rewind
      v-flex(d-flex xs3 sm2 md1)
        v-btn(flat icon color="primary")
          v-icon play_arrow
      v-flex(d-flex xs3 sm2 md1)
        v-btn(flat icon color="primary")
          v-icon pause
      v-flex(d-flex xs3 sm2 md1)
        v-btn(flat icon color="primary")
          v-icon stop
      v-flex(d-flex xs3 sm2 md1)
        v-btn(flat icon color="primary")
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
    onchange (value) {
      console.info('slider_changed', value)

      this.dragStartValue = null
      this.$socket.sendObj({ mutation: 'seek', value: value })
      this.$store.commit('elapsed', { value: value })
    },
    onmousedown () {
      console.info('drag_start')
      this.dragStartValue = this.$store.state.websocket.socket.elapsed
    },
    onmouseup () {
      console.info('drag_end')
      this.dragStartValue = null
    }
  }
}
</script>
