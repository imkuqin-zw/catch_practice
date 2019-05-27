package main

import (
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
	"shop/inventory-srv/model"
	"shop/inventory-srv/service"
	z "shop/plugins/zap"
	"strings"
	"time"

	inventory "shop/inventory-srv/proto/inventory"
)

var (
	appName  string
	etcdAddr string
	cfg      = &inventoryCfg{}
	log      = z.GetLogger()
)

type inventoryCfg struct {
	common.AppCfg
}

func main() {
	svr := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:        "cfg_name",
				Usage:       "service config name",
				Value:       "inventory-srv",
				Destination: &appName,
			},
			cli.StringFlag{
				Name:        "cfg_addr",
				Usage:       "config etcd address",
				Value:       "192.168.2.118:2379",
				Destination: &etcdAddr,
			},
		),
	)
	// Initialise Cmd
	svr.Init()

	initCfg()
	micReg := etcdv3.NewRegistry(registryOptions)
	// Initialise service
	svr.Init(
		micro.Name(cfg.Name),
		micro.Version(cfg.Version),
		micro.Registry(micReg),
		micro.RegisterInterval(cfg.RegInterval),
		micro.RegisterTTL(cfg.RegTTL),
		micro.RegisterInterval(cfg.RegTTL),
		micro.Address(cfg.Address),
		micro.Action(func(context *cli.Context) {
			//init model
			model.Init()
			//init service
			service.Init()
			//init handler
			handler.Init()
		}),
	)

	//Register Handler
	inventory.RegisterInventoryHandler(svr.Server(), new(handler.Inventory))

	// Run service
	if err := svr.Run(); err != nil {
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
		etcd.WithAddress(strings.Split(etcdAddr, ",")...),
		etcd.WithPrefix("zw.com/shop"),
	)
	basic.Init(
		config.WithSource(source),
		config.WithApp(appName),
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
