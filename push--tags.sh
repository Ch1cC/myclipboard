
#!/bin/bash

# git ls-remote --tags origin | awk '{print $2}' | grep -v '\^{}$' | xargs -I{} git push origin --delete {}
# git tag -l | xargs git tag -d
git pull
git tag -l
echo "请输入你的Tag："
read NEW_TAG
git tag -a ${NEW_TAG} -m "build ${NEW_TAG}"

git pull
git push origin ${NEW_TAG}

echo "标签 ${NEW_TAG} 已创建并推送到远程仓库"
