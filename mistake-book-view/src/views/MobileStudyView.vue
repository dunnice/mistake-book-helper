<template>
  <div class="mobile-study-page">
    <div style="max-width: 580px; margin: 0 auto;">
      <!-- 页面顶部品牌 Logo -->
      <div style="display: flex; align-items: center; justify-content: center; gap: 8px; margin-bottom: 16px;">
        <img :src="logoUrl" alt="错题本" style="width: 32px; height: 32px; border-radius: 8px; object-fit: contain; box-shadow: 0 2px 8px rgba(0,0,0,0.06);" />
        <span style="font-size: 17px; font-weight: 700; color: var(--text); letter-spacing: 0.5px;">错题本 · 移动打卡</span>
      </div>

      <!-- 若没有指定题目 ID，提供快速输入 ID 模拟扫码 -->
      <div v-if="!currentQuestionId" class="panel" style="background: var(--surface-solid); border: 1px solid var(--border-soft); border-radius: var(--radius-lg); padding: 24px; text-align: center;">
        <Smartphone :size="40" style="color: var(--primary); margin-bottom: 12px;" />
        <h3 style="font-size: var(--font-size-lg); margin-bottom: 6px;">移动端极简做题模式</h3>
        <p style="font-size: var(--font-size-sm); color: var(--muted); margin-bottom: 18px;">
          此页面专门针对手机扫码设计（免输密码）。你可以输入错题编号快速进入答题体验：
        </p>
        <div style="display: flex; gap: 8px; justify-content: center; max-width: 320px; margin: 0 auto;">
          <input
            v-model="inputQId"
            type="number"
            placeholder="错题 ID（如 1）"
            style="flex: 1;"
            @keyup.enter="loadQuestion(Number(inputQId))"
          />
          <button type="button" class="primary-button" @click="loadQuestion(Number(inputQId))">开始答题</button>
        </div>
      </div>

      <!-- 答题卡片 -->
      <div v-else-if="question" class="panel" style="background: var(--surface-solid); border: 1px solid var(--border-soft); border-radius: var(--radius-lg); padding: 24px; box-shadow: var(--shadow-md);">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
          <span class="tag-badge tag-subject">{{ question.subject }}</span>
          <span style="font-size: var(--font-size-xs); color: var(--muted);">{{ question.source || '错题强化' }}</span>
        </div>

        <!-- 题干 -->
        <div style="font-size: 16px; line-height: 1.8; color: var(--text); margin-bottom: 16px; white-space: pre-wrap;">
          {{ question.stem }}
        </div>

        <!-- 原图 -->
        <div v-if="question.image_url" style="margin-bottom: 16px;">
          <img :src="getImageUrl(question.image_url)" style="max-width: 100%; border-radius: 8px; border: 1px solid var(--border);" alt="错题图片" />
        </div>

        <!-- 选择题选项 -->
        <div v-if="question.options && question.options.length" class="option-list" style="margin-bottom: 20px;">
          <div
            v-for="(opt, idx) in question.options"
            :key="idx"
            class="option-item"
            style="padding: 10px 14px;"
          >
            <strong style="color: var(--primary); margin-right: 8px;">{{ 'ABCD'[idx] || idx + 1 }}.</strong>
            {{ opt }}
          </div>
        </div>

        <!-- 答案揭晓与详细解析区 -->
        <div v-if="showAnswer" style="margin-top: 20px; padding-top: 20px; border-top: 1px dashed var(--border);">
          <div style="display: flex; gap: 12px; margin-bottom: 16px;">
            <div style="flex: 1;">
              <span style="font-size: var(--font-size-xs); color: var(--muted); display: block; margin-bottom: 4px;">正确答案</span>
              <div class="answer-badge-correct" style="font-size: 16px; width: 100%; text-align: center;">
                {{ question.answer }}
              </div>
            </div>
          </div>

          <div style="margin-bottom: 20px;">
            <span style="font-size: var(--font-size-xs); color: var(--muted); display: block; margin-bottom: 6px;">AI 考点解析</span>
            <div class="analysis-box">
              {{ question.analysis || '暂无详细分析' }}
            </div>
          </div>

          <!-- 提交重做结果 -->
          <div v-if="!submitted">
            <p style="text-align: center; font-size: var(--font-size-sm); color: var(--text-secondary); margin-bottom: 12px;">
              对照答案后，这道题本次做对了吗？
            </p>
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px;">
              <button type="button" class="success-button" style="padding: 12px; height: 42px;" @click="submitResult('对')">
                <Check :size="16" /> 我做对了
              </button>
              <button type="button" class="danger-button" style="padding: 12px; height: 42px;" @click="submitResult('错')">
                <X :size="16" /> 还是做错了
              </button>
            </div>
          </div>
          <div v-else style="text-align: center; padding: 20px; background: var(--success-soft); border-radius: var(--radius-md); border: 1px solid rgba(21, 128, 61, 0.2);">
            <div style="font-size: 24px; color: var(--success); margin-bottom: 6px;">✔</div>
            <div style="font-weight: var(--font-weight-bold); color: var(--success);">已成功记录本次重做！</div>
            <div style="font-size: var(--font-size-xs); color: var(--muted); margin-top: 4px;">当前掌握状态: {{ question.status }}</div>
          </div>
        </div>

        <button
          v-else
          type="button"
          class="primary-button"
          style="width: 100%; height: 42px; font-size: 15px; margin-top: 10px;"
          @click="showAnswer = true"
        >
          先独立思考做完，点击查看答案与解析
        </button>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="errorMsg" class="panel" style="text-align: center; padding: 30px; color: var(--danger);">
        {{ errorMsg }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { Smartphone, Check, X } from 'lucide-vue-next';
import { api, getImageUrl } from '../api/client';
import type { Question } from '../types';
import logoUrl from '../assets/logo.png';

const route = useRoute();
const currentQuestionId = ref<number | null>(null);
const inputQId = ref('');
const question = ref<Question | null>(null);
const showAnswer = ref(false);
const submitted = ref(false);
const errorMsg = ref('');

async function loadQuestion(id: number) {
  if (!id) return;
  currentQuestionId.value = id;
  showAnswer.value = false;
  submitted.value = false;
  errorMsg.value = '';

  try {
    const res = await api.getMobileQuestion(id);
    question.value = res.question;
  } catch (err: any) {
    errorMsg.value = err.message || '加载错题失败';
  }
}

async function submitResult(result: '对' | '错') {
  if (!currentQuestionId.value) return;
  try {
    const res = await api.submitMobileReview(currentQuestionId.value, result);
    submitted.value = true;
    if (question.value && res.question) {
      question.value.status = res.question.status;
    }
  } catch (err: any) {
    alert(err.message || '提交失败');
  }
}

onMounted(() => {
  const paramId = route.params.id;
  if (paramId) {
    loadQuestion(Number(paramId));
  }
});
</script>
