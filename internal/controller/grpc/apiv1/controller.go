package apiv1

import (
	proto "github.com/RenterRus/dwld-ftp-sender/docs/proto/v1"

	"github.com/RenterRus/dwld-ftp-sender/internal/usecase"
)

type V1 struct {
	proto.SenderServer

	u usecase.Downloader
}
