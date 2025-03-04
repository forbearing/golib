package gops

import (
	"context"
	"fmt"

	"github.com/forbearing/golib/config"
	"github.com/google/gops/agent"
	"go.uber.org/zap"
)

func Run() error {
	if config.App.EnableGops {
		if err := agent.Listen(agent.Options{
			Addr:            fmt.Sprintf("%s:%d", config.App.GopsListen, config.App.GopsPort),
			ShutdownCleanup: true,
			ConfigDir:       "/tmp/gops",
		}); err != nil {
			return err
		}
	}
	zap.S().Infow("successfully start gops", "listen", config.App.GopsListen, "port", config.App.GopsPort)
	<-context.Background().Done()
	return nil
}
