package fs

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/NatsuiroGinga/mydocker/constant"
	"github.com/sirupsen/logrus"
)

const mountPointIndex = 4

// getCgroupPath 找到cgroup在文件系统中的绝对路径
/*
实际就是将根目录和cgroup名称拼接成一个路径。
如果指定了自动创建，就先检测一下是否存在，如果对应的目录不存在，则说明cgroup不存在，这里就给创建一个
*/
func getCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	// 不需要自动创建就直接返回
	cgroupRoot := findCgroupMountPoint(subsystem)
	absPath := path.Join(cgroupRoot, cgroupPath)
	if !autoCreate {
		return absPath, nil
	}

	// 指定自动创建时判断是否存在
	_, err := os.Stat(absPath)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		err = os.Mkdir(absPath, constant.Perm0755)
		return absPath, err
	}

	// 其他错误或者没有错误都直接返回，如果err=nil,那么errors.Wrap(err, "")也会是nil
	return absPath, errors.Join(err, errors.New("create cgroup"))
}

// findCgroupMountPoint 通过/proc/self/mountinfo找出挂载了某个subsystem的hierarchy cgroup根节点所在的目录
func findCgroupMountPoint(subsystem string) string {
	// /proc/self/mountinfo 为当前进程的 mountinfo 信息
	// 可以直接通过 cat /proc/self/mountinfo 命令查看
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	// 这里主要根据各种字符串处理来找到目标位置
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// txt 大概是这样的：104 85 0:20 / /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime - cgroup cgroup rw,memory
		txt := scanner.Text()
		// 然后按照空格分割
		fields := strings.Split(txt, " ")
		// 对最后一个元素按逗号进行分割，这里的最后一个元素就是 rw,memory
		// 其中的的 memory 就表示这是一个 memory subsystem
		subsystems := strings.Split(fields[len(fields)-1], ",")
		for _, opt := range subsystems {
			if opt == subsystem {
				// 如果等于指定的 subsystem，那么就返回这个挂载点跟目录，就是第四个元素，
				// 这里就是`/sys/fs/cgroup/memory`,即我们要找的根目录
				return fields[mountPointIndex]
			}
		}
	}

	if err = scanner.Err(); err != nil {
		logrus.Error("read err:", err)
		return ""
	}
	return ""
}
