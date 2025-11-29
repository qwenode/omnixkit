package kitrouter

import (
	"net/http"
	"sync"

	"connectrpc.com/connect"
	"github.com/gin-gonic/gin"
	"github.com/qwenode/omnixkit/kitcodec"
)

type CreateServiceFunc func(interceptors []connect.HandlerOption) (string, http.Handler)

// Option 函数式选项
type Option func(*Adapter)

// WithGuestMiddlewares 设置未登录中间件
func WithGuestMiddlewares(middlewares ...gin.HandlerFunc) Option {
	return func(a *Adapter) {
		a.guestMiddlewares = middlewares
	}
}

// WithAuthMiddlewares 设置需登录中间件
func WithAuthMiddlewares(middlewares ...gin.HandlerFunc) Option {
	return func(a *Adapter) {
		a.authMiddlewares = middlewares
	}
}

// WithInterceptors 设置connect拦截器
func WithInterceptors(interceptors ...connect.HandlerOption) Option {
	return func(a *Adapter) {
		a.interceptors = interceptors
	}
}

// route 统一路由结构
type route struct {
	method   string
	path     string
	handler  gin.HandlerFunc
	isCustom bool
	custom   func(route *gin.RouterGroup)
}

// Adapter 路由适配器
type Adapter struct {
	guestMiddlewares []gin.HandlerFunc
	authMiddlewares  []gin.HandlerFunc
	interceptors     []connect.HandlerOption
	guestRoutes      []route
	authRoutes       []route
}

var (
	instance     *Adapter
	authBuilder  *RouteBuilder
	guestBuilder *RouteBuilder
	once         sync.Once
)

// Bootstrap 初始化路由适配器（单例模式，只允许初始化一次）
func Bootstrap(opts ...Option) {
	if instance != nil {
		panic("Router adapter already initialized")
	}
	once.Do(func() {
		instance = &Adapter{
			interceptors: []connect.HandlerOption{kitcodec.WithProtoJSON()},
			guestRoutes:  make([]route, 0, 50),
			authRoutes:   make([]route, 0, 50),
		}
		for _, opt := range opts {
			opt(instance)
		}
		authBuilder = &RouteBuilder{adapter: instance, isAuth: true}
		guestBuilder = &RouteBuilder{adapter: instance, isAuth: false}
	})
}

// get 获取路由适配器实例
func get() *Adapter {
	if instance == nil {
		panic("Router adapter not initialized. Call Bootstrap() first.")
	}
	return instance
}

// Auth 返回需登录路由构建器
func Auth() *RouteBuilder {
	get() // 确保已初始化
	return authBuilder
}

// Guest 返回无需登录路由构建器
func Guest() *RouteBuilder {
	get() // 确保已初始化
	return guestBuilder
}

// Mount 加载路由到gin引擎
func Mount(engine *gin.Engine) {
	get().mount(engine)
}

// Interceptors 获取拦截器配置
func Interceptors() []connect.HandlerOption {
	return get().interceptors
}

// RouteBuilder 路由构建器
type RouteBuilder struct {
	adapter *Adapter
	isAuth  bool
}

// POST 添加POST服务
func (b *RouteBuilder) POST(callback CreateServiceFunc) *RouteBuilder {
	p, h := callback(b.adapter.interceptors)
	r := route{
		method:  http.MethodPost,
		path:    p + "*any",
		handler: gin.WrapH(h),
	}
	if b.isAuth {
		b.adapter.authRoutes = append(b.adapter.authRoutes, r)
	} else {
		b.adapter.guestRoutes = append(b.adapter.guestRoutes, r)
	}
	return b
}

// GET 添加GET服务
func (b *RouteBuilder) GET(callback CreateServiceFunc) *RouteBuilder {
	p, h := callback(b.adapter.interceptors)
	r := route{
		method:  http.MethodGet,
		path:    p,
		handler: gin.WrapH(h),
	}
	if b.isAuth {
		b.adapter.authRoutes = append(b.adapter.authRoutes, r)
	} else {
		b.adapter.guestRoutes = append(b.adapter.guestRoutes, r)
	}
	return b
}

// Custom 添加自定义路由
func (b *RouteBuilder) Custom(callback func(route *gin.RouterGroup)) *RouteBuilder {
	r := route{
		isCustom: true,
		custom:   callback,
	}
	if b.isAuth {
		b.adapter.authRoutes = append(b.adapter.authRoutes, r)
	} else {
		b.adapter.guestRoutes = append(b.adapter.guestRoutes, r)
	}
	return b
}

// mount 加载路由到gin引擎
func (a *Adapter) mount(engine *gin.Engine) {
	guestR := engine.Group("", a.guestMiddlewares...)
	authR := engine.Group("", a.authMiddlewares...)

	for _, r := range a.guestRoutes {
		if r.isCustom {
			r.custom(guestR)
		} else {
			guestR.Handle(r.method, r.path, r.handler)
		}
	}
	for _, r := range a.authRoutes {
		if r.isCustom {
			r.custom(authR)
		} else {
			authR.Handle(r.method, r.path, r.handler)
		}
	}
}
