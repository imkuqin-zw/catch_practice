package service

import (
	"context"
	"github.com/go-redis/redis"
	"shop/inventory-srv/model"
)

func (s *Service) GetInventory(goodsId uint) (inventory uint64, err error) {
	inventory, err = model.GetInventoryFromCache(goodsId)
	if err == nil || err != redis.Nil {
		return
	}

}

func (s *Service) queHandle(ctx context.Context, q *InventoryQueue) {
	for {
		msg := q.Pop()
		s.msgHandle(msg)
		s.elemNumMux.Lock()
		s.curElemNum--
		s.elemNumMux.Unlock()
		select {
		case <-ctx.Done():
			break
		default:
		}
	}
}
