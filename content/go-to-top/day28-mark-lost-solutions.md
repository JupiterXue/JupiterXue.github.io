---
title: "Go 翻山越岭——GC 标记丢失解决初步方案"
date: 2021-11-19T21:31:07+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

Go 语言采用三色标记算法来实现垃圾回收，但标记过程中存在对象漏标情况，最终导致变量被意外回收掉。为了解决这个问题我们需要一个理论基础——强弱三色性。

# 强弱三色性

强三色不变性和弱三色不变性是解决标记流程中对象漏标或错标的方案，还是看图好理解：![强弱三色性](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111192153481.png)

<center>图 1：强弱三色不变性流程</center>

左边的强三色不变性意味着所有黑色对象不能白色对象，比如 A → B 就是种非法的情况。只要让程序中不出现这种情况，黑色只能指向灰色，灰色可任意指向，那么 GC 流程一定不会有漏标或错标的情况的，这是一种比较严格的限制。



右边的弱三色不变性是说黑色对象可以指向白色对象，是合法的，但是必须保证还有指向灰色对象的可达路径，比如 A → E，前提是 A → C → E。因为有灰色可达路径，所以黑色对象指向白色对象是合法的，这种方式适当地放宽了强三色不变性的条件，保证了我们的算法可以顺利都执行完。



以上都只是理论，我们要实现这两种三色不变性，最终还是回到了 writebarrier 或 readbarrier（Go 语言中没有）。**需要注意**：在 GC 中的 writebarrier 和并发中的 writebarrier 不是一个东西。GC 中 barrier 的本质是：指针在编辑之前，可以插入一段代码 snippet of code insert before pointer modify，而 writebarrier 就是指的这段代码：

![image-20211119224915061](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111192249194.png)

<center>图 2：chan.go 中 writebarrier 的实现位置</center>

图中用红色下划线勾出来的意思是说，我们去触发 runtime.gcWriteBarrier 函数时，就是一个指针修改操作之前插入的一个函数调用。我们可以在源码中直接搜索 gcWriteBarrier 找到这个函数。因为反编译也能够看出这句函数调用的源码位置，所以我们可进入确认一下，它一定是一个指针的操作。

![image-20211119225429746](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111192254821.png)

<center>图 3：gcWriteBarrier 源码对应位置</center>

可以看到在源码的 252 行，确实是将 waiting 指针指向了空。



## 强弱三色性的实现例子

以上说了强弱三色不变性和 writebarrier，还都偏抽象、偏基础的，再来看一下强弱三色不变性分别通过哪些 writebarrier 来实现的。首先来一组实现例子：Dijkstra barrier 和 Yuasa barrier。



迪杰斯特拉 Dijkstra 是通过插入 barrier 来实现强三色不变性。（又是他！这个大牛还发明了很多让我们学计算机中痛不欲生的理论）

![image-20211119230051917](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111192300979.png)

<center>图 4：Dijkstra barrier</center>

这个函数中需要注意的是 slot 和 ptr。slot 其实是指类似图 3 中的 waiting 指针，ptr 相当于 waiting 所接收的、等号右侧的值。这段函数的意思是如果想修改 slot ，让它指向新的 ptr，需要先把这个指向的对象先标记为灰。



汤浅 Yuasa 是通过删除 barrier 来实现的弱三色不变性：

![image-20211119231116908](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111192311965.png)

<center>图 5：Yuasa barrier</center>

这段函数的意思是如果想修改 slot ，让它指向新的 ptr，需要先把待指向对象的原来对象先标记为灰。



这两个 writebarrier 实现的区别在于：一个是去标记新对象，一个是标记以前指向的对象。所以左边叫做插入 writebarrier，右边叫做删除 writebarrier。



理论来说，我们有了一个强三色不变性的算法，在所有指针操作中都加上 Dijkstra barrier，就可以避免对象丢失了。那是不是就能屡试不爽，在程序中只用这种实现方式，不再需要使用 Tyasa barrier 了？其实是不达成的。



因为在 Go 语言中对栈的操作频率是极其高的，但是 Go 官方设计者在设计的时候其实是不希望栈上的操作去加 barrier 的。如果加了很多 barrier，那么成本就非常高了。所以 Go 语言在栈上指针的操作都是不加任何 barrier 的。



也就意味着上面所说的两种 barrier 都是有问题的，两种 barrier 都遇到了水土不服的情况，这两者算法的自述如下：

- Dijkstra 算法说：我在标记和用户进程并发的时候，会出现栈上的黑色对象指向白色对象的情况（这个白色对象之前被堆上某对象所指），最终 GC 会把这些白色对象之前的堆对象回收丢弃掉。
- Yuasa 算法说：我也是在标记和用户进程并发的时候，不过我会出现堆上的黑色对象指向堆上白色的情况（这个白色对象之前被栈上某对象所指），最终 GC 会把这些白色对象之前的栈对象回收丢弃掉。



OK，留下一个未解决的悬念，下期文章聊一聊 Go 语言中的高质量 barrier。
