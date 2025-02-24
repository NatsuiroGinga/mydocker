package cgroups

import (
	"github.com/NatsuiroGinga/mydocker/cgroups/fs"
	"github.com/NatsuiroGinga/mydocker/cgroups/resource"
	"github.com/sirupsen/logrus"
)

type CgroupManagerV1 struct {
	// cgroup在hierarchy中的路径 相当于创建的cgroup目录相对于root cgroup目录的路径
	Path string
	// 资源配置
	res        *resource.ResourceConfig
	Subsystems []resource.Subsystem
}

func NewCgroupManagerV1(path string) *CgroupManagerV1 {
	return &CgroupManagerV1{
		Path:       path,
		Subsystems: fs.SubsystemsIns,
	}
}

func (manager *CgroupManagerV1) Apply(pid int) error {
	for _, sys := range manager.Subsystems {
		err := sys.Apply(manager.Path, pid)
		if err != nil {
			logrus.Errorf("apply subsystem:%s err:%s", sys.Name(), err)
		}
	}
	return nil
}

func (manager *CgroupManagerV1) Set(res *resource.ResourceConfig) error {
	for _, sys := range manager.Subsystems {
		err := sys.Set(manager.Path, res)
		if err != nil {
			logrus.Errorf("apply subsystem:%s err:%s", sys.Name(), err)
		}
	}
	return nil
}

func (manager *CgroupManagerV1) Destroy() error {
	for _, sys := range manager.Subsystems {
		err := sys.Remove(manager.Path)
		if err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
