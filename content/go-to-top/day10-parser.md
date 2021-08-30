---
title: "Go 翻山越岭——解释器"
date: 2021-08-22T22:51:00+08:00
category: tech
tags:
    - Go
keywords:
    - 实践
---

上期文章遗留了一个问题“三个 panic 在 runtime 代码的具体哪个位置输出”，如果通过汇编调试工具找出来。今天来动手实践，并且尝试解决另一个业务问题。

# 问题回顾

> 分别有三种情况的 panic：
>
> 1. 往已经关闭的 channel中写入数据
> 2. 关闭一个是 nil 值的 channel
> 3. 关闭一个已经是关闭状态的 channel

要找出它们在 runtime 代码中具体位置，首先我们先写几个小型代码，来完全模拟这几种情况。



情况1：往已经关闭的 channel中写入数据，文件名  send_to_close.go

```go
package main

func main() {
    var ch chan int
    close(ch)
    ch <- 1
}
```

```bash
go tool compile -S send_to_close.go | grep "send_to_close.go:6"
```

![image-20210823233649497](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108240001413.png)

通过反汇编和文本搜索，找到了关键信息 **runtime.chansend1(SB)**



情况2：关闭一个是 nil 值的 channel，文件名  close_nil.go

```go	
package main

func main() {
    var ch chan int
    ch = nil
    close(ch)
}
```

```bash
go tool compile -S close_nil.go | grep "close_nil.go:6"
```

![image-20210823234414243](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108240001274.png)



情况3：关闭一个已经是关闭状态的 channel，文件名  close_closed.go

```go
package main

func main() {
    var ch chan int
    close(ch)
    close(ch)
}
```

```bash
go tool compile -S close_closed.go | grep "close_closed.go:6"
```

![image-20210823233943750](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108240000241.png)

这里找到的关键信息是 **runtime.closechan(SB)**



其实我们还有一种投机的方法来找到以上三个情况。假设我们经过多年的调试，已经知道 channel 的实现一定是在 runtime下 chan.go 这个文件里。在 Goland 中双击 shift 点击 Files 文件搜索 "chan.go"，进入这个文件后 ctrl+f 搜索关键字 "panic" 就找到了。如图所示：

![image-20210823234924519](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108240000237.png)

并且再根据函数名找找关闭 channel 的函数。



其实还有第三种方法，在反汇编工具 dlv 中，panic 默认就是一个断点。



OK，GO 语言语法分析实现已经帮我们解决了业务3的问题，我们继续往前走。

# Parser 场景

> 场景4，有一个中台服务，划分了许多用户等级，具体规则如下：（以下场景都没遇到过，不过可以来做分析）
>
> 初级会员，发帖数 > 10
>
> 中级会员，充值 > 1000 RMB
>
> 高级会员，发帖数 > 100 ，充值 > 10000 RMB
>
> 如果其中项目数量 = 几百，那么每个项目都应该有自己的会员规则吗？也就是说，我们需要单独定制这么几百个代码吗？

在这个业务场景中，我们会用到 Parser 解析器，这个 Parser 其实是编译器前端的东西，它主要的职责是把文本最终变成类似 AST 抽象语法树的东西。



假设数据：

1. 用户信息 posts=143，invest=2042
2. 规则：post>100 && invest > 10000    invest > 1000

![内置 AST 工具 - 简单的规则引擎](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108242225467.png)

Go 语言中有个内置的 Parser 能够帮助我们做些工作：简单的规则引擎。这里可以简单理解为最上方的一句布尔表达式，以语法树结构展示。下方图示暂时不做详解，后续有机会再来详细分析。这里我们先来看另外一个 Parser 场景，能够更好地理解，为什么要使用这种 AST 抽象语法树的结构。



如果大家用过 Beego ，代码中这段类似的注释 comment 一定不陌生。

![image-20210823235449303](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108240000484.png)

这段代码中会从 Go 的注释到 swagger 的 json 文件。其实不光是这段代码，我们只要能够读到这些 comment，就可以转换成为符合 swagger 规范的 json。

![image-20210823235627367](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108240000213.png)



那么问题来了，什么是 swagger？swagger 是开源界做得比较好的一个规范和完整的框架，用于生成、描述、调用和可视化 RESTful 风格的 Web 服务。 即 **API 文档生成工具**。



我们不了解 swagger 可能不知道注释是怎么生成代码的，但我们懂一点点 Go 语言内置的工具，那么我们自己来做从注释到 json 的工具也就《简单》了。



代码不多，来跟我动手实践吧，主文件，文件名 main.go

```go
package main

// this is the comment for main function
func main() {
}
```

Parser 工具，文件名 output.go

```go
package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "./main.go", nil, parser.Mode(4))

	for _, d := range f.Comments {
		ast.Print(fset, d)
	}
}
```

执行这段代码，我们也能够实现一小段注释到 json 的功能：

```bash
go run output.go
```

![image-20210824224559554](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108242246501.png)

第一段代码不用多说，第二段代码用到了 Go 语言自带的工具 ast、parser 和 token。第 10 行定义一个文件符号集合，11 行解析单个Go源文件的源代码并返回对应的ast.File节点。14 行打印从AST节点x到w开始的（子）树。最后我们得到了ast 格式的 json 数据结构。

# 总结

今天我们学到了，**三种方法找到 channel 中的 panic**：

1. 反汇编调试工具 go tool compile 或者 go tool objdump
2. 根据常识，在 Goland 中全路径搜索 runtime 下的 chan.go 文件并搜索关键词 "panic"。
3. Go 语言调试工具 dlv，直接运行，自动在 panic 处打断点停下。



还学到了 **什么是 Parser ，以及一段小 Parser demo 的应用**。



下期预告：讲解更加复杂的 Parser 应用场景。

