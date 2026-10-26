package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"mistakebook/internal/model"

	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	APIKey string
}

func New(db *gorm.DB) *Handler {
	key := os.Getenv("API_KEY")
	if key == "" {
		key = "changeme-please-set-API_KEY"
	}
	return &Handler{DB: db, APIKey: key}
}

// auth API_KEY 鉴权（单用户）
func (h *Handler) auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-Key") != h.APIKey && c.Query("k") != h.APIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的 API Key"})
			return
		}
		c.Next()
	}
}

func (h *Handler) Register(r *gin.Engine, basePath string) {
	api := r.Group(basePath + "/api/v1")
	api.Use(h.auth())
	{
		api.POST("/questions", h.CreateQuestions)
		api.GET("/questions", h.ListQuestions)
		api.GET("/questions/:id", h.GetQuestion)
		api.PUT("/questions/:id", h.UpdateQuestion)
		api.DELETE("/questions/:id", h.DeleteQuestion)
		api.POST("/questions/:id/reviews", h.AddReview)
		api.GET("/stats", h.Stats)
		api.GET("/qrcode/:id", h.QRCode)
		api.POST("/images", h.UploadImage)
	}

	// 移动端扫码（免 API_KEY，只读题目 + 提交重做结果）
	m := r.Group(basePath + "/api/v1/m")
	{
		m.GET("/:id", h.MobileGet)
		m.POST("/:id/review", h.MobileReview)
	}
}

// ---------------- 错题 CRUD ----------------

type questionInput struct {
	Subject         string   `json:"subject"`
	Source          string   `json:"source"`
	ExamDate        string   `json:"exam_date"`
	Stem            string   `json:"stem"`
	Options         []string `json:"options"`
	Answer          string   `json:"answer"`
	MyAnswer        string   `json:"my_answer"`
	Analysis        string   `json:"analysis"`
	KnowledgePoints []string `json:"knowledge_points"`
	Difficulty      int      `json:"difficulty"`
	ErrorReason     string   `json:"error_reason"`
	Status          string   `json:"status"`
	ImageURL        string   `json:"image_url"`
	ImaMediaID      string   `json:"ima_media_id"`
}

// CreateQuestions 批量入库（AI Agent 智能调用），支持单个对象或数组
func (h *Handler) CreateQuestions(c *gin.Context) {
	body, _ := c.GetRawData()
	var list []questionInput
	if err := json.Unmarshal(body, &list); err != nil {
		var single questionInput
		if err := json.Unmarshal(body, &single); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON 格式错误: " + err.Error()})
			return
		}
		list = []questionInput{single}
	}
	var created []model.Question
	for _, in := range list {
		if in.Stem == "" {
			continue
		}
		if in.Status == "" {
			in.Status = "待重做"
		}
		q := model.Question{
			Subject: in.Subject, Source: in.Source, ExamDate: in.ExamDate,
			Stem: in.Stem, Options: model.Options(in.Options), Answer: in.Answer,
			MyAnswer: in.MyAnswer, Analysis: in.Analysis,
			KnowledgePoints: model.KPs(in.KnowledgePoints),
			Difficulty:      in.Difficulty, ErrorReason: in.ErrorReason,
			Status: in.Status, ImageURL: in.ImageURL, ImaMediaID: in.ImaMediaID,
		}
		if err := h.DB.Create(&q).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		created = append(created, q)
	}
	c.JSON(http.StatusOK, gin.H{"created": len(created), "questions": created})
}

func (h *Handler) ListQuestions(c *gin.Context) {
	q := h.DB.Order("created_at DESC")
	if s := c.Query("subject"); s != "" {
		q = q.Where("subject = ?", s)
	}
	if s := c.Query("status"); s != "" {
		q = q.Where("status = ?", s)
	}
	if kp := c.Query("kp"); kp != "" {
		q = q.Where("knowledge_points LIKE ?", "%"+kp+"%")
	}
	if kw := c.Query("q"); kw != "" {
		kw = "%" + kw + "%"
		q = q.Where("stem LIKE ? OR answer LIKE ? OR analysis LIKE ? OR source LIKE ?", kw, kw, kw, kw)
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	var total int64
	q.Model(&model.Question{}).Count(&total)
	var qs []model.Question
	q.Limit(limit).Offset(offset).Find(&qs)
	c.JSON(http.StatusOK, gin.H{"total": total, "questions": qs})
}

func (h *Handler) GetQuestion(c *gin.Context) {
	var q model.Question
	if err := h.DB.First(&q, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该错题"})
		return
	}
	var reviews []model.Review
	h.DB.Where("question_id = ?", q.ID).Order("reviewed_at DESC").Find(&reviews)
	c.JSON(http.StatusOK, gin.H{"question": q, "reviews": reviews})
}

func (h *Handler) UpdateQuestion(c *gin.Context) {
	var q model.Question
	if err := h.DB.First(&q, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该错题"})
		return
	}
	var in map[string]interface{}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	delete(in, "id")
	for _, field := range []string{"knowledge_points", "options"} {
		if v, ok := in[field]; ok {
			if raw, err := json.Marshal(v); err == nil {
				in[field] = string(raw)
			}
		}
	}
	h.DB.Model(&q).Updates(in)
	c.JSON(http.StatusOK, gin.H{"question": q})
}

func (h *Handler) DeleteQuestion(c *gin.Context) {
	id := c.Param("id")
	h.DB.Where("question_id = ?", id).Delete(&model.Review{})
	h.DB.Delete(&model.Question{}, id)
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// ---------------- 重做记录 ----------------

type reviewInput struct {
	Result   string `json:"result"` // 对/错
	MyAnswer string `json:"my_answer"`
	Note     string `json:"note"`
}

func (h *Handler) applyReview(questionID uint, in reviewInput) (model.Question, error) {
	var q model.Question
	if err := h.DB.First(&q, questionID).Error; err != nil {
		return q, err
	}
	h.DB.Create(&model.Review{
		QuestionID: questionID, Result: in.Result, MyAnswer: in.MyAnswer,
		Note: in.Note, ReviewedAt: time.Now(),
	})
	// 最近连续正确次数
	var reviews []model.Review
	h.DB.Where("question_id = ?", questionID).Order("reviewed_at DESC").Limit(model.MasteryStreak).Find(&reviews)
	streak := 0
	for _, r := range reviews {
		if r.Result == "对" {
			streak++
		} else {
			break
		}
	}
	switch {
	case streak >= model.MasteryStreak:
		h.DB.Model(&q).Update("status", "已掌握")
	case in.Result == "错":
		h.DB.Model(&q).Update("status", "待重做")
	default:
		h.DB.Model(&q).Update("status", "已重做")
	}
	h.DB.First(&q, questionID)
	return q, nil
}

func (h *Handler) AddReview(c *gin.Context) {
	h.doReview(c, false)
}

// MobileGet 扫码端只读（免鉴权，不含我的答案与重做历史，防止扫自己卷子时被剧透）
// MobileGet 扫码端只读（免鉴权，不含我的答案与重做历史，防止剧透）
func (h *Handler) MobileGet(c *gin.Context) {
	var q model.Question
	if err := h.DB.First(&q, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该错题"})
		return
	}
	q.MyAnswer = ""
	c.JSON(http.StatusOK, gin.H{"question": q})
}

func (h *Handler) MobileReview(c *gin.Context) {
	h.doReview(c, true)
}

func (h *Handler) doReview(c *gin.Context, mobile bool) {
	id, _ := strconv.Atoi(c.Param("id"))
	var in reviewInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if in.Result != "对" && in.Result != "错" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "result 必须为 对 或 错"})
		return
	}
	q, err := h.applyReview(uint(id), in)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该错题"})
		return
	}
	if mobile {
		c.JSON(http.StatusOK, gin.H{"message": "已记录", "status": q.Status})
		return
	}
	c.JSON(http.StatusOK, gin.H{"question": q})
}

// ---------------- 统计 ----------------

func (h *Handler) Stats(c *gin.Context) {
	var total, pending, mastered int64
	h.DB.Model(&model.Question{}).Count(&total)
	h.DB.Model(&model.Question{}).Where("status = ?", "待重做").Count(&pending)
	h.DB.Model(&model.Question{}).Where("status = ?", "已掌握").Count(&mastered)

	type nameCount struct {
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}
	var bySubject []nameCount
	h.DB.Model(&model.Question{}).Select("subject as name, count(*) as count").Group("subject").Scan(&bySubject)
	var byReason []nameCount
	h.DB.Model(&model.Question{}).Select("error_reason as name, count(*) as count").Group("error_reason").Scan(&byReason)

	// 知识点易错排行：在 Go 侧展开知识点
	var qs []model.Question
	h.DB.Find(&qs)
	type kpStat struct {
		KP      string `json:"kp"`
		Count   int64  `json:"count"`
		Pending int64  `json:"pending"`
	}
	kpMap := map[string]*kpStat{}
	for _, q := range qs {
		for _, kp := range q.KnowledgePoints {
			kp = strings.TrimSpace(kp)
			if kp == "" {
				continue
			}
			s, ok := kpMap[kp]
			if !ok {
				s = &kpStat{KP: kp}
				kpMap[kp] = s
			}
			s.Count++
			if q.Status != "已掌握" {
				s.Pending++
			}
		}
	}
	kpStats := make([]kpStat, 0, len(kpMap))
	for _, s := range kpMap {
		kpStats = append(kpStats, *s)
	}
	for i := 0; i < len(kpStats); i++ {
		for j := i + 1; j < len(kpStats); j++ {
			if kpStats[j].Count > kpStats[i].Count {
				kpStats[i], kpStats[j] = kpStats[j], kpStats[i]
			}
		}
	}

	// 近 30 天重做通过率
	var recentReviews, recentRight int64
	since := time.Now().AddDate(0, 0, -30)
	h.DB.Model(&model.Review{}).Where("reviewed_at >= ?", since).Count(&recentReviews)
	h.DB.Model(&model.Review{}).Where("reviewed_at >= ? AND result = ?", since, "对").Count(&recentRight)
	passRate := 0.0
	if recentReviews > 0 {
		passRate = float64(recentRight) / float64(recentReviews) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total, "pending": pending, "mastered": mastered,
		"by_subject": bySubject, "by_reason": byReason,
		"kp_stats":        kpStats,
		"recent_reviews":  recentReviews,
		"recent_pass_rate": passRate,
	})
}

// ---------------- 图片上传 ----------------

var allowedImgExt = map[string]string{
	".png": "image/png", ".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".webp": "image/webp",
}

// sanitizeSeg 清理目录/文件名片段：去除路径分隔符与不可见字符，避免目录穿越
func sanitizeSeg(s string) string {
	replacer := strings.NewReplacer("/", "-", "\\", "-", "\x00", "", "\n", "", "\r", "", "\t", "")
	s = strings.TrimSpace(replacer.Replace(s))
	if len(s) > 120 {
		s = s[:120]
	}
	return s
}

// UploadImage multipart 图片上传。
// form 字段：file（必填）、folder（学生-学科-学期，如"ice7-英语-八上"）、name（文件名主体，如"20261025-英语练习四"）。
// 同名冲突自动追加 -2、-3… 后缀。返回相对路径 images/{folder}/{file}。
func (h *Handler) UploadImage(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 file 字段（multipart/form-data）"})
		return
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if _, ok := allowedImgExt[ext]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持 png/jpg/jpeg/webp"})
		return
	}
	if fh.Size > 10<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片不能超过 10MB"})
		return
	}
	folder := sanitizeSeg(c.PostForm("folder"))
	name := sanitizeSeg(c.PostForm("name"))
	if folder == "" {
		folder = "未分类"
	}
	if name == "" {
		name = strings.TrimSuffix(fh.Filename, filepath.Ext(fh.Filename))
	}

	dir := path.Join("data", "images", folder)
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	final := name + ext
	if _, err := os.Stat(path.Join(dir, final)); err == nil {
		for i := 2; ; i++ {
			cand := fmt.Sprintf("%s-%d%s", name, i, ext)
			if _, err := os.Stat(path.Join(dir, cand)); err != nil {
				final = cand
				break
			}
		}
	}
	dst, err := os.Create(path.Join(dir, final))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer dst.Close()
	src, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()
	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"url":    path.Join("images", folder, final),
		"folder": folder,
		"file":   final,
	})
}

// ---------------- 二维码 ----------------

func (h *Handler) QRCode(c *gin.Context) {
	base := os.Getenv("PUBLIC_BASE_URL")
	if base == "" {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		base = scheme + "://" + c.Request.Host
	}
	url := fmt.Sprintf("%s/m/%s", strings.TrimRight(base, "/"), c.Param("id"))
	png, err := qrcode.Encode(url, qrcode.Medium, 200)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "image/png", png)
}
