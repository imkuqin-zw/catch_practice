package repository

import (
	"go.uber.org/zap"
	"shop/basic/config"
	"shop/final_consistency/models"
	"shop/final_consistency/repository/mysql"
	z "shop/plugins/zap"
	"sync"
)

var (
	repo Repository
	log  *z.Logger
	once sync.Once
)

const (
	_ = iota
	MYSQL_DB
)

type RepoConf struct {
	DbType int8 `json:"db_type"`
}

type TransactionMsg interface {
	GetTransMsgById(uint64) (*models.TransactionMsg, error)
	InsertTransMsg(m *models.TransactionMsg) error
}

type Repository interface {
	TransactionMsg
}

func Init() {
	once.Do(func() {
		log = z.GetLogger()
		InitRepository()
	})
}

func InitRepository() {
	c := config.C()
	var dbType int8
	if err := c.Path("app/db_type", &dbType); err != nil {
		log.Panic("get app db_type config fault", zap.Error(err))
	}
	switch dbType {
	case MYSQL_DB:
		mysql.Init()
		repo = mysql.NewRepo()
	default:
		log.Panic("unknown db_type", zap.Int8("db_type", dbType))
	}
}

func GetRepo() Repository {
	return repo
}
