# 执行一个交互式命令，让容器能一直后台运行
docker run -d busybox top
# 拿到刚创建的容器的 Id
containerId=$(docker ps --filter "ancestor=busybox:latest"|grep -v IMAGE|awk '{print $1}')
echo "containerId" $containerId
# export 从容器导出
docker export -o busybox.tar $containerId
# 最后将tar包解压
# mkdir busybox
# tar -xvf busybox.tar -C busybox/