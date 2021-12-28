---
title: "读论文《Escape from Escape Analysis of Golang》"
date: 2021-10-25T23:20:11+08:00
category: tech
tags:
    - Go
keywords:
    - 理论
---

首先感谢谢大，推荐了这篇技术论文，我也很乐意读下去，因为一看到论文的标题 “Escape xxx Escape ...”，就立刻联想到了我最喜爱的乐队——逃跑计划。所以这篇论文也是很和我的“胃口”。



最近在练习 Leetcode 的过程中，会去阅读攻略教程，并且边读边看，边打开 process-on 思维导图-大纲模式去记笔记，也越来越能接受用思维导图去学习一个东西，这篇论文我也采取的这种模式，如下：

![Escape from Escape Analysis of](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202110252326025.png)

<center>《Escape from Escape Analysis of Golang》摘录</center>

这篇论文是用 Go 语言实现，通过编译和更改源码来优化堆内存使用和降低垃圾回收内存使用率。并且提出了一个验证方法，将该“解释器”应用于 10 个开源项目、16个字节跳动的工业级项目中，结果表明，在开源项目中，堆分配平均降低了 8.88%，堆使用率平均降低了 8.78%，时间消耗平均降低了9.48%，垃圾回收的暂停累积时间平均降低了 5.64%。在工业项目找到了 452 个可优化的案例。别小看这些百分之几的优化，对于大公司来说，业务量都是在千万用户级别，能够降低百分之几可以说节约了好几斤钞票，又可以给员工发福利了。



对我来说，这篇文章标题很新颖，话题很新颖，论证也很新颖。所以有很大的热情能够读下去，即便是英语也没遇到太大的困难。不过后面的几个章节都还没概括出来，还需要再读一读，以后再多做这样的训练。



最后，地址都放在这里了，大家感兴趣可以康一康。

论文：http://www.wingtecher.com/themes/WingTecherResearch/assets/papers/ICSE20.pdf

Github 源码：https://github.com/wangcong15/escape-from-escape-analysis-of-golang
