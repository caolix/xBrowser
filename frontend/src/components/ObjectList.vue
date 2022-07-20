<template>
  <el-row>
    <el-col :span=4>
      <el-button type="primary" @click="selectObjects">
        <el-icon>
          <UploadFilled/>
        </el-icon>
        Upload
      </el-button>
    </el-col>
    <el-col :span=4>
      <el-button plain type="danger" @click="backward">
        <el-icon>
          <Top/>
        </el-icon>
        Backward
      </el-button>
    </el-col>
    <el-col :span=4>
      <el-input v-model="showPath"></el-input>
    </el-col>
  </el-row>

  <!--Dialog-->
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

  <div class="container">
    <el-table
        :data="objectList"
        max-height="800"
        class="table"
        v-loading="loading"
    >
      <el-table-column prop="key" label="Name" width="400">
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
      <el-table-column label="Operations">
        <template #default="scope">
          <el-button size="small" @click="getObject(tableData[scope.$index].key)">
            Download
          </el-button>
          <el-popconfirm
              confirm-button-text="Yes"
              cancel-button-text="No"
              :icon="QuestionFilled"
              icon-color="#FF0000"
              title="Are you sure to delete this object?"
              @confirm="deleteObject(tableData[scope.$index].key)"
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
</template>

<script>

import {useRoute, useRouter} from "vue-router";
import {computed, onMounted, reactive, ref, toRefs} from "vue";
import {DeleteObject, GetObject, ListObjects, DoPutObject, SelectFiles} from "../../wailsjs/go/app/App";
import {ElMessage} from "element-plus";
import {EventsOn} from "../../wailsjs/runtime";
import {useStore} from "vuex";


export default {
  setup() {
    const route = useRoute()
    const router = useRouter()
    const store = useStore()
    var dialogFormVisible = ref(false)
    const loading = ref(true)
    const bucketName = route.params.bucketName
    let prefix = ref('')
    const basePath = 's3://' + bucketName + '/'
    let showPath = computed(() => {
      return basePath + prefix.value
    })
    const percentage = ref(0)
    const percentageLabel = ref('')

    const colors = [
      {color: '#1989fa', percentage: 100},
      {color: '#5cb87a', percentage: 101},
    ]

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

    const listObjects = (bucketName, marker, prefix, maxKey) => {
      data.tableData = []
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
                type: 'Folder',
                key: k,
              }
              data.tableData.push(folderInfo)
            })
          }

          if (res.contents !== null) {
            res.contents.forEach((c, i) => {
              var k = c.key
              k = k.slice(prefix.length, c.key.length)
              if (k !== '') {
                const objectInfo = {
                  type: 'Object',
                  key: k,
                  last_modified: c.last_modified,
                  size: c.human_size,
                }
                data.tableData.push(objectInfo)
              }
            })
          }
        }
        loading.value = false
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
              name: fp.name,
              key: fp.key,
              source: fp.source,
              size: fp.size,
              human_size: fp.human_size,
              progress: eventProgress
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

            DoPutObject(bucketName, fp.key, fp.source, eventProgress).then(res => {
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

    const putObject = () => {
      var eventDialog = "uploadDialog"
      var eventProgress = "u" + prefix.value + Math.random()
      percentageLabel.value = 'Uploading...'
      EventsOn(eventDialog, () => {
        dialogFormVisible.value = true
      })
      percentage.value = 0
      EventsOn(eventProgress, (data) => {
        percentage.value = data
      })
      PutObject(bucketName, prefix.value, eventDialog, eventProgress).then(res => {
        if (res.err !== '') {
          ElMessage.error(res.err)
          percentageLabel.value = 'Failed'
        } else {
          percentage.value = 100
          listObjects(bucketName, '', prefix.value, 100)
          percentageLabel.value = 'Finished'
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

      GetObject(bucketName, prefix.value + key, true, eventDialog, eventProgress).then(res => {
        if (res.err !== '') {
          ElMessage.error(res.err)
          percentageLabel.value = 'Failed'
        } else {
          percentage.value = 100
          percentageLabel.value = 'Finished'
        }
      })
    }

    const deleteObject = (key) => {
      console.log("deleteObject: " + prefix.value + key)
      DeleteObject(bucketName, prefix.value + key).then(res => {
        if (res.err !== '') {
          ElMessage.error("delete object failed: ", res.err)
        } else {
          listObjects(bucketName, '', prefix.value, 100)
        }
      })
    }

    const toPreview = (data) => {
      if (data.type === 'Folder') {
        prefix.value += data.key
        listObjects(bucketName, '', prefix.value, 100)
      } else {
        console.log(prefix.value, data.key)
      }
    }

    return {
      ...toRefs(data),
      listObjects,
      putObject,
      getObject,
      deleteObject,
      toPreview,
      backward,
      selectObjects,
      dialogFormVisible,
      loading,
      objectList,
      prefix,
      showPath,
      bucketName,
      percentage,
      percentageLabel,
      colors
    }
  }
}
</script>

<style scoped>
.container {
  width: 90%;
  text-align: left;
  margin: 50px auto 20px auto;
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