package container

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/NatsuiroGinga/mydocker/constant"
)

// RecordContainerInfo 记录容器信息, 实现 docker ps 命令
/*
容器创建后，所有需要的信息都被存储到/var/lib/mydocker/containers/{containerID}下，
下面就可以通过读取并遍历这个目录下的容器去实现 mydocker ps 命令了。
*/
func RecordContainerInfo(
	containerPid int,
	commandArray []string,
	containerName,
	containerId string) (*Info, error) {

	if len(containerName) == 0 {
		containerName = containerId
	}

	command := strings.Join(commandArray, "")

	containerInfo := &Info{
		Id:          containerId,
		Pid:         strconv.Itoa(containerPid),
		Command:     command,
		CreatedTime: time.Now().Format(time.DateTime),
		Status:      RUNNING,
		Name:        containerName,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		return containerInfo, errors.Join(err, errors.New("container info marshal failed"))
	}
	jsonStr := string(jsonBytes)

	// 拼接出存储容器信息文件的路径, 如果目录不存在则级联创建
	dirPath := fmt.Sprintf(InfoLocFormat, containerId)
	if err = os.MkdirAll(dirPath, constant.Perm0622); err != nil {
		return containerInfo, errors.Join(err, fmt.Errorf("mkdir %s failed", dirPath))
	}
	// 将容器信息写入文件
	filename := path.Join(dirPath, ConfigName)
	file, err := os.Create(filename)
	if err != nil {
		return containerInfo, errors.Join(err, fmt.Errorf("create file %s failed", filename))
	}
	defer file.Close()

	if _, err := file.WriteString(jsonStr); err != nil {
		return containerInfo, errors.Join(err, fmt.Errorf("write container info to  file %s failed", filename))
	}

	return containerInfo, nil
}

func DeleteContainerInfo(containerID string) error {
	dirPath := fmt.Sprintf(InfoLocFormat, containerID)
	if err := os.RemoveAll(dirPath); err != nil {
		return errors.Join(err, fmt.Errorf("remove dir %s failed", dirPath))
	}
	return nil
}
