package upload

import (
	"fmt"

	"github.com/RenterRus/dwld-ftp-sender/internal/repo/persistent"
	"github.com/RenterRus/dwld-ftp-sender/internal/repo/temporary"
	"github.com/RenterRus/dwld-ftp-sender/internal/usecase"
	"github.com/samber/lo"
)

type downlaoder struct {
	dbRepo    persistent.SQLRepo
	cacheRepo temporary.CacheRepo
}

func NewDownload(dbRepo persistent.SQLRepo, cache temporary.CacheRepo) usecase.Downloader {
	return &downlaoder{
		dbRepo:    dbRepo,
		cacheRepo: cache,
	}
}

func (d *downlaoder) SetToQueue(link string, targerQunatity int32) error {
	if _, err := d.dbRepo.Insert(link, int(targerQunatity)); err != nil {
		return fmt.Errorf("SetToQueue: %w", err)
	}

	return nil
}

func (d *downlaoder) CleanHistory() error {
	if _, err := d.dbRepo.DeleteHistory(); err != nil {
		return fmt.Errorf("CleanHistory: %w", err)
	}

	return nil
}

func (d *downlaoder) Status() (*usecase.StatusResponse, error) {
	resp, err := d.cacheRepo.GetStatus()
	if err != nil {
		return nil, fmt.Errorf("CleanHistory: %w", err)
	}

	links := make([]*usecase.OnWork, 0, len(resp.WorkStatus)*2)

	for link, v := range resp.WorkStatus {
		for file, info := range v {
			links = append(links, &usecase.OnWork{
				Link:           link,
				Filename:       file,
				MoveTo:         info.MoveTo,
				TargetQuantity: int64(info.MaxQuality),
				Procentage:     info.Procentage,
				Status:         info.Status,
				TotalSize:      info.DownloadSize,
				CurrentSize:    info.CurrentSize,
				Message:        info.Message,
			})
		}
	}

	return &usecase.StatusResponse{
		LinksInWork: links,
	}, nil
}

func (d *downlaoder) Queue() ([]*usecase.Task, error) {
	resp, err := d.dbRepo.SelectHistory(nil)
	if err != nil {
		return nil, fmt.Errorf("Queue: %w", err)
	}

	return lo.Map(resp, LinkToTask), nil
}
