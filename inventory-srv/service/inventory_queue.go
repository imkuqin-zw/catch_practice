package service

import (
	"go.uber.org/zap"
	"shop/basic/config"
	"shop/inventory-srv/consts"
	"sync"
)

type InventoryQueueCfg struct {
	QueueNum uint32 `json:"queue_num"`
}

func initQueue() {
	cfg := &InventoryQueueCfg{}
	err := config.C().Path("inventory_queue", cfg)
	if err != nil {
		log.Panic("get inventory_queue config fault", zap.Error(err))
	}
	s.Queue = make(map[uint32]*InventoryQueue)
	for i := uint32(0); i < cfg.QueueNum; i++ {
		s.Queue[i] = NewInventoryQueue(i)
	}
}

type InventoryAction struct {
	Type    int8
	GoodsId uint
	num     int64
}

type InventoryQueueNode struct {
	Action *InventoryAction
	Next   *InventoryQueueNode
	Prev   *InventoryQueueNode
}

type InventoryQueue struct {
	lock   sync.Mutex
	closed bool
	ch     chan bool
	id     uint32
	length int
	read   map[uint]bool
	root   *InventoryQueueNode
	last   *InventoryQueueNode
}

func NewInventoryQueue(id uint32) *InventoryQueue {
	return &InventoryQueue{
		length: 0,
		root:   nil,
		last:   nil,
		id:     id,
		ch:     make(chan bool),
		read:   make(map[uint]bool),
	}
}

func (q *InventoryQueue) Push(action *InventoryAction) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if action.Type == consts.INVENTORY_ACTION_READ {
		if _, ok := q.read[action.GoodsId]; ok {
			return
		}
		q.read[action.GoodsId] = true
	}
	node := &InventoryQueueNode{
		Action: action,
		Next:   nil,
	}
	if q.root == nil {
		q.root, q.last = node, node
	} else {
		node.Prev = q.last
		q.last.Next = node
		q.last = node
	}
	q.length++
	if q.length == 1 && q.closed {
		q.closed = false
		q.ch <- true
	}
}

func (q *InventoryQueue) Pop() *InventoryAction {
	q.lock.Lock()
	if q.IsEmpty() {
		q.closed = true
		q.lock.Unlock()
		<-q.ch
		return q.Pop()
	}
	q.length--
	node := q.root
	q.root = node.Next
	if node.Next != nil {
		node.Next.Prev = nil
		node.Next = nil
	} else {
		q.last = nil
	}
	if node.Action.Type == consts.INVENTORY_ACTION_READ {
		delete(q.read, node.Action.GoodsId)
	}
	q.lock.Unlock()
	return node.Action
}

func (q *InventoryQueue) IsEmpty() bool {
	return q.length == 0
}

func (q *InventoryQueue) Length() int {
	return q.length
}
