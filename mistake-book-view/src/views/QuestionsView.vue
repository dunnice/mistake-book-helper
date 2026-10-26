<template>
  <div class="questions-page">
    <!-- 统计概览条 -->
    <div class="stat-banner">
      <div class="stat-card" @click="setFilterStatus('')" style="cursor: pointer;">
        <div class="stat-card-title">
          <span>累计错题</span>
          <BookOpen :size="16" class="text-muted" />
        </div>
        <div class="stat-card-value">{{ stats.total || 0 }}</div>
      </div>
      <div class="stat-card" @click="setFilterStatus('待重做')" style="cursor: pointer;">
        <div class="stat-card-title">
          <span>待消灭错题</span>
          <AlertCircle :size="16" style="color: var(--danger);" />
        </div>
        <div class="stat-card-value" style="color: var(--danger);">{{ stats.pending || 0 }}</div>
      </div>
      <div class="stat-card" @click="setFilterStatus('已掌握')" style="cursor: pointer;">
        <div class="stat-card-title">
          <span>已掌握题目</span>
          <CheckCircle2 :size="16" style="color: var(--success);" />
        </div>
        <div class="stat-card-value" style="color: var(--success);">{{ stats.mastered || 0 }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-card-title">
          <span>近30天重做通过率</span>
          <TrendingUp :size="16" style="color: var(--primary);" />
        </div>
        <div class="stat-card-value" style="color: var(--primary);">
          {{ (stats.recent_pass_rate || 0).toFixed(0) }}%
        </div>
      </div>
    </div>

    <!-- 过滤器与工具栏 -->
    <div class="filter-bar">
      <!-- 科目过滤 -->
      <div class="filter-row">
        <span class="filter-label">科目</span>
        <div class="chip-list">
          <span
            :class="['chip-item', { active: filters.subject === '' }]"
            @click="setFilterSubject('')"
          >全部</span>
          <span
            v-for="sub in subjects"
            :key="sub"
            :class="['chip-item', { active: filters.subject === sub }]"
            @click="setFilterSubject(sub)"
          >{{ sub }}</span>
        </div>
      </div>

      <!-- 状态与搜索 -->
      <div class="filter-row">
        <span class="filter-label">状态</span>
        <div class="chip-list">
          <span
            :class="['chip-item', { active: filters.status === '' }]"
            @click="setFilterStatus('')"
          >全部</span>
          <span
            :class="['chip-item', { active: filters.status === '待重做' }]"
            @click="setFilterStatus('待重做')"
          >待重做</span>
          <span
            :class="['chip-item', { active: filters.status === '已重做' }]"
            @click="setFilterStatus('已重做')"
          >已重做</span>
          <span
            :class="['chip-item', { active: filters.status === '已掌握' }]"
            @click="setFilterStatus('已掌握')"
          >已掌握</span>
        </div>

        <div class="search-input-wrap" style="max-width: 160px;">
          <input
            v-model="filters.kp"
            type="text"
            placeholder="知识点..."
            @keyup.enter="fetchQuestions"
          />
        </div>

        <div class="search-input-wrap has-icon">
          <Search :size="15" />
          <input
            v-model="filters.q"
            class="search-input"
            type="text"
            placeholder="搜索题干、来源、解析..."
            @keyup.enter="fetchQuestions"
          />
        </div>

        <button type="button" class="primary-button" @click="fetchQuestions">
          <Search :size="14" /> 查询
        </button>
        <button type="button" class="ghost-button" @click="resetFilters">重置</button>
      </div>
    </div>

    <!-- 列表顶部操作条 -->
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
      <div style="font-size: var(--font-size-sm); color: var(--muted);">
        共找到 <strong style="color: var(--text);">{{ questions.length }}</strong> 道错题
      </div>
      <div style="display: flex; gap: 8px;">
        <button type="button" class="primary-button" @click="openCreateModal">
          <Plus :size="15" /> 手动录入错题
        </button>
        <button type="button" class="icon-button" @click="fetchQuestions" title="刷新列表">
          <RefreshCw :size="15" :class="{ 'spin-anim': loading }" />
        </button>
      </div>
    </div>

    <!-- 加载出错提示 -->
    <div v-if="loadError" class="panel" style="background: var(--danger-soft); border: 1px solid rgba(220, 38, 38, 0.2); border-radius: var(--radius-lg); padding: 18px 20px; margin-bottom: 20px; display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 12px;">
      <div style="display: flex; align-items: center; gap: 10px;">
        <AlertCircle :size="20" style="color: var(--danger); flex-shrink: 0;" />
        <div>
          <div style="font-weight: 600; color: var(--danger); font-size: var(--font-size-sm);">加载错题数据失败</div>
          <div style="font-size: var(--font-size-xs); color: #7f1d1d; margin-top: 2px;">{{ loadError }}</div>
        </div>
      </div>
      <div style="display: flex; gap: 8px;">
        <button type="button" class="ghost-button compact-button" @click="fetchQuestions">重试</button>
        <button type="button" class="primary-button compact-button" @click="$router.push('/settings')">系统设置 (API Key)</button>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="!loading && !loadError && questions.length === 0" class="panel" style="text-align: center; padding: 60px 20px;">
      <div style="font-size: 40px; margin-bottom: 12px;">📚</div>
      <h3 style="font-size: var(--font-size-lg); color: var(--text); margin-bottom: 6px;">暂无匹配错题</h3>
      <p style="color: var(--muted); font-size: var(--font-size-sm); max-width: 440px; margin: 0 auto 16px;">
        你可以点击上方「手动录入错题」，或使用 AI 智能助手拍摄试卷后说「收集到错题本」，错题会自动解析入库。
      </p>
      <button type="button" class="primary-button" @click="openCreateModal">立即录入第一道错题</button>
    </div>

    <!-- 错题卡片网格 -->
    <div v-else class="question-grid">
      <div
        v-for="q in questions"
        :key="q.id"
        class="question-card"
        @click="openDetail(q.id)"
      >
        <div>
          <div class="question-card-header">
            <div style="display: flex; gap: 6px; flex-wrap: wrap;">
              <span class="tag-badge tag-subject">{{ q.subject || '未分类' }}</span>
              <span v-if="q.error_reason" class="tag-badge tag-reason">{{ q.error_reason }}</span>
              <span
                v-for="kp in (q.knowledge_points || []).slice(0, 2)"
                :key="kp"
                class="tag-badge tag-kp"
              >{{ kp }}</span>
              <span v-if="(q.knowledge_points || []).length > 2" class="tag-badge tag-kp">
                +{{ q.knowledge_points!.length - 2 }}
              </span>
            </div>
            <span :class="['status-pill', getStatusClass(q.status)]">
              {{ q.status }}
            </span>
          </div>

          <div class="question-stem">{{ q.stem }}</div>
        </div>

        <div class="question-card-footer">
          <span>{{ q.source || '练习题' }}</span>
          <span>{{ formatDate(q.created_at) }}</span>
        </div>
      </div>
    </div>

    <!-- 侧边抽屉：错题详情 -->
    <Teleport to="body">
      <div
        v-if="detailVisible && currentDetail"
        class="drawer-overlay"
        @mousedown.self="detailVisible = false"
      >
        <div class="drawer-slide" style="max-width: 680px;">
          <aside class="drawer">
            <div class="drawer-header">
              <div class="drawer-header-top">
                <span class="eyebrow">{{ currentDetail.question.subject }} · 错题详情</span>
                <button type="button" class="drawer-close-btn" aria-label="关闭详情" @click="detailVisible = false">
                  <X :size="18" />
                </button>
              </div>
              <h2 style="margin: 6px 0 0; font-size: 20px; font-weight: 700; color: var(--text);">
                错题 #{{ currentDetail.question.id }}
              </h2>
              <p style="margin: 4px 0 0; font-size: 13px; color: var(--muted);">
                来源：{{ currentDetail.question.source || '日常纠错' }}
              </p>
            </div>

          <div class="drawer-body">
            <!-- 状态与标签 -->
            <div class="detail-section" style="display: flex; gap: 8px; align-items: center; flex-wrap: wrap;">
              <span :class="['status-pill', getStatusClass(currentDetail.question.status)]">
                当前状态：{{ currentDetail.question.status }}
              </span>
              <span class="tag-badge tag-subject">{{ currentDetail.question.subject }}</span>
              <span v-if="currentDetail.question.error_reason" class="tag-badge tag-reason">
                错因：{{ currentDetail.question.error_reason }}
              </span>
              <span v-if="currentDetail.question.difficulty" class="tag-badge tag-kp">
                难度：{{ '★'.repeat(currentDetail.question.difficulty) }}
              </span>
            </div>

            <!-- 题干 -->
            <div class="detail-section">
              <div class="detail-label">
                <FileText :size="15" /> 题干内容
              </div>
              <div class="detail-stem-box">{{ currentDetail.question.stem }}</div>
            </div>

            <!-- 选项（若为选择题） -->
            <div
              v-if="currentDetail.question.options && currentDetail.question.options.length"
              class="detail-section"
            >
              <div class="detail-label">
                <List :size="15" /> 选择题选项
              </div>
              <div class="option-list">
                <div
                  v-for="(opt, idx) in currentDetail.question.options"
                  :key="idx"
                  class="option-item"
                >
                  <strong style="color: var(--primary); margin-right: 6px;">
                    {{ 'ABCD'[idx] || idx + 1 }}.
                  </strong>
                  {{ opt }}
                </div>
              </div>
            </div>

            <!-- 答案对比 -->
            <div class="detail-section" style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px;">
              <div>
                <div class="detail-label">
                  <XCircle :size="15" style="color: var(--danger);" /> 我的错误答案
                </div>
                <div class="answer-badge-wrong" style="width: 100%; text-align: center;">
                  {{ currentDetail.question.my_answer || '未填写' }}
                </div>
              </div>
              <div>
                <div class="detail-label">
                  <CheckCircle2 :size="15" style="color: var(--success);" /> 正确答案
                </div>
                <div class="answer-badge-correct" style="width: 100%; text-align: center;">
                  {{ currentDetail.question.answer || '—' }}
                </div>
              </div>
            </div>

            <!-- 原题原图 -->
            <div v-if="currentDetail.question.image_url" class="detail-section">
              <div class="detail-label">
                <Image :size="15" /> 原卷截图
              </div>
              <div style="border-radius: 8px; overflow: hidden; border: 1px solid var(--border);">
                <img
                  :src="getImageUrl(currentDetail.question.image_url)"
                  style="max-width: 100%; display: block;"
                  alt="错题截图"
                />
              </div>
            </div>

            <!-- 知识点标签 -->
            <div
              v-if="currentDetail.question.knowledge_points && currentDetail.question.knowledge_points.length"
              class="detail-section"
            >
              <div class="detail-label">
                <TagIcon :size="15" /> 考察知识点
              </div>
              <div style="display: flex; gap: 6px; flex-wrap: wrap;">
                <span
                  v-for="kp in currentDetail.question.knowledge_points"
                  :key="kp"
                  class="tag-badge tag-kp"
                  style="font-size: var(--font-size-sm); padding: 4px 10px;"
                >{{ kp }}</span>
              </div>
            </div>

            <!-- AI 详解 -->
            <div class="detail-section">
              <div class="detail-label">
                <Sparkles :size="15" style="color: var(--primary);" /> 考点精析与思路拆解
              </div>
              <div class="analysis-box">
                {{ currentDetail.question.analysis || '暂无详细分析，可手动补充或通过 AI 生成。' }}
              </div>
            </div>

            <!-- 重做历史记录 -->
            <div class="detail-section">
              <div class="detail-label" style="justify-content: space-between;">
                <div style="display: flex; align-items: center; gap: 6px;">
                  <History :size="15" /> 重做记录（{{ currentDetail.reviews.length }} 次）
                </div>
                <span style="font-size: var(--font-size-xs); color: var(--muted);">连续正确 2 次自动标为已掌握</span>
              </div>

              <div v-if="currentDetail.reviews.length === 0" style="color: var(--muted); font-size: var(--font-size-sm); padding: 12px 0;">
                暂无重做记录，快在下方记录练习成果吧！
              </div>
              <div v-else>
                <div
                  v-for="rev in currentDetail.reviews"
                  :key="rev.id"
                  class="review-history-item"
                >
                  <div style="display: flex; align-items: center; gap: 8px;">
                    <span :class="['tag-badge', rev.result === '对' ? 'tag-subject' : 'tag-reason']">
                      {{ rev.result === '对' ? '做对 ✔' : '做错 ✘' }}
                    </span>
                    <span v-if="rev.note" style="color: var(--text-secondary);">{{ rev.note }}</span>
                  </div>
                  <span style="color: var(--muted); font-size: var(--font-size-xs);">
                    {{ formatDateTime(rev.reviewed_at) }}
                  </span>
                </div>
              </div>
            </div>

            <!-- 手机答题二维码 -->
            <div class="detail-section" style="background: var(--surface-muted); padding: 14px; border-radius: var(--radius-md);">
              <div style="display: flex; align-items: center; justify-content: space-between;">
                <div>
                  <div style="font-weight: var(--font-weight-medium); font-size: var(--font-size-sm); color: var(--text);">
                    手机扫码免密做题
                  </div>
                  <div style="font-size: var(--font-size-xs); color: var(--muted); margin-top: 2px;">
                    手机微信或系统相机扫码，支持隐藏答案自主练习
                  </div>
                </div>
                <button type="button" class="ghost-button compact-button" @click="showQrModal = true">
                  <QrCode :size="14" /> 查看二维码
                </button>
              </div>
            </div>
          </div>

          <!-- 抽屉底部操作条 -->
          <div class="drawer-footer" style="display: flex; justify-content: space-between; align-items: center;">
            <div style="display: flex; gap: 8px;">
              <button type="button" class="danger-button" @click="handleDelete(currentDetail.question.id)">
                <Trash2 :size="14" /> 删除
              </button>
              <button type="button" class="ghost-button" @click="openEditModal(currentDetail.question)">
                <Edit3 :size="14" /> 编辑
              </button>
            </div>
            <div style="display: flex; gap: 8px;">
              <button type="button" class="success-button" @click="handleRecordReview('对')">
                <Check :size="14" /> 记一次做对
              </button>
              <button type="button" class="danger-button" @click="handleRecordReview('错')">
                <X :size="14" /> 记一次做错
              </button>
            </div>
          </div>
        </aside>
      </div>
    </div>
  </Teleport>

    <!-- 侧边抽屉：新建/编辑错题 -->
    <Teleport to="body">
      <div
        v-if="formVisible"
        class="drawer-overlay"
        @mousedown.self="formVisible = false"
      >
        <div class="drawer-slide" style="max-width: 640px;">
          <aside class="drawer">
            <div class="drawer-header">
              <div class="drawer-header-top">
                <span class="eyebrow">{{ isEditing ? '编辑错题' : '新增录入' }}</span>
                <button type="button" class="drawer-close-btn" aria-label="关闭抽屉" @click="formVisible = false">
                  <X :size="18" />
                </button>
              </div>
              <h2 style="margin: 6px 0 0; font-size: 20px; font-weight: 700; color: var(--text);">
                {{ isEditing ? '编辑错题 #' + formData.id : '录入新错题' }}
              </h2>
            </div>

          <div class="drawer-body">
            <form class="drawer-form" @submit.prevent="submitForm">
              <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-bottom: 14px;">
                <div>
                  <label class="detail-label">科目 *</label>
                  <select v-model="formData.subject" required style="width: 100%;">
                    <option v-for="s in subjects" :key="s" :value="s">{{ s }}</option>
                  </select>
                </div>
                <div>
                  <label class="detail-label">掌握状态</label>
                  <select v-model="formData.status" style="width: 100%;">
                    <option value="待重做">待重做</option>
                    <option value="已重做">已重做</option>
                    <option value="已掌握">已掌握</option>
                  </select>
                </div>
              </div>

              <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-bottom: 14px;">
                <div>
                  <label class="detail-label">来源试卷/练习</label>
                  <input v-model="formData.source" type="text" placeholder="如：期中模拟卷一 第12题" style="width: 100%;" />
                </div>
                <div>
                  <label class="detail-label">错因分类</label>
                  <select v-model="formData.error_reason" style="width: 100%;">
                    <option value="">请选择</option>
                    <option value="概念不清">概念不清</option>
                    <option value="计算失误">计算失误</option>
                    <option value="审题错误">审题错误</option>
                    <option value="方法不会">方法不会</option>
                    <option value="粗心马虎">粗心马虎</option>
                  </select>
                </div>
              </div>

              <div style="margin-bottom: 14px;">
                <label class="detail-label">题干内容 *</label>
                <textarea
                  v-model="formData.stem"
                  required
                  rows="4"
                  placeholder="输入题目完整描述..."
                  style="width: 100%; resize: vertical;"
                ></textarea>
              </div>

              <!-- 选项（可增删） -->
              <div style="margin-bottom: 14px;">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px;">
                  <label class="detail-label" style="margin-bottom: 0;">选择题选项（选填）</label>
                  <button type="button" class="ghost-button compact-button" @click="addOption">
                    + 增加选项
                  </button>
                </div>
                <div v-for="(_, idx) in formData.options" :key="idx" style="display: flex; gap: 8px; margin-bottom: 6px;">
                  <span style="font-weight: 600; line-height: 34px; width: 20px;">{{ 'ABCD'[idx] || idx + 1 }}.</span>
                  <input v-model="formData.options[idx]" type="text" placeholder="选项内容" style="flex: 1;" />
                  <button type="button" class="icon-button icon-danger" @click="removeOption(idx)" title="删除此项">
                    <Trash2 :size="14" />
                  </button>
                </div>
              </div>

              <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-bottom: 14px;">
                <div>
                  <label class="detail-label">正确答案 *</label>
                  <input v-model="formData.answer" required type="text" placeholder="如：C 或 42" style="width: 100%;" />
                </div>
                <div>
                  <label class="detail-label">我的错误答案</label>
                  <input v-model="formData.my_answer" type="text" placeholder="如：A" style="width: 100%;" />
                </div>
              </div>

              <div style="margin-bottom: 14px;">
                <label class="detail-label">知识点标签（用逗号分隔）</label>
                <input v-model="kpInput" type="text" placeholder="如：一元二次方程, 判别式应用" style="width: 100%;" />
              </div>

              <div style="margin-bottom: 14px;">
                <label class="detail-label">考点精析与解题思路</label>
                <textarea
                  v-model="formData.analysis"
                  rows="4"
                  placeholder="输入解题核心逻辑、易错陷阱或 AI 详尽解析..."
                  style="width: 100%; resize: vertical;"
                ></textarea>
              </div>

              <div style="margin-bottom: 14px;">
                <label class="detail-label">原题截图（图片文件）</label>
                <input type="file" accept="image/*" @change="onImageFileChange" />
                <div v-if="formData.image_url" style="margin-top: 8px; font-size: var(--font-size-xs); color: var(--muted);">
                  已关联图片路径: {{ formData.image_url }}
                </div>
              </div>
            </form>
          </div>

          <div class="drawer-footer" style="display: flex; justify-content: flex-end; gap: 10px;">
            <button type="button" class="ghost-button" @click="formVisible = false">取消</button>
            <button type="button" class="primary-button" :disabled="submitting" @click="submitForm">
              {{ submitting ? '保存中...' : '确认保存' }}
            </button>
          </div>
        </aside>
      </div>
    </div>
  </Teleport>

    <!-- 弹窗：二维码展示 -->
    <Teleport to="body">
      <div v-if="showQrModal && currentDetail" class="dialog-overlay" @click.self="showQrModal = false">
        <div class="dialog-box" style="text-align: center; max-width: 360px; padding: 24px;">
          <h3 style="margin-bottom: 6px;">错题 #{{ currentDetail.question.id }} 扫码答题</h3>
          <p style="font-size: var(--font-size-xs); color: var(--muted); margin-bottom: 16px;">
            手机微信或系统相机扫描此二维码，立即在移动端进行无密重做练习。
          </p>
          <div style="background: #fff; padding: 12px; display: inline-block; border-radius: 8px; border: 1px solid var(--border);">
            <img :src="api.getQrCodeUrl(currentDetail.question.id)" style="width: 180px; height: 180px; display: block;" alt="答题二维码" />
          </div>
          <div style="margin-top: 16px;">
            <button type="button" class="primary-button" style="width: 100%;" @click="showQrModal = false">完成</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import {
  BookOpen, AlertCircle, CheckCircle2, TrendingUp, Search, Plus, RefreshCw,
  X, FileText, List, XCircle, Image, Tag as TagIcon, Sparkles, History,
  Trash2, Edit3, Check, QrCode
} from 'lucide-vue-next';
import { api, getImageUrl } from '../api/client';
import type { Question, QuestionDetail, StatsOverview } from '../types';

const subjects = ['数学', '物理', '化学', '英语', '语文', '道法', '历史', '地理', '生物'];

const questions = ref<Question[]>([]);
const stats = ref<Partial<StatsOverview>>({});
const loading = ref(false);

const filters = reactive({
  subject: '',
  status: '',
  kp: '',
  q: ''
});

// 抽屉与弹窗状态
const detailVisible = ref(false);
const currentDetail = ref<QuestionDetail | null>(null);
const formVisible = ref(false);
const isEditing = ref(false);
const submitting = ref(false);
const showQrModal = ref(false);

// 表单数据
const formData = reactive<Partial<Question> & { options: string[] }>({
  subject: '数学',
  source: '',
  stem: '',
  options: [],
  answer: '',
  my_answer: '',
  analysis: '',
  knowledge_points: [],
  difficulty: 3,
  error_reason: '概念不清',
  status: '待重做',
  image_url: ''
});
const kpInput = ref('');
let imageFile: File | null = null;
const loadError = ref('');

async function fetchQuestions() {
  loading.value = true;
  loadError.value = '';
  try {
    const res = await api.listQuestions(filters);
    questions.value = res.questions || [];
  } catch (err: any) {
    console.error(err);
    loadError.value = err.message || '加载错题失败，请检查网络或配置 API Key';
  } finally {
    loading.value = false;
  }
}

async function fetchStats() {
  try {
    const s = await api.getStats();
    stats.value = s;
  } catch (err) {
    console.error(err);
  }
}

function setFilterSubject(sub: string) {
  filters.subject = sub;
  fetchQuestions();
}

function setFilterStatus(st: string) {
  filters.status = st;
  fetchQuestions();
}

function resetFilters() {
  filters.subject = '';
  filters.status = '';
  filters.kp = '';
  filters.q = '';
  fetchQuestions();
}

async function openDetail(id: number) {
  try {
    currentDetail.value = await api.getQuestion(id);
    detailVisible.value = true;
  } catch (err: any) {
    alert(err.message || '加载详情失败');
  }
}

async function handleRecordReview(result: '对' | '错') {
  if (!currentDetail.value) return;
  const qId = currentDetail.value.question.id;
  try {
    await api.addReview(qId, result);
    // 重新拉取详情和列表
    currentDetail.value = await api.getQuestion(qId);
    fetchQuestions();
    fetchStats();
  } catch (err: any) {
    alert(err.message || '记录失败');
  }
}

async function handleDelete(id: number) {
  if (!confirm(`确认彻底删除错题 #${id} 吗？`)) return;
  try {
    await api.deleteQuestion(id);
    detailVisible.value = false;
    currentDetail.value = null;
    fetchQuestions();
    fetchStats();
  } catch (err: any) {
    alert(err.message || '删除失败');
  }
}

function openCreateModal() {
  isEditing.value = false;
  Object.assign(formData, {
    id: undefined,
    subject: filters.subject || '数学',
    source: '',
    stem: '',
    options: [],
    answer: '',
    my_answer: '',
    analysis: '',
    knowledge_points: [],
    difficulty: 3,
    error_reason: '',
    status: '待重做',
    image_url: ''
  });
  kpInput.value = '';
  imageFile = null;
  formVisible.value = true;
}

function openEditModal(q: Question) {
  isEditing.value = true;
  Object.assign(formData, {
    id: q.id,
    subject: q.subject,
    source: q.source,
    stem: q.stem,
    options: [...(q.options || [])],
    answer: q.answer,
    my_answer: q.my_answer,
    analysis: q.analysis,
    knowledge_points: [...(q.knowledge_points || [])],
    difficulty: q.difficulty || 3,
    error_reason: q.error_reason || '',
    status: q.status,
    image_url: q.image_url
  });
  kpInput.value = (q.knowledge_points || []).join(', ');
  imageFile = null;
  detailVisible.value = false;
  formVisible.value = true;
}

function addOption() {
  formData.options.push('');
}

function removeOption(idx: number) {
  formData.options.splice(idx, 1);
}

function onImageFileChange(e: Event) {
  const target = e.target as HTMLInputElement;
  if (target.files && target.files[0]) {
    imageFile = target.files[0];
  }
}

async function submitForm() {
  if (!formData.stem || !formData.stem.trim()) {
    alert('请填写题干内容');
    return;
  }
  if (!formData.answer || !formData.answer.trim()) {
    alert('请填写正确答案');
    return;
  }

  submitting.value = true;
  try {
    // 上传图片（若有）
    if (imageFile) {
      const upRes = await api.uploadImage(imageFile, formData.subject || '错题');
      formData.image_url = upRes.url;
    }

    // 解析知识点
    formData.knowledge_points = kpInput.value
      .split(/[,，、\s]+/)
      .map(k => k.trim())
      .filter(Boolean);

    if (isEditing.value && formData.id) {
      await api.updateQuestion(formData.id, formData);
    } else {
      await api.createQuestion(formData);
    }

    formVisible.value = false;
    fetchQuestions();
    fetchStats();
  } catch (err: any) {
    alert(err.message || '保存失败');
  } finally {
    submitting.value = false;
  }
}

function getStatusClass(status: string) {
  if (status === '待重做') return 'status-pending';
  if (status === '已重做') return 'status-reworked';
  if (status === '已掌握') return 'status-mastered';
  return '';
}

function formatDate(dateStr?: string) {
  if (!dateStr) return '';
  return dateStr.slice(0, 10);
}

function formatDateTime(dateStr?: string) {
  if (!dateStr) return '';
  return dateStr.slice(0, 16).replace('T', ' ');
}

onMounted(() => {
  fetchQuestions();
  fetchStats();
});
</script>

<style scoped>
.spin-anim {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  backdrop-filter: blur(4px);
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.dialog-box {
  background: var(--surface-solid);
  border-radius: var(--radius-panel);
  box-shadow: var(--shadow-popover);
  border: 1px solid var(--border-soft);
}
</style>
