mkdir cgroup-test
mount -t cgroup -o cpu cgroup-test ./cgroup-test
cd cgroup-test || exit
ls
