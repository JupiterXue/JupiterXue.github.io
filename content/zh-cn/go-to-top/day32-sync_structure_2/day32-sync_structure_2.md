---
title: "Go 翻山越岭——并发数据结构（2）"
date: 2021-11-29T22:39:57+08:00
category: tech
tags:
    - Go
keywords:
    - 理论

---

本期文章继续拆解 Go 语言中并发部分的一些内置数据结构，了解它们的运行逻辑，梳理其中的要点。

## sync.Pool

上期说到 sync.Pool 结构，它的应用场合非常多，不过主要是用于 GC 标记阶段消耗大量 CPU 或者进程 RSS 占用过高情况。并且通过梳理源码发现它的逻辑是和缓存的机制很像是，都是多级缓存结构的。上期也留下了一个悬念，当 sync.Once 发生 GC 时，sync.Pool 的代码逻辑有有怎样的变化呢？还是看图直观点：

![sync.Pool GC 时逻辑](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111292258979.svg)

<center>图 1：sync.Pool GC 时的逻辑</center>

可以看到，当发生 GC 的时候，相当于将 local 和 localSize 直接向下平移，local 赋值给了 victim，localSize 赋值给了 P。如果之前 victim 和 victimSize 有值的话就直接丢弃掉。通过这种方式将以前不需要的、多出来的对象淘汰掉。如果把其中的值刚替换掉，并且又需要从 sync.Pool 中 get 获取对象的时候，local 为空，因此就会去 victim 中去寻找值。其他的逻辑，其实和之前的逻辑是完全一致的。最终，如果 victim 找不到值，就会去 P 中去找值。



一个小插曲。sync.Pool 最早的实现中，操作 shared 是有锁的。而从 Go 1.13 开始，这个锁就直接被干掉了，变成了右下方看起来有点诡异的双端链表，它就可以无锁进行操作。



并且早期的 sync.Pool 在进行 GC 的时候会将其中的对象完全清空掉的。如果我们的程序对 sync.Pool 有重度使用的情况，那么在每次 GC 完成之后，如果正好发生了一个流量请求的数量波动，就会造成应用程序大量地阻塞在一个锁上，会有短时间的延迟波动。

## semaphore

**信号量 （semaphore）**，是和所有锁相关的实现基础和所有同步原语的基础设施。可以看看这幅图：

![semaphore (2)](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111292325509.svg)

<center>图 2：semaphore 逻辑</center>

semaphore 它的底层其实可以理解为一个个 waitlink 的链表（先看最下方）。不过这个链表需要在 runtime 中去维护，所以设计了一个 semtable（最上方）。semtable 就是一个定长的数组，这里的长度为 251，也就是锁信号量的地址移进来的。其中每一个元素都是 treap 的数据结构，从字面上来看，“treap”其实就是“tree”和“heap”的合成词，也就是说它既是一棵树也是一个堆：

- 它是按照 lock addr 排列的二叉搜索树。如果想要加锁，锁中有个信号成员，成员有个地址，通过这个地址最终找到它在二叉搜索树的哪个节点上。并且每个节点就是一个链表。
- 既然是二叉搜索树的结构，那么为了保持平衡，所以会为每个节点赋值一个 ticket。
- 这个 ticket 会在每次初始化的时候通过 fastrand 来生成。

如果我们有任何一个锁发生了阻塞，最终都是阻塞在以上这个结构中。最终找到它的节点以后，它们是挂载在等待链表之后的。



小结，整个逻辑就是通过 semtable 进行查找，然后在 treap 树中查找，最终把找到的节点挂载在对应元素后面。

## sync.Mutex

**普通写锁 （sync.Mutex）**有两个字段，状态 state 和信号量 sema，其实也就是 atomic 锁和 semaphore 信号量结合实现的。其逻辑如下：

![sync.Mutex 逻辑](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111292341422.svg)

<center>图 2：sync.Mutex 逻辑</center>

state 结构中，会将最后三位作为标记位：

- 最后一位表示当前锁是否已上锁、是否已经被抢占。
- 倒数第二位表示当前 mutex 是否有被唤醒并且在等待锁的 G
- 倒数第三位表示当前的互斥量是否要进入饥饿模式，如果置为 1，那么之后来的 Goroutine 不能去抢占这个锁，必须在 semaphore 结构中的 treap 找到位置，并且排队在之前的 sudog 后面。
  - 注意这里的饥饿模式，其实是为了实现锁的公平性。官方在设计锁的时候，给了 Goroutine 一个优先级，也就是最新来抢锁的 Goroutine 的优先级是最高的。源码中是一个自旋的逻辑，一直在检查能否抢到一把锁，如果能够抢到就直接返回 lock 结构，让用户执行 lock 之后的代码。而如果当前 mutex 已经饥饿了，那说明在队列中已经有 Goroutine 等待了很长时间了，那么新来的 Goroutine 能够拿到锁就不太合理了。

整个 state 是 32 位的结构，除去最后三位，剩下 29 位就是所有等待 Goroutine 的计数。

我们在看代码的时候可能也会遇到这个将最后三位最为状态的优化结构，所以它的计数每次都会向右移动三位。



除了 state，还有一个 sema，如果了解得多，可以发现：几乎所有和锁相关的数据结构，最终都会有这个成员。这个成员最终对应的就是锁的链表，这个链表也就是 treap 树上的一个节点。



OK，下期文章继续再对几个常见并发结构做梳理。
