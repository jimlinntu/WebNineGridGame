<template>
  <div id="app">
    <div id="nav">
      <router-link to="/about">關於</router-link> |
      <router-link v-if="!isAdmin" to="/gridgame">九宮格遊戲</router-link> |
      <router-link to="/login">登入/登出</router-link> |
      <!-- (Soft block) Only allow admin to access this link-->
      <router-link v-if="isAdmin" to="/admin">管理者介面</router-link>
    </div>
    <router-view/>
  </div>
</template>

<script>
export default {
    name: "App",
    computed:{
      isAdmin(){
        return this.$store.state.auth_token.indexOf('admin') !== -1
      }
    }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}
</style>
