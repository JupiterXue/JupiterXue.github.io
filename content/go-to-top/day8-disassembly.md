---
title: "Go 翻山越岭——反编译2"
date: 2021-08-21T15:25:05+08:00
category: tech
tags:
    - Go
keywords:
    - 实践
---

Go 语言常用的两个查看汇编代码的工具是 **go tool compile 和 go tool objdump**。上期文章说到通过 go tool compile -S 可以查看实现相同功能，两个不同的代码在性能上有什么区别，可从反编译出的汇编代码看出，二者的性能一样。今天再来聊一聊另外一个工具的使用。



> 场景3，怎么找到 make 和 new 这种 Go 语言自带数据结构的具体实现？

在解决这种类似问题之前，我们其实可以**查阅官方的资料 spec**，国外源地址：https://golang.org/ref/spec ，国内看这个：https://golang.google.cn/ref/spec



这个 spec 讲述了 Go 语言内部的语法可以怎么用。在 spec 中出现的东西，也就是官方认为是正确的东西。而如果一个用法没有在 spec 中提到，那么我们就没有办法去依赖输出结论。



举个例子，我们用到一个函数，这个函数用到了指针和 Error，如果这个返回的 Error 是非空的时候，我就不能依赖它的指针返回值。我们要依赖这个值，就需要 Error 是空。这种情况就是**语言的确定性**，如果我们依赖了有 Error 的结果，那么就是非确定性的。这个非确定性，我们会经常看到一个名词叫做 **Undefined Behavior**。我们初学 Go 语言，或者习惯于写动态语言，如 Python 的同学，会非常依赖这种写代码习惯，这其实是非常危险的。它可能会给我们带来隐藏的线上 Bug，很多时候都是难发现的，特别是在做语言 SDK 升级的时候，可能会导致比较严重的问题。



所以我们要去查询 Go 语言内部的函数或者结构的用法的话，一定要去看官方的资料。不要去看网上的博客来作为自己的结论，因为别人写的不一定对，尤其是当我们使用了错误的博客内容来操作，会给业务带来更多麻烦。



OK，打开 spec 网址，我们来看看，哥们，你这瓜多少钱一斤？

![image-20210821223659673](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108212337904.png)

spec 告诉我们，make 可以用在三种结构上的：slice，map，channel。也就说明，我们要研究 make 的实现，就去看这三种结构上执行 make 具体会执行哪个函数就好。



我们来实现这段代码，文件名为 make.go

```go
package main

func main() {
    // make slice
    // 空间开的比较大，是为了让这个 slice 分配到堆上。空间小的话会默认分配到栈上，而栈上的 slice 和堆上的 slice 底层实现会不一样。
    var sl = make([]int, 100000)
    println(sl)
    
    // make channel
    var ch = make(chan int, 5)
    println(ch)
    
    // make map
    var m = make(map[int]int, 22)
    println(m)
}
```

```bash
>>> go build make.go && go tool objdump ./make | grep -E "make.go:6|make.go:10|make.go:14"
```

命令说明，这里 go build make.go 是编译我们的 Go 程序代码，生成 ELF 可执行文件 make。然后用管道符 | 接住这个 ELF 文件，并用到反编译工具 go tool objdump 来反编译可执行文件 make。然后用管道符 | 接住结果，并用文本搜索工具 grep -E 指令来分割多个 pattern，以实现 OR 操作。



![image-20210821231924986](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108212337735.png)

初学Go语言中的汇编，我们不需要看懂所有的指令，我们只需要看代码在 runtime 中执行什么函数就好了。这里我们只需要关注 make.go 代码中第 6 行、第 10 行和第 14 行的结果就行了。接下来我们就去找汇编中 runtime.makeslice(SB)，runtime.makechan(SB)，runtime.makemap(SB) 等汇编指令都在干什么就好了。（下期讲如何看）



同样的方法，我们再来学习如何找 new 的实现，还是先问一下 spec，这瓜保熟嘛？

![image-20210821232228535](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108212338479.png)



如下文件名为 new.go （由于输出内容稍微复杂，建议大家跳读）

```go
package main

type person struct { age int }

func main() {
    var a = new(int)
    var b = new(person)
    var c = new(chan int)
    var d = new(map[int]int)
    
    println(a, b, c, d)
}
```

```bash
>>> go tool compile -N -S  new.go | grep -E "new.go:6|new.go:7|new.go:8|new.go:9"
```

![image-20210821232741718](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108212338842.png)

由于官方对 new 优化比较多，会随着版本变化，最新的版本和现在我用的 go1.14.12 差别还是挺大的。从这里我们大概知道，new 的底层实现会返回一个指针，具体的类型会根据不同数据类型进行不同的细节操作。（这个过程不太好找，需要去看编译器内部的代码，曹大建议放弃）



最后，我们只要用了合适的查看汇编的方法，就不会迷失在庞大的 runtime 代码当中，也不会迷失在别人的博客结论中。值得注意的是，虽然汇编底层会随着 Go 版本变化，网上的源码分析文章过一两年就过时了，但我们只要记住结论、知道工具以后，随时都能自己手动去做实验，通过实现去拿到当前版本的结论，肯定是最靠谱的结论。



预告：下期文章将会讲如何使用汇编调试工具来进行语法实现的具体分析。



__参考资料__

[1] Golang specification

https://golang.google.cn/ref/spec

[2] grep -E命令总结

https://blog.csdn.net/hl980703/article/details/90228724
