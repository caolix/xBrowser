import {createApp} from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import router from "./router";
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import {createStore} from 'vuex'

const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
// 创建一个新的 store 实例
const store = createStore({
    state: {
        //uploadList param {
        // name: task.name,
        // bucket: task.bucket,
        // key: task.key,
        // source: task.source,
        // uploadedSize: task.uploadedSize,
        // size: task.size,
        // humanSize: task.humanSize,
        // isPending: true,
        // progress: task.taskId,
        // accountId:  task.accountId,
        // isMultipart: task.isMultipart,
        // partSize: task.partSize,
        // taskId: task.taskId,
        // uploadId: task.uploadId
        //  }
        uploadList: [],
        uploadProgress: {}, // eventProgressId -> value
        downloadList: [],
        listTask: false
    },
    mutations: {
        addToUploadList(state, payload) {
            state.uploadList.push(payload.file)
            state.uploadProgress[payload.progress] = 0
        },
        updateProgress(state, payload) {
            state.uploadProgress[payload.progress] = payload.data
        },
        removeUploadListParam(state, payload) {
            state.uploadList.splice(payload.index, 1)
            delete state.uploadProgress[payload.progress]
        },
    }
})

app.use(ElementPlus).use(router).use(store).mount('#app')
