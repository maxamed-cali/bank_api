package jobs

import (
	"bank/services"
	"fmt"
	"time"
)

func StartAutoExpireJob() {
	ticker := time.NewTicker(1 * time.Minute) // adjust interval as needed

	go func() {
		for {
			select {
			case <-ticker.C:
				services.AutoExpireRequests()
				fmt.Println("Auto-expire job ran at:", time.Now())
			}
		}
	}()
}