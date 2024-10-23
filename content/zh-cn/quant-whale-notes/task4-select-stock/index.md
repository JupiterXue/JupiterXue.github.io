---
title: "投资与量化-任务4-量化选股"
date: 2024-10-23
category: learn
math: true
toc: false
tags:
- Quant
keywords:
- 思考
---


传送门：[第四章-量化选股](https://datawhalechina.github.io/whale-quant/#/./ch04_%E9%87%8F%E5%8C%96%E9%80%89%E8%82%A1%E7%AD%96%E7%95%A5/ch04_%E9%87%8F%E5%8C%96%E9%80%89%E8%82%A1%E7%AD%96%E7%95%A5)


## 思考
- 基础有效市场理论模型的讨论和研究，短期市场3-5天无法预测价格，长期3-5年可预测走势。可以说短期的“市场先生”行为和长期的价格-价值均值回归。因此量化优势就是在短期内捕捉到市场无效的低价格和长期看低价格回归高价格的股票。
- APT 模型解释了市场回归均衡的一个决定因素，套利行为，如果市场未达到均衡状态市场就有无风险套利机会。单看套利行为好像没什么，单纯赚，但其核心和对冲很相像。就是寻找更低风险，甚至零风险的高收益的组合，而不是单一的看多看空行为。比如06年4月到07年2月，认沽行权价、认购行权价和公司利益形成了对冲，当时的五粮液、包钢认购权证的走势就是无风险套利行为。
- 2010年交易所推出 level2，可了解日内高频行情信息，是 ALPHA 因子研究的重要数据来源。高频数据有更高数据密度和不同数据结构，但也可以定义为时间序列来处理。除了单日数据量大，时间单位可能也不整齐。

## 笔记
- **高频量化**处理注意事项：
  - **分段处理**。往往按表的或按时间窗口划分。
  - **低频化**。将间隔不均匀的时间序列变为均匀时间间隔。
  - **时间窗口避免跨交易日**。集合竞价和隔夜信息会增加随机性，给抽样带来干扰。
- 因子中性化（正交化），用于剔除已有常见因子影响的方法。具体操作中：
  - 市值中性化，对市值取对数然后进行线性回归
  - 行业中性化，对行业均值做差和回归取残差。
- 多因子选股模型的核心思想：通过多个因子的组合来选择股票，因子可以通过历史数据来计算，然后用来预测未来的股票表现，以获取更全面更稳定的预测，包括：
  - **基本面因子**：市盈率PE、市净率PB、营收增长
  - **技术分析因子**：动量、波动率
  - **宏观经济因子**：利率、通货膨胀
- 多因子模型步骤
  1. 确定目标和约束条件。目标收益率、风险水平、组合数目限制、行业比例限制。
  2. 选择因子并计算。根据目标和约束条件，选择合适因子，如PE、PB等。
  3. 异常值处理。检查数据中异常值和报错数据，进行处理和滤除。
  4. 因子标准化。不同因子取值范围差异很大，进行标准化处理，如去均值和缩放。
  5. 确定因子权重。通常用统计方法，如主成分分析来确定。
  6. 构建多因子模型。结合因子值和权重，建立评分模型。
  7. 股票筛选和组合优化。根据评分筛选股票，并进行组合优化，获得符合目标和约束条件的组合。
  8. 回测和调整模型。用历史数据回测多因子模型效果，根据结果调整和改进。
