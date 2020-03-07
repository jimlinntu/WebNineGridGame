<template>
    <b-container>
        <b-row>
            <b-col class="slot" cols="4" v-for="i in 9" :key="i">
                <draggable class="draggable_slot" v-model="slots[i-1]" group="people" 
                        @start="drag=true" @end="drag=false"
                        @add="add">
                    <div class="slot_list" v-for="element in slots[i-1]" :key="element">{{ element }}</div>
                </draggable>
            </b-col>
        </b-row>
        --------------
        <b-row>
            <draggable v-model="unselectedNumbers" group="people" 
                        @start="drag=true" @end="drag=false"
                        :move="checkSlotAdd">

                <div v-for="element in unselectedNumbers" :key="element">{{ element }}</div>
            </draggable>
        </b-row>
    </b-container>
</template>

<script>
import draggable from 'vuedraggable'

var emptySign = "X"

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
    checkSlotAdd(evt){
        // check if this slot (list) contains X
        let targetList = evt.relatedContext.list
        window.console.log(evt)
        if(targetList.length == 1){
            // it must be our 9 grids
            if(targetList[0] === emptySign){
                window.console.log("[*] This is a slot")
                // remove emptySign
                
                return true
            }
        }
        return true
    },
    add(evt){
        window.console.log(evt)
    }
  }
}
</script>

<style>
.slot {
    border-style: solid;
    height: 50px;
}

.draggable_slot {
    width: 100%;
    height: 100%;
}

.slot_list {
    position: absolute;
    height: 100%;
    top: 20%;
    left: 50%;
}
</style>
