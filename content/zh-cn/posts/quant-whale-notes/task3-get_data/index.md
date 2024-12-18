---
title: "投资与量化-任务3-获取数据"
date: 2024-10-22
category: learn
math: true
toc: false
tags:
- Quant
keywords:
- 思考
---


传送门：[第三章-获取数据](https://datawhalechina.github.io/whale-quant/#/./ch03_%E8%82%A1%E7%A5%A8%E6%95%B0%E6%8D%AE%E8%8E%B7%E5%8F%96/ch03_%E8%82%A1%E7%A5%A8%E6%95%B0%E6%8D%AE%E8%8E%B7%E5%8F%96)


# 思考
- 获取数据 BaoStock。主要练习以 Python 开源金融数据分析包 BaoStock 来获取数据，开启量化 coding 第一步。
- 数据处理 Pandas。开源的金融数据类型几乎都解析为 Pandas 格式数据类型，便于处理和分析。注意区分Series 和 Dataframes，type 查询其类型，时间类型有专门的，索引非常重要要重视，join可做数据视图，bool可观察可作条件
- 本章以代码基础练习为主，梳理重用的方法和数据类型，更好的方式还是到类似聚宽的金融平台去作一个完整的 demo，跑出获取数据+均线系统+绘图+回测（平台已自带）。算是一个完整的练习。（读完一遍聚宽入门量化的文章才发现，吃透双均线策略是走出新手村的四个任务之一，所以还是多打好基础）
- 关于获取数据，聚宽有一文章专门讲解：[获取数据-聚宽](https://www.joinquant.com/view/community/detail/c688e86342b472f380c8fb9fc58eec54)
- “某些知识点忘了很正常，回头再看就行，用什么去学什么学习效率更高。”引用聚宽文章，个人认为更重要的是，了解有什么知识，建立索引（做笔记或记号），下次再遇到类似问题的时候通过索引找到知识然后用起来。一个查理芒格的隐喻说的好：在拿锤子的人眼里，全世界都是钉子。学习的目的就是如此，避免手里和眼里只有一个“锤子”。
# 拓展
以下为在本章基础上的其他学习补充：来源聚宽博客学习：
- **回测**：让计算机根据一段时间历史数据模拟执行策略，根据结果评价并改进策略。代码层面来讲，就是用一段时间历史真实行情数据来验证一个确定的交易策略在这段时间内表现如何。
- **模拟交易**：让计算机根据实际行情模拟执行该策略一段时间，根据结果评价并改进策略。
- 更**客观的评估策略指标**：年化收益率、最大回撤率、夏普比率。
- 运行回测，执行回测并得到结果，回测结果：收益率、年化最大收益率、最大回撤、夏普比率、下单记录、持仓记录。
- 编译运行，简化版运行回测，只有收益率，速度更快。
- 常用下单函数：
  - order(code, amount) 买一定数量股票
  - order_target(code, amount) 股票调仓，高于卖出，不高不低就不动
  - order_value(code, value)买卖一定价值量。
  - order_target_value，股票调仓至一定价值量。
  - content：存储当前策略运行时间点、持有股票及其数量、持仓成本，一些常用
  - 当前时间 context.current_dt
  - 当前时间格式化：context.current_dt.strftime("%Y-%m-%d")
  - 前一个交易日：context.previous_date
  - 当前可用金：context.portfolio.positions_value
  - 累计收益 context.portfolio.returns
  - 当前持有股票 context.portfolio.positions.keys()
  - 当前持有某股票的开仓均价：context.portfolio.positions['code'].avg_cost
  - 某股票当前现价：context.portfolio.positions['code'].price
  - 当前持有某股票的可卖持仓量：context.portfolio.positions['code'.closeable_amount
- **止损**：狭义是亏损一定幅度下单卖出，减少进一步亏损。广义止损在思路上衍生复杂的减少亏损方法。
- 最大回撤越小越好
- 足够多的交易次数让回测结果更有说服力，不能直接看到，可通过盈利次数与亏损次数相加得到，或每日买卖大致看出。
- 有数理基础，进一步学习 Alpha 与 Beta 构造思路与过程
- **夏普比率**，单位风险收益，越大越好
- 未来函数，策略利用了历史当时无法得到的信息，造成回测结果极大失真。一般人工排查，注意时间。更好方法是用策略建立模拟交易，运行几天，多数未来函数问题都能被发现。模拟交易不可能引入未来数据，往往引入未来函数的策略都无法成功运行模拟交易。
- 判断是否引入未来函数的经验法则：策略回测收益率极高，回撤极低，且各个时间段表现特别好，仿佛发现了交易圣杯，此时大概率引入了未来函数
- [enable_profile 性能分析](https://www.joinquant.com/help/api/help?name=api_old#%E6%80%A7%E8%83%BD%E5%88%86%E6%9E%90%E2%99%A0)
  - 复制粘贴放到策略代码的第一行，回测就有性能分析结果
- **过拟合**：回测好，模拟盘或实盘烂，很可能参数过拟合。
  - 核心思想：过度细致地解读样本数据，从而没有认识到本质的规律，使策略失去普适性，对原本数据表现极为优异，对非原样本数据外情况有效性大大降低。
  - 案例：背卷子答案成绩极高，换题依旧烂
  - 鲁棒性：策略好坏对参数变化的敏感性。建议，控制参数梳理，多测几组参数大致看下参数变化对策略的影响，另外考虑进行样本外才是，用一份样本回测所选参数，另一份样本数据回测看参数在样本外情况的表现。
- 策略有时效性。最大回撤历史新高，收益率不增减少。都是现象，判断策略失效方法：
  - 逻辑基础不成立。涨跌停制，某行业发展、某资源全球持续稀缺背景下。如果制度调整，新政策发布，科技进步，策略自然失效。
  - 资金过大。市场流动性不足承载，买入价更高、卖出价更低。
  - 市场中相似策略过多。策略竞争导致赚钱变难。交易行业非常注意保密且不适合分享，有志者培养自学能力。
  - 寄生策略。有心人发现你的策略并足够程度的监测时，针对你的策略制定策略，寄生在你策略上，压缩盈利空间。
- 往往收益能力与抗风险能力互相制约不能兼顾，建议是，达到基本盈利能力，极力追求低风险，理由是盈利水平往往可以朝桐光增加资本家来提高。
- 量化不适合分享，自学必不可少。
- 彻底读懂下面四篇文章，作为走出新手村的一个历练：
  - [双均线策略](https://www.joinquant.com/post/1398)  （MACD）
  - [彼得林奇成功投资](https://www.joinquant.com/post/1957)  （PEG）
  - [Fama-French 三因子火锅](https://joinquant.com/post/1668) （多因子）
  - [凯利公式](https://www.joinquant.com/post/1311) （期望）
- [编程入门资料](https://joinquant.com/post/10760)

## 量化进阶
- 基本模式：有一个灵感，研究与实现，成败获得新认知
- 灵感来源：论文、研报、聚宽论坛、知乎量化主题问答、宽客博客（多数国外）、各类书籍
- 广泛阅读：不限量化，还有财务分析、行为金融、熟悉、统计学、机器学习、心理学、物理学、语音识别
- 学习时批判眼光看，独立思考，警惕盲从。量化交易很难有特别有价值的东西出现在大众眼前。
- 直接参与到市场中。市场对时下新闻的反应、对量价走势的瞬时反应、不同股票或人群的反应、你的反应等

## 职业化
- 学历只是知识储量与提效手段，你应该被录取的理由是什么？怎么把这个理由告诉他们？比如有一个好策略想入职带他们赚钱，或者有相当的策略实现能力或因子发掘能力。
- 独立开发。多数公司独立出来，或学生阶段转独立职业量化，但现金流压力较大后来成立公司。
- 聚宽栏目，策略被看中要谈买断、发文章要被录取，注意保护自己的成果。