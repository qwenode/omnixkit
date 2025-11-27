package kitctx

import (
    "context"

    "connectrpc.com/connect"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/qwenode/omnixkit/kitjwt"
)

const ginJwtClaimsKey = "_omnixkit_jwt"

// NewAuthMiddleware 创建 JWT 认证中间件
// 示例:
//
//	router.Use(kitctx.NewAuthMiddleware[*types.JwtAdminClaims]())
func NewAuthMiddleware[T jwt.Claims]() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenHeader := c.GetHeader("Authorization")
        claims, err := kitjwt.Get[T]().Parse(tokenHeader, jwt.WithExpirationRequired())
        if err != nil {
            _ = connect.NewErrorWriter().Write(c.Writer, c.Request, NewUnauthenticatedErr(err))
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
//	claims := kitctx.GetClaims[*types.JwtAdminClaims](c)
func GetClaims[T jwt.Claims](c context.Context) (T, error) {
    value := c.Value(ginJwtClaimsKey)
    if value == nil {
        var zero T
        return zero, NewUnauthenticatedErr(nil)
    }
    claims, ok := value.(T)
    if !ok {
        var zero T
        return zero, NewUnauthenticatedErr(nil)
    }
    return claims, nil
}
