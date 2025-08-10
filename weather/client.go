 
// weather/client.go
package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Structs pour décoder les réponses JSON de l'API.
// Elles ne contiennent que les champs qui nous intéressent.

// CurrentWeatherData représente la météo actuelle.
type CurrentWeatherData struct {
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Name string `json:"name"`
}

// ForecastData représente les prévisions météo.
type ForecastData struct {
	List []struct {
		Dt   int64 `json:"dt"` // Date et heure de la prévision (timestamp)
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"list"`
}

// Client gère la communication avec l'API OpenWeatherMap.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient crée un nouveau client météo.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetCurrentWeather récupère la météo actuelle pour une ville donnée.
func (c *Client) GetCurrentWeather(city string) (*CurrentWeatherData, error) {
	apiURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=fr", url.QueryEscape(city), c.apiKey)

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la requête API : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a retourné un statut inattendu : %s", resp.Status)
	}

	var data CurrentWeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de la réponse JSON : %w", err)
	}

	return &data, nil
}

// GetForecast récupère les prévisions pour les 5 prochains jours.
func (c *Client) GetForecast(city string) (*ForecastData, error) {
    // L'API de prévision renvoie des données toutes les 3 heures. Nous allons filtrer pour n'en garder qu'une par jour.
	apiURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric&lang=fr", url.QueryEscape(city), c.apiKey)

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la requête API de prévision : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API de prévision a retourné un statut inattendu : %s", resp.Status)
	}

	var data ForecastData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de la réponse JSON de prévision : %w", err)
	}

    // Filtrer pour n'avoir qu'une prévision par jour (celle de midi)
    dailyForecasts := ForecastData{}
    seenDays := make(map[string]bool)
    for _, forecast := range data.List {
        day := time.Unix(forecast.Dt, 0).Format("2006-01-02")
        // On ne garde que la prévision la plus proche de midi.
        if !seenDays[day] && time.Unix(forecast.Dt, 0).Hour() >= 12 {
            dailyForecasts.List = append(dailyForecasts.List, forecast)
            seenDays[day] = true
        }
    }


	return &dailyForecasts, nil
}

// GetIcon récupère l'image de l'icône météo.
func (c *Client) GetIcon(iconCode string) ([]byte, error) {
	iconURL := fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", iconCode)
	resp, err := c.httpClient.Get(iconURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
