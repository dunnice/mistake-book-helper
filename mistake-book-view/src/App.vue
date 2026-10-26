<template>
  <!-- 如果是手机扫码专用的 /m/:id 页面，提供沉浸式纯移动端视图 -->
  <div v-if="isPureMobileRoute" style="min-height: 100vh; padding: 16px; background: var(--app-bg);">
    <router-view />
  </div>

  <!-- 常规桌面端完整 AppShell 工作台 -->
  <AppShell v-else>
    <router-view />
  </AppShell>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import AppShell from './components/layout/AppShell.vue';

const route = useRoute();
const isPureMobileRoute = computed(() => {
  return route.path.startsWith('/m/');
});
</script>
