package elastic_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/araddon/dateparse"
	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/elastic"
	"github.com/forbearing/golib/logger/zap"
	"github.com/forbearing/golib/util"
	"github.com/stretchr/testify/assert"
)

const Index = "test"

func init() {
	config.SetConfigFile("../../config.ini")
	bootstrap.Register(
		config.Init,
		zap.Init,
		elastic.Init,
	)
	if err := bootstrap.Init(); err != nil {
		panic(err)
	}
}

func TestIndex(t *testing.T) {
	settings := map[string]interface{}{
		"number_of_shards":   3,
		"number_of_replicas": 2,
	}

	mappings := map[string]interface{}{
		"properties": map[string]interface{}{
			"title": map[string]interface{}{
				"type": "text",
			},
			"content": map[string]interface{}{
				"type": "text",
			},
			"date": map[string]interface{}{
				"type": "date",
			},
		},
	}
	idx1 := "my_index"
	idx2 := "my_index2"

	assert.NoError(t, elastic.Index.Create(idx1, &elastic.IndexOption{settings, mappings}))
	assert.NoError(t, elastic.Index.Create(idx2))
	exists, err := elastic.Index.Exists(idx1)
	assert.Equal(t, true, exists)
	assert.NoError(t, err)

	exists, err = elastic.Index.Exists(idx2)
	assert.Equal(t, true, exists)
	assert.NoError(t, err)

	assert.NoError(t, elastic.Index.Delete(idx1))
	assert.NoError(t, elastic.Index.Delete(idx2))
}

func TestDocumentGet(t *testing.T) {
	doc, err := elastic.Document.Get(context.TODO(), Index, "message_recv_7143038995084115996_7274598307442327556_7424788642731753476", nil)
	assert.NoError(t, err)
	fmt.Println(doc)
}

func TestDocumentSearch(t *testing.T) {
	req := &elastic.SearchRequest{
		Query: map[string]any{
			"bool": map[string]any{
				"must": []map[string]any{
					{
						"match": map[string]any{
							"message_text": "hello",
						},
					},
					{
						"term": map[string]any{
							"type": "message_send",
						},
					},
				},
			},
		},
		Size: 10,
	}

	// (type.keyword : "message_send" or type.keyword : "message_recv" or type.keyword : "message_ack" ) and message_user_id.keyword : "7143038995084115996" and message_text: hello
	req = &elastic.SearchRequest{
		Query: map[string]any{
			"bool": map[string]any{
				"must": []map[string]any{
					{
						"bool": map[string]any{
							"should": []map[string]any{
								{
									"term": map[string]any{
										"type.keyword": "message_send",
									},
								},
								{
									"term": map[string]any{
										"type.keyword": "message_recv",
									},
								},
								{
									"term": map[string]any{
										"type.keyword": "message_ack",
									},
								},
							},
							"minimum_should_match": 1,
						},
					},
					{
						"term": map[string]any{
							"message_user_id.keyword": "7143038995084115996",
						},
					},
					{
						"match": map[string]any{
							"message_text": "hello",
						},
					},
				},
			},
		},
		Sort: []map[string]any{
			{
				"created_at": map[string]any{
					"order": "desc",
				},
			},
		},
	}
	req = &elastic.SearchRequest{
		Size: 1,
		Sort: []map[string]any{
			{
				"created_at": map[string]any{
					"order": "desc",
				},
			},
		},
	}

	resp, err := elastic.Document.Search(context.Background(), Index, req)
	assert.NoError(t, err)
	fmt.Println(resp.Total, len(resp.Hits))

	formatHit := func(hit elastic.SearchHit) {
		for k, v := range hit.Source {
			if k == "message_text" {
				fmt.Println("message_text:", v)
			}
		}
	}
	for _, hit := range resp.Hits {
		formatHit(hit)
	}
	for len(resp.Hits) > 0 {
		{
			lastHit := resp.Hits[len(resp.Hits)-1]
			req.SearchAfter = []any{lastHit.Source["created_at"]}
			// 如果是按照ID进行排序，则传入ID值
			// 如果是安装 created_at 进行排序，则传入 created_at 值

			resp, err = elastic.Document.Search(context.Background(), Index, req)
			assert.NoError(t, err)
			for i := range resp.Hits {
				formatHit(resp.Hits[i])
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func TestDocumentSearchNormal(t *testing.T) {
	keyword := "hello"
	size := 2
	// 普通搜索
	req := &elastic.SearchRequest{
		Query: map[string]any{
			"bool": map[string]any{
				"must": []map[string]any{
					{
						"match": map[string]any{
							"message_text": keyword,
						},
					},
				},
			},
		},
		Size: size,
	}

	// 执行搜索

	result, err := elastic.Document.Search(context.TODO(), Index, req)
	assert.NoError(t, err)
	// 打印搜索结果
	for _, hit := range result.Hits {
		fmt.Println(hit.ID)
	}
}

func TestDocumentSearchTimeRange(t *testing.T) {
	size := 2
	startTime := time.Now().Add(-24 * 30 * time.Hour)
	endTime := time.Now()

	req := &elastic.SearchRequest{
		Query: map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{
						"range": map[string]any{
							"created_at": map[string]any{
								"gte":    startTime.Format(time.RFC3339),
								"lte":    endTime.Format(time.RFC3339),
								"format": "strict_date_optional_time||epoch_millis",
							},
						},
					},
				},
			},
		},
		Size: size,
		Sort: []map[string]any{
			{
				"created_at": map[string]any{
					"order": "desc",
				},
			},
		},
	}

	result, err := elastic.Document.Search(context.TODO(), Index, req)
	assert.NoError(t, err)
	for _, hit := range result.Hits {
		fmt.Println(hit.ID)
	}

	query := elastic.NewQueryBuilder().
		TimeRange("created_at", startTime, endTime).
		Size(size).
		Sort("created_at", elastic.Desc)
	req2, err := query.Build()
	assert.NoError(t, err)

	result2, err := elastic.Document.Search(context.TODO(), Index, req2)
	assert.NoError(t, err)
	for _, hit := range result2.Hits {
		fmt.Println(hit.ID)
	}
}

func TestDocumentSearchAfter(t *testing.T) {
	keyword := "hello"
	size := 2
	dateStr := "2024-10-29T10:38:06.085+08:00"
	date, _ := dateparse.ParseLocal(dateStr)
	// search after 分页
	req := &elastic.SearchRequest{
		Query: map[string]any{
			"bool": map[string]any{
				"must": []map[string]any{
					{
						"match": map[string]any{
							"message_text": keyword,
						},
					},
				},
			},
		},
		Size: size,
		Sort: []map[string]any{
			{
				"created_at": map[string]any{
					"order": "desc",
				},
			},
		},
		SearchAfter: []any{date},
	}

	// 执行搜索

	result, err := elastic.Document.Search(context.TODO(), Index, req)
	assert.NoError(t, err)
	for _, hit := range result.Hits {
		fmt.Println(hit.ID)
	}

	query := elastic.NewQueryBuilder().
		Match("message_text", keyword).
		Sort("created_at", elastic.Desc).
		SearchAfter(date).
		Size(size).
		Source("created_at", "type", "chat_type")
	req2, err := query.Build()
	assert.NoError(t, err)

	result2, err := elastic.Document.Search(context.TODO(), Index, req2)
	assert.NoError(t, err)
	for _, hit := range result2.Hits {
		fmt.Println(hit)
	}
}

func TestDocumentBoolQueryBuilder(t *testing.T) {
	userId := "7336820045630406684"
	peerUserId := "7156029937089069057"
	{

		query := elastic.NewQueryBuilder().
			// 基础条件（and 部分）
			Term("chat_type.keyword", "direct").
			Term("type.keyword", "message_send").
			// 嵌套的 or 条件
			Bool(func(bq *elastic.QueryBuilder) {
				// 第一组条件 (message_user_id: A and message_peer_user_id: B)
				bq.Should(elastic.NewQueryBuilder().Term("message_user_id.keyword", userId).Term("message_peer_user_id.keyword", peerUserId).BuildQuery())
				// 第二组条件 (message_user_id: B and message_peer_user_id: A)
				bq.Should(elastic.NewQueryBuilder().Term("message_user_id.keyword", peerUserId).Term("message_peer_user_id.keyword", userId).BuildQuery())
				// 设置 should 至少匹配一个条件
				bq.MinimumShouldMatch(1)
			}).Size(1000)

		if _, err := query.Build(); err != nil {
			t.Fatal(err)
		}
		fmt.Println(query.String())
	}

	{

		query := elastic.NewQueryBuilder().
			Term("chat_type.keyword", "direct").
			Bool(func(bq *elastic.QueryBuilder) {
				bq.Should(elastic.NewQueryBuilder().Term("type.keyword", "message_send").BuildQuery())
				bq.Should(elastic.NewQueryBuilder().Term("type.keyword", "message_recv").BuildQuery())
				bq.MinimumShouldMatch(1)
			}).
			Bool(func(bq *elastic.QueryBuilder) {
				bq.Should(elastic.NewQueryBuilder().Term("message_user_id.keyword", userId).Term("message_peer_user_id.keyword", peerUserId).BuildQuery())
				bq.Should(elastic.NewQueryBuilder().Term("message_user_id.keyword", peerUserId).Term("message_peer_user_id.keyword", userId).BuildQuery())
				bq.MinimumShouldMatch(1)
			})

		if _, err := query.Build(); err != nil {
			t.Fatal(err)
		}
		fmt.Println(query.String())
	}

	{
		query := elastic.NewQueryBuilder().
			Term("chat_type.keyword", "direct").
			Bool(func(bq *elastic.QueryBuilder) {
				bq.Should(elastic.NewQueryBuilder().Term("type.keyword", "message_send").Term("message_user_id.keyword", userId).Term("message_peer_user_id.keyword", peerUserId).BuildQuery())
				bq.Should(elastic.NewQueryBuilder().Term("type.keyword", "message_recv").Term("message_user_id.keyword", userId).Term("message_peer_user_id.keyword", peerUserId).BuildQuery())
				bq.MinimumShouldMatch(1)
			})
		if _, err := query.Build(); err != nil {
			t.Fatal(err)
		}
		fmt.Println(query.String())
	}
}

func TestDocumentSearchAsc(t *testing.T) {
	dateStr := "2024-10-29T10:34:35.991+08:00" // message_text: "abcdefg"
	date, _ := dateparse.ParseAny(dateStr)
	fmt.Println("date:", date)

	query := elastic.NewQueryBuilder().
		Size(2).
		Source("message_text", "created_at").
		Sort("created_at", elastic.Asc).
		TermNot("message_text.keyword", "").
		Term("type.keyword", "message_send").
		SearchAfter(date)

	res, err := elastic.Document.Search(context.Background(), Index, query.BuildForce())
	assert.NoError(t, err)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(res)
}

func TestDocumentSearchDesc(t *testing.T) {
	dateStr := "2024-10-29T10:34:35.991+08:00" // message_text: "abcdefg"
	date, _ := dateparse.ParseAny(dateStr)
	fmt.Println("date:", date)

	query := elastic.NewQueryBuilder().
		Size(2).
		Source("message_text", "created_at").
		Sort("created_at", elastic.Desc).
		TermNot("message_text.keyword", "").
		Term("type.keyword", "message_send").
		SearchAfter(date)

	res, err := elastic.Document.Search(context.Background(), Index, query.BuildForce())
	assert.NoError(t, err)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(res)
}

func TestDocumentSearchNext(t *testing.T) {
	dateStr := "2024-10-29T10:34:35.991+08:00"
	date, _ := dateparse.ParseAny(dateStr)
	fmt.Println("date:", date)

	query := elastic.NewQueryBuilder().
		Size(2).
		Source("message_text", "created_at").
		Sort("created_at", elastic.Desc).
		TermNot("message_text.keyword", "").
		Term("type.keyword", "message_send").
		SearchAfter(date)

	// fmt.Println(query)
	res, err := elastic.Document.SearchNext(context.Background(), Index, query.BuildForce(), 3)
	assert.NoError(t, err)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(res)
}

func TestDocumentSearchPrev(t *testing.T) {
	dateStr := "2024-10-29T10:34:35.991+08:00"
	date, _ := dateparse.ParseAny(dateStr)
	fmt.Println("date:", date)

	query := elastic.NewQueryBuilder().
		Size(2).
		Source("message_text", "created_at").
		Sort("created_at", elastic.Desc).
		TermNot("message_text.keyword", "").
		Term("type.keyword", "message_send").
		SearchAfter(date)

	fmt.Println(query)
	res, err := elastic.Document.SearchPrev(context.Background(), Index, query.BuildForce(), 4)
	assert.NoError(t, err)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(res)
}

func TestDocumentQueryBuilderMatchPharseOptions(t *testing.T) {
	query := elastic.NewQueryBuilder().
		Size(2).
		Source("message_text", "created_at").
		Sort("created_at", elastic.Desc).
		TermNot("message_text.keyword", "").
		Bool(func(qb *elastic.QueryBuilder) {
			qb.Should(elastic.NewQueryBuilder().Term("type.keyword", "message_send").BuildQuery())
			qb.Should(elastic.NewQueryBuilder().Term("type.keyword", "message_recv").BuildQuery())
			qb.MinimumShouldMatch(1)
		}).
		MatchPhraseOptions("message_text", elastic.MatchPhraseOptions{
			Query: "hello",
			Slop:  util.Pointer(1),
		})
	fmt.Println(query.String())
	res, err := elastic.Document.Search(context.Background(), Index, query.BuildForce())
	assert.NoError(t, err)
	_ = res
	// enc := json.NewEncoder(os.Stdout)
	// enc.SetIndent("", "  ")
	// enc.Encode(res)
}

func TestDocumentQueryBuilderAggs(t *testing.T) {
	now := time.Now()
	query := elastic.NewQueryBuilder().
		Size(0).
		TimeRange("created_at", now.Add(-24*time.Hour*365), now).
		TermNot("message_chat_name.keyword", "").
		AggsTerm("unique_message_chat_name", "message_chat_name.keyword", 1000)
	query2 := elastic.NewQueryBuilder().
		Size(0).
		TimeRange("created_at", now.Add(-24*time.Hour*365), now).
		TermNot("message_chat_name.keyword", "").
		AggsTerm("unique_message_chat_name", "message_chat_name.keyword", 1000, "_count", "asc")
	query3 := elastic.NewQueryBuilder().
		Size(0).
		TimeRange("created_at", now.Add(-24*time.Hour*365), now).
		TermNot("message_user_name.keyword", "").
		AggsTerm("unique_message_user_name", "message_user_name.keyword", 1000, "_count", "asc")
	query4 := elastic.NewQueryBuilder().
		Size(0).
		TimeRange("created_at", now.Add(-24*time.Hour*365), now).
		TermNot("message_peer_user_name.keyword", "").
		AggsTerm("unique_message_peer_user_name", "message_peer_user_name.keyword", 1000, "_count", "asc")

	fmt.Println(query.String())
	fmt.Println(query2.String())
	fmt.Println(query3.String())
	fmt.Println(query4.String())
	res, err := elastic.Document.Search(context.Background(), Index, query4.BuildForce())
	assert.NoError(t, err)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(res)
}
