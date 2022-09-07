package handlers

import (
	"github.com/luktech-labs/wax-go-sdk"
	"net/http"

	"github.com/luukkk/wax-gaming-aggregator/cmd/server/handlers/farmersworld"
	"github.com/luukkk/wax-gaming-aggregator/cmd/server/handlers/farmingtales"
	"github.com/luukkk/wax-gaming-aggregator/foundation/web"
)

func Handlers(waxSdk *wax.Sdk) http.Handler {
	server := web.NewServer()
	farmersworld.Handle("/farmers-world", server, waxSdk)
	farmingtales.Handle("/farming-tales", server, waxSdk)

	return server
}
