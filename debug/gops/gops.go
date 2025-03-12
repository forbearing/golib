package gops

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/forbearing/golib/config"
	"github.com/google/gops/agent"
	"go.uber.org/zap"
)

var (
	tempDir string
	running bool
	mu      sync.Mutex
)

func Run() error {
	mu.Lock()
	defer mu.Unlock()
	if !config.App.GopsEnable {
		return nil
	}

	tempDir = config.GetTempdir()
	if len(tempDir) == 0 {
		tempDir = "/tmp/gops"
	} else {
		tempDir = filepath.Join(tempDir, "gops")
	}

	if err := agent.Listen(agent.Options{
		Addr: fmt.Sprintf("%s:%d", config.App.GopsListen, config.App.GopsPort),
		// 千万千万别将 ShutdownCleanup 设置为 true, 如果设置为 true, gops 捕捉信号并 os.Exit(1) 调试了我一下午.
		ShutdownCleanup: false,
		ConfigDir:       tempDir,
	}); err != nil {
		zap.S().Errorw("gops agent startup failed", "err", err)
		return err
	}

	running = true
	zap.S().Infow("gops agent started", "listen", config.App.GopsListen, "port", config.App.GopsPort)
	return nil
}

func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if !running {
		return
	}

	zap.S().Infow("gops agent shutdown initiated")
	agent.Close()
	running = false
	zap.S().Infow("gops agent shutdown completed")
}
