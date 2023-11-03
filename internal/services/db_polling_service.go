package service

import (
	"context"
	"sync"
	"time"
)

type pollingService struct {
	interval time.Duration
}

func NewPollingService(interval time.Duration) PollingService {
	return &pollingService{

		interval: interval,
	}
}

type PollingService interface {
	Polling(ctx context.Context, wg *sync.WaitGroup, interval time.Duration, stopChan <-chan struct{})
}
