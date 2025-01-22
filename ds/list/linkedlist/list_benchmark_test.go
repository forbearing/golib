package linkedlist_test

import (
	golist "container/list"
	"fmt"
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/forbearing/golib/ds/list/linkedlist"
)

func prepareLinkedList(b *testing.B, size int, safe bool) *linkedlist.List[int] {
	slices := make([]int, 0, size)
	for i := range size {
		slices = append(slices, i)
	}

	var list *linkedlist.List[int]
	var err error
	if safe {
		list, err = linkedlist.NewFromSlice(slices, linkedlist.WithSafe[int]())
	} else {
		list, err = linkedlist.NewFromSlice(slices)
	}
	if err != nil {
		b.Fatalf("failed to create list: %v", err)
	}
	return list
}

func prepareStdList(_ *testing.B, size int) *golist.List {
	list := golist.New()
	for i := range size {
		list.PushBack(i)
	}
	return list
}

func benchmarkPushBack(b *testing.B, size int) {
	var mu sync.Mutex

	b.Run("unsafe custom", func(b *testing.B) {
		l := prepareLinkedList(b, 0, false)
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.PushBack(0)
			}
		}
	})
	b.Run("unsafe std", func(b *testing.B) {
		l := prepareStdList(b, 0)
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.PushBack(0)
			}
		}
	})

	b.Run("safe custom", func(b *testing.B) {
		l := prepareLinkedList(b, 0, true)
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.PushBack(0)
			}
		}
	})
	b.Run("safe std", func(b *testing.B) {
		l := prepareStdList(b, 0)
		b.ResetTimer()
		for range b.N {
			for range size {
				mu.Lock()
				_ = l.PushBack(0)
				mu.Unlock()
			}
		}
	})

	b.Run("safe conc custom", func(b *testing.B) {
		l := prepareLinkedList(b, 0, true)
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for range size {
					_ = l.PushBack(0)
				}
			}
		})
	})

	b.Run("safe conc std", func(b *testing.B) {
		l := prepareStdList(b, 0)
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for range size {
					mu.Lock()
					_ = l.PushBack(0)
					mu.Unlock()
				}
			}
		})
	})
}

func BenchmarkLinkedList_PushBack(b *testing.B) {
	for _, size := range []int{10, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkPushBack(b, size)
		})
	}
}

func benchmarkPushFront(b *testing.B, size int) {
	var mu sync.Mutex

	b.Run("unsafe custom", func(b *testing.B) {
		l := prepareLinkedList(b, 0, false)
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.PushFront(0)
			}
		}
	})
	b.Run("unsafe std", func(b *testing.B) {
		l := prepareStdList(b, 0)
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.PushFront(0)
			}
		}
	})

	b.Run("safe custom", func(b *testing.B) {
		l := prepareLinkedList(b, 0, true)
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.PushFront(0)
			}
		}
	})
	b.Run("safe std", func(b *testing.B) {
		l := prepareStdList(b, 0)
		b.ResetTimer()
		for range b.N {
			for range size {
				mu.Lock()
				_ = l.PushFront(0)
				mu.Unlock()
			}
		}
	})

	b.Run("safe conc custom", func(b *testing.B) {
		l := prepareLinkedList(b, 0, true)
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for range size {
					_ = l.PushFront(0)
				}
			}
		})
	})
	b.Run("safe conc std", func(b *testing.B) {
		l := prepareStdList(b, 0)
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for range size {
					mu.Lock()
					_ = l.PushFront(0)
					mu.Unlock()
				}
			}
		})
	})
}

func BenchmarkLinkedList_PushFront(b *testing.B) {
	for _, size := range []int{10, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkPushFront(b, size)
		})
	}
}

func benchmarkInsertAfter(b *testing.B, size int) {
	var mu sync.Mutex

	b.Run("unsafe custom", func(b *testing.B) {
		l := prepareLinkedList(b, 0, false)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		n := l.Head.Next
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.InsertAfter(n, 0)
			}
		}
	})
	b.Run("unsafe std", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		e := l.Front().Next()
		b.ResetTimer()
		for range b.N {
			for range size {
				l.InsertAfter(e, &golist.Element{Value: 0})
			}
		}
	})

	b.Run("safe custom", func(b *testing.B) {
		l := prepareLinkedList(b, 0, true)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		n := l.Head.Next
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.InsertAfter(n, 0)
			}
		}
	})
	b.Run("safe std", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		e := l.Front().Next()
		b.ResetTimer()
		for range b.N {
			for range size {
				mu.Lock()
				l.InsertAfter(e, &golist.Element{Value: 0})
				mu.Unlock()
			}
		}
	})

	b.Run("safe conc custom", func(b *testing.B) {
		l := prepareLinkedList(b, 0, true)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		n := l.Head.Next
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for range size {
					_ = l.InsertAfter(n, 0)
				}
			}
		})
	})
	b.Run("safe conc std", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		e := l.Front().Next()
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for range size {
					mu.Lock()
					l.InsertAfter(e, &golist.Element{Value: 0})
					mu.Unlock()
				}
			}
		})
	})
}

func BenchmarkLinkedListInsertAfter(b *testing.B) {
	for _, size := range []int{10, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkInsertAfter(b, size)
		})
	}
}

func benchmarkInsertBefore(b *testing.B, size int) {
	var mu sync.Mutex

	b.Run("unsafe custom", func(b *testing.B) {
		l := prepareLinkedList(b, 0, false)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		n := l.Head.Next
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.InsertBefore(n, 0)
			}
		}
	})
	b.Run("unsafe std", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		e := l.Front().Next()
		b.ResetTimer()
		for range b.N {
			for range size {
				l.InsertAfter(e, &golist.Element{Value: 0})
			}
		}
	})

	b.Run("safe custom", func(b *testing.B) {
		l := prepareLinkedList(b, 0, true)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		n := l.Head.Next
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.InsertBefore(n, 0)
			}
		}
	})
	b.Run("safe std", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		e := l.Front().Next()
		b.ResetTimer()
		for range b.N {
			for range size {
				mu.Lock()
				l.InsertAfter(e, &golist.Element{Value: 0})
				mu.Unlock()
			}
		}
	})

	b.Run("safe conc custom", func(b *testing.B) {
		l := prepareLinkedList(b, 0, true)
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		n := l.Head.Next
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for range size {
					_ = l.InsertBefore(n, 0)
				}
			}
		})
	})
	b.Run("safe conc std", func(b *testing.B) {
		l := golist.New()
		l.PushBack(0)
		l.PushBack(1)
		l.PushBack(2)
		e := l.Front().Next()
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for range size {
					mu.Lock()
					l.InsertAfter(e, &golist.Element{Value: 0})
					mu.Unlock()
				}
			}
		})
	})
}

func BenchmarkLinkedListInsertBefore(b *testing.B) {
	for _, size := range []int{10, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkInsertBefore(b, size)
		})
	}
}

func benchmarkPopBack(b *testing.B, size int) {
	var mu sync.Mutex

	b.Run("unsafe custom", func(b *testing.B) {
		l := prepareLinkedList(b, b.N*size, false)
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.PopBack()
			}
		}
	})
	b.Run("unsafe std", func(b *testing.B) {
		l := prepareStdList(b, b.N*size)
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.Remove(l.Back())
			}
		}
	})

	b.Run("safe custom", func(b *testing.B) {
		l := prepareLinkedList(b, b.N*size, true)
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.PopBack()
			}
		}
	})
	b.Run("safe std", func(b *testing.B) {
		l := prepareStdList(b, b.N*size)
		b.ResetTimer()
		for range b.N {
			for range size {
				mu.Lock()
				_ = l.Remove(l.Back())
				mu.Unlock()
			}
		}
	})

	// b.Run("safe conc custom", func(b *testing.B) {
	// 	l := prepareLinkedList(b, size, true)
	// 	b.ResetTimer()
	// 	b.RunParallel(func(p *testing.PB) {
	// 		for p.Next() {
	// 			for range size {
	// 				_ = l.PopBack()
	// 			}
	// 		}
	// 	})
	// })
	// b.Run("safe conc std", func(b *testing.B) {
	// 	l := prepareStdList(b, size)
	// 	b.ResetTimer()
	// 	b.RunParallel(func(p *testing.PB) {
	// 		for p.Next() {
	// 			for range size {
	// 				mu.Lock()
	// 				_ = l.Remove(l.Back())
	// 				mu.Unlock()
	// 			}
	// 		}
	// 	})
	// })
}

func BenchmarkLinkedList_PopBack(b *testing.B) {
	for _, size := range []int{100} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkPopBack(b, size)
		})
	}
}

func benchmarkPopFront(b *testing.B, size int) {
	var mu sync.Mutex

	b.Run("unsafe custom", func(b *testing.B) {
		l := prepareLinkedList(b, b.N*size, false)
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.PopFront()
			}
		}
	})
	b.Run("unsafe std", func(b *testing.B) {
		l := prepareStdList(b, b.N*size)
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.Remove(l.Front())
			}
		}
	})

	b.Run("safe custom", func(b *testing.B) {
		l := prepareLinkedList(b, b.N*size, true)
		b.ResetTimer()
		for range b.N {
			for range size {
				_ = l.PopFront()
			}
		}
	})
	b.Run("safe std", func(b *testing.B) {
		l := prepareStdList(b, b.N*size)
		b.ResetTimer()
		for range b.N {
			for range size {
				mu.Lock()
				_ = l.Remove(l.Front())
				mu.Unlock()
			}
		}
	})

	// b.Run("safe conc custom", func(b *testing.B) {
	// 	l := prepareLinkedList(b, size, true)
	// 	b.ResetTimer()
	// 	b.RunParallel(func(p *testing.PB) {
	// 		for p.Next() {
	// 			for range size {
	// 				_ = l.PopFront()
	// 			}
	// 		}
	// 	})
	// })
	// b.Run("safe conc std", func(b *testing.B) {
	// 	l := prepareStdList(b, size)
	// 	b.ResetTimer()
	// 	b.RunParallel(func(p *testing.PB) {
	// 		for p.Next() {
	// 			for range size {
	// 				mu.Lock()
	// 				_ = l.Remove(l.Front())
	// 				mu.Unlock()
	// 			}
	// 		}
	// 	})
	// })
}

func BenchmarkLinkedList_PopFront(b *testing.B) {
	for _, size := range []int{100} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkPopFront(b, size)
		})
	}
}

func stdListFind(list *golist.List, v any, equal func(int, int) bool) (_v any) {
	for e := list.Front(); e != nil; e = e.Next() {
		if equal(v.(int), e.Value.(int)) {
			return v
		}
	}
	return
}

func benchmarkFind(b *testing.B, size int) {
	var mu sync.Mutex
	equal := func(a, b int) bool { return a == b }

	b.Run("unsafe custom", func(b *testing.B) {
		l := prepareLinkedList(b, size, false)
		b.ResetTimer()
		for range b.N {
			for i := range size {
				_ = l.Find(i, equal)
			}
		}
	})
	b.Run("unsafe std", func(b *testing.B) {
		l := prepareStdList(b, size)
		b.ResetTimer()
		for range b.N {
			for i := range size {
				_ = stdListFind(l, i, equal)
			}
		}
	})

	b.Run("safe custom", func(b *testing.B) {
		l := prepareLinkedList(b, size, true)
		b.ResetTimer()
		for range b.N {
			for i := range size {
				_ = l.Find(i, equal)
			}
		}
	})
	b.Run("safe std", func(b *testing.B) {
		l := prepareStdList(b, size)
		b.ResetTimer()
		for range b.N {
			for i := range size {
				mu.Lock()
				_ = stdListFind(l, i, equal)
				mu.Unlock()
			}
		}
	})

	b.Run("safe conc custom", func(b *testing.B) {
		l := prepareLinkedList(b, size, true)
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for i := range size {
					_ = l.Find(i, equal)
				}
			}
		})
	})
	b.Run("safe conc std", func(b *testing.B) {
		l := prepareStdList(b, size)
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for i := range size {
					mu.Lock()
					_ = stdListFind(l, i, equal)
					mu.Unlock()
				}
			}
		})
	})
}

func BenchmarkLinkedList_Find(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkFind(b, size)
		})
	}
}

func benchmarkReverse(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		l := prepareLinkedList(b, size, false)
		b.ResetTimer()
		for range b.N {
			l.Reverse()
		}
	})

	b.Run("safe", func(b *testing.B) {
		l := prepareLinkedList(b, size, true)
		b.ResetTimer()
		for range b.N {
			l.Reverse()
		}
	})

	b.Run("safe conc", func(b *testing.B) {
		l := prepareLinkedList(b, size, true)
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				l.Reverse()
			}
		})
	})
}

func BenchmarkLinkedList_Reverse(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkReverse(b, size)
		})
	}
}

func benchmarkMerge(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		l1 := prepareLinkedList(b, size, false)
		b.ResetTimer()
		for range b.N {
			l2 := prepareLinkedList(b, size, false)
			l1.Merge(l2)
		}
	})

	b.Run("safe", func(b *testing.B) {
		l1 := prepareLinkedList(b, size, true)
		b.ResetTimer()
		for range b.N {
			l2 := prepareLinkedList(b, size, true)
			l1.Merge(l2)
		}
	})
}

func BenchmarkLinkedList_Merge(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkMerge(b, size)
		})
	}
}

func benchmarkMergeSorted(b *testing.B, size int) {
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
		slices := make([]int, 0, size)
		r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
		for range size {
			slices = append(slices, r.IntN(size))
		}
		l1, err := linkedlist.NewFromSlice(slices)
		if err != nil {
			b.Fatalf("failed to create list: %v", err)
		}
		b.ResetTimer()
		for range b.N {
			b.StopTimer()
			l2, err := linkedlist.NewFromSlice(slices)
			if err != nil {
				b.Fatalf("failed to create list: %v", err)
			}
			b.StartTimer()
			l1.MergeSorted(l2, cmp)
		}
	})

	b.Run("safe", func(b *testing.B) {
		slices := make([]int, 0, size)
		r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
		for range size {
			slices = append(slices, r.IntN(size))
		}
		l1, err := linkedlist.NewFromSlice(slices, linkedlist.WithSafe[int]())
		if err != nil {
			b.Fatalf("failed to create list: %v", err)
		}
		b.ResetTimer()
		for range b.N {
			b.StopTimer()
			l2, err := linkedlist.NewFromSlice(slices, linkedlist.WithSafe[int]())
			if err != nil {
				b.Fatalf("failed to create list: %v", err)
			}
			b.StartTimer()
			l1.MergeSorted(l2, cmp)
		}
	})
}

func BenchmarkLinkedList_MergeSorted(b *testing.B) {
	// go test -bench 'MergeSorted' ./ds/list/linkedlist/ -benchtime=1000x
	for _, size := range []int{10, 100} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkMergeSorted(b, size)
		})
	}
}
