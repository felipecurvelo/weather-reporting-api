package resources

import (
	"net/http"

	"github.com/felipecurvelo/weather-reporting-api/pkg/authorizer"
	"github.com/felipecurvelo/weather-reporting-api/pkg/weathermanager"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/felipecurvelo/weather-reporting-api/pkg/internalerror"
	"github.com/julienschmidt/httprouter"
)

type Weather struct {
	api.ResourceBase
	router *httprouter.Router
}

type weatherResponseModel struct {
	Message string `json:"message"`
}

func (weather *Weather) SaveCityWeather(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	auth := authorizer.FromContext(ctx)
	if auth == nil {
		e := internalerror.New("Internal Server Error")
		weather.SetResponse(http.StatusInternalServerError, e, w)
		return
	}

	weatherMgr := weathermanager.FromContext(ctx)
	if weatherMgr == nil {
		e := internalerror.New("Internal Server Error")
		weather.SetResponse(http.StatusInternalServerError, e, w)
		return
	}

	err := weather.ValidateAuthToken(ctx, r)
	if err != nil {
		weather.SetResponse(http.StatusInternalServerError, err, w)
		return
	}

	weatherMgr.SaveWeather("vancouver", map[string]int{
		"2020-0-18": 15,
	})

	weather.SetResponse(http.StatusOK, weatherResponseModel{
		"The weather was saved succesfully!",
	}, w)
}

func (weather *Weather) Register(router *httprouter.Router) {
	weather.router = router
	weather.router.POST("/save/", weather.SaveCityWeather)
}
