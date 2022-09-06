import {createApp} from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import router from "./router";
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import {createStore} from 'vuex'

const app = createApp(App)
// app.directive('loadmore', {
//     mounted(el, binding) {
//         let tbody = el.querySelector(".el-table__body-wrapper");
//         el.tableInfiniteScrollFn = function () {
//             if (this.scrollHeight - this.scrollTop - parseInt(this.style.height)  === 0) {
//                 binding.value();
//             }
//         };
//         tbody.addEventListener("scroll", el.tableInfiniteScrollFn);
//         tbody = undefined;
//     },
// })
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
        uploadSuccessCount: 0,
        uploadTotalSize: 0,
        uploadProgress: {}, // eventProgressId -> value

        downloadList: [],
        downloadProgress: {}, // eventProgressId -> value
    },
    mutations: {
        addToUploadList(state, payload) {
            state.uploadList.push(payload.file)
            state.uploadProgress[payload.progress] = 0
            state.uploadTotalSize += payload.file.size
        },
        updateUploadStatus(state, payload) {
            state.uploadList[payload.index].status = payload.status
        },
        updateUploadProgress(state, payload) {
            if (payload.data === 100) {
                state.uploadSuccessCount ++
            }
            state.uploadProgress[payload.progress] = payload.data
        },
        clearUploadFinished(state) {
            var start = 0
            var delCount = 0
            var uploadTotalSize = 0
            for (let i = 0; i < state.uploadList.length; i++) {
                if (state.uploadProgress[state.uploadList[i].taskId] === 100) {
                    delCount ++
                    uploadTotalSize += state.uploadList[i].size
                } else if (delCount === 0){
                    start ++
                } else {
                    state.uploadTotalSize -= uploadTotalSize
                    state.uploadSuccessCount -= delCount
                    state.uploadList.splice(start, delCount)
                    i = start
                    uploadTotalSize = 0
                    delCount = 0
                }
            }
            if (delCount > 0) {
                state.uploadTotalSize -= uploadTotalSize
                state.uploadSuccessCount -= delCount
                state.uploadList.splice(start, delCount)
            }
            // FIXME: add lock to update array
            if (state.uploadSuccessCount < 0) {
                state.uploadSuccessCount = 0
            }
            if (state.uploadTotalSize < 0) {
                state.uploadTotalSize = 0
            }

        },
        removeUploadListParam(state, payload) {
            state.uploadTotalSize -= state.uploadList[payload.index].size
            state.uploadList.splice(payload.index, 1)
            if (state.uploadProgress[payload.progress] === 100) {
                state.uploadSuccessCount --
            }
            delete state.uploadProgress[payload.progress]
        },

        addToDownloadList(state, payload) {
            state.downloadList.push(payload.file)
            state.downloadProgress[payload.progress] = 0
        },
        updateDownloadStatus(state, payload) {
            state.downloadList[payload.index].status = payload.status
        },
        updateDownloadProgress(state, payload) {
            state.downloadProgress[payload.progress] = payload.data
        },
        removeDownloadListParam(state, payload) {
            state.downloadList.splice(payload.index, 1)
            delete state.downloadProgress[payload.progress]
        },
    }
})

app.use(ElementPlus).use(router).use(store).mount('#app')
