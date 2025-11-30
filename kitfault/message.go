package kitfault
//
//// Fault 定义错误消息的接口
//type Fault interface {
//    GetHint() string
//    SetHint(hint string) Fault
//    Clone () Fault
//}
//
//// FaultHolder 持有 Fault 的对象接口（使用泛型支持 protobuf 生成的具体类型）
//type FaultHolder[T Fault] interface {
//    GetFaultMessage() T
//    SetFaultMessage(msg T)
//}
//
//// FaultConstructor 用于创建 Fault 实例的构造函数类型
//type FaultConstructor[T Fault] func(hint string) T
//
//var constructorAny any
//var FaultModel Fault
//// Bootstrap 初始化 Fault 构造函数（启动时调用一次）
//// 示例:
////
////	kitfault.Bootstrap(func(hint string) *msgpb.FaultMessage {
////	    return &msgpb.FaultMessage{Hint: hint}
////	})
//func Bootstrap[T Fault](fn FaultConstructor[T]) {
//    constructorAny = fn
//}
//
//func getConstructor[T Fault]() FaultConstructor[T] {
//    return constructorAny.(FaultConstructor[T])
//}
//
//// IsHalted 判断流程是否应该停止并返回
//func IsHalted[T Fault](holder FaultHolder[T]) bool {
//    if holder == nil {
//        return false
//    }
//    fault := holder.GetFaultMessage()
//    // 使用 any 转换检查 nil（泛型类型的零值检查）
//    if any(fault) == nil {
//        return false
//    }
//    return fault.GetHint() != ""
//}
//
//// Halt 停止执行后续流程并设置错误消息
//func Halt(holder FaultHolder[Fault], hint string) {
//    
//    holder.SetFaultMessage(FaultModel.Clone().SetHint(hint))
//}
//
//// Forward 判断是否有错误，如果有，将 from 的错误传递到 to
//func Forward[T Fault](from FaultHolder[T], to FaultHolder[T]) bool {
//    if from == nil {
//        return false
//    }
//    fault := from.GetFaultMessage()
//    if any(fault) == nil {
//        return false
//    }
//    if fault.GetHint() != "" {
//        to.SetFaultMessage(fault)
//        return true
//    }
//    return false
//}
