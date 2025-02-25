package main

import (
	"os"
	"strings"
	"sync"

	"github.com/NatsuiroGinga/mydocker/cgroups"
	"github.com/NatsuiroGinga/mydocker/cgroups/resource"
	"github.com/NatsuiroGinga/mydocker/container"
	"github.com/sirupsen/logrus"
)

// Run 执行具体 command
/*
这里的Start方法是真正开始前面创建好的command的调用，它首先会clone出来一个namespace隔离的
进程，然后在子进程中，调用/proc/self/exe,也就是调用自己，发送init参数，调用我们写的init方法，
去初始化容器的一些资源。
*/
func Run(tty bool, comArray []string, res *resource.ResourceConfig, containerName, imageName, volume string, envs []string) {
	var seed string
	if len(containerName) > 0 {
		seed = containerName
	} else {
		seed = imageName
	}

	containerId := container.GenerateContainerID(seed) // 生成容器 id
	logrus.Infof("containerID: %s", containerId)
	cmd, writePipe := container.NewParentProcess(tty, containerId, imageName, volume, envs)

	if cmd == nil {
		logrus.Errorf("new parent process error")
		return
	}
	if err := cmd.Start(); err != nil {
		logrus.Errorf("run parent.Start err:%v", err)
	}

	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	cgroupManager.Set(res)
	cgroupManager.Apply(cmd.Process.Pid)

	// 在子进程创建后通过管道来发送参数
	sendInitCommand(comArray, writePipe)

	if tty {
		cmd.Wait() // 前台运行，等待容器进程结束
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	// 然后创建一个 goroutine 来处理后台运行的清理工作
	go func() {
		defer wg.Done()

		if !tty {
			// 等待子进程退出
			_, _ = cmd.Process.Wait()
		}

		// 清理工作
		container.DeleteWorkSpace(containerId, volume)
		// container.DeleteContainerInfo(containerId)
		// if net != "" {
		// 	network.Disconnect(net, containerInfo)
		// }

		// 销毁 cgroup
		cgroupManager.Destroy()
	}()
	wg.Wait()
}

// sendInitCommand 通过writePipe将指令发送给子进程
func sendInitCommand(comArray []string, writePipe *os.File) {
	defer writePipe.Close()

	command := strings.Join(comArray, " ")
	logrus.Infof("all command is [%s]", command)

	_, err := writePipe.WriteString(command)
	if err != nil {
		logrus.Error(err)
	}
}
