package container

import (
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

// RunContainerInitProcess 启动容器的init进程
func RunContainerInitProcess(command string, args []string) error {
	logrus.Infof("command %s", command)

	// systemd 加入linux之后, mount namespace 就变成 shared by default,
	// 所以你必须显示声明你要这个新的mount namespace独立
	// 即 mount proc 之前先把所有挂载点的传播类型改为 private，避免本 namespace 中的挂载事件外泄。
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")

	// MS_NOEXEC 在本文件系统许运行其程序。
	// MS_NOSUID 在本系统中运行程序的时候， 允许 set-user-ID set-group-ID
	// MS_NOD 这个参数是自 Linux 2.4 ，所有 mount 的系统都会默认设定的参数。
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	_ = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf("%s", err.Error())
	}
	return nil
}
