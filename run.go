package main

import (
	"os"
	"strings"

	"github.com/NatsuiroGinga/mydocker/container"
	"github.com/sirupsen/logrus"
)

// Run 执行具体 command
/*
这里的Start方法是真正开始前面创建好的command的调用，它首先会clone出来一个namespace隔离的
进程，然后在子进程中，调用/proc/self/exe,也就是调用自己，发送init参数，调用我们写的init方法，
去初始化容器的一些资源。
*/
func Run(tty bool, comArray []string) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Errorf("Run parent.Start err:%v", err)
	}
	// 在子进程创建后通过管道来发送参数
	sendInitCommand(comArray, writePipe)
	_ = parent.Wait()
	os.Exit(-1)
}

// sendInitCommand 通过writePipe将指令发送给子进程
func sendInitCommand(comArray []string, writePipe *os.File) {
	defer writePipe.Close()

	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)

	_, err := writePipe.WriteString(command)
	if err != nil {
		logrus.Error(err)
	}
}
