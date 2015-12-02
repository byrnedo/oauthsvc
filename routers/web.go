package routers
import (
	"github.com/gorilla/mux"
	"github.com/byrnedo/apibase/controllers"
	"github.com/justinas/alice"
	"net/http"
"github.com/byrnedo/oauthsvc/controllers/web"
	"github.com/RangelReale/osin"
)


func InitWeb(server *osin.Server){
	var rtr = mux.NewRouter().StrictSlash(true)
	controllers.RegisterMuxRoutes(rtr, web.NewOauthController(server))

	//alice is a tiny package to chain middlewares.
	mChain := alice.New().Then(rtr)

	http.Handle("/", mChain)
}
