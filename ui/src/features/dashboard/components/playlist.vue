<template lang="pug">
v-card.playlist
  v-card-title
    .title Playlist
  v-card-text
    virtual-list(:size="this.size" :remain="this.buffer" :onscroll="onscroll" :tobottom="tobottom")
      ul(v-for="(playlistitem, index) in playlistitems" :index="index" :key="playlistitem.Pos")
        v-flex(d-flex :style="style")
          v-layout(align-center justify-center row fill-height)
            v-flex(d-flex xs12 sm12 md4)
              | {{ playlistitem.Artist }}
            v-flex(d-flex xs12 sm12 md8)
              v-layout(style="align-items: center;")
                v-flex.text-xs-left(md9)
                  | {{ playlistitem.Title }}
                v-flex.text-xs-left(md1)
                  v-btn(flat icon color="primary")
                    v-icon play_arrow
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
      // px size of items
      size: 40,
      start: 0,
      end: 0,
      // preload item count
      buffer: 20,
      // initial load item count
      initialBuffer: 40,
      // save loaded state to refresh items
      bufferedStart: 0,
      bufferedEnd: 0
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
      if (this.end <= 0) {
        this.end = this.initialBuffer
      }
      console.info('playlistversion', this.start, this.end)
      this.$socket.sendObj({ mutation: 'playlistupdate', value: [this.start, this.end] })
    }
  },

  created () {
    this.$socket.sendObj({ mutation: 'playlistupdate', value: [0, this.initialBuffer] })
  },

  methods: {
    tobottom () {
      let end = this.end + (this.buffer * 4)
      this.$socket.sendObj({ mutation: 'playlistupdate', value: [this.end, end] })
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
