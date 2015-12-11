package web
// Osin:
// Copyright (c) 2013, Rangel Reale
// All rights reserved.
// modifications:
// Copyright (c) 2015, Donal Byrne
import (
	"github.com/byrnedo/apibase/routes"
	"net/http"
	. "github.com/byrnedo/apibase/logger"
	"github.com/RangelReale/osin"
	"fmt"
	"net/url"
	"html/template"
	"encoding/json"
	"github.com/byrnedo/oauthsvc/msgspec"
	"github.com/byrnedo/apibase/natsio"
	"time"
	"github.com/byrnedo/usersvc/msgspec/mq"
	"github.com/byrnedo/apibase/controllers"
	"github.com/julienschmidt/httprouter"
	"github.com/byrnedo/apibase/natsio/protobuf"
)

type loginViewData struct {
	ClientID string
	PostURL string
}

type OauthController struct {
	controllers.JsonController
	NatsCon *natsio.Nats
	NatsRequestTimeout time.Duration
	Server *osin.Server
}

func NewOauthController(natsCon *natsio.Nats, server *osin.Server) *OauthController{
	return &OauthController{
		NatsCon:natsCon,
		NatsRequestTimeout: 5*time.Second,
		Server: server,
	}
}


func (oC *OauthController) GetRoutes() []*routes.WebRoute{
	return []*routes.WebRoute{
		routes.NewWebRoute("LoginForm", "/api/v1/authorize", routes.GET, oC.Authorize),
		routes.NewWebRoute("PostCredentials", "/api/v1/authorize", routes.POST, oC.Authorize),
		routes.NewWebRoute("Token", "/api/v1/token", routes.GET, oC.Token),
		routes.NewWebRoute("Info", "/api/v1/info", routes.GET, oC.Info),
	}
}

func (oC *OauthController) Authorize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		resp = oC.Server.NewResponse()
	)
	defer resp.Close()

	if ar := oC.Server.HandleAuthorizeRequest(resp, r); ar != nil {
		if !oC.doAuth(r) {
			renderLoginPage(ar,w,r)
			return
		}
		ar.Authorized = true
		oC.Server.FinishAuthorizeRequest(resp, r, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		Error.Printf("ERROR: %s\n", resp.InternalError)
	}
	osin.OutputJSON(resp, w, r)
}

func (oC *OauthController) Token(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := oC.Server.NewResponse()
	defer resp.Close()

	if ar := oC.Server.HandleAccessRequest(resp, r); ar != nil {
		ar.Authorized = true
		oC.Server.FinishAccessRequest(resp, r, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		Error.Printf("ERROR: %s\n", resp.InternalError)
	}
	osin.OutputJSON(resp, w, r)

}

func (oC *OauthController) Info(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := oC.Server.NewResponse()
	defer resp.Close()

	if ir := oC.Server.HandleInfoRequest(resp, r); ir != nil {
		oC.Server.FinishInfoRequest(resp, r, ir)
	}
	osin.OutputJSON(resp, w, r)

}

func (oC *OauthController) doAuth(r *http.Request) (result bool) {
	if r.Method != "POST" {
		return false
	}
	// talk to data source here.
	if oC.AcceptsJson(r) {
		result = oC.doJSONAuth(r)
	} else {
		oC.doFormAuth(r)
		result = oC.doFormAuth(r)
	}
	return
}

func (oC *OauthController) sendAuthRequest(user string, pass string) bool {

	data := mq.NewAuthenticateUserRequest(&mq.InnerAuthenticateUserRequest{Username: &user, Password: &pass})

	response := mq.InnerAuthenticateUserResponse{}

	if err := oC.NatsCon.Request(mq.AuthenticateUserSubject,&protobuf.NatsContext{},data, &response,oC.NatsRequestTimeout); err != nil {
		Error.Println("Failed to make nats request to user svc:", err.Error())
		return false
	}
	Info.Println("Got authenticate reseponse:", response)
	return response.GetAuthenticated()
}

func (oC *OauthController) doJSONAuth(r *http.Request) bool {
	var (
		d = json.NewDecoder(r.Body)
		creds = &msgspec.AuthorizeRequest{}
	)
	if err := d.Decode(creds); err != nil {
		Error.Println("Failed to decode json:" + err.Error())
		return false
	}


	return oC.sendAuthRequest(creds.User, creds.Password)
}

func (oC *OauthController) doFormAuth(r *http.Request) bool {
	r.ParseForm()
	user := r.Form.Get("user")
	password := r.Form.Get("password")


	return oC.sendAuthRequest(user, password)
}

func renderLoginPage(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) {
	var (
		err error
		t *template.Template
	)
	if t, err = template.ParseFiles("./views/login.html"); err != nil {
		Error.Println("Failed to parse template:" + err.Error())
	}
	t.Execute(w, loginViewData{
		ClientID: ar.Client.GetId(),
		PostURL: fmt.Sprintf("/api/v1/authorize?response_type=%s&client_id=%s&state=%s&redirect_uri=%s",
			ar.Type, ar.Client.GetId(), ar.State, url.QueryEscape(ar.RedirectUri)),
	})
}
