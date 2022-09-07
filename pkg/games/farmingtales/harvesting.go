package farmingtales

import (
	"context"
	"github.com/luktech-labs/wax-go-sdk"
	"time"

	"github.com/pkg/errors"
)

type HarvestingTasks struct {
	HarvestableAnimalIDs []string `json:"harvestable_animal_ids"`
}

func GetHarvestingTasks(ctx context.Context, waxSdk *wax.Sdk, account string) (HarvestingTasks, error) {
	accAnimals, err := GetAccountAnimals(ctx, waxSdk, account)
	if err != nil {
		return HarvestingTasks{}, errors.Wrap(err, "getting account animals")
	}

	var animalIDsToHarvest []string

	for _, animal := range accAnimals {
		if animal.LastHarvest+int64(animal.Cooldown) > time.Now().Unix() {
			continue
		}

		animalIDsToHarvest = append(animalIDsToHarvest, animal.ID)
	}

	return HarvestingTasks{HarvestableAnimalIDs: animalIDsToHarvest}, nil
}
