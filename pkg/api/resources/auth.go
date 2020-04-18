package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/felipecurvelo/weather-reporting-api/pkg/authorizer"
	"github.com/felipecurvelo/weather-reporting-api/pkg/internalerror"
	"github.com/julienschmidt/httprouter"
)

type Auth struct {
	api.ResourceBase
	router *httprouter.Router
}

type authRequestModel struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type authResponseModel struct {
	Token string `json:"token"`
}

func (a *Auth) Authorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.SetResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	var requestModel authRequestModel
	err = json.Unmarshal(b, &requestModel)
	if err != nil {
		e := internalerror.New(fmt.Sprintf("Invalid request body (%s)", err.Error()))
		a.SetResponse(http.StatusBadRequest, e, w)
		return
	}

	if requestModel.Name != "kirang" {
		e := internalerror.New("Invalid Name")
		a.SetResponse(http.StatusBadRequest, e, w)
		return
	}

	auth := authorizer.FromContext(r.Context())
	if auth == nil {
		e := internalerror.New("Internal Server Error")
		a.SetResponse(http.StatusInternalServerError, e, w)
		return
	}

	token := auth.GenerateAccessToken()

	a.SetResponse(http.StatusOK, authResponseModel{
		Token: token,
	}, w)
}

func (res *Auth) Register(router *httprouter.Router) {
	res.router = router
	res.router.POST("/auth/", res.Authorize)
}
