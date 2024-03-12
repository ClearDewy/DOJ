// Command go-judge will starts a http server that receives command to run.go
// programs inside a sandbox.
package main

import (
	"context"
	crypto_rand "crypto/rand"
	"doj-go/JudgeServer/pb"
	"doj-go/JudgeServer/sandbox/config"
	"doj-go/JudgeServer/sandbox/env"
	"doj-go/JudgeServer/sandbox/env/pool"
	"doj-go/JudgeServer/sandbox/envexec"
	"doj-go/JudgeServer/sandbox/filestore"
	grpcexecutor "doj-go/JudgeServer/sandbox/grpc_executor"
	"doj-go/JudgeServer/sandbox/worker"
	"encoding/binary"
	"flag"
	"log"
	math_rand "math/rand"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const GRPC_ADDR = ":5051"

var logger *zap.Logger

func main() {
	conf := loadConf()
	initLogger(conf)
	defer logger.Sync()
	logger.Sugar().Infof("config loaded: %+v", conf)
	initRand()
	warnIfNotLinux()

	// Init environment pool
	fs, fsCleanUp := newFilsStore(conf)
	b, _ := newEnvBuilder(conf)
	envPool := newEnvPool(b, conf.EnableMetrics)
	prefork(envPool, conf.PreFork)
	work := newWorker(conf, envPool, fs)
	work.Start()
	logger.Sugar().Infof("Started worker with parallelism=%d, workdir=%s, timeLimitCheckInterval=%v",
		conf.Parallelism, conf.ServerDir, conf.TimeLimitCheckerInterval)

	servers := []initFunc{
		cleanUpWorker(work),
		cleanUpFs(fsCleanUp),
		initGRPCServer(conf, work, fs),
	}

	// Gracefully shutdown, with signal / HTTP server / gRPC server / Monitor HTTP server
	sig := make(chan os.Signal, 1+len(servers))

	// worker and fs clean up func
	stops := []stopFunc{}
	for _, s := range servers {
		start, stop := s()
		if start != nil {
			go func() {
				start()
				sig <- os.Interrupt
			}()
		}
		if stop != nil {
			stops = append(stops, stop)
		}
	}

	// background force GC worker
	newForceGCWorker(conf)

	// Graceful shutdown...
	signal.Notify(sig, os.Interrupt)
	<-sig
	signal.Reset(os.Interrupt)

	logger.Sugar().Info("Shutting Down...")

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()

	var eg errgroup.Group
	for _, s := range stops {
		s := s
		eg.Go(func() error {
			return s(ctx)
		})
	}

	go func() {
		logger.Sugar().Info("Shutdown Finished ", eg.Wait())
		cancel()
	}()
	<-ctx.Done()
}

func warnIfNotLinux() {
	if runtime.GOOS != "linux" {
		logger.Sugar().Warn("Platform is ", runtime.GOOS)
		logger.Sugar().Warn("Please notice that the primary supporting platform is Linux")
		logger.Sugar().Warn("Windows and macOS(darwin) support are only recommended in development environment")
	}
}

func loadConf() *config.Config {
	var conf config.Config
	if err := conf.Load(); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		log.Fatalln("load config failed ", err)
	}
	return &conf
}

type stopFunc func(ctx context.Context) error
type initFunc func() (start func(), cleanUp stopFunc)

func cleanUpWorker(work worker.Worker) initFunc {
	return func() (start func(), cleanUp stopFunc) {
		return nil, func(ctx context.Context) error {
			work.Shutdown()
			logger.Sugar().Info("Worker shutdown")
			return nil
		}
	}
}

func cleanUpFs(fsCleanUp func() error) initFunc {
	return func() (start func(), cleanUp stopFunc) {
		if fsCleanUp == nil {
			return nil, nil
		}
		return nil, func(ctx context.Context) error {
			err := fsCleanUp()
			logger.Sugar().Info("FileStore cleaned up")
			return err
		}
	}
}

func initGRPCServer(conf *config.Config, work worker.Worker, fs filestore.FileStore) initFunc {
	return func() (start func(), cleanUp stopFunc) {
		// Init gRPC server
		esServer := grpcexecutor.New(work, fs, conf.SrcPrefix, logger)
		grpcServer := newGRPCServer(conf, esServer)

		return func() {
				lis, err := newListener(GRPC_ADDR)
				if err != nil {
					logger.Sugar().Error("gRPC listen failed: ", err)
					return
				}
				logger.Sugar().Info("Starting gRPC server at ", GRPC_ADDR, " with listener ", printListener(lis))
				logger.Sugar().Info("gRPC server stopped: ", grpcServer.Serve(lis))
			}, func(ctx context.Context) error {
				grpcServer.GracefulStop()
				logger.Sugar().Info("GRPC server shutdown")
				return nil
			}
	}
}

func initLogger(conf *config.Config) {
	if conf.Silent {
		logger = zap.NewNop()
		return
	}

	var err error
	if conf.Release {
		logger, err = zap.NewProduction()
	} else {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		if !conf.EnableDebug {
			config.Level.SetLevel(zap.InfoLevel)
		}
		logger, err = config.Build()
	}
	if err != nil {
		log.Fatalln("init logger failed ", err)
	}
}

func initRand() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		logger.Fatal("random generator init failed ", zap.Error(err))
	}
	sd := int64(binary.LittleEndian.Uint64(b[:]))
	logger.Sugar().Infof("random seed: %d", sd)
	math_rand.Seed(sd)
}

func prefork(envPool worker.EnvironmentPool, prefork int) {
	if prefork <= 0 {
		return
	}
	logger.Sugar().Info("create ", prefork, " prefork containers")
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

func newGRPCServer(conf *config.Config, esServer pb.ExecutorServer) *grpc.Server {
	grpc_zap.ReplaceGrpcLoggerV2(logger)
	streamMiddleware := []grpc.StreamServerInterceptor{
		grpc_prometheus.StreamServerInterceptor,
		grpc_zap.StreamServerInterceptor(logger),
		grpc_recovery.StreamServerInterceptor(),
	}
	unaryMiddleware := []grpc.UnaryServerInterceptor{
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(logger),
		grpc_recovery.UnaryServerInterceptor(),
	}
	if conf.AuthToken != "" {
		authFunc := grpcTokenAuth(conf.AuthToken)
		streamMiddleware = append(streamMiddleware, grpc_auth.StreamServerInterceptor(authFunc))
		unaryMiddleware = append(unaryMiddleware, grpc_auth.UnaryServerInterceptor(authFunc))
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamMiddleware...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryMiddleware...)),
	)
	pb.RegisterExecutorServer(grpcServer, esServer)
	grpc_prometheus.Register(grpcServer)
	grpc_prometheus.EnableHandlingTimeHistogram()
	return grpcServer
}

func grpcTokenAuth(token string) func(context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		reqToken, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}
		if reqToken != token {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}
		return ctx, nil
	}
}

func newFilsStore(conf *config.Config) (filestore.FileStore, func() error) {
	const timeoutCheckInterval = 15 * time.Second
	var cleanUp func() error

	var fs filestore.FileStore
	if conf.ServerDir == "" {
		if runtime.GOOS == "linux" {
			conf.ServerDir = "/dev/shm"
		} else {
			conf.ServerDir = os.TempDir()
		}
		var err error
		conf.ServerDir, err = os.MkdirTemp(conf.ServerDir, "go-judge")
		if err != nil {
			logger.Sugar().Fatal("failed to create file store temp dir", err)
		}
		cleanUp = func() error {
			return os.RemoveAll(conf.ServerDir)
		}
	}
	os.MkdirAll(conf.ServerDir, 0755)
	fs = filestore.NewFileLocalStore(conf.ServerDir)
	if conf.EnableDebug {
		fs = newMetricsFileStore(fs)
	}
	if conf.FileTimeout > 0 {
		fs = filestore.NewTimeout(fs, conf.FileTimeout, timeoutCheckInterval)
	}
	return fs, cleanUp
}

func newEnvBuilder(conf *config.Config) (pool.EnvBuilder, map[string]any) {
	b, param, err := env.NewBuilder(env.Config{
		ContainerInitPath:  conf.ContainerInitPath,
		MountConf:          conf.MountConf,
		TmpFsParam:         conf.TmpFsParam,
		NetShare:           conf.NetShare,
		CgroupPrefix:       conf.CgroupPrefix,
		Cpuset:             conf.Cpuset,
		ContainerCredStart: conf.ContainerCredStart,
		EnableCPURate:      conf.EnableCPURate,
		CPUCfsPeriod:       conf.CPUCfsPeriod,
		SeccompConf:        conf.SeccompConf,
		Logger:             logger.Sugar(),
	})
	if err != nil {
		logger.Sugar().Fatal("create environment builder failed ", err)
	}
	if conf.EnableMetrics {
		b = &metriceEnvBuilder{b}
	}
	return b, param
}

func newEnvPool(b pool.EnvBuilder, enableMetrics bool) worker.EnvironmentPool {
	p := pool.NewPool(b)
	if enableMetrics {
		p = &metricsEnvPool{p}
	}
	return p
}

func newWorker(conf *config.Config, envPool worker.EnvironmentPool, fs filestore.FileStore) worker.Worker {
	return worker.New(worker.Config{
		FileStore:             fs,
		EnvironmentPool:       envPool,
		Parallelism:           conf.Parallelism,
		WorkDir:               conf.ServerDir,
		TimeLimitTickInterval: conf.TimeLimitCheckerInterval,
		ExtraMemoryLimit:      *conf.ExtraMemoryLimit,
		OutputLimit:           *conf.OutputLimit,
		CopyOutLimit:          *conf.CopyOutLimit,
		OpenFileLimit:         uint64(conf.OpenFileLimit),
		ExecObserver:          execObserve,
	})
}

func newForceGCWorker(conf *config.Config) {
	go func() {
		ticker := time.NewTicker(conf.ForceGCInterval)
		for {
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			if mem.HeapInuse > uint64(*conf.ForceGCTarget) {
				logger.Sugar().Infof("Force GC as heap_in_use(%v) > target(%v)",
					envexec.Size(mem.HeapInuse), *conf.ForceGCTarget)
				runtime.GC()
				debug.FreeOSMemory()
			}
			<-ticker.C
		}
	}()
}
