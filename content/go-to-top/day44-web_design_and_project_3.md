---
title: "Go 翻山越岭——web 框架原理（3）"
date: 2021-12-28T22:15:02+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

没想到吧，今天继续写点技术文章。上期文章我们说到，路由 Router 的本质其实是从字符串匹配到用户函数的过程；相似字符串对应不同的路由，中间经过 Radix Tree 构造字典树进而实现。这期文章继续讲讲常见 Web 框架中的组件。

# Validator

在没有拦截器（Validator）之前，可能我们写的一个关于，注册请求的代码是这样的：

```go
type RegisterReq struct {
    Username             string    `json:"username"`
    PasswordNew      string    `josn:"password_new"`
    PasswordRepeat  string    `json:"password_repeat"`
    Email                   string   `json:email`
}

func register(req RegisterReq) error {
    if len(req.Username) > 0 {
        if len(req.PasswordNew) > 0 && len(req.PasswordRepeat) > 0 {
            if req.PasswordNew == req.PassowrdRepeat {
                createUser()
                return nil
            } else {
                return errors.New("invalid email")
            }
        } else {
            return errors.New("password and password reinput must be loger than 0")
        }
    } else {
        return errors.New("length of username cannot be 0")
    }
}
```

这段代码的意思是说，我们想要创建一个账户。

- 首先，用户名的长度需要大于 0；
- 其次，密码的长度大于 0，并且重复确认输入的密码的长度也大于 0；
- 然后，密码和重复确认的密码必须相同；
- 再然后，注册用到的 email 必须是合法的。

这段代码看起来就像开了花，但其实是很多公司的日常开发写的。但我们如果这种代码写多了不仅感到繁琐，看起来更是感到不舒服。



于是我们会去学习一些重构的理念之后，比如 early return / guard cluuse 的重构思想后，重构改造写出的代码可能就是这样的：

```go
func register(req RegisterReq) error {
    if len(req.Username) == 0 {
        return errors.New("length of username cannot be 0")
    }
    
    if len(req.PassowrdNew) == 0 || len(reqPasswordRepeat) == 0 {
        return errors.New("password and password reinput must be loger than 0")
    }
    
    if req.PassowordNEw != req.PassowrdRepeat {
        return errors.New("passowrd and reinput must be the same")
    }
    
    if emailFormatValid(req.Email) {
        return errors.New("invalid email")
    }
    
    createUser()
    return nil
}
```

重构之后的代码看起来相比之前的代码是更加清晰、有条理些。



但其实这种代码还能通过 validator 进行改造优化，像这种字段校验的脏活、累活，就尽量不要自己手动去写，可以用一个包来实现，在 go-playground 包中，有一个 validator，代码如下：

```go
import "github.com/go-playground/validator/v10"

type RegisterReq struct {
    // 字符串的 gt=0 表示长度必须 > 0, gt means greater than
    Username            string    `validate:"gt=0"`
    // 同上
    PasswordNew     string    `validate:"gt=0"`
    // eqfield 跨字段相等校验
    PasswordRepeat  string    `validate:"eqfield=PasswordNew"`
    // 合法 email 格式校验
    Email                   string    `validate:"email"`
}

var validate = validator.New()

func validateFunc(req RegisterReq) error {
    err := validate.Struct(req)
    if err != nil {
        domSomething()
        return err
    }
    ......
}
```

- 我们前面提到的校验规则都可以写在当前的 RegisterRequest 结构中某个字段 的 tag 里。比如我们要求 Username 的长度必须大于 0，就是 gt=0，英文意思是 greater than 0。

- 我们可以在字段和字段之间去做校验，比如我们输入的密码和重复输入的密码必须相等；比如我们输入的 email 必须是内置规则下的合法的 email。

- 下方的函数 validateFunc 就是在使用全局定义的 validate。将 RegisterRequest 对象进行 validate.Struct 校验匹配，然后依次校验每个字段是否符合规则，如果不符合规则就返回错误信息，这些错误信息其实我们也可以去做定制的。



像这种 validator 的应用是比较广泛的，我们常见的 gin 框架里面，其实是把这个组件集成在了一起。虽然平常不怎么看得到，但底层就是这么个设计。



如果不是作为外部领域的应用，我们还可以将validator 单独拎出来使用。比如防御性的编程之类的。

## Validator 基本原理

我们除了知道怎么去运用 validator 去做编程开发，还可以进一步了解它的基本原理。



在此之前，我们需要意识到，我们日常工作中定义的 struct 可以这样理解：

```go
type Nested struct {
    Email    string `validaate:"email"`
}

type T struct {
    Age    int `validaate:"eq=10"
    Nested Nested`
}
```

即，我们可以把任意多个 struct 理解成一棵树 。



在 type T 结构中，有一个 Nested 字段，它对应着一个内嵌的结构体 Nested，Nested 结构体中有一个 Email 字段。然后我们将整个结构体的树状图画出来，就是这样的结构：

![validator 基本原理-树-深度优先遍历 (1)](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202112282321720.png)

<center>validator 基本原理-树结构例子</center>

- struct T 是一个根节点；
- 左边是 Age 字段；
- 右边是内嵌的结构体 Nested；
- 内嵌结构体有着自己的 Email 字段。

如果我们想要完成类似 validator 的功能，就需要对这棵树做深度优先遍历，或者广度优先遍历。



既然我们可以通过深度优先遍历的方式，遍历到起始到内嵌结构体所有的字段。那么我们就可以利用反射的 API 做到这件事，可以访问到字段并且可以拿到对应字段的注释和注释中的 tag。这句话反过来说也成立，既然可以拿到所有字段，那么通过深度优先遍历的方式可以扫描到所有字段。



__参考资料__

[1] 《Go 高级编程》5.4 validator请求校验

https://github.com/chai2010/advanced-go-programming-book/blob/master/ch5-web/ch5-04-validator.md#543-%E5%8E%9F%E7%90%86



OK，下期文章继续讲解请求绑定的组件 request binder 的基本原理和实现！
