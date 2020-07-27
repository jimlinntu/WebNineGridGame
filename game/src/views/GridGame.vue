<template>
    <b-container>
        <b-row v-if="hasGridNumbers" class="justify-content-md-center">
            <b-col cols="12"><b-button @click="refetch">更新九宮格狀態</b-button></b-col>
            <b-col cols="12"><hr></b-col>
            <b-col cols="12"><h1>我們的九宮格:</h1></b-col>
        </b-row>
        <b-row v-if="hasGridNumbers" class="justify-content-md-center">
            <h3>我們目前選的題目是: 第 {{ questionIndex }} 號(金黃色) (顯示 ? 代表你還沒選下一題) </h3> 
        </b-row>
        <b-row v-if="hasGridNumbers">
            <!-- add onclick to expand question -->
            <b-col class="slot" :class="{
                    'question_index': index === $store.state.questionIndex,
                    'question_finished': $store.state.question_finished_mask[index]
                }" cols="3" v-for="(number, index) in $store.state.gridNumbers" :key="number" @click="selectQuestion($event, index)">
                {{ number }}
            </b-col>
        </b-row>
        <hr v-if="hasGridNumbers">
        <Question v-if="hasGridNumbers"></Question>

        <!---->
        <b-row v-if=!hasGridNumbers class="text-center">
            <b-col cols="12"><h1>請填入九宮格:</h1></b-col>
        </b-row>
        <b-row v-if="!hasGridNumbers">
            <b-col class="slot" cols="3" v-for="i in num_grid" :key="i">
                <draggable class="draggable_slot" v-model="slots[i-1]" group="people" 
                        :move="checkNumberMove">
                    <div class="slots" v-for="element in slots[i-1]" :key="element">{{ element }}</div>
                </draggable>
            </b-col>
        </b-row>
        <hr v-if="!hasGridNumbers">
        <b-row v-if="!hasGridNumbers">
            <b-col cols="12"><h1>候選數字:</h1></b-col>
        </b-row>
        <b-row v-if="!hasGridNumbers">
            <draggable v-model="unselectedNumbers" group="people" class="unselected_draggable_numbers"
                        :move="checkNumberMove">
                <span class="unselectedNum" v-for="element in unselectedNumbers" :key="element">{{ element }}</span>
            </draggable>
        </b-row>
        <hr v-if="!hasGridNumbers">
        <b-row v-if="!hasGridNumbers">
            <b-col cols="12"><b-button size="lg" @click.prevent="submitNineGrids">提交九宮格</b-button></b-col>
        </b-row>
    </b-container>
</template>

<script>
import draggable from 'vuedraggable'
import Question from '@/components/Question'

var num_grid = 16;

export default {
  name: 'GridGame',
  components: {
    draggable,
    Question
  },
  data () {
    let unselectedNumbers = []
    let slots = []
    for(let i = 0; i < 20; i++){
      unselectedNumbers.push(i+1);
    }
    for(let i = 0; i < num_grid; i++){
        slots.push([])
    }
    return {
      unselectedNumbers: unselectedNumbers,
      slots: slots,
      num_grid: num_grid
    }
  },
  computed: {
    hasGridNumbers(){
        if(this.$store.state.gridNumbers.length === num_grid){
            return true
        }else return false
    },
    questionIndex(){
        let index = this.$store.state.questionIndex
        if(index == -1){
            return "?"
        }
        let gridNumbers = this.$store.state.gridNumbers
        if(gridNumbers === undefined || gridNumbers.length !== num_grid){
            return "?"
        }
        return gridNumbers[index]
    }
  },
  methods: {
    checkNumberMove(evt){
        // check if this slot (list) contains X
        let targetList = evt.relatedContext.list
        let targetComponent = evt.relatedContext.component
        // if the target already contains one number and it is a slot, disable this move
        if(targetList.length == 1 && targetComponent.$el.className === "draggable_slot"){
            return false
        }
        return true
    },
    submitNineGrids(event){
        // check num_grid grids are full
        window.console.log("Hello")
        let grids_are_full = true
        let selectedNumbers = []
        for(let i = 0; i < num_grid; i++){
            // if one of the slots is empty, grids are not full
            if(this.slots[i].length !== 1){
                grids_are_full = false
                break
            }
            else selectedNumbers.push(this.slots[i][0])
        }

        if(grids_are_full){
            // dispatch action to Vuex
            this.$store.dispatch("submitGridNumbers", {
                    gridNumbers: selectedNumbers
                })
            
        }else{
            alert("請將九宮格填滿後再提交!")
        }
    },
    selectQuestion(evt, index){
        // if he haven't choose a question
        let question_finished_mask = this.$store.state.question_finished_mask
        if(question_finished_mask.length !== num_grid){
            console.log("[*] question_finished_mask looks weird: ", question_finished_mask)
            return
        }
        if(this.$store.state.questionIndex !== -1){
            console.log("[*] This question have not yet finished, you cannot choose another question!")
            return
        }
        let gridNumbers = this.$store.state.gridNumbers
        // question_finished_mask[index] must be false
        if(question_finished_mask[index] === false && confirm("你確定要選("+ gridNumbers[index] + ")號嗎? (選完這題, 在解完這題之前是不能換題的!)")){
            // send this question index to the backend server
            console.log("[*] Sending selected index to the backend server")
            this.$store.dispatch("selectQuestion", {
                    questionIndex: index
                })
        }else{
            console.log("[!] Cancel selectQuestion function...")
        }
    },
    refetch(){
        // refetch all information
        this.$store.dispatch("getGridNumbers")
    }
  }
}
</script>

<style scoped>
.slot {
    border-style: solid;
    height: 50px;
}

.draggable_slot {
    width: 100%;
    height: 100%;
}

.slots {
    position: absolute;
    height: 100%;
    top: 20%;
    left: 50%;
    cursor: grab;
}

.unselected_draggable_numbers {
    width: 100%;
}

.unselectedNum {
    margin-left: 10px;
    margin-right: 10px;
    padding: 5px;
    border: double;
    cursor: grab;
    display: inline-block;
}

.question_index{
    color: rgb(255,165,0);
    background-color: rgb(255, 165, 0, 0.3);
}

.question_finished{
    color: rgb(0, 128, 0);
    background-color: rgb(0, 128, 0, 0.3);
}

</style>
