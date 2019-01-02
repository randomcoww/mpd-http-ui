import Vue from 'vue'
import Router from 'vue-router'
import routes from './routes'
import store from '@/store'

Vue.use(Router)

/**
 * The Router instance containing all the routes for the application.
 */
const router = new Router({
  base: '/app',
  // mode: 'history',  // <-- uncomment to turn on history mode (preferred)
  routes: routes.map(route => ({
    name: route.name,
    path: route.path,
    component: route.component,
    beforeEnter: (to, from, next) => {
      // Setup some per-page stuff.
      document.title = route.title
      store.dispatch('common/updateTitle', route.title)
      store.dispatch('common/updateLayout', route.layout)

      next()
    }
  }))
})

export default router
