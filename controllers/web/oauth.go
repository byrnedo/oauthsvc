package web
import (
	"github.com/byrnedo/apibase/routes"
	"net/http"
	. "github.com/byrnedo/apibase/logger"
	"github.com/RangelReale/osin"
	"fmt"
	"net/url"
)


type OauthController struct {
	Server *osin.Server
}

func NewOauthController(server *osin.Server) *OauthController{
	return &OauthController{server}
}


func (oC *OauthController) GetRoutes() []*routes.WebRoute{
	return []*routes.WebRoute{
		routes.NewWebRoute("LoginForm", "/api/v1/authorize", routes.GET, oC.Authorize),
		routes.NewWebRoute("PostCredentials", "/api/v1/authorize", routes.POST, oC.Authorize),
		routes.NewWebRoute("Token", "/api/v1/token", routes.GET, oC.Token),
		routes.NewWebRoute("Info", "/api/v1/info", routes.GET, oC.Info),
	}
}

func (oC *OauthController) Authorize(w http.ResponseWriter, r *http.Request) {
	resp := oC.Server.NewResponse()
	defer resp.Close()

	if ar := oC.Server.HandleAuthorizeRequest(resp, r); ar != nil {
		if !doAuth(r) {
			RenderLoginPage(ar,w,r)
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

func (oC *OauthController) Token(w http.ResponseWriter, r *http.Request) {
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

func (oC *OauthController) Info(w http.ResponseWriter, r *http.Request) {
	resp := oC.Server.NewResponse()
	defer resp.Close()

	if ir := oC.Server.HandleInfoRequest(resp, r); ir != nil {
		oC.Server.FinishInfoRequest(resp, r, ir)
	}
	osin.OutputJSON(resp, w, r)

}

func doAuth(r *http.Request) bool {
	r.ParseForm()
	// talk to data source here.
	if r.Method == "POST" && r.Form.Get("login") == "test" && r.Form.Get("password") == "test" {
		return true
	}
	return false
}

func RenderLoginPage(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) bool {

	w.Write([]byte("<html><body>"))

	w.Write([]byte(fmt.Sprintf("LOGIN %s (use test/test)<br/>", ar.Client.GetId())))
	w.Write([]byte(fmt.Sprintf("<form action=\"/api/v1/authorize?response_type=%s&client_id=%s&state=%s&redirect_uri=%s\" method=\"POST\">",
		ar.Type, ar.Client.GetId(), ar.State, url.QueryEscape(ar.RedirectUri))))

	w.Write([]byte("Login: <input type=\"text\" name=\"login\" /><br/>"))
	w.Write([]byte("Password: <input type=\"password\" name=\"password\" /><br/>"))
	w.Write([]byte("<input type=\"submit\"/>"))

	w.Write([]byte("</form>"))

	w.Write([]byte("</body></html>"))

	return false
}
