package farmersworld

import (
	"context"
	"github.com/luktech-labs/wax-go-sdk"
	"time"

	"github.com/pkg/errors"
)

type TableRowsGetter interface {
	GetTableRows(wax.GetTableRowsPayload, interface{}) error
	GetTableRowsContext(context.Context, wax.GetTableRowsPayload, interface{}) error
}

type MiningTasks struct {
	EnergyToRefill    int      `json:"energy_to_refill"`
	RepairableToolIDs []string `json:"repairable_tool_ids"`
	ClaimableToolIDs  []string `json:"claimable_tool_ids"`
}

func GetMiningTasks(ctx context.Context, waxSdk TableRowsGetter, account string) (MiningTasks, error) {
	accTools, err := GetAccountTools(ctx, waxSdk, account)
	if err != nil {
		return MiningTasks{}, errors.Wrap(err, "getting account tools")
	}

	var totalEnergyConsumed int
	var claimableToolIDs []string
	var repairableToolIDs []string

	for _, tool := range accTools {
		if tool.NextAvailability < time.Now().Unix() {
			claimableToolIDs = append(claimableToolIDs, tool.ID)
			totalEnergyConsumed += tool.EnergyConsumed
		}
		if tool.CurrentDurability < tool.DurabilityConsumed {
			repairableToolIDs = append(repairableToolIDs, tool.ID)
		}
	}

	resources, err := GetResources(ctx, waxSdk, account)
	if err != nil {
		return MiningTasks{}, errors.Wrap(err, "getting account resources")
	}

	var energyToRefill int
	if resources.CurrentEnergy < totalEnergyConsumed {
		energyToRefill = resources.MaxEnergy - resources.CurrentEnergy
	}

	at := MiningTasks{
		EnergyToRefill:    energyToRefill,
		RepairableToolIDs: repairableToolIDs,
		ClaimableToolIDs:  claimableToolIDs,
	}

	return at, nil
}
