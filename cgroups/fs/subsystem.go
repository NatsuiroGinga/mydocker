package fs

import "github.com/NatsuiroGinga/mydocker/cgroups/resource"

// SubsystemsIns 通过不同的subsystem初始化实例创建资源限制处理链数组
var SubsystemsIns = []resource.Subsystem{
	&CpusetSubSystem{},
	&MemorySubSystem{},
	&CPUSubsystem{},
}
