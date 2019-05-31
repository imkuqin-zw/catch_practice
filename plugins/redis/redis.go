package redis

import (
	r "github.com/go-redis/redis"
	"github.com/micro/go-log"
	"shop/basic"
	"shop/basic/config"
	"strings"
	"sync"
	"time"
)

var (
	client r.Cmdable
	once   sync.Once
)

func init() {
	basic.Register(initRedis)
}

type redis struct {
	Addr         []string       `json:"addr"`
	Pwd          string         `json:"password"`
	DBNum        int            `json:"db_num"`
	PoolSize     int            `json:"pool_size"`
	DialTimeout  time.Duration  `json:"dial_timeout"`
	ReadTimeout  time.Duration  `json:"read_timeout"`
	WriteTimeout time.Duration  `json:"write_timeout"`
	MinIdleConns int            `json:"min_idle_conns"`
	MaxRetries   int            `json:"max_retries"`
	Sentinel     *RedisSentinel `json:"sentinel"`
}

type RedisSentinel struct {
	Master string   `json:"master"`
	XNodes []string `json:"nodes"`
	nodes  []string
}

// Nodes redis 哨兵节点列表
func (s *RedisSentinel) GetNodes() []string {
	if len(s.XNodes) != 0 {
		for _, v := range s.XNodes {
			v = strings.TrimSpace(v)
			s.nodes = append(s.nodes, v)
		}
	}
	return s.nodes
}

func initRedis() {
	once.Do(func() {
		c := config.C()
		cfg := &redis{}
		err := c.App("redis", cfg)
		if err != nil {
			log.Fatalf("[initRedis] %s", err)
			panic(err)
		}
		connReids(cfg)
		log.Logf("[initRedis] redis initializing completed")
	})
}

func initSingle(redisConfig *redis) {
	client = r.NewClient(&r.Options{
		Addr:         redisConfig.Addr[0],
		Password:     redisConfig.Pwd,   // no password set
		DB:           redisConfig.DBNum, // use default DB
		PoolSize:     redisConfig.PoolSize,
		DialTimeout:  redisConfig.DialTimeout,
		ReadTimeout:  redisConfig.ReadTimeout,
		WriteTimeout: redisConfig.WriteTimeout,
		MinIdleConns: redisConfig.MinIdleConns,
		MaxRetries:   redisConfig.MaxRetries,
	})
}

func initSentinel(redisConfig *redis) {
	client = r.NewFailoverClient(&r.FailoverOptions{
		MasterName:    redisConfig.Sentinel.Master,
		SentinelAddrs: redisConfig.Sentinel.GetNodes(),
		DB:            redisConfig.DBNum,
		Password:      redisConfig.Pwd,
		PoolSize:      redisConfig.PoolSize,
		DialTimeout:   redisConfig.DialTimeout,
		ReadTimeout:   redisConfig.ReadTimeout,
		WriteTimeout:  redisConfig.WriteTimeout,
		MinIdleConns:  redisConfig.MinIdleConns,
		MaxRetries:    redisConfig.MaxRetries,
	})
}

func initCluster(redisConfig *redis) {
	client = r.NewClusterClient(&r.ClusterOptions{
		Addrs:        redisConfig.Addr,
		Password:     redisConfig.Pwd,
		PoolSize:     redisConfig.PoolSize,
		DialTimeout:  redisConfig.DialTimeout,
		ReadTimeout:  redisConfig.ReadTimeout,
		WriteTimeout: redisConfig.WriteTimeout,
		MinIdleConns: redisConfig.MinIdleConns,
		MaxRetries:   redisConfig.MaxRetries,
	})
}

func connReids(cfg *redis) {
	if cfg.Sentinel != nil {
		log.Log("[initRedis] init redis，sentinel mode")
		initSentinel(cfg)
	} else {
		addrCount := len(cfg.Addr)
		if addrCount == 0 {
			log.Fatalf("[initRedis] %s", "not found redis addr in config")
		}
		if addrCount > 1 {
			log.Log("[initRedis] init redis，normal mode")
			initSingle(cfg)
		} else {
			log.Log("[initRedis] init redis，cluster mode")
			initCluster(cfg)
		}
	}
	log.Log("[initRedis] init Redis，check ping")
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Logf("[initRedis] init Redis，check Ping %s", pong)
}

// Redis 获取redis
func GetRedis() r.Cmdable {
	return client
}
