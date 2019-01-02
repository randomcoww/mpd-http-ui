// Common State.
const defaults = {
  playlist: {
    visible: true
  },
  library: {
    visible: false
  },
  title: '',
  layout: 'DefaultLayout',
  dialog: {
    visible: false,
    text: ''
  },
  snackbar: {
    show: false,
    visible: true,
    text: '',
    timeout: 6000,
    color: ''
  },
  error: {
    code: null,
    level: null,
    message: ''
  }
}

// Global module loaded on first app load.
export default {
  namespaced: true,

  state: Object.assign({}, defaults),

  mutations: {
    togglePlaylist (state, options) {
      state.playlist = Object.assign({}, defaults.playlist, options)
    },

    toggleLibrary (state, options) {
      state.library = Object.assign({}, defaults.library, options)
    },

    updateTitle (state, title) {
      state.title = title
    },

    updateLayout (state, layout) {
      state.layout = layout
    },

    updateSnackbar (state, options) {
      state.snackbar = Object.assign({}, defaults.snackbar, options)
    },

    error (state, options) {
      state.error = Object.assign({}, defaults.error, options)
    },

    clear (state) {
      state = Object.assign({}, defaults)
    }
  },

  actions: {
    toggleLibrary ({ commit }, options) {
      commit('toggleLibrary', options)
    },

    togglePlaylist ({ commit }, options) {
      commit('togglePlaylist', options)
    },

    updateTitle ({ commit }, title) {
      commit('updateTitle', title)
    },

    updateLayout ({ commit }, layout) {
      commit('updateLayout', layout)
    },

    updateSnackbar ({ commit }, options) {
      commit('updateSnackbar', options)
    }
  }
}
