---
title: "Go 翻山越岭——编译原理基础"
date: 2021-08-15T09:38:51+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

Go 程序，从源代码到可执行文件，通常要经过两个步骤：编译 → 链接。看似简单，但其中蕴藏着极大的能量，尤其在面对不能一时半会就解决的业务场景中，稍微深挖编译与链接的过程就有奇效。今天我们就来聊聊 “Go 语言中编译原理的基础”。



上期文章提到了 **8 个常见 Go 业务场景**，虽说都是业务场景，但在代码性能、数据结构转换和系统设计中偏细颗粒度的实践中我们通常感到茫然和繁琐。这是因为在平时开发中，编译器会自动帮我们屏蔽掉代码编译的过程，我们只用关注最后的 ELF 可执行文件。那么我们就来看看，编译器会帮我们做些什么工作，进一步理解我们所写的程序发生了什么样的转化。



# 编译原理概述

我们可以把编译分为两个部分：编译器前端和编译器后端，如图所示：（这些都是我们软件工程专业课上的基础知识）

![图片](https://mmbiz.qpic.cn/mmbiz_png/MzRUv767TY1wicxlSz515XqKcibIoYHicxXHYuVQnibib4SWN9ZpsuBdLH0AvicQr0025TKnEN2PtlakkrAp2Wca0pXg/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

**词法分析（Lexical Analysis）**是计算机科学中将字符序列转换为标记（token）序列的过程。

```go
package main

import "fmt"

func main() {
    println(1 + 2)
}
```

我们这里有一段简单代码，通过词法分析的方法，转化为 token 就是如下结果：

![图片](https://mmbiz.qpic.cn/mmbiz_png/MzRUv767TY1wicxlSz515XqKcibIoYHicxXP7m8Ugza9KCQ93XPIN0qllyNLsTmIl06sBtbkvdiaQOaf0C6kA4IVQQ/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

我们在 Go 语言中没有分号，但其实像很多其他编程语言一样，是需要的。



**语法分析**（syntactic analysis，又称 parsing）是根据某种给定的形式文法对由单词序列（如英语单词序列）构成的输入文本进行分析并确定其语法结构的一种过程。当我们对 Go 源代码进行词法分析后，会形成上图的 token 流，我们想要把这个再转换成汇编还是不行，还需要转换一步，转换为另一种数据结构——AST 语法树，才能够用计算机的逻辑去处理。



**抽象语法树**（Abstract Syntax Tree，AST），或简称语法树（Syntax tree），是源代码语法结构的一种抽象表示。



这里提供一个在线编译 AST 的网站：https://astexplorer.net/，将上面的代码粘贴进来，就可以发现如下结果：

![图片](https://mmbiz.qpic.cn/mmbiz_png/MzRUv767TY1wicxlSz515XqKcibIoYHicxXHiaNN4DHmAU0tHzvhAgJjibR0BlfUQRAZZ09tLyMKNjD189kicu6NktrA/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

左边是代码，右边会把所有代码相应地转换成了树结构呈现出来。（建议大家动手玩一玩，里面会有语法高亮和代码与树一一匹配的功能。如果未来你有机会写一些解释器，很可能会用到这个）



**语义分析**（Semantic Analysis）是对结构上正确的源程序进行上下文有关性质的审查，进行类型审查。比如我们写了类似这样的代码：

```go
package main

import "fmt"

func main() {
    var x int = "abc"
    println(x)
}
```

作为一个找茬小能手（当然不是买瓜），可以发现，字符串 "abc" 是不能赋值为 int 变量 x，这时候我们直接编译 go build 就会出现如下问题：

![图片](https://mmbiz.qpic.cn/mmbiz_png/MzRUv767TY1wicxlSz515XqKcibIoYHicxXuJE6Z057MCia5QyK2PXBCKkMJ65NfQ3kvoEHvnGRic3wIKsRibIOHMHCQ/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

这个过程往往是在抽象语法树 AST 上做类型检查，不让你执行成功。



我们有了 AST 以后，要转换为汇编，还是不可以，还需要经过两步：中间代码（SSA）的生成与优化。



SSA（Single Static Assignment） 即静态单赋值，是一种中间表示形式。 之所以称之为单赋值，是因为每个名字在SSA中仅被赋值一次.

SSA是一种高效的数据流分析技术，目前几乎所有的现代编译器。其两大特点是：

1. Static：每个变量只能赋值一次（因此应该叫做常理更合适）；
2. Single：每个表达式只能做一个简单运算，对于复杂的表达式要做拆分。

我们有个例子：

```go
func foo(a, b int) int {
    c := 8
    return a*4 + b*c // 这句代码不能直接转化为汇编
}
```

我们来手动做一次简化：

```go
func foo(a, b int) int {
    c := 8
    t0 := a * 4
    t1 := b * c
    t2 := t0 + t1
    return t2 // 这个就能直接转化为汇编
}
```

手动当然不现实，那么我们有什么方法能实现自动转 SSA 呢？

1. 在线查看 SSA 过程https://golang.design/gossa
2. GOSSAFUNC=funcname go build x.go



OK，以上的编译原理是不是有点头大？从源代码到汇编，我们需要这么多过程，那么有没有什么工具能够直接从 Go 源代码到汇编呢？当然，也有：看好我施展魔法啦：https://godbolt.org/

![图片](https://mmbiz.qpic.cn/mmbiz_png/MzRUv767TY1wicxlSz515XqKcibIoYHicxXV2IvMOTUfdnN5gu92V8fmngLW89LpibUzZy0fGE9MficPb0oazLQQMEQ/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)



__参考资料__

[1] 维基百科搜索关键词

[2] 从词法分析角度看 Go 代码的组成

https://studygolang.com/articles/24419?fr=sidebar

[3] 静态单赋值（SSA,Static Single-Assignment)

https://blog.csdn.net/manok/article/details/89851085

