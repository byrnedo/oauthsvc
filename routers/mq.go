package routers

import (
	"github.com/byrnedo/apibase/controllers"
	"github.com/byrnedo/oauthsvc/controllers/mq"
	"github.com/byrnedo/oauthsvc/osinserver"
	"github.com/byrnedo/apibase/natsio/defaultnats"
)

func init() {

	controllers.SubscribeNatsRoutes(defaultnats.Conn, "oauth_svc_worker", mq.NewOauthController(defaultnats.Conn, osinserver.Server))
}
