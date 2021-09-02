echo 编译测试
hugo
if [ $? -ne 0 ]; then
    echo Failed !!!
else
    echo Successfully
fi

git add .
git commit -m 更新博客
git push
