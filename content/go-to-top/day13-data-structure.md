<<<<<<< HEAD
---
title: "Day13 Data Structure"
date: 2021-08-29T14:47:52+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

**内置数据结构是一门编程语言的基础核心**，了解基本语法就能够让我们进行简单的开发。今天开始，就来系统地研究 Go 语言内置数据结构。



首先，先来看看所有的Go 语言所有内置数据结构都有哪些，如下图所示，列出了思维导图：

![未命名文件](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108291453518.png)



既然是系统地研究，这里就需要用到暴力破解的思想——把每个数据结构都遍历吃透。



# Channel

之前的文章中提到了通过反汇编调试工具来查看 Go 语言的源码。还提到了三种情况会导致 panic 的关键函数 chansend、chanrecv，下面来进行源码逻辑分析。忽略一些细节实现，来看看 chansend 的流程图：

![chansend 流程图](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108291706892.png)

图上省略了：

1. select dafault 的情况。
2. 逻辑执行时碰到 ch 已 close 的情况。



从流程图，能够清晰地看到 在 chansend 的内置函数中，Go 语言是如何处理我们发送的数据。紧接着我们再来看看 chanrecv 的流程图：

![chanrecv 流程图](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108291706194.png)

图上也省略了：

1. select default 的情况。
2. 逻辑执行时碰到 ch 已经 close 的情况。



对比发现，两个流程差不多，因此 channel 的发送和接收的逻辑都是差不多的，都要判断是否为空，是否阻塞，然后看缓存情况，一个明显不一样的特征是 chansend 要判断满，chanrecv 要判断空。



我们常说Go 语言中 channel 是并发安全的，什么意思呢？从上面的流程可以发现：chansend、chanrecv、closechan 都是要加锁的。即便如此，从代码层面我们还是看不到这些锁，那我们能否通过代码来看，“并发安全”具体而言是什么意思呢？下面罗列了三者的源码：

```go
// chansend 源码
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
    if c ==nil {
        if !block {
            return false
        }
        gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
        throw("unreachable")
    }
    
    lock.(&c.lock)
}

// chanrecv 源码
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
    if c ==nil {
        if !block {
            return false
        }
        gopark(nil, nil, waitReasonChanReceiveNilChan, traceEvGoStop, 2)
        throw("unreachable")
    }
    
    lock.(&c.lock)
}

// closechan 源码
func closechan(c *hchan) {
    if c == nil {
        panic(plainError("close of nil channel"))
    }
    
    lock(&c.lock)
}
```

虽然我们在用户态完全看不到，但我们可以从底层源码看到这个明显的加锁操作，并且这个加锁的操作基本上都是在函数开始加锁，函数结束解锁，除了 goready() 要放在锁外面，其他也没有什么特殊的，这便是达成并发安全的条件之一。



第二个并发安全的条件就是挂起和唤醒。刚刚提到了两个关键的函数 gopark 和 goready，这里我们可以记住两三个结论：

- Sender 挂起，一定是由 receiver（或 close）唤醒
- Receiver 挂起，一定是由 sender（或 close）唤醒
- 结对操作。只要有用到 gopark，那么代码中一定有另外一个地方在执行 goready 

在之前的文章中提到，可接管的阻塞，均是由 gopark 挂起，每一个 gopark 都会对应一个唤醒方。



以上第一二点，在代码中如何体现呢？比如 Sender 挂起，也就是 chansend 中的流程，当 buffer 满了需要执行一个 gopark，把 goroutine 和线程解绑，让 goroutine 进入 sendq 的等待队列。之后在 chanrecv 或 chanclose 的时候，检查等待队列，把 goroutine 从等待队列弹出，再用 goready 把它唤醒，并且 close 操作会把所有 goroutine 唤醒。Receiver 流程和 Sender 是相对的，代码都非常对称，就不逐一赘述。



基于以上第三点，这里总结一下与 gopark 对应的 goready 位置：

- channel send → channel recv/close
- Lock → Unlock
- Read → Read Ready，epoll_wait 返回了该 fd 事件时
- Timer → checkTimers，检查到期唤醒



具体找的方法，可以在 Goland 中双击 shfit 搜索 runtime 中 chan.go 源码，如图所示：

![image-20210829165005782](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108291658456.png)



当我们知道了“与 gopark 对应的 goready 位置”，在自己看 channel 源码的时候就能够抓大放小、有的放矢。

=======
---
title: "Day13 Data Structure"
date: 2021-08-29T14:47:52+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

**内置数据结构是一门编程语言的基础核心**，了解基本语法就能够让我们进行简单的开发。今天开始，就来系统地研究 Go 语言内置数据结构。



首先，先来看看所有的Go 语言所有内置数据结构都有哪些，如下图所示，列出了思维导图：

![未命名文件](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108291453518.png)



既然是系统地研究，这里就需要用到暴力破解的思想——把每个数据结构都遍历吃透。



# Channel

之前的文章中提到了通过反汇编调试工具来查看 Go 语言的源码。还提到了三种情况会导致 panic 的关键函数 chansend、chanrecv，下面来进行源码逻辑分析。忽略一些细节实现，来看看 chansend 的流程图：

![chansend 流程图](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108291706892.png)

图上省略了：

1. select dafault 的情况。
2. 逻辑执行时碰到 ch 已 close 的情况。



从流程图，能够清晰地看到 在 chansend 的内置函数中，Go 语言是如何处理我们发送的数据。紧接着我们再来看看 chanrecv 的流程图：

![chanrecv 流程图](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108291706194.png)

图上也省略了：

1. select default 的情况。
2. 逻辑执行时碰到 ch 已经 close 的情况。



对比发现，两个流程差不多，因此 channel 的发送和接收的逻辑都是差不多的，都要判断是否为空，是否阻塞，然后看缓存情况，一个明显不一样的特征是 chansend 要判断满，chanrecv 要判断空。



我们常说Go 语言中 channel 是并发安全的，什么意思呢？从上面的流程可以发现：chansend、chanrecv、closechan 都是要加锁的。即便如此，从代码层面我们还是看不到这些锁，那我们能否通过代码来看，“并发安全”具体而言是什么意思呢？下面罗列了三者的源码：

```go
// chansend 源码
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
    if c ==nil {
        if !block {
            return false
        }
        gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
        throw("unreachable")
    }
    
    lock.(&c.lock)
}

// chanrecv 源码
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
    if c ==nil {
        if !block {
            return false
        }
        gopark(nil, nil, waitReasonChanReceiveNilChan, traceEvGoStop, 2)
        throw("unreachable")
    }
    
    lock.(&c.lock)
}

// closechan 源码
func closechan(c *hchan) {
    if c == nil {
        panic(plainError("close of nil channel"))
    }
    
    lock(&c.lock)
}
```

虽然我们在用户态完全看不到，但我们可以从底层源码看到这个明显的加锁操作，并且这个加锁的操作基本上都是在函数开始加锁，函数结束解锁，除了 goready() 要放在锁外面，其他也没有什么特殊的，这便是达成并发安全的条件之一。



第二个并发安全的条件就是挂起和唤醒。刚刚提到了两个关键的函数 gopark 和 goready，这里我们可以记住两三个结论：

- Sender 挂起，一定是由 receiver（或 close）唤醒
- Receiver 挂起，一定是由 sender（或 close）唤醒
- 结对操作。只要有用到 gopark，那么代码中一定有另外一个地方在执行 goready 

在之前的文章中提到，可接管的阻塞，均是由 gopark 挂起，每一个 gopark 都会对应一个唤醒方。



以上第一二点，在代码中如何体现呢？比如 Sender 挂起，也就是 chansend 中的流程，当 buffer 满了需要执行一个 gopark，把 goroutine 和线程解绑，让 goroutine 进入 sendq 的等待队列。之后在 chanrecv 或 chanclose 的时候，检查等待队列，把 goroutine 从等待队列弹出，再用 goready 把它唤醒，并且 close 操作会把所有 goroutine 唤醒。Receiver 流程和 Sender 是相对的，代码都非常对称，就不逐一赘述。



基于以上第三点，这里总结一下与 gopark 对应的 goready 位置：

- channel send → channel recv/close
- Lock → Unlock
- Read → Read Ready，epoll_wait 返回了该 fd 事件时
- Timer → checkTimers，检查到期唤醒



具体找的方法，可以在 Goland 中双击 shfit 搜索 runtime 中 chan.go 源码，如图所示：

![image-20210829165005782](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108291658456.png)



当我们知道了“与 gopark 对应的 goready 位置”，在自己看 channel 源码的时候就能够抓大放小、有的放矢。

>>>>>>> 0cae884533b7a5e0cf3eb8be746627712e680e3e
