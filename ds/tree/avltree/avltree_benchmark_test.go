package avltree_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/forbearing/golib/ds/tree/avltree"
)

func createTree(b *testing.B, size int, safe bool) *avltree.Tree[float64, float64] {
	var t *avltree.Tree[float64, float64]
	var err error
	if safe {
		t, err = avltree.NewWithOrderedKeys(avltree.WithSafe[float64, float64]())
	} else {
		t, err = avltree.NewWithOrderedKeys[float64, float64]()
	}
	for i := range size {
		t.Put(float64(i), float64(i))
	}
	if err != nil {
		b.Fatalf("failed to create red-black tree: %v", err)
	}
	return t
}

func benchmark(b *testing.B, hasConcUnsafe bool, sizes []int, do func(t *avltree.Tree[float64, float64])) {
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			b.Run("single unsafe", func(b *testing.B) {
				t := createTree(b, size, false)
				b.ResetTimer()
				for range b.N {
					do(t)
				}
			})
			b.Run("single safe", func(b *testing.B) {
				t := createTree(b, size, true)
				b.ResetTimer()
				for range b.N {
					do(t)
				}
			})

			if hasConcUnsafe {
				b.Run("concur unsafe", func(b *testing.B) {
					t := createTree(b, size, false)
					b.ResetTimer()
					b.RunParallel(func(p *testing.PB) {
						for p.Next() {
							do(t)
						}
					})
				})
			}
			b.Run("concur safe", func(b *testing.B) {
				t := createTree(b, size, true)
				b.ResetTimer()
				b.RunParallel(func(p *testing.PB) {
					for p.Next() {
						do(t)
					}
				})
			})
		})
	}
}

func BenchmarkAVLTreePut(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		t.Put(float64(time.Now().UnixNano()), float64(time.Now().UnixNano()))
	})
}

func BenchmarkAVLTreeGet(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_, _ = t.Get(0)
	})
}

func BenchmarkAVLTreeDelete(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		t.Delete(0)
	})
}

func BenchmarkAVLTreeIsEmpty(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_ = t.IsEmpty()
	})
}

func BenchmarkAVLTreeSize(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_ = t.Size()
	})
}

func BenchmarkAVLTree_Clear(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		t.Clear()
	})
}

func BenchmarkAVLTreeKeys(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_ = t.Keys()
	})
}

func BenchmarkAVLTreeValues(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_ = t.Values()
	})
}

func BenchmarkAVLTree_Min(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_, _, _ = t.Min()
	})
}

func BenchmarkAVLTree_Max(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_, _, _ = t.Max()
	})
}

func BenchmarkAVLTree_Floor(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_, _, _ = t.Floor(0)
	})
}

func BenchmarkAVLTree_Ceiling(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_, _, _ = t.Ceiling(0)
	})
}

func BenchmarkAVLTree_PreOrder(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		t.PreOrder(func(f1, f2 float64) bool {
			_, _ = f1, f2
			return true
		})
	})
}

func BenchmarkAVLTree_InOrder(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		t.InOrder(func(f1, f2 float64) bool {
			_, _ = f1, f2
			return true
		})
	})
}

func BenchmarkAVLTree_PostOrder(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		t.PostOrder(func(f1, f2 float64) bool {
			_, _ = f1, f2
			return true
		})
	})
}

func BenchmarkAVLTree_LevelOrder(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		t.LevelOrder(func(f1, f2 float64) bool {
			_, _ = f1, f2
			return true
		})
	})
}

func BenchmarkAVLTree_String(b *testing.B) {
	benchmark(b, false, []int{1000}, func(t *avltree.Tree[float64, float64]) {
		_ = t.String()
	})
}
