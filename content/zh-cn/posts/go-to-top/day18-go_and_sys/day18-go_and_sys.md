---
title: "GO 翻山越岭——Go 常见系统调用"
date: 2021-09-07T07:28:31+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

系统调用，是操作系统内核为应用提供的 API。今天继续来讲一个系统调用案例和 Go 中常见系统调用



| arch | syscall NR | return | arg0 | arg1 | arg2 | arg3 | arg4 | arg5 |
| :----: | :----: | :--: | :----: | :----: | :----: | :----: | :----: | :----: |
| 213 | epoll_create | man/ cs/ | 0xd5 | int size | - |-|-|-|

来源：https://chromium.googlesource.com/chromiumos/docs/+/master/constants/syscalls.md#x86_64-64_bit



这是 Linux 中的系统调用，编号是 213，我们可以从一下代码了解更多信息：

```ABAP
# define SYS_epoll_create    213

TEXT runtime.epoll_create(SB), NOSPLIT, $0
	MOVL	size+0(FP), DI
	MOVL	$SYS_epoll_create, AX
	SYSCALL
	MOVL	AX, ret+8(FP)
	RET
```

编号 213 按照调用规约会被存储在 rax 寄存器中，也就是这里的 AX 寄存器。epoll_create 只有一个阐述传递，也就是 int 类型的 size。SYSCALL 直接进入内核去了。



因此，我们需要知道 SYSCALL 之后发生了什么，即我们写好应用程序 Application program 执行了哪些 SYSCALL 逻辑。从下图可以得知：

![SYSCALL 之后发生了什么](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109072232527.png)

当我们调用了 syscall.EpollCreate 之后进入了 syscall 包中的一段代码，然后根据一定的规则翻译成 runtime.RawSyscall，里面会有一堆准备参数和逻辑，完了之后就是执行 SYSCALL 指令。上期说的我们写的代码处于 CPU 分级保护域 ring-3，执行这个指令之后就能帮助我们从 ring-3 切换到 ring-0，然后去执行内核相关的代码了。ring-0 内核模式下什么操作都能够进行，没有特别的权限限制。内核也有相关的系统调用表（在最上面的链接中），会去调用具体的系统调用实现。

> 推荐一个可以查看层内核的在线工具，做 C/C++ 的同学都可以看：code.woboq.org



# Go 常见系统调用

我想，大多数同学和我一样，想问一个问题：了解了系统调用对我们有什么用呢？按理来说，我们还要了解几个部分，Go 的实现，Go 怎么和操作系统交互。知道了 GC 就知道 Go 怎么和系统进行交互。知道了内存相关就知道内核 flag 或 syscall 被修改后为什么会导致应用层行为出这么大的问题。像这些业务问题我们大概知道一点，就不会在做应用的时候完全束手无策。



在 Go 代码中，也有支持系统调用的地方，比如：

| Go 代码                     | 系统调用                                     |
| --------------------------- | -------------------------------------------- |
| os.GetPid()                 | getpid()                                     |
| println("hello world")      | write(2, "hello world", 11)                  |
| startm -> newm -> newosproc | clone(child_stack=0xc43003c00, flags=......) |

系统调用 write(2, "hello world", 11) 中，2 代表 stderr，11 代表字符串长度 strlen。我们在 Linux 中最常见的就是 stdin、stdout、stderr，如果不主动做修改，他们分别对应数字 0, 1, 2. 这里可能会问到，为什么会传递字符串长度，在 C 语言中传递字符串就需要传指针，因为指针不知道字符串的长度。而 Go 语言虽然看起来只是传递了数据结构，但底层还是会做展开。比如传入了一个数组，就会展开为：一个地址、一个长度 len、一个容量 cup，写的是一个参数，底层就展开为了三个参数。



启动线程代码 startm -> newm -> newosproc  中，newosproc  最终调用了系统调用 clone。



然后再看下不同操作系统下常见的系统调用

|  Types of System Calls  | Windows                                                      | Linux                                        |
| :---------------------: | ------------------------------------------------------------ | -------------------------------------------- |
|     Process Control     | CreateProcess()<br />ExitProcess()<br />WaitForSingleObject() | fork()<br />exit()<br />wait()               |
|     File Management     | CreateFile()<br />ReadFile()<br />WriteFile()<br />CloseHandle() | open()<br />read()<br />write()<br />close() |
|    Device Management    | SetConsoleMode()<br />ReadConsole()<br />WriteConsole()      | ioctl()<br />read()<br />write()             |
| Information Maintenance | GetCurrentProcessID()<br />SetTimer()<br />Sleep()           | getpid()<br />alarm()<br />sleep()           |
|      Communication      | CreatePipe()<br />CreateFileMapping()<br />MapViewOfFile()   | pipe()<br />shmget()<br />mmap()             |

可以发现，不同操作系统下系统调用是不一样的，所以在不同操作系统上看到的也不一样。了解这些的意义也是在于，我们想要做超过自身能力意外的事情就必须求助于操作系统。



总结，了解系统调用后可以发现这个一个现象：我们的应用程序是被关在权限的监狱中，只能做很少的一些事情，其他的必须拜托操作系统来做。



下期预告：利用工具去观察系统调用。

