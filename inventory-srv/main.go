package main

import (
	"flag"
	"github.com/micro/cli"
	"github.com/micro/go-config/source/etcd"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"go.uber.org/zap"
	"shop/basic"
	"shop/basic/common"
	"shop/basic/config"
	"shop/inventory-srv/handler"
	z "shop/plugins/zap"
	"strings"
	"time"

	inventory "shop/inventory-srv/proto/inventory"
)

var (
	appName  = flag.String("cfg_name", "inventory-srv", "service config name")
	etcdAddr = flag.String("cfg_addr", "127.0.0.1:2379", "config etcd address")
	cfg      = &inventoryCfg{}
	log      = z.GetLogger()
)

type inventoryCfg struct {
	common.AppCfg
}

func main() {
	flag.Parse()
	initCfg()
	micReg := etcdv3.NewRegistry(registryOptions)

	// New Service
	service := micro.NewService(
		micro.Name(cfg.Name),
		micro.Version(cfg.Version),
		micro.Registry(micReg),
		micro.RegisterInterval(cfg.RegInterval),
		micro.RegisterTTL(cfg.RegTTL),
		micro.Address(cfg.Address),
	)

	// Initialise service
	service.Init(
		micro.Action(func(context *cli.Context) {
			//// 初始化handler
			//model.Init()
			// 初始化handler
			handler.Init()
		}),
	)

	//Register Handler
	inventory.RegisterInventoryHandler(service.Server(), new(handler.Inventory))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal("service fault", zap.Error(err))
	}
}

func registryOptions(opts *registry.Options) {
	etcdCfg := &common.EtcdCfg{}
	err := config.C().App("etcd", etcdCfg)
	if err != nil {
		log.Panic("get etcd config fault", zap.Error(err))
	}
	opts.Timeout = time.Second * 5
	opts.Addrs = etcdCfg.Addrs
}

func initCfg() {
	source := etcd.NewSource(
		etcd.WithAddress(strings.Split(*etcdAddr, ",")...),
		etcd.WithPrefix("zw.com/shop"),
	)
	basic.Init(
		config.WithSource(source),
		config.WithApp(*appName),
	)
	log.Info("[initCfg] init config completed")
	initAppCfg()
	return
}

func initAppCfg() {
	err := config.C().Path("app", cfg)
	if err != nil {
		log.Panic("get app config fault", zap.Error(err))
	}
	if cfg.RegTTL <= 0 {
		cfg.RegTTL = time.Second * 15
	}
	if cfg.RegInterval > cfg.RegTTL {
		cfg.RegInterval = cfg.RegTTL - 5
		if cfg.RegInterval <= 0 {
			cfg.RegTTL = time.Second * 15
			cfg.RegInterval = time.Second * 10
		}
	}
	return
}
