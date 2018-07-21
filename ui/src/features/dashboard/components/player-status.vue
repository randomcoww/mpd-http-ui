<template lang="pug">
v-card-text
  v-slider(:value="seek")
</template>

<script>

export default {
  computed: {
    seek () {
      return this.$store.state.websocket.socket.seek
    }
  }
}
</script>
