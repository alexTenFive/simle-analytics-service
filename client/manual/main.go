package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const apiBase = "http://localhost:8888"

var requests = flag.Int("requests", 10000, "request per minute")

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	requestPerTime := time.Minute / time.Duration(*requests)

	client := http.Client{
		Timeout: time.Second,
	}

	processedRequests := 0
	start := time.Now()
	defer func() {
		fmt.Printf("processed requests=%d; time=%s\n", processedRequests, time.Since(start))
	}()
	for i := 0; i < *requests; i++ {
		after := time.After(requestPerTime)

		ba, _ := json.Marshal(map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"value":     rand.Intn(1000) + 1,
		})
		body := bytes.NewBuffer(ba)
		resp, err := client.Post(fmt.Sprintf("%s/send", apiBase), "application/json", body)
		if err != nil {
			fmt.Printf("request err: %s\n", err)
			continue
		}
		data := make(map[string]interface{}, 1)
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			fmt.Printf("unmarshal: %s\n", err)
		}
		processedRequests++
		<-after
	}

}
