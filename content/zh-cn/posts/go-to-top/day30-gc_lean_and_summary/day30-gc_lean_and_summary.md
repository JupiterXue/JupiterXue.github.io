---
title: "Go 翻山越岭——GC 补充与总结"
date: 2021-11-24T21:33:37+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

最后一篇关于 Go 语言内存管理与垃圾回收的文章，再对之前流程中做一些补充，然后对这期系列文章做个总结。

## GC 补充

在 GC 标记流程中有一个环节能够辅助标记：

- Goroutine 中有 gcAssistBytes 字段。
- 当后台 gcWorker 标记时，会累积 credit，记录在 gcController.gbScanCredit 中

![image-20211124233351199](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111242333365.png)

<center>图 1：协助标记源码-判定</center>

- Goroutine 想执行内存分配，要先尝试去 gcController.bgScanCredit 中去借债，如果借到了足够的债，那么就不用协助标记。
- 如果借不到，那就先协助标记，标记完成后再去分配内存。

![image-20211124233409903](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111242334952.png)

<center>图 2：协助标记源码-分配内存</center>



在标记流程阶段，堆上对象可能出现引用交叉情况：

![标记流程-堆对象引用交叉情况](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111242342645.png)

<center>图 3：堆上对象引用交叉情况</center>

- 一个是 isMarked 剪枝

![image-20211124233904984](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111242339052.png)

<center>图 4：isMarked 源码</center>

- 另一个是 atomic.Or8

![image-20211124233915764](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111242339838.png)

<center>图 5：atomic.Or8 源码</center>



再补充一些零零碎碎的小知识点：

- GC 的 CPU 控制目标是整体的 25%。
- 当 P = 4 * N 时，只要启动 N 个 wroker 就可以使用。
- 当 P 无法被 4 整除时，需要吃苦耐劳的 gcMarkWorker 来帮助做一部分工作：
  - 作为全局 GC 员工 Dedicated worker，需要一直干活，知道被抢占。
  - 作为兼职 GC 员工 Fractional worker，达到业绩目标（fractionalUtilizationGoal）时，可以主动让出。
  - 另外一种 IDLE 模式。在调度循环中发现找不到可执行的 g ，并且有标记任务没有完成的情况下，就可以开启 IDLE 模式去帮忙。
- Worker 运行模式位于：\_p_.gcMarkWorkerMode。
- 栈本身的内存：newstack、shrinkstack。
- 使用 allocManual 和 freeManual 相当于手动管理内存，不计入 heap_inuse 和 heap_sys，而是计入 stackinuse 和 stacksys。
- 栈上变量的内存变化：SP 移动销毁，简单快速。



# 总结

最后来做个总结，还是简单回顾一下垃圾回收代码的流程。

1. **gcStart → gcBgMarkWorker && gcBgRootPrepare**。GC 触发后，进入 gcStart 函数，这个函数会启动所有 P 相应的后台 MarkWorker 并且进入休眠状态，同时也会准备好整个 GC 标记流程的根节点对象。
2. **schedule → findRunnableGCWorker → gcBgMarkWorker**。在调度流程 schedule 中，findRunnableGCWorker 会去唤醒适宜数量的 gcBgMarkWorker
3. **gcBgMarkWorker → gcDrain  → scanobject → greyobject**。这些 gcBgMarkWorker 被唤醒之后执行自己的工作，进入 gcDrain  排空 GC 任务，然后执行广度优先遍历算法 scanobject，扫描这些对象并且标灰 greyobject（也就是将 gcMarkbit 置为 1 并放入 gcw 队列中，set mark bit and put to gcw）
4. **gcMarkTermination → gcSweep → sweepg && scvg → sweep → wake bgsweep && bgscaveng**。gcBgMarkWorker 在调用 gcMarkDone 去排空各种 wbBuf 后，就会使用分布式检查算法 termination 结束流程。进入gcMarkTermination，调用 gcSweep 唤醒最终清扫的 Goroutine，即后台沉睡的 sweepg 和还给操作系统的 Goroutine，即环内存 scvg，最后再去执行清扫 sweep，然后 wake 唤醒 bgsweep 和 bgscavenge。



可以看到 Go 语言的标记流程部分包含广度优先遍历算法，意味着 GC 消耗的 CPU 主要是和对象的数量相关的。所以我们在做优化的时候，抓住这个重点，想方设法地做对象复用和降低全局堆上分配的对象数。



梳理清楚 GC 的整个流程是挺有意思的，如果去看一些关于垃圾回收的书籍，书中都会把这样类似的流程也梳理出来。但其实我们并不需要把整个 GC 的流程都背下来的，因为对工作的帮助没那么大。实际上这个 GC 算法能够维护的也就那么官方的几个写了十五、二十年代码、牛得不行的人。（在 GC 中还有涉及有锁和无锁并发的算法，这种比理解调度的困难程度还高一些）



假如在平时工作中，我们遇到了 GC 的 bug，多半只能等官方来维护了。而我们普通人如果能够把流程和理论研究明白，把别人写的代码也看懂，已经很不错了。
