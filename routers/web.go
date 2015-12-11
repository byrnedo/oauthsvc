package routers
import (
	"github.com/byrnedo/apibase/controllers"
	"github.com/justinas/alice"
	"net/http"
"github.com/byrnedo/oauthsvc/controllers/web"
	"github.com/RangelReale/osin"
	"github.com/byrnedo/apibase/middleware"
	"github.com/byrnedo/apibase/natsio"
	"github.com/julienschmidt/httprouter"
)


func InitWeb(natsCon *natsio.Nats, server *osin.Server){
	var rtr = httprouter.New()
	controllers.RegisterRoutes(rtr, web.NewOauthController(natsCon, server))

	//alice is a tiny package to chain middlewares.
	mChain := alice.New(middleware.LogTime).Then(rtr)

	http.Handle("/", mChain)
}
