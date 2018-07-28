<template lang="pug">
v-card
  v-toolbar.my-appbar(
    app
    flat
    dense
  )
    v-btn(icon ripple @click="playprev")
      v-icon fast_rewind
    v-btn(icon ripple @click="playid(-1)")
      v-icon play_arrow
    v-btn(icon ripple @click="playnext")
      v-icon fast_forward

    v-toolbar-title
      | {{ currentsong.Artist || 'No Artist' }}/{{ currentsong.Title || 'No Title' }}
    v-spacer
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

    playid (id) {
      this.$socket.sendObj({ mutation: 'playid', value: parseInt(id) })
    },

    playnext () {
      this.$socket.sendObj({ mutation: 'playnext' })
    },

    playprev () {
      this.$socket.sendObj({ mutation: 'playprev' })
    }
  }
}
</script>

<style lang="stylus">
</style>
