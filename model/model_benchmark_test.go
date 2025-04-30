package model

import "testing"

func BenchmarkHasRequestResponse(b *testing.B) {
	b.Run("HasRequest", func(b *testing.B) {
		for b.Loop() {
			HasRequest[*User]()
		}
	})
	b.Run("HasResponse", func(b *testing.B) {
		for b.Loop() {
			HasResponse[*User]()
		}
	})
	b.Run("NewRequest", func(b *testing.B) {
		for b.Loop() {
			_ = NewRequest[*User]()
		}
	})
	b.Run("NewResponse", func(b *testing.B) {
		for b.Loop() {
			_ = NewResponse[*User]()
		}
	})
}
