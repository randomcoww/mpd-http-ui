const defaults = {
  socket: {
    isConnected: false,
    reconnectError: false,
    status: {},
    playlist: [],
    search: [],
    currentSong: {},
    elapsed: null,
    duration: null,
    databaseUpdateIndex: 0
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

    playlistadd (state, message) {
      let startPos = message.value[0]
      let addLength = message.value[1]
      for (var i = 0; i < addLength; i++) {
        state.socket.playlist.splice(startPos + i, 1, {})
      }
    },

    playlistchange (state, message) {
      let startPos = message.value[0]
      let deleteLength = message.value[1]
      let trimStartPos = message.value[2]
      let trimLength = message.value[3]
      for (var i = 0; i < deleteLength; i++) {
        state.socket.playlist.splice(startPos + i, 1, {})
      }
      if (trimLength > 0) {
        state.socket.playlist.splice(trimStartPos, trimLength)
      }
    },

    status (state, message) {
      state.socket.status = message.value
    },

    currentsong (state, message) {
      state.socket.currentSong = message.value
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

    playlistquery (state, message) {
      // state.socket.playlist = message.value
      message.value.map(v => {
        state.socket.playlist.splice(v.Pos, 1, v)
      })
    },

    updatedb (state, message) {
      state.socket.databaseUpdateIndex += 1
    }
  }
}

export default websocket
