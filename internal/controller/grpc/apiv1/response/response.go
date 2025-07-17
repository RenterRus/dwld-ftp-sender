package response

import (
	proto "github.com/RenterRus/dwld-ftp-sender/docs/proto/v1"

	"github.com/RenterRus/dwld-ftp-sender/internal/usecase"
)

func TasksToLinks(task *usecase.Task) *proto.Task {
	return &proto.Task{
		Link:          task.Link,
		TargetQuality: task.MaxQuantity,
		Status:        task.Status,
		Name:          task.Name,
		Message:       task.Message,
	}
}
