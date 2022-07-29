<template>
  <MainMenu @drawerVisible="changeVisible"/>
  <TaskList :visible="visible" @cancelVisible="changeVisible"/>
  <div style="margin-top: 10px">
    <router-view @changeVisible="changeVisible"/>
  </div>
</template>

<script>
import MainMenu from  "./MainMenu.vue"
import TaskList from "./TaskList.vue";
import {ref} from "vue";

export default {
  components: {MainMenu, TaskList},
  setup() {
    const visible = ref(false)
    const changeVisible = (bool) => {
      visible.value = bool
    }
    return {
      visible,
      changeVisible
    }
  }
}
</script>

<style scoped>

</style>