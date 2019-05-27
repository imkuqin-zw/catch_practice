package service

import (
	"github.com/go-redis/redis"
	"shop/inventory-srv/model"
)

func GetInventory(goodsId uint) (inventory uint64, err error) {
	inventory, err = model.GetInventoryFromCache(goodsId)
	if err == nil || err != redis.Nil {
		return
	}

}
