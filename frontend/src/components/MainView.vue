<template>
  <MainMenu @drawerVisible="changeVisible"/>
  <TaskList :visible="visible" @cancelVisible="changeVisible"/>
  <div style="margin-top: 10px">
    <router-view @changeVisible="changeVisible"/>
  </div>
</template>

<script>
import MainMenu from "./MainMenu.vue"
import TaskList from "./TaskList.vue";
import {ref} from "vue";
import {LoadAllUploadTasks} from "../../wailsjs/go/app/App";
import {EventsOn, LogDebug} from "../../wailsjs/runtime";
import {useStore} from "vuex";

export default {
  components: {MainMenu, TaskList},
  setup() {
    const visible = ref(false)
    const changeVisible = (bool) => {
      visible.value = bool
    }
    const store = useStore()
    LoadAllUploadTasks().then((tasks) => {
      if (tasks.length !== 0) {
        tasks.forEach((task, i) => {
          if (task.status === 0) {
            task.status = 1
          }
          const file = {
            name: task.name,
            bucket: task.bucket,
            key: task.key,
            source: task.source,
            uploadedSize: task.uploadedSize,
            size: task.size,
            humanSize: task.humanSize,
            status: task.status,
            progress: task.taskId,
            accountId:  task.accountId,
            isMultipart: task.isMultipart,
            partSize: task.partSize,
            taskId: task.taskId,
            uploadId: task.uploadId
          }
          const payload = {
            file: file,
            progress: task.taskId
          }

          LogDebug("Load task:" + file.bucket + "/" + file.key)
          store.commit('addToUploadList', payload)

          const progressPayload = {
            data: task.uploadedSize,
            progress: task.taskId
          }
          store.commit('updateProgress', progressPayload)

          EventsOn(task.taskId, (data) => {
            const progressPayload = {
              data: data,
              progress: task.taskId
            }
            store.commit('updateProgress', progressPayload)
          })

        })


      }
    })
    return {
      visible,
      changeVisible
    }
  }
}
</script>

<style scoped>

</style>