---
title: "Go 翻山越岭——内置数据结构-Context"
date: 2021-09-05T15:48:38+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

Go 语言在 1.16 版本之后加入了新的内置数据结构 Context，虽然在代码中使用都比较简单，但语言内部还是做了许多区分，今天来分析一下 Context。



# Context

虽然在使用 context 的时候，看起来都是 context.* 的结构，但 Go 语言内部做了这样一个区分：



![Context 操作](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109051711081.png)

<center>context 操作</center>



具体说明如下：

- emptyCtx：所有 ctx 类型的根
- valueCtx：主要为了在 ctx 中嵌入上下文数据，一个简单的 k 和 v 结构，同一个 ctx 内只支持一对 kv，需要更多的 kv 的话，会形成树形结构。
- cancelCtx：取消程序的执行树
- timerCtx：在 cancelCtx 上包了一层，支持基于时间的 cancel



这里有个例子

```go
package main

import (
    "context"
    "fmt"
)

type orderID int

func main() {
    var x = context.TODO()
    x = context.WithValue(x, orderID(1), "1234")
    x = context.WithValue(x, orderID(2), "2345")
    x = context.WithValue(x, orderID(3), "3456")
    fmt.Println(x.Value(orderID(2)))

}
```

具体流程如下，前一个节点是后一个节点的 parent

valueCtx{ k: 3, v: 3456 }  →  valueCtx{ k: 2, v: 2345 } → valueCtx{ k: 1, v: 1234} → emptyCtx



看起来像是链表，但其实 Context 更像是一棵树，这里例子更直观些：

```go
package main

import (
    "context"
    "fmt"
)

type orderID int

func main() {
    var x = context.TODO()
    x = context.WithValue(x, orderID(1), "1234")
    x = context.WithValue(x, orderID(2), "2345")
    
    y = context.WithValue(x, orderID(3), "4567")
    y = context.WithValue(x, orderID(4), "3456")
    
    fmt.Println(x.Value(orderID(3)))
    fmt.Println(y.Value(orderID(3)))
}
```

x: valueCtx{ k: 3, v: 3456 }  →

​														→  valueCtx{ k: 2, v: 2345 } → valueCtx{ k: 1, v: 1234} → emptyCtx

y: valueCtx{ k: 3, v: 4567 }  →

这样看起来就更像树结构了。



我们在开发中可能会遇到， context 影响 goroutine，伪代码如下：

```go
func () {
    go func1() {} ()
    go func2() {
        go func3() {} ()
    }
}
```

其实这也是一个树形结构，这里做了个简化就是：假设我们在每一个 goroutine 中都创建了和它对应的 context，那么可以认为我们的根就是对应的最外层的函数。内部每启动一个 goroutine ，就会对应以下这样的节点：



![未命名文件 (3)](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109051821621.png)

<center>context 影响 goroutine</center>

这种结构正好和我们程序执行的顺序是匹配的，最终生成一个 context 树。



假如我们想进一步，生成了 context 树之后，做些取消操作：比如 goroutine 都没执行完，不过想要取消下图红框所示 cancelCtx，意味着取消后，后续的子 goroutine 都需要取消，就需要调用最外部的 cancelCtx。就相当于，父节点取消时，可以传导到所有子节点中。



![未命名文件 (4)](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109051821759.png)

<center>context 局部取消</center>

但这种方式也会有一种问题，context 的操作都需要在 goroutine 中去做配合的。也就是说我们在 goroutine 中一定要用 select 然后用 ctx.dump 的操作。如果不写就相当于完全不配合外部取消操作，也就不监听外部取消通知，理论上 goroutine 还是配合不了。所以，即便有了 context 这种侵入式用法，还是需要写一些 goroutine 处理的相关逻辑。虽然整体上不是很好用，但相比其他语言，我能够简单粗暴地中断某个执行还是很不错的，可以不用考虑执行现场中怎么恢复、怎么清理当时分配的资源。



Go 语言中，不仅有这种基于 context 和 goroutine 中 select  ctx.dump 的中断方式，在内部还能够通过信号来中断。而且信号中断没有给用户暴露相关的手段，用户能够操作的也就只有常见的让程序直接停止和退出的情况。



这其实也是一种抢占式调度的方式，相当于可以从执行的汇编流中，任意的位置中断。然后把现场保存，在后续流程中如果需要就恢复回来。



最后说下 context 的社区现状。因为 context 不是很好用，社区中不少人在骂，这种侵入方式到底好不好。据说 Go 2.0.x 以后，用实现简单的 context 来做取消中断还是比较方便，而信号就没那么方便了。
