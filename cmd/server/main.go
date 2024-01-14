package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GrebenschikovDI/metalsys.git/internal/common/logger"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/repository"
	pb "github.com/GrebenschikovDI/metalsys.git/internal/proto"
	"github.com/GrebenschikovDI/metalsys.git/internal/server/config"
	"github.com/GrebenschikovDI/metalsys.git/internal/server/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/server/grpcserver"
	"github.com/GrebenschikovDI/metalsys.git/internal/server/storages"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const dirPath = "sql/migrations"

var (
	Version = "N/A"
	Date    = "N/A"
	Commit  = "N/A"
)

func printBuildInfo() {
	fmt.Printf("Build version: %s\n", Version)
	fmt.Printf("Build date: %s\n", Date)
	fmt.Printf("Build commit: %s\n", Commit)
}

func main() {
	printBuildInfo()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Info("Error loading config", zap.Error(err))
	}
	var storage repository.Repository
	connStr := cfg.GetDsn()
	if connStr == "" {
		storage = storages.NewMemStorage()
	} else {
		db, err := storages.InitDB(context.Background(), connStr, dirPath)
		if err != nil {
			fmt.Println("NO DB")
		}
		storage = db
	}
	filePath := cfg.GetFileStoragePath()
	err = storages.LoadMetrics(cfg.GetRestore(), filePath, storage)
	if err != nil {
		logger.Log.Info("Error reading from file", zap.String("name", filePath))
	}
	interval := cfg.GetStoreInterval()

	go func() {
		for {
			time.Sleep(interval)
			err := storages.SaveMetrics(filePath, storage)
			if err != nil {
				logger.Log.Info("Error writing in file", zap.String("name", filePath))
			}
		}
	}()

	go func() {
		pprofRouter := controllers.PprofRouter()
		err := http.ListenAndServe(":9091", pprofRouter)
		if err != nil {
			logger.Log.Fatal("Error with profiler")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		runGrpc(storage)
	}()

	if err := run(ctx, storage, *cfg); err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	cancel()

	<-ctx.Done()
}

func run(ctx context.Context, storage repository.Repository, cfg config.ServerConfig) error {
	if err := logger.Initialize("info"); err != nil {
		return err
	}

	ct := controllers.NewControllerContext(storage, cfg)
	router := controllers.MetricsRouter(ct)
	address := cfg.GetServerAddress()

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	logger.Log.Info("Running server", zap.String("address", address))

	go func() {
		// Ожидание сигнала завершения или ошибки при запуске сервера
		select {
		case <-ctx.Done():
			// Отмена контекста
			return
		case err := <-runServer(server):
			if !errors.Is(err, http.ErrServerClosed) {
				// Если сервер завершился по причине ошибки, логгируем ее
				logger.Log.Fatal("Error within server init", zap.Error(err))
			}
		}
	}()

	return nil
}

func runServer(server *http.Server) <-chan error {
	errCh := make(chan error)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		} else {
			errCh <- nil
		}
		close(errCh)
	}()
	return errCh
}

func runGrpc(storage repository.Repository) {
	grpcServer := grpc.NewServer()
	grpcMetric := grpcserver.CreateGrpcServer(storage)
	pb.RegisterMetricsServiceServer(grpcServer, grpcMetric)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("gRPC server is listening on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
