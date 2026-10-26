import { createRouter, createWebHistory } from 'vue-router';
import QuestionsView from '../views/QuestionsView.vue';
import StatsView from '../views/StatsView.vue';
import PrintView from '../views/PrintView.vue';
import MobileStudyView from '../views/MobileStudyView.vue';
import SettingsView from '../views/SettingsView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'questions',
      component: QuestionsView
    },
    {
      path: '/stats',
      name: 'stats',
      component: StatsView
    },
    {
      path: '/print',
      name: 'print',
      component: PrintView
    },
    {
      path: '/mobile',
      name: 'mobile',
      component: MobileStudyView
    },
    {
      path: '/m/:id',
      name: 'mobile-question',
      component: MobileStudyView
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsView
    }
  ]
});

export default router;
