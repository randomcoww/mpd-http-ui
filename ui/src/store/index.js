import Vue from 'vue'
import Vuex from 'vuex'
import websocket from '@/websocket/store'
import common from './common'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: { common, websocket }
})
