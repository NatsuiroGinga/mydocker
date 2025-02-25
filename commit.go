package main

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/NatsuiroGinga/mydocker/utils"
	"github.com/sirupsen/logrus"
)

var ErrImageAlreadyExists = errors.New("image already exists")

func commitContainer(containerID string, imageName string) error {
	mntPath := utils.GetMerged(containerID)
	if len(imageName) == 0 {
		imageName = containerID
	}
	imageTar := utils.GetImage(imageName)
	exist, err := utils.PathExists(imageTar)

	if err != nil {
		return errors.Join(err, fmt.Errorf("check is image [%s/%s] exist failed", imageName, imageTar))
	}

	if exist {
		return ErrImageAlreadyExists
	}

	logrus.Infof("commitContainer imageTar:%s", imageTar)

	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntPath, ".").CombinedOutput(); err != nil {
		return errors.Join(err, fmt.Errorf("tar folder %s failed", mntPath))
	}

	return nil
}
