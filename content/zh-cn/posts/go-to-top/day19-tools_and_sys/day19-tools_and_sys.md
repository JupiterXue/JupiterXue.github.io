---
title: "GO 翻山越岭——系统调用调试工具"
date: 2021-09-09T22:43:41+08:00
category: tech
tags:
    - Go
keywords:
    - 实践
---

**了解系统调用就能在业务发生问题时有更多招数，见招拆招**。知道了系统调用的概念，有时候我们还需要用工具去实际观察进程是如何发起系统调用，对操作系统做了哪些请求，今天就来聊一聊。

# 观察系统调用工具

在 Linux 中我们常用用观察系统调用的工具是 strace，在 macOS 系统上常用的是 dtruss



案例1，通过 strace 可以看到一个 Go 进程的启动过程到底调用了哪些系统调用，这里有个例子（具体代码不用关注）：



![image-20210909230308944](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109092303708.png)

<center>系统调用示例</center>

使用 strace 指令 + 可执行文件x，看到了有哪些系统调用函数被调用。execve 执行哪个二进制文件，arch_prctl、sched_getaffinity 设计 CPU 亲和度，oepnat 加载相关文件，mmap 系统映射，gettid 获取线程 ID。



案例2，通过 strace 还可以查看一些软件的做了哪些系统调用，比如查看 nginx 的：



![image-20210909231712937](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109092317098.png)

<center>查看nginx 系统调用</center>

我们知道 nginx 在平时不服务的时候都是阻塞的状态，即阻塞在某个系统调用上。使用 strace 指令 + -f 参数 + nginx 可以观察它启动的过程。如果这个软件需要创建多个进程，就需要这里的 -f 的 flag 参数，跟踪所有刚启动线程创建的其他进程。最后我们可以发现 nginx 阻塞在 pid 为 224 的 epoll_wait 上。

注意：在 docker 中用这个命令可能会遇到一些问题，可以参考这个链接《Why strace doesn't work in Docker》https://jvns.ca/blog/2020/04/29/why-strace-doesnt-work-in-docker/



案例3，通过 strace 观察一个 Go 语言 Hello world 程序生命周期中系统调用情况。



![image-20210909232448820](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109092324898.png)



使用 -c 的 flag 参数，可以看到有哪些系统调用。



在往期文章中讲解“Go 内置数据结构-Timer”时，我们在对 Go 1.13 和 Go1.14 进行对比时用到的就是 strace。



如果官方说对 Go 的 runtime 中做了许多优化，比如 time.Sleep 中用到的 syscall 相比以前少了很多，那么我们同样可以用这个工具去查看。通过运行一定次数 time.Sleep 看前后的结果是否真的有什么区别。



需要注意，strace 的实现其实也依赖了 ptrace 这个 syscall，所以本质上我们要去追踪一个 syscall 还需要依赖一个 syscall（有点套娃）。另外，调试器（如 delve）也是大量使用了 ptrace。因此 **ptrace 不仅可以做跟踪工具，还能做调试工具**。

# Go 语言中的系统调用

了解了这么多基础概念，还是要回到我们的 Go 语言，看看它有哪些内部系统调用。



首先我们要区分一下，在系统调用中有一个简单分类：有一些是阻塞的，有一些是非阻塞的。



![未命名文件](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109092340838.png)

<center>syscall 分类</center>

假如我们对一个网络 fd 去做 read 操作，而连接中没有数据，就会阻塞在这里。其实 read 还是有费阻塞模式的，所以也不太好区分。read 操作的底层通过调用 io_control 之类的系统调用，变成了非阻塞的的调用方式。



在我们用户状态去做阻塞应用的话，其实底层都会转化为非阻塞调用。



还有一部分系统调用没办法做成非阻塞调用。这种情况一定会阻塞。在往期文章讲“调度”的时候，对于 runtime 的 Go 来说没办法去做拦截的一类阻塞。



这两种类型的系统调用最终在底层也会被翻译成不同的函数调用，在底层的函数调用中还是有一些细节区分的，像后面有数字 6 的就是一种派生函数。



问题来了，为啥还要有个 6 后缀？这里其实也是遵循了 Linux 编程中的常见规范。



![image-20210909234809637](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109092348692.png)

<center>syscall 带数字</center>

很多 Linux 函数都是带有数字的，6 其实说的是 6 个参数。很多系统接口也有类似的命名发，如 wait4，accept4。



这里列举一下和 Syscall 相关代码所在的具体位置：

1. OS 相关的基础⽂件，在 syscall package 中：https://golang.org/src/syscall/syscall_linux.go
2. 使⽤脚本⽣成的⽂件，在 syscall package 中：https://golang.org/src/syscall/zsyscall_linux_386.go
3. 不对⽤户暴露的特殊 syscall，不受调度影响，在 runtime 中：https://golang.org/src/runtime/sys_linux_amd64.s



我们知道了系统调用分为阻塞和非阻塞的。还需要注意的是：阻塞的系统调⽤，有特殊的逻辑去把 P 和 M 解绑定，即修改 P 的状态：running -> syscall。这样在 sysmon 中才能发现这个 P 已经在 syscall 状态阻塞了。不过解绑定也是有前提的，如图红色两行代码：



![image-20210909235424347](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109092354413.png)

<center>P 和 M 解绑定约束</center>

# 系统调用发展历史

最后，系统调用的发展历史，可以用一张图简单说明。



![image-20210909235643954](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109092356006.png)

<center>系统调用发展历史简要说明</center>



__参考资料__

[1] Linux strace命令, 博客园

https://www.cnblogs.com/ggjucheng/archive/2012/01/08/2316692.html

[2] [译]strace的10个命令, 鸟窝

https://colobu.com/2021/04/30/strace-commands-for-troubleshooting-and-debugging-linux/
