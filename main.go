package main

// http://localhost:14001

import (
	"fmt"
	"github.com/byrnedo/apibase/config"
	. "github.com/byrnedo/apibase/logger"
	_ "github.com/byrnedo/oauthsvc/routers"
	"net/http"
	"github.com/byrnedo/apibase/helpers/envhelp"
)

func main() {

	var (
		host string
		port int
		err  error
	)

	host = config.Conf.GetDefaultString("http.host", "localhost")
	if port, err = envhelp.GetOrInt("PORT", int(config.Conf.GetDefaultInt("http.port", 9999))); err != nil {
		panic(err.Error())
	}

	var listenAddr = fmt.Sprintf("%s:%d", host, port)
	Info.Printf("listening on " + listenAddr)
	if err = http.ListenAndServe(listenAddr, nil); err != nil {
		panic("Failed to start server:" + err.Error())
	}

}
