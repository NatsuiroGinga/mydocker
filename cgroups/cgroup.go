package cgroups

import (
	"github.com/NatsuiroGinga/mydocker/cgroups/resource"
	"github.com/sirupsen/logrus"
)

// CgroupManager 来统一管理各个 subsystem。
type CgroupManager interface {
	// Apply 将进程pid加入到这个cgroup中
	Apply(pid int) error

	// Set 设置cgroup资源限制
	Set(res *resource.ResourceConfig) error

	// Destroy 释放cgroup
	Destroy() error
}

// path是cgroup在hierarchy中的路径 相当于创建的cgroup目录相对于root cgroup目录的路径
func NewCgroupManager(path string) CgroupManager {
	if IsCgroup2UnifiedMode() {
		logrus.Infof("use cgroup v2")
		return NewCgroupManagerV2(path)
	}
	logrus.Infof("use cgroup v1")
	return NewCgroupManagerV1(path)
}
