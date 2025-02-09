package avltree_test

import (
	"fmt"
	"testing"

	"github.com/forbearing/golib/ds/tree/avltree"
	"github.com/stretchr/testify/assert"
)

func intCmp(a, b int) int {
	return a - b
}

func TestAVLTree_New(t *testing.T) {
	tt, err := avltree.New[int, string](intCmp)
	assert.NoError(t, err)
	assert.NotNil(t, tt)
	assert.Equal(t, 0, tt.Size())

	m := map[int]string{
		10: "ten",
		20: "twenty",
		5:  "five",
	}
	tt, err = avltree.NewFromMap(m, intCmp)
	assert.NoError(t, err)
	assert.NotNil(t, tt)
	assert.Equal(t, 3, tt.Size())

	tt, err = avltree.NewFromMapWithOrderedKeys(m)
	assert.NoError(t, err)
	assert.NotNil(t, tt)
	assert.Equal(t, 3, tt.Size())

	slice := []int{10, 20, 5}
	tt2, err := avltree.NewFromSlice(slice)
	assert.NoError(t, err)
	assert.NotNil(t, tt2)
	assert.Equal(t, 3, tt2.Size())
}

func TestAVLTree_PutAndGet(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(10, "ten")
	tree.Put(20, "twenty")
	tree.Put(5, "five")

	val, found := tree.Get(10)
	assert.True(t, found)
	assert.Equal(t, "ten", val)

	val, found = tree.Get(20)
	assert.True(t, found)
	assert.Equal(t, "twenty", val)

	val, found = tree.Get(5)
	assert.True(t, found)
	assert.Equal(t, "five", val)

	// 查找不存在的键
	_, found = tree.Get(100)
	assert.False(t, found)
}

func TestAVLTree_PutDuplicateKey(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(10, "ten")
	tree.Put(10, "updated-ten")

	val, found := tree.Get(10)
	assert.True(t, found)
	assert.Equal(t, "updated-ten", val)

	assert.Equal(t, 1, tree.Size())
}

func TestAVLTree_Delete(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(10, "ten")
	tree.Put(20, "twenty")
	tree.Put(5, "five")

	tree.Delete(10)
	_, found := tree.Get(10)
	assert.False(t, found)

	tree.Delete(100)
	assert.Equal(t, 2, tree.Size())

	tree.Delete(5)
	tree.Delete(20)
	assert.True(t, tree.IsEmpty())
}

func TestAVLTree_SizeHeight(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	assert.Equal(t, 0, tree.Size())
	assert.Equal(t, 0, tree.Height())

	tree.Put(10, "ten")
	tree.Put(20, "twenty")
	tree.Put(5, "five")

	assert.Equal(t, 3, tree.Size())
	assert.Greater(t, tree.Height(), 0)
}

func TestAVLTree_KeysAndValues(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(10, "ten")
	tree.Put(20, "twenty")
	tree.Put(5, "five")

	keys := tree.Keys()
	values := tree.Values()
	assert.Equal(t, []int{5, 10, 20}, keys)
	assert.Equal(t, []string{"five", "ten", "twenty"}, values)

	maxNode := tree.Max()
	assert.NotNil(t, maxNode)
	assert.Equal(t, 20, maxNode.Key)
	assert.Equal(t, "twenty", maxNode.Value)
}

func TestAVLTree_MinMax(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(10, "ten")
	tree.Put(20, "twenty")
	tree.Put(5, "five")

	minNode := tree.Min()
	assert.NotNil(t, minNode)
	assert.Equal(t, 5, minNode.Key)
	assert.Equal(t, "five", minNode.Value)

	maxNode := tree.Max()
	assert.NotNil(t, maxNode)
	assert.Equal(t, 20, maxNode.Key)
	assert.Equal(t, "twenty", maxNode.Value)
}

func TestAVLTree_Balance(t *testing.T) {
	tree, _ := avltree.New[int, int](intCmp)

	// 插入会导致不平衡的情况
	tree.Put(30, 30)
	tree.Put(20, 20)
	tree.Put(40, 40)
	tree.Put(10, 10) // 触发平衡调整
	tree.Put(25, 25) // 触发平衡调整
	tree.Put(35, 35) // 触发平衡调整
	tree.Put(50, 50) // 触发平衡调整

	assert.Equal(t, 7, tree.Size())
	assert.Greater(t, tree.Height(), 0)
}

func TestAVLTree_FloorAndCeiling(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(10, "ten")
	tree.Put(20, "twenty")
	tree.Put(30, "thirty")

	node, found := tree.Floor(25)
	assert.True(t, found)
	assert.Equal(t, 20, node.Key)

	node, found = tree.Ceiling(15)
	assert.True(t, found)
	assert.Equal(t, 20, node.Key)

	node, found = tree.Floor(5)
	assert.False(t, found)

	node, found = tree.Ceiling(35)
	assert.False(t, found)
}

func TestAVLTree_Clear(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(10, "ten")
	tree.Put(20, "twenty")

	tree.Clear()
	assert.True(t, tree.IsEmpty())
}

func TestAVLTree_Preorder(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(40, "forty")
	tree.Put(20, "twenty")
	tree.Put(60, "sixty")
	tree.Put(10, "ten")
	tree.Put(30, "thirty")
	tree.Put(50, "fifty")
	tree.Put(70, "seventy")

	expectedKeys := []int{40, 20, 10, 30, 60, 50, 70}
	expectedValues := []string{"forty", "twenty", "ten", "thirty", "sixty", "fifty", "seventy"}

	var keys []int
	var values []string
	tree.Preorder(func(k int, v string) {
		keys = append(keys, k)
		values = append(values, v)
	})

	assert.Equal(t, expectedKeys, keys)
	assert.Equal(t, expectedValues, values)
}

func TestAVLTree_PreorderChan(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(40, "forty")
	tree.Put(20, "twenty")
	tree.Put(60, "sixty")
	tree.Put(10, "ten")
	tree.Put(30, "thirty")
	tree.Put(50, "fifty")
	tree.Put(70, "seventy")

	expectedKeys := []int{40, 20, 10, 30, 60, 50, 70}
	expectedValues := []string{"forty", "twenty", "ten", "thirty", "sixty", "fifty", "seventy"}

	var keys []int
	var values []string
	for node := range tree.PreorderChan() {
		keys = append(keys, node.Key)
		values = append(values, node.Value)
	}

	assert.Equal(t, expectedKeys, keys)
	assert.Equal(t, expectedValues, values)
}

func TestAVLTree_Inorder(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(40, "forty")
	tree.Put(20, "twenty")
	tree.Put(60, "sixty")
	tree.Put(10, "ten")
	tree.Put(30, "thirty")
	tree.Put(50, "fifty")
	tree.Put(70, "seventy")

	expectedKeys := []int{10, 20, 30, 40, 50, 60, 70}
	expectedValues := []string{"ten", "twenty", "thirty", "forty", "fifty", "sixty", "seventy"}

	var keys []int
	var values []string
	tree.Inorder(func(k int, v string) {
		keys = append(keys, k)
		values = append(values, v)
	})

	assert.Equal(t, expectedKeys, keys)
	assert.Equal(t, expectedValues, values)
}

func TestAVLTree_InorderChan(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(40, "forty")
	tree.Put(20, "twenty")
	tree.Put(60, "sixty")
	tree.Put(10, "ten")
	tree.Put(30, "thirty")
	tree.Put(50, "fifty")
	tree.Put(70, "seventy")

	expectedKeys := []int{10, 20, 30, 40, 50, 60, 70}
	expectedValues := []string{"ten", "twenty", "thirty", "forty", "fifty", "sixty", "seventy"}

	var keys []int
	var values []string
	for node := range tree.InorderChan() {
		keys = append(keys, node.Key)
		values = append(values, node.Value)
	}

	assert.Equal(t, expectedKeys, keys)
	assert.Equal(t, expectedValues, values)
}

func TestAVLTree_Postorder(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(40, "forty")
	tree.Put(20, "twenty")
	tree.Put(60, "sixty")
	tree.Put(10, "ten")
	tree.Put(30, "thirty")
	tree.Put(50, "fifty")
	tree.Put(70, "seventy")

	expectedKeys := []int{10, 30, 20, 50, 70, 60, 40}
	expectedValues := []string{"ten", "thirty", "twenty", "fifty", "seventy", "sixty", "forty"}

	var keys []int
	var values []string
	tree.Postorder(func(k int, v string) {
		keys = append(keys, k)
		values = append(values, v)
	})

	assert.Equal(t, expectedKeys, keys)
	assert.Equal(t, expectedValues, values)
}

func TestAVLTree_PostorderChan(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(40, "forty")
	tree.Put(20, "twenty")
	tree.Put(60, "sixty")
	tree.Put(10, "ten")
	tree.Put(30, "thirty")
	tree.Put(50, "fifty")
	tree.Put(70, "seventy")

	expectedKeys := []int{10, 30, 20, 50, 70, 60, 40}
	expectedValues := []string{"ten", "thirty", "twenty", "fifty", "seventy", "sixty", "forty"}

	var keys []int
	var values []string
	for node := range tree.PostorderChan() {
		keys = append(keys, node.Key)
		values = append(values, node.Value)
	}

	assert.Equal(t, expectedKeys, keys)
	assert.Equal(t, expectedValues, values)
}

func TestAVLTree_LevelOrder(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(40, "forty")
	tree.Put(20, "twenty")
	tree.Put(60, "sixty")
	tree.Put(10, "ten")
	tree.Put(30, "thirty")
	tree.Put(50, "fifty")
	tree.Put(70, "seventy")

	expectedKeys := []int{40, 20, 60, 10, 30, 50, 70}
	expectedValues := []string{"forty", "twenty", "sixty", "ten", "thirty", "fifty", "seventy"}

	var keys []int
	var values []string
	tree.LevelOrder(func(k int, v string) {
		keys = append(keys, k)
		values = append(values, v)
	})

	assert.Equal(t, expectedKeys, keys)
	assert.Equal(t, expectedValues, values)
}

func TestAVLTree_LevelOrderChan(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(40, "forty")
	tree.Put(20, "twenty")
	tree.Put(60, "sixty")
	tree.Put(10, "ten")
	tree.Put(30, "thirty")
	tree.Put(50, "fifty")
	tree.Put(70, "seventy")

	expectedKeys := []int{40, 20, 60, 10, 30, 50, 70}
	expectedValues := []string{"forty", "twenty", "sixty", "ten", "thirty", "fifty", "seventy"}

	var keys []int
	var values []string
	for node := range tree.LevelOrderChan() {
		keys = append(keys, node.Key)
		values = append(values, node.Value)
	}

	assert.Equal(t, expectedKeys, keys)
	assert.Equal(t, expectedValues, values)
}

func TestAVLTree_String(t *testing.T) {
	tree, _ := avltree.New[int, string](intCmp)

	tree.Put(40, "forty")
	tree.Put(20, "twenty")
	tree.Put(60, "sixty")
	tree.Put(10, "ten")
	tree.Put(30, "thirty")
	tree.Put(50, "fifty")
	tree.Put(70, "seventy")

	fmt.Println(tree.String())
}
