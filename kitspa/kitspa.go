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
	// AssetsCacheControl 静态资源（AssetsPath）Cache-Control，默认 "public, max-age=31536000, immutable"
	AssetsCacheControl string
	// IndexCacheControl 首页/SPA 回退（index.html）Cache-Control，默认 "no-cache"
	IndexCacheControl string
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
	if cfg.AssetsCacheControl == "" {
		cfg.AssetsCacheControl = "public, max-age=31536000, immutable"
	}
	if cfg.IndexCacheControl == "" {
		cfg.IndexCacheControl = "no-cache"
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
	assetsHandler := func(c *gin.Context) {
		if cfg.AssetsCacheControl != "" {
			c.Writer.Header().Set("Cache-Control", cfg.AssetsCacheControl)
		}
		staticServer.ServeHTTP(c.Writer, c.Request)
	}
	r.GET(cfg.AssetsPath+"/*filepath", assetsHandler)
	r.HEAD(cfg.AssetsPath+"/*filepath", assetsHandler)

	// 首页路由
	indexHandler := func(c *gin.Context) {
		file, err := cfg.FS.ReadFile(indexPath)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if cfg.IndexCacheControl != "" {
			c.Writer.Header().Set("Cache-Control", cfg.IndexCacheControl)
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", file)
	}
	r.GET("/", indexHandler)
	r.HEAD("/", indexHandler)

	// SPA 路由回退
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusNotFound)
			return
		}

		path := c.Request.URL.Path

		// 检查是否需要阻止
		for _, prefix := range cfg.BlockPrefixes {
			if strings.HasPrefix(path, prefix) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		}

		indexHandler(c)
	})

	return nil
}
