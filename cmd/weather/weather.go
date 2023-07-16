package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
	"yandex/internal/db"
	"yandex/internal/telegram"
)

type Weather struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	CurrentWeather       struct {
		Temperature   float64 `json:"temperature"`
		Windspeed     float64 `json:"windspeed"`
		Winddirection float64 `json:"winddirection"`
		Weathercode   int     `json:"weathercode"`
		IsDay         int     `json:"is_day"`
		Time          int     `json:"time"`
	} `json:"current_weather"`
	DailyUnits struct {
		Time             string `json:"time"`
		Weathercode      string `json:"weathercode"`
		Temperature2MMax string `json:"temperature_2m_max"`
		Temperature2MMin string `json:"temperature_2m_min"`
		Sunrise          string `json:"sunrise"`
		Sunset           string `json:"sunset"`
	} `json:"daily_units"`
	Daily struct {
		Time             []int     `json:"time"`
		Weathercode      []int     `json:"weathercode"`
		Temperature2MMax []float64 `json:"temperature_2m_max"`
		Temperature2MMin []float64 `json:"temperature_2m_min"`
		Sunrise          []int     `json:"sunrise"`
		Sunset           []int     `json:"sunset"`
	} `json:"daily"`
}

var weatherCode = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "partly cloudy",
	3:  "overcast",
	45: "Fog",
	48: "depositing rime fog",
	51: "Drizzle: Light",
	53: "Drizzle: moderate",
	55: "Drizzle: dense intensity",
	56: "Freezing Drizzle: Light",
	57: "Freezing Drizzle:  dense intensity",
	61: "Rain: Slight",
	63: "Rain: moderate",
	65: "Rain: heavy intensity",
	66: "Freezing Rain: Light",
	67: "Freezing Rain: heavy intensity",
	71: "Snow fall: Slight",
	73: "Snow fall: moderate",
	75: "Snow fall: heavy intensity",
	77: "Snow grains",
	80: "Rain showers: Slight",
	81: "Rain showers: moderate",
	82: "Rain showers: violent",
	85: "Snow showers slight",
	86: "Snow showers heavy",
	95: "Thunderstorm: Slight or moderate",
	96: "Thunderstorm with slight",
	99: "Thunderstorm with heavy hail",
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	get, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=55.6302&longitude=37.6045&daily=weathercode,temperature_2m_max,temperature_2m_min,sunrise,sunset&current_weather=true&windspeed_unit=ms&timeformat=unixtime&timezone=Europe%2FMoscow&forecast_days=3")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer get.Body.Close()

	getResp, _ := io.ReadAll(get.Body)
	weather := Weather{}
	err = json.Unmarshal(getResp, &weather)
	if err != nil {
		log.Fatal(err)
		return
	}

	loc, _ := time.LoadLocation(weather.Timezone)
	sunrise := time.Unix(int64(weather.Daily.Sunrise[0]), 0)
	sunset := time.Unix(int64(weather.Daily.Sunset[0]), 0)

	var message string
	message += fmt.Sprint(time.Now().In(loc).Format("02.01.2006")) + " \n"
	message += "Текущая температура: " + fmt.Sprint(weather.CurrentWeather.Temperature) + "°C \n"
	message += "Погода: " + fmt.Sprint(weatherCode[weather.CurrentWeather.Weathercode]) + " \n"
	message += "Скорость ветра: " + fmt.Sprint(weather.CurrentWeather.Windspeed) + " m/s \n"
	message += "Максимальная температура: " + fmt.Sprint(weather.Daily.Temperature2MMax[0]) + "°C \n"
	message += "Минимальная температура: " + fmt.Sprint(weather.Daily.Temperature2MMin[0]) + "°C \n"
	message += "Расвет: " + fmt.Sprint(sunrise.Format("15:04")) + " \n"
	message += "Закат: " + fmt.Sprint(sunset.Format("15:04")) + " \n"
	message += "Световой день: " + fmt.Sprint(sunset.Sub(sunrise))

	wg := sync.WaitGroup{}
	users := db.GetDigitalUser()
	for _, userId := range users {
		wg.Add(1)
		go telegram.SendTelegram(userId, url.QueryEscape(message), &wg)
	}

	wg.Wait()
}
