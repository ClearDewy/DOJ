/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package sandbox

import (
	"context"
	"doj-go/JudgeServer/config"
	"github.com/ClearDewy/go-pkg/logrus"
	"github.com/criyle/go-judge/env"
	"github.com/criyle/go-judge/env/pool"
	"github.com/criyle/go-judge/envexec"
	"github.com/criyle/go-judge/filestore"
	"github.com/criyle/go-judge/pb"
	"github.com/criyle/go-judge/worker"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

const (
	//SandBoxGrpcAddr = "127.0.0.1:5051"
	MaxProcLimit = 128
)

var (
	Client pb.ExecutorServer
)

func Init() (func() error, func(ctx context.Context) error) {
	//if runtime.GOOS == "linux" {
	//	logrus.Info("容器初始化")
	//	container.Init()
	//}
	// Init environment pool
	fs, fsCleanUp := newFilesStore()
	b, _ := newEnvBuilder()
	envPool := pool.NewPool(b)
	prefork(envPool, config.Conf.PreFork)
	work := newWorker(envPool, fs)
	work.Start()
	logrus.Infof("Started worker with parallelism=%d, workDir=%s, timeLimitCheckInterval=%v",
		config.Conf.Parallelism, config.Conf.Dir, config.Conf.TimeLimitCheckerInterval)

	Client = New(work, fs, config.Conf.SrcPrefix)

	newForceGCWorker()

	return nil, func(ctx context.Context) error {
		work.Shutdown()
		logrus.Info("Worker shutdown")
		err := fsCleanUp()
		logrus.Info("FileStore cleaned up")
		return err
	}
}

func prefork(envPool worker.EnvironmentPool, prefork int) {
	if prefork <= 0 {
		return
	}
	logrus.Info("create ", prefork, " prefork containers")
	m := make([]envexec.Environment, 0, prefork)
	for i := 0; i < prefork; i++ {
		e, err := envPool.Get()
		if err != nil {
			log.Fatalln("prefork environment failed ", err)
		}
		m = append(m, e)
	}
	for _, e := range m {
		envPool.Put(e)
	}
}

func newFilesStore() (filestore.FileStore, func() error) {
	const timeoutCheckInterval = 15 * time.Second
	var cleanUp func() error

	var fs filestore.FileStore
	if config.Conf.Dir == "" {
		if runtime.GOOS == "linux" {
			config.Conf.Dir = "/dev/shm"
		} else {
			config.Conf.Dir = os.TempDir()
		}
		var err error
		config.Conf.Dir, err = os.MkdirTemp(config.Conf.Dir, "go-judge")
		if err != nil {
			logrus.Fatal("failed to create file store temp dir", err)
		}
		cleanUp = func() error {
			return os.RemoveAll(config.Conf.Dir)
		}
	}
	os.MkdirAll(config.Conf.Dir, 0755)
	fs = filestore.NewFileLocalStore(config.Conf.Dir)

	if config.Conf.FileTimeout > 0 {
		fs = filestore.NewTimeout(fs, config.Conf.FileTimeout, timeoutCheckInterval)
	}
	return fs, cleanUp
}

func newEnvBuilder() (pool.EnvBuilder, map[string]any) {
	b, param, err := env.NewBuilder(env.Config{
		ContainerInitPath:  config.Conf.ContainerInitPath,
		MountConf:          config.Conf.MountConf,
		TmpFsParam:         config.Conf.TmpFsParam,
		NetShare:           config.Conf.NetShare,
		CgroupPrefix:       config.Conf.CgroupPrefix,
		Cpuset:             config.Conf.Cpuset,
		ContainerCredStart: config.Conf.ContainerCredStart,
		EnableCPURate:      config.Conf.EnableCPURate,
		CPUCfsPeriod:       config.Conf.CPUCfsPeriod,
		SeccompConf:        config.Conf.SeccompConf,
		Logger:             logrus.StandardLogger(),
	})
	if err != nil {
		logrus.Fatal("create environment builder failed ", err)
	}
	return b, param
}

func newWorker(envPool worker.EnvironmentPool, fs filestore.FileStore) worker.Worker {
	return worker.New(worker.Config{
		FileStore:             fs,
		EnvironmentPool:       envPool,
		Parallelism:           config.Conf.Parallelism,
		WorkDir:               config.Conf.Dir,
		TimeLimitTickInterval: config.Conf.TimeLimitCheckerInterval,
		ExtraMemoryLimit:      *config.Conf.ExtraMemoryLimit,
		OutputLimit:           *config.Conf.OutputLimit,
		CopyOutLimit:          *config.Conf.CopyOutLimit,
		OpenFileLimit:         uint64(config.Conf.OpenFileLimit),
	})
}

func newForceGCWorker() {
	go func() {
		ticker := time.NewTicker(config.Conf.ForceGCInterval)
		for {
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			if mem.HeapInuse > config.Conf.ForceGCTarget.Byte() {
				logrus.Infof("Force GC as heap_in_use(%v) > target(%v)",
					envexec.Size(mem.HeapInuse), config.Conf.ForceGCTarget)
				runtime.GC()
				debug.FreeOSMemory()
			}
			<-ticker.C
		}
	}()
}
