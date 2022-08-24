<template>
  <div class="container2">
    <el-button type="primary" @click="backward">
      <el-icon>
        <ArrowLeftBold/>
      </el-icon>
      Back
    </el-button>
    <el-input v-model="showPath" width="80%"></el-input>
  </div>
  <div class="container2">
    <div style="text-align: left; width: 50%">
      <el-dropdown split-button type="primary" @click="selectUploadObjects" @command="handleUploadCommand">
        <el-icon style="padding-right: 6px">
          <UploadFilled/>
        </el-icon>
        Upload
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="UploadFiles">
              <el-icon color="#409EFF">
                <DocumentAdd/>
              </el-icon>
              Upload Files
            </el-dropdown-item>
            <el-dropdown-item command="UploadFolder">
              <el-icon color="#409EFF">
                <FolderAdd/>
              </el-icon>
              Upload Folder
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <el-button type="primary" @click="dialogCreateDirVisible = true">
        <el-icon style="padding-right: 6px">
          <FolderAdd/>
        </el-icon>
        Create Folder
      </el-button>

      <!--      More-->
      <el-dropdown class="more-button" trigger="click" @command="handleMoreCommand">
        <el-button color="#606266" :dark="isDark" plain :disabled="disableMoreButton">
          More
          <el-icon class="el-icon--right">
            <arrow-down/>
          </el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="RemoveFiles">
              <el-icon style="bottom: 1px" color="#F56C6C">
                <CloseBold/>
              </el-icon>
              Remove
            </el-dropdown-item>
            <el-dropdown-item command="DownloadFiles">
              <el-icon style="bottom: 1px" color="#67C23A">
                <Download/>
              </el-icon>
              Download
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

    </div>
    <div style="text-align: right;  width: 50%">
      <el-input class="input-with-select" placeholder="Search by prefix" v-model="subPrefix">
        <template #append>
          <el-button type="primary" @click="searchBySubPrefix">
            <el-icon>
              <Search/>
            </el-icon>
          </el-button>
        </template>
      </el-input>
    </div>
  </div>

  <!--Dir Dialog-->
  <el-dialog v-model="dialogCreateDirVisible" title="Create Directory">
    <span>Folder Name:</span>
    <el-input v-model="newFolderName" autocomplete="off"/>
    <template #footer>
      <span class="dialog-footer">
         <el-button type="primary" @click="createDir"
         >Create</el-button
         >
        <el-button @click="dialogCreateDirVisible = false">Cancel</el-button>
      </span>
    </template>
  </el-dialog>

  <!--GetObject Dialog-->
  <el-dialog v-model="dialogFormVisible" title="Progress">
    <el-progress type="dashboard" :percentage="percentage" :color="colors">
      <template #default="{ percentage }">
        <span class="percentage-value">{{ percentage }}%</span>
        <span class="percentage-label">{{ percentageLabel }}</span>
      </template>
    </el-progress>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogFormVisible = false">Cancel</el-button>
      </span>
    </template>
  </el-dialog>

  <!--Delete Dialog-->
  <el-dialog v-model="dialogDeleteVisible" title="Delete">
    <span>Delete Process:</span>
    <el-progress
        :text-inside="true"
        :stroke-width="20"
        :percentage="deletePercentage"
        :color="colors"
    >
      <span>{{ deletePercentageLabel }} </span>
    </el-progress>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogDeleteVisible = false">Cancel</el-button>
      </span>
    </template>
  </el-dialog>

  <!-- Table Data -->
  <div class="container">
    <el-table
        :data="objectList"
        max-height="800"
        class="table"
        v-loading="loading"
        ref="multipleTableRef"
        @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55"/>
      <el-table-column prop="key" label="Name" width="400">
        <template #default="scope">
          <el-icon style="top:4px" v-if="tableData[scope.$index].type==='Folder'">
            <FolderOpened/>
          </el-icon>
          <el-icon style="top:4px" v-else>
            <Document/>
          </el-icon>
          <el-tooltip
              class="box-item"
              effect="dark"
              placement="right"
          >
            <template #content>
              {{ tableData[scope.$index].key }}
            </template>
            <el-button
                link
                type="primary"
                size="small"
                @click="clickKey(tableData[scope.$index])"
            >
              {{ tableData[scope.$index].showed }}
            </el-button>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column prop="type" label="Type"></el-table-column>
      <el-table-column prop="humanSize" label="Size"></el-table-column>
      <el-table-column label="Operations" fixed="right" align="right" width="200">
        <template #default="scope">
          <el-button size="small" @click="getObject(tableData[scope.$index].key)"
                     v-if="tableData[scope.$index].type==='Object'">
            Download
          </el-button>
          <el-popconfirm
              confirm-button-text="Yes"
              cancel-button-text="No"
              :icon="QuestionFilled"
              icon-color="#FF0000"
              title="Are you sure to delete this object?"
              @confirm="deleteObject(tableData[scope.$index].key, tableData[scope.$index].type)"
          >
            <template #reference>
              <el-button
                  size="small"
                  type="danger"
              >
                Remove
              </el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
  </div>

  <!--  Page -->
  <el-button-group>
    <el-button type="primary" :icon="ArrowLeft" :disabled="markerStack.length === 0" @click="toPreviousPage">Previous
      Page
    </el-button>
    <el-button type="primary" :disabled="!isTruncate" @click="toNextPage">
      Next Page
      <el-icon class="el-icon--right">
        <ArrowRight/>
      </el-icon>
    </el-button>
  </el-button-group>
</template>

<script lang="ts">

import {useRoute, useRouter} from "vue-router";
import {computed, onMounted, reactive, ref, toRefs} from "vue";
import {
  DeleteObject,
  DeleteObjects,
  DoPutObject,
  DoUploadFolder,
  GetObject,
  ListObjects,
  PutDir,
  SelectDownloadPath,
  SelectUploadFiles,
  SelectUploadFolder
} from "../../wailsjs/go/app/App";
import {ElMessage, ElMessageBox, ElTable} from "element-plus";
import {EventsOn} from "../../wailsjs/runtime";
import {useStore} from "vuex";
import {app} from "../../wailsjs/go/models";
import DeleteKey = app.DeleteKey;


export default {
  setup(props, context) {
    const TypeFolder = 'Folder'
    const TypeObject = 'Object'
    const route = useRoute()
    const router = useRouter()
    const store = useStore()
    var dialogFormVisible = ref(false)
    var dialogCreateDirVisible = ref(false)
    var dialogDeleteVisible = ref(false)
    const loading = ref(true)
    const bucketName = route.params.bucketName
    let prefix = ref('')
    const subPrefix = ref('')
    const basePath = 's3://' + bucketName + '/'
    let showPath = computed(() => {
      return basePath + prefix.value
    })
    const percentage = ref(0)
    const percentageLabel = ref('')
    const newFolderName = ref('')
    const disableMoreButton = ref(true)

    var marker = ''
    var nextMarker = ''
    const isTruncate = ref(false)
    const markerStack = ref([])

    const deleteData = reactive({
          count: 0,
          success: 0
        }
    )
    const deletePercentage = computed(() => {
      if (deleteData.count === 0) {
        return 0
      } else {
        return (deleteData.success * 100 / deleteData.count)
      }
    })
    const deletePercentageLabel = computed(() => {
      return deleteData.success + " / " + deleteData.count
    })

    const colors = [
      {color: '#1989fa', percentage: 100},
      {color: '#5cb87a', percentage: 101},
    ]

    const handleUploadCommand = (command: string | number | object) => {
      if (command === "UploadFiles") {
        selectUploadObjects()
      } else if (command === "UploadFolder") {
        selectUploadFolder()
      }
    }

    const handleMoreCommand = (command: string | number | object) => {
      if (command === "RemoveFiles") {
        confirmDeleteObjects()
      } else if (command === "DownloadFiles") {
        selectDownloadDir()
      }
    }

    interface SelectedObject {
      type: string
      key: string
      lastModified: string
      size: string
      humanSize: string
    }

    const multipleTableRef = ref<InstanceType<typeof ElTable>>()
    const multipleSelection = ref<SelectedObject[]>([])

    const handleSelectionChange = (val: SelectedObject[]) => {
      disableMoreButton.value = val.length === 0;
      multipleSelection.value = val
      console.log(multipleSelection.value)
    }

    const backward = () => {
      if (prefix.value === '') {
        router.push('/main')
      } else {
        var p = prefix.value
        var i = p.slice(0, p.length - 1).lastIndexOf('/')
        prefix.value = p.slice(0, i + 1)
        marker = ''
        listObjects(bucketName, marker, prefix.value, 100)
      }
    }

    const objectList = computed(() => {
      return data.tableData
    })

    let data = reactive({
      tableData: [],
    })

    onMounted(() => {
      listObjects(bucketName, marker, prefix.value, 100)
    })

    const listObjects = async (bucketName, marker, prefix, maxKey) => {
      var tableData = []
      var p = prefix
      var i = p.lastIndexOf('/')
      var folderPrefix = p.slice(0, i + 1)
      ListObjects(bucketName, marker, prefix, maxKey).then(res => {
        loading.value = true
        if (res.err !== '') {
          ElMessage.error(res)
        } else {
          if (res.prefixes !== null) {
            res.prefixes.forEach((p, i) => {
              var k = p
              var showed = k
              k = k.slice(folderPrefix.length, p.length)
              if (k.length > 40) {
                showed = k.slice(0, 40) + "..."
              }
              if (k === '') {
                k = '/'
              }
              const folderInfo = {
                type: TypeFolder,
                key: k,
                showed: showed,
              }
              tableData.push(folderInfo)
            })
          }

          if (res.contents !== null) {
            res.contents.forEach((c, i) => {
              var k = c.key
              var showed = k
              k = k.slice(folderPrefix.length, c.key.length)
              if (k.length > 40) {
                showed = k.slice(0, 30) + "..."
              }
              if (k !== '') {
                const objectInfo = {
                  type: TypeObject,
                  key: k,
                  showed: showed,
                  lastModified: c.lastModified,
                  size: c.size,
                  humanSize: c.humanSize
                }
                tableData.push(objectInfo)
              }
            })
          }
          nextMarker = res.nextMarker
          isTruncate.value = res.isTruncated
          data.tableData = tableData
        }
        loading.value = false
      })
    }

    const toPreviousPage = () => {
      marker = markerStack.value.pop()
      listObjects(bucketName, marker, prefix.value, 100)
    }

    const toNextPage = () => {
      markerStack.value.push(marker)
      marker = nextMarker
      listObjects(bucketName, marker, prefix.value, 100)
    }

    const createDir = () => {
      PutDir(<string>bucketName, newFolderName.value, prefix.value).then(res => {
        if (res.err !== '') {
          ElMessage.error(res.err)
        } else {
          dialogCreateDirVisible.value = false
          listObjects(bucketName, '', prefix.value, 100)
          newFolderName.value = ''
        }
      })
    }

    const selectUploadFolder = () => {
      SelectUploadFolder().then((res) => {
        if (res.err !== '') {
          ElMessage.error(res.err)
        } else {
          if (res.path === '') {
            return
          }
          var eventWalkPath = "w" + res.path + Math.random()
          // open TaskList
          context.emit('changeVisible', true)
          EventsOn(eventWalkPath, (fp) => {
            if (fp.type === TypeFolder) {
              PutDir(<string>bucketName, fp.key, prefix.value).then(res => {
                if (res.err !== '') {
                  ElMessage.error(res.err)
                } else {
                  dialogCreateDirVisible.value = false
                  listObjects(bucketName, '', prefix.value, 100)
                  newFolderName.value = ''
                }
              })
            } else {
              var eventProgress = "u" + fp.key + Math.random()
              const file = {
                accountId: fp.accountId,
                type: fp.type,
                bucket: bucketName,
                name: fp.name,
                key: fp.key,
                source: fp.source,
                size: fp.size,
                humanSize: fp.humanSize,
                status: 0,  // 0-PENDING, 1-PAUSE, 2-ERROR, 3-FINISH
                taskId: eventProgress
              }
              const payload = {
                file: file,
                progress: eventProgress
              }
              store.commit('addToUploadList', payload)
              // begin to upload
              EventsOn(eventProgress, (data) => {
                const payload = {
                  data: data,
                  progress: eventProgress
                }
                store.commit('updateUploadProgress', payload)
              })
              DoPutObject(<string>bucketName, fp.key, fp.source, eventProgress).then(res => {
                if (res.err !== '') {
                  ElMessage.error(res.err)
                } else {
                  const payload = {
                    data: 100,
                    progress: eventProgress
                  }
                  store.commit('updateUploadProgress', payload)
                  listObjects(bucketName, '', prefix.value, 100)
                }
              })
            }
          })
          DoUploadFolder(prefix.value, res.path, eventWalkPath)
        }
      })
    }

    const selectUploadObjects = () => {
      SelectUploadFiles(prefix.value).then(res => {
        if (res.err !== '') {
          ElMessage.error(res.err)
        } else {
          res.files.forEach((fp, i) => {
            var eventProgress = "u" + fp.key + Math.random()
            const file = {
              accountId: fp.accountId,
              bucket: bucketName,
              name: fp.name,
              key: fp.key,
              source: fp.source,
              size: fp.size,
              humanSize: fp.humanSize,
              status: 0,  // 0-PENDING, 1-PAUSE, 2-ERROR, 3-FINISH
              taskId: eventProgress
            }
            const payload = {
              file: file,
              progress: eventProgress
            }
            store.commit('addToUploadList', payload)
            // begin to upload
            EventsOn(eventProgress, (data) => {
              const payload = {
                data: data,
                progress: eventProgress
              }
              store.commit('updateUploadProgress', payload)
            })

            // open TaskList
            context.emit('changeVisible', true)
            DoPutObject(<string>bucketName, fp.key, fp.source, eventProgress).then(res => {
              if (res.err !== '') {
                ElMessage.error(res.err)
              } else {
                const payload = {
                  data: 100,
                  progress: eventProgress
                }
                store.commit('updateUploadProgress', payload)
                listObjects(bucketName, '', prefix.value, 100)
              }
            })
          })
        }
      })
    }

    const selectDownloadDir = () => {
      SelectDownloadPath().then((res) => {
        if (res.err !== '') {
          ElMessage.error(res.err)
        } else {
          multipleSelection.value.forEach((v, i) => {
            var eventProgress = "d" + prefix.value + v.key + Math.random()
            const downloadTask = {
              bucket: bucketName,
              key: v.key,
              size: v.size,
              humanSize: v.humanSize,
              status: 0,  // 0-PENDING, 1-PAUSE, 2-ERROR, 3-FINISH
              taskId: eventProgress
            }
            const payload = {
              file: downloadTask,
              progress: eventProgress
            }
            store.commit('addToDownloadList', payload)
            // begin to download
            EventsOn(eventProgress, (data) => {
              const payload = {
                data: data,
                progress: eventProgress
              }
              store.commit('updateDownloadProgress', payload)
            })

            // open TaskList and download label
            context.emit('changeVisible', true)
            context.emit('setTaskTabName', "download")
          })
        }
      })
    }

    const getObject = (key) => {
      var eventDialog = "downloadDialog"
      var eventProgress = "d" + prefix.value + key + Math.random()
      percentageLabel.value = 'Downloading...'
      EventsOn(eventDialog, () => {
        dialogFormVisible.value = true
      })
      // FIXME: fix range download
      percentage.value = 0
      EventsOn(eventProgress, (data) => {
        percentage.value = data
      })

      GetObject(<string>bucketName, prefix.value + key, true, eventDialog, eventProgress).then(res => {
        if (res.err !== '') {
          ElMessage.error(res.err)
          percentageLabel.value = 'Failed'
        } else {
          percentage.value = 100
          percentageLabel.value = 'Finished'
        }
      })
    }

    // delete single object
    const deleteObject = (key, keyType) => {
      var eventDeleteCount = "r" + Math.random()
      var eventDeleteSuccess = eventDeleteCount + "_success"
      deleteData.count = 0
      deleteData.success = 0
      EventsOn(eventDeleteCount, (data) => {
        deleteData.count = data
      })
      EventsOn(eventDeleteSuccess, (data) => {
        deleteData.success = data
      })
      dialogDeleteVisible.value = true
      DeleteObject(<string>bucketName, prefix.value + key, keyType, eventDeleteSuccess, eventDeleteCount)
          .then(res => {
            if (res.err !== '') {
              ElMessage.error("delete object failed")
            } else {
              listObjects(bucketName, '', prefix.value, 100)
            }
          })
    }

    const confirmDeleteObjects = () => {
      ElMessageBox.confirm(
          'xBrowser will permanently delete these files. Continue?',
          'Warning',
          {
            confirmButtonText: 'Delete',
            cancelButtonText: 'Cancel',
            type: 'warning',
          }
      )
          .then(() => {
            deleteObjects()
          })
          .catch(() => {
            ElMessage({
              type: 'info',
              message: 'Delete canceled',
            })
          })
    }

    // delete multiple objects
    const deleteObjects = () => {
      var deleteKeys = []
      var eventDeleteCount = "r" + Math.random()
      var eventDeleteSuccess = eventDeleteCount + "_success"
      deleteData.count = 0
      deleteData.success = 0
      EventsOn(eventDeleteCount, (data) => {
        deleteData.count = data
      })
      EventsOn(eventDeleteSuccess, (data) => {
        deleteData.success = data
      })
      dialogDeleteVisible.value = true
      multipleSelection.value.forEach((v, i) => {
        var k = new DeleteKey()
        k.key = prefix.value + v.key
        k.keyType = v.type
        deleteKeys.push(k)
      })
      DeleteObjects(<string>bucketName, deleteKeys, eventDeleteSuccess, eventDeleteCount)
          .then(res => {
            if (res.err !== '') {
              ElMessage.error("delete objects failed")
            } else {
              listObjects(bucketName, '', prefix.value, 100)
            }
          })
    }

    // into `Folder`, otherwise preview the object
    const clickKey = (data) => {
      if (data.type === TypeFolder) {
        prefix.value += data.key
        marker = ''
        markerStack.value = []
        listObjects(bucketName, marker, prefix.value, 100)
      } else {
        // TODO: show preview dialog
        console.log(prefix.value, data.key)
      }
    }

    const searchBySubPrefix = () => {
      marker = ''
      markerStack.value = []
      listObjects(bucketName, marker, prefix.value + subPrefix.value, 100)
    }

    return {
      ...toRefs(data),
      listObjects,
      toPreviousPage,
      toNextPage,
      searchBySubPrefix,
      getObject,
      deleteObject,
      confirmDeleteObjects,
      deleteObjects,
      createDir,
      clickKey,
      backward,
      selectUploadObjects,
      selectDownloadDir,
      handleSelectionChange,
      handleUploadCommand,
      handleMoreCommand,
      disableMoreButton,
      newFolderName,
      multipleTableRef,
      dialogFormVisible,
      dialogCreateDirVisible,
      dialogDeleteVisible,
      loading,
      objectList,
      prefix,
      subPrefix,
      showPath,
      bucketName,
      percentage,
      percentageLabel,
      deletePercentage,
      deletePercentageLabel,
      colors,
      markerStack,
      isTruncate,
    }
  }
}
</script>

<style scoped>
.container {
  width: 90%;
  text-align: left;
  margin: 6px auto 1% auto;
  border: 1px solid #EEE;
}

.container2 {
  width: 90%;
  display: flex;
  margin: 6px auto 1% auto;
}

.percentage-value {
  display: block;
  margin-top: 10px;
  font-size: 28px;
}

.percentage-label {
  display: block;
  margin-top: 10px;
  font-size: 12px;
}

.input-with-select {
  width: 50%;
  border-radius: 0;
}

.more-button .el-dropdown-link {
  cursor: pointer;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
}

</style>