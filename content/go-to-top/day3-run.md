<<<<<<< HEAD
---
title: "Go 翻山越岭——调试1"
date: 2021-08-09T09:38:51+08:00
category: tech
tags:
    - Go
keywords:
    - Coding
---

今天来动手实践看看，Go 的底层是如何运作、如何进行调度的，通过调试一段简单的代码，我将带你体验 Go 语言如何接触底层的知识。不会很复杂的，跟着我一步步走肯定都能理解。



既然要进行代码实践操作，首先要考虑的就是Go语言编译器（也就是Go）、Go 编程环境以及 Go 代码的调试环境，这里我们需要用到的版本：

> 操作系统：CentOS7 或其他 Linux 环境
>
> Go：1.14.12
>
> Go 调试工具：readelf

是不是有点复杂？如果版本和我的不一样，那么实验结果可能不一样。如果版本比较低，有的工具需要自己去找旧版本的去匹配安装。那有没有什么办法可以简化这些实验环境的搭建呢？用上容器技术——Docker。简单说，当你用 Docker 的时候，就像在自己的电脑中装上了我所说版本的实验环境，里面包括了上面所说的所有东西。你可以用我给出的定制环境（我们称之为镜像），也可以用其他人给出的定制环境，并且可以做到无缝切换。就这么说，在你的 windows 上可以运行 CentOS、Ubuntu，甚至还能在 windows 上运行一个微型 windows，是不是很神奇？OK，我们说回今天的实验，在开始之前希望你能去了解 Docker 是怎样操作的，这里推荐我非常喜爱的 B 站 UP 主“遇见狂神说”推出的 [Docker 教程](https://www.bilibili.com/video/BV1og4y1q7M4?from=search&seid=13499423805328211071)。好，就当你已经会基本操作啦，跟我开始敲代码了：

1. 搭建实验环境

```bash
docker run -it xargin/go1.14.12-dev bash
```

这里我们运行了曹大（曹春晖）提供的实验环境，并且进入到这个容器环境中。但是存在一个问题，我们在这里面写的代码在这个容器销毁后就没有了，所以我们需要对容器做一个映射，在容器中写了文件，在我们本地也有一份记录。

```bash
# 挂载目录，文件也可以生成
docker run -it -v /root/project/dockerProject:/root/project xargin/go1.14.12-dev /bin/bash

```

2. 创建目录与文件

在home目录下创建文件夹所有实验目录 work，创建第一个实验目录 ch。（这里的样式和大家不一样是因为 安装了 zsh ，你不一定也要按照，如果感兴趣，请自行搜索安装，当前配色主题是 ys）

![img](C:\Users\Xfavor\AppData\Local\YNote\data\weixinobU7VjplR8GhsqZn6qPSNHGlLvOg\2a19397ad8e14472b1f4c9358d3a8ed7\clipboard.png)

实验环境准备完毕。

2. 编写一段简单代码

我们来打开一个叫做 hello.go 的文件，写一段我们最熟悉的代码：

```bash
vi hello.go
```

```go
// hello.go 文件中写入
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```

3. 编译代码

生成可执行程序并查看：

```bash
go build hello.go

ls
```

![img](C:\Users\Xfavor\AppData\Local\YNote\data\weixinobU7VjplR8GhsqZn6qPSNHGlLvOg\03dcb6b3af0b4816b55c2babbe182f59\clipboard.png)

这里的绿色 hello，就是我们 Go 语言编译后的可执行程序，通过这么几步操作我们其实是实现了这样的一个流程：

> 文本 → 编译 → 二进制可执行文件

这里大家需要注意，在不同的操作系统上，可执行文件的格式（规范）是不一样的。

> Linux 可执行文件叫做 ELF
>
> Windows 可执行文件叫做 PE
>
> MacOS 可执行文件叫做 Mach-O

今天我们用到的是 Linux，所以这里着重讲一下 Linux 的可执行文件 ELF(Executable and Linkable Format) 

> ELF 由三部分组成：ELF header、Section header、Sections

> Linux 执行 ELF 文件流程：解析 ELF header → 加载 文件内容至内存 → 从 entry point 开始执行代码

4. 尝试使用 readelf 工具

我们一般通过工具 readelf 来找到 entry point address，通过以下命令来实现：

```bash
readelf -h hello
```

![image-20210809235305494](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210809235305494.png)

Ok，通过这么几步简单的流程，我们就找到了进入 Go 底层的入口了，下期再为大家讲解底层有什么东西。

=======
---
title: "Go 翻山越岭——调试1"
date: 2021-08-09T09:38:51+08:00
category: tech
tags:
    - Go
keywords:
    - Coding
---

今天来动手实践看看，Go 的底层是如何运作、如何进行调度的，通过调试一段简单的代码，我将带你体验 Go 语言如何接触底层的知识。不会很复杂的，跟着我一步步走肯定都能理解。



既然要进行代码实践操作，首先要考虑的就是Go语言编译器（也就是Go）、Go 编程环境以及 Go 代码的调试环境，这里我们需要用到的版本：

> 操作系统：CentOS7 或其他 Linux 环境
>
> Go：1.14.12
>
> Go 调试工具：readelf

是不是有点复杂？如果版本和我的不一样，那么实验结果可能不一样。如果版本比较低，有的工具需要自己去找旧版本的去匹配安装。那有没有什么办法可以简化这些实验环境的搭建呢？用上容器技术——Docker。简单说，当你用 Docker 的时候，就像在自己的电脑中装上了我所说版本的实验环境，里面包括了上面所说的所有东西。你可以用我给出的定制环境（我们称之为镜像），也可以用其他人给出的定制环境，并且可以做到无缝切换。就这么说，在你的 windows 上可以运行 CentOS、Ubuntu，甚至还能在 windows 上运行一个微型 windows，是不是很神奇？OK，我们说回今天的实验，在开始之前希望你能去了解 Docker 是怎样操作的，这里推荐我非常喜爱的 B 站 UP 主“遇见狂神说”推出的 [Docker 教程](https://www.bilibili.com/video/BV1og4y1q7M4?from=search&seid=13499423805328211071)。好，就当你已经会基本操作啦，跟我开始敲代码了：

1. 搭建实验环境

```bash
docker run -it xargin/go1.14.12-dev bash
```

这里我们运行了曹大（曹春晖）提供的实验环境，并且进入到这个容器环境中。但是存在一个问题，我们在这里面写的代码在这个容器销毁后就没有了，所以我们需要对容器做一个映射，在容器中写了文件，在我们本地也有一份记录。

```bash
# 挂载目录，文件也可以生成
docker run -it -v /root/project/dockerProject:/root/project xargin/go1.14.12-dev /bin/bash

```

2. 创建目录与文件

在home目录下创建文件夹所有实验目录 work，创建第一个实验目录 ch。（这里的样式和大家不一样是因为 安装了 zsh ，你不一定也要按照，如果感兴趣，请自行搜索安装，当前配色主题是 ys）

![img](C:\Users\Xfavor\AppData\Local\YNote\data\weixinobU7VjplR8GhsqZn6qPSNHGlLvOg\2a19397ad8e14472b1f4c9358d3a8ed7\clipboard.png)

实验环境准备完毕。

2. 编写一段简单代码

我们来打开一个叫做 hello.go 的文件，写一段我们最熟悉的代码：

```bash
vi hello.go
```

```go
// hello.go 文件中写入
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```

3. 编译代码

生成可执行程序并查看：

```bash
go build hello.go

ls
```

![img](C:\Users\Xfavor\AppData\Local\YNote\data\weixinobU7VjplR8GhsqZn6qPSNHGlLvOg\03dcb6b3af0b4816b55c2babbe182f59\clipboard.png)

这里的绿色 hello，就是我们 Go 语言编译后的可执行程序，通过这么几步操作我们其实是实现了这样的一个流程：

> 文本 → 编译 → 二进制可执行文件

这里大家需要注意，在不同的操作系统上，可执行文件的格式（规范）是不一样的。

> Linux 可执行文件叫做 ELF
>
> Windows 可执行文件叫做 PE
>
> MacOS 可执行文件叫做 Mach-O

今天我们用到的是 Linux，所以这里着重讲一下 Linux 的可执行文件 ELF(Executable and Linkable Format) 

> ELF 由三部分组成：ELF header、Section header、Sections

> Linux 执行 ELF 文件流程：解析 ELF header → 加载 文件内容至内存 → 从 entry point 开始执行代码

4. 尝试使用 readelf 工具

我们一般通过工具 readelf 来找到 entry point address，通过以下命令来实现：

```bash
readelf -h hello
```

![image-20210809235305494](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210809235305494.png)

Ok，通过这么几步简单的流程，我们就找到了进入 Go 底层的入口了，下期再为大家讲解底层有什么东西。

>>>>>>> 0cae884533b7a5e0cf3eb8be746627712e680e3e
