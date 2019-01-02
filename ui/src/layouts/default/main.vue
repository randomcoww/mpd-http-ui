<template lang="pug">
  v-app.my-default-layout(dark)
    app-bar()
    v-content(style="padding-top: 48px;")
      transition(name="slide" mode="out-in")
        router-view
    app-dialog()
    app-snackbar()
</template>

<script>
import AppBar from '@/components/app-bar'
import AppDialog from '@/components/app-dialog'
import AppSnackbar from '@/components/app-snackbar'

export default {
  name: 'DefaultLayout',

  components: {
    AppBar,
    AppDialog,
    AppSnackbar
  },

  methods: {

  }
}
</script>

<style lang="stylus">
</style>
