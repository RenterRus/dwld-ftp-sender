package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/RenterRus/dwld-ftp-sender/internal/controller/grpc"
	"github.com/RenterRus/dwld-ftp-sender/internal/repo/persistent"
	"github.com/RenterRus/dwld-ftp-sender/internal/repo/temporary"
	"github.com/RenterRus/dwld-ftp-sender/internal/usecase/upload"
	"github.com/RenterRus/dwld-ftp-sender/pkg/cache"
	"github.com/RenterRus/dwld-ftp-sender/pkg/ftp"
	"github.com/RenterRus/dwld-ftp-sender/pkg/grpcserver"
	"github.com/RenterRus/dwld-ftp-sender/pkg/sqldb"
)

func NewApp(configPath string) error {
	lastSlash := 0
	for i, v := range configPath {
		if v == '/' {
			lastSlash = i
		}
	}

	conf, err := ReadConfig(configPath[:lastSlash], configPath[lastSlash+1:])
	if err != nil {
		return fmt.Errorf("ReadConfig: %w", err)
	}

	dbconn := sqldb.NewDB(conf.PathToDB, conf.NameDB)
	db := persistent.NewSQLRepo(dbconn, conf.Source.WorkPath)
	cc := cache.NewCache(conf.Cache.Host, conf.Cache.Port)
	cache := temporary.NewMemCache(cc)

	downloadUsecases := upload.NewDownload(
		db,
		cache,
	)

	// FTPSender
	ftpSender := ftp.NewSender(ftp.SenderConf{
		Host:       conf.FTP.Addr.Host,
		User:       conf.FTP.User,
		Pass:       conf.FTP.Pass,
		LocalPath:  conf.Source.WorkPath,
		RemotePath: conf.FTP.RemoteDirectory,
		Port:       conf.FTP.Addr.Port,
		SqlRepo:    db,
		Cache:      cache,
		Enable:     conf.FTP.Addr.Enable,
	})

	ctx, cncl := context.WithCancel(context.Background())
	go cache.Revisor(ctx)
	go ftpSender.Start()

	// gRPC Server
	grpcServer := grpcserver.New(grpcserver.Port(conf.GRPC.Host, strconv.Itoa(conf.GRPC.Port)))
	grpc.NewRouter(grpcServer.App, downloadUsecases)
	go grpcServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Printf("app - Run - signal: %s\n", s.String())
	case err = <-grpcServer.Notify():
		log.Fatal(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	cncl()
	cc.Close()
	ftpSender.Stop()
	err = grpcServer.Shutdown()
	dbconn.Close()
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}

	return nil
}
