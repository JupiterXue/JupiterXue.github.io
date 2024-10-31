---
title: "Go 翻山越岭——并发编程模式"
date: 2021-12-03T20:32:49+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

通常我们在写 Go 代码的时候，听得比较多的是 Go 和传统并发模式不一样。Go 语言采用的是 CSP 模式，即 "Do not communicate by sharing memory; instead, share memory by communicating."



这句英文看起来有些玄乎，但其实是在说，除了 runtime 中的 hchan 结构，没有其他结构能够共享内存。而如果我们想要共享数据的话，就需要 chan 去做任务发送，这时候做的其实是一些拷贝的工作。听起来有点抽象，来看一些别人写的代码示例。



第一个例子 Fan-in 扇入，也就是多个分支合并为一个分支，或合并为少数的分支。

![image-20211203230406497](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112032304677.png)

<center>图 1-1：Fan-in 合并多个 channel 例子</center>

![image-20211203230942784](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112032309871.png)

<center>图 1-2：Fan-in 合并多个 channel 例子</center>

这里的代码意思是要把多个 channel 的结果合并到同一个 channel 中去，这样消费者只用去消费一个就好了。函数中第一行是通过判断输入的 channel 数组的长度来实现的：

- 如果是 0，说明递归结束直接返回；
- 如果是 1，就返回当前的 channel；
- 如果是 2，就需要调用 mergeTwo 函数。这个 mergeTwo 函数也就是简单的 for 循环 + select 结合来实现的。
- 最后 default，就是不断对 channel 的长度除 2 ，然后递归地调用 mergeTwo 最终来实现。

一般情况下，我们要合并 channel 的话，数量不会很多几百上千个，常见的是 5-6 个。并且更加实际环境中，每个 channel 还有其具体的名字的，也就只需要写一个 select 就可以了，不用像以上代码中不定长的 channel 实现方式要简单得多。



第二个例子 Or channel，意思是说有多个 channel，其中一个有结果的时候，不管其他 channel 的结果，直接把整体的结果返回。这个一般是在有多个后端向前端响应结果，并且一旦有个后端返回了就不用等另外个后端返回的场景下去用，也就是通过用更多资源的方式来降低延时、提升用户的体验。Or channel 的实现和这种实现类似，代码如下：

![image-20211203231732542](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112032317626.png)

<center>图 2：Or channel 例子</center>

这个代码的逻辑稍微复杂，具体是这样的：

- 如果当前 channel 数组的长度为 0，直接返回空；
- 如果是 1，就返回当前的 channel；
- 如果是 2，就用 select 和两个 case 就解决了。
- 如果是 3个以上多个 channel，就采用 select 的  0、1、2 case，其他的就递归去调用 or 函数来实现。

以上代码也是《Concurrency in Go》书中的代码，也建议大家通过抄一下去学习。



第三个例子，将 channel 串联起来。 

![image-20211203232747010](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112032327076.png)

<center>图 3：Pipeline 例子</center>

这个代码中逻辑比较简单，也就是一个 channel 收到了然后做一些处理，然后做些处理，然后再发往下一个 channel。这种例子一般用于上下游的数据处理，通过将 channel 串联起来，解决速度不匹配的问题。



第四个例子，是由 PingCap 工程师分享的案例：想要通过 Goroutine 并发去提升请求的性能，但是在业务场景下，需要按照顺序，把数据处理后返回给用户。所以既要通过并发提升性能，又能够保证返回的顺序性。

![image-20211203232839228](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112032328302.png)

<center>图 4：并发同时保序例子</center>

这个代码逻辑如下：

- 用户提供了一个 task 任务，按照用户提供的 task 顺序分别发送了两个目标：一个是并发执行任务的 worker，另一个是保证顺序的函数。
  - 这里保证顺序的函数可以任意去实现，如数组或链表都是可以的，只要保证最后在消费的时候，是从这个保证顺序的列表或 channel 去消费就 OK 的。
- 在进行消费的时候 fifo，每次消费一个元素，就执行 t.wait，也就是上方的 sync.WaitGroup。
  - worker 函数中，每次执行完成任务之后执行 t.Done 就 OK。

这个例子比较好玩，网上别人也很少讲到。



OK，其实以上的例子与代码大部分都来自于《Concurrency in Go》和鸟窝在 Gopher China 分享的内容，比较多，这里给一下地址插个眼：

[1]《 Concurrency in Go》.pdf

https://github.com/chapin666/books/blob/master/golang/Concurrency-in-Go.pdf

[2]《 Concurrency in Go-中文版》.pdf

[https://github.com/hapi666/GOBook/blob/master/Concurrency%20in%20Go%E4%B8%AD%E6%96%87%E7%89%88.pdf](https://github.com/hapi666/GOBook/blob/master/Concurrency in Go中文版.pdf)

[3] 《Go并发编程实践》博客总结，鸟窝，

https://colobu.com/2019/04/28/gopher-2019-concurrent-in-action/



OK，下期文章继续总结一下常见的并发 Bug，让我们防患于未然。
