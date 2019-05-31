package service

import (
	z "shop/plugins/zap"
	"sync"
)

var (
	once sync.Once
	log  *z.Logger
	s    *Service
)

type Service struct {
	inventory *Inventory
}

func Init() {
	once.Do(func() {
		log = z.GetLogger()
		s = &Service{}
		initInventory()
	})

}
