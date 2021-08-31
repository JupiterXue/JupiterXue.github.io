---
title: "Github_action 和 Review_dog"
date: 2021-08-25T23:10:23+08:00
tags:
   - 兴趣
keywords:
   - 技术
---

最近在写微信公众号的同时，也在同步开通 Hugo 搭建的个人博客。在曹大的课上了解到了两个开源工具：Github Action 和 Reviewdog，都是和  CI/CD  自动化相关的。而我最早是在 B 站一个 UP 主“遇见狂神说” 的 10 分钟短视频讲解中了解到的 CI/CD。惊奇地发现 Hugo 搭建的博客特别适合用到这两个工具。于是怀着特别的兴趣和好奇，去了解了一下。今天仅简单地来谈谈这两个工具的相关介绍。

# CI/CD

![zz](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108252355294.png)

首先，什么是 CI/CD？其实说的是三件事情：持续集成 (Continuous Integration, CI)、持续交付(Continuous Delivery)、持续部署 (Continuous Deployment)。其实这里单独一件事拎出来做也能做个不错的项目，不过三者搭配食用，风味更佳。



持续集成 (Continuous Integration, CI)，通常是指在程序员合并 merge 代码的时候不断对代码变更进行验证。一般用于集成时跑单元测试或者接口测试，如果未通过，那么工作流会用社交软件或者邮件通知相应的开发者。更重要的是，CI 流程中的测试代码既保证了新代码不会破坏老的业务功能，还保证新代码能编译通过。预期输入应该能够得到预期输出，如不符合的结果要报错飘红。



持续交付(Continuous Delivery)的目的是最小化部署或释放过程中固有的摩擦。它的实现通常能够将构建部署的每个步骤自动化，以便任何时刻能够安全地完成代码发布（理想情况下）。



持续部署 (Continuous Deployment)，是一种更高程度的自动化，无论何时对代码进行重大更改，都会自动进行构建/部署。能够将部署流程平台化，可以按照天、周、双周发布。发布只需要点一个按钮，就可以把代码部署到测试或者线上环境。并且即便每次改动都很小，部署流程也能够在有问题的时候及时发现。目前大多数公司都在使用这种流程。

# GIthub Action

![img](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108252353712.png)

CI/CD 持续集成可以玩很多有意思的东西，比如自动抓取数据、定期测试代码、一键打包项目自动登录远程服务器并发布到第三方服务等等。



作为全世界最大交流平台，GitHub 在这方面有一定话语，把这些操作统称为 actions。GitHub 注意到，由于许多操作在不同项目里面是类似的，且可以实现完全共享。于是想出了一个很妙的点子，允许开发者把每个操作写成独立的脚本文件，存放到代码仓库，这样其他开发者也可以使用。如果你需要某个 action，不必自己写复杂的脚本，直接引用他人写好的 action 就好。因此，整个持续集成过程，就变成了一个 actions 的组合。

# Review Dog

![reviewdog_M312q](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202108252353675.png)

reviewdog 是一种可以与任何代码分析工具持续集成的自动化代码审查工具，无论是什么样的编程语言都可以接入。



reviewdog 提供了一种特殊方法，通过与任何linter工具轻松集成，自动将审查意见发布到代码托管服务。它使用 lint 工具的输出，并在发现需要审查的补丁时将其作为评论发布。正如其名字所说，拥有一只能够做代码评审的狗勾，让你的代码库保持健康。



结语，后面有机会对这两个工具的使用做一个详细的介绍，敬请期待。



__参考资料__

[1] 【狂神说】CI/CD到底是什么？十分钟理解企业级DevOps, Bilibili

https://www.bilibili.com/video/BV1zf4y127vu?from=search&seid=15753722902113310493

[2] 什么是 CI/CD？Linux.cn

https://linux.cn/article-9926-1.html

[3] 什么是CICD，CSDN

https://blog.csdn.net/weixin_44903147/article/details/96291588

[4] GitHub Actions 入门教程, 阮一峰

https://www.ruanyifeng.com/blog/2019/09/getting-started-with-github-actions.html

[5] reviewdog, 曹大博客

 https://xargin.com/add-reviewdog-for-your-project/

[6] Deploy Hugo as a GitHub Pages, Hugo 官方文档

https://gohugo.io/hosting-and-deployment/hosting-on-github/

[7] Hugo+GitHub Action+Github Pages搭建个人博客

https://pingyangblog.com/setup-hugo-blog/
