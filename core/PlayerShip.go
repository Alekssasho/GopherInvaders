package core

import (
	"time"
)

type playerShipData struct {
	fire       chan struct{}
	cancelFire chan struct{}
}

func (p *playerShipData) startSpawner() {
	p.fire = make(chan struct{}, 1)
	p.cancelFire = make(chan struct{})
	go startSpawner(time.Second/2, p.fire, p.cancelFire)
}
