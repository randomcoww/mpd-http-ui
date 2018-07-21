<template lang="pug">
v-card-text
  v-layout(row wrap style="align-items: center;")
    v-flex(d-flex xs12 sm12 md2)
      | Artist
    v-flex(d-flex xs12 sm12 md10)
      v-layout(style="align-items: center;")
        v-flex.text-xs-left(md10) {{ currentsong.Artist || initialsong.Artist }}
        v-flex.text-xs-right(md2)
          v-btn(flat icon color="primary")
            v-icon delete
    v-flex(d-flex xs12 sm12 md2)
      | Title
    v-flex(d-flex xs12 sm12 md10)
      v-layout(style="align-items: center;")
        v-flex.text-xs-left(md10) {{ currentsong.Title || initialsong.Title }}
        v-flex.text-xs-right(md2)
          v-btn(flat icon color="primary")
            v-icon delete
    v-flex(d-flex xs12 sm12 md2)
      | Album
    v-flex(d-flex xs12 sm12 md10)
      v-layout(style="align-items: center;")
        v-flex.text-xs-left(md10) {{ currentsong.Album || initialsong.Album }}
        v-flex.text-xs-right(md2)
          v-btn(flat icon color="primary")
            v-icon delete
    v-flex(d-flex xs12 sm12 md2)
      | File
    v-flex(d-flex xs12 sm12 md10)
      v-layout(style="align-items: center;")
        v-flex.text-xs-left(md10) {{ currentsong.file || initialsong.file }}
        v-flex.text-xs-right(md2)
          v-btn(flat icon color="primary")
            v-icon delete
</template>

<script>
export default {
  data () {
    return {
      initialsong: {},
      errored: false
    }
  },

  computed: {
    currentsong () {
      return this.$store.state.websocket.socket.currentsong
    }
  },

  created () {
    this.$socket.sendObj({ mutation: 'currentsong' })
  }
}
</script>
