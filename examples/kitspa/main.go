package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/qwenode/omnixkit/kitspa"
)

//go:embed dist/*
var staticFS embed.FS

func main() {
	r := gin.Default()

	// 挂载 SPA 服务
	err := kitspa.Mount(r, kitspa.Config{
		FS:      staticFS,
		DistDir: "dist",

		AssetsPath: "/assets",
		IndexFile:  "index.html",

		// 静态资源（js/css/png/jpg...）建议长缓存；如果你的文件名不带 hash（如 app.js），请把 max-age 调短
		AssetsCacheControl: "public, max-age=31536000, immutable",
		// index.html 建议不缓存（确保发布后能及时更新）
		IndexCacheControl: "no-cache",

		BlockPrefixes: []string{"/.", "/admin."},
	})
	if err != nil {
		log.Fatalf("挂载 SPA 失败: %v", err)
	}

	// 可以添加 API 路由
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
	}

	fmt.Println("启动服务器: http://localhost:8080")
	fmt.Println("测试命令:")
	fmt.Println("  curl http://localhost:8080/")
	fmt.Println("  curl http://localhost:8080/about")
	fmt.Println("  curl http://localhost:8080/assets/app.js")
	fmt.Println("  curl -I http://localhost:8080/assets/app.js")
	fmt.Println("  curl -I http://localhost:8080/")
	fmt.Println("  curl http://localhost:8080/api/ping")

	r.Run(":8080")
}
