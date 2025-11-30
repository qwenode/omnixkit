package kitboot

import "github.com/qwenode/omnixkit/kitfault"

// 常用项目启动配置
func Bootstrap(fault kitfault.FaultConstructor) {
    kitfault.Bootstrap(fault)
}
