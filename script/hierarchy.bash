# 挂载

# 1 挂载一个和 cpuset subsystem 关联的 hierarchy 到 ./cg1 目录
# 首先肯定是创建对应目录
mkdir cgroup-test
# 具体挂载操作--参数含义如下
# -t cgroup 表示操作的是 cgroup 类型，
# -o cpuset 表示要关联 cpuset subsystem，可以写0个或多个，0个则是关联全部subsystem，
# cg1 为 cgroup 的名字，
# ./cg1 为挂载目标目录。
mount -t cgroup -o cpu cgroup-test ./cgroup-test

# 2 挂载一颗和所有subsystem关联的cgroup树到cg1目录
# mkdir cg1
# mount -t cgroup cg1 ./cg1

# 3 挂载一颗与cpu和cpuacct subsystem关联的cgroup树到 cg1 目录
# mkdir cg1
# mount -t cgroup -o cpu,cpuacct cg1 ./cg1

# 4 挂载一棵cgroup树，但不关联任何subsystem，这systemd所用到的方式
# mkdir cg1
# mount -t cgroup -o none,name=cg1 cg1 ./cg1

# 删除
# 指定路径来卸载，而不是名字。
# $ umount /path/to/your/hierarchy
# 例如
# umount /sys/fs/cgroup/hierarchy
