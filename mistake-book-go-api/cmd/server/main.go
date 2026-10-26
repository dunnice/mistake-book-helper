package main

import (
	"log"
	"os"
	"path"
	"strings"

	"mistakebook/internal/handler"
	"mistakebook/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openDB() *gorm.DB {
	logCfg := gorm.Config{Logger: logger.Default.LogMode(logger.Warn)}
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// 本地/演示模式：SQLite
		_ = os.MkdirAll("data", 0755)
		db, err := gorm.Open(sqlite.Open("data/mistakebook.db"), &logCfg)
		if err != nil {
			log.Fatalf("打开 SQLite 失败: %v", err)
		}
		log.Println("使用 SQLite: data/mistakebook.db")
		return db
	}
	db, err := gorm.Open(mysql.Open(dsn), &logCfg)
	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}
	log.Println("已连接 MySQL")
	return db
}

func main() {
	db := openDB()
	if err := db.AutoMigrate(&model.Question{}, &model.Review{}); err != nil {
		log.Fatalf("建表失败: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// BASE_PATH 支持子路径挂载（如 nginx 反代 /cuoti/ -> 127.0.0.1:8090，不剥离前缀）
	basePath := strings.TrimRight(os.Getenv("BASE_PATH"), "/")

	// 允许跨域（前端本地开发 Vite 端口 5173 或其他来源）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 查找前端静态资源目录优先级
	staticDir := os.Getenv("STATIC_DIR")
	candidates := []string{
		staticDir,
		"./mistake-book-view/dist",
		"../mistake-book-view/dist",
		"./dist",
		"../dist",
	}
	foundStaticDir := ""
	for _, dir := range candidates {
		if dir != "" {
			if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
				foundStaticDir = dir
				break
			}
		}
	}

	// 查找图片资源目录
	imagesDir := "data/images"
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		if _, err2 := os.Stat("../data/images"); err2 == nil {
			imagesDir = "../data/images"
		}
	}
	_ = os.MkdirAll(imagesDir, 0755)
	r.Static(basePath+"/images", imagesDir)

	if foundStaticDir != "" {
		log.Printf("前端静态资源目录: %s", foundStaticDir)
		assetsDir := foundStaticDir + "/assets"
		if fi, err := os.Stat(assetsDir); err == nil && fi.IsDir() {
			r.Static(basePath+"/assets", assetsDir)
		}
		r.Static(basePath+"/static", foundStaticDir)
		r.StaticFile(basePath+"/logo.png", path.Join(foundStaticDir, "logo.png"))
		r.StaticFile(basePath+"/favicon.png", path.Join(foundStaticDir, "favicon.png"))
		r.StaticFile(basePath+"/apple-touch-icon.png", path.Join(foundStaticDir, "apple-touch-icon.png"))
		r.StaticFile(basePath+"/mistake-book.png", path.Join(foundStaticDir, "mistake-book.png"))
		r.GET(basePath+"/", func(c *gin.Context) { c.File(foundStaticDir + "/index.html") })
		if basePath != "" {
			r.GET(basePath, func(c *gin.Context) { c.Redirect(301, basePath+"/") })
			r.StaticFile("/favicon.ico", path.Join(foundStaticDir, "favicon.png"))
		}
		r.GET(basePath+"/print", func(c *gin.Context) {
			if _, err := os.Stat(foundStaticDir + "/print.html"); err == nil {
				c.File(foundStaticDir + "/print.html")
			} else {
				c.File(foundStaticDir + "/index.html")
			}
		})
		r.GET(basePath+"/m/:id", func(c *gin.Context) {
			if _, err := os.Stat(foundStaticDir + "/mobile.html"); err == nil {
				c.File(foundStaticDir + "/mobile.html")
			} else {
				c.File(foundStaticDir + "/index.html")
			}
		})

		// SPA Fallback: 非 API 请求回退到 index.html，支持 Vue 路由的 History 模式
		r.NoRoute(func(c *gin.Context) {
			p := c.Request.URL.Path
			if strings.HasPrefix(p, basePath+"/api/") || strings.HasPrefix(p, basePath+"/images/") {
				c.JSON(404, gin.H{"error": "接口未找到"})
				return
			}
			c.File(foundStaticDir + "/index.html")
		})
	}

	// 健康检查
	r.GET(basePath+"/healthz", func(c *gin.Context) { c.String(200, "ok") })

	h := handler.New(db)
	h.Register(r, basePath)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("错题本服务已启动: http://0.0.0.0:%s (API_KEY=%s...)", port, strings.Split(os.Getenv("API_KEY"), "")[0]+"***")
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
