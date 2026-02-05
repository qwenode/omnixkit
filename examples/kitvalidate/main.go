package main

import (
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/qwenode/omnixkit/kitvalidate"
)

func main() {
	// ============================================
	// kitvalidate 使用示例
	// 这是一个用于 Connect RPC 的 protovalidate 验证拦截器
	// ============================================

	// 1. 创建基本的验证拦截器（使用默认配置）
	basicInterceptor, err := kitvalidate.NewInterceptor()
	if err != nil {
		log.Fatalf("创建验证拦截器失败: %v", err)
	}
	fmt.Printf("基本拦截器创建成功: %T\n", basicInterceptor)

	// 2. 创建带有自定义 ErrorDetailBuilder 的拦截器
	// ErrorDetailBuilder 用于将验证错误转换为 connect.ErrorDetail
	customBuilder := kitvalidate.ErrorDetailBuilderFunc(
		func(errors []kitvalidate.ValidationError) (*connect.ErrorDetail, error) {
			// 这里可以自定义错误详情的构建逻辑
			// 例如：将验证错误转换为自定义的 protobuf 消息
			fmt.Println("验证错误列表:")
			for _, e := range errors {
				fmt.Printf("  - 字段: %s, 消息: %s\n", e.Field, e.Message)
			}

			// 示例：返回 nil 表示不添加额外的错误详情
			// 实际使用时，可以返回自定义的 protobuf 消息
			// 例如:
			// fieldErr := &yourpb.ValidateMessages{
			//     Fields: make([]*yourpb.ValidateMessage, 0, len(errors)),
			// }
			// for _, e := range errors {
			//     fieldErr.Fields = append(fieldErr.Fields, &yourpb.ValidateMessage{
			//         Field:   e.Field,
			//         Message: e.Message,
			//     })
			// }
			// return connect.NewErrorDetail(fieldErr)

			return nil, nil
		},
	)

	interceptorWithBuilder, err := kitvalidate.NewInterceptor(
		kitvalidate.WithErrorDetailBuilder(customBuilder),
	)
	if err != nil {
		log.Fatalf("创建自定义拦截器失败: %v", err)
	}
	fmt.Printf("自定义拦截器创建成功: %T\n", interceptorWithBuilder)

	// 3. 在 Connect RPC handler 中使用拦截器
	// 以下是示例代码，展示如何将拦截器应用到 Connect handler

	// 创建带有验证拦截器的 handler 选项
	handlerOptions := connect.WithInterceptors(basicInterceptor)

	// 使用示例（需要根据你的 proto 定义调整）:
	//
	// path, handler := yourpbconnect.NewYourServiceHandler(
	//     &YourServiceImpl{},
	//     handlerOptions,
	// )
	//
	// mux := http.NewServeMux()
	// mux.Handle(path, handler)
	// http.ListenAndServe(":8080", mux)

	fmt.Println("\n=== 使用说明 ===")
	fmt.Println("1. 定义带有 protovalidate 规则的 proto 文件")
	fmt.Println("2. 生成 Go 代码")
	fmt.Println("3. 创建 kitvalidate.Interceptor")
	fmt.Println("4. 将拦截器添加到 Connect handler")
	fmt.Println()
	fmt.Println("Proto 文件示例 (使用 buf.build/bufbuild/protovalidate):")
	fmt.Print(`
syntax = "proto3";

import "buf/validate/validate.proto";

message CreateUserRequest {
    string email = 1 [(buf.validate.field).string.email = true];
    string name = 2 [(buf.validate.field).string = {
        min_len: 2,
        max_len: 100
    }];
    int32 age = 3 [(buf.validate.field).int32 = {
        gte: 0,
        lte: 150
    }];
}
`)

	// 4. 完整的服务器示例
	fmt.Println("=== 完整服务器示例 ===")

	// 创建 HTTP mux
	mux := http.NewServeMux()

	// 健康检查端点
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 在实际项目中，你会这样注册 Connect handler:
	// path, handler := yourpbconnect.NewYourServiceHandler(
	//     &YourServiceImpl{},
	//     connect.WithInterceptors(interceptorWithBuilder),
	// )
	// mux.Handle(path, handler)

	fmt.Println("服务器配置完成")
	fmt.Printf("Handler 选项: %T\n", handlerOptions)
	fmt.Println()
	fmt.Println("当请求不符合 protovalidate 规则时，拦截器会:")
	fmt.Println("1. 返回 connect.CodeInvalidArgument 错误")
	fmt.Println("2. 如果配置了 ErrorDetailBuilder，会添加自定义错误详情")
	fmt.Println()
	fmt.Println("支持的 Connect 模式:")
	fmt.Println("- Unary RPC (一元调用)")
	fmt.Println("- Streaming Client (客户端流)")
	fmt.Println("- Streaming Handler (服务端流)")
}
