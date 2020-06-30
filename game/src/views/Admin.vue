<template>
  <b-container>
    <b-row>
      <b-col cols="6">
        <b-button @click="getAll">獲取最新組別狀態</b-button>
      </b-col>
      <b-col cols="6">
        <b-button @click="resetAll" variant="danger">重設所有玩家的九宮格和答案(非常危險!)</b-button>
      </b-col>
    </b-row>
    <hr>
    <!-- Loop over each row -->
    <template v-for="(user, index) in users">
        <Board :team="user.account" 
              :gridNumbers="user.gridnumbers" :questionIndex="user.questionindex"
              :question_finished_mask="user.questionfinishedmask"
              :answertext="user.answertext" :answerbase64str="user.answerbase64str"
              :question="questions[index]">
        </Board>
        <hr>
        <TimeStatus v-if="questions[index] !== null" :elapsedseconds="elapsedseconds_s[index]"></TimeStatus>
        <hr>
        <PetitionStatus v-if="questions[index] !== null" :haspetition="haspetitions[index]"></PetitionStatus>
        <hr>
        <RejectStatus v-if="questions[index] !== null" :isrejected="isrejecteds[index]"></RejectStatus>
        <hr>
        <!-- User must have answer some texts or attached an image -->
        <b-row v-if="questions[index] !== null">
            <b-col cols="4"><b-button variant="primary" size="lg" @click="approve($event, user)">核准</b-button></b-col>
            <b-col cols="2"><b-button variant="danger" size="lg" @click="skip($event, user)">跳題</b-button></b-col>
            <b-col cols="2"><h4>(已跳題 {{ numskips[index] }} 次)</h4></b-col>
            <b-col cols="4"><b-button size="lg" @click="reject($event, user)">此題回答錯誤</b-button></b-col>
        </b-row>
        <hr>
    </template>
  </b-container>
</template>

<script>
import Board from '@/components/Board'
import RejectStatus from '@/components/RejectStatus'
import PetitionStatus from '@/components/PetitionStatus'
import TimeStatus from '@/components/TimeStatus'

export default {
  name: "Admin",
  components: {
    Board,
    RejectStatus,
    PetitionStatus,
    TimeStatus,
  },
  data(){
    return {
    }
  },
  methods: {
    getAll(){
      this.$store.dispatch('getAll')
    },
    resetAll(){
      if(!confirm("確定要將所有玩家的 九宮格 清除嗎?(此步無法 undo 唷!)")){
        console.log("[!] 取消 清除 (resetAll)")
        return
      }
      this.$store.dispatch("resetAll")
    },
    approve(evt, user){
      if(!confirm("確定要核准" + user.account + "所回答的答案嗎?")){
        console.log("[*] 取消核准...")
        return
      }
      console.log("[*] Approve " + user.account + "'s answer")
      this.$store.dispatch("approveAnswer", {
        target_account: user.account,
      })

    },
    skip(evt, user){
      if(!confirm("確定要讓" + user.account + "跳題嗎?")){
        console.log("[*] 取消跳題")
        return
      }
      console.log("[*] Skip " + user.account + "'s answer")
      this.$store.dispatch("skipAnswer", {
        target_account: user.account,
      })

    },
    reject(evt, user){
      if(!confirm("確定要拒絕 " + user.account + "的回答嗎?")){
        console.log("[*] 取消拒絕")
        return
      }
      console.log("[*] Reject " + user.account + "'s answer")
      this.$store.dispatch("rejectAnswer", {
        target_account: user.account,
      })

    }
  },
  computed: {
    users(){
      return this.$store.state.users
    },
    questions(){
      return this.$store.state.questions
    },
    isrejecteds(){
      return this.$store.state.isrejecteds
    },
    haspetitions(){
      return this.$store.state.haspetitions
    },
    numskips(){
      return this.$store.state.numskips
    },
    elapsedseconds_s(){
      return this.$store.state.elapsedseconds_s
    }
  }
}
</script>

<style>
</style>
