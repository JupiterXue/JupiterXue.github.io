---
title: "Go 翻山越岭——调度-day2"
date: 2021-08-05T14:45:19+08:00
category: tech
tags:
    - Go
keywords:
    - Scheduler
---

上回说到 [Go 的调度流程本质上是⼀个⽣产-消费流程](https://mp.weixin.qq.com/s?__biz=Mzg3MjYzMzQzMQ==&mid=2247483687&idx=1&sn=a357de9b46d3dbf466b9f52f6e2f413a&scene=21#wechat_redirect)，今天来讲一讲“调度组件与调度循环”，再来回顾一下两个生动的动画 goroutine 的⽣产端, goroutine 的消费端。



当 goroutine 处于生产端时，M 执行调度循环时，必须与一个 P 绑定。并且我们常说的 Work stealing 就是说的 runqsteal -> runqgrab 这个流程。

当 goroutine 处于消费端时，执行的是一个循环：runtime.schedule → execute → runtime.gogo → runtime.goexit → runtime.schedule（回到原点），并且最终 P.schedtick = 0。

初学 scheduler 对于以上的流程感受是比较浅的，再来看看这些符号所代表的含义，就能更好地理解了：

> G: goroutine，计算任务。由需要执行的代码和其上下文组成。（上下文包括：当前代码位置，栈顶、栈底地址，状态等）

> M: machine，系统线程，执行实体，想要在 CPU 上来执行代码，必须有线程，与 C 语言中的线程相同，通过系统调用 clone 来创建。

> P: processor，虚拟处理器，M 必须获得 P 才能执行代码，否则必须陷入休眠（后台监控线程除外），你也可以将其理解为一种 token，有了这个 token，才有在物理 CPU 核心上执行的权利。

上面所说的循环调度流程，都是在正常情况下运作的。而实际业务中我们往往还会遇到其他情况——阻塞。如果程序中有阻塞，那么线程不就全部被堵上，程序就卡住了么？

让我们来看看以下几种情况，在线程发生阻塞的时候，是否会无限地创新线程？（并不会）

案例1：

```
// channel send:
var ch = make(chan int)
ch <- 1

// channel recv:
var ch = make(chan int)
<- ch
```

案例2：

```
// net read
var c net.Conn
var buf = make([]byte, 1024)

// data not ready, block here
n, err := c.Read(buf)

// net write
var c net.Conn
var buf = []byte("hello")

// send buffer full, write blocked
n, err := c.Write(buf)
```

案例3：

```
var (
    ch1 = make(chan int)
    ch2 = make(chan int)
)

// no case ready, block
select {
    case <- ch1:
        println("ch1 ready")
    case <- ch2:
        println("ch2 ready")
}
```

案例4：

```
// common func
time.Sleep(time.Hour)


var l sync.RWMutex
// somebody already grab the lock
// block here
l.Lock()
```

以上这些情况不会阻塞调度循环，而是会把 goroutine 挂起。挂起，其实是让 g 先进某个数据结构，待 ready 后再继续进行，并不会占用线程。这时候，线程会进入 scedule，继续消费队列，执行其他的 g。那么我们如何应对以上的情况呢？正确使用锁。

下面还有三种应用阻塞在锁的情况：

1. 按 lock addr 排列的二叉搜索树。
2. 按 ticker 排列的小顶堆。
3. ticket 是每个 sudog 初始化时用 fastrand 生成的

用了锁又会遇到新的问题：为啥有的等待是 sudog，有的是 g？让我们看官方怎么说：

> sudog represents a g in a wait list, such as for sending/receiving on a chnnel.

> sudog is necessary because the g ↔ synchronization boject relation is many-to-many. A g can be on many wait lists, so there may be many sudogs for one g; and many gs may be waiting on the same synchronization object, so there may be many sudogs for one object.

啥意思？也就是说，一个 g 可能对应多个 sudog，比如一个 g 会同时 select 多个 channel。

呼，好像终于告一段落了。等等，前面这些都是能够被 runtime 拦截到的阻塞。来看看英文解释：

> sysnb: syscall nonblocking

> sys: syscall blocking

一些 runtime 无法拦截的例子：

```
package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
void output(char *str) {
    usleep(1000000);
    printf("%s\n", str);
}
*/

import "c"
import "unsafe"
```

我们在执行 c 代码，或者阻塞在 syscall 上时（这个没有列出来），必须占用一个线程

遇到这种问题，我们聪明的程序员还是有解决办法：sysmon——system monitor，它有着高优先级，能够在专有线程中执行，不需要绑定 P 也能执行

> Check for deadlock situtation.

> The check is based on number of runnint M`s, if 0 → deadlock.

有四个主要注意的地方：

1. checkdead，常见误解为可以检查死锁。
2. netpoll，inject g list to global run queue
3. retake，如果是 syscall 卡了很久，那就把 p 剥离(handoff p)
4. retake，如果是用户 g 运行很久了，那么发信号 SIGURG 去抢占过长时间的 G。

关于调度，就暂时告一段落。理解比较浅显，更多的是在重复课上的内容。后续在实践中积累了更多经验后再来做深入分析。
