package kitjwt

import (
    "crypto"
    "crypto/ed25519"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "fmt"
    "strings"
    "sync"

    "github.com/golang-jwt/jwt/v5"
)

var (
    // errInvalidToken 表示 token 无效
    errInvalidToken = errors.New("token invalid")
    // errUnexpectedSigningMethod 表示签名方法不匹配
    errUnexpectedSigningMethod = errors.New("unexpected signing method")
)

// IsInvalidToken 判断错误是否为 token 无效错误
func IsInvalidToken(err error) bool {
    return errors.Is(err, errInvalidToken)
}

// IsUnexpectedSigningMethod 判断错误是否为签名方法不匹配错误
func IsUnexpectedSigningMethod(err error) bool {
    return errors.Is(err, errUnexpectedSigningMethod)
}

// JWT 是一个通用的 JWT 工具接口，支持 ed25519 签名算法
type JWT[T jwt.Claims] interface {
    // Sign 对 claims 进行签名，返回 JWT token 字符串
    Sign(claims T) (string, error)
    // Parse 解析并验证 JWT token，返回 Claims
    // opts 可选的解析选项，如 jwt.WithExpirationRequired(), jwt.WithLeeway() 等
    Parse(token string, opts ...jwt.ParserOption) (T, error)
}

// jwtImpl 是一个通用的 JWT 工具实现，支持 ed25519 签名算法
type jwtImpl[T jwt.Claims] struct {
    key       ed25519.PrivateKey
    publicKey crypto.PublicKey
    algorithm jwt.SigningMethod
    newClaims func() T // 用于创建新的 claims 实例
}

// create 创建一个新的 JWT 实例
// keyHex: 十六进制编码的 ed25519 私钥（64字节，即128个十六进制字符）
// newClaims: 用于创建新的 claims 实例的函数
func create[T jwt.Claims](keyHex string, newClaims func() T) (*jwtImpl[T], error) {
    bytes, err := hex.DecodeString(keyHex)
    if err != nil {
        return nil, fmt.Errorf("failed to decode jwt key: %w", err)
    }
    if len(bytes) != ed25519.PrivateKeySize {
        return nil, fmt.Errorf("jwt key length incorrect: expected %d bytes, got %d", ed25519.PrivateKeySize, len(bytes))
    }

    privateKey := ed25519.PrivateKey(bytes)
    return &jwtImpl[T]{
        key:       privateKey,
        publicKey: privateKey.Public(),
        algorithm: jwt.SigningMethodEdDSA,
        newClaims: newClaims,
    }, nil
}

// Sign 对 claims 进行签名，返回 JWT token 字符串
func (j *jwtImpl[T]) Sign(claims T) (string, error) {
    token := jwt.NewWithClaims(j.algorithm, claims)
    return token.SignedString(j.key)
}

// Parse 解析并验证 JWT token，返回 Claims
// opts 可选的解析选项，如 jwt.WithExpirationRequired(), jwt.WithLeeway() 等
// 示例:
//
//	claims, err := jwt.Parse(token)
//	或
//	claims, err := jwt.Parse(token, jwt.WithExpirationRequired())
func (j *jwtImpl[T]) Parse(token string, opts ...jwt.ParserOption) (T, error) {
    var zero T
    if len(token) < 7 || !strings.EqualFold(token[:7], "bearer ") {
        return zero, errInvalidToken
    }

    claims := j.newClaims()
    parsed, err := jwt.ParseWithClaims(
        token[7:],
        claims,
        func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
                return nil, errUnexpectedSigningMethod
            }
            return j.publicKey, nil
        },
        opts...,
    )

    if err != nil {
        return zero, errInvalidToken
    }

    if !parsed.Valid {
        return zero, errInvalidToken
    }

    return claims, nil
}

var (
    instance any
    once     sync.Once
    initErr  error
)

// Bootstrap 初始化 JWT 实例（泛型版本）
// keyHex: 十六进制编码的 ed25519 私钥（64字节，即128个十六进制字符）
// newClaims: 用于创建新的 claims 实例的函数
// 示例:
//
//	kitjwt.Bootstrap("your-key-hex", func() *types.JwtAdminClaims {
//	    return &types.JwtAdminClaims{}
//	})
func Bootstrap[T jwt.Claims](keyHex string, newClaims func() T) error {
    if instance != nil {
        panic("JWT instance already initialized")
    }
    once.Do(func() {
        var impl *jwtImpl[T]
        impl, initErr = create(keyHex, newClaims)
        instance = impl
    })
    return initErr
}

// Get 获取 JWT 实例（单例模式）
// 必须先调用 Bootstrap 初始化
// 示例:
//
//	claims, err := kitjwt.Get[*types.JwtAdminClaims]().Parse(token, jwt.WithExpirationRequired())
func Get[T jwt.Claims]() JWT[T] {
    if instance == nil {
        panic("JWT instance not initialized. Call Bootstrap(keyHex, newClaims) first.")
    }
    return instance.(JWT[T])
}

// GenerateKey 生成一个新的 ed25519 私钥，并以十六进制字符串形式返回
func GenerateKey() (string, error) {
    _, privateKey, err := ed25519.GenerateKey(rand.Reader)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(privateKey), nil
}
