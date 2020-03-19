import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
      backend_url: "http://125.227.38.80:17989/",
      auth_status: '未登入',
      auth_token: '',
      gridNumbers: [], // grid numbers retrieve from backend
      question: {}, // question = {description: ..., image: ...}
      questionIndex: -1,
      question_finished_mask: [],
      answer: {}, // answer corresponding to current question
      users: [],
  },
  mutations: {
      set_login_loading(state) {
          state.auth_status = "登入中"
      },
      set_login_success(state, {token}) {
          state.auth_status = "登入成功"
          state.auth_token = token
      },
      set_login_fail(state) {
          state.auth_status = "登入失敗"
      },
      set_logout(state){
          state.auth_status = "已登出"
          state.auth_token = "" // clear the authentication token
      },
      set_grid_numbers_and_questions_and_answer(state, {gridNumbers, questionIndex, 
                                             question_finished_mask, question,
                                             answertext, answerbase64str}){
          state.gridNumbers = gridNumbers
          state.questionIndex = questionIndex
          state.question_finished_mask = question_finished_mask
          state.question = question
          state.answer = {
              answertext: answertext,
              answerbase64str: answerbase64str
          }
      },
      set_questionIdx_and_question(state, {questionIndex, question}){
          state.questionIndex = questionIndex
          state.question = question
      },
      set_answer(state, {answertext, answerbase64str}){
          state.answer = {
            answertext: answertext,
            answerbase64str: answerbase64str,
          }
      },
      set_users_status(state, {users}){
          state.users = users 
      }
  },
  actions: {
      login(context, { account, password }){
          context.commit("set_login_loading")

          // perform asynchronous backend authentication
          Vue.axios.post(this.state.backend_url + "api/auth",{
                account: account,
                password: password
              }
            ).then((response)=>{
              console.log(response)
              context.commit("set_login_success", {
                  token: response.data.token,
              })
          }).catch((error)=>{
              console.log(error)
              context.commit("set_login_fail")
          })
      },
      logout(context){
        context.commit("set_logout")
      },
      submitGridNumbers(context, {gridNumbers}){
        // post gridNumbers to the backend server
        Vue.axios.post(this.state.backend_url + "user/push_gridnumbers", {
                gridnumbers: gridNumbers,
                token: this.state.auth_token
            }).then((response)=>{
                console.log(response)
                // change gridNumbers
                context.commit("set_grid_numbers_and_questions_and_answer", {
                    gridNumbers: response.data.gridNumbers,
                    questionIndex: response.data.questionIndex,
                    question_finished_mask: response.data.question_finished_mask,
                    question: {}, // there is no question
                    answertext: "", // there is no answer yet been submitted
                    answerbase64str: "" // there is no answer yet been submitted
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
            context.commit("set_grid_numbers_and_questions_and_answer", {
                gridNumbers: response.data.gridNumbers,
                questionIndex: response.data.questionIndex,
                question_finished_mask: response.data.question_finished_mask,
                question: response.data.question,
                answertext: response.data.answertext,
                answerbase64str: response.data.answerbase64str
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
            // Get current answer from the backend server
            context.dispatch("getAnswer")
        }).catch((error)=>{
            console.log(error)
        })
      },
      getAnswer(context){
        Vue.axios.post(this.state.backend_url + "user/get_answer", {
            token: this.state.auth_token,
        }).then((response)=>{
            console.log(response)
            // TODO: set current answer
            context.commit("set_answer", {
                answertext: response.data.answertext,
                answerbase64str: response.data.answerbase64str
            })
        }).catch((error)=>{
            console.log(error)
        })
      },
      selectQuestion(context, {questionIndex}){
        Vue.axios.post(this.state.backend_url + "user/select_question", {
          token: this.state.auth_token,
          questionIndex: questionIndex,
        }).then((response)=>{
            console.log(response)
            // update question finished mask
            context.commit("set_questionIdx_and_question", {
                questionIndex: response.data.questionIndex,
                question: response.data.question
            })
        }).catch((error)=>{
            console.log(error)
        })
      },
      // admininistrator can fetch all information
      // (gridnumbers, question mask, question index, answertext, answerbase64str)
      getAll(context){
        Vue.axios.post(this.state.backend_url + "user/get_all", {
            token: this.state.auth_token,
        }).then((response)=>{
            console.log(response)
            context.commit("set_users_status", {
                users: response.data.users
            })
            
        }).catch((error)=>{
            console.log(error)
        })
      },
      // reset all teams' gridnumbers
      resetAll(context){
        Vue.axios.post(this.state.backend_url + "user/reset_all", {
            token: this.state.auth_token,
        }).then((response)=>{
            console.log(response)
        }).catch((error)=>{
            console.log(error)
        })
      },
      approveAnswer(context, {target_account}){
        Vue.axios.post(this.state.backend_url + "user/approve_answer",{
            token: this.state.auth_token,
            account: target_account,
        }).then((response)=>{
          console.log(response)
        }).catch((error)=>{
            console.log(error)
        })
      },
      skipAnswer(context, {target_account}){
        Vue.axios.post(this.state.backend_url + "user/skip_answer", {
            token: this.state.auth_token,
            account: target_account,
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