#!/bin/bash
apt install bridge-utils

# 1. 创建bridge br0
brctl addbr br0

# 2. 将veth的另一端接入bridge
brctl addif br0 veth1
brctl addif br0 veth3

# 3. 查看接入效果
echo "brctl show"
brctl show