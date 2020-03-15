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
      question: {}, // question = {description: ..., image: ...}
      questionIndex: -1,
      question_finished_mask: [],
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
      set_grid_numbers_and_questions(state, {gridNumbers, questionIndex, question_finished_mask}){
          state.gridNumbers = gridNumbers
          state.questionIndex = questionIndex
          state.question_finished_mask = question_finished_mask
          // TODO: question
          
      },
      set_questionIdx_and_question(state, {questionIndex, question}){
          state.questionIndex = questionIndex
          // TODO: question
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
              context.commit("set_login_success", {
                  token: response.data.token,
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
                    questionIndex: response.data.questionIndex,
                    question_finished_mask: response.data.question_finished_mask,
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
                questionIndex: response.data.questionIndex,
                question_finished_mask: response.data.question_finished_mask,
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
                question: {} // TODO
            })
        }).catch((error)=>{
            console.log(error)
        })
      }

  },
  modules: {
  }
})
