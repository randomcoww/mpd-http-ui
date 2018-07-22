<template lang="pug">
v-card
  v-card-title
    .title Library
  v-card-text
    v-text-field(label="Search..." v-model.lazy="databasequery")
    v-layout(row wrap style="align-items: center;")
      template(v-for="searchresult in searchresults")
        v-flex(d-flex xs12 sm12 md4)
          | {{searchresult.artist}}
        v-flex(d-flex xs12 sm12 md8)
          v-layout(style="align-items: center;")
            v-flex.text-xs-left(md10) {{searchresult.title}}
            v-flex.text-xs-right(md2)
              v-btn(flat icon color="primary")
                v-icon delete
</template>

<script>
export default {
  data () {
    return {
      databasequery: null,
      errored: false
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
  }
}
</script>
