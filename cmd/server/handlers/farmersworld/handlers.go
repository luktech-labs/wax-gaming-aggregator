package farmersworld

import (
	"github.com/luktech-labs/wax-go-sdk"
	"net/http"

	"github.com/luukkk/wax-gaming-aggregator/foundation/web"
)

func Handle(prefix string, router *web.Server, waxSdk *wax.Sdk) http.Handler {
	subRouter := router.Subrouter(prefix)

	ah := accountHandler{waxSdk: waxSdk}
	subRouter.Handle(http.MethodGet, "/mining", ah.mining)
	subRouter.Handle(http.MethodGet, "/animals", ah.animals)

	return subRouter
}
