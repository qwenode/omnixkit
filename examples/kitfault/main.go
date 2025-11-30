package main

import (
	"fmt"

	"github.com/qwenode/omnixkit/kitfault"
)

// ========== 模拟 protobuf 生成的结构 ==========

// FaultMessage 模拟 protobuf 生成的 msgpb.FaultMessage
type FaultMessage struct {
	Hint string
}


func (f *FaultMessage) GetHint() string {
	if f == nil {
		return ""
	}
	return f.Hint
}

// Response 模拟 protobuf 生成的 Response，内嵌 FaultMessage
// 现在 GetFaultMessage 返回具体类型 *FaultMessage，与 protobuf 生成的代码一致
type Response struct {
	Data         string
	FaultMessage *FaultMessage
}

func (r *Response) GetFaultMessage() *FaultMessage {
	return r.FaultMessage
}

func (r *Response) SetFaultMessage(msg *FaultMessage) {
	r.FaultMessage = msg
}

// ========== 项目初始化（启动时调用一次） ==========

func init() {
	kitfault.Bootstrap(func(hint string) *FaultMessage {
		return &FaultMessage{
			Hint: hint,
		}
	})
}

// ========== 使用示例 ==========

func main() {
	// 示例1: 检查是否需要停止
	resp := &Response{}
	if kitfault.IsHalted(resp) {
		fmt.Println("流程已停止")
	} else {
		fmt.Println("流程正常")
	}

	// 示例2: 主动设置停止
	kitfault.Halt(resp, "参数验证失败")
	fmt.Printf("手动停止: %s\n", resp.GetFaultMessage().GetHint())

	// 示例3: 错误传递
	resp2 := &Response{}
	if kitfault.Forward(resp, resp2) {
		fmt.Printf("错误已传递: %s\n", resp2.GetFaultMessage().GetHint())
	}
}
