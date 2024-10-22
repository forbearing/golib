package elastic

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/forbearing/golib/logger"
)

// Get retrieves a document from Elasticsearch by its ID.
//
// This function sends a GET request to Elasticsearch to fetch a document
// with the specified ID from the given index.
//
// Parameters:
//   - indexName: The name of the index containing the document.
//   - id: The unique identifier of the document to retrieve.
//
// Returns:
//   - map[string]any: The _source content of the retrieved document.
//   - error: An error object which will be non-nil if any error occurs during the process.
//
// The function handles the following scenarios:
//  1. Successfully retrieves the document and returns its _source content.
//  2. Returns an error if the document is not found.
//  3. Returns an error for any Elasticsearch query execution issues.
//  4. Returns an error if there are problems parsing the response.
//
// Usage example:
//
//	doc, err := Get("my_index", "doc123")
//	if err != nil {
//	    log.Fatalf("Failed to get document: %v", err)
//	}
//	fmt.Printf("Retrieved document: %v\n", doc)
//
// Note: This function assumes that the Elasticsearch client is properly configured and connected.
func (*document) Get(indexName string, id string) (map[string]any, error) {
	if err := _check(); err != nil {
		return nil, err
	}

	begin := time.Now()
	defer func() {
		logger.Elastic.Infow("get document", "index", indexName, "id", id, "cost", time.Since(begin).String())
	}()

	// 执行 Get 请求
	res, err := client.Get(indexName, id)
	if err != nil {
		err = fmt.Errorf("error getting document: %s", err)
		logger.Elastic.Error(err)
		return nil, err
	}
	defer res.Body.Close()

	// 检查响应状态
	if res.IsError() {
		err = fmt.Errorf("error getting document, status: %s", res.Status())
		logger.Elastic.Error(err)
		return nil, err
	}

	// 解析响应体
	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		err = fmt.Errorf("error parsing the response body: %s", err)
		logger.Elastic.Error(err)
		return nil, err
	}

	// 返回文档源数据
	if source, found := result["_source"].(map[string]any); found {
		return source, nil
	}
	return nil, fmt.Errorf("_source not found in the response")
}
