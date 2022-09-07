package farmingtales

import (
	"context"
	"github.com/luktech-labs/wax-go-sdk"
	"github.com/pkg/errors"
)

type AnimalInfo struct {
	ID          string `json:"asset_id"`
	TemplateID  int    `json:"template_id"`
	LastHarvest int64  `json:"last_harvest"`
}

type AnimalConf struct {
	TemplateName string `json:"label"`
	TemplateID   int    `json:"template_id"`
	Cooldown     int    `json:"cooldown"`
	FoodConsumed int    `json:"food"`
}

type Animal struct {
	AnimalInfo
	AnimalConf
}

func GetAccountAnimals(ctx context.Context, waxSdk *wax.Sdk, account string) ([]Animal, error) {
	var animalsInfo []AnimalInfo

	payload := wax.GetTableRowsPayload{
		Json:          true,
		Code:          "farminggames",
		Scope:         "farminggames",
		Table:         "animal",
		IndexPosition: 2,
		KeyType:       "i64",
		Limit:         "500",
		LowerBound:    account,
		UpperBound:    account,
	}

	err := waxSdk.GetTableRowsContext(ctx, payload, &animalsInfo)
	if err != nil {
		return nil, errors.Wrap(err, "getting rows from wax service")
	}

	animalsConf, err := getAnimalConfs(ctx, waxSdk)
	if err != nil {
		return nil, errors.Wrap(err, "getting animal confs")
	}

	return aggregate(animalsInfo, animalsConf), nil
}

func getAnimalConfs(ctx context.Context, waxSdk *wax.Sdk) ([]AnimalConf, error) {
	var animalConfs []AnimalConf

	payload := wax.GetTableRowsPayload{
		Json:  true,
		Code:  "farminggames",
		Scope: "farminggames",
		Table: "confanimal",
		Limit: "200",
	}

	err := waxSdk.GetTableRowsContext(ctx, payload, &animalConfs)
	if err != nil {
		return nil, errors.Wrap(err, "getting rows from wax service")
	}

	return animalConfs, nil
}

type animalConfByTemplateID map[int]AnimalConf

func aggregate(info []AnimalInfo, conf []AnimalConf) []Animal {
	animalConfsByID := make(animalConfByTemplateID)
	for _, animalConf := range conf {
		animalConfsByID[animalConf.TemplateID] = animalConf
	}

	animals := make([]Animal, 0, len(info))
	for _, animal := range info {
		animalConf, found := animalConfsByID[animal.TemplateID]
		if !found {
			continue
		}
		animals = append(animals, Animal{
			AnimalInfo: animal,
			AnimalConf: animalConf,
		})
	}

	return animals
}
