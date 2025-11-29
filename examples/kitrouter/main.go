package main

import (
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/gin-gonic/gin"
	"github.com/qwenode/omnixkit/kitrouter"
)

func main() {
	// ============================================
	// kitrouter 使用示例
	// 这是一个用于 Gin + Connect RPC 的路由适配器
	// 支持区分需登录/无需登录路由，统一管理拦截器
	// ============================================

	// 1. 初始化路由适配器
	kitrouter.Bootstrap(
		// 设置无需登录的中间件
		kitrouter.WithGuestMiddlewares(
			loggerMiddleware("guest"),
		),
		// 设置需登录的中间件
		kitrouter.WithAuthMiddlewares(
			loggerMiddleware("auth"),
			mockAuthMiddleware(), // 模拟认证中间件
		),
		// 设置 Connect 拦截器（可选，默认已包含 kitcodec.WithProtoJSON()）
		// kitrouter.WithInterceptors(yourInterceptor),
	)

	// 2. 注册无需登录的路由
	kitrouter.Guest().
		// 注册 Connect RPC 服务
		POST(mockHealthService).
		// 注册自定义路由
		Custom(func(r *gin.RouterGroup) {
			r.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "pong"})
			})
		})

	// 3. 注册需登录的路由
	kitrouter.Auth().
		// 注册 Connect RPC 服务
		POST(mockUserService).
		// 注册自定义路由
		Custom(func(r *gin.RouterGroup) {
			r.GET("/api/profile", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"user_id":  1001,
					"username": "demo_user",
				})
			})
		})

	// 4. 创建 Gin 引擎并挂载路由
	engine := gin.Default()
	kitrouter.Mount(engine)

	// 5. 获取拦截器配置（用于其他需要的地方）
	interceptors := kitrouter.Interceptors()
	fmt.Printf("已配置的拦截器数量: %d\n", len(interceptors))

	fmt.Println("\n=== 路由结构 ===")
	fmt.Println("无需登录路由 (Guest):")
	fmt.Println("  POST /health.v1.HealthService/*any  - Connect RPC 健康检查服务")
	fmt.Println("  GET  /ping                          - 自定义 ping 端点")
	fmt.Println()
	fmt.Println("需登录路由 (Auth):")
	fmt.Println("  POST /user.v1.UserService/*any      - Connect RPC 用户服务")
	fmt.Println("  GET  /api/profile                   - 自定义用户信息端点")

	fmt.Println("\n=== 测试命令 ===")
	fmt.Println("  curl http://localhost:8080/ping")
	fmt.Println("  curl http://localhost:8080/api/profile")
	fmt.Println("  curl -X POST http://localhost:8080/health.v1.HealthService/Check")

	fmt.Println("\n启动服务器: http://localhost:8080")
	engine.Run(":8080")
}

// ============================================
// 以下是模拟的服务和中间件，实际使用时替换为真实实现
// ============================================

// mockHealthService 模拟健康检查服务
// 实际使用时替换为: healthv1connect.NewHealthServiceHandler
func mockHealthService(interceptors []connect.HandlerOption) (string, http.Handler) {
	return "/health.v1.HealthService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	})
}

// mockUserService 模拟用户服务
// 实际使用时替换为: userv1connect.NewUserServiceHandler
func mockUserService(interceptors []connect.HandlerOption) (string, http.Handler) {
	return "/user.v1.UserService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user_id": 1001, "username": "demo_user"}`))
	})
}

// loggerMiddleware 日志中间件
func loggerMiddleware(tag string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("[%s] %s %s\n", tag, c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

// mockAuthMiddleware 模拟认证中间件
// 实际使用时替换为: kitctx.GinMiddlewareJwtAuth[*YourClaims]()
func mockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 模拟认证检查
		// 实际使用时会验证 JWT token
		c.Next()
	}
}
