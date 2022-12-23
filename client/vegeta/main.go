package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	targetApiUrl  = "http://localhost:8888/send"
	reportsFolder = "./client/vegeta/reports"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	rate := vegeta.Rate{Freq: 150, Per: time.Second}
	duration := 60 * time.Second
	targeter := customTargeter()
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Test service insert load") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewHDRHistogramPlotReporter(&metrics)

	f, err := os.Create(fmt.Sprintf("%s/report_%s.txt", reportsFolder, time.Now().Format("2006-01-02_15:04:05")))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := reporter.Report(f); err != nil {
		fmt.Printf("vegeta report write: %s\n", err)
	}
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func customTargeter() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "POST"
		tgt.URL = targetApiUrl

		payload := fmt.Sprintf(`{
            "timestamp" : %d,
			"value": %d
          }`, time.Now().Unix(), rand.Intn(1000)+1)

		tgt.Body = []byte(payload)

		header := http.Header{}
		header.Add("Accept", "application/json")
		header.Add("Content-Type", "application/json")
		tgt.Header = header

		return nil
	}
}
