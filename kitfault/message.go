package kitfault

// Fault 定义错误消息的接口
type Fault interface {
    GetHint() string
}

// FaultHolder 持有 Fault 的对象接口
type FaultHolder interface {
    GetFaultMessage() Fault
    SetFaultMessage(msg Fault)
}

// FaultConstructor 用于创建 Fault 实例的构造函数类型
type FaultConstructor func(hint string) Fault

var constructor FaultConstructor

// Bootstrap 初始化 Fault 构造函数（启动时调用一次）
// 示例:
//
//	kitfault.Bootstrap(func(hint string) kitfault.Fault {
//	    return &msgpb.FaultMessage{Halt: true, Hint: hint}
//	})
func Bootstrap(fn FaultConstructor) {
    constructor = fn
}

func mustGetConstructor() FaultConstructor {
    if constructor == nil {
        panic("kitfault: not initialized. Call Bootstrap first.")
    }
    return constructor
}

// IsHalted 判断流程是否应该停止并返回
func IsHalted(holder FaultHolder) bool {
    if holder == nil {
        return false
    }
    fault := holder.GetFaultMessage()
    if fault == nil {
        return false
    }
    return fault.GetHint() != ""
}

// Halt 停止执行后续流程并设置错误消息
func Halt(holder FaultHolder, hint string) {
    holder.SetFaultMessage(mustGetConstructor()(hint))
}

// Forward 判断是否有错误，如果有，将 from 的错误传递到 to
func Forward(from FaultHolder, to FaultHolder) bool {
    if from == nil {
        return false
    }
    fault := from.GetFaultMessage()
    if fault == nil {
        return false
    }
    if fault.GetHint() != "" {
        to.SetFaultMessage(fault)
        return true
    }
    return false
}
