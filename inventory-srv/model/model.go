package model

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	m "shop/plugins/mysql"
	r "shop/plugins/redis"
	z "shop/plugins/zap"
	"sync"
)

var (
	once sync.Once
	log  *z.Logger
	db   *gorm.DB
	rd   redis.Cmdable
)

func Init() {
	once.Do(func() {
		log = z.GetLogger()
		db = m.GetMysqlDB()
		rd = r.GetRedis()
	})
}
