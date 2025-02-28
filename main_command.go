package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/NatsuiroGinga/mydocker/cgroups/resource"
	"github.com/NatsuiroGinga/mydocker/container"
	"github.com/NatsuiroGinga/mydocker/network"
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
		cli.StringSliceFlag{ // 增加 -e flag
			Name:  "e",
			Usage: "set environment,e.g. -e name=mydocker",
		},
		cli.StringFlag{
			Name:  "net",
			Usage: "container network，e.g. -net testbr",
		},
		cli.StringSliceFlag{
			Name:  "p",
			Usage: "port mapping,e.g. -p 8080:80 -p 30336:3306",
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
		envs := context.StringSlice("e")

		network := context.String("net")
		portMapping := context.StringSlice("p")

		Run(tty, cmdArray, resConf, containerName, imageName, volume, envs, network, portMapping)
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

		containerName := ctx.Args().Get(0)
		imageName := ctx.Args().Get(1)

		return commitContainer(containerName, imageName)
	}),
}

var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list all the containers",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "a",
			Usage: "show all containers, e.g.: mydocker ps -a",
		},
	},
	Action: cli.ActionFunc(func(ctx *cli.Context) error {
		all := ctx.Bool("a")
		ListContainers(all)
		return nil
	}),
}

var logCommand = cli.Command{
	Name:  "logs",
	Usage: "print logs of a container",
	Action: cli.ActionFunc(func(ctx *cli.Context) error {
		if len(ctx.Args()) == 0 {
			return fmt.Errorf("please input your container name")
		}
		containerName := ctx.Args().Get(0)
		logContainer(containerName)
		return nil
	}),
}

var execCommand = cli.Command{
	Name:  "exec",
	Usage: "exec a command into container, e.g.: mydocker exec 123456789 /bin/sh",
	Action: cli.ActionFunc(func(ctx *cli.Context) error {
		// 如果环境变量存在，说明C代码已经运行过了，即setns系统调用已经执行了，这里就直接返回，避免重复执行
		if os.Getenv(EnvExecPid) != "" {
			log.Infof("pid callback pid %v", os.Getgid())
			return nil
		}
		// 格式：mydocker exec 容器名字 命令，因此至少会有两个参数
		if len(ctx.Args()) < 2 {
			return fmt.Errorf("missing container name or command")
		}
		containerName := ctx.Args().Get(0)
		// 将除了容器名之外的参数作为命令部分
		commandArray := ctx.Args().Tail()
		ExecContainer(containerName, commandArray)
		return nil
	}),
}

var stopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop a container,e.g. mydocker stop 1234567890",
	Action: cli.ActionFunc(func(ctx *cli.Context) error {
		// 输入应该是：mydocker stop [containerID]
		if len(ctx.Args()) == 0 {
			return errors.New("missing container id")
		}
		containerName := ctx.Args().Get(0)
		stopContainer(containerName)
		return nil
	}),
}

var removeCommand = cli.Command{
	Name:  "rm",
	Usage: "remove unused containers,e.g. mydocker rm 1234567890",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "f", // 强制删除
			Usage: "force delete containers, e.g. mydocker rm 1234567890",
		},
	},
	Action: cli.ActionFunc(func(ctx *cli.Context) error {
		if len(ctx.Args()) == 0 {
			return errors.New("missing container id")
		}
		containerId := ctx.Args().Get(0)
		force := ctx.Bool("f")
		removeContainer(containerId, force)
		return nil
	}),
}

var networkCommand = cli.Command{
	Name:  "network",
	Usage: "container network commands",
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a container network",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "driver",
					Usage: "network driver",
				},
				cli.StringFlag{
					Name:  "subnet",
					Usage: "subnet cidr",
				},
			},
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("missing network name")
				}
				driver := context.String("driver")
				subnet := context.String("subnet")
				name := context.Args()[0]

				err := network.CreateNetwork(driver, subnet, name)
				if err != nil {
					return fmt.Errorf("create network error: %+v", err)
				}
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "list container network",
			Action: func(context *cli.Context) error {
				network.ListNetwork()
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "remove container network",
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("missing network name")
				}
				err := network.DeleteNetwork(context.Args()[0])
				if err != nil {
					return fmt.Errorf("remove network error: %+v", err)
				}
				return nil
			},
		},
	},
}
