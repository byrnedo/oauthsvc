package mq

import (
	"github.com/RangelReale/osin"
	"github.com/byrnedo/apibase/natsio"
	r "github.com/byrnedo/apibase/routes"
	"github.com/nats-io/nats"
)

type OauthController struct {
	routes    []*r.NatsRoute
	natsCon   *natsio.Nats
	oauthServ *osin.Server
}

func (c *OauthController) GetRoutes() []*r.NatsRoute {
	return []*r.NatsRoute{
		r.NewNatsRoute("oauth.token.validate", c.Validate),
	}
}

func NewOauthController(nc *natsio.Nats, server *osin.Server) (oC *OauthController) {
	oC = &OauthController{}
	oC.natsCon = nc
	oC.oauthServ = server
	return
}

func (c *OauthController) Validate(m *nats.Msg) {
	c.natsCon.EncCon.Publish(m.Reply, "Not implemented")
}
