package mysql

import (
	"github.com/jinzhu/gorm"
	"shop/final_consistency/repository"
	m "shop/plugins/mysql"
	z "shop/plugins/zap"
	"sync"
)

var (
	once sync.Once
	log  *z.Logger
	db   *gorm.DB
)

func Init() {
	once.Do(func() {
		log = z.GetLogger()
		db = m.GetMysqlDB()
	})
}

type RepoMysql struct {
}

func (b *RepoMysql) insert(tx *gorm.DB, m interface{}) error {
	return tx.Create(m).Error
}

func (b *RepoMysql) GetById(tx *gorm.DB, m interface{}) error {
	return tx.Model(m).First(m).Error
}

func NewRepo() repository.Repository {
	return &RepoMysql{}
}
