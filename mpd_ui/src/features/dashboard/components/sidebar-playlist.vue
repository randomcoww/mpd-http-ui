<template lang="pug">
v-navigation-drawer(
  v-model="isActive"
  app
  fixed
  right
  :width="800"
  v-resize="onResize"
)

  v-toolbar(dense flat)
    v-icon(color="grey") playlist_play
    v-spacer
    v-btn(icon ripple @click="togglePlaylist")
      v-icon close

  v-list
    virtual-list(
      :size="this.size"
      :remain="this.buffer"
      :onscroll="onScroll"
      :tobottom="onScrollBottom"
    )
      div(v-for="(playlistitem, index) in playlistitems" :index="index" :key="playlistitem.Id")
        draggable(v-model="playlistitems" @end="onMoved" :options="{group: 'playlistitems', handle: '.handle'}" :id="index")
          v-list-tile(@click="")
            template(v-if="$vuetify.breakpoint.smAndUp")
              v-list-tile-action.handle
                v-icon(color="primary") drag_handle
            v-list-tile-action
              v-btn(flat icon color="primary" @click="playId(playlistitem.Id)")
                v-icon play_arrow
            v-list-tile-title
              | {{ playlistitem.Artist || 'No Artist' }}
            v-list-tile-title
              | {{ playlistitem.Title || 'No Title' }}
            template(v-if="$vuetify.breakpoint.smAndUp")
              v-list-tile-action
                v-btn(flat icon color="primary" @click="removeId(playlistitem.Id)")
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
      buffer: 10,
      // save loaded state to refresh items
      bufferedStart: 0,
      bufferedEnd: 0,
      drag: false
    }
  },

  computed: {
    isActive: {
      get () {
        return this.$store.state.common.playlist.visible
      },
      set (val) {
        this.$store.dispatch('common/togglePlaylist', { visible: val })
      }
    },

    playlistitems: {
      get: function () {
        return this.$store.state.websocket.socket.playlist
      },
      set: function () {
      }
    }
    // playlistVersion () {
    //   return this.$store.state.websocket.socket.version
    // }
  },

  watch: {
    playlistitems: function () {
      if (this.end <= 0) {
        this.end = this.buffer
      }
      // console.info('playlistversion', this.start, this.end)
      this.$socket.sendObj({ mutation: 'playlistupdate', value: [this.start, this.end] })
    }
  },

  mounted () {
    this.onResize()
    this.$socket.sendObj({ mutation: 'playlistupdate', value: [0, this.buffer * 2] })
  },

  methods: {
    togglePlaylist () {
      this.$store.dispatch('common/togglePlaylist', { visible: !this.$store.state.common.playlist.visible })
    },

    onResize () {
      this.buffer = Math.floor((window.innerHeight - this.size - 10) / this.size)
    },

    playId (id) {
      this.$socket.sendObj({ mutation: 'playid', value: parseInt(id) })
    },

    removeId (id) {
      this.$socket.sendObj({ mutation: 'removeid', value: parseInt(id) })
    },

    onMoved (event, data) {
      // console.info(event)
      let from = parseInt(event.from.id)
      let to = parseInt(event.to.id)

      // console.info('moveitem', from, to)
      this.$socket.sendObj({ mutation: 'playlistmove', value: [from, from + 1, to] })
    },

    onScrollBottom () {
      let end = this.end + this.buffer
      this.$socket.sendObj({ mutation: 'playlistupdate', value: [this.end, end] })
    },

    // playlist may have updated - refresh view as they become visible
    onScroll: _.debounce(function (event, data) {
      this.start = data['start']
      this.end = data['end']

      var i
      for (i = this.start; i <= this.end; i++) {
        if (this.$store.state.websocket.socket.playlist[i] === null) {
          this.$socket.sendObj({ mutation: 'playlistupdate', value: [this.start, this.end] })
          return
        }
      }

      // scroll down refresh
      // if (this.end > this.bufferedEnd) {
      //   let start = this.bufferedEnd
      //   if (start < this.start) {
      //     start = this.start
      //   }

      //   let end = this.end + this.buffer

      //   // console.info('down_request', start, end)
      //   this.$socket.sendObj({ mutation: 'playlistupdate', value: [start, end] })

      //   this.bufferedStart = this.start
      //   this.bufferedEnd = end
      //   return
      // }

      // // scroll up refresh
      // if (this.start < this.bufferedStart) {
      //   // buffer start to avoid making many requests
      //   let start = this.start - this.buffer
      //   if (start < 0) {
      //     start = 0
      //   }

      //   let end = this.bufferedStart
      //   if (end > this.end) {
      //     end = this.end
      //   }

      //   // console.info('up_request', start, end)
      //   this.$socket.sendObj({ mutation: 'playlistupdate', value: [start, end] })

      //   this.bufferedStart = start
      //   this.bufferedEnd = this.end
      // }
    }, 300)
  }
}
</script>

<style lang="stylus">
</style>
