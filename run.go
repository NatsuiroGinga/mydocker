package main

import (
	"os"

	"github.com/NatsuiroGinga/mydocker/container"
	"github.com/sirupsen/logrus"
)

// Run 执行具体 command
/*
这里的Start方法是真正开始前面创建好的command的调用，它首先会clone出来一个namespace隔离的
进程，然后在子进程中，调用/proc/self/exe,也就是调用自己，发送init参数，调用我们写的init方法，
去初始化容器的一些资源。
*/
func Run(tty bool, cmd string) {
	parent := container.NewParentProcess(tty, cmd)
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}
	_ = parent.Wait()
	os.Exit(-1)
}
