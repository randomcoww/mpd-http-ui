<template lang="pug">
v-card-text
  v-layout(row wrap style="align-items: center;")
    template(v-for="playlistitem in playlistitems")
      v-flex(d-flex xs12 sm12 md4)
        | {{ playlistitem.Artist }}
      v-flex(d-flex xs12 sm12 md8)
        v-layout(style="align-items: center;")
          v-flex.text-xs-left(md10) {{ playlistitem.Title }}
          v-flex.text-xs-right(md2)
            v-btn(flat icon color="primary")
              v-icon delete
</template>

<script>
export default {
  computed: {
    playlistitems () {
      return this.$store.state.websocket.socket.playlist
    }
  },

  created () {
    this.$socket.sendObj({ mutation: 'playlist', value: [-1, -1] })
  }
}
</script>
