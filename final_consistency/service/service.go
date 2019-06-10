package service

import (
	"shop/final_consistency/repository"
	z "shop/plugins/zap"
	"sync"
)

var (
	once sync.Once
	log  *z.Logger
	s    *Service
)

type Service struct {
	repository.Repository
}

func Init() {
	once.Do(func() {
		log = z.GetLogger()
		s = &Service{
			repository.GetRepo(),
		}
	})

}
