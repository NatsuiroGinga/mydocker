package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/NatsuiroGinga/mydocker/constant"
	"github.com/sirupsen/logrus"
)

// volumeExtract 通过冒号分割解析volume目录，比如 -v /tmp:/tmp
func volumeExtract(volume string) (sourcePath, destinationPath string, err error) {
	parts := strings.Split(volume, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid volume [%s], must split by `:`", volume)
	}

	sourcePath, destinationPath = parts[0], parts[1]
	if sourcePath == "" || destinationPath == "" {
		return "", "", fmt.Errorf("invalid volume [%s], path can't be empty", volume)
	}

	return sourcePath, destinationPath, nil
}

/*
挂载数据卷的过程如下。

1）首先，创建宿主机文件目录

2）然后，拼接出容器目录在宿主机上的真正目录，格式为：$mntPath/$containerPath

因为之前使用了 pivotRoot 将$mntPath 作为容器 rootfs，因此这里的容器目录也可以按层级拼接最终找到在宿主机上的位置。

3）最后，执行 bind mount 操作，至此对数据卷的处理也就完成了。
*/
// mountVolume 使用 bind mount 挂载 volume
func mountVolume(mntPath, hostPath, containerPath string) {
	// 创建宿主机目录
	if err := os.Mkdir(hostPath, constant.Perm0777); err != nil {
		logrus.Infof("mkdir parent dir %s error. %v", hostPath, err)
	}
	// 拼接出对应的容器目录在宿主机上的的位置，并创建对应目录
	containerPathInHost := path.Join(mntPath, containerPath)
	if err := os.Mkdir(containerPathInHost, constant.Perm0777); err != nil {
		logrus.Infof("mkdir container dir %s error. %v", containerPathInHost, err)
	}
	// 通过bind mount 将宿主机目录挂载到容器目录
	// mount -o bind /hostPath /containerPath
	cmd := exec.Command("mount", "-o", "bind", hostPath, containerPathInHost)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("mount volume failed. %v", err)
	}
}

func umountVolume(mntPath, containerPath string) {
	// mntPath 为容器在宿主机上的挂载点，例如 /root/merged
	// containerPath 为 volume 在容器中对应的目录，例如 /root/tmp
	// containerPathInHost 则是容器中目录在宿主机上的具体位置，例如 /root/merged/root/tmp
	containerPathInHost := path.Join(mntPath, containerPath)
	cmd := exec.Command("umount", containerPathInHost)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Umount volume failed. %v", err)
	}
}
