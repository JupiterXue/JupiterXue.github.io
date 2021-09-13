---
title: "读经典《DDIA》-第二章"
date: 2021-09-11T23:02:05+08:00
category: tech
tags:
    - 经典书
keywords:
    - 图模型
---

已经看完第二章《DDIA》，说下我是怎么抄书的，以及我有什么收获。



0.什么？抄书？没错，我是边抄边写笔记。因为我就和你一样，看不进书（当然你比我更自律点）。开玩笑啦



1.当你不想看书的时候，就抄书；当你抄书抄到不想抄了，就想看了。虽然这句话被 S 老师作为段子来传颂，但我看到了有人抄过化学，有人抄过高数。我没有这种经历，也想亲自去实践一下到底有怎样的魔力。其次，我认为好的书是浅显易懂，深入浅出的，值得揣摩。



2.求其上者得其中，求其中者得其下，求其下着无所得。刚开始抄书是一件有点兴奋的事，因为抄这个动作是比较轻松的，因为理想中这是不怎么动脑子的。实际上由于你的专注，你会不由自主地关注作者在说些什么（虽然有时也会走神不知道作者在说什么）。并且由于抄书，你会开始对作者某个词语进行或多或少的思考。而抄书的困难也正是源于此，认真看和思绪漫游会导致你的阅读速度极大地降低，尤其在你发现抄了一两个小时，连一半都没有的时候，心情是有点绝望的，好在，你还可以明天抄嘛：） giao！



3.实事求是，实地考察。抄书中，你会发现作者写的话有多少冗余、多少用词不当。尤其我看的这本书原版是英文，翻译过来连句法还是英语的语法，根本不符合中文逻辑。所以我会将这句话进行汉译汉，转化为我自己能够理解的。大家都认为看技术不要看中文，因为中文和蹩脚，是的。但我再补充一句，你可以自己把蹩脚化为流畅。



4.简明扼要，提高速度。一直这么慢的抄写肯定不行，我在抄的过程中也逐步掌握了一些技巧：着重摘录抽象句式和实用句式，代码复制（看代码其实是最直接和方便的），详细参数说明忽略，图能动就不抄，不能动就用文字描述，不能用文字就截图、先摘录一二级标题，然后一个个攻坚（会有点成就感）。



5.暴力破解，巨无事细。抄书意味着会对每个地方都看一遍，找出是否有价值的地方。所以会对每个文字做一次简单的“有什么有用，是否摘要”的编码动作，这不同于看书。看书的时候，我很有可能是在想“这个不重要，过”。所以抄一遍下来，老师如果问哪一页哪个知识点，虽然你答不上来在哪里，但你知道，我确实看过，而不是看，过了。



6.建立根据地，逐步发展。即便中文的技术书读起来不是很顺畅，但大体意思作者是能够把握。最重要的是，当我看完了中文，我会对原书作者讲了什么内容有个大概的把握。所以看起英文来就更加流畅。



最后，这篇文章不是给看中文技术书洗白，只是想表达我在参与过程中的收获和感受。按照成熟的技术方法来说，看书、看文档还是以官方为准，以原版语言为准。地道的语言和没有信息损失的一手材料，才是好的学习资料。毕竟，高级食材，我们都吃原材料。



和小伙伴讨论章节内容，一个意外收获：

1. 我说：几种 Datalog 方法都没有听过，没有见过。好像是在说图结构有哪些处理方式，怎么表示更简洁、更高效

   Helios 回复：主要还是场景不一样。

   ![image-20210912001019822](C:\Users\Xfavor\AppData\Roaming\Typora\typora-user-images\image-20210912001019822.png)

   场景说得好，看完就局限在代码层面，一下又跳了出来。

# 第二章笔记

**第二章**

数据模型在软件开发中最重要。不仅影响软件编写方式，还影响解题思路。

多数应用用使用层层叠加的数据模型构建，关键问题：如何用低一层数据模型表示。例如：

\1. 用对象或数据结构以及 API 来建模现实世界。

\2. 用数据模式表示存储结构。如 JSON、XML、数据库表、图，

\3. 用内存、磁盘或网络字节表示 JSON/XML/关系/图数据，进而来查询、搜索、操作。

\4. 在更低层次，用电流、光脉冲、磁场或其他东西表示字节。

复杂应用程序有更多中间层，如 API 的 API。但思想仍一样：提供明确数据模型来因此更低层次复杂性。这个抽象使得不同人员能够参与协作。

选择一个适合数据模型非常重要。因为种类很多，易用但不易支持，可操作但表现差，数据转化有的自然有的麻烦。

**关系模型与文档模型**

最著名数据模型——SQL。1970年 Edgar Codd 提出关系模型：数据为关系，关系是元素的无序集合。

关系数据库起源于商业数据处理，今天来看显得很平常：典型事物处理。

\> 如：将销售或银行交易，航空公司预订，库存管理信息记录在库）和**批处理**（客户发票，工资单，报告）

并且当时数据库迫使开发者必须考虑数据库内部的数据表示形式。关系模型解决的是在实现细节隐藏在更简洁的接口之后。

网络模型和分词模型在 70、80 年代的主要选择，但关系模型随后占据主导。对象数据库在 90 年代由盛而衰。XML 数据出现于 21世纪初。

2010 年，NoSQL 开始萌芽，虽名字没有涉及任何技术，最初只是 Twitter 标签，2009 年非关系数据库开源会上这个术语的出现，迅速得到了传播。NoSQL 被重新解释为不仅是 SQL（Not Only SQL），NoSQL 流行的几个因素：

\- 更优于关系型数据库的课伸缩性，能够处理大数据集、具有高写入吞吐量

\- 免费和开源受到青睐

\- 更好地支持特殊查询操作

\- 人们渴望具有多动态性和表现力的数据模型

不同应用程序有不同需求，最佳技术选择可能不同于另一个用例的最佳技术选择。因此，在未来，关系数据库可能和非关系数据库一起使用——混和持久化（polyglot persistence）

**对象关系不匹配**

多数程序都用面向对象编程语言来开发，导致人们普遍对 SQL 的批评：数据存储在关系表中，需要做个笨拙的转换层，才能获取数据。一般处于代码中对象和表，行，列的数据库模型之间。模型之间不连贯有时称为阻抗不匹配（impedance mismatch）

对象关系映射（object-relational mapping），也即是我们熟悉的 ORM，这个框架可以减少转换层所需的代码，但不能完全屏蔽两个模型间的差异。这里有个领英简历关系模型图

![img](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109112352380.png)

图中可以发现，通过唯一标识符 user_id 来标识。想 first_name 和 last_name 这样的字段美国用户只出现一次，所以在 User 表上建模为列。像我们大多数人的职业生涯中会拥有多份工作或不同教育背景，在用户到用户项目信息之间存在一对多的关系，可以有多种表示方式：

\- 传统 SQL 模型（1999 年前）。将职位、教育和联系信息放在单独表中，对 User 表提供外键引用。

\- 后续的 SQL 标准增加对结构化数据类型和 XML 数据的支持。允许将多值数据存储在单行内，并支持在这些文档内查询和索引。现有支持这种功能的数据库有：Oracle、IBM DB2，MS SQL Server 和 PostgreSQL。JSON 数据类型也有支持的：IBM DB2，MySQL 和 PostgreSQL

\- 将职业、教育和联系信息编码为 JSON 或 XML 文档，将其存储在数据库的文列中，并让应用程序解析结构和内容。这种配置下，通常不能从该数据库中查询该编码列中的值。

像简历这样包含文档的数据结构，用 JSON 表示非常合适，其相比 XML 更简单。以下为一个实例：

\> 面向文档数据库：MongoDB，RethinkDB，CouchDB 和 Espresso。

{  "user_id": 251,  "first_name": "Bill",  "last_name": "Gates",  "summary": "Co-chair of the Bill & Melinda Gates... Active blogger.",  "region_id": "us:91",  "industry_id": 131,  "photo_url": "/p/7/000/253/05b/308dd6e.jpg",  "positions": [    {      "job_title": "Co-chair",      "organization": "Bill & Melinda Gates Foundation"    },    {      "job_title": "Co-founder, Chairman",      "organization": "Microsoft"    }  ],  "education": [    {      "school_name": "Harvard University",      "start": 1973,      "end": 1975    },    {      "school_name": "Lakeside School, Seattle",      "start": null,      "end": null    }  ],  "contact_info": {    "blog": "http://thegatesnotes.com",    "twitter": "http://twitter.com/BillGates"  } }

开发人员认为 JSOn 减少了应用代码和存储层之间的阻抗不匹配，但其实 **JSON 作为数据编码也存在问题**——这通常是因为所获取的数据在存储或者发送的过程中为了节省空间，已经被转译为了其他形式的编码。当我们直接使用 json.dump 操作的时候就会出现一个编码问题的报错，想要解决这个错误的出现，首先需要确认数据中都包含了哪几种数据类型，然后构建对应的类来进行数据的过滤筛查，将这些不符合要求的数据类型的数据进行转译。

JSON 具有良好的局部性（locality）。在前面的关系型示例中获取简历信息，需要执行多个查询（用 user_id 来查表），或者在 User 表和它下面的表之间混乱地执行多表连接。而 JSON 会将所有相关信息都放在同一个地方，一次查询就可以得到。

从用户信息到用户教育信息、职位信息和联系信息，一对多关系隐含了数据中的一个树状结构，而 JSON 表示让这样一个树状结构变得清晰明确。

![img](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202109112352517.png)

**多对一和多对多关系**

如果在用户界面想用自由文本字段来输入表示区域和行业，那么用纯文本字符串存储是合理的。另一个方式是先给出地理区域和行业标准的列表，让用户在下拉列表或自动填充器中进行选择，这样的优势在于：

\- 各个简介之间样式和拼写统一

\- 避免歧义

\- 易于更新。名称只存储在一个地方，如果需要更改（如因政治事件改变城市名称），很容易进行全面更新。

\- 本地化支持。网站被翻译成其他语言时，标准化列表也可以被本地化，即地区和行业也可以用翻译语言来显示

\- 更好地搜索。如搜索关键词“Bill”就会匹配这份简历。

存储 ID 还是文本字符串，本质是冗余副本问题（duplication）。当用 ID 时，对人有意义的信息之存储在一个地方，所有引用都只用关注 ID。而当直接存储文本时，和人有意义的信息会复制在每次使用的地方。

使用 ID 的好处在于，ID 对人没有任何意义，所以不需要改变。而任何对人有意义的东西都可能在未来某个时候改变，也就意味着，如果这些信息被复制，所有的冗余副本都需要更新。这就会导致写入开销增加，存在不一致风险。这个重复性问题在数据库中解决的关键思想是规范化（normalization）。

数据进行规范化需要多对一的关系，这用文档模型不太合适。在关系数据库中，通过 ID 来引用其他表的行是正常的，因为连接容易。而在文档数据库中，一对多结构没必要连接，并且对连接的支持通常来说比较弱。

如果数据库本身不支持连接，那么必须在代码中通过对数据库进行多次查询来模拟连接关系。

**实体关系**

上面图中，organization 和 school name 都是字符串。如果我们想将它们作为实体来引用呢？美国组织、学校或大学都有自己的网站。美国简历可以链接到它提到的组织和学习，并且包括图标和其他信息。

**推荐**

假如你想添加新功能：一个用户为另一个用户写推荐。在用户的简历上显示推荐，并附上推荐用户的姓名和照片。如果推荐人更新他们的照片，那么他们写的任何推荐都需要显示新的照片。所以，推荐应该有作者个人借鉴的引用。

这个新功能需要建立多对多关系。

**文档数据库正在颠覆？**

在多对多关系和连接常规的关系数据库时，文档数据库和 NoSQL 发生了争论。这个问题最早可以追溯到计算机化数据库系统。

IMS 最初设计了相当简单的层次模型（hierarchical model），它和文档数据库用的 JSON 模型有些类似。它将所有数据嵌套地记录在树中。虽然 IMS 能处理好一对多关系，但还是很难处理多对多关系，并且不支持连接。开发者必须去复制非规范化或手动解决一个记录到另一个记录的引用。

因此当时的人们提出了不同的解决方案来解决这个层次模型的局限性。最突出的两个是关系模型（relational model）和网络模型（network model），前者变成了 SQL 统治了世界，后者刚开始很热门后来变得冷门。

**网络模型**

网络模型由一个成为数据库语言会议（CODASYL）的委员会进行了标准化处理，并被不同的数据库厂家实现，它也被称为 CODASYL 模型。

CODASYL 模型是层次模型的推广。在层次模型的树状结构中，每条记录只有一个父节点。而在网络模型中，每个记录可能有多个父节点。

网络模型中记录之间的链接不是外键，而更像编程语言的指针（同时存储在磁盘上），访问记录的唯一方法是从根记录沿着这些链路所形成的路径，也叫做访问路径（access path）

最简单的情况下，访问路径的过程类似遍历链表：从列表头开始，每次查看一条记录，直到找到所需的记录。但在多对多的关系中，数条不同的路径可以达到相同的记录，网络模型的程序员必须跟踪这些不同的访问路径。

CODASYL 中查询是通过遍历记录列和跟随访问路径表在数据库中移动游标来执行。如果记录有多个父节点（即多个来自其他记录的传入指针），则应用程序代码必须跟踪所有的各种关系。甚至 CODASYL 委员会也承认，这就像在 n 维数据空间进行导航。

尽管手动访问路径能够最有效地利用硬件功能（如磁带驱动器，但其搜索速度非常慢），但这使得查询和更新数据库的代码变得复杂不灵活。无论是分层模型还是网络模型，如果没有数据路径，都难办。可以改变访问路径，但必须浏览大量手写数据库查询代码，并且重写来处理新的访问路径。更改应用程序的数据模型是很难的。

**关系模型**

一个关系（表）只是一个元组（行）的集合。如果想读取数据，既没有迷宫似的嵌套结构，也没有复杂的访问路径。也可以选择符合任意条件的行，读取表中任何或所有行。通过指定某些列作为匹配关键字来读取特定行。可以在任何表中插入新的行，而不必担心和其他表的外键关系。

\> 外键约束允许对修改进行限制，但对于关系模型这并不是必选项。即使有约束，外键连接在查询时执行，而 CODASYL 中，连接在插入时高效完成。

在关系数据库中，查询优化器自动决定查询的哪个部分以哪个顺序执行，以及使用哪些索引。这些选择实际上是“访问路径”，但最大的区别在于它们是由查询优化器自动生成的，而不是由开发者生成。

如果想按新的方式查询数据，可以声明一个新的索引，查询会自动使用最合适的索引。而无需更改查询来利用新的索引。关系模型因此使添加应用程序新功能变得更加容易。

关系模型的一个关键洞察：只需要构建一次查询优化器，随后使用该数据库的所有应用程序都可以从中受益。如果没有查询优化器，那么特定查询手动编写访问路径比编写通用优化器更容易。

**与文档数据库相比**

文档数据库可还原为层次模型：父记录存储嵌套记录，而不是在单独的表中。

在多对一和多对多关系中，关系数据库和文档数据库没有根本的不同：相关项目都被唯一的标识符引用，这个标识符在关系模型中被称为外键，在文档模型中称为文档引用。标识符在读取时通过连接或后续查询来解析。

**关系型数据库与文档数据库的对比**

考虑多方面差异：容错性、处理并发性。

支持文档数据模型的主要论据是架构灵活性，因局部性而拥有更好的性能，以及对某些应用程序而言更接近应用程序使用的数据结构。关系模型通过为连接提供更好的支持以及支持一对多和多对多的关系。

**那种数据模型更有助于简化代码？**

如果数据有类似文档的结构，那么文档模型不错。将类似文档结构分解为多个表和关系技术会更复杂。

文档模型有一定局限性：不能直接应用文档中的嵌套项目。但如果文件嵌套不太深，问题也不大。

文档数据库在连接上也比较糟糕。

如果程序需要多对多的关系，文档模型就不太适合。虽然可以通过反规范化来消除对连接的需求，但需要代码做些额外工作以确保数据一致性。代码可以通过向数据库发出多个请求来模拟连接，但速度会更慢，导致更差的性能。

没有完全解决数据模型问题又能够简化代码的方案，因为它取决于数据项之间的关系种类。对高度关联的数据而言，文档模型是极其糟糕的，关系模型是可接受的，而图形模型是最自然的。

**文档模型中的模式灵活性**

多数文档数据库和关系数据库中的 JSON 都不会强制文档中的数据采用那种模式。关系数据库的 XML 支持通过带有可选的模式验证。没有模式意味着可以将任何键和值添加到文档中，并且当读取时，客户端无法保证文档可能包含的字段。

文档数据库有时有叫做无模式（schemaless），名字上有点误导，因为读取数据的代码通常假定某种结构，即存在隐式模式，但不由数据库强制执行。一个更精确的术语是**读模式（schema-on-read**）。（数据结构是隐含的，只有在数据被读取时才被解释），相应的是**写模式（schema-on-write**）（传统的关系数据库方法中，模式明确，且数据库确保所有的数据都符合其模式）。

读模式类似编程语言中的动态（运行时）类型检查，而写模式类似于静态类型检查。就像静态和动态类型检查的相对优点具有很大的争议性一样。

在应用程序想要改变数据格式的情况下，不同方法有明显的区别。比如，把每个用户的全名存储在字段中，而现在想分别存储名字和形式。在文档数据库中，只需要开始写入具有新字段的新文档，并在应用程序中使用代码处理读取旧文档的情况。

另一方面，在 “静态类型”数据库模式中，通常会执行迁移（migration）操作。

模式变更的速度很慢，并且被要求终止。多数数据库可在几毫秒执行 ALTER TABLE，但 MySQL 有时需要花费几分钟甚至几个小时来处理。

大型表上运行 UPDATE 操作在任何数据库上都可能很慢，因为每一行都需要重写。如果难以接受，应用程序可以将 first_name 设置为默认 NULL，并且在读的时候再填充，就像使用文档数据库一样。

当由于某个原因（如，数据是异构的）集合中的项目并不都具有相同的结构时，读模式更具有优势。例如，如果：

\- 存在许多不同类型的对象，将每种类型的对象放在表中是不现实的。

\- 数据的结构由外部系统决定。无法控制外部系统并且它随时可能变化。

基于以上两个原因，模式的坏处大于其益处，无模式文档可能是一个更加自然的数据模型。但是，要记录都具有相同的结构，那么模式是记录并强制这种结构的有效机制。

**查询的数据局限性**

文档通常以单个连续字符串形式进行存储，编码为 JSON，XML 或其二进制变体（如 MongoDB 的 BSON）。如果程序需要访问整个文档（如渲染到网页），那么存储局部性会带来性能优势。如果蒋数据分割到多个表中，就需要多次索引查找才能全部找到，这可能需要更多的磁盘查找并花费很多时间。

局部性仅仅适用于同时需要文档绝大部分内容情况。数据库通常需要加载整个文档，即使只访问小部分，对大型文档的资源也是浪费的。尤其更新文档时，通常需要整个重写。只有不改变文档大小的修改才可以容易地原地执行。所以，通常建议保持相对小的文档，并避免增加文档大小的写入。这些性能限制大大减少了文档数据库的实用场景。

注意，为了局部性而分组集合相关数据的想法并不局限于文档模型。如： Google 的 Spanner 数据库在关系数据模型中提供了同样的局部性属性，允许模式声明一个表的行应该嵌套在父表中。 Oracle 类似地允许使用叫做多表索引集群表（multi-table index cluster tables）的类似特性。Bigtable 数据模型（用于 Cassandra 和 HBase）中的列族（column-family）概念与管理局部性的目的类似。

**文档和关系数据库的融合**

2000 年开始，多数关系数据库（MySQL 除外）都已支持 XML。包括对 XML 文档进行本地修改的功能，以及在 XML 文档中进行索引和查询的功能。这允许应用程序使用与文档数据库应当使用的非常类似的数据模型。

PostgreSQL 自 9.3 版开始，MySQL 自 5.7 开始以及 IBM DB2 自 10.5 开始，都对 JSON 文档提供了类似的支持级别。因为 Web APIs 的 JSON 化趋势，其他关系数据库很可能也会这样。

在文档数据库中，RethinkDB 在查询中支持类似关系的连接，一些 MongoDB 驱动程序可以自动解析数据库引用（有效地执行客户端连接，尽管可能比在数据库中执行连接慢，需要额外的网络返还，并且优化更少）。

关系数据库和文档数据库变得越来越相似。这是好事：数据模型相互补充。如果数据库能够处理类似文档的数据，并能够对其执行关系查询，那么应用程序就可以使用最符合其需求的功能组合。

关系模型和文档模型混合是未来数据库一条很好的路线。

**数据查询语言**

引入关系模型，也就引入了新的数据查询新方法：SQL 是一种声明式查询语言，而 IMS 和 CODASYL 使用 命令式 代码来查询数据库。

许多常用的编程语言是命令式的，如：给定一个动物五种的列表，返回列表中鲨鱼。

function getSharks() {    var sharks = [];    for(var i=0; i<animals.length;i++) {        if(animals[i].family==="Sharks") {            sharks.push(animals[i]);        }    }    return sharks; }

关系代数对应是：sharks = Xfamily=!!sharks!!(animals)

定义的 SQL，也遵循关系代数：SELECT * FROM animals WHERE family='Sharks'

命令式语言告诉计算机以特定顺序执行某些操作。逐行遍历代码，评估条件，更新变量，并决定是否再循环一遍。

在声明式查询语言（如 SQL 或关系代数）中，只需要指定所需数据模式 ——符合条件的结果，以及如何将数据转化（如 排序、分组和集合）。数据库系统的查询优化器决定使用哪些索引和哪些连接方法，以及以那种顺序执行查询。

**声明式查询语言**通常比命令式 API 更加**简洁和容易**，还**隐藏了数据库引擎的实现细节**。

在以上的命令代码中，动物列表以特定顺序出现。如果数据库想要在后台回收未使用的磁盘空间，就需要移动记录，这会改变动物出现的顺序。数据库能否安全执行并且不会中断查询？

SQL 示例不确保任何特定的顺序，因此不在意顺序是否改变。如果查询用命令式的代码来写的话，那么数据库就永远不能确定代码是否依赖于排序。SQL 相当有限的功能性为数据库提供了更多自动优化的空间。

最后，声明式语言往往更**适合并行执行**。CPU 速度通过 core 增加变得更快。

**Web 上的声明式查询**

声明式查询语言的优势不仅限于数据库，还可以用在 web 浏览器。例如有个海洋动物网站。查看鲨鱼页面，并标记选中项目。

<!-- HTML --> <ul>     <li class="selected">         <p>Sharks</p>         <ul>             <li>Great White Shark</li>             <li>Tiger Shark</li>             <li>Hammerhead Shark</li>         </ul>     </li>     <li><p>Whales</p>         <ul>             <li>Blue Whale</li>             <li>Humpback Whale</li>             <li>Fin Whale</li>         </ul>     </li> </ul>  // CSS li.selected > p { 	background-color: blue; }

使用 XSL，可以做类似的事情：

<xsl:template match="li[@class='selected']/p">    <fo:block background-color="blue">        <xsl:apply-templates/>    </fo:block> </xsl:template>

这里用到了 XPath 表达式。相同点在于 CSS 和 XSL 的共同之处在于，它们都用于指定文档样式的声明式语言。

JS 中使用文档对象模型（DOM） API，结果如下：

var liElements = document.getElementsByTagName("li"); for (var i = 0; i < liElements.length; i++) {    if (liElements[i].className === "selected") {        var children = liElements[i].childNodes;        for (var j = 0; j < children.length; j++) {            var child = children[j];            if (child.nodeType === Node.ELEMENT_NODE && child.tagName === "P") {                child.setAttribute("style", "background-color: blue");            }        }    } }

从代码量来说，JS 相比 CSS 和 XSL 有更多，更难理解，还有些严重问题：

\- 如果被选定类被移除（用户点击了不同页面），即使代码重新运行，蓝色背景也不会被移除。因此该项目将保持突出显示，直到整个页面被重新加载。

\- 如果要用新的 API （如`document.getElementsBy ClassName（“selected”`）或者提高性能（document.evaluate()），就必须重写代码。一方面浏览器供应商可以在不破坏兼容性情况下提高 CSS 和 XPath 性能。

在 Web 中，使用声明式 CSS 比用 JS 命令式操作要好得多。类似的，在数据库中，用声明式查询语言 SQL 比使用命令式查询 API 好得多。

**MapReduce 查询**

MapReduce 是由 Google 推广的编程模型，用于在多态机器上批量处理大规模的数据。一些 NoSQL 数据存储（包括 MongoDB 和 CouchDB）支持有限形式的 MapReduce，作为在多个文档中执行读查询的机制。

MapReduce 既不是声明式查询语言，也不是完全命令式查询 API，而使处于二者之间的：查询的逻辑用代码片段表示，代码片段会被处理框架重复调用。它基于 map（也称为 collect）和 reduce（也称为 fold 或 inject）函数，两个函数存在于许多函数式编程语言。

一个例子，能够更好地解释 MapReduce 模型。作为海洋生物学家，当看到海洋中动物时，会在数据库中添加一条观察记录。现在想生成一个报告，说明每个月看到多少鲨鱼。可以在 PostgreSQL 中做这样的查询：

SELECT date_trunc('month', observation_timestamp) AS observation_month, sum(num_animals)                           AS total_animals FROM observations WHERE family = 'Sharks' GROUP BY observation_month;

date_trunc 确定了日历月份，并返回代表该月份开始的另一个时间戳。它将时间戳舍入最近的月份。

这个查询首先过滤观察记录，之显示鲨鱼家族的物种，然后根据它们发生的日历月份对观察记录进行分组。最后将该月的所有观察记录中看到的动物数目加起来。

同样的查询在 MongoDB 这样表示：

db.observations.mapReduce(function map() {        var year = this.observationTimestamp.getFullYear();        var month = this.observationTimestamp.getMonth() + 1;        emit(year + "-" + month, this.numAnimals);    },    function reduce(key, values) {        return Array.sum(values);    },    {        query: {          family: "Sharks"        },        out: "monthlySharkReport"    });

另外一个例子，解释 MapReduce：

{  observationTimestamp: Date.parse(  "Mon, 25 Dec 1995 12:34:56 GMT"),  family: "Sharks",  species: "Carcharodon carcharias",  numAnimals: 3 } {  observationTimestamp: Date.parse("Tue, 12 Dec 1995 16:17:18 GMT"),  family: "Sharks",  species:    "Carcharias taurus",  numAnimals: 4 }

对每个文档都会调用一次`map`函数，结果将是`emit("1995-12",3)`和`emit("1995-12",4)`。随后，以`reduce("1995-12",[3,4])`调用`reduce`函数，将返回`7`。

map 和 reduce 函数在功能上有所限制：它们必须是纯函数——传递给它们的数据作为输入，它们不能执行额外的数据库查询，不能有任何副作用。这个限制允许数据库以任何顺序运行功能，并在失败时重新运行。

map 和 reduce 强大在于：它们可以解析字符串，调用库函数，执行计算。

MapReduce 是一个相当底层的编程模型，用于在集权上执行分布式。因此，更高级的 SQL 可以用一些列  MapReduce 来实现。

\> 注意：SQL 中没有任何内容限制它在单个机器上运行，而 MapReduce 在分布式查询上没有垄断权。

在查询时使用 JS 代码是高级查询的特征，但不限于 MapReduce，一些 SQL 数据库也能用 JS 函数进行扩展。

MapReduce 的可用性问题在于：必须编写两个密切合作的 JS 函数，这通常比写单个查询更难。并且，声明式查询语言为查询优化器提供了更多机会来提高查询的信念。基于这些原因，Mongo 2.2 添加了聚合管道的声明式查询语言支持，一个鲨鱼计数查询的例子如下：

db.observations.aggregate([  { $match: { family: "Sharks" } },  { $group: {    _id: {      year:  { $year:  "$observationTimestamp" },      month: { $month: "$observationTimestamp" }    },    totalAnimals: { $sum: "$numAnimals" } }} ]);

聚合管道语言和 SQL 的子集有类似的能力，但其用基于 JSON 的语法而不是 SQL 的英语句子语法。

**图数据模型**

多数关系是不同数据模型之间具有区别性的重要特征。如果程序中大多关系是一对多关系，或者多数记录间不存在关系，那么使用文档模型也是合适的。

虽然关系模型能够处理多对多关系的简单情况，但随着数据之间的连接变得复杂，数据建模为图形显得更加自然。

一个图由两种对象组成：顶点（vertices、node、entities），和边（edges、relationships、arcs）组成。多数数据可以被建模为一个图。一些典型例子：

\- 社交图谱

\- 网络图谱

\- 公路或铁路网络

\- 汽车导航系统最短路径。

图并不局限于同类数据，图还提供了一种一致方式，用在单个数据存储中存储完全不同类型的对象。如：Facebook维护一个包含许多不同类型的顶点和边的单个图：顶点表示人，地点，事件，签到和用户的评论；边缘表示哪些人是彼此的朋友，哪个签到发生在何处，谁评论了哪条消息，谁参与了哪个事件

两种不同但相关的方法构建和查询图表数据：图模型（由Neo4j，Titan和InfiniteGraph实现）、三元组存储模型（由Datomic，AllegroGraph等实现）。

四种声明式查询语言：Cypher，SPARQL和Datalog，Gremlin

**属性图**

每个顶点（vertex）包括：

\- 唯一标识

\- 一组出边（outgoing edges）

\- 一组入边（ingoing edges）

\- 一组属性（键值对）

每条边包括：

\- 唯一标识符

\- 边的起点/尾部顶点（tail vertex）

\- 边的终点/头部顶点（head vertex）

\- 描述两个顶点之间关系类型的标签

\- 一组属性（键值对）

可以将图存储看作由两个关系表组成：一个存储顶点，一个存储边。

一个例子用 PostgreSQL JSON 数据类型来存储每个顶点或每条边的属性。头部和尾部顶点用来存储每条边，如果想要一组顶点的输入或输出边，可以分别通过 head_vertex 或 tail_vertex 来查询 edges 表。

CREATE TABLE vertices (  vertex_id  INTEGER PRIMARY KEY,  properties JSON ); CREATE TABLE edges (  edge_id     INTEGER PRIMARY KEY,  tail_vertex INTEGER REFERENCES vertices (vertex_id),  head_vertex INTEGER REFERENCES vertices (vertex_id),  label       TEXT,  properties  JSON ); CREATE INDEX edges_tails ON edges (tail_vertex); CREATE INDEX edges_heads ON edges (head_vertex);

这个模型传达的重要信息：

\1. 任何顶点都有一边连接到任何其他顶点。没有模式限制哪种事物可不可以关联。

\2. 给定任何顶点，可以高效地找到它们的入边和出边，从而遍历图，即沿着一系列顶点的路径前后移动。

\3. 通过对不同类型的关系使用不同的标签，可以在一个图中存储集中不同信息，同时仍然保持清晰的数据模型。

**Cypher 查询语言**

Cypher 是属性图的声明式查询语言，为 Neo4j 图形数据库而发明。一个例子：

CREATE (NAmerica:Location {name:'North America', type:'continent'}), (USA:Location      {name:'United States', type:'country'  }), (Idaho:Location    {name:'Idaho',         type:'state'    }), (Lucy:Person       {name:'Lucy' }), (Idaho) -[:WITHIN]->  (USA)  -[:WITHIN]-> (NAmerica), (Lucy)  -[:BORN_IN]-> (Idaho)

这个例子展示了将数据字节表示为 Cypher 查询。每个顶点都有一个像`USA`或`Idaho`这样的符号名称，查询的其他部分可以使用这些名称在顶点之间创建边，使用箭头符号：`（Idaho） - [：WITHIN] ->（USA）`创建一条标记为`WITHIN`的边，`Idaho`为尾节点，`USA`为头节点。

当所有顶点和编都被添加进数据库后，有些问题了。如何找到符合条件的所有顶点，并返回这些顶点的 name 属性？有个关键信息可以把握：顶点拥有一条连到美国任一位置的`BORN_IN`边，和一条连到欧洲的任一位置的`LIVING_IN`边。具体例子如下：

MATCH (person) -[:BORN_IN]->  () -[:WITHIN*0..]-> (us:Location {name:'United States'}), (person) -[:LIVES_IN]-> () -[:WITHIN*0..]-> (eu:Location {name:'Europe'}) RETURN person.name

执行这个查询可能有几种可行的查询。首先会扫描数据库中所有人，检查每个人的出生地和居住地，然后返回符合条件的人。

通常对于声明式查询语言来说，在编写查询语句时，不需要指定执行细节：查询优化程序会自动选择预测效率最高的策略，因此你可以继续编写应用程序的其他部分。

**SQL 中的图查询**

建议在关系数据库中表示图数据。

自 SQL：1999，**查询可边长度遍历路径的思想可以称为：递归公用表达式**（WITH RECURSIVE 语法），一个例子：

WITH RECURSIVE  -- in_usa 包含所有的美国境内的位置ID    in_usa(vertex_id) AS (    SELECT vertex_id FROM vertices WHERE properties ->> 'name' = 'United States'    UNION    SELECT edges.tail_vertex FROM edges      JOIN in_usa ON edges.head_vertex = in_usa.vertex_id      WHERE edges.label = 'within'  ),  -- in_europe 包含所有的欧洲境内的位置ID    in_europe(vertex_id) AS (    SELECT vertex_id FROM vertices WHERE properties ->> 'name' = 'Europe'    UNION    SELECT edges.tail_vertex FROM edges      JOIN in_europe ON edges.head_vertex = in_europe.vertex_id      WHERE edges.label = 'within' ),   -- born_in_usa 包含了所有类型为Person，且出生在美国的顶点    born_in_usa(vertex_id) AS (      SELECT edges.tail_vertex FROM edges        JOIN in_usa ON edges.head_vertex = in_usa.vertex_id        WHERE edges.label = 'born_in' ),   -- lives_in_europe 包含了所有类型为Person，且居住在欧洲的顶点。    lives_in_europe(vertex_id) AS (      SELECT edges.tail_vertex FROM edges        JOIN in_europe ON edges.head_vertex = in_europe.vertex_id        WHERE edges.label = 'lives_in')   SELECT vertices.properties ->> 'name'  FROM vertices    JOIN born_in_usa ON vertices.vertex_id = born_in_usa.vertex_id    JOIN lives_in_europe ON vertices.vertex_id = lives_in_europe.vertex_id;

**三元组存储和 SPARQL**

和属性图模型相同，用不同的词来描述相同的想法。常见结构：主谓宾。

三元组的主语相当于图的顶点。一个 Turtle三元组 例子：

@prefix : <urn:example:>. _:lucy     a       :Person. _:lucy     :name   "Lucy". _:lucy     :bornIn _:idaho. _:idaho    a       :Location. _:idaho    :name   "Idaho". _:idaho    :type   "state". _:idaho    :within _:usa. _:usa      a       :Location _:usa      :name   "United States" _:usa      :type   "country". _:usa      :within _:namerica. _:namerica a       :Location _:namerica :name   "North America" _:namerica :type   :"continent"

一直重复相同主语有点麻烦，但可以用分号来说明统一主语的多个事情，使得可读性增强，例子如下：

@prefix : <urn:example:>. _:lucy      a :Person;   :name "Lucy";          :bornIn _:idaho. _:idaho     a :Location; :name "Idaho";         :type "state";   :within _:usa _:usa       a :Loaction; :name "United States"; :type "country"; :within _:namerica. _:namerica  a :Location; :name "North America"; :type "continent".

**语义网络**

三元数据模型完全独立于语义网络。如 Datomic 是三元组存储，没有声称和三元组有任何关系。

本质上讲，语义网络简单且合理：网站已经将信息发布为文字和图片让人们读，为什么不讲信息作为机器可读的数据也发布给计算机？资源描述框架（RDF）的目的是作为不同网站以一直的格式发布数据的一种机制，运行来自不同的数据自动合成一个数据网络。

**RDF 数据模型**

一个使用 Turtle 语言的 RDF 数据例子。有时候 RDF 也可以用 XML 格式编写，不过会有些啰嗦。

<rdf:RDF xmlns="urn:example:"          xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">     <Location rdf:nodeID="idaho">         <name>Idaho</name>         <type>state</type>         <within>             <Location rdf:nodeID="usa">                 <name>United States</name>                 <type>country</type>                 <within>                     <Location rdf:nodeID="namerica">                         <name>North America</name>                         <type>continent</type>                     </Location>                 </within>             </Location>         </within>     </Location>     <Person rdf:nodeID="lucy">         <name>Lucy</name>         <bornIn rdf:nodeID="idaho"/>     </Person> </rdf:RDF>

RDF 的奇怪在于，它是为了在互联网上交换数据而设计。三元组的主语，位于和并与通常是 URI。

**SPARQL 查询语言**

SPARQL是用于三元组存储的面向 RDF 数据模型的查询语言。类似于 Cypher

一个 SPARQL 例子，相比 Cypher 的代码更加简洁

PREFIX : <urn:example:> SELECT ?personName WHERE {  ?person :name ?personName.  ?person :bornIn  / :within* / :name "United States".  ?person :livesIn / :within* / :name "Europe". }

回顾代码，这两个表达式是等价的：

(person) -[:BORN_IN]-> () -[:WITHIN*0..]-> (location)   # Cypher ?person :bornIn / :within* ?location.                   # SPARQL

由于 RDF 不区分属性和边，而只将它们作为谓语，所以可以用相同语法来匹配属性。

(usa {name:'United States'})   # Cypher ?usa :name "United States".    # SPARQL

虽然 SPARQL 没有语义网，但仍然是很好的查询语言。

\> 小结对比

\1. CODASYL中，数据库有一个模式

\2. CODASYL中，达到特定记录的唯一方法是遍历其中的一个访问路径

\3. CODASYL，记录的后续是一个有序集合，所以数据库的人不得不维持排序,并且插入新记录到数据库的应用程序不得不担心的新记录在这些集合中的位置。

\4. CODASYL中，所有查询都是命令式的，难以编写，并且很容易因架构中的变化而受到破坏

**基础：Datalog**

Datalog 是比 SPARQL 和 Cypher 更古老的语言。

Datalog 类似三元组模式，但进行了泛化，把三元组变为了 谓语，主语，宾语。一个例子：

name(namerica, 'North America'). type(namerica, continent). name(usa, 'United States'). type(usa, country). within(usa, namerica). name(idaho, 'Idaho'). type(idaho, state). within(idaho, usa). name(lucy, 'Lucy'). born_in(lucy, idaho).

看起来有点不同于Cypher或SPARQL的等价物，但是**请不要放弃它**。

另一个相对应的查询例子。

within_recursive(Location, Name) :- name(Location, Name). /* Rule 1 */ within_recursive(Location, Name) :- within(Location, Via), /* Rule 2 */ 								within_recursive(Via, Name). migrated(Name, BornIn, LivingIn) :- name(Person, Name), /* Rule 3 */                                    born_in(Person, BornLoc),                                    within_recursive(BornLoc, BornIn),                                    lives_in(Person, LivingLoc),                                    within_recursive(LivingLoc, LivingIn). ?- migrated(Who, 'United States', 'Europe'). /* Who = 'Lucy'. */

定义规则，以将新谓语告诉数据库：在这里，定义了两个新的谓语，`within_recursive`和`migrated`。这些谓语不是存储在数据库中的三元组中，而是它们是从数据或其他规则派生而来的。规则可以引用其他规则，就像函数可以调用其他函数或者递归地调用自己一样。像这样，复杂的查询可以一次构建其中的一小块。

本章讨论的其他查询语言，需要采取不同的思维方式来思考Datalog方法，但却是一种非常强大的方法，因为规则可以在不同的查询中进行组合和重用。虽然对于简单的一次性查询不太方便，但可以更好地处理数据很复杂的情况。
