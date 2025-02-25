package main

import (
	"fmt"

	"github.com/NatsuiroGinga/mydocker/cgroups/resource"
	"github.com/NatsuiroGinga/mydocker/container"
	"github.com/urfave/cli"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit.
			mydocker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it", // 简单起见，这里把 -i 和 -t 参数合并成一个
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m", // 限制进程内存使用量
			Usage: "memory limit, e.g.: -m 100m",
		},
		cli.IntFlag{
			Name:  "cpu", // 限制cpu使用率
			Usage: "cpu quota, e.g.: -cpu 100",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpu limit, e.g.: -cpuset 2,4",
		},
		cli.StringFlag{ // 数据卷
			Name:  "v",
			Usage: "volume, e.g.: -v /etc/conf:/etc/conf",
		},
		cli.StringFlag{
			Name:  "name, ",
			Usage: "container name, e.g.: --name my_container",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container,run background",
		},
	},
	/*
		这里是run命令执行的真正函数。
		1.判断参数是否包含command
		2.获取用户指定的command
		3.调用Run function去准备启动容器:
	*/
	Action: func(context *cli.Context) error {
		if len(context.Args()) == 0 {
			return fmt.Errorf("missing container command")
		}

		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}

		imageName := cmdArray[0] // 镜像名称
		cmdArray = cmdArray[1:]

		tty := context.Bool("it")
		detach := context.Bool("d")

		if tty && detach {
			return fmt.Errorf("it and d flag can not both provided")
		}

		logrus.Debugf("detach: %v", detach)

		if !detach { // 如果不是指定后台运行，就默认前台运行
			tty = true
		}

		resConf := &resource.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuCfsQuota: context.Int("cpu"),
		}

		logrus.Infof("ResourceConfig: %#v", resConf)

		containerName := context.String("name")

		logrus.Infof("image name: %s", imageName)

		logrus.Infof("containerName: %s", containerName)

		volume := context.String("v")

		Run(tty, cmdArray, resConf, containerName, imageName, volume, nil)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside.",
	/*
		1.获取传递过来的 command 参数
		2.执行容器初始化操作
	*/
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit container to image",
	Action: cli.ActionFunc(func(ctx *cli.Context) error {
		if len(ctx.Args()) == 0 {
			return fmt.Errorf("missing image name")
		}

		containerID := ctx.Args().Get(0)
		imageName := ctx.Args().Get(1)

		return commitContainer(containerID, imageName)
	}),
}

var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list all the containers",
	Action: cli.ActionFunc(func(ctx *cli.Context) error {
		ListContainers()
		return nil
	}),
}
