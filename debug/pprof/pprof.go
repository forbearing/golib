package pprof

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/forbearing/golib/config"
	"go.uber.org/zap"
)

func Run() error {
	if config.App.EnablePprof {
		server := &http.Server{
			Addr:    fmt.Sprintf("%s:%d", config.App.PprofListen, config.App.PprofPort),
			Handler: http.DefaultServeMux,
		}
		zap.S().Infow("successfully start pprof", "listen", config.App.PprofListen, "port", config.App.PprofPort)
		return server.ListenAndServe()
	}
	<-context.Background().Done()
	return nil
}
