import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import api from '@/utils/request'
import './assets/main.css'

const app = createApp(App)

app.use(router)

app.mount('#app')
app.config.globalProperties.$api = api;
