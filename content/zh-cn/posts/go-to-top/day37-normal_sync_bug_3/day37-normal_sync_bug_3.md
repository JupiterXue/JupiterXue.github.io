---
title: "Go 翻山越岭——常见并发 bug（3）"
date: 2021-12-10T21:18:46+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

![img](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112102120452.png)

这只地鼠又来了，说明今天我又来讲 Go 语言并发中常见的一些bug，看代码咯。

## 引用传递

```go
for i := 17; i <= 21; i++ {  // write
    go func() { /* Create a new goroutine */ }   // before modify
    go func(i int) {  // modified
        apiVersion := fmt.Sprintf("v1.%d", i) // read
    }()   // before modify
    }(i)  // modified
}
```

这里说的是，更改代码之前，我们启动一个 goroutine，这个 goroutine 用的是一个闭包，闭包捕获了外面的变量 i，而这个 i 用的还是地址。而迭代器 for 用的也是同一个 i 的地址，所以到 goroutine 执行的时候，最后 Print 出来的 i 就是最后一个了。



这里的修改方式比较简单，就是把 i 当作参数传过去，因为是值传递，也就解决问题了。

## WaitGroup

```go
func (p *peer) send() {
    p.mu.Lock()
    defer p.mu.Unlock()
    switch p.status {
    case idle:
        p.wg.Add(1)   // modified
        go func() {
            p.wg.Add(1)   // before modify
            ......
            p.wg.Done()
        }()
        case stopped:
    }
}

func (p * peer) stop() {
    p.mu.Lock()
    p.status = stopped
    p.mu.Unlock()
    p.wg.Wait()
}
```

这里代码还是涉及到 WaitGroup 的用法，在修改代码之前，Add 是放在 go func 中，有可能 WaitGroup 依旧是 0，WaitGroup 的 Wait 就不需要等待任何 goroutine 就能执行完成，整个程序也就执行结束了。



因此，在启动 goroutine 前要保证 Add 完成，将 Add 放在 go func 之前就能使得整个逻辑在不同条件下正常执行。

## 重复关闭 Channel

```go
select {   // before modify
    case <- c.closed:   // before modify
    default:   // before modify
        Once.Do(func() {   // modified
            close(c.closed)
    })   // modified
}   // before modify
```

这代码在执行并发操作 channel 时，多次关闭同一个 channel。这种情况也是我们平时开发中最常见的问题，重复关闭 channel。



为了解决代码逻辑有误的情况，又得额外去打一些补丁。比如这里的 Once.Do 都是之后修改的代码，说明之前的 select 可能会进入多次，因此就会对这个 channel close 关闭多次。编译器就会抛出 channel panic 的错误。



```go
ticker := time.NewTicker()
for {
    select {   // modified
        case <- stopCh:   // modified
            return   // modified
        default:   // modified
    }   // modified
    f()
    select {
        case <- stopCh:
            return
        case <- ticker:
    }
}
```

最后个代码例子，可能不是很直观。这段代码是中，是另一段 goroutine 代码来向 stopCh 的 channel 发送通知数据，在这段中接受的。



在修改之前，这里的意思其实是程序有可能在执行 f 函数，而这个 f 函数内部逻辑比较复杂，时间复杂度比较高，需要计算半个小时之类的。当外部已经通过 stopCh 通知需要停止 f 函数的逻辑了，也就需要退出整个循环，而不是再回到循环，然后进入到 f 函数中。也就是说 Fn 耗时很久，但进入之前没有判断外部给的 stopCh 中的通知而浪费了算力。



修改后的代码就在 f 函数之前，有一次判断，能够提前退出。不用再等 f 函数执行完成，再来接受和判定终止的通知。

## 小结

到目前为止，通过三期的文章，已经将 Go 语言中并发部分的一些常见的 bug 都梳理了一遍，其中一些代码例子还是挺经典的。作为 Go 开发者，这些问题虽然我们没有全部都遇到，其中一两个发生在自己身上也是很正常不过的。



以上的例子其实都来自于一篇 Go 语言的学术型论文当中。没想到学术界也会写一些工程界关于 bug 的论文，还挺神奇的。所以经过我三篇文章的介绍，给我们的启示是，去研读最新的科技论文，其中不仅有较为前沿的科技理论，也可能会有偏工程性的案例研究。学术无涯，研究无界，保持一颗热爱技术的心，无论是理论还是工程都是有研究的意义和价值的，感谢这样的研究者。我们专注于工程中的工程师也应该向他们多多学习和借鉴，不拘泥于做好自己手头上的事情，也可以适当探讨更好的理论解决方案，为科技世界一同贡献一份力。
