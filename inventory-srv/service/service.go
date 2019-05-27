package service

import (
	z "shop/plugins/zap"
	"sync"
)

var (
	m      sync.RWMutex
	inited bool
	log    *z.Logger
)

func Init() {
	m.Lock()
	defer m.Unlock()

	if inited {
		log.Warn("service initialized")
		return
	}

	log = z.GetLogger()

	inited = true
}
