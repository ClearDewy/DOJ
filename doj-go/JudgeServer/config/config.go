/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package config

import (
	"github.com/criyle/go-judge/envexec"
	"github.com/koding/multiconfig"
	"runtime"
	"time"
)

const (
	ROOT_PATH = "/doj/judge"
)

type CanSet interface {
	Set(string) error
}

var Conf = &Config{}

type Config struct {
	MysqlAddr        string `default:"127.0.0.1:3306"`
	MysqlUsername    string `default:"root"`
	MysqlPassword    string `default:"doj123456"`
	RedisAddr        string `default:"127.0.0.1:6379"`
	RedisPassword    string `default:"doj123456"`
	EtcdAddr         string `flagUsage:"etcd主机" default:"127.0.0.1:2379"`
	EtcdRootPassword string `flagUsage:"etcd root 的密码" default:"doj123456"`
	JudgeServerName  string `flagUsage:"显示在后端的判题机的名称" default:"judge-server-alone"`
	JudgeServerAddr  string `flagUsage:"判题服务器的公网" default:"127.0.0.1:8888"`

	ContainerInitPath  string `flagUsage:"container init path"`
	PreFork            int    `flagUsage:"control # of the prefork workers" default:"1"`
	TmpFsParam         string `flagUsage:"tmpfs mount data (only for default mount with no mount.yaml)" default:"size=128m,nr_inodes=4k"`
	NetShare           bool   `flagUsage:"share net namespace with host"`
	MountConf          string `flagUsage:"specifies mount configuration file" default:"mount.yaml"`
	SeccompConf        string `flagUsage:"specifies seccomp filter" default:"seccomp.yaml"`
	Parallelism        int    `flagUsage:"control the # of concurrency execution (default equal to number of cpu)"`
	CgroupPrefix       string `flagUsage:"control cgroup prefix" default:"gojudge"`
	ContainerCredStart int    `flagUsage:"control the start uid&gid for container (0 uses unprivileged root)" default:"0"`

	// file store
	SrcPrefix []string `flagUsage:"specifies directory prefix for source type copyin (example: -src-prefix=/home,/usr)"`
	Dir       string   `flagUsage:"specifies directory to store file upload / download (in memory by default)"`

	// runner limit
	TimeLimitCheckerInterval time.Duration `flagUsage:"specifies time limit checker interval" default:"100ms"`
	ExtraMemoryLimit         *envexec.Size `flagUsage:"specifies extra memory buffer for check memory limit" default:"16k"`
	OutputLimit              *envexec.Size `flagUsage:"specifies POSIX rlimit for output for each command" default:"256m"`
	CopyOutLimit             *envexec.Size `flagUsage:"specifies default file copy out max" default:"256m"`
	OpenFileLimit            int           `flagUsage:"specifies max open file count" default:"256"`
	Cpuset                   string        `flagUsage:"control the usage of cpuset for all container process"`
	EnableCPURate            bool          `flagUsage:"enable cpu cgroup rate control"`
	CPUCfsPeriod             time.Duration `flagUsage:"set cpu.cfs_period" default:"100ms"`
	FileTimeout              time.Duration `flagUsage:"specified timeout for filestore files"`

	// fix for high memory usage
	ForceGCTarget   *envexec.Size `flagUsage:"specifies force GC trigger heap size" default:"20m"`
	ForceGCInterval time.Duration `flagUsage:"specifies force GC trigger interval" default:"5s"`

	// show version and exit
	Version bool `flagUsage:"show version and exit"`
}

func (c *Config) LoadEnv() error {
	cl := multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.EnvironmentLoader{
			Prefix:    "doj",
			CamelCase: true,
		},
		&multiconfig.FlagLoader{
			Prefix:    "doj",
			CamelCase: true,
		},
	)
	if c.Parallelism <= 0 {
		c.Parallelism = runtime.NumCPU()
	}
	return cl.Load(c)
}
