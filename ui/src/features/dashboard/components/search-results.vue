<template lang="pug">
v-card.searchresults
  v-card-title
    .title Library
  v-card-text
    v-text-field(label="Search..." v-model.lazy="databasequery")
    virtual-list(:size="this.size" :remain="this.buffer" :onscroll="onscroll" :tobottom="tobottom")
      ul(v-for="searchresult in searchresults")
        v-flex(d-flex :style="style")
          v-layout(row wrap style="align-items: center;")

            v-flex(d-flex xs10 sm10 md10)
              v-flex(d-flex xs12 sm12 md4)
                | {{searchresult.artist}}
              v-flex.text-xs-left(md8)
                | {{searchresult.title}}

            v-flex(d-flex xs2 sm2 md2)
              v-layout(style="align-items: center;")
                v-btn(flat icon color="primary")
                  v-icon add
</template>

<script>
import VirtualList from 'vue-virtual-scroll-list'

export default {
  components: {
    VirtualList
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
    searchresults () {
      return this.$store.state.websocket.socket.search
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
    tobottom () {
    },

    onscroll (event, data) {
    }
  }
}
</script>
