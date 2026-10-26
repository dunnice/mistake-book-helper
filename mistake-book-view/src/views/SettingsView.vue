<template>
  <div class="settings-page" style="max-width: 720px; margin: 0 auto;">
    <!-- 凭证与连接配置卡片 -->
    <div class="panel" style="background: var(--surface-solid); border: 1px solid var(--border-soft); border-radius: var(--radius-lg); padding: 24px; margin-bottom: 24px; box-shadow: var(--shadow-sm);">
      <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 16px;">
        <KeyRound :size="20" style="color: var(--primary);" />
        <h3 style="font-size: var(--font-size-lg); font-weight: var(--font-weight-bold); margin: 0;">
          API Key 访问密钥
        </h3>
      </div>
      <p style="font-size: var(--font-size-sm); color: var(--muted); line-height: 1.6; margin-bottom: 16px;">
        错题本接口需通过 <code>X-API-Key</code> 头部进行鉴权。该密钥需与服务端环境变量 <code>API_KEY</code> 保持一致，密钥安全保存在当前浏览器的 localStorage 中。
      </p>

      <div style="margin-bottom: 16px;">
        <label class="detail-label">API Key 密钥</label>
        <div style="position: relative;">
          <input
            v-model="keyInput"
            :type="revealApiKey ? 'text' : 'password'"
            placeholder="输入服务端配置的 API_KEY..."
            style="width: 100%; padding-right: 40px;"
          />
          <button
            type="button"
            class="icon-button"
            style="position: absolute; right: 8px; top: 50%; transform: translateY(-50%);"
            @click="toggleRevealKey"
          >
            <Eye v-if="!revealApiKey" :size="16" />
            <EyeOff v-else :size="16" />
          </button>
        </div>
      </div>

      <div style="margin-bottom: 20px;">
        <label class="detail-label">后端服务地址 (API Base URL)</label>
        <input
          v-model="baseUrlInput"
          type="text"
          placeholder="留空表示同源相对路径（推荐），或输入如 http://127.0.0.1:8080"
          style="width: 100%;"
        />
        <span style="font-size: var(--font-size-xs); color: var(--muted); margin-top: 4px; display: block;">
          单体部署或在同一域名/反代路径时保持为空即可；跨域调试时填入后端完整地址。
        </span>
      </div>

      <div style="display: flex; gap: 12px; align-items: center;">
        <button type="button" class="primary-button" @click="saveSettings">
          <Save :size="15" /> 保存配置
        </button>
        <button type="button" class="ghost-button" :disabled="testing" @click="testConnection">
          <Activity :size="15" /> 测试连通性
        </button>
        <span v-if="testResult" :style="{ color: testSuccess ? 'var(--success)' : 'var(--danger)', fontSize: 'var(--font-size-sm)', fontWeight: 500 }">
          {{ testResult }}
        </span>
      </div>
    </div>

    <!-- 架构与工程规范卡片 -->
    <div class="panel" style="background: var(--surface-solid); border: 1px solid var(--border-soft); border-radius: var(--radius-lg); padding: 24px; box-shadow: var(--shadow-sm);">
      <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 16px;">
        <Cpu :size="20" style="color: var(--primary);" />
        <h3 style="font-size: var(--font-size-lg); font-weight: var(--font-weight-bold); margin: 0;">
          系统架构与规范
        </h3>
      </div>
      <div style="font-size: var(--font-size-sm); color: var(--text-secondary); line-height: 1.8;">
        <p><strong>出品人</strong>：ice7</p>
        <p><strong>前端项目</strong>：<code>mistake-book-view</code> (Vue 3 + TypeScript + Vite + Lucide + ECharts)</p>
        <p><strong>后端项目</strong>：<code>mistake-book-go-api</code> (Go + Gin + GORM + SQLite / MySQL)</p>
        <p><strong>项目文档</strong>：<code>00-doc/</code> (实施方案、设计文档、运维手册、AI Agent 工作流)</p>
        <p><strong>构建脚本</strong>：<code>script/</code> (一键编译前后端与自动化部署)</p>
        <p><strong>AI 智能协同</strong>：手机端拍照发送至 AI 智能助手，触发「收集到错题本」技能，即可全自动清洗入库。</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { KeyRound, Eye, EyeOff, Save, Activity, Cpu } from 'lucide-vue-next';
import { getApiKey, setApiKey, getBaseUrl, setBaseUrl, api } from '../api/client';

const keyInput = ref('');
const baseUrlInput = ref('');
const revealApiKey = ref(false);
const testing = ref(false);
const testResult = ref('');
const testSuccess = ref(false);

function toggleRevealKey() {
  revealApiKey.value = !revealApiKey.value;
}

function saveSettings() {
  setApiKey(keyInput.value);
  setBaseUrl(baseUrlInput.value);
  testResult.value = '配置已保存到本地存储';
  testSuccess.value = true;
}

async function testConnection() {
  saveSettings();
  testing.value = true;
  testResult.value = '正在检测连通性...';
  try {
    const stats = await api.getStats();
    testSuccess.value = true;
    testResult.value = `连接成功！当前题库已收录 ${stats.total || 0} 道错题。`;
  } catch (err: any) {
    testSuccess.value = false;
    testResult.value = `连接失败: ${err.message}`;
  } finally {
    testing.value = false;
  }
}

onMounted(() => {
  keyInput.value = getApiKey();
  baseUrlInput.value = getBaseUrl();
});
</script>
