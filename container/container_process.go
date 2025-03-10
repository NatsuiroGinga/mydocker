package container

import (
	"fmt"
	"hash/fnv"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/NatsuiroGinga/mydocker/constant"
	"github.com/NatsuiroGinga/mydocker/utils"
	"github.com/sirupsen/logrus"
)

const (
	RUNNING       = "running"
	STOP          = "stopped"
	Exit          = "exited"
	InfoLoc       = "/var/lib/mydocker/containers/"
	InfoLocFormat = InfoLoc + "%s/"
	ConfigName    = "config.json"
	LogFile       = "%s-json.log"
	MetaFile      = "/var/lib/mydocker/containers/meta.json"
)

type Info struct {
	Pid         string   `json:"pid"`         // 容器的 init 进程在宿主机上的PID
	Id          string   `json:"id"`          // 容器 ID
	Name        string   `json:"name"`        // 容器名
	Command     string   `json:"command"`     // 容器内 init 运行命令
	CreatedTime string   `json:"createTime"`  // 创建时间
	Status      string   `json:"status"`      // 容器的状态
	Volume      string   `json:"volume"`      // 容器挂载的 volume
	NetworkName string   `json:"networkName"` // 容器所在的网络
	PortMapping []string `json:"portmapping"` // 端口映射
	IP          string   `json:"ip"`          // ip地址
}

// NewParentProcess 创建并返回一个新进程. 注意: 在本函数内进程尚未启动
/*
这里是父进程，也就是当前进程执行的内容。
1.这里的/proc/se1f/exe调用中，/proc/self/ 指的是当前运行进程自己的环境，exec 其实就是自己调用了自己，使用这种方式对创建出来的进程进行初始化
2.后面的args是参数，其中init是传递给本进程的第一个参数，在本例中，其实就是会去调用initCommand去初始化进程的一些环境和资源
3.下面的clone参数就是去fork出来一个新进程，并且使用了namespace隔离新创建的进程和外部环境。
4.如果用户指定了-it参数，就需要把当前进程的输入输出导入到标准输入输出上
*/
func NewParentProcess(tty bool, containerId, imageName, volume string, envs []string) (*exec.Cmd, *os.File) {
	// 创建匿名管道用于传递参数，将readPipe作为子进程的ExtraFiles，子进程从readPipe中读取参数
	// 父进程中则通过writePipe将参数写入管道
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		logrus.Errorf("New pipe error %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else { // 对于后台运行容器，将 stdout、stderr 重定向到日志文件中，便于后续查看
		dirPath := fmt.Sprintf(InfoLocFormat, containerId)
		if err := os.MkdirAll(dirPath, constant.Perm0622); err != nil {
			logrus.Errorf("NewParentProcess mkdir %s error %v", dirPath, err)
			return nil, nil
		}
		stdLogFilePath := dirPath + GetLogfile(containerId)
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			logrus.Errorf("NewParentProcess create file %s error %v", stdLogFilePath, err)
			return nil, nil
		}
		cmd.Stdout = stdLogFile
		cmd.Stderr = stdLogFile
	}

	// 指定 cmd 的工作目录为我们前面准备好的用于存放busybox rootfs的目录
	NewWorkSpace(containerId, imageName, volume)

	cmd.Env = append(os.Environ(), envs...)
	cmd.ExtraFiles = []*os.File{readPipe}
	if len(envs) != 0 {
		cmd.Env = append(cmd.Env, envs...)
	}
	cmd.Dir = utils.GetMerged(containerId)

	return cmd, writePipe
}

// GenerateContainerID 根据容器名生成容器id
func GenerateContainerID(seed string) string {
	generator := fnv.New32()
	generator.Write([]byte(seed))
	generator.Write([]byte(time.Now().String()))
	return strconv.Itoa(int(generator.Sum32()))
}

// GetLogfile build logfile name by containerId
func GetLogfile(containerId string) string {
	return fmt.Sprintf(LogFile, containerId)
}
