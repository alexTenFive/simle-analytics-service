package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	targetApiUrl  = "http://localhost:8888/send"
	reportsFolder = "./cmd/vegeta/reports"
)

var requestsPerSec = flag.Int("rps", 50, "specify requests per second")
var dur = flag.Int("duration", 60, "duration in seconds")

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	rate := vegeta.Rate{Freq: *requestsPerSec, Per: time.Second}
	duration := time.Duration(*dur) * time.Second
	targeter := customTargeter()
	attacker := vegeta.NewAttacker()
	fmt.Printf("starting test with: %s rate\nduration: %s\n", rate, duration)

	var metrics vegeta.Metrics
	f, err := os.Create(fmt.Sprintf("%s/report.bin", reportsFolder))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	encoder := vegeta.NewEncoder(f)

	for res := range attacker.Attack(targeter, rate, duration, "Test service insert load") {
		metrics.Add(res)
		if err = encoder.Encode(res); err != nil {
			fmt.Printf("encode: %s\n", err)
		}
	}
	metrics.Close()

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
