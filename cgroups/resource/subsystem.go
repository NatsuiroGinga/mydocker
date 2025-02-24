package resource

// Subsystem 接口，每个Subsystem可以实现下面的4个接口，
// 这里将cgroup抽象成了path,原因是cgroup在hierarchy的路径，便是虚拟文件系统中的虚拟路径
type Subsystem interface {
	// Name 返回当前Subsystem的名称,比如cpu、memory
	Name() string

	// Set 设置某个cgroup在这个Subsystem中的资源限制
	//
	// 例如 memory subsystem 则需将配置写入 memory.limit_in_bytes 文件
	//
	// cpu subsystem 则是写入 cpu.cfs_period_us 和 cpu.cfs_quota_us
	Set(Path string, res *ResourceConfig) error

	// Apply 将进程添加到某个cgroup中
	Apply(path string, pid int) error

	// Remove 移除某个Cgroup
	Remove(path string) error
}
