package kitrpc

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "net"
    "net/http"
    "net/url"
    "strings"

    "github.com/gin-gonic/gin"
)

const ginContextKey = "_omnixkit_context"

// 让rpc能够访问gin的上下文, 需要在gin的中间件里使用这个中间件
func MiddlewareAdapterGinContext() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := context.WithValue(c.Request.Context(), ginContextKey, c)
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}

// rpc内获取 gin.Context
func GetGinContext(ctx context.Context) (*gin.Context, error) {
    g, ok := ctx.Value(ginContextKey).(*gin.Context)
    if !ok {
        return nil, NewInternal("Context Not Found")
    }
    return g, nil
}

const clientIpKey = "_omnixkit_client_ip"

// 从context中获取客户端IP地址
func GetClientIp(c context.Context) string {
    value := c.Value(clientIpKey)
    if value == nil {
        return "1.1.1.1"
    }
    return value.(string)
}

// 从允许的Header列表中设置客户端IP到Context
func MiddlewareSetClientIp(allowHeaders []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        clientIP := ""
        for _, header := range allowHeaders {
            headerValue := c.GetHeader(header)
            if headerValue == "" {
                continue
            }
            ip := strings.TrimSpace(strings.Split(headerValue, ",")[0])
            if net.ParseIP(ip) != nil {
                clientIP = ip
                break
            }
        }
        if clientIP == "" {
            clientIP = c.RemoteIP()
        }
        c.Set(clientIpKey, clientIP)
        c.Next()
    }
}

func MiddlewareAdapterMethodGet() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Method == http.MethodGet {
            if c.Query("message") == "" {
                queryMap := make(map[string]any)
                for key, values := range c.Request.URL.Query() {
                    if key == "encoding" || key == "message" {
                        continue
                    }
                    if len(values) == 1 {
                        queryMap[key] = values[0]
                    } else {
                        queryMap[key] = values
                    }
                }
                messageJSON, _ := json.Marshal(queryMap)
                newQuery := url.Values{}
                newQuery.Set("message", string(messageJSON))
                c.Request.URL.RawQuery = newQuery.Encode()
            }
            if c.Query("encoding") != "json" {
                c.Request.URL.RawQuery += "&encoding=json"
            }
            c.Request.Body = io.NopCloser(bytes.NewReader([]byte{}))
            c.Request.ContentLength = -1
        }
        c.Next()
    }
}
