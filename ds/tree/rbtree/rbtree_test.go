package rbtree_test

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/forbearing/golib/ds/tree/rbtree"
	"github.com/stretchr/testify/assert"
)

func intComparator(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func TestRedBlackTree_New(t *testing.T) {
	// 测试 New
	tree, err := rbtree.New[int, string](intComparator)
	assert.NoError(t, err)
	assert.NotNil(t, tree)
	assert.True(t, tree.IsEmpty())
	assert.Equal(t, 0, tree.Size())

	// 测试 NewWithOrderedKeys
	orderedTree, err := rbtree.NewWithOrderedKeys[int, string]()
	assert.NoError(t, err)
	assert.NotNil(t, orderedTree)
	assert.True(t, orderedTree.IsEmpty())
	assert.Equal(t, 0, orderedTree.Size())

	// 测试 NewFromMap
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	treeFromMap, err := rbtree.NewFromMap(m, intComparator)
	assert.NoError(t, err)
	assert.NotNil(t, treeFromMap)
	assert.False(t, treeFromMap.IsEmpty())
	assert.Equal(t, len(m), treeFromMap.Size())

	// 测试 NewFromSlice
	slice := []string{"zero", "one", "two", "three"}
	sliceTree, err := rbtree.NewFromSlice(slice)
	assert.NoError(t, err)
	assert.NotNil(t, sliceTree)
	assert.False(t, sliceTree.IsEmpty())
	assert.Equal(t, len(slice), sliceTree.Size())

	fmt.Println(sliceTree)

	// 验证 map 中的 key 都存在于树中
	for k, v := range m {
		val, found := treeFromMap.Get(k)
		assert.True(t, found)
		assert.Equal(t, v, val)
	}

	// 测试 NewFromMapWithOrderedKeys
	orderedTreeFromMap, err := rbtree.NewFromMapWithOrderedKeys(m)
	assert.NoError(t, err)
	assert.NotNil(t, orderedTreeFromMap)
	assert.False(t, orderedTreeFromMap.IsEmpty())
	assert.Equal(t, len(m), orderedTreeFromMap.Size())

	// 验证 map 中的 key 都存在于树中
	for k, v := range m {
		val, found := orderedTreeFromMap.Get(k)
		assert.True(t, found)
		assert.Equal(t, v, val)
	}
}

func TestRedBlackTree_BasicOperations(t *testing.T) {
	tree, err := rbtree.New[int, string](intComparator)
	assert.NoError(t, err)

	// 测试空树
	assert.True(t, tree.IsEmpty())
	assert.Equal(t, 0, tree.Size())

	// 插入元素
	tree.Put(10, "ten")
	tree.Put(20, "twenty")
	tree.Put(5, "five")

	assert.False(t, tree.IsEmpty())
	assert.Equal(t, 3, tree.Size())

	// 获取元素
	val, found := tree.Get(10)
	assert.True(t, found)
	assert.Equal(t, "ten", val)

	// 获取不存在的元素
	val, found = tree.Get(100)
	assert.False(t, found)
	assert.Equal(t, "", val)

	// 删除元素
	tree.Delete(10)
	_, found = tree.Get(10)
	assert.False(t, found)
	assert.Equal(t, 2, tree.Size())

	// 删除不存在的元素
	tree.Delete(100) // 不应引发错误
	assert.Equal(t, 2, tree.Size())
}

func TestRedBlackTree_MinMax(t *testing.T) {
	tree, err := rbtree.New[int, string](intComparator)
	assert.NoError(t, err)

	// 测试空树
	assert.Nil(t, tree.Min())
	assert.Nil(t, tree.Max())

	// 插入元素
	tree.Put(15, "fifteen")
	tree.Put(10, "ten")
	tree.Put(20, "twenty")
	tree.Put(5, "five")

	// 最小值
	assert.Equal(t, 5, tree.Min().Key)
	assert.Equal(t, "five", tree.Min().Value)

	// 最大值
	assert.Equal(t, 20, tree.Max().Key)
	assert.Equal(t, "twenty", tree.Max().Value)
}

func TestRedBlackTree_FloorCeiling(t *testing.T) {
	tree, err := rbtree.New[int, string](intComparator)
	assert.NoError(t, err)
	tree.Put(10, "ten")
	tree.Put(20, "twenty")
	tree.Put(30, "thirty")

	// Floor 测试
	node, found := tree.Floor(25) // 应返回 20
	assert.True(t, found)
	assert.Equal(t, 20, node.Key)

	node, found = tree.Floor(10) // 应返回 10
	assert.True(t, found)
	assert.Equal(t, 10, node.Key)

	node, found = tree.Floor(5) // 不存在
	assert.False(t, found)
	assert.Nil(t, node)

	// Ceiling 测试
	node, found = tree.Ceiling(25) // 应返回 30
	assert.True(t, found)
	assert.Equal(t, 30, node.Key)

	node, found = tree.Ceiling(20) // 应返回 20
	assert.True(t, found)
	assert.Equal(t, 20, node.Key)

	node, found = tree.Ceiling(35) // 不存在
	assert.False(t, found)
	assert.Nil(t, node)
}

func TestRedBlackTree_Clear(t *testing.T) {
	tree, err := rbtree.New[int, string](intComparator)
	assert.NoError(t, err)

	tree.Put(1, "one")
	tree.Put(2, "two")

	assert.Equal(t, 2, tree.Size())

	tree.Clear()

	assert.True(t, tree.IsEmpty())
	assert.Equal(t, 0, tree.Size())
}

func TestRedBlackTree_KeysValues(t *testing.T) {
	tree, err := rbtree.New[int, string](intComparator)
	assert.NoError(t, err)
	tree.Put(3, "three")
	tree.Put(1, "one")
	tree.Put(2, "two")

	// Keys 应按排序顺序返回
	expectedKeys := []int{1, 2, 3}
	assert.Equal(t, expectedKeys, tree.Keys())

	// Values 应按 in-order 顺序返回
	expectedValues := []string{"one", "two", "three"}
	assert.Equal(t, expectedValues, tree.Values())
}

func TestRedBlackTree_Traversals(t *testing.T) {
	tree, err := rbtree.New[int, string](intComparator)
	assert.NoError(t, err)
	tree.Put(10, "ten")
	tree.Put(5, "five")
	tree.Put(15, "fifteen")
	tree.Put(3, "three")
	tree.Put(7, "seven")

	// Preorder: 根 → 左 → 右
	expectedPreorder := []int{10, 5, 3, 7, 15}
	var preorder []int
	for n := range tree.Preorder() {
		preorder = append(preorder, n.Key)
	}
	assert.Equal(t, expectedPreorder, preorder)

	// Inorder: 左 → 根 → 右 (排序)
	expectedInorder := []int{3, 5, 7, 10, 15}
	var inorder []int
	for n := range tree.Inorder() {
		inorder = append(inorder, n.Key)
	}
	assert.Equal(t, expectedInorder, inorder)

	// Postorder: 左 → 右 → 根
	expectedPostorder := []int{3, 7, 5, 15, 10}
	var postorder []int
	for n := range tree.Postorder() {
		postorder = append(postorder, n.Key)
	}
	assert.Equal(t, expectedPostorder, postorder)

	// LevelOrder: 层级遍历
	expectedLevelOrder := []int{10, 5, 15, 3, 7}
	var levelOrder []int
	for n := range tree.LevelOrder() {
		levelOrder = append(levelOrder, n.Key)
	}
	assert.Equal(t, expectedLevelOrder, levelOrder)
}

func TestRedBlackTree_String(t *testing.T) {
	fmt.Println("=== Test Red-Black Tree Visualization ===")

	// 1️⃣ 创建一个 int -> string 的红黑树
	tree, err := rbtree.NewWithOrderedKeys(rbtree.WithColorfulString[int, string]())
	assert.NoError(t, err)
	tree.Put(10, "ten")
	tree.Put(20, "twenty")
	tree.Put(30, "thirty")
	tree.Put(15, "fifteen")
	tree.Put(25, "twenty-five")
	tree.Put(5, "five")
	tree.Put(1, "one")
	tree.Put(7, "seven")
	tree.Put(40, "forty")
	tree.Put(50, "fifty")

	fmt.Println("\n🔹 Red-Black Tree (int -> string):")
	fmt.Println(tree.String())

	// 2️⃣ 创建一个 string -> int 的红黑树
	treeStr, err := rbtree.NewWithOrderedKeys(rbtree.WithColorfulString[string, int]())
	assert.NoError(t, err)
	treeStr.Put("banana", 10)
	treeStr.Put("apple", 5)
	treeStr.Put("cherry", 20)
	treeStr.Put("date", 15)
	treeStr.Put("fig", 25)
	treeStr.Put("grape", 8)
	treeStr.Put("lemon", 30)

	fmt.Println("\n🔹 Red-Black Tree (string -> int):")
	fmt.Println(treeStr.String())

	// 3️⃣ 创建一个 float64 -> string 的红黑树
	treeFloat, err := rbtree.NewWithOrderedKeys(rbtree.WithColorfulString[float64, string]())
	assert.NoError(t, err)

	treeFloat.Put(3.14, "pi")
	treeFloat.Put(2.71, "e")
	treeFloat.Put(1.61, "golden ratio")
	treeFloat.Put(1.41, "sqrt(2)")
	treeFloat.Put(2.23, "sqrt(5)")

	fmt.Println("\n🔹 Red-Black Tree (float64 -> string):")
	fmt.Println(treeFloat.String())

	tt, _ := rbtree.NewWithOrderedKeys(rbtree.WithColorfulString[float64, float64]())
	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
	for range 10000 {
		v := r.Float64()
		tt.Put(v, v)
	}
	fmt.Println(tt.Size(), tt.BlackCount(), tt.RedCount(), tt.LeafCount(), tt.MaxDepth(), tt.MinDepth())
}
