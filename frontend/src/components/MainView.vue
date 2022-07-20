<template>
  <MainMenu/>
  <TaskList/>
  <div style="margin-top: 10px">
    <router-view/>
  </div>
</template>

<script>
import MainMenu from  "./MainMenu.vue"
import TaskList from "./TaskList.vue";

export default {
  components: {MainMenu, TaskList},
}
</script>

<style scoped>

</style>