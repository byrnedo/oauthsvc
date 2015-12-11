package osinserver

import (
	"github.com/RangelReale/osin"
	mgostore "github.com/nguyenxuantuong/osin-mongo-storage"
	"github.com/byrnedo/apibase/db/mongo/defaultmongo"
)

var Server *osin.Server

func init() {
	config := osin.NewServerConfig()
	sstorage := mgostore.NewOAuthStorage(defaultmongo.Conn(), "oauth_osin")

	// MOVE THIS AND MAKE DYNAMIC
	if _, err := sstorage.GetClient("test"); err != nil {
		sstorage.SetClient("test", &osin.DefaultClient{
			Id:          "test",
			Secret:      "superSecret!",
			RedirectUri: "http://localhost:14001/appauth",
		})
	}

	Server = osin.NewServer(config, sstorage)
}
