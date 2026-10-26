import { createApp } from 'vue';
import App from './App.vue';
import router from './router';

// 导入视觉设计系统与业务样式
import './styles/global.css';
import './styles/mistakebook.css';

const app = createApp(App);
app.use(router);
app.mount('#app');
