package container

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/NatsuiroGinga/mydocker/utils"
	"github.com/sirupsen/logrus"
)

// RunContainerInitProcess 启动容器的init进程
func RunContainerInitProcess(command string, args []string) error {
	logrus.Infof("command %s", command)
	mountProc()
	cmdArray := readUserCommand()
	if len(cmdArray) == 0 {
		return errors.New("run container get user command error, cmdArray is nil")
	}

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Errorf("Exec loop path error %v", err)
		return err
	}

	logrus.Infof("Find path %s", path)

	if err := syscall.Exec(path, cmdArray, os.Environ()); err != nil {
		logrus.Errorf("RunContainerInitProcess exec :" + err.Error())
	}
	return nil
}

const fdIndex = 3

// 子进程读数据, 子进程启动后，首先要找到前面通过ExtraFiles 传递过来的 readPipe FD，然后才是数据读取
//
// 1）获取 readPipe FD
//
// 2）读取数据
func readUserCommand() []string {

	/*
		uintptr(3）就是指 index 为3的文件描述符，也就是传递进来的管道的另一端，至于为什么是3，具体解释如下：
		因为每个进程默认都会有3个文件描述符，分别是标准输入、标准输出、标准错误。这3个是子进程一创建的时候就会默认带着的，
		前面通过ExtraFiles方式带过来的 readPipe 理所当然地就成为了第4个。
		在进程中可以通过index方式读取对应的文件，比如
		index0：标准输入
		index1：标准输出
		index2：标准错误
		index3：带过来的第一个FD，也就是readPipe
		由于可以带多个FD过来，所以这里的3就不是固定的了。
		比如像这样：cmd.ExtraFiles = []*os.File{a,b,c,readPipe} 这里带了4个文件过来，分别的index就是3,4,5,6
		那么我们的 readPipe 就是 index6,读取时就要像这样：pipe := os.NewFile(uintptr(6), "pipe")
	*/
	pipe := os.NewFile(uintptr(fdIndex), "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("init read pipe error %v", err)
		return nil
	}

	msgStr := utils.Bytes2String(msg)
	return strings.Split(msgStr, " ")
}

// mountProc 在容器内 mount /proc 文件系统
func mountProc() {
	// systemd 加入linux之后, mount namespace 就变成 shared by default,
	// 所以你必须显示声明你要这个新的mount namespace独立
	// 即 mount proc 之前先把所有挂载点的传播类型改为 private，避免本 namespace 中的挂载事件外泄。
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")

	// MS_NOEXEC 在本文件系统许运行其程序。
	// MS_NOSUID 在本系统中运行程序的时候， 允许 set-user-ID set-group-ID
	// MS_NOD 这个参数是自 Linux 2.4 ，所有 mount 的系统都会默认设定的参数。
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	_ = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
}
