<template lang="pug">
v-card
  v-toolbar.my-appbar(
    app
    flat
    dense
  )
    v-btn(icon ripple @click="toggleLibrary")
      v-icon storage
    v-btn(icon ripple @click="togglePlaylist")
      v-icon playlist_play

    template(v-if="$vuetify.breakpoint.smAndUp")
      v-btn(icon ripple @click="playPrev")
        v-icon(color="primary lighten-1") fast_rewind
      v-btn(icon ripple @click="playId(-1)")
        v-icon(color="primary lighten-1") play_arrow
      v-btn(icon ripple @click="playNext")
        v-icon(color="primary lighten-1") fast_forward

      v-toolbar-title
        | {{ currentSong.Artist || 'No Artist' }}/{{ currentSong.Title || 'No Title' }}
    v-spacer

    template(v-if="$vuetify.breakpoint.smAndUp")
      v-btn(icon ripple @click="removeId(currentSong.Id)")
        v-icon(color="primary lighten-1") delete

    v-menu(bottom left)
      v-btn(icon slot="activator")
        v-icon more_vert
      v-list
        v-list-tile(@click="startDatabaseUpdate")
          v-list-tile-title
            | Start database update
        v-list-tile(@click="clearPlaylist")
          v-list-tile-title
            | Clear playlist

</template>

<script>
export default {

  computed: {
    currentSong () {
      return this.$store.state.websocket.socket.currentSong
    }
  },

  methods: {
    togglePlaylist () {
      this.$store.dispatch('common/togglePlaylist', { visible: !this.$store.state.common.playlist.visible })
    },

    toggleLibrary () {
      this.$store.dispatch('common/toggleLibrary', { visible: !this.$store.state.common.library.visible })
    },

    startDatabaseUpdate () {
      this.$socket.sendObj({ mutation: 'updatedb' })
    },

    clearPlaylist () {
      this.$socket.sendObj({ mutation: 'clear' })
    },

    playId (id) {
      this.$socket.sendObj({ mutation: 'playid', value: parseInt(id) })
    },

    removeId (id) {
      this.$socket.sendObj({ mutation: 'removeid', value: parseInt(id) })
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
