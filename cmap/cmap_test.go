package cmap

import "testing"

func init() {
	if err := Init(); err != nil {
		panic(err)
	}
}

type User struct {
	Name string
	Age  int
	Addr string
}

func Test(t *testing.T) {
	Cache[[]User]().Set("key", []User{{Name: "a", Age: 1}, {Name: "b", Age: 2}})
	Cache[[]*User]().Set("key", []*User{{Name: "a", Age: 1}, {Name: "b", Age: 2}})
	Cache[User]().Set("key", User{Name: "a", Age: 1, Addr: "a"})
	Cache[*User]().Set("key", &User{Name: "a", Age: 1, Addr: "a"})
	Cache[int]().Set("key", 1)
}

// func BenchmarkLruReadHeavy(b *testing.B) { benchmarkLru(b, 0.8) }
// func BenchmarkIntReadHeavy(b *testing.B) { benchmarkInt(b, 0.8) }
//
// func BenchmarkLruWriteHeavy(b *testing.B) { benchmarkLru(b, 0.2) }
// func BenchmarkIntWriteHeavy(b *testing.B) { benchmarkInt(b, 0.2) }
//
// func BenchmarkLruBalanced(b *testing.B) { benchmarkLru(b, 0.5) }
// func BenchmarkIntBalanced(b *testing.B) { benchmarkInt(b, 0.5) }
//
// func benchmarkLru(b *testing.B, readRatio float64) {
// 	size := 10000
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < size; j++ {
// 			if rand.Float64() < readRatio {
// 				Lru.Get(strconv.Itoa(rand.Intn(size)))
// 			} else {
// 				Lru.Set(strconv.Itoa(j), j)
// 			}
// 		}
// 	}
// }
//
// func benchmarkInt(b *testing.B, readRatio float64) {
// 	size := 10000
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < size; j++ {
// 			if rand.Float64() < readRatio {
// 				Int.Get(strconv.Itoa(rand.Intn(size)))
// 			} else {
// 				Int.Set(strconv.Itoa(j), j)
// 			}
// 		}
// 	}
// }
//
// func BenchmarkLruConcurrent(b *testing.B) {
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			key := strconv.Itoa(rand.Intn(1000))
// 			if rand.Float32() < 0.5 {
// 				Lru.Set(key, rand.Int())
// 			} else {
// 				Lru.Get(key)
// 			}
// 		}
// 	})
// }
//
// func BenchmarkIntConcurrent(b *testing.B) {
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			key := strconv.Itoa(rand.Intn(1000))
// 			if rand.Float32() < 0.5 {
// 				Int.Set(key, rand.Int())
// 			} else {
// 				Int.Get(key)
// 			}
// 		}
// 	})
// }
//
// //	func BenchmarkInt(b *testing.B) {
// //		for i := 0; i < b.N; i++ {
// //			Int.Set(strconv.Itoa(i), i)
// //		}
// //		for i := 0; i < b.N; i++ {
// //			Int.Get(strconv.Itoa(i))
// //		}
// //	}
// //
// //	func BenchmarkLru(b *testing.B) {
// //		for i := 0; i < b.N; i++ {
// //			Lru.Set(strconv.Itoa(i), i)
// //		}
// //		for i := 0; i < b.N; i++ {
// //			Lru.Get(strconv.Itoa(i))
// //		}
// //	}
