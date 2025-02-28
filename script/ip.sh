#!/bin/bash

# 1. 为bridge分配ip，激活上线
ip addr ad 172.18.0.1/24 dev br0
ip link set br0 up

# 2. 为容器内的网卡分配ip地址，并激活上线
ip netns exec ns1 ip addr add 172.18.0.2/24 dev veth0
ip netns exec ns1 ip link set veth0 up

ip netns exec ns2 ip addr add 172.18.0.3/24 dev veth2
ip netns exec ns2 ip link set veth2 up

# 3. veth另一端的网卡激活上线
ip link set veth1 up
ip link set veth3 up