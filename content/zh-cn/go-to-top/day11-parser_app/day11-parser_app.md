---
title: "Go 翻山越岭——编译器应用（2）"
date: 2021-08-27T21:04:42+08:00
draft: true
---

之前说到几个 Parser 的应用场景，一个比较直观的是：在 Beego 源代码中了解了业界做得比较好的注释转 json 的工具 swagger。其实我们也可以模仿 swagger 的思路，在一个代码中去定义一些特殊的注释，从另外一个代码中把这些注释扫描出来。今天再来谈一谈 Parser 的一个更复杂的业务场景。

# 小插曲

不久前，图灵社区出版了一本书叫做活文档（living documentation），讲的就是我们可以把一些业务上相关的一些信息放在关键函数的注释中。书中主要提到的一个技巧是，注解。而在 Go 语言中不支持注解，但我们可以退而求其次，在注释中写这些信息，通过类似 swagger 的工具将注释都扫描出来，然后当成普通的字符串，按照一些匹配规则做一些匹配，然后做一些相应的内容。



但其实这种解析注释的思想，不仅仅是转json，只要我们敢想还可以做些不一样的东西。



# 转设计结构

> 场景6，公司非常看好 gRPC，尤其在上 K8s 之后 gPRC 显得更云原生些，我们想从 Thrift 切换到 gRPC，但是之前的项目都是 Thrift，怎么办？
>
> （1）之前的项目已经有几千甚至上万的 Thrift 服务、已经有单独的仓库去管理 Thrift IDL
>
> （2）想提供 gRPC 接口
>
> （3）目前的解决方案是手动把 Thrift IDL 按照 pb 文件的格式模仿抄写一遍，但效率太低。

我们已经有了大量的 Thrif 文件，如下，文件名 a.thrift

```c++
namespace php shared
    
struct ShareStruct {
    1: i32 key
    2: string value
}

service SharedService {
    SharedStruct getStruct(1: i32 key)
}
```

想要转到这样的 gPRC 结构，文件名 b.proto

```go
syntax = "proto3"

message SearchRequest {
    int32 key = 1;
}

message SearchResponse {
    int32 key = 1;
    string value = 2;
}

service SharedService {
    rpc getStruct(SearchRequest) returns (SearchResponse);
}
```

从 Thrift 转 gRPC ，需要写 pb 文件。而从代码结构上来看，Thrift 和 gRPC 结构很相似。我们如果不懂编译原理，但业务又在紧急转型，就得硬着头皮去人肉翻译的工作。而大公司的项目一般来说代码量都是几万到十万行，那么 Thrift 文件也不会少于几百行。人肉硬翻译下去的话自己就会非常痛苦。



我们去社区搜一搜可以发现，有的人已经提供了 Parser——从 Thrift 切换到 gRPC。把 *.thrift 文件解析成语法树，遍历一下，再生成目标的 **.pb 文件就大功告成了。有了 pb 文件，再用 gPRC 的客户端和服务端生成工具再去生成客户端和服务端。最后再将整个流程用代码整合，完全可以做到非人肉介入，全自动化。这效率，不知道比人肉高到哪里去了。



接下来我们再来看另外一个 Parser 的应用场景。

# SQL 审计

> 场景7，我是个 DBA，了解 MySQL 的索引机制，了解哪些数据表是数据量比较大的核心表，我也可以将用户代码中的 SQL 都提前出来，不过我希望在上线的时候能够自动做一些工作，比如做一个拦截器，提醒用户要给表加索引才能添加数据（这个问题在我们实际线上会经常遇到，一旦发生还是比较蠢的）。
>
> （1）我是 SQL 专家
>
> （2）我知道怎么获取到表的索引
> （3）我可以把用户代码里的SQL扫描出来

在 DBA 的身份下，我们需要去做些 SQL 审计的工作。其实这个工作也比较好做，同样的，我们找找可以发现，社群提供了这种业务下的 Parser。我发现了两个：

![image-20210827234431807](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108281451278.png)

左边的 Vitess 是 Yutube 开源的 MySQL 的 proc，里面有一个现成的 SQL Parser 可以直接用，而且它们的 API 非常的简单。PingCap 就是 PingCap 公司出品的非常完备的 SQL Parser 工具，因为它们想在功能上完全对标 MySQL，因此代码也相对比较复杂些。



来看看两个 SQL 审计的例子。就能很直观地知道什么是 SQL 审计了

![未命名文件 (1)](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108281451150.png)

场景1：索引推荐

![未命名文件 (2)](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108281451171.png)

场景2：统一数据接入服务 One SQL to Rule Them ALL



总结，今天了解了两种复杂的 Parser 场景，如果全都自己去实现将会非常麻烦，而通过搜索，了解了业界成熟的技术方案后，能够让我们开拓很大的思路。所以，我们学习技术千万不能闭门造车，多去看看、借鉴优秀的项目实现，让自己在公司中游刃有余，毕竟，省出来的时间，都是自己的。





__参考资料__

[1] 书评 | 如何让开发中的各种文档变活？《活文档》阅读总结

https://blog.csdn.net/turingbooks/article/details/115804565

[2] 活文档（living documentation）工具调研 --- Fitnesse

https://www.cnblogs.com/landhu/p/14109639.html
