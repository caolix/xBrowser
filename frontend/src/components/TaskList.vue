<template>
  <el-drawer v-model="props.visible" :show-close="false" size="40%" @close="closeDrawer">
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
      </div>
    </template>
  </el-drawer>
</template>

<script>
import {ref, computed} from 'vue'
import {useStore} from "vuex";
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
      context.emit('changeVisible', false)
    }
    const showDrawer = () => {
      context.emit('changeVisible', true)
    }
    const closeDrawer = () => {
      context.emit('changeVisible', false)
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
      props,
      colors,
      activeName,
      uploadListData,
      percentageMap,
      cancelClick,
      showDrawer,
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