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

type MemorySubSystem struct {
}

// Name 返回cgroup名字
func (s *MemorySubSystem) Name() string {
	return "memory"
}

// Set 设置cgroupPath对应的cgroup的内存资源限制
func (s *MemorySubSystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}

	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, true)
	if err != nil {
		return err
	}

	// 设置这个cgroup的内存限制，即将限制写入到cgroup对应目录的memory.limit_in_bytes 文件中。
	if err := os.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup memory fail %v", err)
	}
	return nil
}

// Apply 将pid加入到cgroupPath对应的cgroup中
func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
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

// Remove 删除cgroupPath对应的cgroup
func (s *MemorySubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsysCgroupPath)
}
