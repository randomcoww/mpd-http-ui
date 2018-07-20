const defaults = {
  socket: {
    isConnected: false,
    reconnectError: false,
    status: {},
    playlist: null,
    currentsong: {}
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
      // console.info(message.value)
      state.socket.playlist = message.value
    },
    status (state, message) {
      // console.info(message.value)
      state.socket.status = message.value
    },
    currentsong (state, message) {
      // console.info(message.value)
      state.socket.currentsong = message.value
    }
  }
}

export default websocket
