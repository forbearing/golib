package task

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/logger"
)

var (
	tasks []*task
	mu    sync.Mutex

	taskFlag uint64
)

type task struct {
	name     string
	interval time.Duration
	fn       func() error
	ctx      context.Context
	cancel   context.CancelFunc
}

func Init() error {
	Register("runtime visitor", 60*time.Second, visitor)
	// Register("fetch software", 1*time.Hour, fetchPackagesFromSoc)
	// Register("fetch software", 1*time.Hour, fetchProcessesFromSoc)
	// Register("sync ubuntu software to remote", 24*time.Hour, syncUbuntuSoftwareToRemote)

	for i := range tasks {
		t := tasks[i]
		if t == nil {
			logger.Visitor.Warnw("task is nil, skip", "name", t.name, "interval", t.interval.String())
			continue
		}
		if t.interval < time.Second {
			logger.Visitor.Warnw("task interval less than 1 second, skip", "name", t.name, "interval", t.interval.String())
			continue
		}
		if t.fn == nil {
			logger.Visitor.Warnw("task function is nil, skip", "name", t.name, "interval", t.interval.String())
			continue
		}
		go func() {
			logger.Visitor.Infow("starting task", "name", t.name, "interval", t.interval.String())
			if err := t.fn(); err != nil {
				logger.Visitor.Errorw(fmt.Sprintf("finished task with error: %s", err), "name", t.name, "interval", t.interval.String())
				return
			} else {
				logger.Visitor.Infow("finished task", "name", t.name, "interval", t.interval.String())
			}

			ticker := time.NewTicker(t.interval)

			for {
				select {
				case <-t.ctx.Done():
					return
				case <-ticker.C:
					logger.Visitor.Infow("starting task", "name", t.name, "interval", t.interval.String())
					if err := t.fn(); err != nil {
						logger.Visitor.Errorw(fmt.Sprintf("finished task with error: %s", err), "name", t.name, "interval", t.interval.String())
						// return
					} else {
						logger.Visitor.Infow("finished task", "name", t.name, "interval", t.interval.String())
						// return
					}
				}
			}
		}()

	}

	return nil
}

func Register(name string, interval time.Duration, fn func() error) {
	mu.Lock()
	defer mu.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	tasks = append(tasks, &task{name: name, fn: fn, interval: interval, ctx: ctx, cancel: cancel})
}

func visitor() error {
	logger.Visitor.Info("==================== config ====================")
	logger.Visitor.Info(config.App)

	logger.Visitor.Info("==================== runtime ====================")
	rtm := new(runtime.MemStats)
	runtime.ReadMemStats(rtm)
	logger.Visitor.Infow("",
		"Goroutines", runtime.NumGoroutine(),
		"Mallocs", rtm.Mallocs, "Frees", rtm.Frees,
		"LiveObjects", rtm.Mallocs-rtm.Frees, "PauseTotalNs", rtm.PauseTotalNs,
		"NumGC", rtm.NumGC, "LastGC", time.UnixMilli(int64(rtm.LastGC/1_000_000)),
		"HeapObjects", rtm.HeapObjects, "HeapAlloc", rtm.HeapAlloc,
	)
	return nil
}
