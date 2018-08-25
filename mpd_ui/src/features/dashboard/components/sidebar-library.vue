<template lang="pug">
v-navigation-drawer(
  v-model="isActive"
  app
  fixed
  left
  :width="1000"
  temporary=true
  v-resize="onResize"
)

  v-toolbar(dense flat)
    v-icon(color="grey") storage
    v-spacer
    v-flex(xs8)
      v-text-field(append-icon="search" v-model="searchText" hide-details single-line v-model.lazy="databaseQuery")
    v-btn(icon ripple @click="toggleLibrary")
      v-icon close

  v-list
    virtual-list(
      :size="this.size"
      :remain="this.buffer"
      :onscroll="onScroll"
      :tobottom="onScrollBottom"
      :debounce="50"
      :bench="this.buffer"
    )
      div(v-for="(searchResult, index) in searchResults" :index="index" :key="searchResult.file")
        v-list-tile(@click="")
          v-list-tile-title
            | {{ searchResult.artist || 'No Artist' }}
          v-list-tile-title
            | {{ searchResult.title || 'No Title' }}
          v-list-tile-title
            | {{ searchResult.album || 'No Album' }}
</template>

<script>
import VirtualList from 'vue-virtual-scroll-list'
import _ from 'lodash'
// import draggable from 'vuedraggable'

export default {
  components: {
    // draggable,
    VirtualList
  },

  data () {
    return {
      // px size of items
      size: 48,
      end: 0,
      // preload item count
      buffer: 10,
      // initial load item count
      requestStart: 0,
      requestCount: 40,
      // search query
      databaseQuery: null
    }
  },

  computed: {
    isActive: {
      get () {
        return this.$store.state.common.library.visible
      },
      set (val) {
        this.$store.dispatch('common/toggleLibrary', { visible: val })
      }
    },

    searchResults: {
      get: function () {
        return this.$store.state.websocket.socket.search
      },
      set: function () {
      }
    }
  },

  watch: {
    databaseQuery: _.debounce(function () {
      let requestCount = this.buffer * 2

      this.requestStart = 0
      this.sendSearch(this.requestStart, requestCount)
      this.requestStart += requestCount
    }, 300)
  },

  mounted () {
    this.onResize()
  },

  methods: {
    toggleLibrary () {
      this.$store.dispatch('common/toggleLibrary', { visible: !this.$store.state.common.library.visible })
    },

    onResize () {
      this.buffer = Math.floor((window.innerHeight - this.size - 10) / this.size)
    },

    sendSearch (start, count) {
      // console.info('search', start, count)
      this.$socket.sendObj({ mutation: 'search', value: [this.databaseQuery, start, count] })
    },

    showSnackMessage (msg) {
      this.$store.dispatch('common/updateSnackbar', { show: true, text: msg })
    },

    addPathToPlaylist (path, position) {
      position = parseInt(position)
      if (Number.isInteger(position)) {
        // console.info('AddToPlaylist', path, position)
        this.$socket.sendObj({ mutation: 'addpath', value: [path, position] })
        // show added message
        this.showSnackMessage('Added ' + path)
      }
    },

    onDragStart (event) {
      // console.info('onDragStart', event)
    },

    onDragEnd (event) {
      // console.info('onDragEnd', event.from.id, event.to.id)
      this.addPathToPlaylist(event.from.id, event.to.id)
    },

    onScrollBottom () {
      if (this.end > this.requestStart - this.requestCount) {
        this.sendSearch(this.requestStart, this.requestCount)
        this.requestStart += this.requestCount
      }
    },

    onMoved (event) {
      this.addPathToPlaylist(event.from.id, event.to.id)
    },

    onScroll (event, data) {
      this.end = data['end']
    }
  }
}
</script>

<style lang="stylus">
</style>
