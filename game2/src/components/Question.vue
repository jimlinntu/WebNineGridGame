<template>
    <b-container class="fluid">
      <b-row class="justify-content-md-center">
        <h1>問題<b-icon-question></b-icon-question></h1>
      </b-row>
      <b-row class="text-center">
        <b-col> </b-col>
      </b-row>
      <b-row>
        <b-col cols="12">
          <b-form-input v-model="answer.text" placeholder="請輸入答案"></b-form-input>
        </b-col>
      </b-row>
      <hr>
      <b-row>
        <b-col cols="12">
          <b-form-file v-model="answer.file" :state="Boolean(answer.file)" placeholder="請上傳照片"></b-form-file>
        </b-col>
      </b-row>
      <hr>
      <b-row class="text-center">
        <b-col><b-button @click.prevent="submitAnswer">提交答案</b-button></b-col>
      </b-row>
      <hr>
      <b-row class="text-center">
        <b-col class="previous_answer">之前已提交的回答為: {{ getPreviousAnswer }}</b-col>
      </b-row>
      <img :src="answer.base64_str"/>
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
      }
    }
  },
  computed: {
    getPreviousAnswer(){
      // get previous answer from backend
      return "previous answers"
    }
  },
  methods:{
    async submitAnswer(){
      let base64_str = null
      if(this.answer.file !== null){
        base64_str = await getBase64(this.answer.file)
      }
      this.$store.dispatch("submitAnswer", {
        text: this.answer.text,
        base64_str: base64_str
      })
    }
  }
}
</script>

<style scoped>
.previous_answer {
  font-size: 20px;
}
</style>
