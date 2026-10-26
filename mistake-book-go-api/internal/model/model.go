package model

import "time"

// 知识点数组，JSON 存储
type KPs []string

// 选项数组，JSON 存储
type Options []string

// Question 错题主表
type Question struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Subject         string    `gorm:"index;size:20" json:"subject"`            // 数学/物理/化学/英语/语文/道法/历史/地理/生物
	Source          string    `gorm:"size:200" json:"source"`                  // 来源（卷子名 + 日期）
	ExamDate        string    `gorm:"size:20" json:"exam_date"`                // 考试日期 YYYY-MM-DD
	Stem            string    `gorm:"type:text" json:"stem"`                   // 题干
	Options         Options   `gorm:"type:text;serializer:json" json:"options"`                // 选项（选择题）
	Answer          string    `gorm:"type:text" json:"answer"`                 // 正确答案
	MyAnswer        string    `gorm:"type:text" json:"my_answer"`              // 我的错误答案
	Analysis        string    `gorm:"type:text" json:"analysis"`               // AI 详解
	KnowledgePoints KPs       `gorm:"type:text;serializer:json" json:"knowledge_points"`       // 知识点标签
	Difficulty      int       `json:"difficulty"`                              // 1-5
	ErrorReason     string    `gorm:"size:50" json:"error_reason"`             // 概念不清/计算失误/审题错误/方法不会/粗心
	Status          string    `gorm:"index;size:20;default:待重做" json:"status"` // 待重做/已重做/已掌握
	ImageURL        string    `gorm:"size:300" json:"image_url"`               // 原题截图相对路径 images/{folder}/{file}
	ImaMediaID      string    `gorm:"size:200" json:"ima_media_id"`            // ima 知识库 media id（双写关联）
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// MasteryStreak 连续正确次数达到该值则标记已掌握
const MasteryStreak = 2
type Review struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	QuestionID uint      `gorm:"index" json:"question_id"`
	Result     string    `gorm:"size:10" json:"result"` // 对/错
	MyAnswer   string    `gorm:"type:text" json:"my_answer"`
	Note       string    `gorm:"size:300" json:"note"`
	ReviewedAt time.Time `json:"reviewed_at"`
}
