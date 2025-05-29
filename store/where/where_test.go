package where

import (
	"context"
	"reflect"
	"testing"
)

func TestOptions_P(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected Options
	}{
		{
			name:     "Valid page and pageSize",
			page:     2,
			pageSize: 10,
			expected: Options{Offset: 10, Limit: 10},
		},
		{
			name:     "Page is zero",
			page:     0,
			pageSize: 10,
			expected: Options{Offset: 0, Limit: 10},
		},
		{
			name:     "PageSize is zero",
			page:     2,
			pageSize: 0,
			expected: Options{Offset: -1, Limit: defaultLimit},
		},
		{
			name:     "Negative page and pageSize",
			page:     -1,
			pageSize: -5,
			expected: Options{Offset: 0, Limit: defaultLimit},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &Options{}
			options.P(tt.page, tt.pageSize)

			if options.Offset != tt.expected.Offset || options.Limit != tt.expected.Limit {
				t.Errorf("Expected Offset: %d, Limit: %d, got Offset: %d, Limit: %d",
					tt.expected.Offset, tt.expected.Limit, options.Offset, options.Limit)
			}
		})
	}
}

func TestOptions_Q(t *testing.T) {
	tests := []struct {
		name     string
		query    interface{}
		args     []interface{}
		expected []Query
	}{
		{
			name:     "Single query with args",
			query:    "name = ?",
			args:     []interface{}{"John"},
			expected: []Query{{Query: "name = ?", Args: []interface{}{"John"}}},
		},
		{
			name:     "Multiple queries",
			query:    "age > ?",
			args:     []interface{}{30},
			expected: []Query{{Query: "age > ?", Args: []interface{}{30}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &Options{}
			options.Q(tt.query, tt.args...)

			if !reflect.DeepEqual(options.Queries, tt.expected) {
				t.Errorf("Expected Queries: %v, got: %v", tt.expected, options.Queries)
			}
		})
	}
}

func TestOptions_T(t *testing.T) {
	RegisterTenant("tenant_id", func(ctx context.Context) string {
		return "12345"
	})

	tests := []struct {
		name     string
		ctx      context.Context
		expected map[any]any
	}{
		{
			name:     "Tenant value added to filters",
			ctx:      context.Background(),
			expected: map[any]any{"tenant_id": "12345"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &Options{Filters: map[any]any{}}
			options.T(tt.ctx)

			if !reflect.DeepEqual(options.Filters, tt.expected) {
				t.Errorf("Expected Filters: %v, got: %v", tt.expected, options.Filters)
			}
		})
	}
}

func TestOptions_F(t *testing.T) {
	tests := []struct {
		name       string
		initial    map[any]any
		input      []any
		expected   map[any]any
		shouldSkip bool
	}{
		{
			name:     "Valid key-value pairs",
			initial:  map[any]any{},
			input:    []any{"key1", "value1", "key2", "value2"},
			expected: map[any]any{"key1": "value1", "key2": "value2"},
		},
		{
			name:       "Uneven key-value pairs",
			initial:    map[any]any{},
			input:      []any{"key1", "value1", "key2"},
			expected:   map[any]any{},
			shouldSkip: true,
		},
		{
			name:     "Merge with existing filters",
			initial:  map[any]any{"existingKey": "existingValue"},
			input:    []any{"key1", "value1"},
			expected: map[any]any{"existingKey": "existingValue", "key1": "value1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &Options{
				Filters: tt.initial,
			}

			options.F(tt.input...)

			if tt.shouldSkip {
				if !reflect.DeepEqual(options.Filters, tt.initial) {
					t.Errorf("Expected filters to remain unchanged, got: %v", options.Filters)
				}
				return
			}

			if !reflect.DeepEqual(options.Filters, tt.expected) {
				t.Errorf("Expected filters: %v, got: %v", tt.expected, options.Filters)
			}
		})
	}
}

func TestOptions_O(t *testing.T) {
	tests := []struct {
		name     string
		offset   int
		expected Options
	}{
		{
			name:     "Valid offset",
			offset:   5,
			expected: Options{Offset: 5, Limit: defaultLimit},
		},
		{
			name:     "Negative offset",
			offset:   -3,
			expected: Options{Offset: 0, Limit: defaultLimit},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &Options{}
			options.O(tt.offset)

			if options.Offset != tt.expected.Offset || options.Limit != tt.expected.Limit {
				t.Errorf("Expected Offset: %d, Limit: %d, got Offset: %d, Limit: %d",
					tt.expected.Offset, tt.expected.Limit, options.Offset, options.Limit)
			}
		})
	}

}

func TestOptions_L(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		expected Options
	}{
		{
			name:     "Valid limit",
			limit:    10,
			expected: Options{Offset: 0, Limit: 10},
		},
		{
			name:     "Zero limit",
			limit:    0,
			expected: Options{Offset: 0, Limit: defaultLimit},
		},
		{
			name:     "Negative limit",
			limit:    -5,
			expected: Options{Offset: 0, Limit: defaultLimit},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &Options{}
			options.L(tt.limit)

			if options.Offset != tt.expected.Offset || options.Limit != tt.expected.Limit {
				t.Errorf("Expected Offset: %d, Limit: %d, got Offset: %d, Limit: %d",
					tt.expected.Offset, tt.expected.Limit, options.Offset, options.Limit)
			}
		})
	}
}

func TestOptions_WithOffset(t *testing.T) {
	tests := []struct {
		name     string
		offset   int64
		expected int
	}{
		{
			name:     "Valid offset",
			offset:   10,
			expected: 10,
		},
		{
			name:     "Negative offset",
			offset:   -5,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithOffset(tt.offset)
			options := &Options{}
			option(options)

			if options.Offset != tt.expected {
				t.Errorf("Expected Offset: %d, got: %d", tt.expected, options.Offset)
			}
		})
	}
}
func TestOptions_WithLimit(t *testing.T) {
	tests := []struct {
		name     string
		limit    int64
		expected int
	}{
		{
			name:     "Valid limit",
			limit:    20,
			expected: 20,
		},
		{
			name:     "Zero limit",
			limit:    0,
			expected: defaultLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithLimit(tt.limit)
			options := &Options{}
			option(options)

			if options.Limit != tt.expected {
				t.Errorf("Expected Limit: %d, got: %d", tt.expected, options.Limit)
			}
		})
	}
}
func TestOptions_WithPage(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		pageSize   int
		expected   Options
		shouldSkip bool
	}{
		{
			name:     "Valid page and pageSize",
			page:     2,
			pageSize: 10,
			expected: Options{Offset: 10, Limit: 10},
		},
		{
			name:       "Page is zero",
			page:       0,
			pageSize:   10,
			expected:   Options{Offset: 0, Limit: 10},
			shouldSkip: true,
		},
		{
			name:     "PageSize is zero",
			page:     2,
			pageSize: 0,
			expected: Options{Offset: 10, Limit: defaultLimit},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithPage(tt.page, tt.pageSize)
			options := &Options{}
			option(options)

			if tt.shouldSkip {
				if options.Offset != 0 || options.Limit != defaultLimit {
					t.Errorf("Expected Offset: 0, Limit: %d, got Offset: %d, Limit: %d",
						defaultLimit, options.Offset, options.Limit)
				}
				return
			}

			if options.Offset != tt.expected.Offset || options.Limit != tt.expected.Limit {
				t.Errorf("Expected Offset: %d, Limit: %d, got Offset: %d, Limit: %d",
					tt.expected.Offset, tt.expected.Limit, options.Offset, options.Limit)
			}
		})
	}
}
func TestOptions_WithFilter(t *testing.T) {
	tests := []struct {
		name     string
		filter   map[any]any
		expected map[any]any
	}{
		{
			name:     "Single filter",
			filter:   map[any]any{"key1": "value1"},
			expected: map[any]any{"key1": "value1"},
		},
		{
			name:     "Multiple filters",
			filter:   map[any]any{"key1": "value1", "key2": "value2"},
			expected: map[any]any{"key1": "value1", "key2": "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithFilter(tt.filter)
			options := &Options{}
			option(options)

			if !reflect.DeepEqual(options.Filters, tt.expected) {
				t.Errorf("Expected Filters: %v, got: %v", tt.expected, options.Filters)
			}
		})
	}
}
