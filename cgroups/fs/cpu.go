package fs

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/NatsuiroGinga/mydocker/cgroups/resource"
	"github.com/NatsuiroGinga/mydocker/constant"
)

const (
	PeriodDefault = 100000
	Percent       = 100
)

type CPUSubsystem struct {
}

func (sys *CPUSubsystem) Name() string {
	return "cpu"
}

func (sys *CPUSubsystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if len(cgroupPath) == 0 {
		return nil
	}
	subsysCgroupPath, err := getCgroupPath(sys.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	// cpu.cfs_period_us & cpu.cfs_quota_us 控制的是CPU使用时间，单位是微秒，比如每1秒钟，这个进程只能使用200ms，相当于只能用20%的CPU
	if res.CpuCfsQuota != 0 {
		// cpu.cfs_period_us 默认为100000，即100ms
		if err := os.WriteFile(path.Join(subsysCgroupPath, "cpu.cfs_period_us"), []byte(strconv.Itoa(PeriodDefault)), constant.Perm0644); err != nil {
			return fmt.Errorf("set cgroup cpu share fail %v", err)
		}
		// cpu.cfs_quota_us 则根据用户传递的参数来控制，比如参数为20，就是限制为20%CPU，所以把cpu.cfs_quota_us设置为cpu.cfs_period_us的20%就行
		// 这里只是简单的计算了下，并没有处理一些特殊情况，比如负数什么的
		if err = os.WriteFile(path.Join(subsysCgroupPath, "cpu.cfs_quota_us"), []byte(strconv.Itoa(PeriodDefault/Percent*res.CpuCfsQuota)), constant.Perm0644); err != nil {
			return fmt.Errorf("set cgroup cpu share fail %v", err)
		}
	}

	return nil
}

func (s *CPUSubsystem) Apply(cgroupPath string, pid int) error {
	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		return errors.Join(err, fmt.Errorf("get cgroup %s", cgroupPath))
	}

	// 打开cgroup的tasks文件，追加模式
	tasks, err := os.OpenFile(
		path.Join(subsysCgroupPath, "tasks"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		constant.Perm0644,
	)
	if err != nil {
		return fmt.Errorf("open cgroup tasks file failed: %v", err)
	}
	defer tasks.Close()

	// 将pid转换为字符串并写入tasks文件
	pidStr := strconv.Itoa(pid) + "\n"
	if _, err := tasks.Write([]byte(pidStr)); err != nil {
		return fmt.Errorf("append pid to cgroup tasks file failed: %v", err)
	}

	return nil
}

func (s *CPUSubsystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsysCgroupPath)
}
