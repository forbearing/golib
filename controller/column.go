package controller

import (
	"fmt"
	"regexp"
	"strings"

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
		"owner_status",
		"label_status",
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
		"type",
		"os",
		"hostname",
		"user",
		"ip_addresses",
		"mac_addresses",
		"machine_id",
		// "name",
		// "version",
		// "maintainer",
		// "description",
		// "home_page",
	}
	columnSoftwrarePurchased = []string{
		"name",
		"version",
		"platform",
		"purchase_user",
		"purchase_method",
		"purchase_date",
		"expire_date",
		"license_type",
		"license_quantity",
		"vendor",
		"cost",
	}
	columnSoftwrarePurchasedAssignment = []string{
		"type",
		"machine_id",
		"hostname",
		"user",
		"ip_addresses",
		"mac_addresses",
		"name",
		"version",
	}
	columnSoftwareReport = []string{
		"name",
		"version",
		"status",
		"type",
		"os",
		"hostname",
		"user",
		"ip_addresses",
		"mac_addresses",
		"machine_id",
		"maintainer",
		"description",
		"home_page",
	}
	columnSoftwareCatalog = []string{
		"label",
		"os",
		"name",
		"version",
		"maintainer",
		// "description",
		// "home_page",
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
		// cs.Asset(c)
		cs.GetColumns(c, "assets", columnAsset)
	case "distribute":
		// cs.Distribute(c)
		cs.GetColumns(c, "distributes", columnDistribute)
	case "back":
		// cs.Back(c)
		cs.GetColumns(c, "backs", columnBack)
	case "change":
		// cs.Change(c)
		cs.GetColumns(c, "changes", columnChange)
	case "check":
		// cs.Check(c)
		cs.GetColumns(c, "checks", columnCheck)
	case "software":
		// cs.Software(c)
		cs.GetColumns(c, "softwares", columnSoftware)
	case "software_purchased":
		// cs.SoftwarePurchased(c)
		cs.GetColumns(c, "software_purchaseds", columnSoftwrarePurchased)
	case "software_purchased_assignment":
		// cs.SoftwarePurchasedAssignment(c)
		cs.GetColumns(c, "software_purchased_assignments", columnSoftwrarePurchasedAssignment)
	case "software_report":
		// cs.SoftwareReport(c)
		cs.GetColumns(c, "software_reports", columnSoftwareReport)
	case "software_catalog":
		// cs.SoftwareCatalog(c)
		cs.GetColumns(c, "software_catalogs", columnSoftwareCatalog)
	case "feishu_approval":
		// cs.FeishuApproval(c)
		cs.GetColumns(c, "feishu_approvals", columnFeishuApproval)
	case "feishu_event_approval":
		// cs.FeishuEventApproval(c)
		cs.GetColumns(c, "feishu_event_approvals", columnFeishuEventApproval)
	default:
		zap.S().Warn("unknow id: ", c.Param(PARAM_ID))
		ResponseJSON(c, CodeSuccess)
	}
}

func (cs *column) GetColumns(c *gin.Context, tableName string, columns []string) {
	columnRes, err := queryColumnsWithQuery(tableName, columns, c.Request.URL.Query())
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnRes)
}

func (cs *column) Asset(c *gin.Context) {
	columnRes, err := queryColumnsWithQuery("assets", columnAsset, c.Request.URL.Query())
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnRes)
}

func (cs *column) Distribute(c *gin.Context) {
	columnRes, err := queryColumnsWithQuery("distributes", columnDistribute, c.Request.URL.Query())
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnRes)
}

func (cs *column) Back(c *gin.Context) {
	columnRes, err := queryColumnsWithQuery("backs", columnBack, c.Request.URL.Query())
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnRes)
}

func (cs *column) Change(c *gin.Context) {
	columnRes, err := queryColumnsWithQuery("changes", columnChange, c.Request.URL.Query())
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnRes)
}

func (cs *column) Check(c *gin.Context) {
	columnsRes, err := queryColumnsWithQuery("checks", columnCheck, c.Request.URL.Query())
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnsRes)
}

func (cs *column) Software(c *gin.Context) {
	columnsRes, err := queryColumnsWithQuery("softwares", columnSoftware, c.Request.URL.Query())
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnsRes)
}

func (cs *column) SoftwarePurchased(c *gin.Context) {
	columnsRes, err := queryColumnsWithQuery("software_purchaseds", columnSoftwrarePurchased, c.Request.URL.Query())
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnsRes)
}

func (cs *column) SoftwarePurchasedAssignment(c *gin.Context) {
	columnsRes, err := queryColumnsWithQuery("software_purchased_assignments", columnSoftwrarePurchasedAssignment, c.Request.URL.Query())
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnsRes)
}

func (cs *column) SoftwareReport(c *gin.Context) {
	columnsRes, err := queryColumnsWithQuery("software_reports", columnSoftwareReport, c.Request.URL.Query())
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnsRes)
}

func (cs *column) SoftwareCatalog(c *gin.Context) {
	columnsRes, err := queryColumnsWithQuery("software_catalogs", columnSoftwareCatalog, c.Request.URL.Query())
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnsRes)
}

func (cs *column) FeishuApproval(c *gin.Context) {
	columnsRes, err := queryColumnsWithQuery("feishu_approvals", columnFeishuApproval, c.Request.URL.Query())
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, columnsRes)
}

func (cs *column) FeishuEventApproval(c *gin.Context) {
	columnsRes, err := queryColumnsWithQuery("feishu_event_approvals", columnFeishuEventApproval, c.Request.URL.Query())
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
			// 前端过滤出空值并且 _fuzzy=true 时,没有任何过滤作用
			// 前端过滤出空值并且 _fuzzy=false 时,查询不到任何结果
			if len(name) == 0 {
				continue
			}
			results = append(results, name)
		}
		cr[column] = results
	}
	return cr, nil
}

func queryColumnsWithQuery(table string, columns []string, query map[string][]string) (map[string][]string, error) {
	cr := make(map[string][]string)
	sql := `SELECT %s FROM %s WHERE deleted_at IS NULL %s GROUP BY %s`

	var queryBuilder strings.Builder
	for k, v := range query { // v eg: [process,package,]
		if len(k) > 0 && len(strings.Join(v, "")) > 0 {
			items := make([]string, 0)
			for _, item := range v {
				if len(item) > 0 && strings.TrimSpace(item) != "," {
					for _, _item := range strings.Split(item, ",") {
						if len(strings.TrimSpace(_item)) > 0 {
							items = append(items, strings.TrimSpace(_item))
						}
					}
				}
			}

			var out strings.Builder
			for i, item := range items {
				switch i {
				case 0:
					if len(items) == 1 {
						out.WriteString(fmt.Sprintf(`('%s')`, regexp.QuoteMeta(strings.TrimSpace(item))))
					} else {
						out.WriteString(fmt.Sprintf(`('%s'`, regexp.QuoteMeta(strings.TrimSpace(item))))
					}
				case len(items) - 1:
					out.WriteString(fmt.Sprintf(`,'%s')`, regexp.QuoteMeta(strings.TrimSpace(item))))
				default:
					out.WriteString(fmt.Sprintf(`,'%s'`, regexp.QuoteMeta(strings.TrimSpace(item))))
				}
			}
			if len(strings.TrimSpace(out.String())) > 0 {
				queryBuilder.WriteString(fmt.Sprintf(" AND `%s` IN %s", k, strings.TrimSpace(out.String())))
			}
		}
	}

	for _, column := range columns {
		statement := fmt.Sprintf(sql, column, table, queryBuilder.String(), column)
		// fmt.Println("---------------------", statement)
		rows, err := mysql.Default.Raw(statement).Rows()
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
			// 前端过滤出空值并且 _fuzzy=true 时,没有任何过滤作用
			// 前端过滤出空值并且 _fuzzy=false 时,查询不到任何结果
			if len(name) == 0 {
				continue
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
