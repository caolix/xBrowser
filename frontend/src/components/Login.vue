<template>
  <main>
<!--        <div>-->
<!--          <img id="logo-uni" alt="Wails logo" src="../assets/images/uni2.png"/>-->
<!--        </div>-->
    <div style="display: flex;">
      <img id="logo" alt="Wails logo" src="../assets/images/bear_good.gif"/>
      <img id="logo2" alt="Wails logo" src="../assets/images/logo1.jpeg"/>
    </div>
    <div class="container" :style="{
          boxShadow: `var(--el-box-shadow-light)`,
        }">
      <div class="input">
        <span style="color: #ff0000"> * </span><span class="demo-input-label">Endpoint:</span>
        <el-row :gutter="10">
          <el-input
              v-model="data.endpoint"
              placeholder="ENDPOINT"
              class="endpoint-input"
          />
          <el-checkbox v-model="data.usessl" label="HTTPS" size="large"/>
        </el-row>
        <span style="color: #ff0000"> * </span><span class="demo-input-label">AccessKey:</span>
        <el-row :gutter="10">
          <el-input
              v-model="data.ak"
              placeholder="ACCESS_KEY"
          />
        </el-row>
        <span style="color: #ff0000"> * </span><span class="demo-input-label">SecretKey:</span>
        <el-row :gutter="10">
          <el-input
              v-model="data.sk"
              placeholder="SECRET_KEY"
          />
        </el-row>
        <span class="demo-input-label">Preset Path:</span>
        <el-row :gutter="10">
          <el-input
              v-model="data.prepath"
              placeholder="Optional. format: s3://bucket/dir/"
          />
        </el-row>
        <span class="demo-input-label">Description:</span>
        <el-row :gutter="10">
          <el-input
              v-model="data.remark"
              placeholder="Optional. Up to 30 words."
          />
        </el-row>
        <el-row :gutter="10">
          <el-col :span="16">
            <div class="demo-input-label">
              <el-checkbox v-model="data.save" label="Remember" size="large"/>
            </div>
          </el-col>
        </el-row>
        <el-button type="primary" class="submit" @click="login">Login</el-button>
      </div>
    </div>
  </main>
</template>


<script>

import {reactive} from 'vue'
import {CheckDbError, LoadLatestLoginInfo, Login} from '../../wailsjs/go/app/App'
import {useRouter} from "vue-router";
import {db} from '../../wailsjs/go/models'
import {ElLoading, ElMessage} from 'element-plus'

export default {
  setup() {
    // NOTE: this code MUST before than ElLoading
    const router = useRouter()
    const loading = ElLoading.service({
      lock: true,
      text: 'Loading...',
      background: 'rgba(0, 0, 0, 0.7)',
    })

    const data = reactive({
      ak: "",
      sk: "",
      endpoint: "",
      prepath: "",
      remark: "",
      save: false,
      usessl: false,
    })

    CheckDbError().then(err => {
      if (err !== "") {
        ElMessage({
          showClose: true,
          message: err,
          type: 'warning',
          duration: 5000,
        })
      }
    })
    LoadLatestLoginInfo().then(info => {
          data.ak = info.ak
          data.sk = info.sk
          data.endpoint = info.endpoint
          data.prepath = info.prepath
          data.remark = info.remark
          data.usessl = info.useSSL
          loading.close()
        }
    )

    const login = () => {
      let info = new db.LoginInfo()
      info.endpoint = data.endpoint
      info.ak = data.ak
      info.sk = data.sk
      info.prepath = data.prepath
      info.remark = data.remark
      Login(info, data.save, data.usessl).then((res) => {
        if (res === "") {
          router.push({path: '/main'})
        } else {
          ElMessage.error(res)
        }
      })
    }

    return {
      loading,
      data,
      login,
    }
  }
}


</script>

<style scoped>

.el-checkbox {
  color: #000000;
}

.submit {
  margin: 0 auto;
  display: flex;
  width: 40%;
  height: 10%;
}

.demo-input-label {
  display: inline-block;
  width: 75%;
  margin-top: 2%;
  color: #000000;
  height: 2%;
  font-size: 85%;
}

.container {
  width: 60%;
  text-align: left;
  margin: 1% auto 2% auto;
  border: 1px solid #EEE;
}

.endpoint-input {
  width: 80%;
  margin-right: 2%;
}

.input {
  margin: 0 2%;
  padding-bottom: 5%;
}

.input-box .btn {
  width: 60%;
  height: 3%;
  line-height: 3%;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 3px;
  line-height: 3px;
  padding: 0 1px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}
</style>
