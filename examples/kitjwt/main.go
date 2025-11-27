package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/qwenode/omnixkit/kitjwt"
)

// JwtAdminClaims 自定义 claims 结构体
type JwtAdminClaims struct {
	jwt.RegisteredClaims
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func main() {
	// 1. 生成密钥（仅用于演示，生产环境应使用固定密钥）
	keyHex, err := kitjwt.GenerateKey()
	if err != nil {
		log.Fatalf("生成密钥失败: %v", err)
	}
	fmt.Printf("生成的密钥: %s\n", keyHex)

	// 2. Bootstrap 初始化（指定 claims 类型）
	err = kitjwt.Bootstrap(keyHex, func() *JwtAdminClaims {
		return &JwtAdminClaims{}
	})
	if err != nil {
		log.Fatalf("初始化 JWT 失败: %v", err)
	}

	// 3. Sign 签名示例
	claims := &JwtAdminClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "admin",
		},
		UserID:   1001,
		Username: "admin",
		Role:     "super_admin",
	}

	token, err := kitjwt.Get[*JwtAdminClaims]().Sign(claims)
	if err != nil {
		log.Fatalf("签名失败: %v", err)
	}
	fmt.Printf("生成的 Token: %s\n", token)

	// 4. Parse 解析示例
	parsedClaims, err := kitjwt.Get[*JwtAdminClaims]().Parse(
		"Bearer "+token,
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		log.Fatalf("解析失败: %v", err)
	}
	fmt.Printf("解析的 Claims: UserID=%d, Username=%s, Role=%s\n",
		parsedClaims.UserID, parsedClaims.Username, parsedClaims.Role)

	// 5. Gin 中间件使用示例
	r := gin.Default()

	// 公开路由
	r.GET("/login", func(c *gin.Context) {
		// 模拟登录，返回 token
		c.JSON(200, gin.H{"token": token})
	})

	// 需要认证的路由组
	auth := r.Group("/api")
	auth.Use(kitjwt.NewAuthMiddleware[*JwtAdminClaims]())
	{
		auth.GET("/profile", func(c *gin.Context) {
			// 从 context 获取 getClaims
			getClaims := kitjwt.GetClaims[*JwtAdminClaims](c)
			c.JSON(200, gin.H{
				"user_id":  getClaims.UserID,
				"username": getClaims.Username,
				"role":     getClaims.Role,
			})
		})
	}

	fmt.Println("\n启动服务器: http://localhost:8080")
	fmt.Println("测试命令:")
	fmt.Printf("  curl http://localhost:8080/login\n")
	fmt.Printf("  curl -H 'Authorization: Bearer %s' http://localhost:8080/api/profile\n", token)

	// 注意：根据用户规则，这里不实际运行服务器
	r.Run(":8080")
}

