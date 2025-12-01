package kitspa

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Config SPA 服务配置
type Config struct {
	// FS 嵌入的文件系统
	FS embed.FS
	// DistDir 构建产物目录，默认 "dist"
	DistDir string
	// AssetsPath 静态资源路径前缀，默认 "/assets"
	AssetsPath string
	// IndexFile 入口文件，默认 "index.html"
	IndexFile string
	// BlockPrefixes 需要阻止的路径前缀，默认 ["/.", "/admin."]
	BlockPrefixes []string
}

// Mount 挂载 SPA 服务到 Gin 路由
func Mount(r *gin.Engine, cfg Config) error {
	// 设置默认值
	if cfg.DistDir == "" {
		cfg.DistDir = "dist"
	}
	if cfg.AssetsPath == "" {
		cfg.AssetsPath = "/assets"
	}
	if cfg.IndexFile == "" {
		cfg.IndexFile = "index.html"
	}
	if cfg.BlockPrefixes == nil {
		cfg.BlockPrefixes = []string{"/.", "/admin."}
	}

	// 创建子文件系统
	sub, err := fs.Sub(cfg.FS, cfg.DistDir)
	if err != nil {
		return err
	}

	staticServer := http.FileServer(http.FS(sub))
	indexPath := cfg.DistDir + "/" + cfg.IndexFile

	// 静态资源路由
	r.GET(cfg.AssetsPath+"/*filepath", func(c *gin.Context) {
		staticServer.ServeHTTP(c.Writer, c.Request)
	})

	// 首页路由
	r.GET("/", func(c *gin.Context) {
		file, err := cfg.FS.ReadFile(indexPath)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Writer.Write(file)
		c.Writer.Flush()
	})

	// SPA 路由回退
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 检查是否需要阻止
		for _, prefix := range cfg.BlockPrefixes {
			if strings.HasPrefix(path, prefix) {
				c.Abort()
				return
			}
		}

		// 返回 index.html
		file, err := cfg.FS.ReadFile(indexPath)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Writer.Write(file)
		c.Writer.Flush()
	})

	return nil
}
