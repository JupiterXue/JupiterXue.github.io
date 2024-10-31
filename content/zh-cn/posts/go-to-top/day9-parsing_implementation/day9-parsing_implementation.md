---
title: "Go 翻山越岭——语法分析实现"
date: 2021-08-22T08:34:33+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
typora-copy-images-to: upload
---

之前的文章提到：Go 程序编译和链接涉及到工具 go tool compile 、go tool objdump 和 SSA func  build 将源代码转化成静态单赋值形式的中间代码。这些工具都可以直接上手玩一玩，即便对编译原理的理论不熟悉，在实践后还是会逐步明白的，Go 的工程化在这方面做得挺不错。今天来讲一讲“GO 语言语法分析的具体实现”



# 调试工具回顾

通过命令 readelf  -h 可以查看 ELF 可执行文件的头信息，发现调试入口的十六进制地址码，然后用 dlv 来调试 Go 程序。



在 dlv 中，打断点有三种常用方式：

1. b * 地址
2. b 函数名
3. b 文件名:行数



指令 c 是从一个断点跳到另一个断点。如果打多个断点，那么可以做连续代码的跳转。



指令 si 是单步调试，调试汇编时常用于使用 si 到 jmp 目标位置，即一步步跳转。



我们用工具 go tool objdump 来做反汇编，而它输出的是 plan9 形式的汇编。其实在 dlv 中内置了反汇编工具，disass，不过它输出的是另一种形式的汇编。这里，我们可以掌握多种调试工具，平时就用自己擅长点的，而遇到了没弄明白的，也许另一个工具换来使用有不一样的效果。



之前我们都是遇到问题搜资料，我想系统学习 dlv 怎么办呢？就像上期文章讲到的，还是去看官方文档：https://github.com/go-delve/delve/tree/master/Documentation/cli



简单浏览一下官方文档可以发现，官方也在逐步添加新的功能。比如说，现在有一个地址，我可以直接用 dlv 中的 x 指令查看一段连续内存里存储的值，这个有点像 gdb 中的 x（另一个调试工具，用于查看内存地址的值）。这个指令在 runtime 中有些开头是 len，然后跟着 unsafe pointer 之类的，我能看到 unsafe pointer ，但它后面的结构直接调试可能看不到，这时候用到 x 指令就可以看到它后面内存里存储的是什么值了。



# 语法实现分析

我们已经储备了很多的调试工具，现在就是让这些工具派上用场的时候了。



来分析一下我们写的第一个 go func 简单程序，文件名为 gofunc.go

```go
package main

import "time"

func main() {
    go func() {
        println("hello world")
    }()
    
    time.Sleep(time.Second * 5)
}
```

在通过反编译和文本过滤命令

```bash
go tool compile -S gofunc.go | grep "gofunc.go:6"
```

![image-20210822154943813](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108221618914.png)

通过定位第 6 行代码结果，我们找到了被翻译出来的关键结果：runtime.newproc(SB)。这里就将之前文章中提到的创建 goroutine 就串起来了。



因此，这里做个小结，所有的 runtime 函数都可以通过上述方式来找关键信息。



如果我们对 channel 的 send  和 recv 感兴趣，也可以如法炮制，文件名 chan.go

```go
package main

func main() {
    var a = make(chan int, 1)
    a <- 666
    
    x := <-a
    println(x)
}
```

```bash
go tool compile -S chan.go | grep -E "chan.go:5|chan.go:7"
```

![image-20210822155941901](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108221618915.png)

从图中可知，我们要找的关键函数就是：runtime.chansend1(SB), runtime.chanrecv1(SB)。



我们除了这种普通的 chan 收发，在平时应用中多和 select 联合起来用，这里还是一个简单例子，文件名 nonblock_recv.go。

```go
package main

func main() {
    var ch1 = make(chan int)
    select {
        case <- ch1:
    default:
    }
}
```

我们从 select 中获取 channel 的值，如果当前 channel 阻塞了，我直接走 default 分支跳出来了。来看看反编译后的结果

```bash
go tool compile -S nonblock_recv.go | grep "nonblock_recv.go:6"
```

![image-20210822160730403](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108221618916.png)

这里的关键信息 nbrecv，不是 nb recv 而是我们代码 nonblock_recv 对应的底层信息，完美对应上了。



今天，我们举了三个例子来学习如何进行 Go 语法实现分析，学会了没？这里留下个问题，用上今天的方法来找到以下三个实现，摘自书籍 《concurrency in go》

![image-20210822161317456](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108221617934.png)

分别有三种情况的 panic：

1. 往已经关闭的 channel中写入数据
2. 关闭一个是 nil 值的 channel
3. 关闭一个已经是关闭状态的 channel



下期预告：讲述以上三个 panic 在 runtime 代码的具体哪个位置输出，关键信息是什么。
