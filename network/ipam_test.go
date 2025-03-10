package network

import (
	"net"
	"testing"
)

func Test_Allocate(t *testing.T) {
	_, ipNet, _ := net.ParseCIDR("192.168.0.2/24")
	ip, err := ipAllocator.Allocate(ipNet)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("alloc ip: %v", ip)
}

func TestRelease(t *testing.T) {
	ip, ipNet, _ := net.ParseCIDR("192.168.0.1/24")
	err := ipAllocator.Release(ipNet, &ip)
	if err != nil {
		t.Fatal(err)
	}
}
