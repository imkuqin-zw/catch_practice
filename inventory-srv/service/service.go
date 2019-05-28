package service

import (
	z "shop/plugins/zap"
	"sync"
)

var (
	m      sync.RWMutex
	inited bool
	log    *z.Logger
	s      *Service
)

type Service struct {
	inventory *Inventory
}

func Init() {
	m.Lock()
	defer m.Unlock()

	if inited {
		log.Warn("service initialized")
		return
	}
	log = z.GetLogger()

	s = &Service{}
	initInventory()

	inited = true
}
