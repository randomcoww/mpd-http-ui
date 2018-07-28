<template lang="pug">
v-navigation-drawer.my-sidebar(
  v-model="isActive"
  fixed
  :mobile-break-point="1200"
  app
  right
  width="500"
)

  v-toolbar(dense flat)
    v-toolbar-title
      | Library
    v-spacer
    v-flex(xs8)
      v-text-field(append-icon="search" v-model="searchText" hide-details single-line v-model.lazy="databasequery")

  v-list
    virtual-list(
      :size="this.size"
      :remain="this.buffer"
      :onscroll="onscroll"
      :tobottom="tobottom"
    )
      div(v-for="(searchresult, index) in searchresults" :index="index" :key="searchresult.file")
        draggable(v-model="searchresults" @end="onmoved" :options="{group: 'playlistitems'}" :id="searchresult.file")
          v-list-tile(@click="")
            v-list-tile-action
              v-btn(flat icon color="primary" @click="addpath(searchresult.file, -1)")
                v-icon add
            v-list-tile-title
              | {{ searchresult.artist || 'No Artist' }}
            v-list-tile-title
              | {{ searchresult.title || 'No Title' }}
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
        return this.$store.state.common.sidebar.visible
      },
      set (val) {
        this.$store.dispatch('common/updateSidebar', { visible: val })
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
      this.requestStart = 0
      this.sendSearch(this.requestStart, this.requestCount)
      this.requestStart += this.requestCount
    }, 300)
  },

  mounted () {
    window.addEventListener('resize', this.onresize)
    this.onresize()
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.onresize)
  },

  methods: {
    onresize: _.debounce(function () {
      this.buffer = Math.floor((window.innerHeight - this.size - 10) / this.size)
    }, 300),

    sendSearch (start, count) {
      console.info('search', start, count)
      this.$socket.sendObj({ mutation: 'search', value: [this.databasequery, start, count] })
    },

    addpath (path, position) {
      console.info('addpath', path, position)

      position = parseInt(position)
      if (Number.isInteger(position)) {
        this.$socket.sendObj({ mutation: 'addpath', value: [path, position] })
      }
    },

    onmoved (event) {
      this.addpath(event.from.id, event.to.id)
    },

    tobottom () {
      if (this.end > this.requestStart - this.requestCount) {
        this.sendSearch(this.requestStart, this.requestCount)
        this.requestStart += this.requestCount
      }
    },

    onscroll (event, data) {
      this.end = data['end']
    }
  }
}
</script>

<style lang="stylus">
</style>
