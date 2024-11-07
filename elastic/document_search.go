package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/forbearing/golib/logger"
)

const (
	DefaultTimeFormat = "strict_date_optional_time||epoch_millis"
	DefaultTimeField  = "created_at"
	DefaultSortOrder  = "desc"
)

// SearchRequest represents the search parameters
type SearchRequest struct {
	Query       map[string]any   `json:"query,omitempty"`
	From        int              `json:"from,omitempty"`
	Size        int              `json:"size,omitempty"`
	Sort        []map[string]any `json:"sort,omitempty"`
	Source      []string         `json:"_source,omitempty"`
	SearchAfter []any            `json:"search_after,omitempty"`
}

// SearchResult represents the search response
type SearchResult struct {
	Total    int64       `json:"total"`
	MaxScore *float64    `json:"max_score,omitempty"`
	Hits     []SearchHit `json:"hits"`
}

// SearchHit represents a single hit in search results
type SearchHit struct {
	ID     string         `json:"_id"`
	Score  *float64       `json:"_score,omitempty"`
	Source map[string]any `json:"_source"`
}

// Search performs a search operation on the specified index
// default size is 10
func (*document) Search(ctx context.Context, indexName string, req *SearchRequest) (*SearchResult, error) {
	if err := _check(); err != nil {
		return nil, fmt.Errorf("elasticsearch client check: %w", err)
	}
	if indexName == "" {
		return nil, fmt.Errorf("index name cannot be empty")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	// Set default values if not provided
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.From < 0 {
		req.From = 0
	}
	// SearchAfter require Sort to be set.
	if len(req.Sort) == 0 {
		req.Sort = []map[string]any{
			{
				"_doc": map[string]any{ // 使用 _doc 作为第二排序字段
					"order": "asc",
				},
			},
		}
	}

	begin := time.Now()
	logger := logger.Elastic.With(
		"index", indexName,
		"from", strconv.Itoa(req.From),
		"size", strconv.Itoa(req.Size),
	)
	defer func() {
		logger.Infow("search completed", "cost", time.Since(begin).String())
	}()

	// Convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search request: %w", err)
	}

	// Perform search request
	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(indexName),
		client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		logger.Errorw("failed to execute search", "error", err)
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		logger.Errorw("elasticsearch error response",
			"status", res.Status(),
			"body", string(body),
		)
		return nil, fmt.Errorf("elasticsearch error [%s]: %s", res.Status(), string(body))
	}
	var esRes map[string]any
	if err := json.NewDecoder(res.Body).Decode(&esRes); err != nil {
		logger.Errorw("failed to decode response", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return parseSearchResult(esRes)
}

// parseSearchResult
// 辅助函数：解析搜索结果
func parseSearchResult(esRes map[string]any) (*SearchResult, error) {
	var (
		ok       bool
		id       string
		err      error
		total    int64
		maxScore float64
		hits     map[string]any
		hitMap   map[string]any
		source   map[string]any
		hitsList []any
	)

	// Extract hits with type assertion safety
	if hits, ok = esRes["hits"].(map[string]any); !ok {
		return nil, fmt.Errorf("invalid response format: hits not found or invalid type")
	}
	// Process search result with safe type assertions
	if total, err = extractTotal(hits); err != nil {
		return nil, fmt.Errorf("failed to extract total: %w", err)
	}
	result := &SearchResult{Total: total}
	// Safely extract max_score
	if maxScore, ok = hits["max_score"].(float64); ok {
		result.MaxScore = &maxScore
	}
	// Process hits with safe type assertions
	if hitsList, ok = hits["hits"].([]any); !ok {
		return nil, fmt.Errorf("invalid response format: hits list not found or invalid type")
	}
	result.Hits = make([]SearchHit, len(hitsList))

	for i, hit := range hitsList {
		if hitMap, ok = hit.(map[string]any); !ok {
			return nil, fmt.Errorf("invalid hit format at index %d", i)
		}
		if id, ok = hitMap["_id"].(string); !ok {
			return nil, fmt.Errorf("invalid or missing _id at index %d", i)
		}
		if source, ok = hitMap["_source"].(map[string]any); !ok {
			return nil, fmt.Errorf("invalid or missing _source at index %d", i)
		}
		var score *float64
		if scoreVal, ok := hitMap["_score"].(float64); ok {
			score = &scoreVal
		}
		result.Hits[i] = SearchHit{
			ID:     id,
			Source: source,
			Score:  score,
		}
	}
	return result, nil
}

// 辅助函数：安全地提取总数
func extractTotal(hits map[string]any) (int64, error) {
	total, ok := hits["total"].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("invalid total format")
	}

	value, ok := total["value"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid total value format")
	}

	return int64(value), nil
}
