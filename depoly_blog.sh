#!/bin/bash

myblog_dir=/opt/myblog      #博客网站根目录
public=/opt/myblog/public   #生成的git静态网页
msg="build site `date`"         #git commit的备注信息

if [ -e $public ];then
        cd $myblog_dir
        rm -rf $public  #删除原有静态文件
        hugo --theme=m10c --baseUrl="https://YourName.github.io/" --buildDrafts     #指定主题编译成静态文件，存放在public
        cd $public
        git init
        git add .

        if [ $# -eq 1 ];then
          msg = "$1"
        fi
        git commit -m "$msg"
        git remote add origin https://github.com/YourName/YourName.github.io.git	#添加仓库源
        git push -u origin master  #推送到GitHub
        if [ $? -eq 0 ];then
                echo "文章已更新至github！"
        fi
fi
