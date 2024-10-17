---
title: "Go 翻山越岭——标记流程"
date: 2021-11-15T09:07:46+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---
**Go 语言使用并发的三色标记（tricolor marking）算法作为其垃圾回收算法**，并且在 Go 中的 GC 没有分代，没有压缩，没有对象移动。整个 GC 流程如图：

![gc](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111152027238.png)

<center>图 1：GC 全流程</center>

整个图看起来觉得很复杂，代码很长，但其实里面的关键节点不是很多。左上角可以看到三种 GC 的触发的方式：后台触发，用户主动触发，内存分配触发。



在触发之后会进入 runtime.gcStart 函数。因为很多内存分配器都会进来，但最终只有一个能够进去，所以它会抢一个全局锁。然后启动所有后端标记的 worker 并且把它和每个 P 绑定。每个 P 都会对应一个 BgMarkWorker，但不一定每个 worker 都能全部启动，其实是按比例启动的。所以在 cBgMarkStartWorkers 函数中初始化所有 worker 都会让它们进入一个休眠状态 goPark，它会在执行到 schedule 函数中的 finRunnableGCWorker 时之前休眠的 GC 标记线程唤醒，之后是从 G runing 状态变为 G waiting 状态。



注意，在标记的时候会控制所有标记协程的 CPU，所以不能把所有的 P 上协程全部唤醒，如果全唤醒意味着 CPU 占用中标记就占用了了80% 到 90%，应用程序也就跑不起来了。



**gcDrain 是整个标记的主流程**。其中，mark 标记根节点，gc work 全局和本地的平衡，scan 对象并进行标灰。整个清扫过程完成之后会，会有一个分布式完成算法进行判断，然后会进入 gcMarkDone 函数。



gcMarkDone 阶段会去把剩余的 write buffer 或 gc worker 里面的工作都清理掉。整个函数中还可能会有回退情况循环情况。标记完成之后就是 gcMarkTermniation ，也就是 GC Mark 的终极阶段。



gcMarkTermniation 会调用 gcSweep ，去唤醒 sweep 的协程。这个 sweep 协程其实是在整个程序初始化的时候就已经启动了，但在后台是一直休眠的。GC 完成之后唤醒就能进行清扫工作了。



大致的代码流程就是这样，再来进行一些细节 QA。



# GC 标记流程

首先我们要问自己几个个问题：

**1、标记对象从哪里来？**

一个是后台的标记协程 gcMarkWorker 里 ，另一个是应用程序在 GC 标记过程中会并发地去写或删除堆上的指针，即 mutator write/delete heap pointers。



**2、标记对象到哪里去？**

这两个对象写入的目标其实是不一样的。前者是写入 work buffer，后者是写入 write barrier buffer。具体代码如下：

- Work buffer
  - 本地 worker buffer => p.gcw
  - 全局 work buffer => runtime.work.full
- Write barrier buffer
  - 本地 write barrier buffer => p.wbBuf



再来看一个图：

![未命名文件](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111152157603.png)

<center>图 2：scanobject/greyobject 流程</center>

可以看到两个 gc 相关的结构 gcw 和 wbBuf。gcw 也就是 gc work 的意思，这个 work 包含两个数组，一个叫 工作队列1 wbuf1 一个叫工作队列2 wbuf2。前面说到的标记也就是往 wbuf1 中存入值，这个值也就是指针的值；如果 wbuf1 写满了以后会把 wbuf1 和 wbuf2 先交换一下（没想到吧，我也没想到），再继续存 wbuf2；如果 wbuf2 也写满了，就会把这些数据推送到全局唯一的队列 Globally unique work；在标记过程中也会遇到 buffer 中没有值的情况，这时就会从全局把数据拿回来。如果没有工作执行的情况，还会从 wbBuf (write barrier buffer) 中去操作。



**3、GC 消费对象是哪些？**

说到了 GC 的队列和 GC 的来源，也需要了解 GC 的消费者是哪些。其实也是 work buffer，它会不断消费队列，按着广度优先遍历算法，直到所有堆上的指针都标记完了，认为工作就完成了。



OK，说了很多关于标记的偏逻辑上的，具体怎么标记下回继续分析，看一下 Go 的三色标记具体是指什么。
