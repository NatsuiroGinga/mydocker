package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/NatsuiroGinga/mydocker/cgroups"
	"github.com/NatsuiroGinga/mydocker/cgroups/resource"
	"github.com/NatsuiroGinga/mydocker/container"
	"github.com/NatsuiroGinga/mydocker/network"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Run 执行具体 command
/*
这里的Start方法是真正开始前面创建好的command的调用，它首先会clone出来一个namespace隔离的
进程，然后在子进程中，调用/proc/self/exe,也就是调用自己，发送init参数，调用我们写的init方法，
去初始化容器的一些资源。
*/
func Run(tty bool, comArray []string, res *resource.ResourceConfig, containerName, imageName, volume string, envs []string, net string, portMapping []string) {
	// 生成容器 id
	containerId := container.GenerateContainerID(imageName)

	logrus.Infof("containerID: %s", containerId)
	cmd, writePipe := container.NewParentProcess(tty, containerId, imageName, volume, envs)

	if cmd == nil {
		logrus.Errorf("new parent process error")
		return
	}
	// 启动子进程
	if err := cmd.Start(); err != nil {
		logrus.Errorf("run parent.Start err:%v", err)
	}

	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	cgroupManager.Set(res)
	cgroupManager.Apply(cmd.Process.Pid)
	var containerIP string
	// 如果指定了网络信息则进行配置
	if net != "" {
		// config container network
		containerInfo := &container.Info{
			Id:          containerId,
			Pid:         strconv.Itoa(cmd.Process.Pid),
			Name:        containerName,
			PortMapping: portMapping,
		}
		ip, err := network.Connect(net, containerInfo)
		if err != nil {
			log.Errorf("Error Connect Network %v", err)
			return
		}
		containerIP = ip.String()
	}
	
	// 记录容器信息， 写入/var/lib/mydocker/[containerId]/config.json中
	containerInfo, err := container.RecordContainerInfo(cmd.Process.Pid, comArray, containerName, containerId, volume, containerIP)
	if err != nil {
		logrus.Errorf("Record container info error %v", err)
		return
	}


	// 在子进程创建后通过管道来发送参数
	sendInitCommand(comArray, writePipe)

	if tty { // // 如果是tty，那么父进程等待，就是前台运行，否则就是跳过，实现后台运行
		cmd.Wait() // 前台运行，等待容器进程结束
	}

	// 然后创建一个 goroutine 来处理后台运行的清理工作
	go func() {
		if !tty {
			// 等待子进程退出
			_, _ = cmd.Process.Wait()
		}

		// 清理工作
		container.DeleteWorkSpace(containerId, volume)
		container.DeleteContainerInfo(containerId)
		if net != "" {
			network.Disconnect(net, containerInfo)
		}
		// 销毁 cgroup
		cgroupManager.Destroy()
	}()
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
