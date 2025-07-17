package grpc

import (
	"github.com/RenterRus/dwld-ftp-sender/internal/usecase"

	v1 "github.com/RenterRus/dwld-ftp-sender/internal/controller/grpc/apiv1"

	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouter(app *pbgrpc.Server, usecases usecase.Downloader) {
	v1.NewDownloadRoutes(app, usecases)
	reflection.Register(app)
}
