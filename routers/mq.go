package routers
import (
	"github.com/byrnedo/apibase/natsio"
	"github.com/byrnedo/apibase/controllers"
	"github.com/byrnedo/oauthsvc/controllers/mq"
	"github.com/RangelReale/osin"
)


func InitMq(natsCon *natsio.Nats, server *osin.Server) {
	controllers.SubscribeNatsRoutes(natsCon, "oauth_svc_worker", mq.NewOauthController(natsCon.EncCon, server))
}
