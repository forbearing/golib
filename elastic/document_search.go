package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// Search performs a simple Search on a single field.
//
// Parameters:
//   - indexName: The name of the index to Search.
//   - field: The field to Search in.
//   - value: The value to Search for.
//
// Returns:
//   - []map[string]any: A slice of documents matching the Search criteria.
//   - error: An error object which will be non-nil if any error occurs during the process.
func (*document) Search(indexName, field, value string, pagination ...*Pagination) ([]map[string]any, int64, error) {
	if err := _check(); err != nil {
		return nil, 0, err
	}

	from, size := _paginationValue(pagination...)
	query := map[string]any{
		"query": map[string]any{
			"match": map[string]any{
				field: value,
			},
		},
	}

	return _searchWithPagination(indexName, query, from, size)
}

// SearchMulti performs a search across multiple fields.
//
// Parameters:
//   - indexName: The name of the index to search.
//   - queries: A map where keys are field names and values are search terms.
//
// Returns:
//   - []map[string]any: A slice of documents matching the search criteria.
//   - error: An error object which will be non-nil if any error occurs during the process.
func (*document) SearchMulti(indexName string, queries map[string]string, pagination ...*Pagination) ([]map[string]any, int64, error) {
	if err := _check(); err != nil {
		return nil, 0, err
	}

	from, size := _paginationValue(pagination...)
	shouldClauses := make([]map[string]any, 0, len(queries))
	for field, value := range queries {
		shouldClauses = append(shouldClauses, map[string]any{
			"match": map[string]any{
				field: value,
			},
		})
	}
	query := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"should": shouldClauses,
			},
		},
	}

	return _searchWithPagination(indexName, query, from, size)
}

// SearchRange performs a range query and returns paginated results.
//
// Parameters:
//   - indexName: The name of the Elasticsearch index to search.
//   - field: The name of the field to perform the range query on.
//   - gte: The value for "greater than or equal to". Can be any type. If nil, no lower bound is set.
//   - lte: The value for "less than or equal to". Can be any type. If nil, no upper bound is set.
//   - pagination: Optional pagination parameters. If provided, should be a pointer to a Pagination struct.
//     If not provided or nil, default pagination settings will be used (page 1, 1000 items per page).
//
// Returns:
//   - []map[string]any: A slice containing the query results, where each result is a map.
//   - int64: The total number of documents matching the query criteria.
//   - error: An error if one occurred, otherwise nil.
//
// Description:
//
//	This function executes a range query on the specified Elasticsearch index. It allows users to specify
//	a field and a range of values for that field (greater than or equal to and/or less than or equal to).
//	The function supports pagination, which can be controlled via the optional Pagination parameter.
//	If no pagination parameter is provided, default settings will be used.
//
// Example:
//
//	results, total, err := SearchRange("my_index", "age", 18, 30, &Pagination{Page: 1, Size: 20})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Found %d results, total %d matches\n", len(results), total)
//
// NOTE:
//   - If both gte and lte are nil, it will return all non-null values for the specified field.
//   - This function uses the From/Size pagination method, which may have performance issues for large datasets.
//   - if field type is time/date, format must be RFC3339
func (*document) SearchRange(indexName, field string, gte, lte any, pagination ...*Pagination) ([]map[string]any, int64, error) {
	if err := _check(); err != nil {
		return nil, 0, err
	}

	from, size := _paginationValue(pagination...)
	rangeQuery := make(map[string]any)
	if gte != nil {
		rangeQuery["gte"] = gte
	}
	if lte != nil {
		rangeQuery["lte"] = lte
	}
	query := map[string]any{
		"query": map[string]any{
			"range": map[string]any{
				field: rangeQuery,
			},
		},
	}

	return _searchWithPagination(indexName, query, from, size)
}

// SearchRaw performs a raw Elasticsearch query and returns paginated results.
//
// Parameters:
//   - indexName: The name of the Elasticsearch index to search.
//   - query: A map[string]any representing the raw Elasticsearch query.
//   - pagination: Optional pagination parameters. If provided, should be a pointer to a Pagination struct.
//     If not provided or nil, default pagination settings will be used (page 1, 1000 items per page).
//
// Returns:
//   - []map[string]any: A slice containing the query results, where each result is a map.
//   - int64: The total number of documents matching the query criteria.
//   - error: An error if one occurred, otherwise nil.
//
// Description:
//
//	This function executes a raw Elasticsearch query on the specified index. It allows users to provide
//	any valid Elasticsearch query structure. The function supports pagination, which can be controlled
//	via the optional Pagination parameter. If no pagination parameter is provided, default settings will be used.
//
// Example:
//
//	query := map[string]any{
//	    "query": map[string]any{
//	        "match": map[string]any{
//	            "title": "elasticsearch",
//	        },
//	    },
//	}
//	results, total, err := SearchRaw("my_index", query, &Pagination{Page: 1, Size: 20})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Found %d results, total %d matches\n", len(results), total)
//
// NOTE:
//   - This function uses the From/Size pagination method, which may have performance issues for large datasets.
//   - The provided query should be a valid Elasticsearch query structure.
func (*document) SearchRaw(indexName string, query map[string]any, pagination ...*Pagination) ([]map[string]any, int64, error) {
	if err := _check(); err != nil {
		return nil, 0, err
	}
	from, size := _paginationValue(pagination...)
	return _searchWithPagination(indexName, query, from, size)
}

func _paginationValue(pagination ...*Pagination) (from int, size int) {
	pageSize := 1000
	pageNum := 1
	from = (pageNum - 1) * pageSize
	if len(pagination) > 0 {
		if pagination[0] != nil {
			pageSize = pagination[0].Size
			pageNum = pagination[0].Page
			from = (pageNum - 1) * pageSize
		}
	}
	return from, pageSize
}

// _search is a helper function to execute the _search request.
func _search(indexName string, query map[string]any) ([]map[string]any, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(indexName),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error executing search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search request failed: %s", res.Status())
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	hits, _ := result["hits"].(map[string]any)["hits"].([]any)
	docs := make([]map[string]any, len(hits))

	for i, hit := range hits {
		source, _ := hit.(map[string]any)["_source"].(map[string]any)
		docs[i] = source
	}

	return docs, nil
}

// _searchWithPagination is a helper function to execute the _search request.
func _searchWithPagination(indexName string, query map[string]any, from, size int) ([]map[string]any, int64, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(indexName),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithFrom(from),
		client.Search.WithSize(size),
		client.Search.WithPretty(),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("error executing search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search request failed: %s", res.Status())
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("error parsing the response body: %w", err)
	}

	hitsMap, _ := result["hits"].(map[string]any)
	hits, _ := hitsMap["hits"].([]any)
	totalHits, _ := hitsMap["total"].(map[string]any)["value"].(float64)

	docs := make([]map[string]any, len(hits))

	for i, hit := range hits {
		source, _ := hit.(map[string]any)["_source"].(map[string]any)
		docs[i] = source
	}

	return docs, int64(totalHits), nil
}
