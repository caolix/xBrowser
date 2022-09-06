<template>
  <el-drawer v-model="props.visible" :show-close="false" size="40%" @close="closeDrawer">
    <template #default>
      <el-tabs v-model="props.tabName" class="demo-tabs" @tab-click="handleClick">
        <div style="text-align: right;margin-bottom: 8px">
          <el-button
              plain
              type="info"
              @click="clearUploadFinished"
          >
            <el-icon color="#67C23A" style="margin-right: 3px">
              <Finished/>
            </el-icon>
            Clear Finished
          </el-button>
        </div>
        <el-tab-pane name="upload">
          <template #label>
            <span>Upload</span>
            <el-badge :value="uploadListCount" class="item">
            </el-badge>

          </template>
          <div class="infinite-list-wrapper" style="overflow: auto">
            <ul
                v-infinite-scroll="load"
                class="list"
                :infinite-scroll-disabled="disabled"
            >
              <li v-for="(task, i) in uploadListData" :key="i" class="list-item">
                <el-row>
                  <el-col :span="12" class="hidden-name">{{ task.name }}
                    <el-progress :percentage="uploadPercentageMap[task.taskId]" :color="colors"/>
                  </el-col>
                  <el-col :span="4" class="hidden-name">{{ task.humanSize }}</el-col>
                  <el-col :span="8" style="text-align: right">
                    <el-button
                        v-if="task.status===0"
                        plain
                        type="success"
                        @click="resumeUpload(task, i)"
                        style="width: 36px;"
                    >
                      <el-icon color="#67C23A">
                        <CaretRight/>
                      </el-icon>
                    </el-button>
                    <el-button
                        type="info"
                        plain
                        @click="removeUpload(task, i)"
                        style="width: 36px;"
                    >
                      <el-icon color="#F56C6C">
                        <CloseBold/>
                      </el-icon>
                    </el-button>
                  </el-col>
                </el-row>
              </li>
            </ul>
            <p v-if="loading" style="color: black">Loading...</p>
            <p v-if="noMore"></p>
          </div>
          <div style="text-align: left">
            <span style="color: black">Data: {{ uploadTotalSize }}</span>
            <span style="color: red; margin-left: 10px">Failed: 0</span>
          </div>
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

            <el-table-column fixed="right" align="right">
              <template #default="scope">
                <el-button
                    v-if="downloadListData[scope.$index].status===1"
                    type="success"
                    @click="resumeDownload(downloadListData[scope.$index], scope.$index)"
                >
                  <el-icon>
                    <CaretRight/>
                  </el-icon>
                </el-button>
                <el-button
                    type="danger"
                    plain
                    @click="removeDownload(downloadListData[scope.$index], scope.$index)"
                >
                  <el-icon>
                    <Delete/>
                  </el-icon>
                </el-button>
              </template>
            </el-table-column>
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

<script lang="ts">
import {computed, ref} from 'vue'
import {useStore} from "vuex";
import {RemoveDownloadTask, RemoveUploadTask, ResumeUploadTask} from "../../wailsjs/go/app/App";
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

    const uploadSuccessCount = computed(() => {
      return store.state.uploadSuccessCount
    })

    const uploadListCount = computed(() => {
      return uploadSuccessCount.value + "/" + store.state.uploadList.length
    })

    const uploadTotalSize = computed(() => {
      return humanReadableFilesize(store.state.uploadTotalSize)
    })

    const InitListItemCount = ref(20)
    var hasLoadedOnce = false
    const listCount = ref(20)

    const listItemCount = computed(() => {
      if (!hasLoadedOnce) {
        if (InitListItemCount.value >= store.state.uploadList.length) {
          listCount.value = store.state.uploadList.length
          return listCount
        }
        return InitListItemCount
      }
      return listCount
    })

    const loading = ref(false)
    const noMore = computed(() => listCount.value >= store.state.uploadList.length)
    const disabled = computed(() => loading.value || noMore.value)
    const load = () => {
      if (!hasLoadedOnce) {
        hasLoadedOnce = true
      }
      loading.value = true
      setTimeout(() => {
        listCount.value += 20
        console.log("load", listCount.value)
        loading.value = false
      }, 1000)
    }

    const uploadListData = computed(() => {
      var v = listItemCount.value.value
      console.log("listv:", listItemCount.value.value)
      return store.state.uploadList.slice(0, v)
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

    const clearUploadFinished = () => {
      store.commit('clearUploadFinished')
    }

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
        }
      })
    }

    const resumeDownload = (task, index) => {

    }

    const removeUpload = (task, index) => {
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
      store.commit('removeUploadListParam', payload)
    }

    const removeDownload = (task, index) => {
      if (downloadPercentageMap[task.taskId] !== 100) {
        var t = new db.DownloadTask()
        t.key = task.key
        t.size = task.size
        t.name = task.name
        t.accountId = task.accountId
        t.bucket = task.bucket
        t.humanSize = task.humanSize
        t.dest = task.dest
        t.status = task.status
        t.taskId = task.taskId
        RemoveDownloadTask(t).then((res) => {
          if (res.err !== '') {
            ElMessage.error(res.err)
          }
        })
      }

      const payload = {
        index: index,
        progress: task.taskId
      }
      store.commit('removeDownloadListParam', payload)
    }

    function humanReadableFilesize(size) {
      var units = new Array("B", "KB", "MB", "GB", "TB", "PB");
      var mod = 1024.0;
      var i = 0;
      while (size >= mod) {
        size /= mod;
        i++;
      }
      //return Math.round(size) + units[i];
      return formatNum(size, 1) + units[i];
    }

    //格式化数字类型,保留小数点后几位,非四舍五入
    //size:值
    //n:保留位数
    function formatNum(size, n) {
      var sizeStr = size.toString();
      if (sizeStr.lastIndexOf('.') > -1)
        return sizeStr.substring(0, sizeStr.toString().indexOf('.') + 1 + n);
      else
        return sizeStr;
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
      clearUploadFinished,
      uploadSuccessCount,
      uploadListCount,
      uploadTotalSize,
      listCount,
      loading,
      noMore,
      disabled,
      load,
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

.item {
  margin-left: 3px;
}

.hidden-name {
  text-align: left;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
}

.infinite-list-wrapper {
  height: 800px;
  width: 100%;
}

.infinite-list-wrapper .list {
  padding: 0;
  margin: 0;
  list-style: none;
}

.infinite-list-wrapper .list-item {
  display: block;
  height: 40px;
  color: black;
  border-bottom: 1px solid #ddd
}

.infinite-list-wrapper .list-item + .list-item {
  margin-top: 0px;
}
</style>