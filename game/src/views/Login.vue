<template>
    <b-container>
      <b-row>
        <b-col cols="12"><h1>登入</h1></b-col>
      </b-row>
      <b-row v-if="$store.state.auth_token === ''">
        <b-col cols="12">
            <b-form class="login_form" @submit.prevent="login">
              <b-form-group description="請輸入剛剛發下的帳號">
                <b-form-input type="text" v-model="loginForm.account" placeholder="帳號" required>
                </b-form-input>
              </b-form-group>
              <b-form-group description="請輸入剛剛發下的密碼">
                <b-form-input type="password" v-model="loginForm.password" placholder="密碼" required>
                </b-form-input>
              </b-form-group>
              <b-button type="submit" variant="dark">登入</b-button>
            </b-form>
        </b-col>
      </b-row>
      <p>登入狀態: {{ status }} </p>
      <!-- if the authentication token is evaluated true, show logout button-->
      <b-button v-if="$store.state.auth_token" @click="logout">登出</b-button>
    </b-container>
</template>

<script>

export default {
  name: 'Login',
  data() {
    return {
      loginForm: {
        account: '',
        password: ''
      }
    }
  },
  methods: {
    login(evt){
      console.log("[*] User login form: ", this.loginForm)
      let account = this.loginForm.account
      let password = this.loginForm.password
      // perform backend authentication
      this.$store.dispatch("login", {
        account: account,
        password: password
      })
      // Clear loginForm
      this.loginForm.account = ''
      this.loginForm.password = ''
    },
    logout(evt){
      console.log("[*] User trys to log out...") 
      this.$store.dispatch("logout")
    }
  },
  computed: {
    status(){
      return this.$store.state.auth_status
    }
  }
}


</script>
<style scoped>
.login_form {
    width: 50%;
    margin: auto; /* center*/
}
</style>
