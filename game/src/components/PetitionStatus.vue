<template>
    <b-row class="text-center">
      <b-col cols="12"><h3>卡關請願跳題狀態</h3></b-col>
      <b-col cols="12" v-if="haspetition === false">
          <b-icon-chat></b-icon-chat>尚未請願跳題!加油
      </b-col>
      <b-col cols="12" v-else>
          <b-icon-check></b-icon-check>已請願跳題!等待主辦人批准跳題
      </b-col>
    </b-row>
</template>
<script>
export default {
  name: "PetitionStatus",
  props: {
    haspetition: Boolean,
  },
  data(){
    return {
    }
  }
}
</script>
