package cgroups

import (
	"errors"

	"github.com/NatsuiroGinga/mydocker/cgroups/fs2"
	"github.com/NatsuiroGinga/mydocker/cgroups/resource"
	"github.com/sirupsen/logrus"
)

type CgroupManagerV2 struct {
	Path       string
	Resource   *resource.ResourceConfig
	Subsystems []resource.Subsystem
}

func NewCgroupManagerV2(path string) *CgroupManagerV2 {
	return &CgroupManagerV2{
		Path:       path,
		Subsystems: fs2.Subsystems,
	}
}

// Apply 将进程pid加入到这个cgroup中
func (manager *CgroupManagerV2) Apply(pid int) error {
	if len(manager.Subsystems) > 0 {
		if err := manager.Subsystems[0].Apply(manager.Path, pid); err != nil {
			logrus.Errorf("apply pid [%d] to cgroup [%s] err:%s", pid, manager.Path, err)
		}
	}
	return nil
}

// Set 设置cgroup资源限制
func (manager *CgroupManagerV2) Set(res *resource.ResourceConfig) error {
	for _, subSysIns := range manager.Subsystems {
		err := subSysIns.Set(manager.Path, res)
		if err != nil {
			logrus.Errorf("apply subsystem:%s err:%s", subSysIns.Name(), err)
		}
	}
	return nil
}

// Destroy 释放cgroup
func (manager *CgroupManagerV2) Destroy() error {
	if len(manager.Subsystems) > 0 {
		manager.Subsystems[0].Remove(manager.Path)
		logrus.Infof("remove cgroup [%s] success", manager.Path)
	}
	return errors.New("fail to destroy cgroup")
}
