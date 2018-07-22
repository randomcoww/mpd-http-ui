<template lang="pug">
v-card.playlist
  v-card-title
    .title Playlist
  v-card-text
    virtual-list(:size="64" :remain="this.buffer" :bench="this.bench" :onscroll="onscroll" :tobottom="tobottom" :totop="totop")
      ul(v-for="playlistitem in playlistitems")
        v-flex(d-flex)
          v-layout(align-center justify-center row fill-height)
            v-flex(d-flex xs12 sm12 md4)
              | {{ playlistitem.Artist }}
            v-flex(d-flex xs12 sm12 md8)
              v-layout(style="align-items: center;")
                v-flex.text-xs-left(md10) {{ playlistitem.Title }}
                v-flex.text-xs-left(md2)
                  v-btn(flat icon color="primary")
                    v-icon delete
</template>

<script>
import VirtualList from 'vue-virtual-scroll-list'

export default {
  components: {
    VirtualList
  },

  data () {
    return {
      bench: 100,
      start: 0,
      end: 40,
      buffer: 12,
      bufferedStart: 0,
      bufferedEnd: 0
    }
  },

  computed: {
    playlistitems () {
      return this.$store.state.websocket.socket.playlist
    }
  },

  created () {
    let end = this.end + this.bench
    this.$socket.sendObj({ mutation: 'playlist', value: [this.start, end] })
  },

  methods: {
    tobottom () {
      let end = this.end + this.bench
      console.info('request', this.end, end)
      this.$socket.sendObj({ mutation: 'playlist', value: [this.end, end] })
    },

    totop () {
      let start = this.start - this.bench
      if (start < 0) {
        start = 0
      }
      console.info('request', start, this.start)
      this.$socket.sendObj({ mutation: 'playlist', value: [start, this.start] })
    },

    onscroll (event, data) {
      this.start = data['start']
      this.end = data['end']

      if (
        this.start > this.bufferedEnd - this.buffer ||
        this.end < this.bufferedStart + this.buffer
      ) {
        this.bufferedStart = this.start
        this.bufferedEnd = this.end

        console.info('request', this.start, this.end)
        this.$socket.sendObj({ mutation: 'playlist', value: [this.start, this.end] })
      }
    }
  }
}
</script>

<style lang="stylus">

</style>
