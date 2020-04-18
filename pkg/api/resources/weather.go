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

type saveWeatherResponseModel struct {
	Message string `json:"message"`
}

type getWeatherResponseModel struct {
	WeatherReport map[string]interface{} `json:"weather_report"`
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

	weather.SetResponse(http.StatusOK, saveWeatherResponseModel{
		"The weather was saved succesfully!",
	}, w)
}

func (weather *Weather) GetCityWeather(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	city := "vancouver"
	vancouverWeather, err := weatherMgr.GetWeather(city)
	if err != nil {
		e := internalerror.New(err.Error())
		weather.SetResponse(http.StatusBadRequest, e, w)
		return
	}

	weather.SetResponse(http.StatusOK, getWeatherResponseModel{
		WeatherReport: map[string]interface{}{
			city: vancouverWeather,
		},
	}, w)
}

func (weather *Weather) DeleteCityWeather(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	err = weatherMgr.DeleteWeather("vancouver")
	if err != nil {
		weather.SetResponse(http.StatusInternalServerError, err, w)
		return
	}

	weather.SetResponse(http.StatusOK, saveWeatherResponseModel{
		"The weather was deleted succesfully!",
	}, w)
}

func (weather *Weather) Register(router *httprouter.Router) {
	weather.router = router
	weather.router.POST("/weather/", weather.SaveCityWeather)
	weather.router.GET("/weather/", weather.GetCityWeather)
	weather.router.DELETE("/weather/", weather.DeleteCityWeather)
}
