package farmingtales

import (
	"context"
	"github.com/luktech-labs/wax-go-sdk"
	"net/http"

	"github.com/luukkk/wax-gaming-aggregator/foundation/web"
	"github.com/luukkk/wax-gaming-aggregator/pkg/games/farmingtales"

	"github.com/pkg/errors"
)

type accountHandler struct {
	waxSdk *wax.Sdk
}

func (h *accountHandler) harvesting(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	accounts, ok := r.URL.Query()["account"]

	if !ok || len(accounts[0]) < 1 {
		err := errors.Errorf("Url param account is missing")
		return web.RespondError(ctx, w, err, http.StatusBadRequest)
	}

	account := accounts[0]

	miningTasks, err := farmingtales.GetHarvestingTasks(ctx, h.waxSdk, account)
	if err != nil {
		return web.RespondError(ctx, w, err, http.StatusInternalServerError)
	}

	return web.Respond(ctx, w, miningTasks, http.StatusOK)
}
