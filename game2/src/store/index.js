import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

function getBase64(file){
  let reader = new FileReader();
  reader.readAsDataURL(file);
  reader.onload = function (){
    console.log(reader.result)
  };
  reader.onerror = function (error){
    console.log("Error: ", error)
  };
}

export default new Vuex.Store({
  state: {
      backend_url: "http://125.227.38.80:17989/",
      auth_status: '未登入',
      auth_token: '',
      gridNumbers: [], // grid numbers retrieve from backend
      questionOrder: [],
      questionDescriptions: [],
      questionIndex: -1,
  },
  mutations: {
      set_login_loading(state) {
          state.auth_status = "登入中"
      },
      set_login_success(state, {token , gridNumbers}) {
          state.auth_status = "登入成功"
          state.auth_token = token
          state.gridNumbers = gridNumbers
      },
      set_login_fail(state) {
          state.auth_status = "登入失敗"
      },
      set_grid_numbers_and_questions(state, {gridNumbers, 
            questionOrder, questionIndex}){
          state.gridNumbers = gridNumbers
          state.questionOrder = questionOrder
          state.questionIndex = questionIndex
      },
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
              context.commit("set_login_success", {
                  token: response.data.token,
                  gridNumbers: response.data.gridNumbers
              })
          }).catch((error)=>{
              console.log(error)
              context.commit("set_login_fail")
          })
      },
      submitGridNumbers(context, {gridNumbers}){
        // post gridNumbers to the backend server
        Vue.axios.post(this.state.backend_url + "user/push_gridnumbers", {
                gridnumbers: gridNumbers,
                token: this.state.auth_token
            }).then((response)=>{
                console.log(response)
                // change gridNumbers
                context.commit("set_grid_numbers_and_questions", {
                    gridNumbers: response.data.gridNumbers,
                    questionOrder: response.data.questionOrder,
                    questionIndex: response.data.questionIndex,
                })
                
            }).catch((error)=>{
                console.log(error)
            })
        
      },
      getGridNumbers(context){
        Vue.axios.post(this.state.backend_url + "user/get_gridnumbers", {
              token: this.state.auth_token
        }).then((response)=>{
            console.log(response)
            // assign gridNumbers
            context.commit("set_grid_numbers_and_questions", {
                gridNumbers: response.data.gridNumbers,
                questionOrder: response.data.questionOrder,
                questionIndex: response.data.questionIndex
            })

        }).catch((error)=>{
            console.log(error)
        })
      },
      submitAnswer(context, {text, base64_str}){
        Vue.axios.post(this.state.backend_url + "user/push_answer", {
            token: this.state.auth_token,
            answertext: text,
            answerbase64str: base64_str
        }).then((response)=>{
            console.log(response)
        }).catch((error)=>{
            console.log(error)
        })
      }
  },
  modules: {
  }
})
