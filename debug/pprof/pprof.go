package pprof

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"github.com/forbearing/golib/config"
	"go.uber.org/zap"
)

var server *http.Server

func Run() error {
	if !config.App.PprofEnable {
		return nil
	}
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.App.PprofListen, config.App.PprofPort),
		Handler: http.DefaultServeMux,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zap.S().Errorw("failed to start pprof server", "err", err)
		return err
	}
	zap.S().Infow("pprof server started", "listen", config.App.PprofListen, "port", config.App.PprofPort)

	return nil
}

func Stop() {
	if server == nil {
		return
	}

	zap.S().Infow("pprof server shutdown initiated")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		zap.S().Errorw("pprof server shutdown failed", "err", err)
	} else {
		zap.S().Infow("pprof server shutdown completed")
	}
	server = nil
}
