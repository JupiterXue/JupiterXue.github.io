---
title: "Go 翻山越岭——常见并发 bug（2）"
date: 2021-12-08T16:20:11+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---
上期文章说到了一些常见并发 bug 的场景，本期文章继续聊聊 Go 语言中常见的一些 bug，继续写 bug（当然不是），并通过代码案例来讲解。

![img](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112081853217.png)

## K8S

这里有段一篇论文的代码，这篇论文总结了业界比较知名的 bug，有一些是 k8s、docker 等耳熟能详的应用。因此，我们也可以发现，即便是 Google 的工程师，他们还是会写出来一些 bug，

```go
func finishReq(timeout time.Duration) r ob {
    ch := make(chan ob)  // before modify
    ch := make(chan ob, 1)  // modified
    go func() {
        result := fn()
        ch <- result // block
    } ()
    select {
        case result = <- ch:
        	return result
        case <- time.After(timeout):
        	return nil
    }
}
```

这段代码的思路其实很简单：我现在想往外部发一个请求，同时要控制超时。如果请求超时那么就会给用户返回一个空。



这里前两行分别有两段注释，说的是官方修改前的代码 channel 没有缓冲，还有段是修改后的代码 channel 中加入了 1 的缓冲。如果是未修改前的代码，在整个程序跑起来后，我们可以自己先想一下可能会发生什么问题？



在启动了 goroutine 去向远端发起请求的时候，如果发生了超时，就会触发 select 直接返回空。说明 channel 的对端已经没有消费者在等待结果，而没有缓冲时，channel 的 buffer 为 0，那么往 channel 发送结果的 goroutine 一定会阻塞，也就意味着这个 goroutine 永远释放不了，最终造成 goroutine 泄露。也就是超时一次就泄露一个 goroutine，如果超时越多，泄露的 goroutine 涨幅也越多。

## sync.WaitGroup

这个是 sync.WaitGroup 的 bug

```go
var group sync.WaitGroup
group.Add(len(pm.plugins))
for _, p := range pm.plugins {
    go func(p *plugin) {
        defer group.Dont()
    }
    group.Wait()   // before modify
}
group.Wait()   // modified
```

看到代码大家可以思考一些这段代码有什么问题。



在此之前，我们先回顾一下 sync.WaitGroup，它的基本用法中 Add、Done 和 Wait 其实是不能并发的。虽然 Add 和 Done 可以并发，但 Wait 是不可以和另外两个一起并发。



这段代码中，我们原意是想让所有逻辑都运行完之后再退出，而事实上 Wait 和 Done 产生并发了，也就有可能没有执行等待，直接执行下方的逻辑，并且 Wait 永远退出不了。因此必须把 Wait 拿到外面去，让 Wait 单独地执行，而不是并发执行。

## context.WithCancel

这个是个 context 的 bug

```go
hctx, hcancel := context.WithCancel(ctx)  // before modify
var hctx context.Context   // modified
var hcancel context.CancelFunc   // modified
if timeout > 0 {
    hctx, hcancel = context.WithTimeout(ctx, timeout)
} else {   // modified
    hctx, hcancel = context.WithCancel(ctx)   // modified
}
```

老代码的逻辑是如果有 timeout > 0，就把原来的 context 覆盖掉就结束了。大家想想有没有什么问题？



在第一次生成 context 的时候，其实在底层会生成 goroutine 的。而当 timeout > 0 直接覆盖掉原来的 context ，就相当于原来 context 中的 goroutine 完全没有办法去做控制了，最终导致这个 goroutine 泄露了。所以这种覆盖的方式其实是有问题的。



后期维护的代码也比较直观，根据是否有 timeout，分别生成两个 不同的 context，并且不再有覆盖的操作。

## 死锁

最后一个死锁的例子：

```go
func goroutine1() {
    m.Lock()
    ch <- request // blocks,    before modify
    select {   // modified
        case ch <- request   // modified
        defautl:   // modified
    }   // modified
    m.Unlock()
}

func goroutine2() {
    for {
        m.Lock()  // blocks
        m.Unlock()
        request <- ch
    }
}
```

这段代码看起来较为简单，第一个 goroutine 中，全局变量 m 先抢占锁，然后 channel 发送值；第二个 goroutine 中，全局对象 m 抢占锁，再释放锁，channel 接收值。也就是说，如果先进入 goroutine1，这个 channel 没有缓冲，就会直接卡住。



修改后的代码是说，如果 channel 还不能暂时发送、对端没有阻塞了的话，那么就可以进入  default 结束，并且全局对象 m 释放锁。这样 goroutine2 接收端就可以进去。如果接收端可以进去了，发送端也便可以发数据了。这个例子算是个基本的死锁逻辑了。



OK，对于并发 bug，我们看得越多越容易在实际项目开发中定位到问题，下期文章再来看一些例子并做个小总结。
