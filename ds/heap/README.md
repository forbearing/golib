## **常见的堆数据结构**

堆是一种特殊的树，按照存储结构和维护规则，主要包括以下类型：

### **1. 二叉堆（Binary Heap）**

-   特点：

    -   **完全二叉树**（每一层从左到右填满，最后一层可以不满，但必须左对齐）。

    -   堆序性质（Heap Order Property）

        ：

        -   **最大堆（Max Heap）**：父节点值 **≥** 子节点值（根是最大值）。
        -   **最小堆（Min Heap）**：父节点值 **≤** 子节点值（根是最小值）。

-   应用：

    -   **优先队列**（Priority Queue）。
    -   **堆排序（Heap Sort）**。

### **2. 二项堆（Binomial Heap）**

-   特点：
    -   由 **二项树（Binomial Tree）** 组成，具有类似二进制表示的结构。
    -   **支持快速合并（Merge）操作**。
-   应用：
    -   **动态合并优先队列**（合并效率高于二叉堆）。
    -   **Dijkstra 最短路径算法**。

### **3. 斐波那契堆（Fibonacci Heap）**

-   特点：
    -   由 **多个根节点的堆序森林** 组成。
    -   具有 **摊还复杂度 O(1) 的插入操作**。
    -   **删除最小值操作** 比二项堆更快。
-   应用：
    -   **Dijkstra、Prim 等最短路径算法**（适用于稀疏图）。
    -   **高效合并多个堆**。

### **4. 配对堆（Pairing Heap）**

-   特点：
    -   类似斐波那契堆，但结构更简单。
    -   插入、合并操作比二叉堆快。
-   应用：
    -   **优先队列**（适用于需要频繁合并的小规模应用）。

### **5. 左偏堆（Skew Heap）**

-   特点：
    -   **自适应** 的二叉堆，允许不完全平衡。
    -   **合并操作比二叉堆快**。
-   应用：
    -   **适用于动态合并的优先队列**。



## **堆和树的对比**

| **数据结构**              | **是否是树** | **是否是二叉树** | **是否平衡**         | **主要用途**  |
| ------------------------- | ------------ | ---------------- | -------------------- | ------------- |
| **普通二叉树**            | ✅ 是         | ✅ 是             | ❌ 可能不平衡         | 基础结构      |
| **二叉搜索树（BST）**     | ✅ 是         | ✅ 是             | ❌ 可能不平衡         | 快速搜索      |
| **AVL 树**                | ✅ 是         | ✅ 是             | ✅ 严格平衡           | 高效搜索      |
| **红黑树**                | ✅ 是         | ✅ 是             | ✅ 近似平衡           | STL map/set   |
| **B 树**                  | ✅ 是         | ❌ 否             | ✅ 平衡               | 数据库索引    |
| **Trie**                  | ✅ 是         | ❌ 否             | ✅                    | 字符串查找    |
| **二叉堆（Binary Heap）** | ✅ 是         | ✅ 是             | ✅ 通过完全二叉树性质 | 优先队列      |
| **斐波那契堆**            | ✅ 是         | ❌ 否             | ✅                    | Dijkstra 算法 |





## **总结**

1.  **树（Tree）是一种通用的数据结构**，包括 **二叉树、BST、B 树、Trie** 等多种变体。

2.  **堆（Heap）是树的一种特殊形式**，通常是 **完全二叉树**，并且满足 **堆序性质**（最小堆/最大堆）。

3.  树的种类

    ：

    -   **二叉树**（普通二叉树、二叉搜索树 BST）。
    -   **平衡二叉树**（AVL 树、红黑树）。
    -   **多叉树**（B 树、B+ 树）。
    -   **特殊树**（Trie、线段树）。

4.  堆的种类

    ：

    -   **二叉堆（Binary Heap）**：用于优先队列。
    -   **二项堆、斐波那契堆**：用于图算法（Dijkstra）。
    -   **配对堆、左偏堆**：适用于动态合并的优先队列。

堆是树的一种特殊形式，它用于优先级管理，而树结构在数据存储和检索方面应用广泛。