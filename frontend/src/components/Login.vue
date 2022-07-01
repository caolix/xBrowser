<template>
  <main>
    <img id="logo" alt="Wails logo" src="../assets/images/comeon.gif"/>
    <div class="container">
      <div class="input">
        <div id="input" class="input-box">
          <div class="demo-input-suffix">
            <span style="color: #ff0000"> * </span><span class="demo-input-label">Endpoint:</span>
            <el-row :gutter="10">
              <el-input
                  v-model="data.endpoint"
                  placeholder="ENDPOINT"
              />
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
      </div>
    </div>
  </main>
</template>


<script>

import {reactive} from 'vue'
import {CheckDbError, LoadLatestLoginInfo, Login} from '../../wailsjs/go/app/App'
import {useRouter} from "vue-router";
import {ElMessage} from 'element-plus'

export default {
  setup: function () {
    const data = reactive({
      ak: "",
      sk: "",
      endpoint: "",
      prepath: "",
      remark: "",
      save: false,
    })
    const router = useRouter()
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
          data.ak = info.accessKey
          data.sk = info.secretKey
          data.endpoint = info.endpoint
          data.prepath = info.prepath
          data.remark = info.remark
        }
    )

    function login() {
      let info = new db.LoginInfo()
      info.endpoint = data.endpoint
      info.accessKey = data.accessKey
      info.secretKey = data.secretKey
      info.prepath = data.prepath
      info.remark = data.remark
      Login(info, data.save).then((res) => {
        if (res === "") {
          router.push("/main")
          return
        }
        ElMessage.error(res)
      })
    }

    return {
      data,
      router,
      login,
    }
  }
}


</script>

<style scoped>

.el-checkbox {
  color: #ffffff;
}

.submit {
  margin-top: 40px;
  margin-left: -10px;
  margin-right: -10px;
  width: 340px;
}

.demo-input-label {
  display: inline-block;
  width: 130px;
  margin-top: 20px;
}

.container {
  width: 420px;
  text-align: left;
  margin: 50px auto 20px auto;
  min-height: 600px;
}

.input {
  margin: 0 50px;
}

.input-box .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
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
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
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
