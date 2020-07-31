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
          <img v-if="answer.base64_str !== null" :src="answer.base64_str"/>
        </b-col>
      </b-row>
      <hr>
      <b-row class="text-center" v-if="$store.state.question.description">
        <b-col cols="6"><b-button size="lg" @click.prevent="submitAnswer">提交答案</b-button></b-col>
        <b-col cols="6"><b-button variant="danger" size="lg" @click.prevent="petitionSkipQuestion">我卡關了,我想跳題!</b-button></b-col>
      </b-row>
      <hr>
        <b-row class="text-center" v-if="getUploadStatus !== ''">
          <b-col cols="12">{{ getUploadStatus }}</b-col>
        </b-row>
      <hr>
      <b-row class="text-center">
        <b-col cols="12"><h4>之前已提交的答案為:</h4></b-col>
        <b-col cols="12" class="previous_answer">{{ getCurrentAnswer.answertext }} </b-col>
        <b-col cols="12">
          <img v-if="getCurrentAnswer.answerbase64str" :src="getCurrentAnswer.answerbase64str"/>
        </b-col>
      </b-row>
      <hr>
      <!-- Only show the elapsed seconds when this user has chosen a question!-->
      <TimeStatus v-if="$store.state.question.description" :elapsedseconds="$store.state.elapsedseconds"></TimeStatus>
      <hr>
      <PetitionStatus :haspetition="$store.state.haspetition"></PetitionStatus>
      <hr>
      <RejectStatus :isrejected="$store.state.isrejected"></RejectStatus>
    </b-container>
</template>

<script>

import RejectStatus from '@/components/RejectStatus'
import PetitionStatus from '@/components/PetitionStatus'
import TimeStatus from '@/components/TimeStatus'

function getBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = error => reject(error);
  });
}

var max_height = 1024
var image_quality = 0.6
function getCompressedBase64(file){
    // Reference: https://github.com/josefrichter/resize/blob/master/public/preprocess.js
    return new Promise((resolve, reject) => {
        let rejectFunc = (error) => {
            reject(error);
        }

        const reader = new FileReader();
        reader.readAsArrayBuffer(file);
        reader.onload = function(event){
            let blob = new Blob([event.target.result]);
            window.URL = window.URL || window.webkitURL;
            let blobURL = window.URL.createObjectURL(blob);
            let image = new Image();
            image.src = blobURL;
            image.onload = function(){
                let width = image.width;
                let height = image.height;
                let aspect_ratio = width / height;

                if(height > max_height){
                    height = max_height;
                    width = Math.round(max_height * aspect_ratio);
                }

                let resized = resize(image, width, height);

                resolve(resized);
            }
            image.onerror = rejectFunc;
        }

        reader.onerror = rejectFunc;
    })
}

function resize(img, width, height){
    let canvas = document.createElement("canvas");
    canvas.width = width;
    canvas.height = height;
    let ctx = canvas.getContext("2d");
    ctx.drawImage(img, 0, 0, width, height);
    return canvas.toDataURL("image/jpeg", image_quality);
}

export default {
  name: "Question",
  components: {
    RejectStatus,
    PetitionStatus,
    TimeStatus
  },  
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
        getCompressedBase64(newFile).then((data) =>{
            this.answer.base64_str = data;
        });
        return;
      }else{
        this.answer.base64_str = null
      }
    }
  },
  computed: {
    getCurrentAnswer(){
      // get previous answer from backend
      return this.$store.state.answer
    },
    getUploadStatus(){
      return this.$store.state.upload_status;
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
    },
    petitionSkipQuestion(evt){
      if(!confirm("確定要請求跳題嗎?(注意, 請求之後不一定就會馬上跳題, 主辦人會視情況允許跳題)")){
        console.log("取消跳題")
        return
      }
      console.log("請願跳題中...")
      this.$store.dispatch("petitionSkipQuestion")
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
