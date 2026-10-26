export interface Question {
  id: number;
  subject: string;
  source: string;
  exam_date?: string;
  stem: string;
  options?: string[];
  answer: string;
  my_answer?: string;
  analysis?: string;
  knowledge_points?: string[];
  difficulty?: number;
  error_reason?: string;
  status: '待重做' | '已重做' | '已掌握' | string;
  image_url?: string;
  ima_media_id?: string;
  created_at: string;
  updated_at?: string;
}

export interface Review {
  id: number;
  question_id: number;
  result: '对' | '错' | string;
  my_answer?: string;
  note?: string;
  reviewed_at: string;
}

export interface QuestionDetail {
  question: Question;
  reviews: Review[];
}

export interface SubjectStat {
  name: string;
  count: number;
}

export interface ReasonStat {
  name: string;
  count: number;
}

export interface KpStat {
  kp: string;
  count: number;
}

export interface StatsOverview {
  total: number;
  pending: number;
  mastered: number;
  recent_pass_rate: number;
  by_subject: SubjectStat[];
  by_reason: ReasonStat[];
  kp_stats: KpStat[];
}

export interface QuestionFilters {
  subject: string;
  status: string;
  kp: string;
  q: string;
}
