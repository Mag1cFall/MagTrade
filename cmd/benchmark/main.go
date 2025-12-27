// 效能基準測試工具
//
// 本檔案是獨立的 HTTP 壓力測試工具
// 用於測試 API 端點的 QPS、延遲、成功率等指標
// 執行：go run cmd/benchmark/main.go -c 100 -n 1000
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Result 單次請求結果
type Result struct {
	Success  bool
	Duration time.Duration
	Status   int
}

// Stats 統計結果
type Stats struct {
	Total     int64
	Success   int64
	Failed    int64
	TotalTime time.Duration
	MinTime   time.Duration
	MaxTime   time.Duration
	Latencies []time.Duration // 用於計算百分位
}

func main() {
	baseURL := flag.String("url", "http://localhost:8080", "Base URL")
	concurrency := flag.Int("c", 100, "Concurrency level")
	requests := flag.Int("n", 1000, "Total requests")
	flag.Parse()

	fmt.Println("MagTrade Benchmark Tool (Go)")
	fmt.Printf("Target: %s\n", *baseURL)
	fmt.Printf("Concurrency: %d\n", *concurrency)
	fmt.Printf("Total Requests: %d\n", *requests)
	fmt.Println()

	// 測試端點列表
	endpoints := []struct {
		Name string
		Path string
	}{
		{"Health Check", "/health"},
		{"Products", "/api/v1/products"},
		{"Flash Sales", "/api/v1/flash-sales"},
		{"Stock Check", "/api/v1/flash-sales/1/stock"},
	}

	results := make(map[string]*Stats)

	for _, ep := range endpoints {
		url := *baseURL + ep.Path
		stats := benchmark(url, *concurrency, *requests)
		results[ep.Name] = stats
		printStats(ep.Name, stats)
		fmt.Println()
	}

	// 總結報告
	fmt.Println("==========================================")
	fmt.Println("  Summary")
	fmt.Println("==========================================")
	fmt.Println()
	fmt.Printf("| %-15s | %8s | %8s | %8s | %8s | %6s |\n", "Endpoint", "QPS", "Avg(ms)", "P95(ms)", "P99(ms)", "Failed")
	fmt.Println("|-----------------|----------|----------|----------|----------|--------|")

	for _, ep := range endpoints {
		s := results[ep.Name]
		qps := float64(s.Success) / s.TotalTime.Seconds()
		avg := float64(0)
		p95 := time.Duration(0)
		p99 := time.Duration(0)

		if len(s.Latencies) > 0 {
			sort.Slice(s.Latencies, func(i, j int) bool {
				return s.Latencies[i] < s.Latencies[j]
			})

			var sum time.Duration
			for _, l := range s.Latencies {
				sum += l
			}
			avg = float64(sum.Milliseconds()) / float64(len(s.Latencies))
			p95 = s.Latencies[int(float64(len(s.Latencies))*0.95)]
			p99 = s.Latencies[int(float64(len(s.Latencies))*0.99)]
		}

		fmt.Printf("| %-15s | %8.2f | %8.2f | %8d | %8d | %6d |\n",
			ep.Name, qps, avg, p95.Milliseconds(), p99.Milliseconds(), s.Failed)
	}
	fmt.Println()
}

// benchmark 執行壓力測試
func benchmark(url string, concurrency, total int) *Stats {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        concurrency * 2,
			MaxIdleConnsPerHost: concurrency * 2,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	var wg sync.WaitGroup
	resultChan := make(chan Result, total)

	// 分配請求數
	requestsPerWorker := total / concurrency
	remainder := total % concurrency

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		reqs := requestsPerWorker
		if i < remainder {
			reqs++
		}

		wg.Add(1)
		go func(numReqs int) {
			defer wg.Done()
			for j := 0; j < numReqs; j++ {
				result := makeRequest(client, url)
				resultChan <- result
			}
		}(reqs)
	}

	wg.Wait()
	totalTime := time.Since(start)
	close(resultChan)

	// 統計結果
	stats := &Stats{
		TotalTime: totalTime,
		MinTime:   time.Hour,
		Latencies: make([]time.Duration, 0, total),
	}

	for r := range resultChan {
		stats.Total++
		if r.Success {
			stats.Success++
			stats.Latencies = append(stats.Latencies, r.Duration)
			if r.Duration < stats.MinTime {
				stats.MinTime = r.Duration
			}
			if r.Duration > stats.MaxTime {
				stats.MaxTime = r.Duration
			}
		} else {
			stats.Failed++
		}
	}

	return stats
}

// makeRequest 發送單次請求
func makeRequest(client *http.Client, url string) Result {
	start := time.Now()

	resp, err := client.Get(url)
	if err != nil {
		return Result{Success: false, Duration: time.Since(start)}
	}
	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body) // 消耗 body

	return Result{
		Success:  resp.StatusCode == http.StatusOK,
		Duration: time.Since(start),
		Status:   resp.StatusCode,
	}
}

// printStats 輸出單個端點統計
func printStats(name string, s *Stats) {
	fmt.Printf("Testing: %s\n", name)

	qps := float64(s.Success) / s.TotalTime.Seconds()

	fmt.Printf("  Total:    %d requests in %.2fs\n", s.Total, s.TotalTime.Seconds())
	fmt.Printf("  Success:  %d\n", s.Success)
	fmt.Printf("  Failed:   %d\n", s.Failed)
	fmt.Printf("  QPS:      %.2f req/s\n", qps)

	if len(s.Latencies) > 0 {
		sort.Slice(s.Latencies, func(i, j int) bool {
			return s.Latencies[i] < s.Latencies[j]
		})

		var sum time.Duration
		for _, l := range s.Latencies {
			sum += l
		}
		avg := sum / time.Duration(len(s.Latencies))
		p50 := s.Latencies[len(s.Latencies)/2]
		p95 := s.Latencies[int(float64(len(s.Latencies))*0.95)]
		p99 := s.Latencies[int(float64(len(s.Latencies))*0.99)]

		fmt.Printf("  Avg:      %v\n", avg)
		fmt.Printf("  Min:      %v\n", s.MinTime)
		fmt.Printf("  Max:      %v\n", s.MaxTime)
		fmt.Printf("  P50:      %v\n", p50)
		fmt.Printf("  P95:      %v\n", p95)
		fmt.Printf("  P99:      %v\n", p99)
	}
}
