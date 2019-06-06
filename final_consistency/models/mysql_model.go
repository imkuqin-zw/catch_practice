package models

import (
	"github.com/jinzhu/gorm"
	m "shop/plugins/mysql"
)

type mysqlModel struct {
	db *gorm.DB
}

func initMysqlDb() *mysqlModel {
	return &mysqlModel{
		db: m.GetMysqlDB(),
	}
}

func (m *mysqlModel) InsertTransactionMsg(msg *TransactionMsg) error {
	return m.db.Create(msg).Error
}

func (m *mysqlModel) GetTransactionMsg(msg *TransactionMsg) (res *TransactionMsg, err error) {
	res = &TransactionMsg{}
	err = m.db.Where(msg).First(res).Error
	return
	m.db.Begin()
}
