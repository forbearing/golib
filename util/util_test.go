package util

import (
	"fmt"
	"testing"
)

func TestUtil(t *testing.T) {
	fmt.Println(UUID())
	fmt.Println(LightUUID())
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
		LightUUID()
	}
}
