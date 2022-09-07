package farmersworld

import (
	"context"
	"github.com/luktech-labs/wax-go-sdk"
	"github.com/pkg/errors"
)

type ToolInfo struct {
	ID                string `json:"asset_id"`
	TemplateID        int    `json:"template_id"`
	MaxDurability     int    `json:"durability"`
	CurrentDurability int    `json:"current_durability"`
	NextAvailability  int64  `json:"next_availability"`
}

type ToolConf struct {
	TemplateName       string `json:"template_name"`
	TemplateID         int    `json:"template_id"`
	Type               string `json:"type"`
	EnergyConsumed     int    `json:"energy_consumed"`
	DurabilityConsumed int    `json:"durability_consumed"`
}

type Tool struct {
	ToolInfo
	ToolConf
}

func GetAccountTools(ctx context.Context, waxSdk TableRowsGetter, account string) ([]Tool, error) {
	payload := wax.GetTableRowsPayload{
		Json:          true,
		Code:          "farmersworld",
		Scope:         "farmersworld",
		Table:         "tools",
		IndexPosition: 2,
		KeyType:       "i64",
		Limit:         "200",
		LowerBound:    account,
		UpperBound:    account,
	}

	var toolsInfo []ToolInfo
	err := waxSdk.GetTableRowsContext(ctx, payload, &toolsInfo)
	if err != nil {
		return nil, errors.Wrap(err, "getting rows from wax service")
	}

	toolConfs, err := getToolConfs(ctx, waxSdk)
	if err != nil {
		return nil, errors.Wrap(err, "getting tool confs")
	}

	return aggregate(toolsInfo, toolConfs), nil
}

func getToolConfs(ctx context.Context, waxSdk TableRowsGetter) ([]ToolConf, error) {
	var toolConfs []ToolConf

	payload := wax.GetTableRowsPayload{
		Json:  true,
		Code:  "farmersworld",
		Scope: "farmersworld",
		Table: "toolconfs",
		Limit: "200",
	}

	err := waxSdk.GetTableRowsContext(ctx, payload, &toolConfs)
	if err != nil {
		return nil, errors.Wrap(err, "getting rows from wax service")
	}

	return toolConfs, nil
}

type toolConfByTemplateID map[int]ToolConf

func aggregate(info []ToolInfo, conf []ToolConf) []Tool {
	toolConfsByID := make(toolConfByTemplateID)
	for _, toolConf := range conf {
		toolConfsByID[toolConf.TemplateID] = toolConf
	}

	tools := make([]Tool, 0, len(info))
	for _, tool := range info {
		toolConf, found := toolConfsByID[tool.TemplateID]
		if !found {
			continue
		}

		tools = append(tools, Tool{
			ToolInfo: tool,
			ToolConf: toolConf,
		})
	}

	return tools
}
