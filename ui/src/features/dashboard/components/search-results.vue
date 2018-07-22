<template lang="pug">
v-card
  v-card-title
    .title Library
  v-card-text
    v-text-field(label="Search..." v-model.lazy="databasequery")
    virtual-list(:size="64" :remain="this.buffer" :onscroll="onscroll" :tobottom="tobottom" :totop="totop")
      ul(v-for="searchresult in searchresults")
        v-flex(d-flex)
          v-layout(row wrap style="align-items: center;")
            v-flex(d-flex xs12 sm12 md4)
              | {{searchresult.artist}}
            v-flex(d-flex xs12 sm12 md8)
              v-layout(style="align-items: center;")
                v-flex.text-xs-left(md10) {{searchresult.title}}
                v-flex.text-xs-left(md2)
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
      databasequery: null,
      errored: false,
      buffer: 15
    }
  },

  computed: {
    searchresults () {
      return this.$store.state.websocket.socket.search
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

    totop () {
    },

    onscroll (event, data) {
    }
  }
}
</script>
