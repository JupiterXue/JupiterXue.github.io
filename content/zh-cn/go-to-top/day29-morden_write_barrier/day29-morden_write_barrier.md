---
title: "Go 翻山越岭——现代化标记丢失解决方案"
date: 2021-11-22T11:49:45+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---
在 Go 语言 GC 中，为避免对象丢失，可以在所有指针操作中加 Dijkstra barrier，而 Go 官方设计者为了降低 Go 运行环境中的成本，不希望栈上的操作频率很高，所以限制了 Go 不能在栈上操作指针时加 barrier。因此 Dijkstra barrier 和 Yuasa barrier 都失效了。



因为 Go 语言不在栈上去加一个 writebarrier，所以任何一个 barrier 单独拿出来都没办法解决 Go 语言 GC 标记的正确性问题。为了解决这个问题，机智的语言设计者就把 Dijkstra barrier 和 Yuasa barrier 结合起来，做出了集大成者——Hybrid barrier，并且写了一个提案 preposal。网上有很多博客会去解析 writebarrier 的源码，其实就是这张图中的，中和两种算法的集合体。

![image-20211122224523656](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111222245159.png)

<center>图 1：混合屏障 Hybrid barrier</center>

这里的意思是说，把要修改的指针所指向的老的对象先标灰，然后判断当前栈是否是灰色。如果是灰色就去把新来的对象也标灰。最后再做个赋值操作。



这个算法虽然集成了两种算法思想，但两个算法所遇到的问题正好被有效避免了。虽然作者写个提案 preposal，但后来因为没有时间就没有去实现。后来他们还发现，如果要检查栈的颜色的工作话还需要做一些同步操作，而这个 stack check 动作维护成本太高，所以暂时放弃了，至今也没实现。



但我们亲自去看 Go 语言源码可以发现，它不是内部实现，代码片段如下：

```go
// slot is the destination in Go code.
// ptr is the value that goes into the slot in Go code.
func RealityWB(slot *unsafe.Pointer, ptr unsafe.Pointer) {
    shade(*slot)
    shade(ptr)
    *slot = ptr
}
```

可以看到，Go 源码中其实是没有任何 if 判断。如果发送了任何堆上指针的变化，那么一定要把原来的对象和新指向的对象都标灰。然后再执行赋值操作。



在分析具体流程之前，先梳理一下 Dijkstra barrier 和 Yuasa barrier 的逻辑流程。![Dijkstra-inseetion barrier](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111222303362.png)

<center>图 2：Dijkstra insertion barrier 流程</center>

Dijkstra barrier 是说，在堆上对象中，想要把 A 对象指向 C，直接操作就会把黑色对象指向白色对象，但在代码中其实是有两步，先标灰 shade(ptr)，再赋值 *slot = ptr 让指针指向 C。



在标记过程中，还会有栈上对象指向堆上的白色对象情况。如果我们在只有插入 barrier 的情况下，不管栈。那么有可能因为 GC 没有检测到 D 对象，最终导致 D 对象被回收。这便是这个算法的缺陷。这种情况下我们需要在所有堆对象完成标记后，对栈做一次扫描，将这些对象全部变成黑色，也就不会导致误回收的情况。但是这会导致第二个 stw 的时间非常长、成本较高。



而 Yuasa barrier 是说，如果栈上有个指针 X 指向堆上的白色对象 C，同时还有一个灰色对象 A 能够可达这个白色对象，就形成了弱三色不变性。最后在标记流程中，对象 C 便不会被丢失。 ![Yuasa deletion barrier (1)](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111222321718.png)

<center>图 3：Yuasa deletion barrier 流程</center>

A → B → C 这个路径有可能被修改掉，就可能出现问题。比如灰色的 A 对象指向了白色的 D 对象，并把 B 对象标灰。虽然保证了所有堆上的指针连接在断开前被置为灰，使得扫描结束后不再需要 STW 了，但无法防止栈上指向的白色对象断开连接后，被堆上的黑色对象引用。



以上，了解这两种插入（Dijkstra insertion barrier）和删除（Yuasa barrier）的 barrier，再了解到它们的缺陷，最终知道为什么我的 Go 语言中会采用将两种算法合并起来，做出了混合的 barrier。



OK，下期终于到了 Go 语言内存管理和垃圾回收的最后一个章节——Go 协助标记和总结 ，累积写了 11 篇，终于粗略地过了一遍。
