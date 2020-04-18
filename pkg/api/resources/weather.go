package resources

import (
	"net/http"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/julienschmidt/httprouter"
)

type Weather struct {
	api.ResourceBase
	router *httprouter.Router
}

func (weather *Weather) FirstEndpoint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	responseMessage := "Welcome! This is the first endpoint working!"
	weather.SetResponse(http.StatusOK, responseMessage, w)
}

func (res *Weather) Register(router *httprouter.Router) {
	res.router = router
	res.router.GET("/first_endpoint/", res.FirstEndpoint)
}
