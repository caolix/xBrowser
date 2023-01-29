<template>
  <!--Settings Dialog-->
  <el-dialog v-model="props.visible" title="Settings" @close="onClose">
    <el-form :model="form" label-width="180px">
      <el-form-item label="Upload Concurrency">
        <el-input-number
            v-model="form.uploadConcurrency"
            :min="1"
            :max="10"
            controls-position="right"
        />
        <span style="margin-left: 20px" ty>(1 - 10)</span>
      </el-form-item>
      <el-form-item label="Download Concurrency">
        <el-input-number
            v-model="form.downloadConcurrency"
            :min="1"
            :max="10"
            controls-position="right"
        />
        <span style="margin-left: 20px" ty>(1 - 10)</span>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">Submit</el-button>
        <el-button @click="onClose">Cancel</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script lang="ts">
import {onMounted, reactive} from 'vue'
import {LoadSettings, UpdateSettings} from "../../wailsjs/go/app/App";
import {models} from "../../wailsjs/go/models";
import {ElMessage} from "element-plus";

export default {
  name: "Settings",
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, context) {
    // do not use same name with ref
    const form = reactive({
      uploadConcurrency: 1,
      downloadConcurrency: 1,
    })

    onMounted(() => {
      LoadSettings().then((res) => {
        form.uploadConcurrency = res.uploadConcurrency
        form.downloadConcurrency = res.downloadConcurrency
      })
    })

    const onSubmit = () => {
      var s = new models.Settings()
      s.uploadConcurrency = form.uploadConcurrency
      s.downloadConcurrency = form.downloadConcurrency
      UpdateSettings(s).then((res) => {
        if (res.err !== '') {
          ElMessage.error(res.err)
        }
      })
      context.emit('cancelVisible', false)
    }

    const onClose = () => {
      context.emit('cancelVisible', false)
    }

    const handleUploadConcurrencyChange = (value: number) => {

    }

    return {
      props,
      onSubmit,
      onClose,
      handleUploadConcurrencyChange,
      form,
    }
  }
}


</script>

<style scoped>

</style>