package arraylist_test

import (
	"fmt"
	"testing"

	"github.com/forbearing/golib/ds/list/arraylist"
)

func prepareList(b *testing.B, size int, safe bool) *arraylist.List[int] {
	var list *arraylist.List[int]
	var err error

	if safe {
		list, err = arraylist.New(intEqual, arraylist.WithSafe[int]())
	} else {
		list, err = arraylist.New(intEqual)
	}
	if err != nil {
		b.Fatalf("failed to create list: %v", err)
	}
	for i := range size {
		list.Append(i)
	}

	return list
}

func benchmarkGet(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		if size <= 0 {
			size = 1
		}
		b.ResetTimer()
		for range b.N {
			for i := range size {
				_, _ = list.Get(i)
			}
		}
	})

	b.Run("safe", func(b *testing.B) {
		list := prepareList(b, size, true)
		if size <= 0 {
			size = 1
		}
		b.ResetTimer()
		for range b.N {
			for i := range size {
				_, _ = list.Get(i)
			}
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list := prepareList(b, size, true)
		if size <= 0 {
			size = 1
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for i := range size {
					_, _ = list.Get(i)
				}
			}
		})
	})
}

func BenchmarkArrayList_Get(b *testing.B) {
	for _, size := range []int{0, 1, 10, 100, 1000, 10000, 100000, 1000000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkGet(b, size)
		})
	}
}

func benchmarkAppend(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		list, err := arraylist.New(intEqual)
		if err != nil {
			b.Fatalf("failed to create list: %v", err)
		}
		b.ResetTimer()
		for range b.N {
			for i := range size {
				list.Append(i)
			}
		}
	})

	b.Run("safe", func(b *testing.B) {
		list, err := arraylist.New(intEqual, arraylist.WithSafe[int]())
		if err != nil {
			b.Fatalf("failed to create list: %v", err)
		}
		b.ResetTimer()
		for range b.N {
			for i := range size {
				list.Append(i)
			}
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list, err := arraylist.New(intEqual, arraylist.WithSafe[int]())
		if err != nil {
			b.Fatalf("failed to create list: %v", err)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for i := range size {
					list.Append(i)
				}
			}
		})
	})
}

func BenchmarkArrayList_Append(b *testing.B) {
	for _, size := range []int{0, 10, 100, 1000, 10000, 100000, 1000000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkAppend(b, size)
		})
	}
}

func BenchmarkArrayList_Insert(b *testing.B) {
	// go test -bench 'Insert' ./ds/list/arraylist -benchtime 100000x
	size := 1000
	b.Run("unsafe", func(b *testing.B) {
		list, _ := arraylist.New(intEqual)
		for i := 0; i < size; i++ {
			list.Append(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			list.Insert(i%size, i)
		}
	})

	b.Run("safe", func(b *testing.B) {
		list, _ := arraylist.New(intEqual, arraylist.WithSafe[int]())
		for i := 0; i < size; i++ {
			list.Append(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			list.Insert(i%size, i)
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list, _ := arraylist.New(intEqual, arraylist.WithSafe[int]())
		for i := 0; i < size; i++ {
			list.Append(i)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				list.Insert(10, 0)
			}
		})
	})
}

func benchmarkSet(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		if size <= 0 {
			size = 1
		}
		b.ResetTimer()
		for range b.N {
			for i := range size {
				list.Set(i%size, i)
			}
		}
	})

	b.Run("safe", func(b *testing.B) {
		list := prepareList(b, size, true)
		if size <= 0 {
			size = 1
		}
		b.ResetTimer()
		for range b.N {
			for i := range size {
				list.Set(i%size, i)
			}
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list := prepareList(b, size, true)
		if size <= 0 {
			size = 1
		}
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				for i := range size {
					list.Set(i%size, i)
				}
			}
		})
	})
}

func BenchmarkArrayList_Set(b *testing.B) {
	for _, size := range []int{0, 10, 100, 1000, 10000, 100000, 1000000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkSet(b, size)
		})
	}
}

func benchmarkRemove(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		b.ResetTimer()
		for range b.N {
			for i := range size {
				list.Remove(i)
			}
		}
	})

	b.Run("safe", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		for range b.N {
			for i := range size {
				list.Remove(i)
			}
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for i := range size {
					list.Remove(i)
				}
			}
		})
	})
}

func BenchmarkArrayList_Remove(b *testing.B) {
	for _, size := range []int{0, 10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkRemove(b, size)
		})
	}
}

func benchmarkRemoveAt(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		if size <= 0 {
			size = 1
		}
		b.ResetTimer()
		for range b.N {
			for i := range size {
				list.RemoveAt(i)
			}
		}
	})

	b.Run("safe", func(b *testing.B) {
		list := prepareList(b, size, true)
		if size <= 0 {
			size = 1
		}
		b.ResetTimer()
		for range b.N {
			for i := range size {
				list.RemoveAt(i)
			}
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list := prepareList(b, size, true)
		if size <= 0 {
			size = 1
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for range size {
					list.RemoveAt(0)
				}
			}
		})
	})
}

func BenchmarkArrayList_RemoveAt(b *testing.B) {
	for _, size := range []int{0, 10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkRemoveAt(b, size)
		})
	}
}

func BenchmarkArrayList_Clear(b *testing.B) {
	size := 1000
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		b.ResetTimer()
		for range b.N {
			list.Clear()
		}
	})

	b.Run("safe", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		for range b.N {
			list.Clear()
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				list.Clear()
			}
		})
	})
}

func benchmarkContains(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		b.ResetTimer()
		for range b.N {
			for i := range size {
				list.Contains(i)
			}
		}
	})

	b.Run("safe", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		for range b.N {
			for i := range size {
				list.Contains(i)
			}
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for i := range size {
					list.Contains(i)
				}
			}
		})
	})
}

func BenchmarkArrayList_Contains(b *testing.B) {
	for _, size := range []int{0, 10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("siz-%d", size), func(b *testing.B) {
			benchmarkContains(b, size)
		})
	}
}

func benchmarkValues(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		b.ResetTimer()
		for range b.N {
			_ = list.Values()
		}
	})

	b.Run("safe", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		for range b.N {
			_ = list.Values()
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for range size {
					_ = list.Values()
				}
			}
		})
	})
}

func BenchmarkArrayList_Values(b *testing.B) {
	for _, size := range []int{100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkValues(b, size)
		})
	}
}

func benchmarkIndexOf(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		b.ResetTimer()
		for range b.N {
			for i := range size {
				list.IndexOf(i)
			}
		}
	})

	b.Run("safe", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		for range b.N {
			for i := range size {
				list.IndexOf(i)
			}
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for i := range size {
					list.IndexOf(i)
				}
			}
		})
	})
}

func BenchmarkArrayList_IndexOf(b *testing.B) {
	for _, size := range []int{100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkIndexOf(b, size)
		})
	}
}

func benchmarkIsEmpty(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		b.ResetTimer()
		for range b.N {
			for range size {
				list.IsEmpty()
			}
		}
	})

	b.Run("safe", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		for range b.N {
			for range size {
				list.IsEmpty()
			}
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for range size {
					list.IsEmpty()
				}
			}
		})
	})
}

func BenchmarkArrayList_IsEmpty(b *testing.B) {
	for _, size := range []int{100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkIsEmpty(b, size)
		})
	}
}

func benchmarkLen(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		b.ResetTimer()
		for range b.N {
			for range size {
				list.Len()
			}
		}
	})

	b.Run("safe", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		for range b.N {
			for range size {
				list.Len()
			}
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for range size {
					list.Len()
				}
			}
		})
	})
}

func BenchmarkArrayList_Len(b *testing.B) {
	for _, size := range []int{100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkLen(b, size)
		})
	}
}

func benchmarkSort(b *testing.B, size int) {
	cmp := func(a, b int) int {
		return a - b
	}
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		b.ResetTimer()
		for range b.N {
			for range size {
				list.Sort(cmp)
			}
		}
	})

	b.Run("safe", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		for range b.N {
			for range size {
				list.Sort(cmp)
			}
		}
	})
}

func BenchmarkArrayList_Sort(b *testing.B) {
	for _, size := range []int{10, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkSort(b, size)
		})
	}
}

func benchmarkSwap(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		b.ResetTimer()
		for range b.N {
			for range size {
				list.Swap(0, size-1)
			}
		}
	})

	b.Run("safe", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		for range b.N {
			for range size {
				list.Swap(0, size-1)
			}
		}
	})

	b.Run("safe-conc", func(b *testing.B) {
		list := prepareList(b, size, true)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for range size {
					list.Swap(0, size-1)
				}
			}
		})
	})
}

func BenchmarkArrayList_Swap(b *testing.B) {
	for _, size := range []int{100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkSwap(b, size)
		})
	}
}

func benchmarkRange(b *testing.B, size int) {
	b.Run("unsafe", func(b *testing.B) {
		list := prepareList(b, size, false)
		b.ResetTimer()
		for range b.N {
			for range size {
				list.Range(func(v int) bool {
					_ = v
					return true
				})
			}
		}
	})

	b.Run("safe", func(b *testing.B) {
		l := prepareList(b, size, true)
		b.ResetTimer()
		for range b.N {
			for range size {
				l.Range(func(v int) bool {
					_ = v
					return true
				})
			}
		}
	})
}

func BenchmarkArrayList_Range(b *testing.B) {
	for _, size := range []int{1, 10, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchmarkRange(b, size)
		})
	}
}
