<template lang="pug">
v-navigation-drawer(
  v-model="isActive"
  app
  fixed
  left
  :width="800"
  temporary=true
  v-resize="onResize"
)

  v-toolbar(dense flat)
    v-icon(color="grey") storage
    v-spacer
    v-flex(xs8)
      v-text-field(append-icon="search" v-model="searchText" hide-details single-line v-model.lazy="databasequery")
    v-btn(icon ripple @click="toggleLibrary")
      v-icon close

  v-list
    virtual-list(
      :size="this.size"
      :remain="this.buffer"
      :onscroll="onScroll"
      :tobottom="onScrollBottom"
    )
      div(v-for="(searchresult, index) in searchresults" :index="index" :key="searchresult.file")
        draggable(v-model="searchresults" @end="onMoved" :options="{group: 'playlistitems'}" :id="searchresult.file")
          v-list-tile(@click="")
            v-list-tile-action
              v-btn(flat icon color="primary" @click="addPathToPlaylist(searchresult.file, -1)")
                v-icon add
            v-list-tile-title
              | {{ searchresult.artist || 'No Artist' }}
            v-list-tile-title
              | {{ searchresult.title || 'No Title' }}
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
      databasequery: null
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

    searchresults: {
      get: function () {
        return this.$store.state.websocket.socket.search
      },
      set: function () {
      }
    },
    style () {
      return {
        // 'height': this.size + 'px'
      }
    }
  },

  watch: {
    databasequery: _.debounce(function () {
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
      this.$socket.sendObj({ mutation: 'search', value: [this.databasequery, start, count] })
    },

    showSnackMessage (msg) {
      this.$store.dispatch('common/updateSnackbar', { show: true, text: msg })
    },

    addPathToPlaylist (path, position) {
      position = parseInt(position)
      if (Number.isInteger(position)) {
        console.info('AddToPlaylist', path, position)
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
