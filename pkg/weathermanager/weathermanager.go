package weathermanager

import "context"

type WeatherManager interface {
	SaveWeather(string, map[string]int) error
	GetWeather(string) ([]map[string]int, error)
	DeleteWeather(string) error
}

type MainWeatherManager struct {
	validToken string
	weathers   map[string][]map[string]int
}

func (m *MainWeatherManager) SaveWeather(city string, temperatures map[string]int) error {
	m.weathers[city] = []map[string]int{
		temperatures,
	}
	return nil
}

func (m *MainWeatherManager) GetWeather(city string) ([]map[string]int, error) {
	return m.weathers[city], nil
}

func (m *MainWeatherManager) DeleteWeather(city string) error {
	delete(m.weathers, city)
	return nil
}

func New() *MainWeatherManager {
	return &MainWeatherManager{
		weathers: map[string][]map[string]int{},
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
