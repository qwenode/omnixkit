package kitset

type PageSize interface {
    GetSize() int32
    GetPage() int32
    SetSize(size int32)
    SetPage(page int32)
}

// 根据分页参数计算 offset 和 limit
// page 从 1 开始，size 为每页数量
// 返回值用于数据库查询的 OFFSET 和 LIMIT
func PageToOffsetLimit(input PageSize) (_offset int, _limit int) {
    page := int(input.GetPage())
    size := int(input.GetSize())
    if page < 1 {
        page = 1
    }
    if size < 1 {
        size = 20 // 默认每页 10 条
    }
    return (page - 1) * size, size
}
// 根据分页参数计算 offset 和 limit，限制最大 size 为 1000
func PageToOffsetLimitDefault(input PageSize) (_offset int, _limit int) {
    if input.GetSize() > 1000 {
        input.SetSize(1000)
    }
    return PageToOffsetLimit(input)
}
// 根据分页参数计算 offset 和 limit，限制最大 size 和 page
func PageToOffsetLimitMax(input PageSize,maxSize int32,maxPage int32) (_offset int, _limit int){
    if input.GetSize() > maxSize {
        input.SetSize(maxSize)
    }
    if input.GetPage() > maxPage {
        input.SetPage(maxPage)
    }
    return PageToOffsetLimit(input)
}