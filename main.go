package main

// http://localhost:14001

import (
	"fmt"
	"github.com/RangelReale/osin"
	"net/http"
	"github.com/byrnedo/apibase"
	"github.com/byrnedo/apibase/db/mongo"
	"github.com/byrnedo/apibase/env"
	"github.com/byrnedo/apibase/natsio"
	"time"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/oauthsvc/routers"
	mgostore "github.com/nguyenxuantuong/osin-mongo-storage"
)

var (
	server *osin.Server
)

func init() {

	apibase.Init()

	mongo.Init(env.GetOr("MONGO_URL", apibase.Conf.GetDefaultString("mongo.url", "")), Trace)

	config := osin.NewServerConfig()
	sstorage := mgostore.NewOAuthStorage(mongo.Conn(), "oauth_osin")

	// MOVE THIS AND MAKE DYNAMIC
	if _, err := sstorage.GetClient("test"); err != nil {
		sstorage.SetClient("test", &osin.DefaultClient{
			Id:          "test",
			Secret:      "superSecret!",
			RedirectUri: "http://localhost:14001/appauth",
		})
	}

	server = osin.NewServer(config, sstorage)

	natsOpts := natsio.NewNatsOptions(func(n *natsio.NatsOptions) error {
		n.Url = env.GetOr("NATS_URL", apibase.Conf.GetDefaultString("nats.url", "nats://localhost:4222"))
		n.Timeout = 10 * time.Second
		return nil
	})

	natsCon, err := natsOpts.Connect()
	if err != nil {
		panic("Failed to connect to nats:" + err.Error())
	}

	routers.InitMq(natsCon, server)

	routers.InitWeb(server)

}
func main() {

	var (
		host string
		port int
		err error
	)

	host = apibase.Conf.GetDefaultString("http.host", "localhost")
	if port, err = env.GetOrInt("PORT", int(apibase.Conf.GetDefaultInt("http.port", 9999))); err != nil {
		panic(err.Error())
	}

	var listenAddr = fmt.Sprintf("%s:%d", host, port)
	Info.Printf("listening on " + listenAddr)
	if err = http.ListenAndServe(listenAddr, nil);err != nil {
		panic("Failed to start server:"+err.Error())
	}

}
