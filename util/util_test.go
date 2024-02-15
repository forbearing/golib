package util

import (
	"fmt"
	"testing"
)

func TestUtil(t *testing.T) {
	fmt.Println(UUID())
	fmt.Println(IndexedUUID())
}
