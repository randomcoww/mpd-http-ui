<template lang="pug">
v-snackbar(
  :bottom="true"
  :right="true"
  :timeout='$store.state.common.snackbar.timeout'
  v-model='snackbarActive')
  | {{ $store.state.common.snackbar.text | truncate(100) }}
  v-btn(dark='' flat='' @click.native='snackbarActive = false')
    v-icon close
</template>

<script>
export default {
  name: 'DefaultSnackbar',

  filters: {
    truncate: function (string, value) {
      return string.substring(0, value) + '...'
    }
  },

  data () {
    return {
    }
  },

  computed: {
    snackbarActive: {
      get () {
        return this.$store.state.common.snackbar.show
      },
      set (val) {
        this.$store.dispatch('common/updateSnackbar', { show: val })
      }
    }
  }
}
</script>

<style lang="stylus">

</style>
