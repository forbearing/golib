package types

type ControllerConfig[M Model] struct {
	DB        any // only support *gorm.DB
	TableName string
	ParamName string
}
