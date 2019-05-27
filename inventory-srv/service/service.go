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
	Queue      map[uint32]*InventoryQueue
	maxElemNum uint32
	curElemNum uint32
	elemNumMux sync.Mutex
	ctx        context.Context
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
	initQueue()

	inited = true
}
