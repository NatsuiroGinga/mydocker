package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"syscall"

	"github.com/NatsuiroGinga/mydocker/constant"
	"github.com/NatsuiroGinga/mydocker/container"
	log "github.com/sirupsen/logrus"
)

/*
stopContainer 负责停止容器的运行

主要分三步：

1.首先根据 ContainerId 找到之前记录的容器信息的文件并拿到容器具体信息，主要是 PID

2.然后调用 Kill 命令，给指定 PID 发送 SIGTERM

3.最后更新容器状态为 stop 并写回记录容器信息的文件.
*/
func stopContainer(containerId string) {
	// 1. 根据containerId查询容器信息
	containerInfo, err := container.GetContainerInfoById(containerId)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerId, err)
		return
	}
	pidInt, err := strconv.Atoi(containerInfo.Pid)
	if err != nil {
		log.Errorf("Conver pid from string to int error %v", err)
		return
	}
	// 2. 发送SIGTERM信号
	if err = syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		log.Errorf("Stop container %s error %v", containerId, err)
		return
	}
	// 3. 修改容器信息，将容器置为STOP状态，并清空PID
	containerInfo.Status = container.STOP
	containerInfo.Pid = ""
	contentBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Json marshal %s error %v", containerId, err)
		return
	}
	// 4. 重新写回存储容器信息的文件
	dirPath := fmt.Sprintf(container.InfoLocFormat, containerId)
	configFilePath := path.Join(dirPath, container.ConfigName)
	if err := os.WriteFile(configFilePath, contentBytes, constant.Perm0622); err != nil {
		log.Errorf("Write file %s error:%v", configFilePath, err)
	}
}
