package main

import (
	"fmt"

	"github.com/qwenode/omnixkit/kitfault"
)

// ========== 模拟 protobuf 生成的结构 ==========

// FaultMessage 模拟 protobuf 生成的 msgpb.FaultMessage
type FaultMessage struct {
	Halt bool
	Hint string
}

func (f *FaultMessage) GetHalt() bool {
	if f == nil {
		return false
	}
	return f.Halt
}

func (f *FaultMessage) GetHint() string {
	if f == nil {
		return ""
	}
	return f.Hint
}

// Response 模拟 protobuf 生成的 Response，内嵌 FaultMessage
type Response struct {
	Data         string
	FaultMessage *FaultMessage
}

func (r *Response) GetFaultMessage() kitfault.Fault {
	if r.FaultMessage == nil {
		return nil
	}
	return r.FaultMessage
}

func (r *Response) SetFaultMessage(msg kitfault.Fault) {
	if fm, ok := msg.(*FaultMessage); ok {
		r.FaultMessage = fm
	}
}

// 确保 Response 实现 Accessor 接口
var _ kitfault.Accessor = (*Response)(nil)

// ========== 项目初始化（启动时调用一次） ==========

func init() {
	kitfault.Bootstrap(func(halt bool, hint string) kitfault.Fault {
		return &FaultMessage{
			Halt: halt,
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
	if kitfault.Transfer(resp, resp2) {
		fmt.Printf("错误已传递: %s\n", resp2.GetFaultMessage().GetHint())
	}

	// 示例4: GetOrCreate
	resp3 := &Response{}
	fmt.Printf("resp3.FaultMessage 初始为 nil: %v\n", resp3.FaultMessage == nil)
	kitfault.GetOrCreate(resp3)
	fmt.Printf("调用 GetOrCreate 后: %v\n", resp3.FaultMessage != nil)
}
