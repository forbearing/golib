package util

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUtil(t *testing.T) {
	fmt.Println(UUID())
	fmt.Println(RequestID())
	fmt.Println(IndexedUUID())
}

func BenchmarkUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUID()
	}
}

func BenchmarkIndexedUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IndexedUUID()
	}
}

func BenchmarkLightUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RequestID()
	}
}

// func TestTping(t *testing.T) {
// 	type ip_port struct {
// 		ip   string
// 		port int
// 	}
//
// 	onelineIpPorts := []ip_port{
// 		// {"127.0.0.1", 22},
// 		{"1.1.1.1", 53},
// 		{"8.8.8.8", 53},
// 		{"8.8.4.4", 53},
// 	}
// 	offlineIpPorts := []ip_port{
// 		{"172.16.1.1", 22},
// 		{"192.168.1.1", 22},
// 		{"172.16.1.1", 22},
// 	}
// 	for _, ipp := range onelineIpPorts {
// 		isOnline := Tcping(ipp.ip, ipp.port, 1*time.Second)
// 		fmt.Println(isOnline)
// 		assert.Equal(t, isOnline, true)
// 	}
//
// 	for _, ipp := range offlineIpPorts {
// 		isOnline := Tcping(ipp.ip, ipp.port, 1*time.Second)
// 		fmt.Println(isOnline)
// 		assert.Equal(t, isOnline, false)
// 	}
// }

// func TestPing(t *testing.T) {
// 	onelineIps := []string{
// 		"127.0.0.1",
// 		"1.1.1.1",
// 		"8.8.8.8",
// 	}
// 	offlineIps := []string{
// 		"172.16.1.1",
// 		"192.168.1.1",
// 		"127.0.0.2",
// 	}
// 	for _, ip := range onelineIps {
// 		isOnline, err := Ping(ip, 1*time.Second)
// 		fmt.Println(isOnline, err)
// 		assert.NoError(t, err)
// 		assert.Equal(t, isOnline, true)
// 	}
//
// 	for _, ip := range offlineIps {
// 		isOnline, err := Ping(ip, 1*time.Second)
// 		fmt.Println(isOnline, err)
// 		assert.NoError(t, err)
// 		assert.Equal(t, isOnline, false)
// 	}
// }

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}
	assert.True(t, Contains(slice, "a"))
	assert.True(t, Contains(slice, "b"))
	assert.True(t, Contains(slice, "c"))
	assert.False(t, Contains(slice, "d"))

	slice2 := []int{1, 2, 3}
	assert.True(t, Contains(slice2, 1))
	assert.True(t, Contains(slice2, 2))
	assert.True(t, Contains(slice2, 3))
	assert.False(t, Contains(slice2, 4))
}

func TestFileExists(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test")
	defer os.RemoveAll(tmpFile.Name())
	defer tmpFile.Close()

	assert.NoError(t, err)

	assert.Equal(t, true, FileExists("/tmp"))
	assert.Equal(t, true, FileExists("/tmp/"))
	assert.Equal(t, true, FileExists(tmpFile.Name()))
	assert.Equal(t, false, FileExists(tmpFile.Name()+"---"))
}

func TestRound(t *testing.T) {
	tests := []struct {
		name      string
		value     any // float32 or float64
		precision uint
		want      any // float32 or float64
	}{
		{
			name:      "float64 positive round down",
			value:     3.14159,
			precision: 3,
			want:      3.142,
		},
		{
			name:      "float64 positive round up",
			value:     3.14859,
			precision: 3,
			want:      3.149,
		},
		{
			name:      "float32 positive round down",
			value:     float32(3.14159),
			precision: 3,
			want:      float32(3.142),
		},
		{
			name:      "float64 negative round down",
			value:     -3.14159,
			precision: 3,
			want:      -3.142,
		},
		{
			name:      "float32 negative round down",
			value:     float32(-3.14159),
			precision: 3,
			want:      float32(-3.142),
		},
		{
			name:      "float64 zero precision",
			value:     3.14159,
			precision: 0,
			want:      3.0,
		},
		{
			name:      "float64 large number",
			value:     123456.789,
			precision: 2,
			want:      123456.79,
		},
		{
			name:      "float32 small number",
			value:     float32(0.0000123),
			precision: 7,
			want:      float32(0.0000123),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got any
			switch v := tt.value.(type) {
			case float64:
				got = Round(v, tt.precision)
			case float32:
				got = Round(v, tt.precision)
			}

			// compare result
			switch want := tt.want.(type) {
			case float64:
				got64 := got.(float64)
				if math.Abs(got64-want) > 1e-10 {
					t.Errorf("Round() = %v, want %v", got64, want)
				}
			case float32:
				got32 := got.(float32)
				if math.Abs(float64(got32-want)) > 1e-6 {
					t.Errorf("Round() = %v, want %v", got32, want)
				}
			}
		})
	}
}

func TestHashID(t *testing.T) {
	hashID := HashID("user", "email", "address")
	fmt.Println(hashID)
}

func TestFormatDurationMilliseconds(t *testing.T) {
	tests := []struct {
		duration  time.Duration
		precision int
		expected  string
	}{
		{1234567 * time.Nanosecond, 2, "1.23ms"},
		{1234567 * time.Nanosecond, 3, "1.235ms"},
		{1500 * time.Millisecond, 2, "1500.00ms"},
		{1 * time.Second, 0, "1000ms"},
		{1 * time.Second, -1, "1000.00ms"}, // negative precision, default 2
		{0, 2, "0.00ms"},
		{2500 * time.Microsecond, 4, "2.5000ms"},
	}

	for _, tt := range tests {
		got := FormatDurationMilliseconds(tt.duration, tt.precision)
		if got != tt.expected {
			t.Errorf("FormatDurationMilliseconds(%v, %d) = %s; want %s", tt.duration, tt.precision, got, tt.expected)
		}
	}
}

func BenchmarkRound(b *testing.B) {
	b.Run("float64", func(b *testing.B) {
		value := 3.14159
		for i := 0; i < b.N; i++ {
			Round(value, 3)
		}
	})

	b.Run("float32", func(b *testing.B) {
		value := float32(3.14159)
		for i := 0; i < b.N; i++ {
			Round(value, 3)
		}
	})
}

func ExampleRound() {
	fmt.Printf("%.3f\n", Round(3.14159, 3))
	fmt.Printf("%.2f\n", Round(2.71828, 2))
	fmt.Printf("%.1f\n", Round(-3.14159, 1))
	// Output:
	// 3.142
	// 2.72
	// -3.1
}
