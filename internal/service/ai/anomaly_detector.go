// 異常檢測服務
//
// 本檔案實現秒殺場景的反作弊檢測
// 檢測類型：IP 頻率限制、反應時間異常、多設備異常
// 使用記憶體儲存檢測狀態，定期清理過期資料
package ai

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Mag1cFall/magtrade/internal/model"
	"go.uber.org/zap"
)

// AnomalyDetector 異常檢測器
type AnomalyDetector struct {
	log *zap.Logger

	requestCounts     map[string]*requestCounter // IP:活動ID → 請求計數
	requestCountsLock sync.RWMutex

	userDevices     map[int64]map[string]time.Time // 使用者ID → UserAgent → 最後活動時間
	userDevicesLock sync.RWMutex

	config *AnomalyDetectionConfig
}

// AnomalyDetectionConfig 檢測配置
type AnomalyDetectionConfig struct {
	IPRequestLimit       int           // IP 請求上限
	IPRequestWindow      time.Duration // IP 請求統計視窗
	MultiDeviceThreshold int           // 多設備閾值
	MinHumanReaction     time.Duration // 最小人類反應時間
	BatchPatternWindow   time.Duration // 批次模式檢測視窗
}

// requestCounter 請求計數器
type requestCounter struct {
	count     int
	firstSeen time.Time
}

func NewAnomalyDetector(log *zap.Logger) *AnomalyDetector {
	return &AnomalyDetector{
		log:           log,
		requestCounts: make(map[string]*requestCounter),
		userDevices:   make(map[int64]map[string]time.Time),
		config: &AnomalyDetectionConfig{
			IPRequestLimit:       100, // 每分鐘 100 次
			IPRequestWindow:      time.Minute,
			MultiDeviceThreshold: 3,                      // 3 個設備以上異常
			MinHumanReaction:     200 * time.Millisecond, // 人類最低反應時間
			BatchPatternWindow:   time.Minute,
		},
	}
}

// RequestEvent 請求事件
type RequestEvent struct {
	UserID       int64
	IP           string
	UserAgent    string
	FlashSaleID  int64
	ResponseTime time.Duration // 客戶端反應時間
	Timestamp    time.Time
}

// Detect 執行異常檢測（返回警報或 nil）
func (d *AnomalyDetector) Detect(ctx context.Context, event *RequestEvent) *model.AnomalyAlert {
	// 檢查 IP 頻率
	if alert := d.checkIPRateLimit(event); alert != nil {
		return alert
	}

	// 檢查反應時間
	if alert := d.checkReactionTime(event); alert != nil {
		return alert
	}

	// 檢查多設備異常
	if alert := d.checkMultiDevice(event); alert != nil {
		return alert
	}

	return nil
}

// checkIPRateLimit 檢查 IP 請求頻率
func (d *AnomalyDetector) checkIPRateLimit(event *RequestEvent) *model.AnomalyAlert {
	d.requestCountsLock.Lock()
	defer d.requestCountsLock.Unlock()

	key := fmt.Sprintf("%s:%d", event.IP, event.FlashSaleID)

	counter, exists := d.requestCounts[key]
	if !exists {
		d.requestCounts[key] = &requestCounter{
			count:     1,
			firstSeen: event.Timestamp,
		}
		return nil
	}

	// 超出時間視窗，重置計數
	if event.Timestamp.Sub(counter.firstSeen) > d.config.IPRequestWindow {
		d.requestCounts[key] = &requestCounter{
			count:     1,
			firstSeen: event.Timestamp,
		}
		return nil
	}

	counter.count++

	if counter.count > d.config.IPRequestLimit {
		d.log.Warn("IP rate limit exceeded",
			zap.String("ip", event.IP),
			zap.Int("count", counter.count),
		)

		return &model.AnomalyAlert{
			AlertType: "ip_rate_limit_exceeded",
			Severity:  "high",
			Details: map[string]interface{}{
				"ip":            event.IP,
				"request_count": counter.count,
				"window":        d.config.IPRequestWindow.String(),
				"flash_sale_id": event.FlashSaleID,
			},
			AutoAction: "rate_limit_applied",
			Timestamp:  time.Now(),
		}
	}

	return nil
}

// checkReactionTime 檢查反應時間（過快可能是機器人）
func (d *AnomalyDetector) checkReactionTime(event *RequestEvent) *model.AnomalyAlert {
	if event.ResponseTime > 0 && event.ResponseTime < d.config.MinHumanReaction {
		d.log.Warn("suspicious reaction time",
			zap.Int64("user_id", event.UserID),
			zap.Duration("response_time", event.ResponseTime),
		)

		return &model.AnomalyAlert{
			AlertType: "suspected_bot",
			Severity:  "medium",
			Details: map[string]interface{}{
				"user_id":           event.UserID,
				"response_time_ms":  event.ResponseTime.Milliseconds(),
				"min_human_time_ms": d.config.MinHumanReaction.Milliseconds(),
				"behavior":          "响应时间低于人类正常反应速度",
			},
			AutoAction: "verification_required",
			Timestamp:  time.Now(),
		}
	}

	return nil
}

// checkMultiDevice 檢查多設備異常
func (d *AnomalyDetector) checkMultiDevice(event *RequestEvent) *model.AnomalyAlert {
	d.userDevicesLock.Lock()
	defer d.userDevicesLock.Unlock()

	devices, exists := d.userDevices[event.UserID]
	if !exists {
		devices = make(map[string]time.Time)
		d.userDevices[event.UserID] = devices
	}

	devices[event.UserAgent] = event.Timestamp

	// 清理過期設備記錄
	for ua, lastSeen := range devices {
		if event.Timestamp.Sub(lastSeen) > 10*time.Minute {
			delete(devices, ua)
		}
	}

	if len(devices) > d.config.MultiDeviceThreshold {
		d.log.Warn("multiple devices detected",
			zap.Int64("user_id", event.UserID),
			zap.Int("device_count", len(devices)),
		)

		return &model.AnomalyAlert{
			AlertType: "multi_device_anomaly",
			Severity:  "medium",
			Details: map[string]interface{}{
				"user_id":      event.UserID,
				"device_count": len(devices),
				"threshold":    d.config.MultiDeviceThreshold,
			},
			AutoAction: "monitoring",
			Timestamp:  time.Now(),
		}
	}

	return nil
}

// RecordRequest 記錄請求（不執行檢測）
func (d *AnomalyDetector) RecordRequest(event *RequestEvent) {
	d.requestCountsLock.Lock()
	key := fmt.Sprintf("%s:%d", event.IP, event.FlashSaleID)
	counter, exists := d.requestCounts[key]
	if !exists {
		d.requestCounts[key] = &requestCounter{
			count:     1,
			firstSeen: event.Timestamp,
		}
	} else {
		counter.count++
	}
	d.requestCountsLock.Unlock()

	d.userDevicesLock.Lock()
	devices, exists := d.userDevices[event.UserID]
	if !exists {
		devices = make(map[string]time.Time)
		d.userDevices[event.UserID] = devices
	}
	devices[event.UserAgent] = event.Timestamp
	d.userDevicesLock.Unlock()
}

// Cleanup 清理過期資料（應定期呼叫）
func (d *AnomalyDetector) Cleanup() {
	d.requestCountsLock.Lock()
	now := time.Now()
	for key, counter := range d.requestCounts {
		if now.Sub(counter.firstSeen) > d.config.IPRequestWindow*2 {
			delete(d.requestCounts, key)
		}
	}
	d.requestCountsLock.Unlock()

	d.userDevicesLock.Lock()
	for userID, devices := range d.userDevices {
		for ua, lastSeen := range devices {
			if now.Sub(lastSeen) > 30*time.Minute {
				delete(devices, ua)
			}
		}
		if len(devices) == 0 {
			delete(d.userDevices, userID)
		}
	}
	d.userDevicesLock.Unlock()
}
