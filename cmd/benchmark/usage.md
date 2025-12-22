# 默认测试
.\bin\benchmark.exe

# 自定义并发和请求数
.\bin\benchmark.exe -c 100 -n 5000

# 自定义目标地址
.\bin\benchmark.exe -url http://your-server:8080 -c 50 -n 1000

# 示例测试
````
PS C:\Users\2\Desktop\MagTrade> .\bin\benchmark.exe -c 50 -n 2000
==========================================
  MagTrade Benchmark Tool (Go)
==========================================
Target: http://localhost:8080
Concurrency: 50
Total Requests: 2000

Testing: Health Check
  Total:    2000 requests in 0.31s
  Success:  2000
  Failed:   0
  QPS:      6369.43 req/s
  Avg:      7.246979ms
  Min:      0s
  Max:      44.0031ms
  P50:      7.0003ms
  P95:      21.9975ms
  P99:      29.9988ms

Testing: Products
  Total:    2000 requests in 0.84s
  Success:  2000
  Failed:   0
  QPS:      2393.54 req/s
  Avg:      19.758791ms
  Min:      0s
  Max:      565.582ms
  P50:      7.0106ms
  P95:      15.0046ms
  P99:      550.5824ms

Testing: Flash Sales
  Total:    2000 requests in 0.68s
  Success:  2000
  Failed:   0
  QPS:      2954.47 req/s
  Avg:      15.320991ms
  Min:      0s
  Max:      459.9386ms
  P50:      7.0017ms
  P95:      24.0067ms
  P99:      423.9377ms

Testing: Stock Check
  Failed:   0
  QPS:      5076.15 req/s
  Avg:      7.540979ms
  Min:      0s
  Max:      236.999ms
  P50:      4.9981ms
  P95:      11.0044ms
  P99:      143.0001ms

==========================================
  Summary
==========================================

| Endpoint        |      QPS |  Avg(ms) |  P95(ms) |  P99(ms) | Failed |
|-----------------|----------|----------|----------|----------|--------|
| Health Check    |  6369.43 |     7.25 |       21 |       29 |      0 |
| Products        |  2393.54 |    19.76 |       15 |      550 |      0 |
| Flash Sales     |  2954.47 |    15.32 |       24 |      423 |      0 |
| Stock Check     |  5076.15 |     7.54 |       11 |      143 |      0 |
````