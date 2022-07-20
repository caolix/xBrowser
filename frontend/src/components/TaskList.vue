<template>
  <el-drawer v-model="taskDrawer" :show-close="false" size="40%">
    <template #default>
      <el-tabs v-model="activeName" class="demo-tabs" @tab-click="handleClick">
        <el-tab-pane label="Upload" name="upload">
          <el-table :data="uploadListData" style="width: 100%">
            <el-table-column prop="key"  width="150" />
            <el-table-column prop="human_size"  width="100" />
            <el-table-column>
              <template #default="scope">
                <el-progress :percentage="percentageMap[scope.row.progress]" :color="colors" />
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
        <el-button type="primary" @click="confirmClick">confirm</el-button>
      </div>
    </template>
  </el-drawer>

  <el-button type="primary" class="task" @click="taskDrawer = true">Task</el-button>
</template>

<script>
import {ref, computed, reactive} from 'vue'
import {useStore} from "vuex";
export default {
  name: "TaskList",
  setup() {
    const activeName = ref('upload')
    const taskDrawer = ref(false)
    const store = useStore()
    const cancelClick = () => {
      taskDrawer.value = false
    }
    const confirmClick = () => {
      console.log(store.state.uploadProgress)
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

    return{
      colors,
      activeName,
      taskDrawer,
      uploadListData,
      percentageMap,
      cancelClick,
      confirmClick
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