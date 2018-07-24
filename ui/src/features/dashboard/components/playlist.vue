<template lang="pug">
v-card.playlist
  v-toolbar(dark)
    v-toolbar-side-icon
    v-toolbar-title
      | Playlist

  v-list
    virtual-list(:size="this.size" :remain="this.buffer" :onscroll="onscroll" :tobottom="tobottom")
      div(v-for="(playlistitem, index) in playlistitems" :index="index" :key="playlistitem.Pos")
        draggable(v-model="playlistitems" @end="onmoved" :options="{group: 'playlistitems'}" :id="index")
          v-list-tile(@click="")
            v-list-tile-title
              | {{ playlistitem.Artist || 'No Artist' }}
            v-list-tile-title
              | {{ playlistitem.Title || 'No Title' }}
            v-list-action
              v-btn(flat icon color="primary" @click="playid(playlistitem.Id)")
                v-icon play_arrow
            v-list-action
              v-btn(flat icon color="primary" @click="removeid(playlistitem.Id)")
                v-icon delete
</template>

<script>
import VirtualList from 'vue-virtual-scroll-list'
import _ from 'lodash'
import draggable from 'vuedraggable'

export default {
  components: {
    VirtualList,
    draggable
  },

  data () {
    return {
      // px size of items
      size: 48,
      start: 0,
      end: 0,
      // preload item count
      buffer: 17,
      // initial load item count
      initialBuffer: 40,
      // save loaded state to refresh items
      bufferedStart: 0,
      bufferedEnd: 0,
      drag: false
    }
  },

  computed: {
    playlistitems: {
      get: function () {
        return this.$store.state.websocket.socket.playlist
      },
      set: function () {
      }
    },
    playlistversion () {
      return this.$store.state.websocket.socket.version
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
    playid (id) {
      this.$socket.sendObj({ mutation: 'playid', value: parseInt(id) })
    },

    removeid (id) {
      this.$socket.sendObj({ mutation: 'removeid', value: parseInt(id) })
    },

    onmoved (event, data) {
      // console.info(event)
      let from = parseInt(event.from.id)
      let to = parseInt(event.to.id)

      console.info('moveitem', from, to)
      this.$socket.sendObj({ mutation: 'playlistmove', value: [from, from + 1, to] })
    },

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
