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
        uploadList: [],
        uploadProgress: {},
        downloadList: [],
    },
    mutations: {
        addToUploadList(state, payload) {
            state.uploadList.push(payload.file)
            state.uploadProgress[payload.progress] = 0
        },
        updateProgress(state, payload) {
            state.uploadProgress[payload.progress] = payload.data
        }
    }
})

app.use(ElementPlus).use(router).use(store).mount('#app')
