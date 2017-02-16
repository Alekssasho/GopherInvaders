package core

import (
	"math/rand"
	"time"
)

type enemyShipData struct {
	startX float32
	delta  float32

	fire       chan struct{}
	cancelFire chan struct{}
}

func (e *enemyShipData) startSpawner() {
	e.fire = make(chan struct{}, 1)
	e.cancelFire = make(chan struct{}, 1)
	go func() {
		// some random sleep to reduce patterns
		time.Sleep(time.Duration(rand.Float32() * float32(time.Second) * 5))
		e.fire <- struct{}{}
		startSpawner(
			time.Second*2+time.Duration(float32(time.Second*2)*rand.Float32()),
			e.fire,
			e.cancelFire)
	}()
}
