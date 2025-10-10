package task

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/forbearing/gst/logger"
	"github.com/forbearing/gst/util"
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
	Register(runtimestats, 60*time.Second, "runtime stats")

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
			logger.Task.Errorw(fmt.Sprintf("finished task with error: %s", err), "name", t.name, "interval", t.interval.String(), "cost", util.FormatDurationSmart(time.Since(begin)))
		} else {
			logger.Task.Infow("finished task", "name", t.name, "interval", t.interval.String(), "cost", util.FormatDurationSmart(time.Since(begin)))
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
					logger.Task.Errorw(fmt.Sprintf("finished task with error: %s", err), "name", t.name, "interval", t.interval.String(), "cost", util.FormatDurationSmart(time.Since(begin)))
					// return
				} else {
					logger.Task.Infow("finished task", "name", t.name, "interval", t.interval.String(), "cost", util.FormatDurationSmart(time.Since(begin)))
					// return
				}
			}
		}
	}()
}

func runtimestats() error {
	rtm := new(runtime.MemStats)
	runtime.ReadMemStats(rtm)

	// 基本运行时信息
	logger.Runtime.Infow("Basic Runtime Info",
		"GoVersion", runtime.Version(),
		"GOMAXPROCS", runtime.GOMAXPROCS(0),
		"NumCPU", runtime.NumCPU(),
		"Goroutines", runtime.NumGoroutine(),
		"CGOCalls", runtime.NumCgoCall(),
		"OS", runtime.GOOS,
		"Arch", runtime.GOARCH,
		"Compiler", runtime.Compiler,
	)

	// 内存分配和GC统计
	logger.Runtime.Infow("Memory Allocation Stats",
		"Mallocs", rtm.Mallocs,
		"Frees", rtm.Frees,
		"LiveObjects", rtm.Mallocs-rtm.Frees,
		"TotalAlloc", rtm.TotalAlloc,
		"Sys", rtm.Sys,
		"Lookups", rtm.Lookups,
		"Alloc", rtm.Alloc, // 当前分配的内存
	)

	// 堆内存统计
	logger.Runtime.Infow("Heap Memory Stats",
		"HeapAlloc", rtm.HeapAlloc,
		"HeapSys", rtm.HeapSys,
		"HeapIdle", rtm.HeapIdle,
		"HeapInuse", rtm.HeapInuse,
		"HeapReleased", rtm.HeapReleased,
		"HeapObjects", rtm.HeapObjects,
	)

	// 栈内存统计
	logger.Runtime.Infow("Stack Memory Stats",
		"StackInuse", rtm.StackInuse,
		"StackSys", rtm.StackSys,
		"MSpanInuse", rtm.MSpanInuse,
		"MSpanSys", rtm.MSpanSys,
		"MCacheInuse", rtm.MCacheInuse,
		"MCacheSys", rtm.MCacheSys,
	)

	// GC统计
	logger.Runtime.Infow("GC Stats",
		"NumGC", rtm.NumGC,
		"LastGC", time.UnixMilli(int64(rtm.LastGC/1_000_000)), //nolint:gosec
		"PauseTotalNs", rtm.PauseTotalNs,
		"PauseNs", rtm.PauseNs[(rtm.NumGC%256)], // 最近一次GC暂停时间
		"PauseEnd", rtm.PauseEnd[(rtm.NumGC%256)], // 最近一次GC暂停结束时间
		"GCCPUFraction", rtm.GCCPUFraction, // GC占用CPU时间的比例
		"EnableGC", rtm.EnableGC,
		"DebugGC", rtm.DebugGC,
		"NumForcedGC", rtm.NumForcedGC, // 强制GC的次数
	)

	// GC暂停历史记录（最近几次）
	gcHistory := make(map[string]any)
	for i := 0; i < int(rtm.NumGC) && i < 5; i++ {
		idx := int(rtm.NumGC-uint32(i)) % 256 //nolint:gosec
		gcHistory[fmt.Sprintf("GC-%d-PauseNs", i+1)] = rtm.PauseNs[idx]
		gcHistory[fmt.Sprintf("GC-%d-End", i+1)] = time.UnixMilli(int64(rtm.PauseEnd[idx] / 1_000_000)) //nolint:gosec
	}
	logger.Runtime.Infow("Recent GC History", "gcHistory", gcHistory)

	// 其他内存统计
	logger.Runtime.Infow("Other Memory Stats",
		"BuckHashSys", rtm.BuckHashSys,
		"GCSys", rtm.GCSys,
		"OtherSys", rtm.OtherSys,
		"NextGC", rtm.NextGC,
	)

	// 进程信息 (cross-platform)
	getProcessStats()

	// 互斥锁争用情况（仅在设置了GODEBUG=mutexprofile=1时可用）
	logger.Runtime.Info("Note: For mutex contention profiling, run with GODEBUG=mutexprofile=1")

	return nil
}
