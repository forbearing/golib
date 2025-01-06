package client

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/forbearing/golib/types/consts"
	"gorm.io/gorm"
)

type Query struct {
	ID string `json:"id" gorm:"primaryKey" schema:"id" url:"-"`

	CreatedBy string     `json:"created_by,omitempty" schema:"created_by" gorm:"index" url:"-"`
	UpdatedBy string     `json:"updated_by,omitempty" schema:"updated_by" gorm:"index" url:"-"`
	CreatedAt *time.Time `json:"created_at,omitempty" schema:"-" gorm:"index" url:"-"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" schema:"-" gorm:"index" url:"-"`
	Remark    *string    `json:"remark,omitempty" gorm:"size:10240" schema:"-" url:"-"` // 如果需要支持 PATCH 更新,则必须是指针类型
	Order     *uint      `json:"order,omitempty" schema:"-" url:"-"`

	// Query parameter
	Page       uint    `json:"-" gorm:"-" schema:"page" url:"page,omitempty"`                 // Query parameter, eg: "page=2"
	Size       uint    `json:"-" gorm:"-" schema:"size" url:"size,omitempty"`                 // Query parameter, eg: "size=10"
	Expand     *string `json:"-" gorm:"-" schema:"_expand" url:"_expand,omitempty"`           // Query parameter, eg: "_expand=children,parent".
	Depth      *uint   `json:"-" gorm:"-" schema:"_depth" url:"_depth,omitempty"`             // Query parameter, eg: "_depth=3".
	Fuzzy      *bool   `json:"-" gorm:"-" schema:"_fuzzy" url:"_fuzzy,omitempty"`             // Query parameter, eg: "_fuzzy=true"
	SortBy     string  `json:"-" gorm:"-" schema:"_sortby" url:"_sortby,omitempty"`           // Query parameter, eg: "_sortby=name"
	NoCache    bool    `json:"-" gorm:"-" schema:"_nocache" url:"_nocache,omitempty"`         // Query parameter: eg: "_nocache=false"
	ColumnName string  `json:"-" gorm:"-" schema:"_column_name" url:"_column_name,omitempty"` // Query parameter: eg: "_column_name=created_at"
	StartTime  string  `json:"-" gorm:"-" schema:"_start_time" url:"_start_time,omitempty"`   // Query parameter: eg: "_start_time=2024-04-29+23:59:59"
	EndTime    string  `json:"-" gorm:"-" schema:"_end_time" url:"_end_time,omitempty"`       // Query parameter: eg: "_end_time=2024-04-29+23:59:59"
	Or         *bool   `json:"-" gorm:"-" schema:"_or" url:"_or,omitempty"`                   // query parameter: eg: "_or=true"
	Index      string  `json:"-" gorm:"-" schema:"_index" url:"_index,omitempty"`             // Query parameter: eg: "_index=name"
	Select     string  `json:"-" gorm:"-" schema:"_select" url:"_select,omitempty"`           // Query parameter: eg: "_select=field1,field2"

	gorm.Model `json:"-" schema:"-" url:"-"`
}

func WithQuery(_keyValues ...any) Option {
	if len(_keyValues) == 0 || len(_keyValues) == 1 {
		return func(_ *Client) {}
	}
	keyValues := make([]string, 0, len(_keyValues))
	for i := range _keyValues {
		val := reflect.ValueOf(_keyValues[i])
		if val.Kind() == reflect.Ptr && !val.IsNil() {
			val = val.Elem()
		}
		switch val.Kind() {
		case reflect.String:
			if str := strings.TrimSpace(val.String()); len(str) > 0 {
				keyValues = append(keyValues, str)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			keyValues = append(keyValues, fmt.Sprint(val.Interface()))
		case reflect.Float32, reflect.Float64:
			keyValues = append(keyValues, fmt.Sprintf("%.6f", val.Float()))
		case reflect.Bool:
			keyValues = append(keyValues, fmt.Sprint(val.Bool()))
		}
	}

	length := len(keyValues)
	if length == 0 || length == 1 {
		return func(_ *Client) {}
	}
	if length%2 != 0 {
		length--
	}

	var queryBuilder strings.Builder
	queryBuilder.Grow(length * 8)
	for i := 0; i < length; i += 2 {
		if i > 0 {
			queryBuilder.WriteByte('&')
		}
		key := url.QueryEscape(keyValues[i])
		value := url.QueryEscape(keyValues[i+1])
		queryBuilder.WriteString(key)
		queryBuilder.WriteByte('=')
		queryBuilder.WriteString(value)
	}

	return func(c *Client) {
		c.queryRaw = queryBuilder.String()
	}
}

func WithQueryPagination(page, size uint) Option {
	return func(c *Client) {
		if c.query == nil {
			c.query = new(Query)
		}
		if page == 0 {
			page = 1
		}
		if size == 0 {
			size = 10
		}
		c.query.Page = page
		c.query.Size = size
	}
}

func WithQueryExpand(expand string, depth uint) Option {
	return func(c *Client) {
		if c.query == nil {
			c.query = new(Query)
		}
		if expand = strings.TrimSpace(expand); len(expand) == 0 {
			return
		}
		c.query.Expand = &expand
		c.query.Depth = &depth
	}
}

func WithQueryFuzzy(fuzzy bool) Option {
	return func(c *Client) {
		if c.query == nil {
			c.query = new(Query)
		}
		c.query.Fuzzy = &fuzzy
	}
}

func WithQuerySortby(sortby string) Option {
	return func(c *Client) {
		if sortby = strings.TrimSpace(sortby); len(sortby) == 0 {
			return
		}
		if c.query == nil {
			c.query = new(Query)
		}
		c.query.SortBy = sortby
	}
}

func WithQueryNocache(nocache bool) Option {
	return func(c *Client) {
		if c.query == nil {
			c.query = new(Query)
		}
		c.query.NoCache = nocache
	}
}

func WithQueryTimeRange(columeName string, start, end time.Time) Option {
	return func(c *Client) {
		if c.query == nil {
			c.query = new(Query)
		}
		if columeName = strings.TrimSpace(columeName); len(columeName) == 0 {
			return
		}
		if start.IsZero() || end.IsZero() {
			return
		}
		if start.After(end) {
			start, end = end, start
		}
		c.query.ColumnName = columeName
		c.query.StartTime = start.Format(consts.DATE_TIME_LAYOUT)
		c.query.EndTime = end.Format(consts.DATE_TIME_LAYOUT)
	}
}

func WithQueryOr(or bool) Option {
	return func(c *Client) {
		if c.query == nil {
			c.query = new(Query)
		}
		c.query.Or = &or
	}
}

func WithQueryIndex(index string) Option {
	return func(c *Client) {
		if index = strings.TrimSpace(index); len(index) == 0 {
			return
		}
		if c.query == nil {
			c.query = new(Query)
		}
		c.query.Index = index
	}
}

func WithQuerySelect(selects ...string) Option {
	return func(c *Client) {
		_selects := make([]string, 0, len(selects))
		for i := range selects {
			if len(strings.TrimSpace(selects[i])) != 0 {
				_selects = append(_selects, strings.TrimSpace(selects[i]))
			}
		}
		if len(_selects) == 0 {
			return
		}
		if c.query == nil {
			c.query = new(Query)
		}
		c.query.Select = strings.Join(_selects, ",")
	}
}
