package resources

import (
	"fmt"
	"net/http"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/felipecurvelo/weather-reporting-api/pkg/authorizer"
	"github.com/julienschmidt/httprouter"
)

type Weather struct {
	api.ResourceBase
	router *httprouter.Router
}

func (weather *Weather) FirstEndpoint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	auth := authorizer.FromContext(r.Context())
	if auth == nil {
		weather.SetResponse(http.StatusInternalServerError, "Internal Server Error", w)
		return
	}
	responseMessage := fmt.Sprintf("Welcome %s! This is the first endpoint working!", auth.GenerateAccessToken())
	weather.SetResponse(http.StatusOK, responseMessage, w)
}

func (weather *Weather) Register(router *httprouter.Router) {
	weather.router = router
	weather.router.GET("/first_endpoint/", weather.FirstEndpoint)
}
