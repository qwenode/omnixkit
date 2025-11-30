package main
//
//import (
//    "fmt"
//
//    "github.com/qwenode/omnixkit/examples/kitfault/tyy"
//    "github.com/qwenode/omnixkit/kitfault"
//)
//
//// ========== 模拟 protobuf 生成的结构 ==========
//
//// ========== 项目初始化（启动时调用一次） ==========
//
//func init() {
//
//}
//
//// ========== 使用示例 ==========
//
//func main() {
//    // 示例1: 检查是否需要停止
//    resp := &tyy.Response{}
//    //if kitfault.IsHalted(resp) {
//    //    fmt.Println("流程已停止")
//    //} else {
//    //    fmt.Println("流程正常")
//    //}
//    //kitfault.Bootstrap(func(hint string) *types.FaultMessage {
//    //    return &types.FaultMessage{
//    //        Hint: hint,
//    //    }
//    //})
//    // 示例2: 主动设置停止
//    kitfault.Halt(resp, "参数验证失败")
//    fmt.Printf("手动停止: %s\n", resp.GetFaultMessage().GetHint())
//
//    // 示例3: 错误传递
//    resp2 := &tyy.Response{}
//    if kitfault.Forward(resp, resp2) {
//        fmt.Printf("错误已传递: %s\n", resp2.GetFaultMessage().GetHint())
//    }
//}
