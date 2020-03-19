<template>
    <b-container class="fluid">
      <b-row>
        <b-col cols="12"><h1>問題<b-icon-question></b-icon-question></h1></b-col>
      </b-row>
      <b-row class="text-center">
        <b-col cols="12" v-if="$store.state.question.description"><h3> {{ $store.state.question.description }} </h3></b-col>
        <b-col cols="12" v-if="$store.state.question.image"><img :src="'data:image/png;base64, ' + $store.state.question.image"/></b-col>
      </b-row>
      <hr>
      <b-row v-if="$store.state.question.description">
        <b-col cols="12">
          <b-form-input ref="input" v-model="answer.text" placeholder="請輸入答案"></b-form-input>
        </b-col>
      </b-row>
      <hr>
      <b-row v-if="$store.state.question.description">
        <b-col cols="12"><b-alert show variant="danger">照片 只能上傳 png 或是 jpg (jpeg) 格式唷!(如果底下有正確顯示就可以了)</b-alert></b-col>
        <b-col cols="12">
          <b-form-file v-model="answer.file" :state="Boolean(answer.file)" placeholder="請上傳照片"></b-form-file>
        </b-col>
        <b-col cols="12">
          <img :src="answer.base64_str"/>
        </b-col>
      </b-row>
      <hr>
      <b-row class="text-center" v-if="$store.state.question.description">
        <b-col><b-button @click.prevent="submitAnswer">提交答案</b-button></b-col>
      </b-row>
      <hr>
      <b-row class="text-center">
        <b-col cols="12"><h4>之前已提交的答案為:</h4></b-col>
        <b-col cols="12"class="previous_answer">{{ getCurrentAnswer.answertext }} </b-col>
        <b-col cols="12"><img :src="getCurrentAnswer.answerbase64str"/></b-col>
      </b-row>
    </b-container>
</template>

<script>

function getBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = error => reject(error);
  });
}

export default {
  name: "Question",
  data(){
    return {
      answer: {
        text: "",
        file: null,
        base64_str: null
      },
    }
  },
  watch: {
    'answer.file' (newFile){
      if(newFile !== null){
        getBase64(newFile).then((data) =>{
          this.answer.base64_str = data
        })
      }else{
        this.answer.base64_str = null
      }
    }
  },
  computed: {
    getCurrentAnswer(){
      // get previous answer from backend
      return this.$store.state.answer
    }
  },
  methods:{
    async submitAnswer(evt){
      // questionIndex must not be -1
      console.log("[*] Submit questionIndex == : ", this.$store.state.questionIndex)
      if(this.$store.state.questionIndex === -1){
        return
      }
      let base64_str = null
      if(this.answer.file !== null){
        base64_str = await getBase64(this.answer.file)
      }
      // Wait for this promise to be resolved
      this.$store.dispatch("submitAnswer", {
        text: this.answer.text,
        base64_str: base64_str
      })
      // Reset file to null
      this.answer.file = null
      // Reset text
      this.answer.text = ""
    }
  }
}
</script>

<style scoped>
.previous_answer {
  font-size: 20px;
}

img {
    width: 70%;
    height: auto;
}
</style>
