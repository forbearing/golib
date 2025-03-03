package statsviz

import (
	"context"
	"fmt"
	"net/http"

	"github.com/arl/statsviz"
	"github.com/forbearing/golib/config"
	"go.uber.org/zap"
)

func Run() error {
	if config.App.EnableStatsviz {
		mux := http.NewServeMux()
		statsviz.Register(mux)
		zap.S().Infow("successfully start statsviz", "listen", config.App.ServerConfig.StatsvizList, "port", config.App.ServerConfig.StatsvizPort)
		return http.ListenAndServe(fmt.Sprintf("%s:%d", config.App.ServerConfig.StatsvizList, config.App.ServerConfig.StatsvizPort), mux)
	}

	<-context.Background().Done()
	return nil
}
