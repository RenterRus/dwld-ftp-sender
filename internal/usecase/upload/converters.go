package upload

import (
	"strconv"

	"github.com/RenterRus/dwld-ftp-sender/internal/usecase"

	"github.com/RenterRus/dwld-ftp-sender/internal/repo/persistent"
)

func LinkToTask(item persistent.LinkModel, _ int) *usecase.Task {
	return &usecase.Task{
		Link:        item.Link,
		MaxQuantity: strconv.Itoa(item.TargetQuantity),
		Status:      item.WorkStatus,
		Name:        item.Filename,
		Message:     item.Message,
	}
}
