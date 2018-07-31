<template lang="pug">
  v-app.my-default-layout(dark)
    app-sidebar
    app-bar
    v-content(style="padding-top: 48px;")
      transition(name="slide" mode="out-in")
        router-view
    app-dialog
    app-snackbar
</template>

<script>
import AppBar from '@/components/app-bar'
import AppSidebar from '@/components/app-sidebar'
import AppDialog from '@/components/app-dialog'
import AppSnackbar from '@/components/app-snackbar'

export default {
  name: 'DefaultLayout',

  components: { AppBar, AppSidebar, AppDialog, AppSnackbar },

  methods: {

  }
}
</script>

<style lang="stylus">
</style>
