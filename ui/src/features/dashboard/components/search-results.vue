<template lang="pug">
v-card.searchresults
  v-card-title
    .title Library
  v-card-text
    v-text-field(label="Search..." v-model.lazy="databasequery")
    virtual-list(:size="this.size" :remain="this.buffer" :onscroll="onscroll" :tobottom="tobottom")
      div(v-for="(searchresult, index) in searchresults" :index="index" :key="searchresult.file")
        draggable(v-model="searchresults" @end="onmoved" :options="{group: {name: 'playlistitems'}}" :id="searchresult.file")
          v-flex(d-flex :style="style")
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
      databasequery: null,
      errored: false,
      buffer: 25
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
    databasequery (after, before) {
      this.$socket.sendObj({ mutation: 'search', value: [this.databasequery, 100] })
    }
  },

  methods: {
    addpath (path, position) {
      console.info('addpath', path, position)
      this.$socket.sendObj({ mutation: 'addpath', value: [path, parseInt(position)] })
    },

    onmoved (event, data) {
      this.addpath(event.from.id, event.to.id)
    },

    tobottom () {
    },

    onscroll (event, data) {
    }
  }
}
</script>
