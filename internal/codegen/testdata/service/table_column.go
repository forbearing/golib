package service

import (
	"errors"
	"fmt"

	"github.com/forbearing/golib/database"
	model_internal "github.com/forbearing/golib/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	. "github.com/forbearing/golib/util"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

func init() {
	service.Register[*tableColumn]()
}

type tableColumn struct {
	service.Base[*model_internal.TableColumn]
}

func (t *tableColumn) ListBefore(ctx *types.ServiceContext, data *[]*model_internal.TableColumn) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())

	query := new(model_internal.TableColumn)
	schema.NewDecoder().Decode(query, ctx.URL.Query())
	log.Infoz("query", zap.Object("table_column", query))

	if len(query.UserId) == 0 {
		log.Warn("query parameter 'user_id' length is 0")
		query.UserId = "-"
	}
	if len(query.TableName) == 0 {
		err := errors.New("query parameter 'table_name' length is 0")
		log.Error(err)
		return err
	}

	// TODO: 检查是否是有效用户

	switch query.TableName {
	case "machines":
		handleTableColumns(ctx, log, query, DefaultMachineColumns)
	}

	return nil
}

func handleTableColumns(
	ctx *types.ServiceContext,
	log types.Logger,
	query *model_internal.TableColumn,
	defaults []*model_internal.TableColumn,
) error {
	var original []*model_internal.TableColumn
	var columns []*model_internal.TableColumn
	if err := database.Database[*model_internal.TableColumn]().WithQuery(&model_internal.TableColumn{UserId: query.UserId, TableName: query.TableName}).List(&original); err != nil {
		log.Error(err)
		return err
	}
	// 如果在数据库中没有 asset 表的列信息数据, 或者少于指定条列的数据,
	// 比如 assets 有21条列信息, 但是数据库中只有19条,那么也需要再创建下
	if len(original) != len(defaults) {
		columns = make([]*model_internal.TableColumn, len(defaults))
		copy(columns, defaults)
		for i := range columns {
			columns[i].UserId = query.UserId // 别忘了 user_id, model.DefaultAssetColumns 是没有设置 user_id 的
			columns[i].ID = fmt.Sprintf("%s-%s-%s", query.UserId, query.TableName, columns[i].Key)
		}
		if err := database.Database[*model_internal.TableColumn]().Create(columns...); err != nil {
			log.Error(err)
			return err
		}
	}
	// 重置列, 也就是重新创建一下, 但是必须确保只能重置自己的配置.
	if len(ctx.URL.Query().Get("reset")) > 0 {
		if query.UserId != ctx.UserId {
			err := errors.New("forbidden reset other user table column configurations")
			log.Errorw(err.Error(), "source", ctx.UserId, "target", query.UserId)
			return err
		}
		columns = make([]*model_internal.TableColumn, len(defaults))
		copy(columns, defaults)
		for i := range columns {
			columns[i].UserId = query.UserId // 别忘了 user_id, model.DefaultAssetColumns 是没有设置 user_id 的
			columns[i].ID = fmt.Sprintf("%s-%s-%s", query.UserId, query.TableName, columns[i].Key)
		}
		// 创建之前先全部全删除一下
		old := make([]*model_internal.TableColumn, 0)
		if err := database.Database[*model_internal.TableColumn]().WithLimit(-1).WithQuery(&model_internal.TableColumn{UserId: ctx.UserId}).List(&old); err != nil {
			log.Error(err)
			return err
		}
		log.Infoz("successfully find", zap.Int("total", len(old)))
		if err := database.Database[*model_internal.TableColumn]().Delete(old...); err != nil {
			log.Error(err)
			return err
		}
		log.Infoz("successfully delete", zap.Int("total", len(old)))
		if err := database.Database[*model_internal.TableColumn]().Create(columns...); err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

const noneId = "none"

var DefaultMachineColumns = []*model_internal.TableColumn{
	{UserId: noneId, TableName: "machines", Name: "ID", Key: "id", Visiable: ValueOf(true), Sequence: ValueOf[uint](1), Width: ValueOf[uint](100)},
	{UserId: noneId, TableName: "machines", Name: "主机名", Key: "hostname", Visiable: ValueOf(true), Sequence: ValueOf[uint](2), Width: ValueOf[uint](200), Fixed: ValueOf(string(model_internal.FIXED_LEFT))},
	{UserId: noneId, TableName: "machines", Name: "IP地址", Key: "ip_address", Visiable: ValueOf(true), Sequence: ValueOf[uint](3), Width: ValueOf[uint](200)},
	{UserId: noneId, TableName: "machines", Name: "邮箱", Key: "email", Visiable: ValueOf(true), Sequence: ValueOf[uint](4), Width: ValueOf[uint](200)},
	{UserId: noneId, TableName: "machines", Name: "所有者", Key: "owner", Visiable: ValueOf(true), Sequence: ValueOf[uint](5), Width: ValueOf[uint](100)},

	{UserId: noneId, TableName: "machines", Name: "CPU", Key: "cpu", Visiable: ValueOf(true), Sequence: ValueOf[uint](6), Width: ValueOf[uint](100)},
	{UserId: noneId, TableName: "machines", Name: "内存", Key: "memory", Visiable: ValueOf(true), Sequence: ValueOf[uint](7), Width: ValueOf[uint](100)},
	{UserId: noneId, TableName: "machines", Name: "存储", Key: "storage", Visiable: ValueOf(true), Sequence: ValueOf[uint](8), Width: ValueOf[uint](100)},
	{UserId: noneId, TableName: "machines", Name: "网络", Key: "network", Visiable: ValueOf(true), Sequence: ValueOf[uint](9), Width: ValueOf[uint](100)},

	{UserId: noneId, TableName: "machines", Name: "成本", Key: "cost", Visiable: ValueOf(true), Sequence: ValueOf[uint](10), Width: ValueOf[uint](100)},
	{UserId: noneId, TableName: "machines", Name: "部署应用", Key: "deployed_apps", Visiable: ValueOf(true), Sequence: ValueOf[uint](11), Width: ValueOf[uint](100)},
	{UserId: noneId, TableName: "machines", Name: "所属项目", Key: "project_name", Visiable: ValueOf(true), Sequence: ValueOf[uint](12), Width: ValueOf[uint](150)},
	{UserId: noneId, TableName: "machines", Name: "供应商", Key: "vendor_name", Visiable: ValueOf(true), Sequence: ValueOf[uint](13), Width: ValueOf[uint](150)},

	{UserId: noneId, TableName: "machines", Name: "备注", Key: "remark", Visiable: ValueOf(true), Sequence: ValueOf[uint](20), Width: ValueOf[uint](200)},
	{UserId: noneId, TableName: "machines", Name: "创建时间", Key: "created_at", Visiable: ValueOf(true), Sequence: ValueOf[uint](21), Width: ValueOf[uint](160)},
	{UserId: noneId, TableName: "machines", Name: "更新时间", Key: "updated_at", Visiable: ValueOf(true), Sequence: ValueOf[uint](22), Width: ValueOf[uint](160)},
	{UserId: noneId, TableName: "machines", Name: "创建人", Key: "created_by", Visiable: ValueOf(true), Sequence: ValueOf[uint](23), Width: ValueOf[uint](100)},
	{UserId: noneId, TableName: "machines", Name: "更新人", Key: "updated_by", Visiable: ValueOf(true), Sequence: ValueOf[uint](24), Width: ValueOf[uint](100)},
	{UserId: noneId, TableName: "machines", Name: "操作", Key: "operation", Visiable: ValueOf(true), Sequence: ValueOf[uint](25), Width: ValueOf[uint](170), Fixed: ValueOf(string(model_internal.FIXED_RIGHT))},
}
