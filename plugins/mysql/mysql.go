package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-log"
	"shop/basic"
	"shop/basic/config"
	"sync"
)

var (
	db     *gorm.DB
	inited bool
	m      sync.RWMutex
)

func init() {
	basic.Register(initDB)
}

func initDB() {
	m.Lock()
	defer m.Unlock()

	if inited {
		err := fmt.Errorf("[initMysql] mysql initialized")
		log.Logf(err.Error())
		return
	}
	initMysql()
	inited = true
}

// Mysql mySQL 配置
type Mysql struct {
	URL               string `json:"url"`
	MaxIdleConnection int    `json:"maxIdleConnection"`
	MaxOpenConnection int    `json:"maxOpenConnection"`
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
	db.DB().SetMaxIdleConns(cfg.MaxIdleConnection)
	db.DB().SetMaxOpenConns(cfg.MaxOpenConnection)
	if err = db.DB().Ping(); err != nil {
		log.Fatal(err)
	}
	log.Logf("[initMysql] Mysql init success")
}

// GetDB 获取db
func GetMysqlDB() *gorm.DB {
	return db
}
