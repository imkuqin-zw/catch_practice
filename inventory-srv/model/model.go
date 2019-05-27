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
	lock   sync.RWMutex
	inited bool
	log    *z.Logger
	db     *gorm.DB
	rd     redis.Cmdable
)

func Init() {
	lock.Lock()
	defer lock.Unlock()

	if inited {
		log.Warn("handle initialized")
		return
	}

	log = z.GetLogger()
	db = m.GetMysqlDB()
	rd = r.GetRedis()

	inited = true
}
