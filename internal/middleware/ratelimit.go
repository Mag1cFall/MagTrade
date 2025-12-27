// 限流中間件
//
// 本檔案使用令牌桶（Token Bucket）算法實現 IP 級別限流
// IPRateLimiter：全域 IP 限流
// FlashSaleRateLimiter：秒殺活動專用限流（按 IP+活動ID）
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/Mag1cFall/magtrade/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter IP 限流器
// 每個 IP 維護獨立的令牌桶
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter // IP → Limiter 映射
	mu       sync.RWMutex             // 讀寫鎖保護併發訪問
	r        rate.Limit               // 每秒令牌生成速率
	b        int                      // 令牌桶容量（允許突發）
}

// NewIPRateLimiter 建立 IP 限流器
func NewIPRateLimiter(cfg *config.RateLimitConfig) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        rate.Limit(cfg.RequestsPerSecond),
		b:        cfg.Burst,
	}
}

// GetLimiter 取得指定 IP 的限流器（惰性建立）
func (l *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	l.mu.RLock()
	limiter, exists := l.limiters[ip]
	l.mu.RUnlock()

	if exists {
		return limiter
	}

	// Double-check 模式防止重複建立
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists = l.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(l.r, l.b)
		l.limiters[ip] = limiter
	}

	return limiter
}

// Cleanup 清理所有限流器（定期調用防止記憶體洩漏）
func (l *IPRateLimiter) Cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.limiters = make(map[string]*rate.Limiter)
}

// RateLimit 全域 IP 限流中間件
func RateLimit(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)

		if !l.Allow() { // 嘗試取令牌，無令牌則拒絕
			response.TooManyRequests(c, "rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}

// FlashSaleRateLimiter 秒殺專用限流器
// 按「IP + 活動ID」組合限流，每秒最多 10 次請求
type FlashSaleRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
}

func NewFlashSaleRateLimiter() *FlashSaleRateLimiter {
	return &FlashSaleRateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

// GetLimiter 取得「IP:活動ID」組合的限流器
func (l *FlashSaleRateLimiter) GetLimiter(key string) *rate.Limiter {
	l.mu.RLock()
	limiter, exists := l.limiters[key]
	l.mu.RUnlock()

	if exists {
		return limiter
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists = l.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Second), 10) // 每秒 10 次
		l.limiters[key] = limiter
	}

	return limiter
}

// FlashSaleRateLimit 秒殺限流中間件
func FlashSaleRateLimit(limiter *FlashSaleRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		flashSaleID := c.Param("id")
		key := ip + ":" + flashSaleID // 組合 Key

		l := limiter.GetLimiter(key)

		if !l.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后重试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
