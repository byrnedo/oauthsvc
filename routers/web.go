package routers

import (
	"github.com/byrnedo/apibase/controllers"
	"github.com/byrnedo/apibase/middleware"
	"github.com/byrnedo/oauthsvc/controllers/web"
	"github.com/byrnedo/oauthsvc/osinserver"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"github.com/byrnedo/apibase/natsio/defaultnats"
)

func init() {
	var rtr = httprouter.New()
	controllers.RegisterRoutes(rtr, web.NewOauthController(defaultnats.Conn, osinserver.Server))

	//alice is a tiny package to chain middlewares.
	mChain := alice.New(middleware.LogTime).Then(rtr)

	http.Handle("/", mChain)
}
