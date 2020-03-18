<template>
  <b-container>
    <b-row>
      <b-col cols="12">
        <b-button @click="getAll">獲取最新組別狀態</b-button>
      </b-col>
    </b-row>
    <hr>
    <!-- Loop over each row -->
    <template v-for="(user, index) in users">
        <Board :team="user.account" 
              :gridNumbers="user.gridnumbers" :questionIndex="user.questionindex"
              :question_finished_mask="user.questionfinishedmask"
              :answertext="user.answertext" :answerbase64str="user.answerbase64str">
        </Board>
        <hr>
        <!-- User must have answer some texts or attached an image -->
        <b-row v-if="user.answertext || user.answerbase64str">
            <b-col cols="6"><b-button variant="primary" size="lg" @click="approve($event, user)">核准</b-button></b-col>
            <b-col cols="6"><b-button variant="danger" size="lg" @click="skip($event, user)">跳題</b-button></b-col>
        </b-row>
        <hr>
    </template>
  </b-container>
</template>

<script>
import Board from '@/components/Board'

export default {
  name: "Admin",
  components: {
    Board,
  },
  data(){
    return {
    }
  },
  methods: {
    getAll(){
      this.$store.dispatch('getAll')
    },
    approve(evt, user){
      console.log("[*] Approve " + user.account + "'s answer")
      this.$store.dispatch("approveAnswer", {
        target_account: user.account,
      })
    },
    skip(evt, user){
      console.log("[*] Skip " + user.account + "'s answer")
      this.$store.dispatch("skipAnswer", {
        target_account: user.account,
      })
    }
  },
  computed: {
    users(){
      return this.$store.state.users
    }
  }
}
</script>

<style>
</style>
