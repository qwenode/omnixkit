package kitboot

import "github.com/qwenode/omnixkit/kitfault"

// 常用项目启动配置
func Bootstrap[T kitfault.Fault](fault kitfault.FaultConstructor[T]) {
    kitfault.Bootstrap(fault)
}
