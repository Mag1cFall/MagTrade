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

type Result struct {
	Success  bool
	Duration time.Duration
	Status   int
}

type Stats struct {
	Total     int64
	Success   int64
	Failed    int64
	TotalTime time.Duration
	MinTime   time.Duration
	MaxTime   time.Duration
	Latencies []time.Duration
}

func main() {
	baseURL := flag.String("url", "http://localhost:8080", "Base URL")
	concurrency := flag.Int("c", 100, "Concurrency level")
	requests := flag.Int("n", 1000, "Total requests")
	flag.Parse()

	fmt.Println("==========================================")
	fmt.Println("  MagTrade Benchmark Tool (Go)")
	fmt.Println("==========================================")
	fmt.Printf("Target: %s\n", *baseURL)
	fmt.Printf("Concurrency: %d\n", *concurrency)
	fmt.Printf("Total Requests: %d\n", *requests)
	fmt.Println()

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

func makeRequest(client *http.Client, url string) Result {
	start := time.Now()

	resp, err := client.Get(url)
	if err != nil {
		return Result{Success: false, Duration: time.Since(start)}
	}
	defer resp.Body.Close()

	io.Copy(io.Discard, resp.Body)

	return Result{
		Success:  resp.StatusCode == http.StatusOK,
		Duration: time.Since(start),
		Status:   resp.StatusCode,
	}
}

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
