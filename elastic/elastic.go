package elastic

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/logger"
	"go.uber.org/zap"
)

const (
	TIMESTAMP_FILE    = "/tmp/timestamp"
	defaultSearchSize = 10000
)

type (
	document struct{}
	index    struct{}
)

var (
	client   *elasticsearch.Client
	esConfig elasticsearch.Config
	timeout  = 10 * time.Second

	Document = new(document)
	Index    = new(index)
)

// Init initializes the global elasticsearch client.
// It reads elasticsearch configuration from config.App.ElasticsearchConfig.
// If elasticsearch not enabled, it returns nil.
// The functions also starts a background goroutines to ensure connection health.
func Init() (err error) {
	cfg := config.App.ElasticsearchConfig
	if !cfg.Enable {
		return nil
	}
	esConfig = makeESConfig(cfg)
	if client, err = elasticsearch.NewClient(esConfig); err != nil {
		return err
	}
	ticker := time.NewTicker(timeout + 10*time.Second)
	go func() {
		for range ticker.C {
			_ensureConnection()
		}
	}()
	zap.S().Infow("successfully connect to elasticsearch", "hosts", cfg.Hosts)
	return nil
}

// New creates a new elasticsearch client instance with the given configuration.
func New(cfg config.ElasticsearchConfig) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(makeESConfig(cfg))
}

func makeESConfig(cfg config.ElasticsearchConfig) elasticsearch.Config {
	return elasticsearch.Config{
		Addresses:              strings.Split(cfg.Hosts, ","),
		Username:               cfg.Username,
		Password:               cfg.Password,
		CloudID:                cfg.CloudID,
		APIKey:                 cfg.APIKey,
		ServiceToken:           cfg.ServiceToken,
		CertificateFingerprint: cfg.CertificateFingerprint,
		Transport:              &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
}

// _ensureConnection checks the connection and reconnects if necessary
func _ensureConnection() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	logger.Elastic.Info("check elasticsearch connection")
	defer cancel()
	if _, err := client.Ping(client.Ping.WithContext(ctx)); err != nil {
		logger.Elastic.Warn("elasticsearch connection maybe broken, try to reconnect: %v", err)
		if newClient, err := elasticsearch.NewClient(esConfig); err != nil {
			logger.Elastic.Error("reconnect to elasticsearch error: %v", err)
		} else {
			client = newClient
		}
	}
}

// _check will check the client and return an error if it's nil or invalid.
func _check() error {
	if client == nil || client == new(elasticsearch.Client) {
		return errors.New("elasticsearch client is nil")
	}
	return nil
}

func Client() *elasticsearch.Client { return client }

// SearchTimestamp
// WithIndex("winlog*"), // 使用通配符查询多个索引
func SearchTimestamp(index string, size ...int) ([]byte, error) {
	_size := defaultSearchSize
	if len(size) > 0 {
		if size[0] > 0 {
			_size = size[0]
		}
	}

	// query := `
	// {
	//   "size":3000,
	//   "sort": [{"@timestamp":{"order":"asc"}}],
	//   "query": {
	//     "bool": {
	//       "must": [
	// 			{"range":{"@timestamp":{"gte":"2024-02-27T01:00:00Z","lte":"2024-02-27T23:00:00Z"}}},
	// 			{"term":{"event.action.keyword":"Logon"}}
	//     ]
	//     }
	//   }
	// }`

	var (
		err                error
		timestampEnd       time.Time
		timestampStart     time.Time
		timestampStartData []byte
	)

	now := time.Now().UTC()
	timestampEnd = time.Date(2099, now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
	if timestampStartData, err = os.ReadFile(TIMESTAMP_FILE); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			timestampStart = time.Date(now.Year(), now.Month(), now.Day()-1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC) // 提前一天
			fmt.Println("------------------- touch file and write time: ", timestampStart.Format(time.RFC3339))
			if err = os.WriteFile(TIMESTAMP_FILE, []byte(timestampStart.Format(time.RFC3339)), 0o644); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if timestampStart, err = time.ParseInLocation(time.RFC3339, string(bytes.TrimSpace(timestampStartData)), time.UTC); err != nil {
			return nil, err
		}
	}

	queryFormat := `{
  "size":%d,
  "sort": [{"@timestamp":{"order":"asc"}}],
  "query": {
	"bool": {
	  "must": [
			{"range":{"@timestamp":{"gte":"%s","lte":"%s"}}},
			{"wildcard":{"event.action.keyword": "*"}}
	]
	}
  }
}`

	query := fmt.Sprintf(queryFormat, _size, timestampStart.Format(time.RFC3339), timestampEnd.Format(time.RFC3339))
	// fmt.Println(query)

	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		// client.Search.WithIndex("winlog*"), // 使用通配符查询多个索引
		client.Search.WithIndex(index), // 使用通配符查询多个索引
		client.Search.WithBody(strings.NewReader(query)),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)
	if err != nil {
		panic(err)
	}
	if res.IsError() {
		fmt.Println("------------------- error")
		fmt.Println(res.String())
		return nil, errors.New(res.Status())
	}

	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

type Pagination struct {
	Page int // page number
	Size int // page size
}

func handleError(fn func() (*esapi.Response, error), format string) error {
	res, err := fn()
	if err != nil {
		return fmt.Errorf(format, err.Error())
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf(format, res.String())
	}
	return nil
}

/*
操作	方法	描述
索引文档	es.Index()	创建或替换文档
更新文档	es.Update()	部分更新已存在的文档
获取文档	es.Get()	检索特定文档
删除文档	es.Delete()	从索引中删除文档
批量操作	es.Bulk()	在单个请求中执行多个索引/更新/删除操作
搜索	es.Search()	在一个或多个索引中搜索文档
创建索引	es.Indices.Create()	创建新的索引
删除索引	es.Indices.Delete()	删除一个或多个索引
索引别名	es.Indices.PutAlias()	为索引创建或更新别名
刷新索引	es.Indices.Refresh()	刷新一个或多个索引
获取映射	es.Indices.GetMapping()	获取一个或多个索引的映射
更新映射	es.Indices.PutMapping()	更新一个或多个索引的映射


es.Index(): 用于创建新文档或替换现有文档。如果文档不存在，它会被创建；如果存在，则会被完全替换。
es.Update(): 用于部分更新已存在的文档。你可以添加、修改或删除文档中的特定字段，而不影响其他字段。
es.Get(): 通过索引名和文档ID检索特定文档。可以获取整个文档或指定字段。
es.Delete(): 从指定索引中删除特定文档。
es.Bulk(): 允许在单个API调用中执行多个操作，如批量索引、更新或删除文档，提高效率。
es.Search(): 执行搜索查询，可以在一个或多个索引中搜索文档。支持各种查询类型和聚合。
es.Indices.Create(): 创建新的索引，可以指定设置和映射。
es.Indices.Delete(): 删除一个或多个索引及其所有数据。
es.Indices.PutAlias(): 为一个或多个索引创建或更新别名，便于索引的逻辑分组或无缝切换。
es.Indices.Refresh(): 刷新索引，使最近的更改对搜索可见。
es.Indices.GetMapping(): 获取一个或多个索引的映射信息，包括字段类型和索引选项。
es.Indices.PutMapping(): 更新一个或多个索引的映射，允许添加新字段或修改现有字段的映射。
*/
