---
title: "Go 翻山越岭——常见并发 bug"
date: 2021-12-06T22:30:43+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

这期文章来聊一聊 Go 语言中常见的并发 bug 有哪些，并不是说要写 bug，而是在出现这些情况的时候，有相应的解决方案。

![img](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112062349044.png)

## 死锁

死锁应该是最常见的，来看看代码，应该是最直观的：

```go
func x() {
    a.RLock()
    defer a.RUnlock()
    
    y()
}

func y() {
    a.RLock()
    defer a.RUnlock()
    
    // do some logic
}
```

这里举的例子都是对对象 a 加锁，可以认为它是个全局变量。a.RLock() 加了读锁，defer 延迟解锁 a.RUnlock()，这时候调用了一个函数 y。进入函数 y 以后，再次执行 a.RLock()和 defer 解锁 a.RUlock()。这个过程相当于递归地执行了两次 a.RLock()，就会发生死锁。



还有段较为简单的代码：

```go
func x() {
    a.Lock()
    b.Lock()
}

func y() {
    b.Lock()
    a.Lock()
}
```

有两个函数 x 和 y，都有 a.Lock() 和 b.Lock()，只不过顺序不一样。在主程序中，并发执行这两个函数时，就会出现循环等待，没有办法对锁进行抢占，也就出现了死锁。



面对死锁问题，我们常见的解决办法是通过 pprof 进入 goroutine 页面查看。Go 语言自带了 pprof 工具，只要进程不挂掉就可以进入 goroutine 页面插件。但是有个问题是，如果死锁是偶发的情况，那就需要想办法在整个进程崩溃掉之前完整地把 goroutine 的栈保存下来（又是另一个话题了）。

### 并发读写

Map concurrent writes/reads，Map 的并发读写也是容易出现死锁的情况。



既然 Map 并发读写容易出问题，那么需要我们注意的是，就是尽量不用 Map 进行并发读写（今日份废话文学）。某些场景下，用 sync.Map 不容易写出死锁代码，那么就值得用。但是 sync.Map 有点尴尬，在 Go 1.17 以前，因为 Go 语言没有泛型，所以它的 sync.Map 的签名并不是特别友好，都是 interface 并且需要去做断言。



小结，如果要使用并发读写，在某些场景下对 sync.Map 简单做个封装就好了。如果用 Map 一定不要用并发读写。



说着容易但其实并不容易，尤其在老版本的 Go 语言中。在日志中，把 http.Request 中的 context 打印一下时，整个程序就有可能出现偶发性的 Map 并发报错。这个问题是因为当初 Go 官方在写 http 库的时候，会在 context 里写个 Map，然后用这个 Map 去维护连接之类的东西。但是，用户其实是不知道这件事情的，用户会把 context 打印在日志中，在打印时会调用 context.String 方法，这个 String 方法会用反射一步步往下遍历，相当于是个递归的过程，最终回到之前那个 Map，因而非常可能发生 Map 的并发读写问题。毕竟我们在打印日志的时候一般不会考虑说对 context 做保护的情况。后来应该是做过修改的，所以现在我们可以在日志中输出 context 是没什么问题的。



即便整个程序做了些很完备的处理，还是需要在线上崩溃时输出 stderr，将这个 stderr 重定向到一个单独文件中。事故发生之后就能够去查这个文件，否则发生了什么事情都不知道。大多数公司里面是通过 supervisor 对进程做了托管，它能够把程序中的 stdout 和 stderr 都重定向到相应的文件当中。如果为了稳定性考虑，也需要把文件中的一些 fatal 或 panic 都监控起来。再严谨些会对针对这些 panic 向运维或用户发送些 p0 或 p1 的报警。

## Channel 关闭

Channel 在关闭的时候很容易导致 panic，写出复杂的代码时，自己可能也不知道自己为什么要去那里 close channel 或重复地关闭一个 channel。因此，在我们使用 channel 进行编程的时候，就需要有意地去考虑 **Channel closing principle**：

1. M receivers, one sender, the sender says "no more sends" by closing the data channel.
2. One receiver, N senders, the only receiver says "please stop sending more" by closing an additional signal channel.
3. M receivers, N senders, any one of them says "lets` end the game" by notifying a moderator to close an additional signal channel.

这是出自 go101 所描述的原则：https://go101.org/article/channel-closing.html。来加以说明一下：

- 第一条原则说的是，我们有 M 个消费者，1 个生产者的时候，需要生产者主动地通知消费者“数据已经发送完了”，大家可以散伙了，也就是可以关闭数据 channel 了。
- 第二天原则是说：一个消费者，N 个生产者。这种情况就需要在 channel 以外去做一个 down channel。通过消费者去关闭 down channel 然后告诉生产者“channel 的对端已经没有人在消费了，你们需要停止当前的生产动作”，同时也需要对 sender 的代码中去写一些 select 操作。
- 第三个原则：M 个消费者、N 个生产者。这种情况就需要任意一个出问题的时候去把额外的 down channel 关闭。

这些原则其实只是提供的一些思路，如果不遵循这些思路，也就很容易发生 panic。比如多个 receiver，1 个 sender，直接在 receiver 端把 channel 关闭，那么 sender 发送的时候就会发生 panic。如果是在 receiver 代码中去执行 close ，那么更可能对同一个 channel 关闭多次，这种情况要解决就需要一个标记位比如 Once 只关闭一次，这样代码上会比较麻烦点。



OK，下期文章还是会结合代码来继续聊聊常见的并发 bug，防患于未然。
