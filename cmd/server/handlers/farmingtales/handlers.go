package farmingtales

import (
	"github.com/luktech-labs/wax-go-sdk"
	"net/http"

	"github.com/luukkk/wax-gaming-aggregator/foundation/web"
)

func Handle(prefix string, router *web.Server, waxSdk *wax.Sdk) {
	subRouter := router.Subrouter(prefix)

	ah := accountHandler{waxSdk: waxSdk}
	subRouter.Handle(http.MethodGet, "/harvesting", ah.harvesting)
}
