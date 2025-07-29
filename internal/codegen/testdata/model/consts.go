package model

type Env string

const (
	EnvProd  Env = "prod"
	EnvStg   Env = "stg"
	EnvPre   Env = "pre"
	EnvTest  Env = "test"
	EnvDev   Env = "dev"
	EnvLocal Env = "local"
)
