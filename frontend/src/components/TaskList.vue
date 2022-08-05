<template>
  <el-drawer v-model="props.visible" :show-close="false" size="40%" @close="closeDrawer">
    <template #default>
      <el-tabs v-model="activeName" class="demo-tabs" @tab-click="handleClick">
        <el-tab-pane label="Upload" name="upload">
          <el-table :data="uploadListData" style="width: 100%">

            <el-table-column width="250">
              <template #default="scope">
                <span>{{ scope.row.key }}</span>
                <el-progress :percentage="percentageMap[scope.row.taskId]" :color="colors"/>
              </template>
            </el-table-column>

            <el-table-column prop="humanSize" width="100"/>

            <el-table-column fixed="right" align="right">
              <template #default="scope">
                <el-button
                    type="success"
                    @click="resumeUpload(uploadListData[scope.$index])"
                >
                  <el-icon>
                    <CaretRight/>
                  </el-icon>
                </el-button>
                <el-button
                    type="danger"
                    plain
                    @click="cancelUpload(uploadListData[scope.$index], scope.$index)"
                >
                  <el-icon>
                    <Delete/>
                  </el-icon>
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="Download" name="download">

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
import {computed, ref} from 'vue'
import {useStore} from "vuex";
import {CancelUploadTask, ResumeUploadTask} from "../../wailsjs/go/app/App";
import {ElMessage} from "element-plus";
import {db} from '../../wailsjs/go/models'

export default {
  name: "TaskList",
  props: {
    visible: {
      type: Boolean,
      default: false,
    }
  },
  setup(props, context) {
    const activeName = ref('upload')
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

    const colors = [
      {color: '#1989fa', percentage: 100},
      {color: '#5cb87a', percentage: 101},
    ]

    const percentageMap = computed(() => {
      return store.state.uploadProgress
    })

    const resumeUpload = (task) => {
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
        if (res.err !== '') {
          ElMessage.error(res.err)
        } else {
          const payload = {
            data: 100,
            progress: t.taskId
          }
          store.commit('updateProgress', payload)
          listObjects(bucketName, '', prefix.value, 100)
        }
      })
    }

    const cancelUpload = (task, index) => {
      console.log(task.taskId)
      if (percentageMap[task.taskId] !== 100) {
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
        CancelUploadTask(t).then((res) => {
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

    return {
      props,
      colors,
      activeName,
      uploadListData,
      percentageMap,
      cancelUpload,
      cancelClick,
      resumeUpload,
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