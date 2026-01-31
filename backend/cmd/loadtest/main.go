package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Config
var (
	baseURL     = "http://localhost:3030"
	dropID      = uint64(1) // Default drop ID
	concurrency = 50        // Number of concurrent workers
	duration    = 10 * time.Second
)

func main() {
	flag.StringVar(&baseURL, "url", "http://localhost:3030", "Base URL of the API")
	flag.Uint64Var(&dropID, "drop", 1, "ID of the drop to purchase")
	flag.IntVar(&concurrency, "c", 50, "Number of concurrent workers")
	flag.DurationVar(&duration, "d", 10*time.Second, "Duration of the test")
	flag.Parse()

	log.Printf("🚀 Starting Load Test on %s/api/drops/%d/purchase", baseURL, dropID)
	log.Printf("Workers: %d, Duration: %s", concurrency, duration)

	var (
		successCount uint64
		failCount    uint64
		totalReqs    uint64
	)

	start := time.Now()
	timeout := time.After(duration)
	done := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			defer wg.Done()
			client := &http.Client{
				Transport: &http.Transport{
					MaxIdleConns:        100,
					MaxIdleConnsPerHost: 100,
				},
				Timeout: 5 * time.Second,
			}

			for {
				select {
				case <-done:
					return
				default:
					// Prepare request
					reqBody := map[string]interface{}{
						"quantity": 1,
						"name":     fmt.Sprintf("LoadTest User %d", workerID),
						"phone":    "0987654321",
						"email":    fmt.Sprintf("user%d@example.com", workerID),
						"address":  "123 Load Test St",
						"province": "Hanoi",
						"district": "Cau Giay",
						"ward":     "Dich Vong",
					}
					jsonBody, _ := json.Marshal(reqBody)

					req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/drops/%d/purchase", baseURL, dropID), bytes.NewBuffer(jsonBody))
					req.Header.Set("Content-Type", "application/json")

					resp, err := client.Do(req)
					atomic.AddUint64(&totalReqs, 1)
					if err != nil {
						atomic.AddUint64(&failCount, 1)
						continue
					}

					bodyBytes, _ := io.ReadAll(resp.Body)
					resp.Body.Close()

					if resp.StatusCode == 200 {
						// Parse response to get OrderCode
						var res struct {
							OrderCode int64 `json:"order_code"`
						}
						
						if err := json.Unmarshal(bodyBytes, &res); err != nil {
							// log.Printf("Parse error: %v", err)
						}
						
						// Trigger Webhook immediately to race for stock
						if res.OrderCode != 0 {
						    // Fire Webhook
						    webhookBody := map[string]interface{}{
						        "code": "00",
						        "data": map[string]interface{}{
						            "orderCode": res.OrderCode,
						            "amount": 10000,
						            "status": "PAID",
						        },
						    }
						    wbBytes, _ := json.Marshal(webhookBody)
						    wbReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/limited-drops/webhook/payos", baseURL), bytes.NewBuffer(wbBytes))
						    wbReq.Header.Set("Content-Type", "application/json")
						    
						    wbResp, err := client.Do(wbReq)
						    if err == nil {
						        io.Copy(io.Discard, wbResp.Body)
						        wbResp.Body.Close()
						        
						        if wbResp.StatusCode == 200 {
						             // Webhook Success = REAL Success (Stock claimed or Handle Success)
						             // Wait, if 200 OK, it means ProcessSuccessfulDropPayment returned nil.
						             atomic.AddUint64(&successCount, 1)
						        } else {
						             // 500 Error = Sold Out
						             atomic.AddUint64(&failCount, 1)
						        }
						    } else {
						        atomic.AddUint64(&failCount, 1)
						    }
						} else {
						    // Failed to parse order code, count as fail
						    atomic.AddUint64(&failCount, 1)
						}
					} else {
						// log.Printf("Status: %d", resp.StatusCode)
						atomic.AddUint64(&failCount, 1)
					}
				}
			}
		}(i)
	}

	<-timeout
	close(done)
	wg.Wait()
	elapsed := time.Since(start)

	log.Printf("🏁 Load Test Finished in %s", elapsed)
	log.Printf("Total Requests: %d", totalReqs)
	log.Printf("Success: %d (%.2f%%)", successCount, float64(successCount)/float64(totalReqs)*100)
	log.Printf("Failed: %d (%.2f%%)", failCount, float64(failCount)/float64(totalReqs)*100)
	log.Printf("RPS: %.2f", float64(totalReqs)/elapsed.Seconds())
}
