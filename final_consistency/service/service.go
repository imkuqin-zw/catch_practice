package service

import (
	"shop/final_consistency/models"
	z "shop/plugins/zap"
	"sync"
)

var (
	once sync.Once
	log  *z.Logger
	s    *Service
	db   models.DBModel
)

type Service struct {
	//dbModel
}

func Init() {
	once.Do(func() {
		log = z.GetLogger()
		s = &Service{}
	})

}
