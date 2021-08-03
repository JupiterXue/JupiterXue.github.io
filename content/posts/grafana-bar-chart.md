---
title: "Grafana 条形图实践"
date: 2021-07-29T09:56:39+08:00
draft: true
---

# 前言

可视化平台有很多，最重要的是选择符合业务需求的方案，费用低、学习曲线低、易于上手、集成性好这些都是加分项。众望所归，我选择了 Grafana。

## 特点

让我们康康 Grafana 是怎样的：
1. 免费。Grafana 是一个**开源**的独立日志分析和监视工具。
2. 开发少可直接用。面向分析师和一般使用者，一般**不需要做多少编码**工作就能直接拿来分析数据、搭建可视化系统。
3. 界面炫酷。来看看官网的展示 [Grafana 官网概述](https://grafana.com/grafana/)
4. 丰富集成。Grafana 是跨平台工具，它提供了与各种平台和数据库的集成，持 InfluxDB，AWS，MySQL，PostgreSQL 等。
5. 用户交互体验好。展示连续**实时监控**指标（如 CPU 负载，内存）、对数据提供自定义**实时警报**、提供基于数据库及其查询语法的**命令列界面**。

## 准备工作-采集器

一个可视化平台，重要的是有数据来做支撑，像上面所说的数据库。本文采取具有灵活查询和实时报警构建的时序数据库 **promethues** 来作为数据管理端。

作为圣火的传人，promethues 有着各种采集器小弟，如服务器资源采集器 node-exporter、进程资源采集器 process-exporter、服务资源采集器 blackbox-exporter，
有了这些采集器就可以直接搭建起监控和警报

**[服务器监控](https://grafana.com/grafana/dashboards/8919)**
![img](https://grafana.com/api/dashboards/8919/images/8260/image)

**[进程监控](https://grafana.com/grafana/dashboards/249)**
![img](https://grafana.com/api/dashboards/249/images/439/image)

**[服务监控](https://grafana.com/grafana/dashboards/9965)**
![img](https://grafana.com/api/dashboards/9965/images/6248/image)

最后再配有 docker 微服务，你的监控和警报平台就能够一键部署起来了。
![img](https://img2.baidu.com/it/u=1337695316,4281005371&fm=26&fmt=auto&gp=0.jpg)

## 准备工作-定制化
有了各种采集器，其实能够满足大部分通用监控的需求，但往往我们的业务是不一样的，那么就需要进行定制。

用到了胶水语言 Python 在 promethues 的模块 promethues_client，就能开始定制了！


## 制作一个条形图
其实这个标题，我最先想取 “fuck the bar chart”，其他定制化的数据，表格、饼图、折线图都是能够在粗略熟悉 Grafana 后可以自己上手做的。

而这个条形图困扰了我整整一周的时间，各种找插件、换版本、调格式、选展示方式，都没有达到理想效果：

今天来彻底搞定这个问题！！

### 参考链接

[1] 

[2] [A Complete Guide to Bar Charts](https://chartio.com/learn/charts/bar-chart-complete-guide/)

[3] [A Complete Guide to Stacked Bar Charts](https://chartio.com/learn/charts/stacked-bar-chart-complete-guide/)

[4] [《How to Choose the Right Data Visualization》](https://cdn2.hubspot.net/hubfs/392937/How-To-Choose-The-Right-Data-Visualization.pdf)

# 其他补充

## 1. 匿名访问
```bash
docker run 
    -itd --name grafana -p 3000:3000 
    -e "GF_AUTH_PROXY_ENABLED=true"
    -e "GF_AUTH_ANONYMOUS_ENABLED=true"
    -e "GF_SECURITY_ALLOW_EMBEDDING=true"
    grafana/grafana
```
这里解释一下：GF_SECURITY_ALLOW_EMBEDDING,用于网页嵌套的。如果需要在 iframe 中嵌套 Grafana，那么把这个变量设置为true；


## 2. 全屏化展示
想让 Grafana 分享的面板全屏化，在分享的 URL 后面跟上
```bash
&kiosk=true
```
再加上按钮 **F11**，食用效果更佳。


## 3. 支持 Emoji 啦 😄
虽然 Grafana 的显示都是文本不能显示图片 🤔，但是支持 emoji ！！！✌

[1] [emoji 点击复制版](http://emojihomepage.com/) 👈

[2] [emoji 大全版](https://unicode.org/emoji/charts/full-emoji-list.html) 🔮

# 结语
今天关于 Grafana 的分享就到这里结束啦，如果你感兴趣，欢迎联系我 🤙