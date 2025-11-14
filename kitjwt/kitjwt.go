package kitjwt

import (
    "crypto"
    "crypto/ed25519"
    "encoding/hex"
    "errors"
    "fmt"
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
type JWT interface {
    // Sign 对 claims 进行签名，返回 JWT token 字符串
    Sign(claims jwt.Claims) (string, error)
    // Parse 解析并验证 JWT token，返回 Claims
    // claims 必须是实现了 jwt.Claims 接口的类型指针
    // opts 可选的解析选项，如 jwt.WithExpirationRequired(), jwt.WithLeeway() 等
    Parse(claims jwt.Claims, token string, opts ...jwt.ParserOption) error
}

// jwtImpl 是一个通用的 JWT 工具实现，支持 ed25519 签名算法
type jwtImpl struct {
    key       ed25519.PrivateKey
    publicKey crypto.PublicKey
    algorithm jwt.SigningMethod
}

// create 创建一个新的 JWT 实例
// keyHex: 十六进制编码的 ed25519 密钥（64字节，即128个十六进制字符）
func create(keyHex string) (*jwtImpl, error) {
    bytes, err := hex.DecodeString(keyHex)
    if err != nil {
        return nil, fmt.Errorf("failed to decode jwt key: %w", err)
    }
    if len(bytes) != ed25519.SignatureSize {
        return nil, fmt.Errorf("jwt key length incorrect: expected %d bytes, got %d", ed25519.SignatureSize, len(bytes))
    }

    privateKey := ed25519.NewKeyFromSeed(bytes[:ed25519.SeedSize])
    return &jwtImpl{
        key:       privateKey,
        publicKey: privateKey.Public(),
        algorithm: jwt.SigningMethodEdDSA,
    }, nil
}

// Sign 对 claims 进行签名，返回 JWT token 字符串
func (j *jwtImpl) Sign(claims jwt.Claims) (string, error) {
    token := jwt.NewWithClaims(j.algorithm, claims)
    return token.SignedString(j.key)
}

// Parse 解析并验证 JWT token，返回 Claims
// claims 必须是实现了 jwt.Claims 接口的类型指针
// opts 可选的解析选项，如 jwt.WithExpirationRequired(), jwt.WithLeeway() 等
// 示例: 
//   var claims MyClaims
//   err := jwt.Parse(&claims, token)
//   或
//   claims := &MyClaims{}
//   err := jwt.Parse(claims, token, jwt.WithExpirationRequired())
func (j *jwtImpl) Parse(claims jwt.Claims, token string, opts ...jwt.ParserOption) error {
    parsed, err := jwt.ParseWithClaims(
        token,
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
        return errInvalidToken
    }

    if !parsed.Valid {
        return errInvalidToken
    }

    return nil
}

var (
    instance JWT
    once     sync.Once
    initErr  error
)

// Get 获取 JWT 实例（单例模式）
// keyHex: 十六进制编码的 ed25519 密钥（64字节，即128个十六进制字符）
// 第一次调用时会使用 keyHex 初始化实例，后续调用会直接返回已初始化的实例
func Get(keyHex string) (JWT, error) {
    once.Do(func() {
        var impl *jwtImpl
        impl, initErr = create(keyHex)
        instance = impl
    })
    return instance, initErr
}
