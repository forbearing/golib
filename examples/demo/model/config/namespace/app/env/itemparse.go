package env

import (
	. "github.com/forbearing/golib/dsl"
	"github.com/forbearing/golib/model"
)

type ItemParse struct {
	model.Base
}

type ItemParseRsp struct {
	Result string `json:"result"`
}

func (ItemParse) Design() {
	Endpoint("items-parse")

	Create(func() {
		Enabled(true)
		Service(true)
		Result[*ItemParseRsp]()
	})
}
