package service

import (
	"shop/inventory-srv/consts"
	"sync"
)

const (
	INVENTORY_ACTION_STATUS_WAITING = iota
	INVENTORY_ACTION_STATUS_RUNING
	INVENTORY_ACTION_STATUS_CLOSED
)

type InventoryAction struct {
	mu      sync.RWMutex
	status  int8
	Type    int8
	GoodsId uint
	Num     int64
	result  chan uint64
	err     chan error
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
