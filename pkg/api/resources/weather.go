package resources

import (
	"fmt"
	"net/http"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/felipecurvelo/weather-reporting-api/pkg/auth"
	"github.com/julienschmidt/httprouter"
)

type Weather struct {
	api.ResourceBase
	router *httprouter.Router
}

func (weather *Weather) FirstEndpoint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	a := auth.FromContext(r.Context())
	responseMessage := fmt.Sprintf("Welcome %s! This is the first endpoint working!", a.Name)
	weather.SetResponse(http.StatusOK, responseMessage, w)
}

func (res *Weather) Register(router *httprouter.Router) {
	res.router = router
	res.router.GET("/first_endpoint/", res.FirstEndpoint)
}
