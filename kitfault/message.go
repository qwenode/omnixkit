package kitfault

// Fault 定义错误消息的接口
type Fault interface {
    GetHint() string
}

// Accessor 访问错误消息的接口
type Accessor interface {
    GetFaultMessage() Fault
    SetFaultMessage(msg Fault)
}

// Factory 用于创建 Fault 实例的工厂函数类型
type Factory func(hint string) Fault

var factory Factory

// Bootstrap 初始化 Fault 工厂（启动时调用一次）
// 示例:
//
//	kitfault.Bootstrap(func(hint string) kitfault.Fault {
//	    return &msgpb.FaultMessage{Halt: true, Hint: hint}
//	})
func Bootstrap(f Factory) {
    factory = f
}

func getFactory() Factory {
    if factory == nil {
        panic("kitfault: not initialized. Call Bootstrap first.")
    }
    return factory
}

// IsHalted 判断流程是否应该停止并返回
func IsHalted(response Accessor) bool {
    if response == nil {
        return false
    }
    faultMessage := response.GetFaultMessage()
    if faultMessage == nil {
        return false
    }
    return faultMessage.GetHint() != ""
}

// Halt 停止执行后续流程并设置错误消息
func Halt(response Accessor, hint string) {
    response.SetFaultMessage(getFactory()(hint))
}

// Transfer 判断是否有错误，如果有，将 from 的错误传递到 to
func Transfer(from Accessor, to Accessor) bool {
    if from == nil {
        return false
    }
    message := from.GetFaultMessage()
    if message == nil {
        return false
    }
    if message.GetHint() != "" {
        to.SetFaultMessage(message)
        return true
    }
    return false
}
