<template>
  <el-drawer v-model="props.visible" :show-close="false" size="40%" @close="closeDrawer">
    <template #default>
      <el-tabs v-model="props.tabName" class="demo-tabs" @tab-click="handleClick">
<!--        upload pane -->
        <el-tab-pane label="Upload" name="upload">
          <el-table :data="uploadListData" style="width: 100%">
            <el-table-column width="250">
              <template #default="scope">
                <span>{{ scope.row.key }}</span>
                <el-progress :percentage="uploadPercentageMap[scope.row.taskId]" :color="colors"/>
              </template>
            </el-table-column>

            <el-table-column prop="humanSize" width="100"/>

            <el-table-column fixed="right" align="right">
              <template #default="scope">
                <el-button
                    v-if="uploadListData[scope.$index].status===1"
                    type="success"
                    @click="resumeUpload(uploadListData[scope.$index], scope.$index)"
                >
                  <el-icon>
                    <CaretRight/>
                  </el-icon>
                </el-button>
                <el-button
                    type="danger"
                    plain
                    @click="removeUpload(uploadListData[scope.$index], scope.$index)"
                >
                  <el-icon>
                    <Delete/>
                  </el-icon>
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!--        download pane -->
        <el-tab-pane label="Download" name="download">
          <el-table :data="downloadListData" style="width: 100%">
            <el-table-column width="250">
              <template #default="scope">
                <span>{{ scope.row.key }}</span>
                <el-progress :percentage="downloadPercentageMap[scope.row.taskId]" :color="colors"/>
              </template>
            </el-table-column>
            <el-table-column prop="humanSize" width="100"/>
          </el-table>
        </el-tab-pane>
      </el-tabs>

    </template>
    <template #footer>
      <div style="flex: auto">
        <el-button @click="cancelClick">cancel</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script>
import {computed} from 'vue'
import {useStore} from "vuex";
import {RemoveUploadTask, ResumeUploadTask} from "../../wailsjs/go/app/App";
import {ElMessage} from "element-plus";
import {db} from '../../wailsjs/go/models'

export default {
  name: "TaskList",
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    tabName: {
      type: String,
      default: "upload",
    }
  },
  setup(props, context) {
    const store = useStore()
    const cancelClick = () => {
      context.emit('cancelVisible', false)
    }
    const closeDrawer = () => {
      context.emit('cancelVisible', false)
    }

    const uploadListData = computed(() => {
      return store.state.uploadList
    })

    const downloadListData = computed(() => {
      return store.state.downloadList
    })

    const colors = [
      {color: '#1989fa', percentage: 100},
      {color: '#5cb87a', percentage: 101},
    ]

    const uploadPercentageMap = computed(() => {
      return store.state.uploadProgress
    })

    const downloadPercentageMap = computed(() => {
      return store.state.downloadProgress
    })

    const resumeUpload = (task, index) => {
      var t = new db.UploadTask()
      t.uploadedSize = task.uploadedSize
      t.key = task.key
      t.size = task.size
      t.name = task.name
      t.accountId = task.accountId
      t.bucket = task.bucket
      t.humanSize = task.humanSize
      t.isMultipart = task.isMultipart
      t.partSize = task.partSize
      t.source = task.source
      t.status = task.status
      t.taskId = task.taskId
      t.uploadId = task.uploadId
      ResumeUploadTask(t).then((res) => {
        console.log(t.status)
        if (res.err !== '') {
          ElMessage.error(res.err)
        } else {
          const payload = {
            data: 100,
            progress: t.taskId
          }
          store.commit('updateUploadProgress', payload)
          const statusPayload = {
            index: index,
            status: 3, // finish
          }
          store.commit('updateUploadStatus', statusPayload)
          listObjects(bucketName, '', prefix.value, 100)
        }
      })
    }

    const resumeDownload = (task, index) => {

    }

    const removeUpload = (task, index) => {
      console.log(task.taskId)
      if (uploadPercentageMap[task.taskId] !== 100) {
        var t = new db.UploadTask()
        t.uploadedSize = task.uploadedSize
        t.key = task.key
        t.size = task.size
        t.name = task.name
        t.accountId = task.accountId
        t.bucket = task.bucket
        t.humanSize = task.humanSize
        t.isMultipart = task.isMultipart
        t.partSize = task.partSize
        t.source = task.source
        t.status = task.status
        t.taskId = task.taskId
        t.uploadId = task.uploadId
        RemoveUploadTask(t).then((res) => {
          if (res.err !== '') {
            ElMessage.error(res.err)
          }
        })
      }

      const payload = {
        index: index,
        progress: task.taskId
      }
      console.log("payload:" + payload.index + payload.progress)
      store.commit('removeUploadListParam', payload)
    }

    const removeDownload = (task, index) => {

    }

    return {
      props,
      colors,
      uploadListData,
      downloadListData,
      uploadPercentageMap,
      downloadPercentageMap,
      removeUpload,
      removeDownload,
      cancelClick,
      resumeUpload,
      resumeDownload,
      closeDrawer
    }
  }
}
</script>

<style scope>
.task {
  position: fixed;
  bottom: 0px;
  right: 0px;
}

.el-drawer {
  --el-transition-duration: 0;
  transition: null;
}
</style>