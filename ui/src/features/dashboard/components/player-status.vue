<template lang="pug">
v-card-text
  v-slider(
    :max="seek_duration"
    :value="seek_elaspsed"
    v-on:mousedown="onmousedown"
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

  methods: {
    onchange (value) {
      console.info('slider_changed', value)

      this.dragStartValue = null
      this.$socket.sendObj({ mutation: 'seek', value: value })
      this.$store.commit('elapsed', { value: value })
    },
    onmousedown () {
      console.info('dragging')
      this.dragStartValue = this.$store.state.websocket.socket.elapsed
    }
  }
}
</script>
