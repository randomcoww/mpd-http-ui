<template lang="pug">
v-card(dense height="100%")
  player-status()
  v-divider(inset)
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

</style>
