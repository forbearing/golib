package gops

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/forbearing/golib/config"
	"github.com/google/gops/agent"
	"go.uber.org/zap"
)

func Run() error {
	if !config.App.EnableGops {
		return nil
	}
	log := zap.S()
	err := agent.Listen(agent.Options{
		Addr:            fmt.Sprintf("%s:%d", config.App.GopsListen, config.App.GopsPort),
		ShutdownCleanup: true,
		ConfigDir:       "/tmp/gops",
	})
	if err != nil {
		log.Errorw("gops agent startup failed", "err", err)
		return err
	}

	log.Infow("gops agent started", "listen", config.App.GopsListen, "port", config.App.GopsPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-quit
	log.Infow("gops agent shutdown initiated", "signal", sig)
	agent.Close()
	log.Infow("gops agent shutdown completed", "signal", sig)

	return nil
}
