<template>
  <div class="print-page">
    <!-- 打印控制面板（打印时自动隐藏） -->
    <div class="filter-bar print-controls" style="margin-bottom: 24px;">
      <div class="filter-row" style="justify-content: space-between;">
        <div style="display: flex; gap: 10px; align-items: center; flex-wrap: wrap;">
          <span class="filter-label">科目</span>
          <div class="chip-list">
            <span
              :class="['chip-item', { active: selectedSubject === '' }]"
              @click="selectedSubject = ''; loadData()"
            >全部</span>
            <span
              v-for="s in subjects"
              :key="s"
              :class="['chip-item', { active: selectedSubject === s }]"
              @click="selectedSubject = s; loadData()"
            >{{ s }}</span>
          </div>
        </div>

        <div style="display: flex; gap: 10px; align-items: center;">
          <label style="font-size: var(--font-size-sm); color: var(--text-secondary); display: flex; align-items: center; gap: 6px;">
            抽取数量:
            <input
              v-model.number="limitCount"
              type="number"
              min="1"
              max="50"
              style="width: 70px; padding: 4px 8px;"
              @change="loadData"
            />
          </label>
          <button type="button" class="primary-button" @click="triggerPrint">
            <Printer :size="15" /> 发送打印 (A4)
          </button>
        </div>
      </div>

      <div class="filter-row">
        <span class="filter-label">知识点</span>
        <div class="search-input-wrap" style="max-width: 240px;">
          <input
            v-model="kpFilter"
            type="text"
            placeholder="如：立体几何..."
            @keyup.enter="loadData"
          />
        </div>
        <button type="button" class="ghost-button" @click="loadData">
          <RefreshCw :size="14" /> 重新组卷
        </button>
        <span style="font-size: var(--font-size-xs); color: var(--muted); margin-left: 8px;">
          💡 每道题右上角印有独立二维码，学生做完后手机扫码即可对照 AI 详细解析并打卡。
        </span>
      </div>
    </div>

    <!-- A4 纸张排版仿真区 -->
    <div class="print-sheet">
      <div style="text-align: center; border-bottom: 2px solid #111827; padding-bottom: 16px; margin-bottom: 24px;">
        <h2 style="font-size: 24px; font-weight: bold; margin: 0 0 6px; letter-spacing: 2px;">
          {{ selectedSubject ? `${selectedSubject}学科 · ` : '' }}错题强化巩固练习卷
        </h2>
        <div style="font-size: 13px; color: #4b5563; display: flex; justify-content: center; gap: 24px;">
          <span>卷面题数：{{ questions.length }} 题</span>
          <span>生成时间：{{ printDate }}</span>
          <span>姓名：______________</span>
          <span>得分：________</span>
        </div>
      </div>

      <div v-if="questions.length === 0" style="text-align: center; color: #9ca3af; padding: 60px 0;">
        未检索到符合条件的错题，请调整筛选条件后重试。
      </div>

      <div
        v-for="(q, idx) in questions"
        :key="q.id"
        class="print-qitem"
      >
        <div style="display: flex; justify-content: space-between; align-items: flex-start; gap: 16px;">
          <div style="flex: 1;">
            <div style="font-size: 15px; font-weight: bold; margin-bottom: 6px;">
              第 {{ idx + 1 }} 题.
              <span style="font-size: 12px; font-weight: normal; color: #6b7280; margin-left: 6px;">
                [{{ q.subject }} · {{ q.source || '日常纠错' }}]
              </span>
              <span
                v-if="q.knowledge_points && q.knowledge_points.length"
                style="font-size: 12px; font-weight: normal; color: #4b5563; margin-left: 8px;"
              >
                考点：{{ q.knowledge_points.join('、') }}
              </span>
            </div>

            <!-- 题干 -->
            <div style="font-size: 15px; line-height: 1.8; white-space: pre-wrap; margin-bottom: 8px;">
              {{ q.stem }}
            </div>

            <!-- 选项（若有） -->
            <div
              v-if="q.options && q.options.length"
              style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 8px; margin: 8px 0;"
            >
              <div
                v-for="(opt, optIdx) in q.options"
                :key="optIdx"
                style="font-size: 14px; line-height: 1.6;"
              >
                <strong>{{ 'ABCD'[optIdx] || optIdx + 1 }}.</strong> {{ opt }}
              </div>
            </div>

            <!-- 原题截图（若有） -->
            <div v-if="q.image_url" style="margin: 10px 0; max-width: 400px;">
              <img :src="getImageUrl(q.image_url)" style="max-width: 100%; border: 1px solid #e5e7eb; border-radius: 4px;" alt="题目插图" />
            </div>

            <!-- 预留解题与草稿空白 -->
            <div class="print-answer-area">
              <span style="font-size: 11px; color: #9ca3af; margin: 6px 10px; display: inline-block;">答题与草稿区：</span>
            </div>
          </div>

          <!-- 右上角手机扫码答题二维码 -->
          <div style="text-align: center; min-width: 84px;">
            <img
              :src="api.getQrCodeUrl(q.id)"
              style="width: 80px; height: 80px; border: 1px solid #e5e7eb; border-radius: 4px; display: block;"
              alt="扫码答题"
            />
            <span style="font-size: 10px; color: #6b7280; display: block; margin-top: 3px;">
              扫码看解析
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Printer, RefreshCw } from 'lucide-vue-next';
import { api, getImageUrl } from '../api/client';
import type { Question } from '../types';

const subjects = ['数学', '物理', '化学', '英语', '语文', '道法', '历史', '地理', '生物'];
const selectedSubject = ref('');
const kpFilter = ref('');
const limitCount = ref(10);
const questions = ref<Question[]>([]);
const printDate = new Date().toLocaleDateString('zh-CN');

async function loadData() {
  try {
    const res = await api.listQuestions({
      subject: selectedSubject.value,
      kp: kpFilter.value,
      status: '' // 打印可以默认包含所有或待重做
    });
    const all = res.questions || [];
    // 优先取「待重做」题目
    const sorted = all.sort((a: Question, b: Question) => (a.status === '待重做' ? -1 : 1));
    questions.value = sorted.slice(0, limitCount.value);
  } catch (err) {
    console.error(err);
  }
}

function triggerPrint() {
  window.print();
}

onMounted(() => {
  loadData();
});
</script>
