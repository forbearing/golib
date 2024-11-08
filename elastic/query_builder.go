package elastic

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	defaultSize       = 10
	defaultFrom       = 0
	defaultTimeFormat = "strict_date_optional_time||epoch_millis"

	Asc  Order = "asc"
	Desc Order = "desc"
)

type Order string

// QueryBuilder 用于构建 Elasticsearch 查询
// 支持 must, must_not, should, filter 等查询条件
// 支持分页、排序、字段过滤和 search_after
type QueryBuilder struct {
	must               []map[string]any
	mustNot            []map[string]any
	should             []map[string]any
	filter             []map[string]any
	size               int
	from               int
	sort               []map[string]any
	source             []string
	searchAfter        []any
	minimumShouldMatch any
}

// NewQueryBuilder 创建一个新的查询构建器
// 默认 size=10, from=0
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		size: defaultSize,
		from: defaultFrom,
	}
}

// QueryBuilder 添加新方法
func (qb *QueryBuilder) Bool(fn func(builder *QueryBuilder)) *QueryBuilder {
	nestedBuilder := NewQueryBuilder()
	fn(nestedBuilder)
	query := nestedBuilder.BuildQuery()
	if query != nil { // 只有在query不为nil时才添加
		qb.Must(query)
	}
	return qb
}

func (qb *QueryBuilder) Must(query map[string]any) *QueryBuilder {
	qb.must = append(qb.must, query)
	return qb
}

func (qb *QueryBuilder) MustNot(query map[string]any) *QueryBuilder {
	qb.mustNot = append(qb.mustNot, query)
	return qb
}

func (qb *QueryBuilder) Should(query map[string]any) *QueryBuilder {
	qb.should = append(qb.should, query)
	return qb
}

func (qb *QueryBuilder) MinimumShouldMatch(minimum any) *QueryBuilder {
	qb.minimumShouldMatch = minimum
	return qb
}

func (qb *QueryBuilder) Filter(query map[string]any) *QueryBuilder {
	qb.filter = append(qb.filter, query)
	return qb
}

func (qb *QueryBuilder) Term(field string, value any) *QueryBuilder {
	return qb.Must(map[string]any{
		"term": map[string]any{
			field: value,
		},
	})
}

func (qb *QueryBuilder) TermNot(field string, value any) *QueryBuilder {
	return qb.MustNot(map[string]any{
		"term": map[string]any{
			field: value,
		},
	})
}

func (qb *QueryBuilder) TermShould(field string, value any) *QueryBuilder {
	return qb.Should(map[string]any{
		"term": map[string]any{
			field: value,
		},
	})
}

func (qb *QueryBuilder) Match(field string, value any) *QueryBuilder {
	return qb.Must(map[string]any{
		"match": map[string]any{
			field: value,
		},
	})
}

func (qb *QueryBuilder) MatchNot(field string, value any) *QueryBuilder {
	return qb.MustNot(map[string]any{
		"match": map[string]any{
			field: value,
		},
	})
}

func (qb *QueryBuilder) MatchShould(field string, value any) *QueryBuilder {
	return qb.Should(map[string]any{
		"match": map[string]any{
			field: value,
		},
	})
}

func (qb *QueryBuilder) MatchPhrase(field string, value any) *QueryBuilder {
	return qb.Must(map[string]any{
		"match_phrase": map[string]any{
			field: value,
		},
	})
}

func (qb *QueryBuilder) MatchPhraseNot(field string, value any) *QueryBuilder {
	return qb.MustNot(map[string]any{
		"match_phrase": map[string]any{
			field: value,
		},
	})
}

func (qb *QueryBuilder) MatchPhraseShould(field string, value any) *QueryBuilder {
	return qb.Should(map[string]any{
		"match_phrase": map[string]any{
			field: value,
		},
	})
}

// MatchAll 添加匹配所有文档的查询
func (qb *QueryBuilder) MatchAll() *QueryBuilder {
	return qb.Must(map[string]any{
		"match_all": map[string]any{},
	})
}

func (qb *QueryBuilder) Range(field string, ranges map[string]any) *QueryBuilder {
	return qb.Filter(map[string]any{
		"range": map[string]any{
			field: ranges,
		},
	})
}

func (qb *QueryBuilder) Exists(field string) *QueryBuilder {
	return qb.Must(map[string]any{
		"exists": map[string]any{
			"field": field,
		},
	})
}

func (qb *QueryBuilder) TimeRange(field string, start, end time.Time) *QueryBuilder {
	return qb.Range(field, map[string]any{
		"gte":    start.Format(time.RFC3339),
		"lte":    end.Format(time.RFC3339),
		"format": defaultTimeFormat,
	})
}

func (qb *QueryBuilder) TimeGte(field string, tm time.Time) *QueryBuilder {
	return qb.Range(field, map[string]any{
		"gte":    tm.Format(time.RFC3339),
		"format": defaultTimeFormat,
	})
}

func (qb *QueryBuilder) TimeLte(field string, tm time.Time) *QueryBuilder {
	return qb.Range(field, map[string]any{
		"lte":    tm.Format(time.RFC3339),
		"format": defaultTimeFormat,
	})
}

func (qb *QueryBuilder) Size(size int) *QueryBuilder {
	if size > 0 {
		qb.size = size
	}
	return qb
}

func (qb *QueryBuilder) From(from int) *QueryBuilder {
	if from >= 0 {
		qb.from = from
	}
	return qb
}

func (qb *QueryBuilder) Sort(field string, order Order) *QueryBuilder {
	if field == "" {
		return qb
	}
	if order != Asc && order != Desc {
		order = Desc
	}
	qb.sort = append(qb.sort, map[string]any{
		field: map[string]any{
			"order": order,
		},
	})
	return qb
}

func (qb *QueryBuilder) Validate() error {
	if qb.size < 0 {
		return fmt.Errorf("size cannot be negative")
	}
	if qb.from < 0 {
		return fmt.Errorf("from cannot be negative")
	}

	// 如果使用了 search_after，from 必须为 0
	if len(qb.searchAfter) > 0 && qb.from != 0 {
		return fmt.Errorf("from must be 0 when using search_after")
	}

	return nil
}

// Source 设置返回字段
// 不设置则返回全部字段
// 设置空数组则不返回任何字段
func (qb *QueryBuilder) Source(fields ...string) *QueryBuilder {
	qb.source = fields
	return qb
}

// SearchAfter 设置 search_after 值
func (qb *QueryBuilder) SearchAfter(value ...any) *QueryBuilder {
	qb.searchAfter = value
	return qb
}

func (qb *QueryBuilder) Build() (*SearchRequest, error) {
	if err := qb.Validate(); err != nil {
		return nil, err
	}

	if len(qb.must) == 0 && len(qb.mustNot) == 0 &&
		len(qb.should) == 0 && len(qb.filter) == 0 {
		return &SearchRequest{
			Query: map[string]any{
				"match_all": map[string]any{},
			},
			From:        qb.from,
			Size:        qb.size,
			Sort:        qb.sort,
			Source:      qb.source,
			SearchAfter: qb.searchAfter,
		}, nil
	}

	boolQuery := make(map[string]any)
	if len(qb.must) > 0 {
		boolQuery["must"] = qb.must
	}
	if len(qb.mustNot) > 0 {
		boolQuery["must_not"] = qb.mustNot
	}
	if len(qb.should) > 0 {
		boolQuery["should"] = qb.should
	}
	if len(qb.should) > 0 && qb.minimumShouldMatch != nil {
		boolQuery["minimum_should_match"] = qb.minimumShouldMatch
	}
	if len(qb.filter) > 0 {
		boolQuery["filter"] = qb.filter
	}

	return &SearchRequest{
		Query:       map[string]any{"bool": boolQuery},
		From:        qb.from,
		Size:        qb.size,
		Sort:        qb.sort,
		Source:      qb.source,
		SearchAfter: qb.searchAfter,
	}, nil
}

func (qb *QueryBuilder) BuildQuery() map[string]any {
	req, err := qb.Build()
	if err != nil || req.Query == nil {
		return nil // 构建失败返回 nil
	}
	return req.Query
}

func (qb *QueryBuilder) String() string {
	req, err := qb.Build()
	if err != nil {
		return fmt.Sprintf("invalid query: %v", err)
	}

	bytes, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return fmt.Sprintf("failed to marshal query: %v", err)
	}

	return string(bytes)
}
