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

	inited bool
)

type task struct {
	name     string
	interval time.Duration
	fn       func() error
	ctx      context.Context
	cancel   context.CancelFunc
}

func Init() error {
	Register(visitor, 60*time.Second, "runtime visitor")

	for _, t := range tasks {
		register(t)
	}

	inited = true
	return nil
}

// Register task can be called at any point before or after Init().
func Register(fn func() error, interval time.Duration, name string) {
	mu.Lock()
	defer mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	if inited {
		register(&task{name: name, fn: fn, interval: interval, ctx: ctx, cancel: cancel})
	} else {
		tasks = append(tasks, &task{name: name, fn: fn, interval: interval, ctx: ctx, cancel: cancel})
	}
}

func register(t *task) {
	if t == nil {
		logger.Task.Warnw("task is nil, skip", "name", t.name, "interval", t.interval.String())
		return
	}
	if t.interval < time.Second {
		logger.Task.Warnw("task interval less than 1 second, skip", "name", t.name, "interval", t.interval.String())
		return
	}
	if t.fn == nil {
		logger.Task.Warnw("task function is nil, skip", "name", t.name, "interval", t.interval.String())
		return
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Task.Errorw(fmt.Sprintf("task panic: %s", err), "name", t.name, "interval", t.interval.String())
			}
		}()
		begin := time.Now()
		logger.Task.Infow("starting task", "name", t.name, "interval", t.interval.String())
		if err := t.fn(); err != nil {
			logger.Task.Errorw(fmt.Sprintf("finished task with error: %s", err), "name", t.name, "interval", t.interval.String(), "cost", time.Since(begin).String())
		} else {
			logger.Task.Infow("finished task", "name", t.name, "interval", t.interval.String(), "cost", time.Since(begin).String())
		}

		ticker := time.NewTicker(t.interval)

		for {
			select {
			case <-t.ctx.Done():
				return
			case <-ticker.C:
				begin = time.Now()
				logger.Task.Infow("starting task", "name", t.name, "interval", t.interval.String())
				if err := t.fn(); err != nil {
					logger.Task.Errorw(fmt.Sprintf("finished task with error: %s", err), "name", t.name, "interval", t.interval.String(), "cost", time.Since(begin).String())
					// return
				} else {
					logger.Task.Infow("finished task", "name", t.name, "interval", t.interval.String(), "cost", time.Since(begin).String())
					// return
				}
			}
		}
	}()
}

func visitor() error {
	logger.Visitor.Info("==================== config ====================")
	logger.Visitor.Infof("%+v", config.App)

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
