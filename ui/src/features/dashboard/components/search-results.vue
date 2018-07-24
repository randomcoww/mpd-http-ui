<template lang="pug">
v-card.searchresults
  v-card-title
    .title Library
  v-card-text
    v-text-field(label="Search..." v-model.lazy="databasequery")
    virtual-list(:size="this.size" :remain="this.buffer" :onscroll="onscroll" :tobottom="tobottom")
      div(v-for="(searchresult, index) in searchresults" :index="index" :key="searchresult.file")
        draggable(v-model="searchresults" @end="onmoved" :options="{group: 'playlistitems'}" :id="searchresult.file")
          v-flex(d-flex :style="style")
            v-container
              v-layout(row wrap style="align-items: center;")

                v-flex(d-flex xs10 sm10 md10)
                  v-flex(d-flex xs12 sm12 md4)
                    | {{ searchresult.artist || 'No Artist' }}
                  v-flex.text-xs-left(md8)
                    | {{ searchresult.title || 'No Title' }}

                v-flex(d-flex xs2 sm2 md2)
                  v-layout(style="align-items: center;")
                    v-btn(flat icon color="primary" @click="addpath(searchresult.file, -1)")
                      v-icon add
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
      size: 40,
      end: 0,
      // preload item count
      buffer: 25,
      // initial load item count
      requestStart: 0,
      requestCount: 40,
      // search query
      databasequery: null
    }
  },

  computed: {
    searchresults: {
      get: function () {
        return this.$store.state.websocket.socket.search
      },
      set: function () {
      }
    },
    style () {
      return {
        'height': this.size + 'px'
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

  methods: {
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
