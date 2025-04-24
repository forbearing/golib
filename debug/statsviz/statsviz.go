package statsviz

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/arl/statsviz"
	"github.com/forbearing/golib/config"
	"go.uber.org/zap"
)

var server *http.Server

func Run() error {
	if !config.App.StatsvizEnable {
		return nil
	}

	mux := http.NewServeMux()
	statsviz.Register(mux)
	server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.App.StatsvizListen, config.App.StatsvizPort),
		Handler: mux,
	}

	zap.S().Infow("statsviz server started", "listen", config.App.StatsvizListen, "port", config.App.StatsvizPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zap.S().Errorw("failed to start statsviz server", "err", err)
		return err
	}

	return nil
}

func Stop() {
	if server == nil {
		return
	}

	zap.S().Infow("statsviz server shutdown initiated")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		zap.S().Errorw("statsviz server shutdown failed", "err", err)
	} else {
		zap.S().Infow("statsviz server shutdown completed")
	}
	server = nil
}
