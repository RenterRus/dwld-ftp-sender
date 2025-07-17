package apiv1

import (
	proto "github.com/RenterRus/dwld-ftp-sender/docs/proto/v1"

	"github.com/RenterRus/dwld-ftp-sender/internal/usecase"

	pbgrpc "google.golang.org/grpc"
)

func NewDownloadRoutes(app *pbgrpc.Server, usecases usecase.Downloader) {
	r := &V1{
		u: usecases,
	}

	proto.RegisterSenderServer(app, r)
}
