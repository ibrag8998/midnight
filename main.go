package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"midnight/api"
)

type config struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Year      int     `json:"year"`
	OutputCSV string  `json:"outputCSV"`
}

func main() {
	// read config
	jsonFile, err := os.Open("config.json")
	if err != nil {
		panic(fmt.Errorf("Error opening config.json: %v", err))
	}
	defer jsonFile.Close()

	// parse json
	byteValue, _ := io.ReadAll(jsonFile)
	var cfg config
	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		panic(fmt.Errorf("Error parsing config.json: %v", err))
	}

	// get year data
	prayerApi := api.NewPrayerTimeAPI()

	prayerTimes := prayerApi.GetYearData(cfg.Year, cfg.Latitude, cfg.Longitude)

	err = prayerTimes.WriteCSV(cfg.OutputCSV)
	if err != nil {
		panic(fmt.Errorf("Error writing prayer times csv: %v", err))
	}

	fmt.Println("OK")
}
