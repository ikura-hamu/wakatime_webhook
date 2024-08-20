package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func main() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Printf("Error: failed to load location: %v\n", err)
		os.Exit(1)
	}

	conf := loadConfig()

	s, _ := gocron.NewScheduler(gocron.WithLocation(jst))
	defer func() { _ = s.Shutdown() }()

	_, _ = s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(0, 10, 0),
			),
		),
		gocron.NewTask(
			action, conf,
		),
	)

	s.Start()

	select {}
}

type config struct {
	WebhookID   string `json:"webhookID"`
	Secret      string `json:"secret"`
	WakatimeURL string `json:"wakatimeURL"`
}

func action(conf config) {
	data, err := fetchWakatimeActivity(conf.WakatimeURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	grassCount := data.GrandTotal.Hours + data.GrandTotal.Minutes/30

	text := fmt.Sprintf(`%s
**%s**
%s%s`, data.Range.Date, data.GrandTotal.Text, strings.Repeat(":0x30a14e:", grassCount), strings.Repeat(":0xcccccc:", 24-grassCount))

	err = postTraqWebhook(conf.WebhookID, conf.Secret, text)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(text)

}

func loadConfig() config {
	webhookID, ok := os.LookupEnv("TRAQ_WEBHOOK_ID")
	if !ok {
		fmt.Println("TRAQ_WEBHOOK_ID is not set")
		os.Exit(1)
	}

	secret, ok := os.LookupEnv("TRAQ_WEBHOOK_SECRET")
	if !ok {
		fmt.Println("TRAQ_WEBHOOK_SECRET is not set")
		os.Exit(1)
	}

	wakatimeURL, ok := os.LookupEnv("WAKATIME_URL")
	if !ok {
		fmt.Println("WAKATIME_URL is not set")
		os.Exit(1)
	}

	return config{
		WebhookID:   webhookID,
		Secret:      secret,
		WakatimeURL: wakatimeURL,
	}
}
