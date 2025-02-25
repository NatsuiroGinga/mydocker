package container

import (
	"os"
	"os/exec"

	"github.com/NatsuiroGinga/mydocker/utils"
	"github.com/sirupsen/logrus"
)

// NewWorkSpace Create an Overlay2 filesystem as container root workspace
/*
1）创建lower层
2）创建upper、worker层
3）创建merged目录并挂载overlayFS
4）如果有指定volume则挂载volume // TODO
*/
func NewWorkSpace(containerID, imageName, volume string) {
	createLower(containerID, imageName)
	createDirs(containerID)
	mountOverlayFS(containerID)

	/*
		在原有创建过程最后增加 volume bind 逻辑：

		1）首先判断 volume 是否为空，如果为空，就表示用户并没有使用挂载参数，不做任何处理
		2）如果不为空，则使用 volumeUrlExtract 函数解析 volume 字符串，得到要挂载的宿主机目录和容器目录，并执行 bind mount
	*/
	// 如果指定了volume则还需要mount volume
	if volume != "" {
		mntPath := utils.GetMerged(containerID)
		hostPath, containerPath, err := volumeExtract(volume)
		if err != nil {
			logrus.Errorf("extract volume failed，maybe volume parameter input is not correct，detail:%v", err)
			return
		}
		mountVolume(mntPath, hostPath, containerPath)
	}
}

// DeleteWorkSpace Delete the UFS filesystem while container exit
/*
和创建相反
1）有volume则卸载volume

删除容器文件系统时，先判断是否挂载了 volume，如果挂载了则删除时则需要先 umount volume。

注意：一定要要先 umount volume ，然后再删除目录，否则由于 bind mount 存在，删除临时目录会导致 volume 目录中的数据丢失。

2）卸载并移除merged目录
3）卸载并移除upper、worker层
*/
func DeleteWorkSpace(containerID, volume string) {
	// 如果指定了volume则需要umount volume
	// NOTE: 一定要要先 umount volume ，然后再删除目录，否则由于 bind mount 存在，删除临时目录会导致 volume 目录中的数据丢失。
	if volume != "" {
		_, containerPath, err := volumeExtract(volume)
		if err != nil {
			logrus.Errorf("extract volume failed，maybe volume parameter input is not correct，detail:%v", err)
			return
		}
		mntPath := utils.GetMerged(containerID)
		umountVolume(mntPath, containerPath)
	}

	umountOverlayFS(containerID)
	deleteDirs(containerID)
}

func umountOverlayFS(containerID string) {
	mntPath := utils.GetMerged(containerID)
	cmd := exec.Command("umount", mntPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logrus.Infof("umountOverlayFS,cmd:%v", cmd.String())
	if err := cmd.Run(); err != nil {
		logrus.Errorf("%v", err)
	}
	logrus.Infof("umount overlayfs [%s] success", mntPath)
}

func deleteDirs(containerID string) {
	dirs := []string{
		utils.GetMerged(containerID),
		utils.GetUpper(containerID),
		utils.GetWorker(containerID),
		utils.GetLower(containerID),
		utils.GetRoot(containerID), // root 目录也要删除
	}

	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			logrus.Errorf("Remove dir %s error %v", dir, err)
		}
	}

	logrus.Info("delete dirs success")
}

// createLower 根据 containerID, imageName 准备 lower 层目录
func createLower(containerID, imageName string) {
	// 根据 containerID 拼接出 lower 目录
	// 根据 imageName 找到镜像 tar，并解压到 lower 目录中
	lowerPath := utils.GetLower(containerID)
	imagePath := utils.GetImage(imageName)
	logrus.Infof("lower:%s image.tar:%s", lowerPath, imagePath)
	// 检查目录是否已经存在
	exist, err := utils.PathExists(lowerPath)
	if err != nil {
		logrus.Infof("Fail to judge whether dir %s exists. %v", lowerPath, err)
	}
	// 不存在则创建目录并将image.tar解压到lower文件夹中
	if !exist {
		if err = os.MkdirAll(lowerPath, 0777); err != nil {
			logrus.Errorf("Mkdir dir %s error. %v", lowerPath, err)
		}
		if _, err = exec.Command("tar", "-xvf", imagePath, "-C", lowerPath).CombinedOutput(); err != nil {
			logrus.Errorf("Untar dir %s error %v", lowerPath, err)
		}
	}
}

// createDirs 创建overlayfs需要的的merged、upper、worker目录
func createDirs(containerID string) {
	dirs := []string{
		utils.GetMerged(containerID),
		utils.GetUpper(containerID),
		utils.GetWorker(containerID),
	}

	for _, dir := range dirs {
		if err := os.Mkdir(dir, 0777); err != nil {
			logrus.Errorf("mkdir dir %s error. %v", dir, err)
		}
	}
}

// mountOverlayFS 挂载overlayfs
func mountOverlayFS(containerID string) {
	// 拼接参数
	// e.g. lowerdir=/root/busybox,upperdir=/root/upper,workdir=/root/work
	dirs := utils.GetOverlayFSDirs(utils.GetLower(containerID), utils.GetUpper(containerID), utils.GetWorker(containerID))
	mergedPath := utils.GetMerged(containerID)
	//完整命令：mount -t overlay overlay -o lowerdir=/root/{containerID}/lower,upperdir=/root/{containerID}/upper,workdir=/root/{containerID}/work /root/{containerID}/merged
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mergedPath)
	logrus.Infof("mount overlayfs: [%s]", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("%v", err)
	}
}
