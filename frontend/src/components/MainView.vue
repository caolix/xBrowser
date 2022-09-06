<template>
  <MainMenu @drawerVisible="changeVisible" @settingsVisible="changeSettingsVisible"/>
  <TaskList :visible="visible" @cancelVisible="changeVisible" :tabName="tabName"/>
  <Settings :visible="settingsVisible" @cancelVisible="changeSettingsVisible"/>
  <div style="margin-top: 10px">
    <router-view @changeVisible="changeVisible" @setTaskTabName="setTaskTabName"/>
  </div>
</template>

<script>
import MainMenu from "./MainMenu.vue"
import TaskList from "./TaskList.vue";
import Settings from "./Settings.vue";
import {onMounted, ref} from "vue";
import {LoadAllDownloadTasks, LoadAllUploadTasks} from "../../wailsjs/go/app/App";
import {EventsOn, LogDebug} from "../../wailsjs/runtime";
import {useStore} from "vuex";

export default {
  components: {MainMenu, TaskList, Settings},
  setup() {
    const visible = ref(false)
    const settingsVisible = ref(false)
    const changeVisible = (bool) => {
      visible.value = bool
    }
    const changeSettingsVisible = (bool) => {
      settingsVisible.value = bool
    }
    const tabName = ref("upload")
    const setTaskTabName = (name) => {
      tabName.value = name
    }
    const store = useStore()

    onMounted(() => {
      LoadAllUploadTasks().then((tasks) => {
        if (tasks.length !== 0) {
          tasks.forEach((task, i) => {
            if (task.status === 1 || task.status === 2) {
              task.status = 0 // PAUSE
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
              accountId: task.accountId,
              isMultipart: task.isMultipart,
              partSize: task.partSize,
              taskId: task.taskId,
              uploadId: task.uploadId
            }
            const payload = {
              file: file,
              progress: task.taskId
            }

            console.log("Load task:" + file.bucket + "/" + file.key + " status: " + file.status)
            store.commit('addToUploadList', payload)

            const progressPayload = {
              data: task.uploadedSize,
              progress: task.taskId
            }
            store.commit('updateUploadProgress', progressPayload)

            EventsOn(task.taskId, (data) => {
              const progressPayload = {
                data: data,
                progress: task.taskId
              }
              store.commit('updateUploadProgress', progressPayload)
            })
          })
        }
      })

      LoadAllDownloadTasks().then((tasks) => {
        if (tasks.length !== 0) {
          tasks.forEach((task, i) => {
            if (task.status === 0) {
              task.status = 1
            }
            const file = {
              name: task.name,
              bucket: task.bucket,
              key: task.key,
              dest: task.dest,
              size: task.size,
              humanSize: task.humanSize,
              status: task.status,
              progress: task.taskId,
              accountId: task.accountId,
              taskId: task.taskId,
            }
            const payload = {
              file: file,
              progress: task.taskId
            }

            LogDebug("Load task:" + file.bucket + "/" + file.key + " status: " + file.status)
            store.commit('addToDownloadList', payload)

            const progressPayload = {
              data: 0,
              progress: task.taskId
            }
            store.commit('updateDownloadProgress', progressPayload)

            EventsOn(task.taskId, (data) => {
              const progressPayload = {
                data: data,
                progress: task.taskId
              }
              store.commit('updateDownloadProgress', progressPayload)
            })

          })

        }
      })
    })


    return {
      visible,
      settingsVisible,
      tabName,
      changeVisible,
      changeSettingsVisible,
      setTaskTabName
    }
  }
}
</script>

<style scoped>

</style>