<template lang="pug">
v-container.my-dashboard(
  fluid
  fill-height
  text-xs-center
  v-bind:grid-list-sm="$vuetify.breakpoint.smAndDown"
  v-bind:grid-list-lg="$vuetify.breakpoint.mdAndUp"
)
  v-layout(row wrap)

    // Top Row
    v-flex(d-flex xs12 sm12 md6)
      v-card
        v-card-title
          .title Currently Selected
        v-card-text
          current-song()

    v-flex(d-flex xs12 sm12 md6)
      v-card.my-dashboard__pizza-status
        v-card-title
          .title Status
        v-card-title
          player-status()

    // Bottom Row
    v-flex(d-flex xs12 sm12 md6)
      v-card.my-dashboard__pizza-status
        v-card-title
          .title Playlist
        v-card-text
          playlist()

    v-flex(d-flex xs12 sm12 md6)
      v-card.my-dashboard__pizza-status
        v-card-title
          .title Library
        v-card-title
          search-results()

</template>

<script>
import store from './store' // eslint-disable-line no-unused-vars
import Chart from './components/chart'
import LineChart from './components/line-chart'
import Playlist from './components/playlist'
import CurrentSong from './components/current-song'
import SearchResults from './components/search-results'
import PlayerStatus from './components/player-status'

export default {
  name: 'Dashboard',

  components: {
    Chart,
    LineChart,
    Playlist,
    CurrentSong,
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
