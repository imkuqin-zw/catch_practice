package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	inventory "shop/inventory-srv/proto/inventory"
)

type Inventory struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Inventory) UpdateInventory(ctx context.Context, req *inventory.ReqUpdateInventory, rsp *inventory.InventoryCount) error {
	log.Debug("Received inventory.UpdateInventory request")
	rsp.Count = 5
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Inventory) GetInventory(ctx context.Context, req *inventory.GoodsId, rsp *inventory.InventoryCount) error {
	log.Debug("Received inventory.GetInventory request", zap.Uint32("goods_id", req.GoodsId))
	if req.GoodsId == 0 {
		return fmt.Errorf("goods id must greater than zero")
	}
	rsp.Count = 20
	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Inventory) PingPong(ctx context.Context, stream inventory.Inventory_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Debug("Got ping ", zap.Int64("stroke", req.Stroke))
		if err := stream.Send(&inventory.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
