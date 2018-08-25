<template lang="pug">
v-navigation-drawer(
  v-model="isActive"
  app
  fixed
  right
  :width="1000"
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
      :debounce="50"
      :bench="this.buffer"
    )
      div(v-for="(playlistItem, index) in playlistItems" :index="index" :key="playlistItem.Id")
        v-list-tile(@click="")
          v-list-tile-title
            | {{ playlistItem.Artist || 'No Artist' }}
          v-list-tile-title
            | {{ playlistItem.Title || 'No Artist' }}
          v-list-tile-title
            | {{ playlistItem.Album || 'No Album' }}
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
      buffer: 0,
      // save loaded state to refresh items
      drag: false
    }
  },

  computed: {
    socketReady: _.debounce(function () {
      if (this.$store.state.websocket.socket.isConnected) {
        this.$socket.sendObj({ mutation: 'playlistlengthquery' })
      }
    }, 300),

    isActive: {
      get () {
        return this.$store.state.common.playlist.visible
      },
      set (val) {
        this.$store.dispatch('common/togglePlaylist', { visible: val })
      }
    },

    playlistItems: {
      get: function () {
        return this.$store.state.websocket.socket.playlist
      },
      set: function () {
      }
    }
  },

  watch: {
    playlistItems: function () {
      this.updatePlaylist()
    },
    socketReady: function () {
    }
  },

  mounted () {
    this.onResize()
  },

  methods: {
    togglePlaylist () {
      this.$store.dispatch('common/togglePlaylist', { visible: !this.$store.state.common.playlist.visible })
    },

    onResize () {
      this.buffer = Math.floor((window.innerHeight - this.size - 10) / this.size)
      if (this.end <= this.start) {
        this.end = this.start + this.buffer
      }
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

    updatePlaylist () {
      let i
      let foundNull = false
      let updateStart = this.start
      let updateEnd = this.end

      for (i = updateStart; i <= updateEnd; i++) {
        if (
          !('Pos' in this.$store.state.websocket.socket.playlist[i])
        ) {
          foundNull = true
          updateStart = i
          break
        }
      }

      if (foundNull) {
        for (i = updateEnd; i >= updateStart; i--) {
          if (
            !('Pos' in this.$store.state.websocket.socket.playlist[i])
          ) {
            updateEnd = i
            break
          }
        }

        this.$socket.sendObj({ mutation: 'playlistquery', value: [updateStart, updateEnd + 1] })
      }
    },

    // playlist may have updated - refresh view as they become visible
    onScroll (event, data) {
      this.start = data['start']
      this.end = data['end']

      this.updatePlaylist()
    }
  }
}
</script>

<style lang="stylus">
</style>
