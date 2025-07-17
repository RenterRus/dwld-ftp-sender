package apiv1

import (
	"context"
	"fmt"

	"github.com/RenterRus/dwld-ftp-sender/internal/controller/grpc/apiv1/response"

	proto "github.com/RenterRus/dwld-ftp-sender/docs/proto/v1"

	"github.com/RenterRus/dwld-ftp-sender/internal/usecase"

	"github.com/AlekSi/pointer"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (v *V1) SetToQueue(ctx context.Context, in *proto.SetToQueueRequest) (*emptypb.Empty, error) {
	fmt.Println("================\nSetToQueue")
	defer func() {
		fmt.Println(in)
		fmt.Println("================")
	}()

	if in == nil || pointer.Get(in).Link == "" {
		return nil, fmt.Errorf("SetToQueue: empty request")
	}

	if err := v.u.SetToQueue(in.GetLink(), in.TargetQuality); err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (v *V1) CleanHistory(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	fmt.Println("================\nCleanHistory")
	defer func() {
		fmt.Println("================")
	}()

	if err := v.u.CleanHistory(); err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (v *V1) Status(ctx context.Context, in *emptypb.Empty) (*proto.StatusResponse, error) {
	fmt.Println("================\nStatus")
	defer func() {
		fmt.Println("================")
	}()

	tasks, err := v.u.Status()
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.StatusResponse{
		LinksInWork: lo.Map(tasks.LinksInWork, func(t *usecase.OnWork, _ int) *proto.OnWork {
			return &proto.OnWork{
				Link:           t.Link,
				Filename:       t.Filename,
				MoveTo:         t.MoveTo,
				TargetQuantity: t.TargetQuantity,
				Procentage:     t.Procentage,
				Status:         t.Status,
				CurrentSize:    t.CurrentSize,
				TotalSize:      t.TotalSize,
				Message:        t.Message,
			}
		}),
	}, nil
}

func (v *V1) Queue(ctx context.Context, in *emptypb.Empty) (*proto.HistoryResponse, error) {
	fmt.Println("================\nQueue")
	defer func() {
		fmt.Println("================")
	}()

	tasks, err := v.u.Queue()
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.HistoryResponse{
		Queue: lo.Map(tasks, func(t *usecase.Task, _ int) *proto.Task {
			return response.TasksToLinks(t)
		}),
	}, nil
}

func (v *V1) Healtheck(ctx context.Context, in *emptypb.Empty) (*proto.HealtheckResponse, error) {
	fmt.Println("================\nHealtheck")
	defer func() {
		fmt.Println("================")
	}()

	return &proto.HealtheckResponse{
		Message: pointer.To("OK"),
	}, nil
}
