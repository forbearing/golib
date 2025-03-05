package statsviz

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arl/statsviz"
	"github.com/forbearing/golib/config"
	"go.uber.org/zap"
)

func Run() error {
	if !config.App.EnableStatsviz {
		return nil
	}

	log := zap.S()
	mux := http.NewServeMux()
	statsviz.Register(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.App.StatsvizListen, config.App.StatsvizPort),
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		log.Infow("statsviz server started", "listen", config.App.StatsvizListen, "port", config.App.StatsvizPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorw("failed to start statsviz", "err", err)
		}
	}()

	sig := <-quit
	log.Infow("statsviz shutdown initiated", "signal", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorw("statsviz shutdown failed", "err", err, "signal", sig)
		return err
	}
	log.Infow("statsviz shutdown completed", "signal", sig)
	return nil
}
