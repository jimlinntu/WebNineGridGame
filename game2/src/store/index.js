import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
      backend_url: "http://125.227.38.80:17989/",
      auth_status: '未登入',
      auth_token: ''
  },
  mutations: {
      set_login_loading(state) {
          state.auth_status = "登入中"
      },
      set_login_success(state, token) {
          state.auth_status = "登入成功"
          state.auth_token = token
      },
      set_login_fail(state) {
          state.auth_status = "登入失敗"
      }
  },
  actions: {
      login(context, { account, password }){
          context.commit("set_login_loading")

          // perform asynchronous backend authentication
          Vue.axios.post(this.state.backend_url + "auth",{
                account: account,
                password: password
              }
            ).then((response)=>{
              console.log(response)
              context.commit("set_login_success", "TODOTOKEN")
          }).catch((error)=>{
              console.log(error)
              context.commit("set_login_fail")
          })
      }
  },
  modules: {
  }
})
