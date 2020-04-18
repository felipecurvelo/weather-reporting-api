package resources

import (
	"io/ioutil"
	"net/http"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/julienschmidt/httprouter"
)

type Auth struct {
	api.ResourceBase
	router *httprouter.Router
}

func (auth *Auth) Authorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		auth.SetResponse(http.StatusBadRequest, err.Error(), w)
		return
	}
	auth.SetResponse(http.StatusOK, body, w)
}

func (res *Auth) Register(router *httprouter.Router) {
	res.router = router
	res.router.POST("/auth/", res.Authorize)
}
