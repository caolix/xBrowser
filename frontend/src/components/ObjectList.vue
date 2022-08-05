<template>
  <div class="container">
    <el-space :size="20" alignment="start">
      <el-button type="primary" @click="selectObjects">
        <el-icon style="padding-right: 6px">
          <UploadFilled/>
        </el-icon>
        Upload
      </el-button>
      <el-button type="primary" @click="dialogCreateDirVisible = true">
        <el-icon style="padding-right: 6px">
          <FolderAdd/>
        </el-icon>
        Create Folder
      </el-button>
      <el-button plain type="danger" @click="confirmDeleteObjects" :disabled="disableDeleteButton">
        <el-icon style="padding-right: 6px">
          <Delete/>
        </el-icon>
        Delete
      </el-button>
      <el-space :size="1">
        <el-button type="primary" @click="backward">
          <el-icon>
            <Top/>
          </el-icon>
        </el-button>
        <el-input v-model="showPath"></el-input>
      </el-space>
    </el-space>
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
      <el-table-column prop="key" label="Name" width="400" show-overflow-tooltip="true">
        <template #default="scope">
          <el-icon style="padding-right: 2px; padding-top: 2px" v-if="tableData[scope.$index].type==='Folder'">
            <FolderOpened/>
          </el-icon>
          <el-icon style="padding-right: 2px; padding-top: 2px" v-else>
            <Document/>
          </el-icon>
          <el-button
              link
              type="primary"
              size="small"
              @click="toPreview(tableData[scope.$index])"
          >
            {{ tableData[scope.$index].key }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column prop="type" label="Type" width="200"></el-table-column>
      <el-table-column prop="size" label="Size" width="200"></el-table-column>
      <el-table-column label="Operations" fixed="right" align="right">
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
  <el-button-group>
    <el-button type="primary" :icon="ArrowLeft">Previous Page</el-button>
    <el-button type="primary">
      Next Page<el-icon class="el-icon--right"><ArrowRight /></el-icon>
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
  GetObject,
  ListObjects,
  PutDir,
  SelectFiles
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
    const basePath = 's3://' + bucketName + '/'
    let showPath = computed(() => {
      return basePath + prefix.value
    })
    const percentage = ref(0)
    const percentageLabel = ref('')
    const newFolderName = ref('')
    const disableDeleteButton = ref(true)

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

    interface SelectedObject {
      type: string
      key: string
      last_modified: string
      size: string
    }

    const multipleTableRef = ref<InstanceType<typeof ElTable>>()
    const multipleSelection = ref<SelectedObject[]>([])

    const handleSelectionChange = (val: SelectedObject[]) => {
      disableDeleteButton.value = val.length === 0;
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
        listObjects(bucketName, '', prefix.value, 100)
      }
    }

    const objectList = computed(() => {
      return data.tableData
    })

    let data = reactive({
      tableData: [],
    })

    onMounted(() => {
      listObjects(bucketName, '', prefix.value, 100)
    })

    const listObjects = async (bucketName, marker, prefix, maxKey) => {
      var tableData = []
      ListObjects(bucketName, marker, prefix, maxKey).then(res => {
        loading.value = true
        if (res.err !== '') {
          ElMessage.error(res)
        } else {
          if (res.prefixes !== null) {
            res.prefixes.forEach((p, i) => {
              var k = p
              k = k.slice(prefix.length, p.length)
              if (k === '') {
                k = '/'
              }
              const folderInfo = {
                type: TypeFolder,
                key: k,
              }
              tableData.push(folderInfo)
            })
          }

          if (res.contents !== null) {
            res.contents.forEach((c, i) => {
              var k = c.key
              k = k.slice(prefix.length, c.key.length)
              if (k !== '') {
                const objectInfo = {
                  type: TypeObject,
                  key: k,
                  last_modified: c.last_modified,
                  size: c.human_size,
                }
                tableData.push(objectInfo)
              }
            })
          }
          data.tableData = tableData
        }
        loading.value = false
      })
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

    const selectObjects = () => {
      SelectFiles(prefix.value).then(res => {
        if (res.err !== '') {
          ElMessage.error(res.err)
        } else {
          res.files.forEach((fp, i) => {
            var eventProgress = "u" + fp.key + Math.random()
            const file = {
              bucket: bucketName,
              name: fp.name,
              key: fp.key,
              source: fp.source,
              size: fp.size,
              human_size: fp.human_size,
              isPending: false,
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
              store.commit('updateProgress', payload)
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
                store.commit('updateProgress', payload)
                listObjects(bucketName, '', prefix.value, 100)
              }
            })
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
    const toPreview = (data) => {
      if (data.type === TypeFolder) {
        prefix.value += data.key
        listObjects(bucketName, '', prefix.value, 100)
      } else {
        console.log(prefix.value, data.key)
      }
    }

    return {
      ...toRefs(data),
      listObjects,
      getObject,
      deleteObject,
      confirmDeleteObjects,
      deleteObjects,
      createDir,
      toPreview,
      backward,
      selectObjects,
      handleSelectionChange,
      disableDeleteButton,
      newFolderName,
      multipleTableRef,
      dialogFormVisible,
      dialogCreateDirVisible,
      dialogDeleteVisible,
      loading,
      objectList,
      prefix,
      showPath,
      bucketName,
      percentage,
      percentageLabel,
      deletePercentage,
      deletePercentageLabel,
      colors
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

</style>