package fs2

import "github.com/NatsuiroGinga/mydocker/cgroups/resource"

var Subsystems = []resource.Subsystem{
	&CpuSubSystem{},
	&MemorySubSystem{},
	&CpusetSubSystem{},
}
