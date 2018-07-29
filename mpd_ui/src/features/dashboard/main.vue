<template lang="pug">
v-card(dense)
  v-container.my-dashboard(
    fluid
    text-xs-center
  )
    v-layout(row wrap align-content-center)
      v-flex(d-flex xs12 sm12 md12 height="100%")
        v-layout(row wrap)
          v-flex(d-flex xs12 sm12 md12)
            player-status()
          v-flex(d-flex xs12 sm12 md12)
            playlist()

</template>

<script>
// import store from './store' // eslint-disable-line no-unused-vars
import Playlist from './components/playlist'
import PlayerStatus from './components/player-status'

export default {
  name: 'Dashboard',

  components: {
    Playlist,
    PlayerStatus
  },

  mounted () {
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
