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

func BenchmarkRedBlackTree_PreorderChan(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		for n := range t.PreorderChan() {
			_ = n
		}
	})
}

func BenchmarkRedBlackTree_Preorder(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		t.Preorder(func(f1, f2 float64) {
			_, _ = f1, f2
		})
	})
}

func BenchmarkRedBlackTree_InorderChan(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		for n := range t.InorderChan() {
			_ = n
		}
	})
}

func BenchmarkRedBlackTree_Inorder(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		t.Inorder(func(f1, f2 float64) {
			_, _ = f1, f2
		})
	})
}

func BenchmarkRedBlackTree_PostorderChan(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		for n := range t.PostorderChan() {
			_ = n
		}
	})
}

func BenchmarkRedBlackTree_Postorder(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		t.Postorder(func(f1, f2 float64) {
			_, _ = f1, f2
		})
	})
}

func BenchmarkRedBlackTree_LevelOrderChan(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		for n := range t.LevelOrderChan() {
			_ = n
		}
	})
}

func BenchmarkRedBlackTree_LevelOrder(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		t.LevelOrder(func(f1, f2 float64) {
			_, _ = f1, f2
		})
	})
}

func BenchmarkRedBlackTree_Min(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_ = t.Min()
	})
}

func BenchmarkRedBlackTree_Max(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_ = t.Max()
	})
}

func BenchmarkRedBlackTree_Floor(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_, _ = t.Floor(0)
	})
}

func BenchmarkRedBlackTree_Ceiling(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		_, _ = t.Ceiling(0)
	})
}

func BenchmarkRedBlackTree_Clear(b *testing.B) {
	benchmark(b, false, []int{10, 100000}, func(t *avltree.Tree[float64, float64]) {
		t.Clear()
	})
}
