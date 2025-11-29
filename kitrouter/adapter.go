package kitrouter

import (
    "net/http"

    "connectrpc.com/connect"
    "github.com/gin-gonic/gin"
    "github.com/qwenode/omnixkit/kitcodec"
)

type AuthType int

const (
    AuthRequired AuthType = iota
    AuthGuest
)

type Service struct {
    Method  string
    Path    string
    Handler gin.HandlerFunc
    Auth    AuthType
}

type CustomService struct {
    Auth   AuthType
    Create func(route *gin.RouterGroup)
}

type CreateServiceFunc func(interceptors []connect.HandlerOption) (string, http.Handler)

// RouterConfig 路由配置
type RouterConfig struct {
    GuestMiddlewares []gin.HandlerFunc       // 未登录中间件
    AuthMiddlewares  []gin.HandlerFunc       // 需登录中间件
    Interceptors     []connect.HandlerOption // connect拦截器
}

// Adapter 路由适配器
type Adapter struct {
    config         RouterConfig
    services       []Service
    customServices []CustomService
}

func NewAdapter(config RouterConfig) *Adapter {
    if config.Interceptors == nil {
        config.Interceptors = []connect.HandlerOption{kitcodec.WithProtoJSON()}
    }
    return &Adapter{
        config:         config,
        services:       make([]Service, 0, 50),
        customServices: make([]CustomService, 0),
    }
}

// POST 添加POST服务
func (a *Adapter) POST(auth AuthType, callback CreateServiceFunc) {
    p, h := callback(a.config.Interceptors)
    a.services = append(a.services, Service{
        Method:  http.MethodPost,
        Path:    p + "*any",
        Handler: gin.WrapH(h),
        Auth:    auth,
    })
}

// GET 添加GET服务
func (a *Adapter) GET(auth AuthType, callback CreateServiceFunc) {
    p, h := callback(a.config.Interceptors)
    a.services = append(a.services, Service{
        Method:  http.MethodGet,
        Path:    p,
        Handler: gin.WrapH(h),
        Auth:    auth,
    })
}

// Custom 添加自定义路由
func (a *Adapter) Custom(auth AuthType, callback func(route *gin.RouterGroup)) {
    a.customServices = append(a.customServices, CustomService{
        Auth:   auth,
        Create: callback,
    })
}

// Mount 加载路由到gin引擎
func (a *Adapter) Mount(route *gin.Engine) {
    guestR := route.Group("", a.config.GuestMiddlewares...)
    authR := route.Group("", a.config.AuthMiddlewares...)

    for _, s := range a.services {
        if s.Auth == AuthGuest {
            guestR.Handle(s.Method, s.Path, s.Handler)
        } else {
            authR.Handle(s.Method, s.Path, s.Handler)
        }
    }
    for _, s := range a.customServices {
        if s.Auth == AuthGuest {
            s.Create(guestR)
        } else {
            s.Create(authR)
        }
    }
}

// Interceptors 获取拦截器配置
func (a *Adapter) Interceptors() []connect.HandlerOption {
    return a.config.Interceptors
}
