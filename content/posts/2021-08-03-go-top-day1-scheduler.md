---
title: "Go 翻山越岭——调度-day1"
date: 2021-08-03T09:38:51+08:00
category: tech
tags:
    - Go
keywords:
    - Scheduler
---

### 前言
Go 语言是一门入门容易深入难的编程语言。说得好像其他编程语言就有入门难深入简单的？其实每一门编程语言要深入学习都会有不一样的壁垒与门槛，只是看作者如何挖掘。

<!-- more -->

## Golang 四座大山
作为 Go 语言进阶的一个标识，是磕磕绊绊地翻过 runtime 四座大山，包括：
> Scheduler	：调度器管理所有的 G，M，P，在后台执行调度循环
> 
> Netpoll	：网络轮询负责管理网络 DF 相关的读写、就绪事项
> 
> Memory	：当代码需要内存时，负责内存分配工作
> 
> Garbage	：当内存不再需要时，负责回收内存

而这些模块中，最核心的就是 Scheduler 调度，它负责串联所有的 runtime 流程。“要去到哪里,就从哪里开始”，既然选择学习 Golang，我们就直面这个第一个难题。

## 调度器的发展历史
如果你学习数学感觉困难，那么可以去读数学史。如果学习 Go 语言调度感到困难，那么我们先老了解调度的发展历史。再反过来想一下，如果你是一个领域的专家，是否应该比较熟悉该领域的发展历史？ 哲学家培根说过：“读史使人明智，读诗使人灵秀，数学使人周密，科学使人深刻，伦理使人庄重，逻辑修辞之学使人善辩“。那么就开始我们的调度器发展史：

## 什么是调度[^3],[^4] 
举个例子，十字路口的红绿灯，就是一种调度系统。因为车速过快，人工去做疏导指挥的效率太低而采取的一种自动化的处理流程，为了限制这些车辆不随意行驶，便有了红绿灯调度系统。
Go 调度是为了多个协程能合理的利用线程。这里的协程（goroutine）相当于车辆了，线程相当于十字路口。程序在运行时，会运行很多协程 goroutine，也就是我们常说的并发，为了保障这些协程能够有序快速地在线程上执行，这时候就引入可调度 Scheduler。
一个简短的例子
每当我们写下这样形式的代码，到底发生了什么事情呢？
```go
go func() {
println("hello world in goroutine!")
}
```
这里我们是向 runtime 提交了一个计算任务，并且 func() { xxx } 里包裹的代码，就是这个计算任务的内容。
Go调度流程的本质
也许这里还不清楚，但请记住，Go 的调度流程本质上是一个生产—消费流程。在后面的讲解中你会更加深刻了解到这句话的深意。
这里有两个动画便于你理解：

Goroutine 生产过程[^1],[^2]

Goroutine 消费过程

想要学好编程，就需要下功夫，多投入时间去理解与实践。想要学好 Go 语言，就需要跋山涉水，为你我愿意翻山越岭。
（未完待续）

__相关链接__

[^1]: Go 高级工程师实战营第一课, 草春晖

[^2]: [go-scheduler](https://qcrao.com/ishare/go-scheduler/), 饶全成, 博客

[^3]: [goroutine 和线程的区别](https://golang.design/go-questions/sched/goroutine-vs-thread/)，饶全成, 博客

[^4]: [深入理解Go语言(03)：scheduler调度器 - 基本介绍](https://www.cnblogs.com/jiujuan/p/12735559.html), 博客园