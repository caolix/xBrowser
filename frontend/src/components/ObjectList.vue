<template>
  <el-row>
    <el-col :span=4>
      <el-button type="primary" @click="putObject">
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

  <div class="container">
    <el-table
        :data="objectList"
        max-height="800"
        class="table"
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
          <el-button
              size="small"
              type="danger"
              @click="deleteObject(tableData[scope.$index].key)">
            Remove
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>

import {useRoute, useRouter} from "vue-router";
import {computed, onMounted, reactive, ref, toRefs} from "vue"
import {ListObjects} from "../../wailsjs/go/app/App";
import {ElMessage} from "element-plus";

export default {
  setup() {
    const route = useRoute()
    const router = useRouter()
    const bucketName = route.params.bucketName
    let prefix = ref('')
    const basePath = 's3://' + bucketName + '/'
    let showPath = computed(() => {
      return basePath + prefix.value
    })

    const backward = () => {
      if (prefix.value === '') {
        router.push('/main')
      } else {
        var p = prefix.value
        var i = p.slice(0, p.length - 1).lastIndexOf('/')
        prefix.value = p.slice(0, i+1)
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
      })
    }

    const putObject = () => {

    }

    const getObject = () => {

    }

    const deleteObject = (str) => {


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
      objectList,
      prefix,
      showPath,
      bucketName
    }
  }
}
</script>

<style scoped>
.container {
  width: 1000px;
  text-align: left;
  margin: 50px auto 20px auto;
  border: 1px solid #EEE;
  min-height: 600px;
}
</style>