package common

import (
	"strconv"
	"time"
)

// AppCfg common config
type AppCfg struct {
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	Address     string        `json:"address"`
	Port        int           `json:"port"`
	RegTTL      time.Duration `json:"reg_ttl"`
	RegInterval time.Duration `json:"reg_interval"`
}

func (a *AppCfg) Addr() string {
	return a.Address + ":" + strconv.Itoa(a.Port)
}
