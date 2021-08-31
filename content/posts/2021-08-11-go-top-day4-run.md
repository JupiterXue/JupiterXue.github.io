<<<<<<< HEAD
---
title: "Go 翻山越岭——调试2"
date: 2021-08-11T09:38:51+08:00
category: tech
tags:
    - Go
keywords:
    - 实践
---

上回我们说到，利用调试工具 readelf 读取我们在 Linux 上编译的可执行文件 hello，并且找到了查看 Go 底层的入口 entry point address，今天我们尝试踏进底层调试的大门，看一看我们平时所说的底层，是否真的很困难，今天的主题是“初探Go底层”

[toc]

# 实验环境与回顾

> 操作系统：CentOS 7
>
> Go：1.14.12
>
> 调试工具：readelf、dlv

这里我们同样运行曹大提供的实验环境，并且挂载映射到本地项目，然后进入到这个容器环境中。

```bash
# 挂载目录，文件也可以生成
docker run -it -v /root/project/dockerProject:/root/project xargin/go1.14.12-dev /bin/bash

# 进入项目目录
cd /project/work/ch01
# 查看文件
ls

```

![image-20210811230644860](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811230644860.png)

这里的绿色 hello，就是我们 Go 语言编译后产生的 ELF 格式可执行程序（ELF 不知道是什么，可回顾上期文章）

我们一般通过工具 readelf 来找到进入可执行程序底层的入口 entry point address，通过以下命令来实现，这里的 -h 表示去读取 header

> -h --file-header	Display the ELF file header

```bash
readelf -h hello
```

![image-20210809235305494](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210809235305494.png)

# 开始底层调试

当找我们通过 readelf 找到了调试入口的十六进制地址码，我们还需要用另一个工具 dlv 才能来调试 Golang 程序。再次之前，来简单说明一下什么是 readelf、什么是 dlv。



我们知道，Windows 上可执行文件的格式为 .exe，而 Linux 上可执行文件的格式为 ELF，见名知意，那么 **readelf 是用于查看 ELF 可执行文件的命令工具**，并且一般会有单个或多个参数的组合命令，这里我们只关注 -h 参数命令。



Delv，简称 **dlv 是 Go 语言的第三方调试工具**，其目标是为 Go 提供简单、功能齐全的调试工具。相比于 Go 自带的调试器，dlv 做到了观察代码底层行为的功能。并且也是通过单个或多个命令组合实现一定功能。



好，开始正式调试旅程，跟着我左手右手：

```bash
# 开启对可执行 go 文件的调试
dlv exec ./hello

# 对 go 可执行文件的入口打断点，这里用到我们通过 readelf 工具找到的地址
b *0x455780
```

![image-20210811233208118](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811233208118.png)

可以看到整个可执行程序代码位置了



继续查看代码的运行，通过 c 命令，c 是 continue 英文的简写

![image-20210811233238802](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811233238802.png)



再用 si 命令，执行单步调试。（如果不知道 si 命令是什么意思，可以通过 h 命令，查看帮助，可以找到以下这段话）

> stel-instrution (alias: si)	Single step a single cpu instruction.

这里我们先记住一个 l 命令，可以查看当前执行到哪个位置。

![image-20210811233604465](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811233604465.png)

再用 si 进入下一步

![image-20210811233644567](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811233644567.png)

然后一路输入 si、si、si 可以看到每一步 Go 的底层做了什么，即 Go 语言在汇编层面进行了什么样的指令集。



OK，通过以上的流程，我们搞清楚了 “Go 进程的启动与初始化的流程”。首先通过 entry point 找到 Go 进程的执行入口，然后实现了：文本 → 编译 → 二进制可执行文件

![未命名文件](C:\Users\Xfavor\Desktop\未命名文件.png)

最后，我们可以通过 si 来一步步查看 Go 底层执行的代码，不过这种一直 si 的方式是不是有点不爽？没关系，有好的解决办法——上 Goland。在 Goland 中我们通过 Command+Shift+9，打开底层搜索框，输入 ·rt0_go ，然后点击系统的 rt0_go 的所在位置：

![image-20210811235217024](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811235217024.png)

在 Goland 中，就可以自在地上下滚动去查看汇编代码了，不用再机械地用 si 进行调试了。



这些汇编指令看不懂没有关系，直接找里面的关键函数看就好了。比如不用关注初始化、处理平台的完全不用看，而去关注 ok 后面的代码就可以

![image-20210811235344792](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811235344792.png)



今天初步探索了 Go 底层的调试，是不是感觉也没有那么困难。没错，只要肯花时间，都能看得懂的。



__参考资料__

[1] readelf命令使用说明

https://blog.csdn.net/yfldyxl/article/details/81566279

[2] 译文 使用 Delve 工具调试 Golang 程序, GOCN

https://gocn.vip/topics/12090

[3] 常用的分析ELF文件的命令（readelf、objdump及od）

https://blog.csdn.net/qq_21331015/article/details/103210449
=======
---
title: "Go 翻山越岭——调试2"
date: 2021-08-11T09:38:51+08:00
category: tech
tags:
    - Go
keywords:
    - 实践
---

上回我们说到，利用调试工具 readelf 读取我们在 Linux 上编译的可执行文件 hello，并且找到了查看 Go 底层的入口 entry point address，今天我们尝试踏进底层调试的大门，看一看我们平时所说的底层，是否真的很困难，今天的主题是“初探Go底层”

[toc]

# 实验环境与回顾

> 操作系统：CentOS 7
>
> Go：1.14.12
>
> 调试工具：readelf、dlv

这里我们同样运行曹大提供的实验环境，并且挂载映射到本地项目，然后进入到这个容器环境中。

```bash
# 挂载目录，文件也可以生成
docker run -it -v /root/project/dockerProject:/root/project xargin/go1.14.12-dev /bin/bash

# 进入项目目录
cd /project/work/ch01
# 查看文件
ls

```

![image-20210811230644860](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811230644860.png)

这里的绿色 hello，就是我们 Go 语言编译后产生的 ELF 格式可执行程序（ELF 不知道是什么，可回顾上期文章）

我们一般通过工具 readelf 来找到进入可执行程序底层的入口 entry point address，通过以下命令来实现，这里的 -h 表示去读取 header

> -h --file-header	Display the ELF file header

```bash
readelf -h hello
```

![image-20210809235305494](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210809235305494.png)

# 开始底层调试

当找我们通过 readelf 找到了调试入口的十六进制地址码，我们还需要用另一个工具 dlv 才能来调试 Golang 程序。再次之前，来简单说明一下什么是 readelf、什么是 dlv。



我们知道，Windows 上可执行文件的格式为 .exe，而 Linux 上可执行文件的格式为 ELF，见名知意，那么 **readelf 是用于查看 ELF 可执行文件的命令工具**，并且一般会有单个或多个参数的组合命令，这里我们只关注 -h 参数命令。



Delv，简称 **dlv 是 Go 语言的第三方调试工具**，其目标是为 Go 提供简单、功能齐全的调试工具。相比于 Go 自带的调试器，dlv 做到了观察代码底层行为的功能。并且也是通过单个或多个命令组合实现一定功能。



好，开始正式调试旅程，跟着我左手右手：

```bash
# 开启对可执行 go 文件的调试
dlv exec ./hello

# 对 go 可执行文件的入口打断点，这里用到我们通过 readelf 工具找到的地址
b *0x455780
```

![image-20210811233208118](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811233208118.png)

可以看到整个可执行程序代码位置了



继续查看代码的运行，通过 c 命令，c 是 continue 英文的简写

![image-20210811233238802](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811233238802.png)



再用 si 命令，执行单步调试。（如果不知道 si 命令是什么意思，可以通过 h 命令，查看帮助，可以找到以下这段话）

> stel-instrution (alias: si)	Single step a single cpu instruction.

这里我们先记住一个 l 命令，可以查看当前执行到哪个位置。

![image-20210811233604465](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811233604465.png)

再用 si 进入下一步

![image-20210811233644567](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811233644567.png)

然后一路输入 si、si、si 可以看到每一步 Go 的底层做了什么，即 Go 语言在汇编层面进行了什么样的指令集。



OK，通过以上的流程，我们搞清楚了 “Go 进程的启动与初始化的流程”。首先通过 entry point 找到 Go 进程的执行入口，然后实现了：文本 → 编译 → 二进制可执行文件

![未命名文件](C:\Users\Xfavor\Desktop\未命名文件.png)

最后，我们可以通过 si 来一步步查看 Go 底层执行的代码，不过这种一直 si 的方式是不是有点不爽？没关系，有好的解决办法——上 Goland。在 Goland 中我们通过 Command+Shift+9，打开底层搜索框，输入 ·rt0_go ，然后点击系统的 rt0_go 的所在位置：

![image-20210811235217024](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811235217024.png)

在 Goland 中，就可以自在地上下滚动去查看汇编代码了，不用再机械地用 si 进行调试了。



这些汇编指令看不懂没有关系，直接找里面的关键函数看就好了。比如不用关注初始化、处理平台的完全不用看，而去关注 ok 后面的代码就可以

![image-20210811235344792](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210811235344792.png)



今天初步探索了 Go 底层的调试，是不是感觉也没有那么困难。没错，只要肯花时间，都能看得懂的。



__参考资料__

[1] readelf命令使用说明

https://blog.csdn.net/yfldyxl/article/details/81566279

[2] 译文 使用 Delve 工具调试 Golang 程序, GOCN

https://gocn.vip/topics/12090

[3] 常用的分析ELF文件的命令（readelf、objdump及od）

https://blog.csdn.net/qq_21331015/article/details/103210449
>>>>>>> 0cae884533b7a5e0cf3eb8be746627712e680e3e
