---
title: "Go 翻山越岭——编译与反编译"
date: 2021-08-19T23:02:38+08:00
category: tech
tags:
    - Go
keywords:
    - 实践
---

Go 程序，从源代码到可执行文件，通常要经过两个步骤：编译 → 链接，而**最重要的就是进行虚拟地址重定位（Relocation）**。上期我们了解了GO语言编译过程包括六个细致的编译步骤，而链接是最后的编译过程。知道了整个流程后，今天再来聊一聊 Go 程序在完成编译后，我们还可以通过哪些工具知道一些信息——编译与反编译。



重定位(relocation)是指把符号引用与符号定义链接到一起的过程。当程序调用一个函数时，将从当前运行的指令地址处跳转到一个新的指令地址处去执行。我们平时在编写程序时，只需指明所要调用的函数名即可，对编译器来说就是指明了符号引用。然后在重定位过程中，动态链接器会把函数名与函数实际所在的地址（符号定义）联系到一起，进而使得程序能够知道应该跳转到哪里去。



具体而言，GO 程序在编译后，所有函数地址都是从 0 开始，每天指令是相对函数第一条指令的偏移，并在进行链接后，所有指令都有了全局唯一的虚拟地址。（注意：不是上方框选地址，这里是指某个函数的偏移量，而是箭头下方的地址，虚拟地址）



![image-20210819234132142](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210819234132142.png)



# 编译与反编译工具

我们学习掌握 Go 语言编译原理理论的过程，会相对有些枯燥。不过通过实践，我们能够很好地将所学到的理论融汇贯通，并且即便你对理论不是很熟悉，也能够玩转这些工具。我会提到平时我们在研究 GO 汇编底层经常用到的工具，让我们来看看它是什么，跟着我一起动动手。

```bash
go tool
```

![img](C:\Users\Xfavor\AppData\Local\YNote\data\weixinobU7VjplR8GhsqZn6qPSNHGlLvOg\bf1978b1b2864b888597f55ca77551b3\clipboard.png)

这里我标注出了两个最常用的工具 go tool compile 和 go tool objdump。在我们卷 Go 汇编底层的时候，业内还有这么一句黑话：

> 几乎没有任何一个 Go 汇编底层问题不是用一条 go tool compile 不能解决的，如果不行的话，就用 go tool objdump，总能知道是怎么回事。

再让我们看看这两句指令在具体场合下都是什么意思：

> go tool compile -S main.go  # 反编译代码为汇编代码。

> go tool objdump # 可用于查看任意函数的机器码、汇编指令、偏移。（go 源码下面有个 cmd/internal/goobj包，可以读到.o文件的重定向信息，更好）
>
>  
>
> go tool objdump [-s symregexp] 二进制 # Objdump 打印二进制文件中所有文本符号（代码）的反汇编。如果存在 -s 选项，objdump 只会反汇编名称与正则表达式匹配的符号。



我们来编写一个小案例进行汇编调试实践：

```go
package main

import "fmt"

func main() {
    var a = "hello"
    var b = []byte(a)
    fmt.Println(b)
}
```



通过反汇编工具来查看一下：

```bash
go tool compile -S ./hello.go | grep "hello.go:5"
```

![image-20210819235112398](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210819235112398.png)

这条命令的意思是，产生 .o 目标文件，并且把目标的汇编内容输出出来。后面跟上管道，用 grep 截取出来 hello.go 编译出来第五行的汇编代码。这里大部分汇编代码其实我们不需要看懂，我们只需要它调用了其中某个函数就知道了。



对于初学者来说，学习Go汇编不需要全都懂，只需要一些关键路径在干什么就行。比如这里最关键的是 "runtime.stringtoslicebyte(SB)"，即我们要把 string 转换成 byte 数组，底层会调用这个函数。



掌握了这个方法我们就能解决之前文章中提到的第一个问题：

场景1，这两段代码运行速度怎样？第一个比第二个快？

```go
// 代码1
package main

type person struct {
    age int
}

func main() {
    var a = &person{111}
    println(a)
}
```

```go
// 对比代码2
package main

type person struct {
    age int
}

func main() {
    var b = person {111}
    var a = &b
    println(a)
}
```



我们来看看第一个代码第 8 行编译后变成啥了：

![image-20210819235718660](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210819235718660.png)



再来看看第二个代码第 8 和第 9 行：

![image-20210819235743177](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210819235743177.png)



由此，我们可以得出结论：一行版本的代码和两行版本的代码最终编译出的结果是完全一致的，没有任何区别。



掌握了 go tool 中两个常用工具，我们在 Go 语言底层汇编的道路就打开了。下期我还将讲解场景的 8 个业务场景如何应对。



__参考资料__

[1] 重定位- elf文件类型四Relocation

http://www.360doc.com/content/15/1126/20/7377734_516130511.shtml

[2] go tool 命令

[https://louyuting.blog.csdn.net/article/details/101014450](https://louyuting.blog.csdn.net/article/details/101014450?utm_medium=distribute.pc_relevant.none-task-blog-2~default~BlogCommendFromBaidu~default-3.control&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2~default~BlogCommendFromBaidu~default-3.control)

[3] go tool objdump

https://my.oschina.net/u/4586289/blog/4634710

[4] 命令 objdump

https://golang.google.cn/cmd/objdump/
