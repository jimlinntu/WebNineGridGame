<template>
  <b-row>
    <b-col cols="12"><h3>已經經過:</h3></b-col>
    <b-col cols="12"><b-icon-clock></b-icon-clock> {{ minutes }} 分 {{ seconds }} 秒</b-col>
  </b-row>
</template>

<script>
import { BIconClock} from 'bootstrap-vue'

export default {
  name: "TimeStatus",
  props: ["elapsedseconds"],
  components: {
    BIconClock
  },
  data(){
    return{
    }
  },
  computed: {
    minutes(){
      return Math.floor(this.elapsedseconds / 60)
    },
    seconds(){
      return this.elapsedseconds % 60
    }
  }
}
</script>
