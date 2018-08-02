<template lang="pug">
v-card(dense height="100%")
  sidebar-library()
  sidebar-playlist()
  player-status()
</template>

<script>
// import store from './store' // eslint-disable-line no-unused-vars
import PlayerStatus from './components/player-status'
import SidebarLibrary from './components/sidebar-library'
import SidebarPlaylist from './components/sidebar-playlist'

export default {
  name: 'Dashboard',

  components: {
    PlayerStatus,
    SidebarLibrary,
    SidebarPlaylist
  },

  mounted () {
  }
}
</script>

<style lang="stylus">

</style>
