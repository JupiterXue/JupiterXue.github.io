---
title: "Go 翻山越岭——web 框架原理（4）"
date: 2021-12-28T22:15:36+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

这周尝试连续更技术文。上期说到，Validator 组件能够将复杂的字段校验工作完全自动化完成，只需使用到 go-playground 中的 validator 包就能实现这一功能。今天继续聊聊其他组件。

# Request Binder

这是一段摘自 gin 框架的代码：

```go
// Content-Type MIME of the most common data formats.
const (
    MIMEJSON	  		    = "application/json"
    MIMEHTML			= "text/html"
    MIMEXML				= "application/xml"
    MIMEXML2			 = "text/xml"
    MIMEPlain			       = "text/plain"
    MIMEPOSTFORM		  =  "application/x-www-form-urlencoded"
    MIMEMutipartPOSTForm	= "multipart/form-data"
    MIMEPROTOBUF		  = "application/x-protobuf"
    MIMEMSGPACK			  =  "application/x-msgpack"
    MIMEMSGPACK2		  =  "application/msgpack"
    MIMEYAML			  = "application/x-yaml"
)
```

由于 gin 走的是 HTTP 协议，所以 Request binding 请求绑定，本质上是根据 HTTP header 中的 Contet-Type 的各种类型情况。



根据这个 header 我们可以判断用户传的 body 和 url 是什么东西。  



下面就是我们在代码中，根据 http 的 contentType 做的简单的 switch-case：

```go
func Default(method, contentType string) Binding {
    if method == http.MethodGet {
        return Form
    }
    
    switch contentType {
        case MIMEJSON:
            return JSON
        case MIMEXML, MIMEXML2:
            return XML
        case MIMEPROTOBUF:
            return ProtoBuf
        case MIMEMSGPACK, MIMEMSGPACK2:
            return MsgPack
        case MIMEYAML:
            return YAML
        case MIMEMultipartPOSTForm:
            return FormMultipart
        default:    // case MIMEPOSTForm:
            return Form
    }
}
```

最后，switch-case 返回的是一个 binding 结构。这个过程其实是设计模式中简单工厂模式的应用。



而不同的工厂实现很简单，其实就是某种 codec 的 unmarshal。

![image-20211229234220577](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112292342018.png)

<center>gin 中 Binding 位置</center>

# SQL Binder 

我们在刚接触 Go 语言的时候，一般都是通过标准库的 database.sql 和 mysql 的驱动去写一些 db 相关的代码。

```go
func main() {
    // db 是一个 sql.DB 类型的对象
    // 该对象线程安全，且内部已包含了一个连接池
    // 连接池的选项可以在 sql.DB 的方法中设置，这里为了简单省略了
    db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/hello")
    if err != nil {
        log.Fatil(err)
    }
    defer db.Close()
    
    var (
        id int
        name string
    )
    rows, err := db.Query("select id, name from users where id = ?", 1)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    // 必须把 rows 里面的内容读完，或者显示调用 Close() 方法，
    // 否则在 defer 的 rows.Close() 执行之前，连接永远不会释放
    for rows.Next() {
        err := rows.Scan(&id, &name)
        if err != nil {
            log.Fatal(err)
        }
        log.Println(id, name)
    }
}
```

- 标准库的 API 难用且容易犯错。比如 db 是可以不用关闭的。
- 有无数的新手 Gopher 倒在 sql.Rows 忘记关闭的坑下。

这里我们也可以看出来，Google 工程师平时不怎么写和业务打交道的代码的，这种 API 的设计明显是有问题的。



这就是为什么我们不用标准的 API 去写代码的原因。因此我们在公司里大多数开发工作都需要在这个基础上包一层，不可能直接去用这个。哪怕是对 rows.Close 的操作简单封装一层，也能防止业务开发人员在做开发的时候忘了这件事，这种事故太容易发生了。



还有一个比较麻烦的地方是，如果我们想要将 rows 和我们的结构体绑定也是比较麻烦，必须一行一行地用 rows.Next 遍历读取，再将所需的字段扫描出来，比如数组、字符串、结构体、map 绑定。



因此，我们在写业务代码的时候都会用一些社区维护好的库去做开发，比如 sqlx，以下摘自它的 README：

![image-20211229234313279](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112292343369.png)

<center>sqlx 代码摘选</center>

这段代码是和事物相关，看起来略微复杂，但普通的查询应该比这个更简单些。它给用户提供的 SQL 有以下几个功能，不过总的来说，实现很简单，核心只有一句话：占位符替换。



如果我们想要自己写代码来实现这种类似功能，就是个字符串遍历，只要能把特殊开头的单词提取出来就可以了，⽐如：开头，@开头，$开头。

# go-micro

最后来聊一聊目前最流行的微服务框架 go-micro 实现的插件原理，以下摘自 go 夜读：

![image-20211229234917298](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112292349393.png)

<center>go-micro 插件化原理</center>

它的原理是在每一层定义了都定义了一个接口，也就是说为每个组件强定义了接口。



如果我们亲自去看了它的源码，会发现这个项目写得非常规范，是个非常好的开源项目学习范本。



还需要注意一点，虽然开源框架的设计比较好，但整体的结构还是略微复杂的。我们在学习开源框架的时候不必完整地梳理流程，看每一层的抽象就好了；而如果真正用到这个框架的时候才需要深入去看。



OK，到此为止，常见的 web 框架原理就讲解完了，你学会了吗。
