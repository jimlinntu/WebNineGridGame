import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import { BootstrapVue, BootstrapVueIcons } from 'bootstrap-vue'
import Grid from 'vue-js-grid'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
// SocketIO
import VueSocketIO from 'vue-socket.io'
import SocketIO from 'socket.io-client'

import App from './App.vue'
import router from './router'
import store from './store'

Vue.use(BootstrapVue)
Vue.use(BootstrapVueIcons)
Vue.use(Grid)
Vue.use(VueAxios, axios)
// SocketIO
// TODO: bug...
Vue.use(new VueSocketIO({
    debug: true,
    connection: SocketIO('http://125.227.38.80:17989', {transports: ['polling'], path: "/socket"}),
    vuex:{
        store,
        actionPrefix: 'SOCKET_',
        mutationPrefix: 'SOCKET_'
    },
}))
Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
