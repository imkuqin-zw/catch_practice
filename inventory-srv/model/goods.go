package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"shop/inventory-srv/consts"
)

type Goods struct {
	gorm.Model
	Inventory uint64 `json:"inventory"`
}

func (Goods) TableName() string {
	return "goods"
}

func SetInventoryCache(goodsId uint, inventory uint64) error {
	cacheKey := fmt.Sprintf(consts.GOODS_INFO, goodsId)
	b, err := rd.HSet(cacheKey, "inventory", inventory).Result()
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("set inventory cache fault")
	}
	return nil
}

func DelInventoryCache(goodsId uint) error {
	cacheKey := fmt.Sprintf(consts.GOODS_INFO, goodsId)
	b, err := rd.HDel(cacheKey, "inventory").Result()
	if err != nil {
		return err
	}
	if b == 0 {
		return fmt.Errorf("del inventory cache fault")
	}
	return nil
}

func ChangeInventoryDBByGoods(goods *Goods, num int64) error {
	inventory := int64(goods.Inventory) + num
	if inventory < 0 {
		return fmt.Errorf("inventory shortage")
	}
	err := db.Model(goods).
		Where("inventory = ?", goods.Inventory).
		Update("inventory", uint64(inventory)).Error
	if err != nil {
		return err
	}
	return nil
}

func ChangeInventoryDB(goodsId uint, num int64) (uint64, error) {
	tx := db.Begin()
	goods := &Goods{}
	if err := tx.Select("id, inventory").First(&goods, goodsId).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	inventory := int64(goods.Inventory) + num
	if inventory < 0 {
		return 0, fmt.Errorf("inventory shortage")
	}
	lastInventory := uint64(inventory)
	err := tx.Model(goods).
		Where("inventory = ?", goods.Inventory).
		Update("inventory", lastInventory).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return lastInventory, nil
}

func GetInventoryFromCache(goodsId uint) (uint64, error) {
	cacheKey := fmt.Sprintf(consts.GOODS_INFO, goodsId)
	return rd.HGet(cacheKey, "inventory").Uint64()
}

func GetInventoryFromDB(goodsId uint) (uint64, error) {
	goods := &Goods{}
	if err := db.Select("id, inventory").First(&goods, goodsId).Error; err != nil {
		return 0, err
	}
	return goods.Inventory, nil
}
