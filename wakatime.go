package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"
)

type wakatimeDataRange struct {
	Start    time.Time `json:"start,omitempty"`
	End      time.Time `json:"end,omitempty"`
	Date     string    `json:"date,omitempty"`
	Text     string    `json:"text,omitempty"`
	TimeZone string    `json:"timeZone,omitempty"`
}

type wakatimeDataGrandTotal struct {
	Hours        int    `json:"hours,omitempty"`
	Minutes      int    `json:"minutes,omitempty"`
	TotalSeconds int    `json:"totalSeconds,omitempty"`
	Digital      string `json:"digital,omitempty"`
	Decimal      string `json:"decimal,omitempty"`
	Text         string `json:"text,omitempty"`
}

type wakatimeData struct {
	Range      wakatimeDataRange      `json:"range,omitempty"`
	GrandTotal wakatimeDataGrandTotal `json:"grand_total,omitempty"`
}

func fetchWakatimeActivity(url string) (*wakatimeData, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from wakatime: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data from wakatime: %v", res.Status)
	}

	var body map[string][]wakatimeData
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("failed to decode data from wakatime: %v", err)
	}

	if len(body["data"]) == 0 {
		return nil, fmt.Errorf("no data found from wakatime")
	}

	data := body["data"]

	todayData := slices.MaxFunc(data, func(i, j wakatimeData) int {
		return cmp.Compare(i.Range.Start.Unix(), j.Range.Start.Unix())
	})

	return &todayData, nil
}
