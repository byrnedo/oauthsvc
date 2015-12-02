package routers
import (
	"github.com/byrnedo/apibase/natsio"
	"github.com/byrnedo/apibase/controllers"
	"github.com/byrnedo/oauthsvc/controllers/mq"
)


func InitMq(natsCon *natsio.Nats) {
	controllers.SubscribeNatsRoutes(natsCon, "oauth_svc_worker", mq.NewOauthController(natsCon.EncCon))
}
