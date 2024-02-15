package controller

import (
	"fmt"

	"github.com/forbearing/golib/database/mysql"
	. "github.com/forbearing/golib/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type column struct{}

var Column = new(column)
var (
	columnAsset = []string{
		"status",
		"owner_entity_id",
		"area_id",
		"category_level1_id",
		"category_level2_id",
		"owner_id",
		"current_user",
		"department_level1_id",
		"department_level2_id",
		"brand",
		"brand_model",
	}

	columnDistribute = []string{
		"status",
		"user_id",
		"department_level1_id",
		"department_level2_id",
	}
	columnBack = []string{
		"status",
	}
	columnChange = []string{
		"status",
	}
	columnCheck = []string{
		"status",
		"progress",
	}
	columnSoftware = []string{
		"hostname",
		"machine_id",
		"mac_addresses",
		"user",
		"ip_addresses",
		"type",
	}
	columnFeishuApproval = []string{
		"status",
		"approval_code",
		"approval_name",
		"user_id",
		"department_id",
	}
	columnFeishuEventApproval = []string{
		"type",
		"status",
		"approval_code",
		"user_id",
	}
)

func (cs *column) Get(c *gin.Context) {
	switch c.Param(PARAM_ID) {
	case "asset":
		cs.Asset(c)
	case "distribute":
		cs.Distribute(c)
	case "back":
		cs.Back(c)
	case "change":
		cs.Change(c)
	case "check":
		cs.Check(c)
	case "software":
		cs.Software(c)
	case "feishu_approval":
		cs.FeishuApproval(c)
	case "feishu_event_approval":
		cs.FeishuEventApproval(c)
	default:
		zap.S().Warn("unknow id: ", c.Param(PARAM_ID))
		ResponseJSON(c, CodeSuccess)
	}
}

func (cs *column) Asset(c *gin.Context) {
	columnRes, err := queryColumns("assets", columnAsset)
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnRes)
}

func (cs *column) Distribute(c *gin.Context) {
	columnRes, err := queryColumns("distributes", columnDistribute)
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnRes)
}
func (cs *column) Back(c *gin.Context) {
	columnRes, err := queryColumns("backs", columnBack)
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnRes)
}
func (cs *column) Change(c *gin.Context) {
	columnRes, err := queryColumns("changes", columnChange)
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnRes)
}
func (cs *column) Check(c *gin.Context) {
	columnsRes, err := queryColumns("checks", columnCheck)
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnsRes)
}
func (cs *column) Software(c *gin.Context) {
	columnsRes, err := queryColumns("softwares", columnSoftware)
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnsRes)
}
func (cs *column) FeishuApproval(c *gin.Context) {
	columnsRes, err := queryColumns("feishu_approvals", columnFeishuApproval)
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnsRes)
}
func (cs *column) FeishuEventApproval(c *gin.Context) {
	columnsRes, err := queryColumns("feishu_event_approvals", columnFeishuEventApproval)
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnsRes)
}

// queryColumns 只查询字段有多少种
//
// select category_level2_id from assets group by category_level2_id;
// +--------------------+
// | category_level2_id |
// +--------------------+
// | BJ                 |
// | NU                 |
// | XS                 |
// | ZJ                 |
// +--------------------+
func queryColumns(table string, columns []string) (map[string][]string, error) {
	cr := make(map[string][]string)
	sql := `SELECT %s FROM %s WHERE deleted_at IS NULL GROUP BY %s`
	for _, column := range columns {
		rows, err := mysql.Default.Raw(fmt.Sprintf(sql, column, table, column)).Rows()
		if err != nil {
			zap.S().Error(err)
			return nil, err
		}
		results := make([]string, 0)
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return nil, err
			}
			results = append(results, name)
		}
		cr[column] = results
	}
	return cr, nil
}

// queryColumns 只查询字段有多少种, 并且计算每种字段值的个数
//
// select category_level2_id, count(*) as category_count from assets group by category_level2_id;
// +--------------------+----------------+
// | category_level2_id | category_count |
// +--------------------+----------------+
// | BJ                 |            110 |
// | NU                 |            800 |
// | XS                 |            328 |
// | ZJ                 |            215 |
// +--------------------+----------------+
//
// select department_level2_id, count(*) as department_count from assets group by department_level2_id;
// +-------------------------------------+------------------+
// | department_level2_id                | department_count |
// +-------------------------------------+------------------+
// |                                     |             1236 |
// | od-ea0ed19af82622a997edf6c2aab262bc |               28 |
// | od-9011520298e3aca4f245e075dd873d02 |               10 |
// | od-3a87018f46f9d37fa811503745fc0b05 |                5 |
// | od-60e10a8929373b1ac0aff828dd5cacf8 |               30 |
// | od-198eb3d20e4783518acee52b1bc48356 |               20 |
// | od-ed452e84ca58c26719ea0ca8b8acecdd |                4 |
// | od-1d7f4ac953b109f2a7e2a2366f5f315e |               72 |
// | od-c6bbbc7f089b356cd45396e3443d1558 |                2 |
// | od-39c14e77f3504a8ca05f3681e9d0470b |                3 |
// | od-095e7e716c0a8262b3dad7888eb4776b |               42 |
// | od-7e8d4fb875bed78400bc5bbca88eed0c |                1 |
// +-------------------------------------+------------------+
func queryColumnsAndCount(table string, columns []string) (columnResult, error) {
	cr := make(map[string][]result)
	sql := `SELECT %s, count(*) as count FROM %s where deleted_at IS NULL GROUP BY %s `
	for _, column := range columns {
		rows, err := mysql.Default.Raw(fmt.Sprintf(sql, column, table, column)).Rows()
		if err != nil {
			zap.S().Error(err)
			return nil, err
		}
		results := make([]result, 0)
		for rows.Next() {
			var name string
			var count uint
			if err := rows.Scan(&name, &count); err != nil {
				return nil, err
			}
			results = append(results, result{name, count})
		}
		cr[column] = results
	}
	return cr, nil
}

type columnResult map[string][]result

type result struct {
	Name  string
	Count uint
}
