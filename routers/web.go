package routers
import (
	"github.com/gorilla/mux"
	"github.com/byrnedo/apibase/controllers"
	"github.com/justinas/alice"
	"net/http"
"github.com/byrnedo/oauthsvc/controllers/web"
	"github.com/RangelReale/osin"
	"github.com/byrnedo/apibase/middleware"
	"github.com/byrnedo/apibase/natsio"
)


func InitWeb(natsCon *natsio.Nats, server *osin.Server){
	var rtr = mux.NewRouter().StrictSlash(true)
	controllers.RegisterMuxRoutes(rtr, web.NewOauthController(natsCon, server))

	//alice is a tiny package to chain middlewares.
	mChain := alice.New(middleware.LogTime).Then(rtr)

	http.Handle("/", mChain)
}
