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
    )
      div(v-for="(playlistItem, index) in playlistItems" :index="index" :key="playlistItem.Id")
        draggable(v-model="playlistItems" @end="onMoved" :options="{group: 'playlistItems', handle: '.handle'}" :id="index")
          v-list-tile(@click="")
            template(v-if="$vuetify.breakpoint.smAndUp")
              v-list-tile-action.handle
                v-icon(color="primary") drag_handle
            v-list-tile-action
              v-btn(flat icon color="primary" @click="playId(playlistItem.Id)")
                v-icon play_arrow
            v-list-tile-title
              | {{ playlistItem.Artist || '...' }}
            v-list-tile-title
              | {{ playlistItem.Title || '...' }}
            template(v-if="$vuetify.breakpoint.smAndUp")
              v-list-tile-action
                v-btn(flat icon color="primary" @click="removeId(playlistItem.Id)")
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
      if (this.end <= 0) {
        this.end = this.buffer
      }
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

      for (i = this.start; i <= this.end; i++) {
        if (
          typeof (this.$store.state.websocket.socket.playlist[i]) === 'undefined' ||
          !('Pos' in this.$store.state.websocket.socket.playlist[i])
        ) {
          foundNull = true
          updateStart = i
          break
        }
      }

      if (foundNull) {
        for (i = this.end; i >= updateStart; i--) {
          if (
            typeof (this.$store.state.websocket.socket.playlist[i]) === 'undefined' ||
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
    onScroll: _.debounce(function (event, data) {
      this.start = data['start']
      this.end = data['end']

      this.updatePlaylist()
    }, 100)
  }
}
</script>

<style lang="stylus">
</style>
