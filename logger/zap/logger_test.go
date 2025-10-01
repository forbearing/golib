package zap_test

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/forbearing/gst/config"
	"github.com/forbearing/gst/logger/zap"
	"github.com/forbearing/gst/types"
)

var (
	msg10    = "0000000000"
	msg100   = strings.Repeat(msg10, 10)
	msg1000  = strings.Repeat(msg10, 100)
	msg10000 = strings.Repeat(msg10, 1000)

	keyValues10  = []string{}
	keyValues100 = []string{}
)

func init() {
	// init keyValues10
	for i := 0; i < 10; i++ {
		keyValues10 = append(keyValues10, "key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}
	// init keyValues100
	for i := 0; i < 100; i++ {
		keyValues100 = append(keyValues100, "key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}
}

func createLogger(b *testing.B, filename string) types.Logger {
	os.Setenv(config.LOGGER_FILE, filename)
	os.Setenv(config.LOGGER_DIR, "/tmp/gst")
	if err := config.Init(); err != nil {
		b.Fatal(err)
	}
	l := zap.New()
	return l
}

func TestLogger(b *testing.T) {
	os.Setenv(config.LOGGER_FILE, "")
	os.Setenv(config.LOGGER_DIR, "/tmp/gst")
	if err := config.Init(); err != nil {
		b.Fatal(err)
	}
	l := zap.New()
	l.With("key1", "value1", "key2", "value2").Info("hello world")
}

func BenchmarkLogger_File10(b *testing.B) {
	l := createLogger(b, "test.log")
	b.ResetTimer()
	for range b.N {
		l.Infoz(msg10)
	}
}

func BenchmarkLogger_File100(b *testing.B) {
	l := createLogger(b, "test.log")
	b.ResetTimer()
	for range b.N {
		l.Infoz(msg100)
	}
}

func BenchmarkLogger_File1000(b *testing.B) {
	l := createLogger(b, "test.log")
	b.ResetTimer()
	for range b.N {
		l.Infoz(msg1000)
	}
}

func BenchmarkLogger_File10000(b *testing.B) {
	l := createLogger(b, "test.log")
	b.ResetTimer()
	for range b.N {
		l.Infoz(msg10000)
	}
}

func BenchmarkLogger_Discard10(b *testing.B) {
	l := createLogger(b, "/dev/null")
	b.ResetTimer()
	for range b.N {
		l.Infoz(msg10)
	}
}

func BenchmarkLogger_Discard100(b *testing.B) {
	l := createLogger(b, "/dev/null")
	b.ResetTimer()
	for range b.N {
		l.Infoz(msg100)
	}
}

func BenchmarkLogger_Discard1000(b *testing.B) {
	l := createLogger(b, "/dev/null")
	for range b.N {
		l.Infoz(msg1000)
	}
}

func BenchmarkLogger_Discard10000(b *testing.B) {
	l := createLogger(b, "/dev/null")
	b.ResetTimer()
	for range b.N {
		l.Infoz(msg10000)
	}
}

func BenchmarkLogger_With10(b *testing.B) {
	l := createLogger(b, "test.log")
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		l.With(keyValues10...).Info(msg10)
	}
}

func BenchmarkLogger_With100(b *testing.B) {
	l := createLogger(b, "test.log")
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		l.With(keyValues100...).Info(msg10)
	}
}
