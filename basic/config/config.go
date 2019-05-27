package config

import (
	"fmt"
	"github.com/micro/go-config"
	"github.com/micro/go-log"
	"shop/basic/consts"
	"strings"
	"sync"
)

var (
	m      sync.RWMutex
	inited bool
	c      = &configurator{}
)

type Configurator interface {
	App(string, interface{}) error
	Path(string, interface{}) error
}

type configurator struct {
	conf       config.Config
	appName    string
	appPrefix  []string
	pathPrefix []string
}

func (c *configurator) App(name string, config interface{}) (err error) {
	v := c.conf.Get(append(c.appPrefix, name)...)
	if v != nil {
		err = v.Scan(config)
	} else {
		err = fmt.Errorf("[App] config not found，err：%s", name)
	}
	return
}

func (c *configurator) Path(path string, config interface{}) (err error) {
	v := c.conf.Get(append(c.pathPrefix, path)...)
	if v != nil {
		err = v.Scan(config)
	} else {
		err = fmt.Errorf("[Path] config not found，err：%s", path)
	}
	return
}

func C() Configurator {
	return c
}

func (c *configurator) init(ops Options) (err error) {
	m.Lock()
	defer m.Unlock()
	if inited {
		log.Logf("[init] config has been initialized")
		return
	}
	c.conf = config.NewConfig()
	c.appName = ops.AppName
	c.appPrefix = strings.Split(consts.CFG_PREFIX, "/")
	c.pathPrefix = append(c.appPrefix, c.appName)

	// 加载配置
	err = c.conf.Load(ops.Sources...)
	if err != nil {
		log.Fatal(err)
	}

	//监听配置变化
	//c.watch()
	// 标记已经初始化
	inited = true
	return
}

func (c *configurator) watch() {
	go func() {
		log.Logf("[init] watch config changes")
		// 开始侦听变动事件
		watcher, err := c.conf.Watch()
		if err != nil {
			log.Fatal(err)
		}
		for {
			v, err := watcher.Next()
			if err != nil {
				log.Fatal(err)
			}
			log.Logf("[init] config changed: %v", string(v.Bytes()))
		}
	}()
}

// Init 初始化配置
func Init(opts ...Option) {
	ops := Options{}
	for _, o := range opts {
		o(&ops)
	}

	c = &configurator{}

	c.init(ops)
}
