import type { Question, QuestionDetail, QuestionFilters, StatsOverview } from '../types';

const STORAGE_KEY_API_KEY = 'mb_api_key';
const STORAGE_KEY_BASE_URL = 'mb_base_url';

// 默认兜底 API Key（支持构建期环境变量 VITE_DEFAULT_API_KEY 注入，未配置则留空由用户在界面输入）
const DEFAULT_API_KEY = (import.meta.env.VITE_DEFAULT_API_KEY as string) || '';

export function getApiKey(): string {
  if (typeof window !== 'undefined') {
    // 1. 优先从 url 查询参数读取 ?k=xxx 并缓存（方便快捷扫码或书签跳转）
    const params = new URLSearchParams(window.location.search);
    const k = params.get('k');
    if (k) {
      localStorage.setItem(STORAGE_KEY_API_KEY, k);
      params.delete('k');
      const newSearch = params.toString() ? `?${params.toString()}` : '';
      window.history.replaceState(null, '', window.location.pathname + newSearch);
      return k;
    }

    // 2. 本地缓存已保存的 key
    const cached = localStorage.getItem(STORAGE_KEY_API_KEY);
    if (cached) return cached;

    // 3. 构建时环境变量注入的兜底 key
    if (DEFAULT_API_KEY) {
      localStorage.setItem(STORAGE_KEY_API_KEY, DEFAULT_API_KEY);
      return DEFAULT_API_KEY;
    }
  }
  return '';
}

export function setApiKey(key: string): void {
  if (typeof window !== 'undefined') {
    localStorage.setItem(STORAGE_KEY_API_KEY, key.trim());
    window.dispatchEvent(new CustomEvent('mb-api-key-changed', { detail: key.trim() }));
  }
}

export function getBaseUrl(): string {
  if (typeof window !== 'undefined') {
    // 1. 用户自定义配置的 BaseUrl (如跨域调试)
    const custom = localStorage.getItem(STORAGE_KEY_BASE_URL);
    if (custom) return custom;

    // 2. Vite 构建注入的 BASE_URL (如生产环境 /cuoti/)
    const viteBase = import.meta.env.BASE_URL;
    if (viteBase && viteBase !== '/' && viteBase !== './') {
      return viteBase.replace(/\/+$/, '');
    }

    // 3. 运行时动态感知当前页面子路径 (如 /cuoti/...)
    const pathParts = window.location.pathname.split('/').filter(Boolean);
    if (pathParts.length > 0 && pathParts[0] === 'cuoti') {
      return '/cuoti';
    }
  }
  return '';
}

export function setBaseUrl(url: string): void {
  if (typeof window !== 'undefined') {
    const cleanUrl = url.trim().replace(/\/+$/, '');
    localStorage.setItem(STORAGE_KEY_BASE_URL, cleanUrl);
    window.dispatchEvent(new CustomEvent('mb-base-url-changed', { detail: cleanUrl }));
  }
}

export function resolveUrl(path: string): string {
  const base = getBaseUrl();
  const cleanPath = path.startsWith('/') ? path : `/${path}`;
  if (!base) return cleanPath;
  if (cleanPath.startsWith(base)) return cleanPath;
  return `${base}${cleanPath}`;
}

export function getImageUrl(imagePath?: string): string {
  if (!imagePath) return '';
  if (imagePath.startsWith('http://') || imagePath.startsWith('https://')) {
    return imagePath;
  }
  return resolveUrl(`/${imagePath.replace(/^\/+/, '')}`);
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const apiKey = getApiKey();
  const headers = new Headers(options.headers || {});
  
  if (apiKey) {
    headers.set('X-API-Key', apiKey);
  }
  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }

  const fullUrl = resolveUrl(path);
  const response = await fetch(fullUrl, {
    ...options,
    headers
  });

  if (response.status === 401) {
    if (typeof window !== 'undefined') {
      window.dispatchEvent(new CustomEvent('mb-auth-error'));
    }
    throw new Error('需要有效的 API Key，请在系统设置中配置');
  }

  if (!response.ok) {
    const errorText = await response.text();
    let message = `请求失败: HTTP ${response.status}`;
    try {
      const json = JSON.parse(errorText);
      if (json.error) message = json.error;
    } catch {
      if (errorText) message = errorText;
    }
    throw new Error(message);
  }

  return response.json();
}

export const api = {
  async listQuestions(filters?: Partial<QuestionFilters>): Promise<{ questions: Question[]; count: number }> {
    const query = new URLSearchParams();
    if (filters) {
      if (filters.subject) query.set('subject', filters.subject);
      if (filters.status) query.set('status', filters.status);
      if (filters.kp) query.set('kp', filters.kp);
      if (filters.q) query.set('q', filters.q);
    }
    const qStr = query.toString() ? `?${query.toString()}` : '';
    return request<{ questions: Question[]; count: number }>(`/api/v1/questions${qStr}`);
  },

  async getQuestion(id: number): Promise<QuestionDetail> {
    return request<QuestionDetail>(`/api/v1/questions/${id}`);
  },

  async createQuestion(data: Partial<Question>): Promise<Question> {
    return request<Question>('/api/v1/questions', {
      method: 'POST',
      body: JSON.stringify(data)
    });
  },

  async updateQuestion(id: number, data: Partial<Question>): Promise<Question> {
    return request<Question>(`/api/v1/questions/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    });
  },

  async deleteQuestion(id: number): Promise<{ message: string }> {
    return request<{ message: string }>(`/api/v1/questions/${id}`, {
      method: 'DELETE'
    });
  },

  async addReview(id: number, result: '对' | '错', myAnswer?: string, note?: string): Promise<{ message: string; question: Question }> {
    return request<{ message: string; question: Question }>(`/api/v1/questions/${id}/reviews`, {
      method: 'POST',
      body: JSON.stringify({ result, my_answer: myAnswer, note })
    });
  },

  async getStats(): Promise<StatsOverview> {
    return request<StatsOverview>('/api/v1/stats');
  },

  async uploadImage(file: File, folder = '未分类', name?: string): Promise<{ url: string; folder: string; file: string }> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('folder', folder);
    if (name) formData.append('name', name);

    return request<{ url: string; folder: string; file: string }>('/api/v1/images', {
      method: 'POST',
      body: formData
    });
  },

  getQrCodeUrl(id: number): string {
    const key = getApiKey();
    return resolveUrl(`/api/v1/qrcode/${id}${key ? `?k=${key}` : ''}`);
  },

  // 移动端免鉴权接口
  async getMobileQuestion(id: number): Promise<{ question: Question }> {
    return request<{ question: Question }>(`/api/v1/m/${id}`);
  },

  async submitMobileReview(id: number, result: '对' | '错', myAnswer?: string): Promise<{ message: string; question: Question }> {
    return request<{ message: string; question: Question }>(`/api/v1/m/${id}/review`, {
      method: 'POST',
      body: JSON.stringify({ result, my_answer: myAnswer })
    });
  }
};
