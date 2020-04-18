package resources

import (
	"fmt"
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

type messageResponseModel struct {
	Message string `json:"message"`
}

type weatherReportResponseModel struct {
	WeatherReport map[string]interface{} `json:"weather_report"`
}

type saveWeatherReportRequestModel struct {
	City    string `json:"city"`
	Weather []struct {
		Date        string `json:"date"`
		Temperature int    `json:"temperature"`
	} `json:"weather"`
}

type getWeatherReportRequestModel struct {
	City        string `json:"city"`
	InitialDate string `json:"initial_date"`
	EndDate     string `json:"end_date"`
}

type deleteWeatherReportRequestModel struct {
	City string `json:"city"`
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
		e := internalerror.New(fmt.Sprintf("Error validating auth token (%s)", err.Error()))
		weather.SetResponse(http.StatusUnauthorized, e, w)
		return
	}

	var requestModel saveWeatherReportRequestModel
	err = weather.ParseFromBody(r, &requestModel)
	if err != nil {
		e := internalerror.New(fmt.Sprintf("Error parsing request body (%s)", err.Error()))
		weather.SetResponse(http.StatusInternalServerError, e, w)
		return
	}

	weatherReport := map[string]int{}
	for _, o := range requestModel.Weather {
		weatherReport[o.Date] = o.Temperature
	}

	err = weatherMgr.SaveWeather(requestModel.City, weatherReport)
	if err != nil {
		e := internalerror.New(fmt.Sprintf("Error saving weather (%s)", err.Error()))
		weather.SetResponse(http.StatusBadRequest, e, w)
		return
	}

	weather.SetResponse(http.StatusOK, messageResponseModel{
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
		e := internalerror.New(fmt.Sprintf("Error validating auth token (%s)", err.Error()))
		weather.SetResponse(http.StatusUnauthorized, e, w)
		return
	}

	var requestModel getWeatherReportRequestModel
	err = weather.ParseFromBody(r, &requestModel)
	if err != nil {
		e := internalerror.New(fmt.Sprintf("Error parsing request body (%s)", err.Error()))
		weather.SetResponse(http.StatusInternalServerError, e, w)
		return
	}

	vancouverWeather, err := weatherMgr.GetWeather(
		requestModel.City,
		requestModel.InitialDate,
		requestModel.EndDate,
	)
	if err != nil {
		e := internalerror.New(err.Error())
		weather.SetResponse(http.StatusBadRequest, e, w)
		return
	}

	weather.SetResponse(http.StatusOK, weatherReportResponseModel{
		WeatherReport: map[string]interface{}{
			requestModel.City: vancouverWeather,
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
		e := internalerror.New(fmt.Sprintf("Error validating auth token (%s)", err.Error()))
		weather.SetResponse(http.StatusUnauthorized, e, w)
		return
	}

	var requestModel deleteWeatherReportRequestModel
	err = weather.ParseFromBody(r, &requestModel)
	if err != nil {
		e := internalerror.New(fmt.Sprintf("Error parsing request body (%s)", err.Error()))
		weather.SetResponse(http.StatusInternalServerError, e, w)
		return
	}

	err = weatherMgr.DeleteWeather(requestModel.City)
	if err != nil {
		weather.SetResponse(http.StatusInternalServerError, err, w)
		return
	}

	weather.SetResponse(http.StatusOK, messageResponseModel{
		"The weather was deleted succesfully!",
	}, w)
}

func (weather *Weather) Register(router *httprouter.Router) {
	weather.router = router
	weather.router.POST("/weather/", weather.SaveCityWeather)
	weather.router.GET("/weather/", weather.GetCityWeather)
	weather.router.DELETE("/weather/", weather.DeleteCityWeather)
}
