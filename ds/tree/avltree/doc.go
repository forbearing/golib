package avltree

/*

1.右旋 (Right Rotation, 单旋转)
- 适用情况：LL (Left-Left) 失衡
- 原因：某个结点的 左子树 过高 (BF > 1)，而且新增的结点也在左子树的左侧。
- 解决方案：对 失衡结点进行右旋。
- 旋转前后依然保持：A < B < C
旋转前
        C
       /
      B
     /
    A
旋转后
        B
       / \
      A   C


2.左旋 (Left Rotation, 单旋转)
- 适用情况：RR (Right-Right) 失衡
- 原因：某个结点的 右子树 过高 (BF < -1)，而且新增的结点也在右子树的右侧。
- 解决方案：对 失衡结点进行左旋。
- 旋转前后依然保持：A < B < C
旋转前
    A
     \
      B
       \
        C
旋转后
	B
   / \
  A   C


3.左-右旋 (Left-Right Rotation, 双旋转)
- 适用情况：LR (Left-Right) 失衡
- 原因：某个结点的 左子树 过高 (BF > 1)，但新增的结点在左子树的右侧。
- 解决方案：先对左子树进行 左旋，再对失衡结点进行 右旋。
- 旋转前后依然保持：A < B < C
旋转前
        C
       /
      A
       \
        B
先左旋 (A 左旋)
        C
       /
      B
     /
    A
再右旋 (C 右旋)
        B
       / \
      A   C


4.右-左旋 (Right-Left Rotation, 双旋转)
- 适用情况：RL (Right-Left) 失衡
- 原因：某个结点的 右子树 过高 (BF < -1)，但新增的结点在右子树的左侧。
- 解决方案：先对右子树进行 右旋，再对失衡结点进行 左旋。
- 旋转前后依然保持：A < B < C
旋转前
    A
     \
      C
     /
    B
先右旋 (C 右旋)
    A
     \
      B
       \
        C
再左旋 (A 左旋)
      B
     / \
    A   C



*/
