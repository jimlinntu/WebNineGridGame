<template>
    <b-container>
        <b-row class="justify-content-md-center">
            <h1>請填入以下九宮格:</h1>
        </b-row>
        <b-row>
            <b-col class="slot" cols="4" v-for="i in 9" :key="i">
                <draggable class="draggable_slot" v-model="slots[i-1]" group="people" 
                        :move="checkNumberMove">
                    <div class="slots" v-for="element in slots[i-1]" :key="element">{{ element }}</div>
                </draggable>
            </b-col>
        </b-row>
        <hr>
        <b-row class="justify-content-md-center">
            <h1>候選數字:</h1>
        </b-row>
        <b-row>
            <draggable v-model="unselectedNumbers" group="people" class="unselected_draggable_numbers"
                        :move="checkNumberMove">
                <span class="unselectedNum" v-for="element in unselectedNumbers" :key="element">{{ element }}</span>
            </draggable>
        </b-row>
        <hr>
        <b-row class="justify-content-md-center">
            <b-button size="lg" @click.prevent="submitNineGrids">提交九宮格</b-button>
        </b-row>
    </b-container>
</template>

<script>
import draggable from 'vuedraggable'

export default {
  name: 'GridGame',
  components: {
    draggable
  },
  data () {
    let unselectedNumbers = []
    let slots = []
    for(let i = 0; i < 15; i++){
      unselectedNumbers.push(i+1);
    }
    for(let i = 0; i < 9; i++){
        slots.push([])
    }
    return {
      selectedNumbers: [],
      unselectedNumbers: unselectedNumbers,
      slots: slots
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
        // check 9 grids are full
        window.console.log("Hello")
        let grids_are_full = true
        for(let i = 0; i < 9; i++){
            if(this.slots[i].length !== 1){
                grids_are_full = false
                break
            }
        }
        if(grids_are_full){
            // send POST request to the backend server
            this.axios.get(this.$store.state.backend_url + "ping").then((response) =>{
                console.log(response)
            })
        }else{
            alert("請將九宮格填滿後再提交!")
        }
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

</style>
