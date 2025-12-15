package kitset

import "testing"

// mockPageSize 实现 PageSize 接口用于测试
type mockPageSize struct {
	page int32
	size int32
}

func (m *mockPageSize) GetPage() int32     { return m.page }
func (m *mockPageSize) GetSize() int32     { return m.size }
func (m *mockPageSize) SetPage(page int32) { m.page = page }
func (m *mockPageSize) SetSize(size int32) { m.size = size }

func TestPageToOffsetLimit(t *testing.T) {
	tests := []struct {
		name       string
		page       int32
		size       int32
		wantOffset int
		wantLimit  int
	}{
		{
			name:       "normal case page=1 size=10",
			page:       1,
			size:       10,
			wantOffset: 0,
			wantLimit:  10,
		},
		{
			name:       "normal case page=2 size=10",
			page:       2,
			size:       10,
			wantOffset: 10,
			wantLimit:  10,
		},
		{
			name:       "normal case page=3 size=20",
			page:       3,
			size:       20,
			wantOffset: 40,
			wantLimit:  20,
		},
		{
			name:       "page=0 should default to 1",
			page:       0,
			size:       10,
			wantOffset: 0,
			wantLimit:  10,
		},
		{
			name:       "negative page should default to 1",
			page:       -5,
			size:       10,
			wantOffset: 0,
			wantLimit:  10,
		},
		{
			name:       "size=0 should default to 20",
			page:       1,
			size:       0,
			wantOffset: 0,
			wantLimit:  20,
		},
		{
			name:       "negative size should default to 20",
			page:       1,
			size:       -10,
			wantOffset: 0,
			wantLimit:  20,
		},
		{
			name:       "both page and size invalid",
			page:       0,
			size:       0,
			wantOffset: 0,
			wantLimit:  20,
		},
		{
			name:       "large page number",
			page:       100,
			size:       50,
			wantOffset: 4950,
			wantLimit:  50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &mockPageSize{page: tt.page, size: tt.size}
			gotOffset, gotLimit := PageToOffsetLimit(input)
			if gotOffset != tt.wantOffset {
				t.Errorf("PageToOffsetLimit() offset = %v, want %v", gotOffset, tt.wantOffset)
			}
			if gotLimit != tt.wantLimit {
				t.Errorf("PageToOffsetLimit() limit = %v, want %v", gotLimit, tt.wantLimit)
			}
		})
	}
}

func TestPageToOffsetLimitDefault(t *testing.T) {
	tests := []struct {
		name       string
		page       int32
		size       int32
		wantOffset int
		wantLimit  int
		wantSize   int32 // 检查 size 是否被修改
	}{
		{
			name:       "size within default max",
			page:       1,
			size:       100,
			wantOffset: 0,
			wantLimit:  100,
			wantSize:   100,
		},
		{
			name:       "size equals default max",
			page:       1,
			size:       1000,
			wantOffset: 0,
			wantLimit:  1000,
			wantSize:   1000,
		},
		{
			name:       "size exceeds default max should be capped",
			page:       1,
			size:       2000,
			wantOffset: 0,
			wantLimit:  1000,
			wantSize:   1000,
		},
		{
			name:       "page 2 with size exceeds max",
			page:       2,
			size:       1500,
			wantOffset: 1000,
			wantLimit:  1000,
			wantSize:   1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &mockPageSize{page: tt.page, size: tt.size}
			gotOffset, gotLimit := PageToOffsetLimitDefault(input)
			if gotOffset != tt.wantOffset {
				t.Errorf("PageToOffsetLimitDefault() offset = %v, want %v", gotOffset, tt.wantOffset)
			}
			if gotLimit != tt.wantLimit {
				t.Errorf("PageToOffsetLimitDefault() limit = %v, want %v", gotLimit, tt.wantLimit)
			}
			if input.GetSize() != tt.wantSize {
				t.Errorf("PageToOffsetLimitDefault() size was modified to %v, want %v", input.GetSize(), tt.wantSize)
			}
		})
	}
}

func TestPageToOffsetLimitMax(t *testing.T) {
	tests := []struct {
		name       string
		page       int32
		size       int32
		maxSize    int32
		maxPage    int32
		wantOffset int
		wantLimit  int
		wantPage   int32
		wantSize   int32
	}{
		{
			name:       "within all limits",
			page:       2,
			size:       50,
			maxSize:    100,
			maxPage:    10,
			wantOffset: 50,
			wantLimit:  50,
			wantPage:   2,
			wantSize:   50,
		},
		{
			name:       "size exceeds max",
			page:       1,
			size:       200,
			maxSize:    100,
			maxPage:    10,
			wantOffset: 0,
			wantLimit:  100,
			wantPage:   1,
			wantSize:   100,
		},
		{
			name:       "page exceeds max",
			page:       15,
			size:       50,
			maxSize:    100,
			maxPage:    10,
			wantOffset: 450,
			wantLimit:  50,
			wantPage:   10,
			wantSize:   50,
		},
		{
			name:       "both page and size exceed max",
			page:       20,
			size:       200,
			maxSize:    100,
			maxPage:    10,
			wantOffset: 900,
			wantLimit:  100,
			wantPage:   10,
			wantSize:   100,
		},
		{
			name:       "equals max values",
			page:       10,
			size:       100,
			maxSize:    100,
			maxPage:    10,
			wantOffset: 900,
			wantLimit:  100,
			wantPage:   10,
			wantSize:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &mockPageSize{page: tt.page, size: tt.size}
			gotOffset, gotLimit := PageToOffsetLimitMax(input, tt.maxSize, tt.maxPage)
			if gotOffset != tt.wantOffset {
				t.Errorf("PageToOffsetLimitMax() offset = %v, want %v", gotOffset, tt.wantOffset)
			}
			if gotLimit != tt.wantLimit {
				t.Errorf("PageToOffsetLimitMax() limit = %v, want %v", gotLimit, tt.wantLimit)
			}
			if input.GetPage() != tt.wantPage {
				t.Errorf("PageToOffsetLimitMax() page was modified to %v, want %v", input.GetPage(), tt.wantPage)
			}
			if input.GetSize() != tt.wantSize {
				t.Errorf("PageToOffsetLimitMax() size was modified to %v, want %v", input.GetSize(), tt.wantSize)
			}
		})
	}
}

func TestGetPageSize(t *testing.T) {
	tests := []struct {
		name     string
		page     int32
		size     int32
		wantPage int
		wantSize int
	}{
		{
			name:     "normal case",
			page:     2,
			size:     30,
			wantPage: 2,
			wantSize: 30,
		},
		{
			name:     "page=1 size=20",
			page:     1,
			size:     20,
			wantPage: 1,
			wantSize: 20,
		},
		{
			name:     "page=0 should default to 1",
			page:     0,
			size:     10,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "negative page should default to 1",
			page:     -3,
			size:     10,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "size=0 should default to 20",
			page:     1,
			size:     0,
			wantPage: 1,
			wantSize: 20,
		},
		{
			name:     "negative size should default to 20",
			page:     1,
			size:     -5,
			wantPage: 1,
			wantSize: 20,
		},
		{
			name:     "both invalid",
			page:     -1,
			size:     -1,
			wantPage: 1,
			wantSize: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &mockPageSize{page: tt.page, size: tt.size}
			gotPage, gotSize := GetPageSize(input)
			if gotPage != tt.wantPage {
				t.Errorf("GetPageSize() page = %v, want %v", gotPage, tt.wantPage)
			}
			if gotSize != tt.wantSize {
				t.Errorf("GetPageSize() size = %v, want %v", gotSize, tt.wantSize)
			}
		})
	}
}

func TestDefaultMaxSize(t *testing.T) {
	// 测试 DefaultMaxSize 的默认值
	if DefaultMaxSize != 1000 {
		t.Errorf("DefaultMaxSize = %v, want 1000", DefaultMaxSize)
	}

	// 测试修改 DefaultMaxSize
	originalMax := DefaultMaxSize
	defer func() { DefaultMaxSize = originalMax }()

	DefaultMaxSize = 500
	input := &mockPageSize{page: 1, size: 800}
	_, gotLimit := PageToOffsetLimitDefault(input)
	if gotLimit != 500 {
		t.Errorf("PageToOffsetLimitDefault() with modified DefaultMaxSize limit = %v, want 500", gotLimit)
	}
}
