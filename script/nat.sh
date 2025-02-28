#!/bin/bash

# 1. 配置容器内路由
# 将bridge设置为“容器”的缺省网关。让非172.18.0.0/24网段的数据包都路由给bridge，这样数据就从“容器”跑到宿主机上来了。
sudo ip netns exec ns1 ip route add default via 172.18.0.1 dev veth0
sudo ip netns exec ns2 ip route add default via 172.18.0.1 dev veth2
echo "ip netns exec ns1 ip route"
ip netns exec ns1 ip route

# 2. 宿主机开启转发功能并配置转发规则
# 在宿主机上配置内核参数，允许IP forwarding，这样才能把网络包转发出去。
sudo sysctl net.ipv4.conf.all.forwarding=1
# 还有就是要配置 iptables FORWARD 规则
iptables -t filter -L FORWARD
iptables -t nat -A POSTROUTING -s 172.18.0.0/24 ! -o br0 -j MASQUERADE

# 3. 外部访问容器需要进行 DNAT，把目的IP地址从宿主机地址转换成容器地址。
# 在nat表的PREROUTING链增加规则，当输入设备不是br0，目的端口为80 时，做目的地址转换，将宿主机IP替换为容器IP。
sudo iptables -t nat -A PREROUTING ! -i br0 -p tcp -m tcp --dport 80 -j DNAT --to-destination 172.18.0.2:80
sudo iptables -t nat -A OUTPUT -p tcp -m tcp --dport 80 -j DNAT --to-destination 172.18.0.2:80
