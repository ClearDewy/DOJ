package config

import (
	"doj-go/JudgeServer/sandbox/envexec"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Config defines go judge server configuration
type Config struct {
	// container
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
	ServerDir string   `flagUsage:"specifies directory to store file upload / download (in memory by default)"`

	// runner limit
	TimeLimitCheckerInterval time.Duration `flagUsage:"specifies time limit checker interval" default:"100ms"`
	ExtraMemoryLimit         *envexec.Size `flagUsage:"specifies extra memory buffer for check memory limit" default:"16k"`
	OutputLimit              *envexec.Size `flagUsage:"specifies POSIX rlimit for output for each command" default:"256m"`
	CopyOutLimit             *envexec.Size `flagUsage:"specifies default file copy out max" default:"256m"`
	OpenFileLimit            int           `flagUsage:"specifies max open file count" default:"256"`
	Cpuset                   string        `flagUsage:"control the usage of cpuset for all containerd process"`
	EnableCPURate            bool          `flagUsage:"enable cpu cgroup rate control"`
	CPUCfsPeriod             time.Duration `flagUsage:"set cpu.cfs_period" default:"100ms"`
	FileTimeout              time.Duration `flagUsage:"specified timeout for filestore files"`

	// server config
	AuthToken     string `flagUsage:"bearer token auth for REST / gRPC"`
	EnableDebug   bool   `flagUsage:"enable debug endpoint"`
	EnableMetrics bool   `flagUsage:"enable promethus metrics endpoint"`

	// logger config
	Release bool `flagUsage:"release level of logs"`
	Silent  bool `flagUsage:"do not print logs"`

	// fix for high memory usage
	ForceGCTarget   *envexec.Size `flagUsage:"specifies force GC trigger heap size" default:"20m"`
	ForceGCInterval time.Duration `flagUsage:"specifies force GC trigger interval" default:"5s"`
}

// Load loads config from flag & environment variables
func (c *Config) Load() error {
	err := loadEnvDefault(c)
	// 自定义变量
	// 如果在容器中运行
	//if os.Getpid() == 1 {
	//	c.Release = true
	//}
	c.EnableDebug = true

	if c.Parallelism <= 0 {
		c.Parallelism = runtime.NumCPU()
	}
	return err
}

func loadEnvDefault(v interface{}) error {
	elem := reflect.ValueOf(v).Elem()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		typeField := elem.Type().Field(i)
		var err error
		envName := camelCaseToEnvVar(typeField.Name)
		if envVal, ok := os.LookupEnv(envName); field.CanSet() {
			if ok {
				err = setFieldWithValue(field, envVal)
			} else {
				if tagValue := typeField.Tag.Get("default"); tagValue != "" {
					err = setFieldWithValue(field, tagValue)
				}
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}
func camelCaseToEnvVar(name string) string {
	// 正则表达式匹配大写字母
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	words := re.FindAllString(name, -1)
	for i := range words {
		words[i] = strings.ToUpper(words[i])
	}
	return strings.Join(words, "_")
}
func setFieldWithValue(field reflect.Value, value string) error {
	if field.CanSet() {
		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int, reflect.Int64:
			// 首先检查字段类型是否为 time.Duration
			if field.Type() == reflect.TypeOf(time.Duration(0)) {
				if duration, err := time.ParseDuration(value); err == nil {
					field.Set(reflect.ValueOf(duration))
				} else {
					return err
				}
			} else {
				// 处理普通的整数
				if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
					field.SetInt(intValue)
				} else {
					return err
				}
			}
		case reflect.Bool:
			if boolValue, err := strconv.ParseBool(value); err == nil {
				field.SetBool(boolValue)
			} else {
				return err
			}

		case reflect.Ptr:
			// 特别处理指向 *envexec.Size 的指针
			if field.Type() == reflect.TypeOf((*envexec.Size)(nil)) {
				sizePtr, ok := field.Interface().(*envexec.Size)
				if !ok || sizePtr == nil {
					// 如果 sizePtr 是 nil，创建一个新的实例
					sizePtr = new(envexec.Size)
					field.Set(reflect.ValueOf(sizePtr))
				}
				sizePtr.Set(value) // 现在可以安全地调用 Set 方法
			}

		default:
			field.Set(reflect.ValueOf(value))
		}
	}
	return nil
}
