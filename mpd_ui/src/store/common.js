// Common State.
const defaults = {
  sidebar: {
    visible: true
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
    updateSidebar (state, options) {
      state.sidebar = Object.assign({}, defaults.sidebar, options)
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
    updateSidebar ({ commit }, options) {
      commit('updateSidebar', options)
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
