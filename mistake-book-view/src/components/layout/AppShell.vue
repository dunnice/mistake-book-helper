<template>
  <div :class="['app-shell', { 'drawer-open': drawerOpen }]">
    <!-- 侧边 Dock 胶囊导航栏 -->
    <aside
      :class="['side-rail', { 'is-docking': hoverIndex !== null }]"
      @mouseleave="hoverIndex = null"
    >
      <button
        type="button"
        class="rail-logo"
        title="错题本工作台"
        aria-label="返回错题本"
        @click="navigate('/')"
      >
        <img :src="logoUrl" alt="错题本" class="rail-logo-image" />
      </button>

      <nav class="rail-dock" aria-label="主导航">
        <button
          v-for="(item, index) in primaryNav"
          :key="item.path"
          type="button"
          :class="['rail-item', { active: isRouteActive(item.path), 'is-hot': hoverIndex === index }]"
          :style="{ '--rail-scale': getScale(index) }"
          @mouseenter="hoverIndex = index"
          @focus="hoverIndex = index"
          @blur="hoverIndex = null"
          @click="navigate(item.path)"
        >
          <span class="rail-item-icon">
            <component :is="item.icon" :size="20" :stroke-width="1.8" />
          </span>
          <span :class="['rail-item-label', { 'is-visible': hoverIndex === index }]">
            {{ item.label }}
          </span>
        </button>
      </nav>

      <div class="rail-spacer"></div>

      <!-- 底部系统设置 -->
      <button
        type="button"
        :class="['rail-item', { active: isRouteActive(settingsNav.path), 'is-hot': hoverIndex === primaryNav.length }]"
        :style="{ '--rail-scale': getScale(primaryNav.length) }"
        @mouseenter="hoverIndex = primaryNav.length"
        @focus="hoverIndex = primaryNav.length"
        @blur="hoverIndex = null"
        @click="navigate(settingsNav.path)"
      >
        <span class="rail-item-icon">
          <component :is="settingsNav.icon" :size="20" :stroke-width="1.8" />
        </span>
        <span :class="['rail-item-label', { 'is-visible': hoverIndex === primaryNav.length }]">
          {{ settingsNav.label }}
        </span>
      </button>
    </aside>

    <!-- 主工作区 -->
    <main class="workspace">
      <header class="topbar">
        <div>
          <h1>{{ currentTitle }}</h1>
          <p>{{ currentSubtitle }}</p>
        </div>
        <div class="topbar-right">
          <slot name="top-actions" />
          <div class="api-status-badge" @click="navigate('/settings')" :title="hasApiKey ? '已配置 API Key' : '未配置 API Key'">
            <span :class="['status-dot', hasApiKey ? 'online' : 'offline']"></span>
            <span class="status-text">{{ hasApiKey ? '已授权' : '待配置密钥' }}</span>
          </div>
        </div>
      </header>

      <div class="workspace-content">
        <slot />
      </div>
    </main>

    <!-- 侧边抽屉 -->
    <Teleport to="body">
      <div
        v-if="drawerOpen"
        class="drawer-overlay"
        @mousedown.self="onOverlayClick"
      >
        <div class="drawer-slide" :style="{ maxWidth: drawerWidth || '640px' }">
          <slot name="drawer" />
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { BookOpen, BarChart3, Printer, Settings, Smartphone } from 'lucide-vue-next';
import { getApiKey } from '../../api/client';
import logoUrl from '../../assets/logo.png';

const props = defineProps<{
  drawerOpen?: boolean;
  drawerWidth?: string;
}>();

const emit = defineEmits<{
  (e: 'close-drawer'): void;
}>();

const router = useRouter();
const route = useRoute();
const hoverIndex = ref<number | null>(null);
const hasApiKey = ref(Boolean(getApiKey()));

const primaryNav = [
  { path: '/', label: '错题库', icon: BookOpen },
  { path: '/stats', label: '统计看板', icon: BarChart3 },
  { path: '/print', label: '打印组卷', icon: Printer },
  { path: '/mobile', label: '移动复习', icon: Smartphone }
];

const settingsNav = { path: '/settings', label: '系统配置', icon: Settings };

function getScale(index: number): number {
  if (hoverIndex.value === null) return 1;
  const dist = Math.abs(index - hoverIndex.value);
  if (dist === 0) return 1.32;
  if (dist === 1) return 1.14;
  if (dist === 2) return 1.05;
  return 1;
}

function isRouteActive(path: string): boolean {
  if (path === '/') return route.path === '/';
  return route.path.startsWith(path);
}

function navigate(path: string) {
  router.push(path);
}

const titlesMap: Record<string, { title: string; subtitle: string }> = {
  '/': { title: '错题本 · 题库工作台', subtitle: '集中管理全科错题，科学分类、智能剖析与循环重做' },
  '/stats': { title: '错题分析与掌握看板', subtitle: '多维度透视薄弱知识点、学科分布与近期重做通过率' },
  '/print': { title: '纸质打印与线下练习', subtitle: '生成高保真 A4 错题练习卷，附带题目独立扫码二维码' },
  '/mobile': { title: '手机端移动复习', subtitle: '扫码即练，无需输入复杂密码，随时随地消灭错题' },
  '/settings': { title: '系统配置', subtitle: '管理 API Key 访问凭证、后端 API 代理地址与云端同步' }
};

const currentTitle = computed(() => {
  return titlesMap[route.path]?.title || '错题本';
});

const currentSubtitle = computed(() => {
  return titlesMap[route.path]?.subtitle || '智能错题管理与复习助手';
});

function onOverlayClick() {
  emit('close-drawer');
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.drawerOpen) {
    emit('close-drawer');
  }
}

function updateKeyStatus() {
  hasApiKey.value = Boolean(getApiKey());
}

onMounted(() => {
  window.addEventListener('keydown', onKeydown);
  window.addEventListener('mb-api-key-changed', updateKeyStatus);
});

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown);
  window.removeEventListener('mb-api-key-changed', updateKeyStatus);
});
</script>

<style scoped>
.api-status-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--border-soft);
  cursor: pointer;
  font-size: var(--font-size-xs);
  transition: all 0.2s ease;
}

.api-status-badge:hover {
  background: var(--surface-muted);
  border-color: var(--border-strong);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot.online {
  background: var(--success);
  box-shadow: 0 0 0 2px rgba(21, 128, 61, 0.2);
}

.status-dot.offline {
  background: var(--danger);
  box-shadow: 0 0 0 2px rgba(220, 38, 38, 0.2);
}

.status-text {
  color: var(--text-secondary);
  font-weight: var(--font-weight-medium);
}

.workspace-content {
  padding-bottom: 40px;
}

.rail-logo-image {
  width: 34px;
  height: 34px;
  object-fit: contain;
  border-radius: 9px;
  display: block;
}
</style>
