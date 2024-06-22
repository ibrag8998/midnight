package api

import (
	"encoding/csv"
	"os"
)

type MidnightTimes map[string]string

func (m MidnightTimes) WriteCSV(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for k, v := range m {
		err = w.Write([]string{k, v})
		if err != nil {
			return err
		}
	}

	return nil
}

type Resp struct {
	Data []RespData `json:"data"`
}

type RespData struct {
	Timings RespTimings `json:"timings"`
	Date    RespDate    `json:"date"`
}

type RespTimings struct {
	Midnight string `json:"Midnight"`
}

type RespDate struct {
	Gregorian RespGregorianDate `json:"gregorian"`
}

type RespGregorianDate struct {
	Date string `json:"date"`
}
