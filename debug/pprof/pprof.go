package pprof

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/forbearing/golib/config"
	"go.uber.org/zap"
)

func Run() error {
	if !config.App.EnablePprof {
		return nil
	}

	log := zap.S()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.App.PprofListen, config.App.PprofPort),
		Handler: http.DefaultServeMux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		log.Infow("pprof server started", "listen", config.App.PprofListen, "port", config.App.PprofPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorw("failed to start pprof", "err", err)
		}
	}()

	sig := <-quit
	log.Infow("pprof shutdown initiated", "signal", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorw("pprof shutdown failed", "err", err, "signal", sig)
		return err
	}
	log.Infow("pprof shutdown completed", "signal", sig)
	return nil
}
