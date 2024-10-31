---
title: "Go 翻山越岭——内置数据结构-Timer"
date: 2021-08-30T23:13:54+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

上期 Go 文章说道了 Go 内置数据结构，从底层源码角度简单地描述了 channel。今天继续分析其他的内置数据结构 Timer。



# Timer

Timer 在 golang 1.14 版本以前比较简单，整个 Timer.go 文件中代码 才 700 行左右。整个代码就是一个数据结构和许多 goroutine。但 golang 到了 1.14 版本以后，由于官方的升级维护就变得非常复杂。所以这里笔者简单讲一下 Timer 的迭代历史，状态机制不会涉及到。



![未命名文件 (1)](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108302359026.png)



从这张图可以看到**，最早的 Timer 就是一个四叉堆**。我们平时写算法的时候二叉堆见得比较多，那是不是理解这个更困难？不是的，我们着重看 Timer 的结构和演进。像最早的 time.Sleep 或者 time.After 这两个函数，实际上我们创建了两个 Timer 类，这个 Timer 会最终加到我们的 runtime 维护的四叉堆里面。这个四叉堆其实是很好理解的，其实是有个基准的。



堆顶要放什么元素，元素如何排列？都是以触发时间为准。也就是说离当前时间最近的，一定是在堆顶的。如果来了一个新的 Timer ，它是在这个时间之后，它就会继续往堆下面走。如果比堆顶小就会涉及到对当前四叉堆的调整了，这个和写二叉堆调整算法很类似。



最老的 Timer 实现全局只有一个四叉堆，这个四叉堆专门启用一个叫做 timerproc 的 goroutine。这个 goroutine 的逻辑也不复杂，就是用一个 for 循环。这个 for 循环会不断地检查堆顶元素是不是已经到期了，如果到期了就会触发，在触发的同时会逐渐地调整堆，直到把所有需要触发的 Timer 都触发完毕为止，继续去休眠。



但这种单一的四叉堆会有一个问题。当前 Go 程序所有的 goroutine 在执行 Timer 相关操作的时候都需要去抢占操作这个堆的全局锁，而其实这个锁都是写锁。如果并发量很高，那么就会导致程序整体的吞吐量下降。全局锁对于任何程序来说，性能影响都比较大，这个问题导致了 go 程序在多核处理器上性能不佳，有人专门提了一个 issues：https://github.com/golang/go/issues/15133。（就 python 来说， 虽然 python 支持多进程，也支持多线程。但因为 GIL 全局解释器锁的存在，python 的多线程程序在同一时间只有一个线程在运行。所以多线程 Python 程序只是并发，而不是并行）



在有人提出了性能问题之后，官方就对这个问题去做了优化。官方的思路是这样的：既然一个堆和一把锁是有性能问题的，是否可以尝试将单个四叉堆扩展为多个四叉堆呢，尽量把并发性能提升。也就是说降低锁的粒度，也可以说是分片锁。来看一下图就能快速理解：



![未命名文件 (2)](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108302358793.png)



官方写死了一些东西，最多只有 64 个四叉堆。但一般不会有这么多。比如 GOMAXPROCS = 8，那么 P 结构有 8 个，每个 P 在执行调度循环的时候会绑定一个线程 M ，线程上正在执行 G。因此，每个 P 下面都会关联到一个 timer 的 G 。



但这种结构又会出现一个问题，社区中还有人提出问题：Go 程序在执行密集计算任务时会导致 timer 唤醒延迟。因此官方又做了改进，从这个版本往后，TImer 极其复杂了。



因此总结一下 Timer 1.14 的变化情况：

- 调整：
  - Timer headp 和 GMP 中的 P 绑定
  - 去除唤醒 goroutine：timerproc
- 检查：
  - 检查 timer 到期在特殊函数 checkTimer 中进行
  - 检查 timer 操作移至调度循环中进行
- 工作窃取：
  - 在 work-stealing 中，会从其他 P 那里偷 timer
- 兜底：
  - runtime.sysmon 中会为 timer 未被触发（timeSleepUntil）兜底，启动新线程



这里悄悄透露一下两个找源码的地方：



![image-20210830235814283](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108302358352.png)

偷 timer



![image-20210830235827727](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108302358770.png)

 runtime.sysmon 兜底

