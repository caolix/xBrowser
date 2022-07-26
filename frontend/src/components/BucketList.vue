<template>

  <el-row>
    <el-col :span=4>
      <el-button type="primary" @click="dialogFormVisible = true">
        <el-icon style="padding-right: 6px"><Plus/></el-icon>
        Create Bucket
      </el-button>
    </el-col>
  </el-row>

  <!--  Create Bucket Dialog-->
  <el-dialog v-model="dialogFormVisible" title="Create Bucket" destroy-on-close>
    <span>Name:</span>
    <el-input v-model="form.bucketname" autocomplete="off"/>
    <template #footer>
      <span class="dialog-footer">
         <el-button type="primary" @click="toCreateBucket"
         >Create</el-button
         >
        <el-button @click="dialogFormVisible = false">Cancel</el-button>
      </span>
    </template>
  </el-dialog>

  <!--  BucketList-->
  <div class="container">
    <el-table
        :data="bucketList"
        max-height="800"
        class="table"
        v-loading="loading"
    >
      <!--      <el-table-column type="index" label="#" width="100"/>-->
      <el-table-column prop="name" label="BucketName" width="200">
        <template #default="scope">
          <el-icon style="padding-right: 2px; padding-top: 2px">
            <Menu/>
          </el-icon>
            <el-button
                link
                size="small"
                @click="toBucket(tableData[scope.$index].name)"
            >
              {{ tableData[scope.$index].name }}
            </el-button>
        </template>
      </el-table-column>
      <el-table-column label="Operations">
        <template #default="scope">
          <el-popconfirm
              confirm-button-text="Yes"
              cancel-button-text="No"
              :icon="QuestionFilled"
              icon-color="#FF0000"
              title="Are you sure to delete this bucket?"
              @confirm="toDeleteBucket(tableData[scope.$index].name)"
          >
            <template #reference>
              <el-button
                  size="small"
                  type="danger"
              >
                <el-icon style="padding-right: 6px">
                  <Delete/>
                </el-icon>
                Delete
              </el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import {computed, onMounted, reactive, ref, toRefs} from "vue";
import {DeleteBucket, ListBuckets, MakeBucket} from "../../wailsjs/go/app/App";
import {ElMessage} from "element-plus";
import {useRouter} from "vue-router"

export default {
  setup() {
    const loading = ref(true)
    let form = reactive({
      bucketname: '',
    })

    var dialogFormVisible = ref(false)

    const router = useRouter()
    const bucketList = computed(() => {
      return data.tableData
    })

    let data = reactive({
      tableData: [],
    })

    onMounted(() => {
      listBuckets()
    })

    const listBuckets = () => {
      ListBuckets().then(res => {
        loading.value = true
        if (res.err !== '') {
          ElMessage.error(res)
        } else {
          res.buckets.forEach((b, i) => {
            const bucketInfo = {
              name: b,
            }
            data.tableData.push(bucketInfo)
          })
        }
        loading.value = false
      })
    }

    const toDeleteBucket = (bucketName) => {
      DeleteBucket(bucketName).then(res => {
        if (res !== '') {
          ElMessage.error(res)
        } else {
          ElMessage.success("delete bucket success")
          data.tableData = []
          listBuckets()
        }
      })
    }

    const toCreateBucket = () => {
      dialogFormVisible.value = false
      MakeBucket(form.bucketname).then(res => {
            if (res !== '') {
              ElMessage.error(res)
            } else {
              ElMessage.success("create bucket success")
              form.bucketname = ''
              data.tableData = []
              listBuckets()
            }
          }
      )
    }

    const toBucket = (b) => {
      router.push({
        path: '/main/bucket/' + b,
      })
    }

    return {
      dialogFormVisible,
      ...toRefs(data),
      form,
      loading,
      toCreateBucket,
      toDeleteBucket,
      toBucket,
      listBuckets,
      bucketList,
    }

  }
}


</script>

<style scoped>
.container {
  width: 90%;
  text-align: left;
  margin: 10px auto 2% auto;
  border: 1px solid #EEE;
}

.el-link {
  margin-right: 8px;
}

.el-link .el-icon--right.el-icon {
  vertical-align: text-bottom;
}

.el-table .el-table__cell {
  padding: 0;
}
</style>