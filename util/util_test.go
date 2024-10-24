package util

import (
	"fmt"
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

func TestTping(t *testing.T) {
	type ip_port struct {
		ip   string
		port int
	}

	onelineIpPorts := []ip_port{
		// {"127.0.0.1", 22},
		{"1.1.1.1", 53},
		{"8.8.8.8", 53},
		{"8.8.4.4", 53},
	}
	offlineIpPorts := []ip_port{
		{"172.16.1.1", 22},
		{"192.168.1.1", 22},
		{"172.16.1.1", 22},
	}
	for _, ipp := range onelineIpPorts {
		isOnline := Tcping(ipp.ip, ipp.port, 1*time.Second)
		fmt.Println(isOnline)
		assert.Equal(t, isOnline, true)
	}

	for _, ipp := range offlineIpPorts {
		isOnline := Tcping(ipp.ip, ipp.port, 1*time.Second)
		fmt.Println(isOnline)
		assert.Equal(t, isOnline, false)
	}
}

func TestPing(t *testing.T) {
	onelineIps := []string{
		"127.0.0.1",
		"1.1.1.1",
		"8.8.8.8",
	}
	offlineIps := []string{
		"172.16.1.1",
		"192.168.1.1",
		"127.0.0.2",
	}
	for _, ip := range onelineIps {
		isOnline, err := Ping(ip, 1*time.Second)
		fmt.Println(isOnline, err)
		assert.NoError(t, err)
		assert.Equal(t, isOnline, true)
	}

	for _, ip := range offlineIps {
		isOnline, err := Ping(ip, 1*time.Second)
		fmt.Println(isOnline, err)
		assert.NoError(t, err)
		assert.Equal(t, isOnline, false)
	}
}

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
