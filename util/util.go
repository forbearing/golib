package util

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
	"unsafe"

	tcping "github.com/cloverstd/tcping/ping"
	"github.com/cockroachdb/errors"
	"github.com/go-ping/ping"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func UUID(prefix ...string) string {
	var id uuid.UUID
	if v7, err := uuid.NewV7(); err == nil {
		id = v7
	} else {
		id = uuid.New()
	}
	if len(prefix) > 0 {
		if len(prefix[0]) > 0 {
			return fmt.Sprintf("%s%s", prefix[0], id.String())
		}
	}
	return id.String()
}

func RequestID() string { return xid.New().String() }
func TraceID() string   { return xid.New().String() }
func SpanID() string    { return xid.New().String() }

func IndexedUUID() string { return ksuid.New().String() }

// Pointer will return a pointer to T with given value.
func Pointer[T comparable](t T) *T {
	if reflect.DeepEqual(t, nil) {
		return new(T)
	}
	return &t
}

// Depointer will return a T with given value.
func Depointer[T comparable](t *T) T {
	if t == nil {
		return *new(T)
	}
	return *t
}

func SafePointer[T any](v T) T {
	if reflect.DeepEqual(v, nil) {
		return *new(T)
	}
	if reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface()) {
		return *new(T)
		// return reflect.Zero(reflect.TypeOf(v)).Interface().(T)
	}
	return v
}

// CharSpliter is the custom split function for bufio.Scanner.
func CharSpliter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}
	if atEOF {
		return len(data), data, nil
	}
	if data[0] == '|' {
		return 1, data[:1], nil
	}
	return 0, nil, nil
}

// SplitByDoublePipe is the custom split function for bufio.Scanner.
func SplitByDoublePipe(data []byte, atEOF bool) (advance int, token []byte, err error) {
	delimiter := []byte("||")

	// Search for the delimiter in the input data
	if i := bytes.Index(data, delimiter); i >= 0 {
		return i + len(delimiter), data[:i], nil
	}

	// If the delimiter is not found, and it's at the end of the input data, return it
	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}

	// If no delimiter is found, return no data and wait for more input
	return 0, nil, nil
}

// RunOrDie will panic when error encountered.
func RunOrDie(fn func() error) {
	if err := fn(); err != nil {
		name := GetFunctionName(fn)
		HandleErr(fmt.Errorf("%s error: %+v", name, err))
	}
}

// HandleErr will call os.Exit() when any error encountered.
func HandleErr(err error, notExit ...bool) {
	var flag bool
	if len(notExit) != 0 {
		flag = notExit[0]
	}
	if err != nil {
		fmt.Println(err)
		if !flag {
			os.Exit(1)
		}
	}
}

// CheckErr just check error and print it.
func CheckErr(err error) {
	HandleErr(err, true)
}

// StringAny format anything to string.
func StringAny(x any) string {
	if x == nil {
		return ""
	}
	if v, ok := x.(fmt.Stringer); ok {
		return v.String()
	}

	switch v := x.(type) {
	case string:
		return v
	case []byte:
		return *(*string)(unsafe.Pointer(&v))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", x)
	default:
		return fmt.Sprintf("%v", x)
	}
}

func GetFunctionName(x any) string {
	switch v := x.(type) {
	case uintptr:
		return runtime.FuncForPC(v).Name()
	default:
		return runtime.FuncForPC(reflect.ValueOf(x).Pointer()).Name()
	}
}

func ParseScheme(req *http.Request) string {
	if scheme := req.Header.Get("x-forwarded-proto"); len(scheme) != 0 {
		return scheme
	}
	if scheme := req.Header.Get("x-forwarded-protocol"); len(scheme) != 0 {
		return scheme
	}
	if ssl := req.Header.Get("x-forwarded-ssl"); ssl == "on" {
		return "https"
	}
	if scheme := req.Header.Get("x-url-scheme"); len(scheme) != 0 {
		return scheme
	}
	if req.TLS != nil {
		return "https"
	}
	return ""
}

// _tcping work like command `tcping`.
func Tcping(host string, port int, timeout time.Duration) bool {
	if timeout < 500*time.Millisecond {
		timeout = 1 * time.Second
	}
	_, _, _, res := _tcping(host, port, 1, 1, timeout)
	return res.SuccessCounter == 1
}

func _tcping(host string, port, count int, interval, timeout time.Duration) (minLatency, maxLatency, avgLatency time.Duration, result *tcping.Result) {
	pinger := tcping.NewTCPing()
	pinger.SetTarget(&tcping.Target{
		Protocol: tcping.TCP,
		Host:     host,
		Port:     port,
		Counter:  count,
		Interval: interval,
		Timeout:  timeout,
	})
	<-pinger.Start()
	if pinger.Result() == nil {
		return
	}
	return pinger.Result().MinDuration, pinger.Result().MaxDuration, pinger.Result().Avg(), pinger.Result()
}

// Ping work like command `ping`.
// If target ip is reachable, return true, nil,
// If target ip is unreachable, return false, nil,
// If error encountered, return false, error.
// More usage see tests in `pkg/util/util_test.go`.
func Ping(ip string, timeout time.Duration) (bool, error) {
	if len(ip) == 0 {
		return false, errors.New("ip is empty")
	}
	if timeout < 500*time.Millisecond {
		timeout = 1 * time.Second
	}
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return false, err
	}
	pinger.Count = 1
	pinger.Timeout = timeout

	err = pinger.Run()
	if err != nil {
		return false, err
	}
	return pinger.Statistics().PacketsSent == pinger.Statistics().PacketsRecv, nil
}

// NoError call fn and always return nil.
func NoError(fn func() error) error {
	if err := fn(); err != nil {
		zap.S().Warn(err)
	}
	return nil
}

// Contains check T in slice.
func Contains[T comparable](slice []T, elem T) bool {
	for i := range slice {
		if slice[i] == elem {
			return true
		}
	}
	return false
}

// CombineError combine error from fns.
func CombineError(fns ...func() error) error {
	errs := make([]error, len(fns))
	for i := range fns {
		if fns[i] == nil {
			continue
		}
		errs[i] = fns[i]()
	}
	return multierr.Combine(errs...)
}

// FileExists check file exists.
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return err == nil
	}
}

// Round returns a rounded version of x with a specified precision.
//
// The precision parameter specifies the number of decimal places to round to.
// Round uses the "round half away from zero" rule to break ties.
//
// Examples:
//
//	Round(3.14159, 2) returns 3.14
//	Round(3.14159, 1) returns 3.1
//	Round(-3.14159, 1) returns -3.1
func Round[T float32 | float64](value T, precision uint) T {
	ratio := math.Pow(10, float64(precision))
	return T(math.Round(float64(value)*ratio) / ratio)
}

// IPv6ToIPv4 converts IPv6 to IPv4 if possible
func IPv6ToIPv4(ipStr string) string {
	// If its ipv4, return.
	if net.ParseIP(ipStr).To4() != nil {
		return ipStr
	}

	// handle IPv6 localhost
	if strings.HasPrefix(ipStr, "::") {
		return "127.0.0.1"
	}

	// handle IPv4-mapped IPv6 addresses
	// eg ::ffff:192.0.2.128 或 ::ffff:c000:280
	if strings.Contains(ipStr, "::ffff:") {
		split := strings.Split(ipStr, "::ffff:")
		if len(split) == 2 {
			if ip := net.ParseIP(split[1]).To4(); ip != nil {
				return ip.String()
			}
		}
	}

	// handle embedded IPv4 address
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ipStr
	}

	ip4 := ip.To4()
	if ip4 != nil {
		return ip4.String()
	}

	return ipStr
}
