package ftp

import (
	"context"

	"github.com/RenterRus/dwld-ftp-sender/internal/controller/ftp"
	"github.com/RenterRus/dwld-ftp-sender/internal/repo/persistent"
	"github.com/RenterRus/dwld-ftp-sender/internal/repo/temporary"
)

type SenderConf struct {
	Host       string
	User       string
	Pass       string
	LocalPath  string
	RemotePath string
	Port       int
	Enable     bool
	SqlRepo    persistent.SQLRepo
	Cache      temporary.CacheRepo
}

type Sender struct {
	FTP    ftp.Sender
	notify chan struct{}
}

func NewSender(conf SenderConf) *Sender {
	return &Sender{
		FTP: ftp.NewFTPSender(&ftp.FTPSenderConf{
			Host:       conf.Host,
			User:       conf.User,
			Pass:       conf.Pass,
			LocalPath:  conf.LocalPath,
			RemotePath: conf.RemotePath,
			Port:       conf.Port,
			SqlRepo:    conf.SqlRepo,
			Cache:      conf.Cache,
			Enable:     conf.Enable,
		}),
		notify: make(chan struct{}, 1),
	}
}

func (s *Sender) Start() {
	ctx, cncl := context.WithCancel(context.Background())
	go func() {
		s.FTP.Loader(ctx)
	}()

	<-s.notify
	cncl()
}

func (s *Sender) Stop() {
	s.notify <- struct{}{}
}
