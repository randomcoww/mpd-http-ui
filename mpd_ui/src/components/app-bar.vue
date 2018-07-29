<template lang="pug">
v-card
  v-toolbar.my-appbar(
    app
    flat
    dense
  )
    v-btn(icon ripple @click="playPrev")
      v-icon(color="primary lighten-1") fast_rewind
    v-btn(icon ripple @click="playId(-1)")
      v-icon(color="primary lighten-1") play_arrow
    v-btn(icon ripple @click="playNext")
      v-icon(color="primary lighten-1") fast_forward

    v-toolbar-title
      | {{ currentsong.Artist || 'No Artist' }}/{{ currentsong.Title || 'No Title' }}
    v-spacer
    v-btn(icon ripple @click="clearPlaylist")
      v-icon(color="primary lighten-1") delete
    v-toolbar-side-icon(@click.stop="toggleSidebar()")

</template>

<script>
export default {

  computed: {
    currentsong () {
      return this.$store.state.websocket.socket.currentsong
    }
  },

  methods: {
    toggleSidebar () {
      this.$store.dispatch('common/updateSidebar', { visible: !this.$store.state.common.sidebar.visible })
    },

    clearPlaylist () {
      this.$socket.sendObj({ mutation: 'clear' })
    },

    playId (id) {
      this.$socket.sendObj({ mutation: 'playid', value: parseInt(id) })
    },

    playNext () {
      this.$socket.sendObj({ mutation: 'playnext' })
    },

    playPrev () {
      this.$socket.sendObj({ mutation: 'playprev' })
    }
  }
}
</script>

<style lang="stylus">
</style>
