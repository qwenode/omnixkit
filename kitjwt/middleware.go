package kitjwt

import (
    "connectrpc.com/connect"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/qwenode/omnixkit/kitctx"
)

const ginJwtClaimsKey = "_omnixkit_jwt"

// NewAuthMiddleware 创建 JWT 认证中间件
// 示例:
//
//	router.Use(kitjwt.NewAuthMiddleware[*types.JwtAdminClaims]())
func NewAuthMiddleware[T jwt.Claims]() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenHeader := c.GetHeader("Authorization")
        claims, err := Get[T]().Parse(tokenHeader, jwt.WithExpirationRequired())
        if err != nil {
            _ = connect.NewErrorWriter().Write(c.Writer, c.Request, kitctx.NewUnauthenticatedErr(err))
            c.Abort()
            return
        }
        c.Set(ginJwtClaimsKey, claims)
        c.Next()
    }
}

// GetClaims 从 gin context 获取 claims
// 示例:
//
//	claims := kitjwt.GetClaims[*types.JwtAdminClaims](c)
func GetClaims[T jwt.Claims](c *gin.Context) T {
    return c.MustGet(ginJwtClaimsKey).(T)
}
