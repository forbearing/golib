package env

import (
	. "github.com/forbearing/gst/dsl"
	"github.com/forbearing/gst/model"
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
