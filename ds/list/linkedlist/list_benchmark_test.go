package linkedlist_test

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	golist "container/list"

	"github.com/forbearing/golib/ds/list/linkedlist"
	"github.com/forbearing/golib/ds/util"
	"github.com/stretchr/testify/assert"
)

func BenchmarkList_PushBack(b *testing.B) {
	b.Run("custom unsafe", func(b *testing.B) {
		l, err := linkedlist.New[int]()
		assert.NoError(b, err)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
	})

	b.Run("custom safe", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
	})

	b.Run("custom safe conc", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				l.PushBack(0)
			}
		})
	})

	b.Run("std unsafe", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
	})

	var mu sync.Mutex
	b.Run("std safe", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mu.Lock()
			l.PushBack(i)
			mu.Unlock()
		}
	})

	b.Run("std safe conc", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				mu.Lock()
				l.PushBack(0)
				mu.Unlock()
			}
		})
	})
}

func BenchmarkList_PushFront(b *testing.B) {
	b.Run("custom unsafe", func(b *testing.B) {
		l, err := linkedlist.New[int]()
		assert.NoError(b, err)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.PushFront(i)
		}
	})

	b.Run("custom safe", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.PushFront(i)
		}
	})

	b.Run("custom safe conc", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				l.PushFront(0)
			}
		})
	})

	b.Run("std unsafe", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.PushFront(i)
		}
	})

	var mu sync.Mutex
	b.Run("std safe", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mu.Lock()
			l.PushFront(i)
			mu.Unlock()
		}
	})

	b.Run("std safe conc", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				mu.Lock()
				l.PushFront(0)
				mu.Unlock()
			}
		})
	})
}

func BenchmarkList_InsertAfter(b *testing.B) {
	b.Run("custom unsafe", func(b *testing.B) {
		l, err := linkedlist.New[int]()
		assert.NoError(b, err)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		n := l.Head.Next
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.InsertAfter(n, i)
		}
	})

	b.Run("custom safe", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		n := l.Head.Next
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.InsertAfter(n, i)
		}
	})

	b.Run("custom safe conc", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		n := l.Head.Next
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				l.InsertAfter(n, 0)
			}
		})
	})

	b.Run("std unsafe", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		e := l.Front().Next()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.InsertAfter(e, &golist.Element{Value: i})
		}
	})

	var mu sync.Mutex
	b.Run("std safe", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		e := l.Front().Next()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mu.Lock()
			l.InsertAfter(e, &golist.Element{Value: i})
			mu.Unlock()
		}
	})

	b.Run("std safe conc", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		e := l.Front().Next()
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				mu.Lock()
				l.InsertAfter(e, &golist.Element{Value: 0})
				mu.Unlock()
			}
		})
	})
}

func BenchmarkList_InsertBefore(b *testing.B) {
	b.Run("custom unsafe", func(b *testing.B) {
		l, err := linkedlist.New[int]()
		assert.NoError(b, err)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		n := l.Head.Next
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.InsertBefore(n, i)
		}
	})

	b.Run("custom safe", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		n := l.Head.Next
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.InsertBefore(n, i)
		}
	})

	b.Run("custom safe conc", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		n := l.Head.Next
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				l.InsertBefore(n, 0)
			}
		})
	})

	b.Run("std unsafe", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		e := l.Front().Next()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.InsertBefore(e, &golist.Element{Value: i})
		}
	})

	var mu sync.Mutex
	b.Run("std safe", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		e := l.Front().Next()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mu.Lock()
			l.InsertBefore(e, &golist.Element{Value: i})
			mu.Unlock()
		}
	})

	b.Run("std safe conc", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		e := l.Front().Next()
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				mu.Lock()
				l.InsertBefore(e, &golist.Element{Value: 0})
				mu.Unlock()
			}
		})
	})
}

func BenchmarkList_PopBack(b *testing.B) {
	b.Run("custom unsafe", func(b *testing.B) {
		l, err := linkedlist.New[int]()
		assert.NoError(b, err)
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.PopBack()
		}
	})

	b.Run("custom safe", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.PopBack()
		}
	})

	b.Run("custom safe conc", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				l.PopBack()
			}
		})
	})

	b.Run("std unsafe", func(b *testing.B) {
		l := golist.New()
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Remove(l.Back())
		}
	})

	var mu sync.Mutex
	b.Run("std safe", func(b *testing.B) {
		l := golist.New()
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mu.Lock()
			l.Remove(l.Back())
			mu.Unlock()
		}
	})

	b.Run("std safe conc", func(b *testing.B) {
		l := golist.New()
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				mu.Lock()
				l.Remove(l.Back())
				mu.Unlock()
			}
		})
	})
}

func BenchmarkList_PopFront(b *testing.B) {
	b.Run("custom unsafe", func(b *testing.B) {
		l, err := linkedlist.New[int]()
		assert.NoError(b, err)
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.PopFront()
		}
	})

	b.Run("custom safe", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.PopFront()
		}
	})

	b.Run("custom safe conc", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				l.PopFront()
			}
		})
	})

	b.Run("std unsafe", func(b *testing.B) {
		l := golist.New()
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Remove(l.Front())
		}
	})

	var mu sync.Mutex
	b.Run("std safe", func(b *testing.B) {
		l := golist.New()
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mu.Lock()
			l.Remove(l.Front())
			mu.Unlock()
		}
	})

	b.Run("std safe conc", func(b *testing.B) {
		l := golist.New()
		for i := 0; i < b.N; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				mu.Lock()
				l.Remove(l.Front())
				mu.Unlock()
			}
		})
	})
}

func BenchmarkList_Find(b *testing.B) {
	equal := func(a, b int) bool { return a == b }
	_ = equal

	b.Run("unsafe", func(b *testing.B) {
		l, err := linkedlist.New[int]()
		assert.NoError(b, err)
		for i := 0; i < 1000; i++ {
			l.PushBack(i)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Find(i%1000, util.Equal)
		}
	})

	b.Run("safe", func(b *testing.B) {
		l, err := linkedlist.New(linkedlist.WithSafe[int]())
		assert.NoError(b, err)
		for i := 0; i < 1000; i++ {
			l.PushBack(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Find(i%1000, util.Equal)
		}
	})
}

func BenchmarkList_Reverse(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	b.Run("unsafe", func(b *testing.B) {
		for _, size := range sizes {
			b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
				l, err := linkedlist.New[int]()
				assert.NoError(b, err)
				for i := 0; i < size; i++ {
					l.PushBack(i)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					l.Reverse()
				}
			})
		}
	})

	b.Run("safe", func(b *testing.B) {
		for _, size := range sizes {
			b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
				l, err := linkedlist.New(linkedlist.WithSafe[int]())
				assert.NoError(b, err)
				for i := 0; i < size; i++ {
					l.PushBack(i)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					l.Reverse()
				}
			})
		}
	})
}

func BenchmarkList_Merge(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	b.Run("unsafe", func(b *testing.B) {
		for _, n := range sizes {
			slices := make([]int, 0, n)
			for i := 0; i < n; i++ {
				slices = append(slices, i)
			}
			l1, err := linkedlist.NewFromSlice(slices)
			assert.NoError(b, err)
			l2, err := linkedlist.NewFromSlice(slices)
			assert.NoError(b, err)
			l1.Merge(l2)
		}
	})

	b.Run("safe", func(b *testing.B) {
		for _, n := range sizes {
			slices := make([]int, 0, n)
			for i := 0; i < n; i++ {
				slices = append(slices, i)
			}
			l1, err := linkedlist.NewFromSlice(slices, linkedlist.WithSafe[int]())
			assert.NoError(b, err)
			l2, err := linkedlist.NewFromSlice(slices, linkedlist.WithSafe[int]())
			assert.NoError(b, err)
			l1.Merge(l2)
		}
	})
}

func BenchmarkList_MergeSorted(b *testing.B) {
	sizes := []int{10, 20, 30}
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	b.Run("unsafe", func(b *testing.B) {
		for _, size := range sizes {
			slices := make([]int, 0, size)
			r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
			for range size {
				slices = append(slices, r.IntN(size))
			}
			l1, err := linkedlist.NewFromSlice(slices)
			assert.NoError(b, err)
			l2, err := linkedlist.NewFromSlice(slices)
			assert.NoError(b, err)
			b.ResetTimer()
			l1.Merge(l2)
		}
	})

	b.Run("safe", func(b *testing.B) {
		for _, size := range sizes {
			slices := make([]int, 0, size)
			r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
			for range size {
				slices = append(slices, r.IntN(size))
			}
			l1, err := linkedlist.NewFromSlice(slices, linkedlist.WithSafe[int]())
			assert.NoError(b, err)
			l2, err := linkedlist.NewFromSlice(slices, linkedlist.WithSafe[int]())
			assert.NoError(b, err)
			b.ResetTimer()
			l1.MergeSorted(l2, cmp)
		}
	})
}
