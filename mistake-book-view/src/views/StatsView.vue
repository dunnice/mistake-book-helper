<template>
  <div class="stats-page">
    <!-- 顶部核心指标条 -->
    <div class="stat-banner">
      <div class="stat-card">
        <div class="stat-card-title">
          <span>累计归档错题</span>
          <BookOpen :size="16" class="text-muted" />
        </div>
        <div class="stat-card-value">{{ stats.total || 0 }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-card-title">
          <span>待攻克题目</span>
          <AlertCircle :size="16" style="color: var(--danger);" />
        </div>
        <div class="stat-card-value" style="color: var(--danger);">{{ stats.pending || 0 }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-card-title">
          <span>已掌握题目</span>
          <CheckCircle2 :size="16" style="color: var(--success);" />
        </div>
        <div class="stat-card-value" style="color: var(--success);">{{ stats.mastered || 0 }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-card-title">
          <span>重做达标率</span>
          <TrendingUp :size="16" style="color: var(--primary);" />
        </div>
        <div class="stat-card-value" style="color: var(--primary);">
          {{ (stats.recent_pass_rate || 0).toFixed(0) }}%
        </div>
      </div>
    </div>

    <!-- 图表网格 -->
    <div style="display: grid; grid-template-columns: 1fr; gap: 20px; margin-bottom: 20px;">
      <!-- 高频易错知识点 -->
      <div class="panel" style="background: var(--surface-solid); border: 1px solid var(--border-soft); border-radius: var(--radius-lg); padding: 20px;">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
          <div>
            <h3 style="font-size: var(--font-size-lg); font-weight: var(--font-weight-bold); margin: 0;">
              高频薄弱知识点 TOP 10
            </h3>
            <p style="margin: 2px 0 0; font-size: var(--font-size-xs); color: var(--muted);">
              重点关注失分集中的知识模块，优先安排针对性巩固
            </p>
          </div>
          <span class="tag-badge tag-subject">易错排行</span>
        </div>
        <div ref="kpChartRef" style="width: 100%; height: 320px;"></div>
      </div>
    </div>

    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 20px;">
      <!-- 学科分布 -->
      <div class="panel" style="background: var(--surface-solid); border: 1px solid var(--border-soft); border-radius: var(--radius-lg); padding: 20px;">
        <h3 style="font-size: var(--font-size-lg); font-weight: var(--font-weight-bold); margin: 0 0 16px;">
          各学科错题占比分布
        </h3>
        <div ref="subChartRef" style="width: 100%; height: 280px;"></div>
      </div>

      <!-- 错因归因 -->
      <div class="panel" style="background: var(--surface-solid); border: 1px solid var(--border-soft); border-radius: var(--radius-lg); padding: 20px;">
        <h3 style="font-size: var(--font-size-lg); font-weight: var(--font-weight-bold); margin: 0 0 16px;">
          错误原因根因剖析
        </h3>
        <div ref="reasonChartRef" style="width: 100%; height: 280px;"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue';
import { BookOpen, AlertCircle, CheckCircle2, TrendingUp } from 'lucide-vue-next';
import * as echarts from 'echarts';
import { api } from '../api/client';
import type { StatsOverview, KpStat, SubjectStat, ReasonStat } from '../types';

const stats = ref<Partial<StatsOverview>>({});
const kpChartRef = ref<HTMLDivElement | null>(null);
const subChartRef = ref<HTMLDivElement | null>(null);
const reasonChartRef = ref<HTMLDivElement | null>(null);

let kpChartInstance: echarts.ECharts | null = null;
let subChartInstance: echarts.ECharts | null = null;
let reasonChartInstance: echarts.ECharts | null = null;

async function loadStats() {
  try {
    const data = await api.getStats();
    stats.value = data;
    await nextTick();
    renderCharts(data);
  } catch (err) {
    console.error(err);
  }
}

function renderCharts(data: StatsOverview) {
  // 1. 知识点排行
  if (kpChartRef.value) {
    if (!kpChartInstance) kpChartInstance = echarts.init(kpChartRef.value);
    const kps = (data.kp_stats || []).slice(0, 10).reverse();
    kpChartInstance.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: { type: 'shadow' }
      },
      grid: {
        left: 140,
        right: 40,
        top: 20,
        bottom: 20
      },
      xAxis: {
        type: 'value',
        axisLine: { show: false },
        splitLine: { lineStyle: { color: '#f1f5f9' } }
      },
      yAxis: {
        type: 'category',
        data: kps.map((k: KpStat) => k.kp || '未知知识点'),
        axisLabel: { width: 120, overflow: 'truncate', color: '#475569' },
        axisLine: { lineStyle: { color: '#e2e8f0' } }
      },
      series: [
        {
          name: '错题数量',
          type: 'bar',
          data: kps.map((k: KpStat) => k.count),
          barMaxWidth: 20,
          itemStyle: {
            borderRadius: [0, 4, 4, 0],
            color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
              { offset: 0, color: '#3b82f6' },
              { offset: 1, color: '#60a5fa' }
            ])
          }
        }
      ]
    });
  }

  // 2. 学科分布
  if (subChartRef.value) {
    if (!subChartInstance) subChartInstance = echarts.init(subChartRef.value);
    const subData = (data.by_subject || []).map((x: SubjectStat) => ({
      name: x.name || '未分类',
      value: x.count
    }));
    subChartInstance.setOption({
      tooltip: { trigger: 'item' },
      legend: { bottom: 0, icon: 'circle' },
      color: ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6', '#ec4899', '#06b6d4', '#64748b'],
      series: [
        {
          name: '错题学科',
          type: 'pie',
          radius: ['45%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 6,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: { show: false },
          data: subData
        }
      ]
    });
  }

  // 3. 错因剖析
  if (reasonChartRef.value) {
    if (!reasonChartInstance) reasonChartInstance = echarts.init(reasonChartRef.value);
    const reasonData = (data.by_reason || []).map((x: ReasonStat) => ({
      name: x.name || '未归类',
      value: x.count
    }));
    reasonChartInstance.setOption({
      tooltip: { trigger: 'item' },
      legend: { bottom: 0, icon: 'circle' },
      color: ['#ef4444', '#f97316', '#eab308', '#6366f1', '#a855f7', '#94a3b8'],
      series: [
        {
          name: '错误原因',
          type: 'pie',
          radius: ['45%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 6,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: { show: false },
          data: reasonData
        }
      ]
    });
  }
}

function handleResize() {
  kpChartInstance?.resize();
  subChartInstance?.resize();
  reasonChartInstance?.resize();
}

onMounted(() => {
  loadStats();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  kpChartInstance?.dispose();
  subChartInstance?.dispose();
  reasonChartInstance?.dispose();
});
</script>
