<template>
  <el-header>
    <el-row>
      <el-col :span="1" style="margin-left: 20px">
        <el-image
            style="width: 60px; height: 60px"
            :src="src"
        ></el-image>
      </el-col>
      <el-col :span="6" :offset="16">
        <el-link style="padding: 10px" href="/document">
          <el-button plain type="success">文档</el-button>
        </el-link>
      </el-col>
    </el-row>
  </el-header>
</template>

<script>
export default {
  setup() {
    const src = ref('../assets/images/bear.png')
    return {
      src
    }
  }
}


</script>

<style scoped>

</style>