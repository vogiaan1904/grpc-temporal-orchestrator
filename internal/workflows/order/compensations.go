package order

import (
	"github.com/vogiaan1904/order-orchestrator/internal/activities"
	orderpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/order"
	"go.temporal.io/sdk/workflow"
)

func releaseInventory(ctx workflow.Context, items []*orderpb.OrderItem) error {
	return workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, getProductActivityOptions()),
		(*activities.ProductActivities).ReleaseInventory,
		items,
	).Get(ctx, nil)
}
