const defaults = {
  socket: {
    isConnected: false,
    reconnectError: false,
    status: {},
    playlist: [],
    search: [],
    currentsong: {},
    elapsed: null,
    duration: null
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
      message.value.map(v => {
        state.socket.playlist.splice(v.Pos, 1, v)
      })
    },
    status (state, message) {
      console.info(message.value)
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
      state.socket.search = message.value
    },
    elapsed (state, message) {
      state.socket.elapsed = message.value
    }
  }
}

export default websocket
