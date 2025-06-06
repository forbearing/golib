package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/types"
)

func (*document) BulkIndex(_ context.Context, indexName string, docs ...types.ESDocumenter) error {
	var (
		buf  bytes.Buffer
		raw  map[string]any
		data []byte
		res  *esapi.Response
		err  error
	)
	start := time.Now()

	// 遍历消息数组
	for i := range docs {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%s" } }%s`, docs[i].GetID(), "\n"))
		if data, err = json.Marshal(docs[i].Document()); err != nil {
			err = errors.New("failed to marshaling document: " + err.Error())
			logger.Elastic.Error(err)
			return err
		}
		data = append(data, "\n"...)
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
	}

	// 执行批量请求
	res, err = client.Bulk(bytes.NewReader(buf.Bytes()), client.Bulk.WithIndex(indexName))
	if err != nil {
		err = fmt.Errorf("failed to execute bulk request: %v", err)
		logger.Elastic.Error(err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		if err = json.NewDecoder(res.Body).Decode(&raw); err != nil {
			err = fmt.Errorf("failed to parse response body: %v", err)
			logger.Elastic.Error(err)
			return err
		}
		err = fmt.Errorf("failed to execute bulk request: %v", raw)
		logger.Elastic.Error(err)
		return err
	}

	var blk map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&blk); err != nil {
		err = fmt.Errorf("failed to parse response body: %v", err)
		logger.Elastic.Error(err)
		return err
	}
	if blk["errors"].(bool) {
		for _, item := range blk["items"].([]interface{}) {
			if idx, ok := item.(map[string]interface{})["index"].(map[string]interface{}); ok {
				if idx["error"] != nil {
					err = fmt.Errorf("error in item: %v", idx["error"])
					logger.Elastic.Error(err)
					return err
				}
			}
		}
	}
	logger.Elastic.Infow("successfully indexed documents", "length", len(docs), "cost", time.Since(start).String())
	return nil
}
