package kitrouter

import (
    "fmt"
    "io/fs"
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
)

// Config 静态文件服务配置
type Config struct {
    EmbedFS        fs.FS         // embed.FS 或其子目录
    IndexFile      string        // 默认 "index.html"
    AssetsPath     string        // 静态资源路径前缀，默认 "/assets"
    IgnorePaths    []string      // 需要忽略的路径前缀，如 "/.", "/admin."
    AssetsCacheAge time.Duration // 静态资源缓存时间，默认 1 年
}

// 挂载静态文件服务到 Gin 路由
func ServeSPA(r *gin.Engine, cfg Config) {
    if cfg.IndexFile == "" {
        cfg.IndexFile = "index.html"
    }
    if cfg.AssetsPath == "" {
        cfg.AssetsPath = "/assets"
    }
    if cfg.AssetsCacheAge == 0 {
        cfg.AssetsCacheAge = 365 * 24 * time.Hour // 默认 1 年
    }

    staticServer := http.FileServer(http.FS(cfg.EmbedFS))
    cacheControl := fmt.Sprintf("public, max-age=%d", int(cfg.AssetsCacheAge.Seconds()))

    // 静态资源
    r.GET(cfg.AssetsPath+"/*filepath", func(c *gin.Context) {
        c.Header("Cache-Control", cacheControl)
        staticServer.ServeHTTP(c.Writer, c.Request)
    })

    // 首页
    r.GET("/", func(c *gin.Context) {
        serveIndex(c, cfg.EmbedFS, cfg.IndexFile)
    })

    // SPA fallback
    r.NoRoute(func(c *gin.Context) {
        path := c.Request.URL.Path
        for _, ignore := range cfg.IgnorePaths {
            if strings.HasPrefix(path, ignore) {
                c.Abort()
                return
            }
        }
        serveIndex(c, cfg.EmbedFS, cfg.IndexFile)
    })
}

func serveIndex(c *gin.Context, embedFS fs.FS, indexFile string) {
    file, err := fs.ReadFile(embedFS, indexFile)
    if err != nil {
        c.AbortWithStatus(http.StatusNotFound)
        return
    }
    c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
    c.Writer.WriteHeader(http.StatusOK)
    c.Writer.Write(file)
}

