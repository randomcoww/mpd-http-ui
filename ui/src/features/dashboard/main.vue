<template lang="pug">
v-container.my-dashboard(
  fluid
  text-xs-center
  v-bind:grid-list-sm="$vuetify.breakpoint.smAndDown"
  v-bind:grid-list-lg="$vuetify.breakpoint.mdAndUp"
)
  v-layout(row wrap align-content-center)
    v-flex(d-flex xs12 sm12 md6 height="100%")
      v-layout(row wrap)
        v-flex(d-flex xs12 sm12 md12)
          player-status()
        v-flex(d-flex xs12 sm12 md12)
          playlist()
    v-flex(d-flex xs12 sm12 md6 max-height="100%")
      v-layout(row wrap)
        v-flex(d-flex xs12 sm12 md12)
          search-results()
</template>

<script>
import store from './store' // eslint-disable-line no-unused-vars
import Playlist from './components/playlist'
import SearchResults from './components/search-results'
import PlayerStatus from './components/player-status'

export default {
  name: 'Dashboard',

  components: {
    Playlist,
    SearchResults,
    PlayerStatus
  },

  data () {
    return {
      test: this.$store.state.dashboard.test,
      date: '2018-05-21'
    }
  },

  mounted () {
  },

  methods: {
    updateTest () {
      this.test++
      this.$store.dispatch('dashboard/updateTest', this.test)
    },
    functionEvents (date) {
      const [,, day] = date.split('-')
      return parseInt(day, 10) % 3 === 0
    }
  }
}
</script>

<style lang="stylus">
.my-dashboard

  &__media
    height: 100%
    margin: 0

  .picker__title
    display: none !important
</style>
