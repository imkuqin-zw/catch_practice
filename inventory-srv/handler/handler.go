package handler

import (
	z "shop/plugins/zap"
	"sync"
)

var (
	log  *z.Logger
	once sync.Once
)

func Init() {
	once.Do(func() {
		log = z.GetLogger()
	})
}
