<template>
  <div id="app">
    <Header/>
    <b-container class="bv-example-row">
        <b-row>
            <b-col><GridGame/></b-col>
        </b-row>
    </b-container>
  </div>
</template>

<script>
import Header from './components/Header.vue'
import GridGame from './components/GridGame.vue'

export default {
  name: 'App',
  components: {
    Header,
    GridGame
  },
  data(){
    return {
        questions: []
    }
  },
  mounted: function(){
    fetch('https://opentdb.com/api.php?amount=10', {
        method: 'get' 
    }).then((response) => {
        return response.json()
    }).then((jsonData) => {
        this.questions = jsonData.results
    })
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
  margin-top: 20px;
}
</style>
