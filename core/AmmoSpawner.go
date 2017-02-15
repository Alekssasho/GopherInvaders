package core

import "time"

func startSpawner(dur time.Duration, receive chan<- struct{}, cancel <-chan struct{}) {
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-cancel:
			ticker.Stop()
			return
		case <-ticker.C:
			receive <- struct{}{}
		}
	}
}
