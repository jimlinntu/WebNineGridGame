import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import GridGame from '../views/GridGame.vue'
import Login from '../views/Login.vue'
import Admin from '../views/Admin.vue'
import store from '../store'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'About',
    component: () => import('../views/About.vue')
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  },
  {
    path: '/gridgame',
    name: 'GridGame',
    component: GridGame,
    beforeEnter: (to, from, next) =>{
        if(store.state.auth_token){
            // TODO: Fetch grid numbers, questions and index from the backend server
            store.dispatch("getGridNumbers")
            next()
            return
        }
        alert("請先登入!")
        // redirect it to login page
        next("/login")
    }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: Admin,
    beforeEnter: (to, from, next) =>{
        // TODO: How to verify this user is an admin
        if(store.state.auth_token){
            next()
            return
        }
        alert("請先登入!")
        // redirect it to login page
        next("/login")
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  }

]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
