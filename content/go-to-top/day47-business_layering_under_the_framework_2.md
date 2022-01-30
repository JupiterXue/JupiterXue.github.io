---
title: "Go 翻山越岭——架构下的业务逻辑（2）"
date: 2021-12-30T19:36:20+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

上期文章讲了后端 API 层面的框架模式演变历史，以及介绍了面向设计 5 大原则 SOLID 中的依赖倒置原则 DIP，为今天的文章奠定了基础，来聊聊怎么具体来写代码。

# Clean Architecture

整洁架构（Clean Architecture），是由 Uncle Bob 在书《Clean Code》中提到的概念，作者也是面向对象编程中 SOLID 原则的提出者。书中还提到了涉及到了依赖倒置的一个概念图，理解之后我们能够实现编程的洋葱模式：

![image-20211231222701004](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112312227207.png)

<center>整洁架构 style</center>

- 最中心 Entities 是指企业的业务规则。
- 外面红色层 Use Cases 是指应用业务逻辑，这层就是上期文章中说到的 logic。
- 黑色箭头是指依赖的方向。整个依赖的方向是从外部到内部的。

可以发现，所有和业务相关的结构体的定义、逻辑都位于整洁架构的内层。外层有控制器 Controller、网关 Gateways、展示 Presenters。最外层是设备 Devices、数据库 DB、外部接口 External Interface、界面 UI、框架 Web。



简单来说，就是业务逻辑和业务实体都放在中心，其他与业务无关的外部资源相关的东西都放在外环上。不会让业务逻辑去依赖任何外部实现，这也是整个整洁架构最出彩的地方：不让业务逻辑和外部框架绑定。如果用 Java 或 .Net 开发是很难避免这个问题，因为它们的 web 框架很强大。而 Golang 中可能不需要考虑这些问题，因为 web 框架并没有多强大，就可以用 Go 语言去实现业务逻辑与外层实现的隔离是比较容易的。

## 整洁架构的核心观点

整洁架构中有以下几点核心观点，摘自博客《The Clean Code Blog》 https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

- **不与框架绑定**：业务不应该与某个 web 框架绑定，应该做到想换就换。
- **可测试**：业务逻辑应该在没有 UI、database 环境、web server 等所有外部环境的前提下做到可测试。
- **不与 UI 绑定**：不与具体的 UI 库绑定，项目应该做到随意切换外部 UI，比如可以将 web UI 替换为 console UI，同时无需业务逻辑做修改。
- **不与数据库绑定**：可以把 Oracle 换成 SQL Server，也可以换成 Mongo，换成 BigTable 或 换成 CouchDB。总之，业务不依赖具体的存储方式。
- **不依赖任何外部代理**：业务应该对外部环境一无所知。

然后我们再来看看，这几个核心观点在业务和代码中是怎样的设计与实现。

### 不与框架绑定

业务代码入口不应与任何协议绑定，同时框架代码（比如 gin.Context）不要入侵到业务层。

```go
// GetNotes is the logic entry for api /user/getAllNotes
func GetNotes(ctx context.Context, req *GetNotesReq) (*GetNotesResp, error) {
    // business logic
    return nil, nil
}
```

GetNotes 是业务模块的入口，如 logic 或 domain model，在这层就不应该有外层的协议，如 http.Request 或 protobuf 中的 grpc.XXX。同时也不应该有任何框架相关的字眼。



既然业务代码不能和框架绑定，意味着未来某天可能会把今天用到的 gin 框架换掉，如果业务代码中引用了，那么在换框架中肯定都会该业务代码，这是有非常大的工作量、不合适的。

### 可测试

在没有 UI、database 环境、web server 等所有外部环境的前提下做到可测试，这在 Go 语言中有一些可应用的方案：

- Sql driver mock for Golang
- Go monkeypatching
- etcd issues 11849
- Redis client Mock
- gomock
- Go 内置 net/http/httptest
- Testify - Thou Shalt Write Tests

### 不与 UI 绑定

在公司中我们做的基本上都是前后端分离的系统，不太容易与 UI 发生绑定

![image-20211231234203190](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112312342291.png)

<center>前后端分离架构</center>

### 不与数据库绑定

不与具体的 UI 库绑定，项目应该做到随意切换外部 UI，比如可以将 web UI 提花那位 console UI，同时无需业务逻辑做修改。



我们就需要借助 DDD 中的 Repo 设计方式：

![image-20211231234359689](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112312343762.png)

<center>DDD 中 Repo 设计</center>

直接抽象出来对象存储接口，然后可以用任何存储方式实现接口，比如 mysql、内存。我们的业务也不依赖存储，而是依赖接口。

### 不依赖任何外部代理

意思是说，没有外部的 agent，项目也能启动。这种方式做到是比较难的。比如我们经常会碰到：

- 没有配置下发的 agent 的环境，系统跑不起来。
- 没有 service mesh 模块，系统跑不起来。
- 没有 metrics 采集模块，系统跑不起来。



最后推荐两个资料，第一个是作者对整洁架构做了总结，告诉了大家怎么使用 clean architecture 去写代码。第二个是他的博文的实际案例：

- 《Clean Architecture, 2 years later》https://eltonminetto.dev/en/post/2020-07-06-clean-architecture-2years-later/
- 《clean-architecture-go-v2》https://github.com/eminetto/clean-architecture-go-v2



OK，下期文章聊聊比较期待的 DDD 模式！
