package main

import (
	"github.com/NatsuiroGinga/mydocker/container"
	log "github.com/sirupsen/logrus"
)

/*
removeContainer 则是 rm 命令的真正实现，根据 Id 拿到容器信息，然后先判断状态:

# STOP 状态，则直接删除

# RUNNING 状态，如果带了 force flag 则先 Stop 然后再删除，否则打印错误信息
*/
func removeContainer(containerId string, force bool) {
	containerInfo, err := container.GetContainerInfoById(containerId)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerId, err)
		return
	}

	switch containerInfo.Status {
	case container.STOP: // STOP状态的容器可以直接删除
		container.DeleteContainerInfo(containerId)
		container.DeleteWorkSpace(containerId, containerInfo.Volume)
	case container.RUNNING: // RUNNING容器如果指定了force则先stop再删除
		if !force {
			log.Errorf("Couldn't remove running container [%s], Stop the container before attempting removal or"+
				" force remove", containerId)
			return
		}
		stopContainer(containerId)
		removeContainer(containerId, force)
	default:
		log.Errorf("Couldn't remove container,invalid status %s", containerInfo.Status)
		return
	}
}
