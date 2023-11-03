package service

import (
	"context"
	"sync"
	"time"
)

func (p *pollingService) Polling(ctx context.Context, wg *sync.WaitGroup, interval time.Duration, stopChan <-chan struct{}) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// fmt.Println("Polling the database...")
			// p.transactionService.PollScheduledTransaction(ctx, time.Now().Truncate(5*time.Minute))

		case <-stopChan:
			return
		}
	}
}
