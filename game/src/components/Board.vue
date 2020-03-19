<template>
  <b-row>
    <b-col cols="12">備註: 1. <b class="gold_color">金黃色</b> 為組別正在解的題目, 2.<b class="green_color">綠色</b> 為組別已經完成的題目</b-col>
    <hr>
    <b-col cols="12"><h1> 組別: {{ team }} </h1></b-col>
    <b-col class="slot" :class="{'question_finished': question_finished_mask[index],
           'question_index': index === questionIndex
          }"
          cols="4" v-for="(number, index) in gridNumbers" :key="index">
      {{ number }}
    </b-col>
    <b-col cols="12" v-if="gridNumbers.length === 0">
      <b-alert show>該組別尚未上傳九宮格</b-alert>
    </b-col>
    <b-col cols="12"><h3>當前該組回答: </h3></b-col>
    <b-col cols="12">答案: {{ answertext }}</b-col>
    <b-col cols="12"><img :src="answerbase64str"/></b-col>
  </b-row>
</template>

<script>

export default {
  name: "Board",
  props: ['team', 'gridNumbers', 'questionIndex', 'question_finished_mask', 
          'answertext', 'answerbase64str'],
  data(){
    return {
    }
  }
}
</script>
<style scoped>

.slot {
    border-style: solid;
    height: 50px;
}

.gold_color {
    color: rgb(255,165,0);
}

.green_color {
    color: rgb(0, 128, 0);
}

.question_index{
    color: rgb(255,165,0);
}

.question_finished{
    color: rgb(0, 128, 0);
}

img {
  width: 70%;
  height: auto;
}

</style>
