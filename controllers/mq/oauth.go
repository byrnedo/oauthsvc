package mq

import (
	r "github.com/byrnedo/apibase/routes"
"github.com/apcera/nats"
)


type OauthController struct {
	routes []*r.NatsRoute
	encCon *nats.EncodedConn
}

func (c *OauthController) GetRoutes() []*r.NatsRoute {
	return []*r.NatsRoute{
		r.NewNatsRoute("oauth.token.validate", c.Validate),
	}
}

func NewOauthController(nc *nats.EncodedConn) (oC *OauthController) {
	oC = &OauthController{}
	oC.encCon = nc
	return
}

func (c *OauthController) Validate(m *nats.Msg) {
	c.encCon.Publish(m.Reply, "Not implemented")
}
