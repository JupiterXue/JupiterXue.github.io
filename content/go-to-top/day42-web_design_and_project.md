---
title: "Go 翻山越岭——web 框架原理"
date: 2021-12-24T19:14:59+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

今天来讨论一下流行、常见的具有 RESTful API 风格的 web 框架中的原理，涵盖了 gin、go-chi 和 ego 的一些组件。我们会发现，它们框架中的组件，或多或少会有以下几个设计：middleware、router、validator、request binder、sql binder，。或许名称不都叫这个，但其实现的逻辑大概是差不多的。下面就来具体讲解它们的实现。

## Middleware

### 实际场景

在讲 middleware 底层实现之前，我们先来看看它在我们实际项目中是怎么用到的：

```go
package main

func hello(wr http.ResponseWriter, r *http.Request) {
    wr.Write([]byte("hello"))
}

func main() {
    http.HandleFunc("/", hello)
    err := http.ListenAndServe(":8080", nil)
    ......
}
```

这里我们有个 hello world 的业务逻辑，监听了 8080 端口，配置了一个 http 的 Handler，只要访问了路由 / 就会进入 hello 的逻辑。内部的逻辑这里只有个简单的写动作，实际场景中还会不断地叠加业务逻辑。



某一天新的业务来了，想要统计每个请求所花费的时间，怎么实现呢？

```go
// mindware/hello_with_time.elapse.go
var logger = log.New(os.Stdout, "", 0)

func hello(we http.ReponseWriter, r *http.Request) {
    timeStart := time.Now()
    we.Write([]byte("hello"))
    timeElapsed := time.Since(timeStart)
    logger.Println(timeElapsed)
}
```

最简单粗暴的就是在 hello 逻辑中手动添加统计时间的过程，在向 http.ResponseWriter 写入  hello 逻辑之前记录当前的时间，在处理完业务逻辑以后，再记录消耗过的时间，并且把中间消耗过的时间打到日志中。



随着业务的迭代，接口肯定会逐渐增加。我们一个模块不可能只有一个诸如 hello 这样简单的接口，还会有各种各样的接口。虽然大多数公司都是微服务架构，但一般一个模块中至少有 10 个以上的接口，这种“笨笨”的办法就不是很适用了。哪怕现在年轻体力好，有 100 个接口，也可以写这样的代码一百遍，随着公司发展壮大后，除了要把请求的耗时写到日志里之外，可能业务还需要将耗时上传到可视化监控当中，代码就会发生如下改动：

```go
func hello(we http.ReponseWriter, r *http.Request) {
    timeStart := time.Now()
    we.Write([]byte("hello"))
    timeElapsed := time.Since(timeStart)
    logger.Println(timeElapsed)
    // 新增耗时上报功能
    metrics.Upload("timeHandler", timeElapsed)
}
```

### 洋葱模式

我们有 100 个接口都有这种功能修改需求，而如果有几百上千个接口，每个接口都要改一遍？所以我们有了中间件的思想，有个框架中叫做 middleware 有点框架中叫 filter (Java 框架中多这么叫，比如 Spring)，中间件的本质其实就是实现了 23 种设计模式中的一种，责任链模式 (Chain of Responsibility Pattern)，或者叫做拦截器模式，又或者叫装饰器模式，还有的地方称之为代理模式、洋葱模式等等，它们的实现都是差不多的东西。



实现责任链模式的基本思路是，把业务代码中功能性和非业务代码中非功能性的代码分离。具体实现如下：

```go
func hello(wr http.ResponseWriter, r *http.Request) {
    wr.Write([]byte("hello"))
}

func timeMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(we http.ResponseWriter, r *http.Request){
        timeStart := time.Now()
        
        // next handler 
        next.ServeHTTP(wr, r)
        
        timeElapsed := time.Since(timeStart)
        logger.Println(timeElapsed)
    })
}

func main() {
    http.HandleFunc("/", timeMiddleware(http.HandlerFunc(hello)))
    err := http.ListenAndServe(":8080", nil)
    ......
}
```

我们要实现统计耗时的中间件，声明一个叫做 timeMiddleware 的函数，它接受一个 http.Handler 的参数，并且返回同类型的参数。其中，传入参数 next 就是实际业务的 Handler，那么我们在返回这个 Handler 之前，先计一个时，在 next.ServeHTTP(wr, r) 执行用户传进来的整套业务逻辑流程，执行完成之后再去记录新的时间和之前的时间做对比，最后记录下整个逻辑消耗的时间日志。



这种过程之后，相当于业务逻辑之中几乎没有非业务逻辑相关的内容了，依旧是原来简洁清爽的 hello。



这段代码在阅读上不是很清晰，他传入了一个 Handler 又返回了一个 Handler。本质上在执行一个具体的 http.Handler 的时候其实执行的是返回的函数逻辑，这个函数会先进行非业务的处理，然后执行业务逻辑，然后再执行业务逻辑。所以在刚刚提到的几个名称当中，其实最合适的还是叫做洋葱模式。

### 多层嵌套

如果我们还是按照这个思路把业务逻辑包裹起来，那么不仅可以包裹一层，还能包裹多层，比如除了打日志，还有超时、限流等，一层层地嵌套。

```go
customizeHandler = logger(timeout(ratelimit(helloHandler)))
```

每一个 middleware 的实现都是把原来的函数包装了一下，返回了新的函数。所以一个套一个的方式是可以完成的，并且我们在最后挂在在真实路由上的一定是最后的 Handler。



而这种多层嵌套结构在项目执行过程中的处理流程是这样的：

![image-20211224201731104](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112242017451.png)

<center>多层嵌套的处理流程</center>

- 程序一开始会进入到最外层的 logger 的 middleware 的上半部分逻辑，执行完成之后会继续往内层的 middleware 执行，然后重复这种逻辑。然后就类似于剥洋葱一样，你是一层层地剥开，进入到最里面的业务逻辑。
- 在业务逻辑执行完成之后，会从这里返回，返回的过程还是类似剥洋葱一样，先执行最内层的 middleware 的下半部分逻辑，执行完成之后会继续往外层的 middleware 执行，不断重复这种逻辑。
- 执行到最外部 http response 的时候，相当于整个处理逻辑就处理完成了。

以上就是 web 框架中如何实现一个、多个 middleware 的完整逻辑流程。

### Handler 原理

在上面的讲解中我们知道了 middleware 可以传入一个 Handler 然后返回一个 Handler。实现这种功能主要是 http.Handler 和 http.HandlerFunc，我们可以看一下源码结构：

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (f HanderFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

http.Handler 它本质上是一个 interface，也就是说当我们的对象实现了 ServeHTTP 的函数功能，也就相当于实现了 Handler 的接口。



HandlerFunc 也是定义在 http 库中，它的签名就是场景的 http 的入口。



在内部是给 Handler 对象实现了 ServeHTTP，也就意味着我们可以返回一个 HandlerFunc 可以作为 http.Handler。虽然我们之前定义的 hello 本身不是 HandlerFunc，但它的签名和 HandlerFunc 的签名不一样，那么会经过临时匿名函数 func 的一个强制类型转换过程。

### 具体实现

我们有了这么多工具以后，如果代码依旧这么写：

```go
customizeHandler = logger(timeout(ratelimit(helloHandler)))
```

看起来不仅丑陋，还不好理解，所以我们大多数框架是这么实现的：

```go
r = NewRouter()
r.Use(logger)
r.Use(timeout)
r.Use(ratelimit)
r.Add("/", helloHandler)
```

这里的 Use 的源码实现也贴出来一下：

```go
type middleware func(http.Handler) http.Handler

type Router struct {
    middlewareChain [] middleware
    mux map[string] http.Handler
}

func NewRouter() *Router {
    return &Router {
        mux: make(map[string]http.Handler),
    }
}

func (r *Router) Use(m middleware) {
    r.middlewareChain = append(r.middlewareChain, m)
}

func (r *Router) Add(route string, h http.Handler) {
    var mergeHandler = h
    
    for i := len(r.middlewareChain) - 1; i >= 0; i-- {
        mergeHandler = r.middlewareChain[i](mergeHandler)
    }
    
    r.mux[route] = mergeHandler
}
```

这段代码应该也可以直接用来执行。其中的原理是在写 Use 的时候，引用了 middleware 并且直接 append 到 middlewareChain 的数组中。并且把已经在middleware 里面的东西按照洋葱的包装方式，一层层地包装完成，最终把业务逻辑包装完成并且挂载到路由上。



OK，下期文章继续讲解 Router 的实现。
