<template lang="pug">
v-card.playlist
  v-card-title
    .title Playlist
  v-card-text
    virtual-list(:size="this.size" :remain="this.buffer" :onscroll="onscroll" :tobottom="tobottom" :totop="totop")
      ul(v-for="(playlistitem, index) in playlistitems" :index="index" :key="playlistitem.Pos")
        v-flex(d-flex :style="style")
          v-layout(align-center justify-center row fill-height)
            v-flex.text-xs-left(d-flex xs1 sm1 md1)
              | {{ index }}
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
import _ from 'lodash'

export default {
  components: {
    VirtualList
  },

  data () {
    return {
      size: 40,
      start: 0,
      end: 40,
      buffer: 20,
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
    },
    style () {
      return {
        'height': this.size + 'px'
      }
    }
  },

  watch: {
    playlistversion: function () {
      if (this.end < 0) {
        this.end = 40
        this.bufferedEnd = 40
      }
      console.info('playlistversion', this.start, this.end)
      this.$socket.sendObj({ mutation: 'playlistupdate', value: [this.start, this.end] })
    }
  },

  created () {
    this.$socket.sendObj({ mutation: 'playlistupdate', value: [0, this.end] })
  },

  methods: {
    tobottom () {
      let end = this.end + (this.buffer * 4)
      this.$socket.sendObj({ mutation: 'playlistupdate', value: [this.end, end] })
    },

    totop () {
    },

    // playlist may have updated - refresh view as they become visible
    onscroll: _.debounce(function (event, data) {
      this.start = data['start']
      this.end = data['end']

      // scroll down refresh
      if (this.end > this.bufferedEnd) {
        let start = this.bufferedEnd
        if (start < this.start) {
          start = this.start
        }

        let end = this.end + this.buffer

        console.info('down_request', start, end)
        this.$socket.sendObj({ mutation: 'playlistupdate', value: [start, end] })

        this.bufferedStart = this.start
        this.bufferedEnd = end
        return
      }

      // scroll up refresh
      if (this.start < this.bufferedStart) {
        // buffer start to avoid making many requests
        let start = this.start - this.buffer
        if (start < 0) {
          start = 0
        }

        let end = this.bufferedStart
        if (end > this.end) {
          end = this.end
        }

        console.info('up_request', start, end)
        this.$socket.sendObj({ mutation: 'playlistupdate', value: [start, end] })

        this.bufferedStart = start
        this.bufferedEnd = this.end
      }
    }, 300)
  }
}
</script>

<style lang="stylus">
</style>
