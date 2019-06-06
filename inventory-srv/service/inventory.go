package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"hash/crc32"
	"shop/basic/config"
	"shop/inventory-srv/consts"
	"shop/inventory-srv/model"
	"sync"
	"time"
)

type InventoryCfg struct {
	QueueNum   uint32 `json:"queue_num"`
	MaxElemNum uint32 `json:"max_elem_num"`
}

type Inventory struct {
	queues     map[uint32]*InventoryQueue
	queueNum   uint32
	maxElemNum uint32
	curElemNum uint32
	elemNumMux sync.Mutex
	ctx        context.Context
}

func initInventory() {
	cfg := &InventoryCfg{}
	err := config.C().Path("inventory_cfg", cfg)
	if err != nil {
		log.Panic("get inventory_queue config fault", zap.Error(err))
	}
	inventory := &Inventory{
		maxElemNum: cfg.MaxElemNum,
		queueNum:   cfg.QueueNum,
	}
	if inventory.queueNum == 0 {
		inventory.queueNum = 100
	}
	if inventory.maxElemNum == 0 {
		inventory.maxElemNum = inventory.queueNum * 10
	}
	inventory.queues = make(map[uint32]*InventoryQueue)
	for i := uint32(0); i < inventory.queueNum; i++ {
		inventory.queues[i] = NewInventoryQueue(i)
		go inventory.queHandle(inventory.ctx, inventory.queues[i])
	}
	s.inventory = inventory
}

func (s *Inventory) queHandle(ctx context.Context, q *InventoryQueue) {
	for {
		act := q.Pop()
		s.inventoryHandle(act)
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

func (s *Inventory) inventoryHandle(act *InventoryAction) {
	if act.Type == consts.INVENTORY_ACTION_READ {
		if err := s.readInventory(act); err != nil {
			log.Error("inventoryHandle read", zap.Error(err), zap.Uint("goods_id", act.GoodsId))
		}
	} else if act.Type == consts.INVENTORY_ACTION_CHANGE {
		s.updateInventory(act)
	}
	return
}

func (s *Inventory) updateInventory(act *InventoryAction) {
	act.mu.Lock()
	if act.status == INVENTORY_ACTION_STATUS_CLOSED {
		return
	}
	act.status = INVENTORY_ACTION_STATUS_RUNING
	act.mu.Unlock()
	if err := model.DelInventoryCache(act.GoodsId); err != nil {
		act.err <- err
		return
	}
	inventory, err := model.ChangeInventoryDB(act.GoodsId, act.Num)
	if err != nil {
		act.err <- err
		return
	}
	act.result <- inventory
	return
}

func (s *Inventory) readInventory(act *InventoryAction) error {
	inventory, err := model.GetInventoryFromDB(act.GoodsId)
	if err != nil {
		return err
	}
	if err = model.SetInventoryCache(act.GoodsId, inventory); err != nil {
		return err
	}
	return nil
}

func (s *Inventory) push(act *InventoryAction) error {
	s.elemNumMux.Lock()
	if s.curElemNum >= s.maxElemNum {
		s.elemNumMux.Unlock()
		return fmt.Errorf("exceeded the maximum number of pending orders")
	}
	s.curElemNum++
	s.elemNumMux.Unlock()

	index := crc32.ChecksumIEEE([]byte(fmt.Sprintf("%d", act.GoodsId))) % s.queueNum
	s.queues[index].Push(act)
	return nil
}

/********** 外部服务 ***********/

func (s *Service) GetInventory(goodsId uint) (uint64, error) {
	inventory, err := model.GetInventoryFromCache(goodsId)
	if err == nil || err != redis.Nil {
		return inventory, err
	}
	err = s.inventory.push(&InventoryAction{Type: consts.INVENTORY_ACTION_READ, GoodsId: goodsId})
	if err != nil {
		log.Warn("GetInventory push fault", zap.Error(err))
	}
	t := time.NewTimer(time.Millisecond * 200)
	for {
		select {
		case <-t.C:
			goto SearchDB
		default:
			inventory, err = model.GetInventoryFromCache(goodsId)
			if err == nil || err != redis.Nil {
				return inventory, err
			}
		}
		time.Sleep(time.Millisecond * 20)
	}
SearchDB:
	return model.GetInventoryFromDB(goodsId)
}

func (s *Service) ChangeInventory(ctx context.Context, goodsId uint, num int64) (uint64, error) {
	act := &InventoryAction{
		status:  INVENTORY_ACTION_STATUS_WAITING,
		Type:    consts.INVENTORY_ACTION_CHANGE,
		GoodsId: goodsId,
		Num:     num,
		result:  make(chan uint64),
		err:     make(chan error),
	}
	err := s.inventory.push(act)
	if err != nil {
		return 0, err
	}
	select {
	case <-ctx.Done():
		act.mu.Lock()
		if act.status == INVENTORY_ACTION_STATUS_WAITING {
			act.status = INVENTORY_ACTION_STATUS_CLOSED
			act.mu.Unlock()
			return 0, fmt.Errorf("time out")
		}
		act.mu.Unlock()
		select {
		case inventory := <-act.result:
			return inventory, nil
		case err := <-act.err:
			return 0, err
		}
	case inventory := <-act.result:
		return inventory, nil
	case err := <-act.err:
		return 0, err
	}
}
