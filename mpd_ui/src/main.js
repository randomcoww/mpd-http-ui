import Vue from 'vue'
import store from './store'
import { sync } from 'vuex-router-sync'
import { router } from './router'
import Vuetify from 'vuetify'
// import URLSearchParams from 'url-search-params'
import App from './app'
import Appbar from './components/app-bar'
import VueNativeSock from 'vue-native-websocket'

Vue.config.productionTip = false

// Sync router to store, as `store.state.route`.
sync(store, router)

// Vuetify
Vue.use(Vuetify, {
  theme: {
    primary: '#21CE99',
    secondary: '#D81B60',
    accent: '#805441'
  }
})

// Websocket
Vue.use(VueNativeSock, process.env.NODE_ENV === 'produciton' ? 'ws://' + location.host + '/ws' : 'ws://localhost:3000/ws', {
  reconnection: true,
  reconnectionDelay: 3000,
  store: store,
  format: 'json'
})

// Styles
require('./styles/scss/main.scss')
require('./styles/stylus/main.styl')

// Global Components
Vue.component('Appbar', Appbar)

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App)
})
