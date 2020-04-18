package resources

import (
	"fmt"
	"net/http"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/felipecurvelo/weather-reporting-api/pkg/authorizer"
	"github.com/felipecurvelo/weather-reporting-api/pkg/internalerror"
	"github.com/julienschmidt/httprouter"
)

type Weather struct {
	api.ResourceBase
	router *httprouter.Router
}

func (weather *Weather) FirstEndpoint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	auth := authorizer.FromContext(r.Context())
	if auth == nil {
		e := internalerror.New("Internal Server Error")
		weather.SetResponse(http.StatusInternalServerError, e, w)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		e := internalerror.New("Empty Token")
		weather.SetResponse(http.StatusInternalServerError, e, w)
		return
	}
	if !auth.ValidateToken(token) {
		e := internalerror.New("Invalid Token")
		weather.SetResponse(http.StatusInternalServerError, e, w)
		return
	}

	responseMessage := fmt.Sprintf("Welcome %s! This is the first endpoint working!", auth.GenerateAccessToken())
	weather.SetResponse(http.StatusOK, responseMessage, w)
}

func (weather *Weather) Register(router *httprouter.Router) {
	weather.router = router
	weather.router.GET("/first_endpoint/", weather.FirstEndpoint)
}
