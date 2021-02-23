package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type locationData struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

func main() {
	fmt.Println("")

	formTemplate := template.Must(template.ParseFiles("form.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			formTemplate.Execute(w, nil)
			return
		}
		url := "http://api.weatherapi.com/v1/current.json?key=8b8d0650a0b74fa088c94105210402&q=" + r.FormValue("city_name")
		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		} else {
			// fmt.Println(res.Body)
			resbyteValue, _ := ioutil.ReadAll(res.Body)
			// fmt.Println(resbyteValue)
			var responseData locationData
			json.Unmarshal(resbyteValue, &responseData)
			// fmt.Println(responseData.Location.Lat)
			var data = struct {
				Success bool
				CelTemp float64
				FarTemp float64
				Lat     float64
				Lon     float64
			}{true, responseData.Current.FeelslikeC, responseData.Current.FeelslikeF, responseData.Location.Lat, responseData.Location.Lon}
			formTemplate.Execute(w, data)
		}
	})

	http.ListenAndServe(":8080", nil)

}
