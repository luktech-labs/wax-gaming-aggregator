package farmersworld

import (
	"context"
	"github.com/luktech-labs/wax-go-sdk"
	"github.com/pkg/errors"
)

type Resources struct {
	CurrentEnergy int `json:"energy"`
	MaxEnergy     int `json:"max_energy"`
}

func GetResources(ctx context.Context, waxSdk TableRowsGetter, account string) (Resources, error) {
	var resources []Resources

	payload := wax.GetTableRowsPayload{
		Json:       true,
		Code:       "farmersworld",
		Scope:      "farmersworld",
		Table:      "accounts",
		LowerBound: account,
		UpperBound: account,
		Limit:      "1",
	}

	err := waxSdk.GetTableRowsContext(ctx, payload, &resources)
	if err != nil {
		return Resources{}, errors.Wrap(err, "getting rows from wax waxSdk")
	}

	if len(resources) == 0 {
		return Resources{}, errors.New("no rows")
	}

	return resources[0], nil
}
