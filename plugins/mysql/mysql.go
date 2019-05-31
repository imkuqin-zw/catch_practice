package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-log"
	"shop/basic"
	"shop/basic/config"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

func init() {
	basic.Register(initDB)
}

func initDB() {
	once.Do(func() {
		initMysql()
	})
}

// Mysql mySQL 配置
type Mysql struct {
	URL               string `json:"url"`
	LogMode           bool   `json:"log_mode"`
	MaxIdleConnection int    `json:"max_idle_connection"`
	MaxOpenConnection int    `json:"max_open_connection"`
}

func initMysql() {
	c := config.C()
	cfg := &Mysql{}
	err := c.App("mysql", cfg)
	if err != nil {
		log.Fatalf("[initMysql] %s", err)
		panic(err)
	}
	db, err = gorm.Open("mysql", cfg.URL)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	db.LogMode(cfg.LogMode)
	if cfg.MaxIdleConnection != 0 {
		db.DB().SetMaxIdleConns(cfg.MaxIdleConnection)
	}
	if cfg.MaxOpenConnection != 0 {
		db.DB().SetMaxOpenConns(cfg.MaxOpenConnection)
	}
	if err = db.DB().Ping(); err != nil {
		log.Fatal(err)
	}
	log.Logf("[initMysql] Mysql init success")
}

// GetDB 获取db
func GetMysqlDB() *gorm.DB {
	return db
}
