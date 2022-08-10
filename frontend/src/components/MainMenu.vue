<template>
  <el-menu
      :default-active="activeIndex"
      class="el-menu-demo"
      mode="horizontal"
      background-color="#545c64"
      text-color="#fff"
      active-text-color="#ffd04b"
      :ellipsis="false"
      @select="handleSelect"
  >
    <el-menu-item index="0">
<!--            <img src="../assets/images/uni2.png" style="width:140px; height:50px; margin-top:6px;"/>-->
      <img src="../assets/images/bear.gif" style="width:50px; height:50px; margin-top:6px;"/>
      <img src="../assets/images/logo4.png" style="width:150px; height:50px; margin-top:6px;"/>
    </el-menu-item>
    <div class="flex-grow"/>
    <el-menu-item index="1" @click="showDrawer">
      <span>TaskList</span>
    </el-menu-item>
    <el-menu-item index="2">Settings</el-menu-item>
    <el-menu-item index="3">About</el-menu-item>
    <el-menu-item index="4" style="height:56px">
      <el-button link type="danger" @click="logout">
        <el-icon style="padding-right: 1px">
          <SwitchButton/>
        </el-icon>
        Logout
      </el-button>
    </el-menu-item>
  </el-menu>
</template>

<script>
import {useRouter} from "vue-router";

export default {
  name: "MainMenu",

  setup(props, context) {
    const router = useRouter()
    const logout = () => {
      router.push({
        path: '/'
      })
    }
    const showDrawer = () => {
      context.emit('drawerVisible', true)
    }
    return {
      showDrawer,
      logout
    }
  }
}
</script>

<style scoped>
.flex-grow {
  flex-grow: 1;
}
</style>