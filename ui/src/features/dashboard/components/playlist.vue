<template lang="pug">
v-card-text
  v-layout(row wrap style="align-items: center;")
    template(v-for="playlistitem in playlistitems || initialplaylistitems")
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
import axios from 'axios'

export default {
  data () {
    return {
      initialplaylistitems: [],
      errored: false
    }
  },

  computed: {
    playlistitems () {
      return this.$store.state.websocket.socket.playlist
    }
  },

  mounted () {
    axios.get('http://localhost:3000/playlist/items', {
      params: {
        start: -1,
        end: -1
      }
    }).then(response => {
      this.initialplaylistitems = response.data
    }).catch(error => {
      console.log(error)
      this.errored = true
    })
  }
}
</script>
