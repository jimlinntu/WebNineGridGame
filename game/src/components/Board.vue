<template>
  <b-row>
    <b-col cols="12">備註: 1. <b class="gold_color">金黃色</b> 為組別正在解的題目, 2.<b class="green_color">綠色</b> 為組別已經完成的題目</b-col>
    <hr>
    <b-col cols="12"><h1> 組別: {{ team }} </h1></b-col>
    <b-col cols="12"><h2> 目前連線數為: {{ num_lines }} </h2></b-col>
    <b-col class="slot" :class="{'question_finished': question_finished_mask[index],
           'question_index': index === questionIndex
          }"
          cols="3" v-for="(number, index) in gridNumbers" :key="index"
          @click="deleteFinished($event, index)">
      {{ number }}
    </b-col>
    <b-col cols="12" v-if="gridNumbers.length === 0">
      <b-alert show>該組別尚未上傳九宮格</b-alert>
    </b-col>
    <b-col v-if="question !== null" cols="12"><h3>問題: {{ question.description }} </h3></b-col>
    <b-col v-if="question !== null && question.base64image" cols="12"><img class="smaller_img" :src="'data:image/png;base64, '+ question.base64image"/></b-col>
    <b-col cols="12"><h3>當前該組回答: </h3></b-col>
    <b-col cols="12">答案: {{ answertext }}</b-col>
    <b-col cols="12"><img :src="answerbase64str"/></b-col>
  </b-row>
</template>

<script>

export default {
  name: "Board",
  props: ['team', 'gridNumbers', 'questionIndex', 'question_finished_mask', 
          'answertext', 'answerbase64str', 'question'],
  data(){
    return {
    }
  },
  methods: {
    deleteFinished(evt, index){
      if(this.question_finished_mask === null) return;
      let len = this.question_finished_mask.length;
      if(this.question_finished_mask[index] === false) return;
      if(!confirm("你確定要刪除這組已經完成的這一題嗎?")) return;

      this.$store.dispatch("deleteFinished", {
        target_account: this.team,
        questionIndex: index,
      })
      return;
    }
  },
  computed: {
    num_lines(){
      if(this.question_finished_mask === null) return 0;
      let len = this.question_finished_mask.length;
      if(len !== 16) return -1;
      let side_len = 4;
      let counter = 0;
      for(let i = 0; i < side_len; i++){
        // Compute horizontal
        let has_line = true;
        for(let j = 0; j < side_len && has_line; j++){
          let idx = side_len * i + j;
          if(this.question_finished_mask[idx] === false){
            has_line = false;
          }
        }
        if(has_line) counter += 1;
        // Compute vertical
        has_line = true;
        for(let j = 0; j < side_len && has_line; j++){
          let idx = side_len * j + i;
          if(this.question_finished_mask[idx] === false){
            has_line = false;
          }
        }
        if(has_line) counter += 1;
      }
      // Compute the main diagonal
      let has_line = true;
      for(let i = 0; i < side_len && has_line; i++){
        let idx = side_len * i + i;
        if(this.question_finished_mask[idx] === false){
          has_line = false;
        }
      }
      if(has_line) counter += 1;
      // Compute the antidiagonal
      has_line = true;
      for(let i = 0; i < side_len && has_line; i++){
        let idx = side_len * i + (side_len - i - 1);
        if(this.question_finished_mask[idx] === false){
          has_line = false;
        }
      }
      if(has_line) counter += 1;
      return counter;
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
    background-color: rgb(255, 165, 0);
}

.green_color {
    color: rgb(0, 128, 0);
    background-color: rgb(0, 128, 0);
}

.question_index{
    color: rgb(255,165,0);
    background-color: rgb(255, 165, 0, 0.3);
}

.question_finished{
    color: rgb(0, 128, 0);
    background-color: rgb(0, 128, 0, 0.3);
}

img {
  width: 70%;
  height: auto;
}

.smaller_img {
  width: 40%;
}

</style>
