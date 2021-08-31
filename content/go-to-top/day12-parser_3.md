---
title: "Day12 Parser_3"
date: 2021-08-28T15:14:37+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

上期说到两种复杂的 Parser场景，通过借鉴业界成熟的技术方案，能够大大节约我们自己的时间。今天再来谈谈最后一个 Parser 应用场景和 Parser 总结。



# 客户定制需求

> 场景8，我们公司是做 ToB 的服务，我们的软件会先编译完再卖出去，比如我们公司是用 Go 开发的，假如我们是卖 gRPC 服务的，而用户想要在此基础上做一些定制，他们也有自己的开发人员，比如说实现内部的 RPC 协议。用户想定制，而我们想不给他们源代码并且能够支持他们的定制需求，这就需要些技术手段来解决。（特别是 Go 语言的生态还不是很完善，并且每个方案都有自己的缺陷）

我们在给客户交接系统的时候，客户想要自己的 RD去做些开发工作（一般指研发工程师）。而我们知道，Go 程序需要编译成二进制文件才能使用，因此让 Go 模块去做热更新是非常难的。为了满足这样的需求，我们应该怎样来做扩展呢？业界目前有四种方案：

1. RPC
2. go-plugin
3. REPL，社区的解释器方案
4. WASM

以上**四种方案都有不足，具体还是要看我们在真实业务场景中能够接受哪种折中的**。



RPC 中，在两个系统之间定义好交互的 API，要做扩展那么就将新系统中所有的扩展做成 API，要实现这个扩展，就在原系统中调用一下就好。但 RPC 自身也有一定**性能问题**，如做些网关之类的扩展就不适合。



go-plugin 在Go 1.8版本就有，但也是有缺陷的：不同版本编译不兼容。用户编译出的 plugin 和我们的原始二进制文件必须要用相同的版本去编译才能通过。



REPL，社区的解释器方案例如 gopherjs、gopherlua，这些都是不错的项目，在可以在线上去用这些东西，并给用户提供局部热更新的功能。但问题就在于，有些业务需求不一定社区有。



WASM（webAssemble）是一种不针对特定平台的二进制格式文件，也是目前比较新的方案。理论上在每个语言直接互相调用没有性能损失的，但 Go 语言版本不是很完善，所以这个方案也是在调研阶段。有些公司只是在做些宣传、



# Parser 总结

OK，到目前为止，我们了解了编译和反编译工具去找 Go 程序的语法具体实现，以及用到 Parser 知识可以在实际公司的业务里面去做些探索。但说到底，如果我们掌握了编译原理的知识就能够在很多场景下找到比较灵活的方案。



对于像笔者一样大学毕业工作不久，因工作需要转到用 Go 语言，而 Go 语言入门又没多久的同学来说。大学期间我们学习编译原理感到 枯燥、困难，略有些脱离实践，因此学习得不扎实，勉强通过。而工作期间又确确实实需要这方面实打实的知识，怎么办？边做边看。既然大学的时光已经过去，一去不复返，那么能抓住的就是现在。业务上需要什么我们去看什么，同时多去逛逛技术社区和论坛，比如汇总整理了几个有意思 Parser 项目的项目：https://github.com/cch123/parser_example，将自己的技术面拓宽，尽可能地边实践，边回顾理论。



之前总结了 8 个 Parser 场景，其实我也只经历了两个的样子。经验非常受到局限，但是我们通过学习，通过学习技术前辈曹大所经历的、所见识的，就能让我们也坐上通往技术前方的高铁。不知道，不擅长没关系；先了解，多了解，以后总还有可能遇到。到时我们就不再是愁思苦想，而是胸有成竹了。



__参考资料__

[1] Go进行wasm编程, 博客园

https://www.cnblogs.com/codingnote/p/12453872.html