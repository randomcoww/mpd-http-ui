<template lang="pug">
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
import axios from 'axios'

export default {
  data () {
    return {
      databasequery: null,
      searchresults: [],
      errored: false
    }
  },

  watch: {
    databasequery (after, before) {
      this.fetch()
    }
  },

  methods: {
    fetch () {
      axios.get('http://localhost:3000/database/search', {
        params: {
          q: this.databasequery,
          size: 100
        }
      }).then(response => {
        this.searchresults = response.data
      }).catch(error => {
        console.log(error)
        this.errored = true
      })
    }
  }
}
</script>
