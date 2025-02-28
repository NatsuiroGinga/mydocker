package network

import (
	"net"

	"github.com/vishvananda/netlink"
)

/*
Network 就是Linux中Bridge的抽象

网络(Netowrk)中可以有多个容器，在同一个网络里的容器可以通过这个网络互相通信。

就像挂载到同一个 Linux Bridge 设备上的网络设备一样， 可以直接通过 Bridge 设备实现网络互连;连接到同一个网络中的容器也可以通过这个网络和网络中别的容器互连。

网络中会包括这个网络相关的配置，比如网络的容器地址段、网络操作所调用的网络驱动等信息。
*/
type Network struct {
	Name    string     // 网络名
	IPRange *net.IPNet // 地址段
	Driver  string     // 网络驱动名
}

/*
Endpoint 是 Linux 中 Veth 的抽象
*/
type Endpoint struct {
	ID          string           `json:"id"`
	Device      netlink.Veth     `json:"dev"`
	IPAddress   net.IP           `json:"ip"`
	MacAddress  net.HardwareAddr `json:"mac"`
	Network     *Network
	PortMapping []string
}

/*
网络驱动(Network Driver) 是一个网络功能中的组件

不同的驱动对网络的创建、连接、销毁的策略不同

通过在创建网络时指定不同的网络驱动来定义使用哪个驱动做网络的配置。
*/
type Driver interface {
	Name() string
	Create(subnet string, name string) (*Network, error)
	Delete(network *Network) error
	Connect(networkName string, endpoint *Endpoint) error // 内部会修改 endpoint.Device，传指针
	Disconnect(endpointID string) error
}

/*
IPAM(IP Address Management) 也是网络功能中的一个组件

用于网络 IP 地址的分配和释放，包括容器的IP地址和网络网关的IP地址
*/
type IPAMer interface {
	Allocate(subnet *net.IPNet) (ip net.IP, err error) // 从指定的 subnet 网段中分配 IP 地址
	Release(subnet *net.IPNet, ipaddr *net.IP) error   // 从指定的 subnet 网段中释放掉指定的 IP 地址。
}
