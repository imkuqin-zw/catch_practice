package models

import (
	z "shop/plugins/zap"
	"sync"
)

type DBModel interface {
	//存储事务消息
	InsertTransactionMsg(*TransactionMsg) error

	//获取事务消息
	GetTransactionMsg(*TransactionMsg) (*TransactionMsg, error)
}

var (
	once sync.Once
	log  *z.Logger
	db   DBModel
)

func Init() {
	once.Do(func() {
		log = z.GetLogger()
		db = initMysqlDb()
	})
}

func GetDBModel() DBModel {
	return db
}
