#!/bin/bash

# 1. 创建出Namespace网络隔离环境来模拟容器行为
ip netns add ns1
ip netns add ns2
ip netns show

# 2. 创建 Veth pairs
ip link add veth0 type veth peer name veth1
ip link add veth2 type veth peer name veth3

# 3. 查看
echo "ip link show"
ip link show

# 4. 将veth的一端放入容器内
ip link set veth0 netns ns1
ip link set veth2 netns ns2
echo "ip link show"
ip link show

# 5. 进入容器查看
echo "show ns1 links"
ip netns exec ns1 ip link show