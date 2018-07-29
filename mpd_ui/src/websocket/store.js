const defaults = {
  socket: {
    isConnected: false,
    reconnectError: false,
    status: {},
    playlist: [],
    search: [],
    currentsong: {},
    elapsed: null,
    duration: null,
    version: null
  }
}

const websocket = {
  // namespaced: true,
  state: Object.assign({}, defaults),

  mutations: {
    SOCKET_ONOPEN (state, event) {
      state.socket.isConnected = true
    },

    SOCKET_ONCLOSE (state, event) {
      state.socket.isConnected = false
    },

    SOCKET_ONERROR (state, event) {
      console.error(state, event)
    },

    // default handler called for all methods
    SOCKET_ONMESSAGE (state, message) {
      console.info(state, message)
    },

    // mutations for reconnect methods
    SOCKET_RECONNECT (state, count) {
      console.info(state, count)
    },

    SOCKET_RECONNECT_ERROR (state) {
      state.socket.reconnectError = true
    },

    playlist (state, message) {
      state.socket.version = message.value[0]
      state.socket.playlist.splice(message.value[1])
    },

    status (state, message) {
      state.socket.status = message.value
    },

    currentsong (state, message) {
      state.socket.currentsong = message.value
    },

    seek (state, message) {
      state.socket.elapsed = message.value[0]
      state.socket.duration = message.value[1]
    },

    search (state, message) {
      // console.info(message.value)
      let results = message.value[0]
      let start = parseInt(message.value[1])
      // console.info('searchstart', start)
      if (results != null) {
        results.map(v => {
          state.socket.search.splice(start, 1, v)
          start++
        })
      }
      state.socket.search.splice(start)
      // console.info('searchend', start)
      // state.socket.search = message.value
    },

    elapsed (state, message) {
      state.socket.elapsed = message.value
    },

    playlistupdate (state, message) {
      // state.socket.playlist = message.value
      message.value.map(v => {
        state.socket.playlist.splice(v.Pos, 1, v)
      })
    }
  }
}

export default websocket
