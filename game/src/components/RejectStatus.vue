<template>
    <b-row class="text-center">
      <b-col cols="12"><h3>答案審核狀態(還未審核 或 答案錯誤):</h3></b-col>
      <b-col cols="12" v-if="isrejected === false"><b-icon-question-square></b-icon-question-square>還未審核</b-col>
      <b-col cols="12" v-else><b-icon-x-circle></b-icon-x-circle>答案錯誤</b-col>
    </b-row>
</template>

<script>
export default {
  name: "RejectStatus",
  props: {
    isrejected: Boolean,
  },
  data(){
    return {
    }
  }
}
</script>
<style>
</style>
