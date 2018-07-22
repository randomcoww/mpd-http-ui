<template lang="pug">
v-card.playlist
  v-card-title
    .title Playlist
  v-card-text
    virtual-list(:size="64" :remain="this.buffer" :onscroll="onscroll" :tobottom="tobottom" :totop="totop")
      ul(v-for="(playlistitem, index) in playlistitems" :index="index" :key="playlistitem.Id")
        v-flex(d-flex)
          v-layout(align-center justify-center row fill-height)
            v-flex.text-xs-left(d-flex xs1 sm1 md1)
              | {{ playlistitem.Pos }}
            v-flex(d-flex xs12 sm12 md3)
              | {{ playlistitem.Artist }}
            v-flex(d-flex xs12 sm12 md8)
              v-layout(style="align-items: center;")
                v-flex.text-xs-left(md10)
                  | {{ playlistitem.Title }}
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
      bufferedEnd: 40
    }
  },

  computed: {
    playlistversion () {
      return this.$store.state.websocket.socket.version
    },
    playlistitems () {
      return this.$store.state.websocket.socket.playlist
    }
  },

  watch: {
    playlistversion: function () {
      this.$socket.sendObj({ mutation: 'playlistupdate', value: [this.start, this.end] })
    }
  },

  created () {
    this.$socket.sendObj({ mutation: 'playlistupdate', value: [0, this.end] })
  },

  methods: {
    tobottom () {
      let end = this.end + this.bench
      // console.info('request', this.end, end)
      this.$socket.sendObj({ mutation: 'playlistupdate', value: [this.end, end] })
    },

    totop () {
    },

    // playlist may have updated - refresh view as they become visible
    onscroll (event, data) {
      this.start = data['start']
      this.end = data['end']

      // scroll down refresh
      if (this.end > this.bufferedEnd) {
        let start = this.bufferedEnd
        if (start < this.start) {
          start = this.start
        }

        // console.info('request', start, this.end)
        this.$socket.sendObj({ mutation: 'playlistupdate', value: [start, this.end] })

        this.bufferedStart = this.start
        this.bufferedEnd = this.end
        return
      }

      // scroll up refresh
      if (this.start < this.bufferedStart) {
        let end = this.bufferedStart
        if (end > this.end) {
          end = this.end
        }

        // console.info('request', this.start, end)
        this.$socket.sendObj({ mutation: 'playlistupdate', value: [this.start, end] })

        this.bufferedEnd = this.end
        this.bufferedStart = this.start
      }
    }
  }
}
</script>

<style lang="stylus">
</style>
