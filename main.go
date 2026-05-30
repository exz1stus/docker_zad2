package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	port       = "8080"
	author     = "Denys Petrov"
	API_ORIGIN = "https://api.open-meteo.com/v1/forecast"
)

type City struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Country struct {
	Country string `json:"country"`
	Cities  []City `json:"cities"`
}

type GeoData struct {
	Countries []Country `json:"countries"`
}

var geoData GeoData

func weatherPage(w http.ResponseWriter, r *http.Request) {
	selCountry := r.URL.Query().Get("country")
	selCity := r.URL.Query().Get("city")

	//Szablon css pobrany z Internetu
	fmt.Fprint(w, `<html><head><style>
		body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 40px; background: #f0f2f5; }
		.container { max-width: 500px; margin: auto; background: white; padding: 30px; border-radius: 15px;  }
		select, button { width: 100%; padding: 12px; margin: 10px 0; border-radius: 8px; border: 1px solid #ccc; font-size: 1rem; }
		button { background: #007bff; color: white; border: none; cursor: pointer; font-weight: bold; }
		button:hover { background: #0056b3; }
		.result { margin-top: 20px; padding: 15px; background: #e7f3ff; border-radius: 10px; border-left: 5px solid #007bff; }
		a { color: #666; text-decoration: none; font-size: 0.9rem; }
	</style></head><body><div class="container">`)

	fmt.Fprintf(w, "<h2>Pogoda</h2>")

	fmt.Fprint(w, `<form method="GET">`)
	fmt.Fprint(w, `<label>1. Select Country:</label>`)
	fmt.Fprint(w, `<select name="country" onchange="this.form.submit()">`)
	fmt.Fprint(w, `<option value="">-- Choose Country --</option>`)
	for _, c := range geoData.Countries {
		selected := ""
		if c.Country == selCountry {
			selected = "selected"
		}
		fmt.Fprintf(w, `<option value="%s" %s>%s</option>`, c.Country, selected, c.Country)
	}
	fmt.Fprint(w, `</select>`)

	if selCountry != "" {
		fmt.Fprint(w, `<label>2. Select City:</label>`)
		fmt.Fprint(w, `<select name="city">`)
		fmt.Fprint(w, `<option value="">-- Choose City --</option>`)

		for _, c := range geoData.Countries {
			if c.Country == selCountry {
				for _, city := range c.Cities {
					selected := ""
					if city.Name == selCity {
						selected = "selected"
					}
					fmt.Fprintf(w, `<option value="%s" %s>%s</option>`, city.Name, selected, city.Name)
				}
			}
		}
		fmt.Fprint(w, `</select>`)
		fmt.Fprint(w, `<button type="submit">Get Weather</button>`)
	}
	fmt.Fprint(w, `</form>`)

	if selCountry != "" && selCity != "" {
		weather, err := getWeather(selCountry, selCity)
		time, err := time.Parse("2006-01-02T15:04", weather.Time)
		if err != nil {
			fmt.Fprintf(w, `<p style="color:red">Error: %v</p>`, err)
		} else {
			fmt.Fprintf(w, `
				<div class="result">
					<h3>%s, %s</h3>
					<p><b>Temp:</b> %.1f°C</p>
					<p><b>Wind:</b> %.1f km/h</p>
					<p><small>Updated: %s</small></p>
				</div>`, selCity, selCountry, weather.Temperature, weather.WindSpeed, time)
		}
	}

	fmt.Fprint(w, `<br><center><a href="/">Clear All</a></center>`)
	fmt.Fprint(w, `</div></body></html>`)
}

type Current struct {
	Time        string  `json:"time"`
	Temperature float32 `json:"temperature_2m"`
	WindSpeed   float32 `json:"wind_speed_10m"`
}

type ApiWeatherResponse struct {
	Current Current `json:"current"`
}

func getGeoCoords(countryName string, cityName string) (float64, float64, error) {
	for _, c := range geoData.Countries {
		if c.Country == countryName {
			for _, city := range c.Cities {
				if city.Name == cityName {
					return city.Latitude, city.Longitude, nil
				}
			}
		}
	}
	return 0, 0, fmt.Errorf("failed to find place")
}

func getWeather(country string, city string) (*Current, error) {
	requestedData := "temperature_2m,wind_speed_10m"
	latitude, longitude, err := getGeoCoords(country, city)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&current=%s", API_ORIGIN, latitude, longitude, requestedData)
	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Error fetching API: %v", err)
	}
	defer res.Body.Close()

	var data ApiWeatherResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("Error decoding JSON:%v", err)
	}

	return &data.Current, nil
}

func main() {
	//reczny healtcheck, zeby nie pobierac curla
	if len(os.Args) > 1 && os.Args[1] == "health" {
		_, err := http.Get("http://localhost:8080/")
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	fmt.Printf("Data uruchomienia: %s\n", time.Now().Format(time.DateTime))
	fmt.Printf("Autor: %s\n", author)
	fmt.Printf("Port TCP: %s\n", port)

	content, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	err = json.Unmarshal(content, &geoData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	http.HandleFunc("/", weatherPage)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
