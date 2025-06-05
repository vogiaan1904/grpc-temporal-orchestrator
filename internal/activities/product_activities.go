package activities

import (
	"context"
	"fmt"
	"log"

	orderpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/order"
	productpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/product"
)

type ProductActivities struct {
	Client productpb.ProductServiceClient
}

func (a *ProductActivities) ReserveInventory(ctx context.Context, items []*orderpb.OrderItem) error {
	rItems := make([]*productpb.ReserveInventoryItem, len(items))
	for i, item := range items {
		rItems[i] = &productpb.ReserveInventoryItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}
	}
	log.Printf("Reserving inventory: %+v", rItems)
	_, err := a.Client.ReserveInventory(ctx, &productpb.ReserveInventoryRequest{
		Items: rItems,
	})
	if err != nil {
		return fmt.Errorf("failed to reserve inventory: %w", err)
	}

	return nil
}

func (a *ProductActivities) ReleaseInventory(ctx context.Context, items []*orderpb.OrderItem) error {
	rItems := make([]*productpb.ReleaseInventoryItem, len(items))
	for i, item := range items {
		rItems[i] = &productpb.ReleaseInventoryItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}
	}
	_, err := a.Client.ReleaseInventory(ctx, &productpb.ReleaseInventoryRequest{
		Items: rItems,
	})
	if err != nil {
		return fmt.Errorf("failed to release inventory: %w", err)
	}

	return nil
}

func (a *ProductActivities) UpdateStock(ctx context.Context, items []*orderpb.OrderItem) error {
	uItems := make([]*productpb.UpdateStockItem, len(items))
	for i, item := range items {
		uItems[i] = &productpb.UpdateStockItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}
	}
	_, err := a.Client.UpdateStock(ctx, &productpb.UpdateStockRequest{
		Items: uItems,
	})
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}
