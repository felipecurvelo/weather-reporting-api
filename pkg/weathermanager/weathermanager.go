package weathermanager

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type WeatherManager interface {
	SaveWeather(string, map[string]int) error
	GetWeather(string, string, string) (map[string]int, error)
	DeleteWeather(string) error
}

type MainWeatherManager struct {
	validToken string
	weathers   map[string]map[string]int
}

func (m *MainWeatherManager) SaveWeather(city string, temperatures map[string]int) error {
	dateLayout := "2006-01-02"
	for k, _ := range temperatures {
		_, err := time.Parse(dateLayout, k)
		if err != nil {
			return fmt.Errorf("Invalid date %s (%s)", k, err.Error())
		}
	}

	m.weathers[strings.ToLower(city)] = temperatures
	return nil
}

func (m *MainWeatherManager) GetWeather(city string, initialDate string, endDate string) (map[string]int, error) {
	if city == "" {
		return nil, fmt.Errorf("Empty city")
	}

	if initialDate == "" {
		return nil, fmt.Errorf("Empty initial date")
	}

	if endDate == "" {
		return nil, fmt.Errorf("Empty end date")
	}

	_, ok := m.weathers[strings.ToLower(city)]
	if !ok {
		return nil, fmt.Errorf("Weather report not found")
	}

	dateLayout := "2006-01-02"
	initial, err := time.Parse(dateLayout, initialDate)
	if err != nil {
		return nil, fmt.Errorf("Invalid initial date (%s)", err.Error())
	}

	end, err := time.Parse(dateLayout, endDate)
	if err != nil {
		return nil, fmt.Errorf("Invalid end date (%s)", err.Error())
	}

	if !initial.Before(end) {
		return nil, fmt.Errorf("Invalid date range")
	}

	temperatures := map[string]int{}
	for k, e := range m.weathers[strings.ToLower(city)] {
		temperatureDate, _ := time.Parse(dateLayout, k)
		if temperatureDate.After(initial) && temperatureDate.Before(end) {
			temperatures[k] = e
		}
	}

	return temperatures, nil
}

func (m *MainWeatherManager) DeleteWeather(city string) error {
	delete(m.weathers, strings.ToLower(city))
	return nil
}

func New() *MainWeatherManager {
	return &MainWeatherManager{
		weathers: map[string]map[string]int{},
	}
}

type contextKey struct{}

func FromContext(ctx context.Context) WeatherManager {
	auth, _ := ctx.Value(contextKey{}).(WeatherManager)
	return auth
}

func NewContext(parentContext context.Context, weatherMgr WeatherManager) context.Context {
	return context.WithValue(parentContext, contextKey{}, weatherMgr)
}
