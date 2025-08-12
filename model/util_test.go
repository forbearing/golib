package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAreTypesEqual(t *testing.T) {
	require.True(t, AreTypesEqual[*User, *User, *User]())
	require.False(t, AreTypesEqual[*User, User, *User]())
	require.False(t, AreTypesEqual[*User, *User, User]())
	require.False(t, AreTypesEqual[*User, User, User]())
	require.False(t, AreTypesEqual[*User, *Menu, *Menu]())
	require.False(t, AreTypesEqual[*User, string, *User]())
	require.False(t, AreTypesEqual[*User, *User, int]())
}

func BenchmarkAreTypesEqual(b *testing.B) {
	b.Run("test1", func(b *testing.B) {
		for b.Loop() {
			AreTypesEqual[*User, *User, *User]()
		}
	})
	b.Run("test2", func(b *testing.B) {
		for b.Loop() {
			AreTypesEqual[*User, User, *User]()
		}
	})
	b.Run("test3", func(b *testing.B) {
		for b.Loop() {
			AreTypesEqual[*User, *User, User]()
		}
	})
	b.Run("test4", func(b *testing.B) {
		for b.Loop() {
			AreTypesEqual[*User, User, User]()
		}
	})
	b.Run("test5", func(b *testing.B) {
		for b.Loop() {
			AreTypesEqual[*User, *Menu, *Menu]()
		}
	})
	b.Run("test6", func(b *testing.B) {
		for b.Loop() {
			AreTypesEqual[*User, string, *User]()
		}
	})
	b.Run("test7", func(b *testing.B) {
		for b.Loop() {
			AreTypesEqual[*User, *User, int]()
		}
	})
}
