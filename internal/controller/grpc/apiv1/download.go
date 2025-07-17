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

func (v *V1) SetToQueue(ctx context.Context, in *proto.ToQueueRequest) (*emptypb.Empty, error) {
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

func (v *V1) Status(ctx context.Context, in *emptypb.Empty) (*proto.LoadStatusResponse, error) {
	fmt.Println("================\nStatus")
	defer func() {
		fmt.Println("================")
	}()

	tasks, err := v.u.Status()
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.LoadStatusResponse{
		LinksInWork: lo.Map(tasks.LinksInWork, func(t *usecase.OnWork, _ int) *proto.FileOnWork {
			return &proto.FileOnWork{
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

func (v *V1) Queue(ctx context.Context, in *emptypb.Empty) (*proto.LoadHistoryResponse, error) {
	fmt.Println("================\nQueue")
	defer func() {
		fmt.Println("================")
	}()

	tasks, err := v.u.Queue()
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.LoadHistoryResponse{
		Queue: lo.Map(tasks, func(t *usecase.Task, _ int) *proto.FileInfo {
			return response.TasksToLinks(t)
		}),
	}, nil
}

func (v *V1) Healtheck(ctx context.Context, in *emptypb.Empty) (*proto.SenderHealtheckResponse, error) {
	fmt.Println("================\nHealtheck")
	defer func() {
		fmt.Println("================")
	}()

	return &proto.SenderHealtheckResponse{
		Message: pointer.To("OK"),
	}, nil
}
