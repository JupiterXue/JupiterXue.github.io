---
title: "沉浸式写作体验——Typora 插图"
date: 2021-11-18T10:16:09+08:00
category: tech
tags:
    - 工具
keywords:
    - 实践
---
最近一位也在写公众号的小伙伴问到公众号有没有什么便捷的方式插入图片，以及是更好的**沉浸式写作体验，边输入边搞定排版**。这个问题其实之前也有两个小伙伴问过我，我也只是丢给她们怎样使用的链接就“放心”了。没想到后来还会有其他小伙伴也可能用到，既然这个分享的动作也是重复的，那么我就来优化一下，动手写一篇手把手的教学。



虽然本文只是在说如何配置和使用 Markdown 编辑器 Typora，但经过下面略微复杂的操作后，可以做到一劳永逸。刚开始的时候有点麻烦，后来每次在 Typora 中写完文章后就能全选然后 Ctrl C 和 Ctrl V 复制粘贴到公众号，填上标题作者、选上配图打开原创，几秒就能发文章了。



# Typora 的安装与使用

## 安装

首先来介绍一下 Typora，打开官网，可以看到非常简洁的界面，找到右上角的下载，根据你的电脑系统平台选择相应的安装包。



![image-20211118230238536](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111182302695.png)

<center>图 1：Typora 下载页面</center>

## 使用

安装好之后，打开 Typora 可以看到一个非常清爽、没有任何干扰的写作界面。



![image-20211118231539419](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111182315518.png)

<center>图 2：Typora 打开的一篇文稿</center>

什么是 Markdown，图中也给出了一定的解释。其实 Markdown 也是专门给程序员而设计的语法，通过类似代码的方式来实现文字、排版的不同样式。并且 Typora 采用所见即所得(what you can see is what you can get) 的方式，写入的 Markdown 语法会被直接解析成相应的格式。如果你想看源代码，使用快捷键 Ctrl /，同样的，再次按下 Ctrl /，就能回到原样了。



![image-20211118231920536](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111182319605.png)

<center>图 3：Typora 源码模式</center>

## 公众号插图的解决方案

公众号可以直接插图片有三种方式，一中是 QQ 或微信截图，一种是网络图（本地如果上传到了云端就可以直接粘贴过去），还有一种是本地拖入图。它们分别有以下的局限：

- 截图能够直接粘贴，但图片大小和位置不能做改动（我是没找到）

- 本地拖入的方式虽然还勉强方便，但同样不能修改大小和位置。

- 而网络上的图或者从本地上传到微信的图，大小固定，更是不能修改大小。

有些图我们想**修改大小、加上边框、直接粘贴**，有没有什么好的办法呢？有的——**Typora+Github + PicGo + jsDelivr**。通过 PicGo 创建稳定图床 URL 链接，存放在免费的云端仓库 Github，再通过 jsDelivr 来加速访问 Github，最终**在 Typora 实现每插入一张图就自动上传到云端、生成可直接复制粘贴的云端图**。



网上有很多具体操作方案，就不重述了。具体配置方法请查看参考资料[《Github + PicGo + jsDelivr 创建稳定、免费图床》][2]，[《Typora如何通过PicGo自动上传图片到图床》][3]。



（如果你按照教程操作）需要注意的是，出于安全考虑，我们一般不会设置永久的 token 令牌去上传图片，没准哪天有人获得了，给我的仓库乱加图呢？因此我们可以设置一个中长期的 token，一般是 90 天。所以在 90 天后，如果你写文章发现图片自动上传失败，就需要来检查一下，顺便再熟悉一下如何操作整个 token 生成流程。不麻烦的。

## Typora 定制化方案

如果希望对文章的一二三多级标题的字体颜色、字体大小以及字体样式做改动，就需要我们有一点前端基础了。（不感兴趣的读者请跳过）这里以 vue 的 css 为例，首先打开文件。

![image-20211118234239692](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111182342778.png)

<center>图 4：打开 Typora 偏好设置，再打开主题文件夹</center>

通过记事本 notepad++ 打开**主题名+.css**，我这里是 vue.css，就可以看到整个主题的样式了。现在我想修改一级标题为橙色、带有居中、带有下划线的标题，二级标题也差不多，就可以这么改动 css：

![image-20211118234849813](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111182348908.png)

<center>图 5：vue.css 自定义样式</center>

最后，记得按 Ctrl S 保存，然后把所有 Typora 关闭，重新打开，您定制的精美排版样式成型了。



以上，我们就完成了图片手动上传、精美样式等完美功能。现在写文章不管是网上图，还是本地图，还是截图，直接放进 Typora 中，它会自动转换为网络图（记得勾选自动上传图片），就可以直接轻松加愉快地发公众号啦！！

![image-20211118235620967](https://cdn.jsdelivr.net/gh/JupiterXue/PictureBed/BlogImg/202111182356058.png)

<center>图 6：Typora 自动上传图片</center>

还有一系列功能都在 Typora 中已经实现，等待你去探索，如目录大纲、代码块、引用块、数学公式、流程图、表格、文章超链接引用等。




如果你对以上提到的功能感兴趣欢迎留言交流，如果本文对你有帮助也欢迎给小编加鸡腿哦。



__参考资料__
[1]: https://www.zybuluo.com/mdeditor	"Markdown 教程，非常友好，左边是源码，右边是映射"

[2]: https://www.jianshu.com/p/4a29945ad69c	"Github + PicGo + jsDelivr 创建稳定、免费图床"
[3]: https://www.jianshu.com/p/a3371dc5802e	"Typora如何通过PicGo自动上传图片到图床"
