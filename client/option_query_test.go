package client_test

import (
	"testing"
	"time"

	"github.com/forbearing/gst/client"
	"github.com/stretchr/testify/assert"
)

var addr = "http://localhost:8080"

func Test_OptionQuery(t *testing.T) {
	t.Run("WithQuery", func(t *testing.T) {
		cli, err := client.New(addr, client.WithQuery("name", "tom", "age", 20, "_sortby", "created_at desc,name asc"))
		assert.NoError(t, err)
		query, err := cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "name=tom&age=20&_sortby=created_at+desc%2Cname+asc", query)

		cli, err = client.New(addr, client.WithQuery("name", "tom", "age", 20, "suname"))
		assert.NoError(t, err)
		query, err = cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "name=tom&age=20", query)

		cli, err = client.New(addr, client.WithQuery("name", "tom", "age", 20, "suname"), client.WithQueryIndex("idx_composite_name_createdat"))
		assert.NoError(t, err)
		query, err = cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "name=tom&age=20&_index=idx_composite_name_createdat", query)
	})

	t.Run("WithQueryPagination", func(t *testing.T) {
		cli, err := client.New(addr, client.WithQueryPagination(1, 10))
		assert.NoError(t, err)
		query, err := cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "page=1&size=10", query)
	})

	t.Run("WithQueryExpand", func(t *testing.T) {
		cli, err := client.New(addr, client.WithQueryExpand("all", 3))
		assert.NoError(t, err)
		query, err := cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "_depth=3&_expand=all", query)

		cli, err = client.New(addr, client.WithQueryExpand("children,parent", 3))
		assert.NoError(t, err)
		query, err = cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "_depth=3&_expand=children%2Cparent", query)
	})

	t.Run("WithQueryFuzzy", func(t *testing.T) {
		cli, err := client.New(addr, client.WithQueryFuzzy(true))
		assert.NoError(t, err)
		query, err := cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "_fuzzy=true", query)
	})

	t.Run("WithQuerySortby", func(t *testing.T) {
		cli, err := client.New(addr, client.WithQuerySortby("created_at desc,id asc"))
		assert.NoError(t, err)
		query, err := cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "_sortby=created_at+desc%2Cid+asc", query)
	})

	t.Run("WithQueryNocache", func(t *testing.T) {
		cli, err := client.New(addr, client.WithQueryNocache(true))
		assert.NoError(t, err)
		query, err := cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "_nocache=true", query)
	})

	t.Run("WithQueryTimeRange", func(t *testing.T) {
		begin := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)
		end := time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local)
		cli, err := client.New(addr, client.WithQueryTimeRange("created_at", begin, end))
		assert.NoError(t, err)
		query, err := cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "_column_name=created_at&_end_time=2022-01-02+00%3A00%3A00&_start_time=2022-01-01+00%3A00%3A00", query)
	})

	t.Run("WithQueryOr", func(t *testing.T) {
		cli, err := client.New(addr, client.WithQueryOr(true))
		assert.NoError(t, err)
		query, err := cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "_or=true", query)
	})

	t.Run("WithQueryIndex", func(t *testing.T) {
		cli, err := client.New(addr, client.WithQueryIndex("idx_composite_name_createdat"))
		assert.NoError(t, err)
		query, err := cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "_index=idx_composite_name_createdat", query)
	})

	t.Run("WithQuerySelect", func(t *testing.T) {
		cli, err := client.New(addr, client.WithQuerySelect("name", "age", ""))
		assert.NoError(t, err)
		query, err := cli.QueryString()
		assert.NoError(t, err)
		assert.Equal(t, "_select=name%2Cage", query)
	})
}
