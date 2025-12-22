package ai

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Mag1cFall/magtrade/internal/model"
	"go.uber.org/zap"
)

type AnomalyDetector struct {
	log *zap.Logger

	requestCounts     map[string]*requestCounter
	requestCountsLock sync.RWMutex

	userDevices     map[int64]map[string]time.Time
	userDevicesLock sync.RWMutex

	config *AnomalyDetectionConfig
}

type AnomalyDetectionConfig struct {
	IPRequestLimit       int
	IPRequestWindow      time.Duration
	MultiDeviceThreshold int
	MinHumanReaction     time.Duration
	BatchPatternWindow   time.Duration
}

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
			IPRequestLimit:       100,
			IPRequestWindow:      time.Minute,
			MultiDeviceThreshold: 3,
			MinHumanReaction:     200 * time.Millisecond,
			BatchPatternWindow:   time.Minute,
		},
	}
}

type RequestEvent struct {
	UserID       int64
	IP           string
	UserAgent    string
	FlashSaleID  int64
	ResponseTime time.Duration
	Timestamp    time.Time
}

func (d *AnomalyDetector) Detect(ctx context.Context, event *RequestEvent) *model.AnomalyAlert {
	if alert := d.checkIPRateLimit(event); alert != nil {
		return alert
	}

	if alert := d.checkReactionTime(event); alert != nil {
		return alert
	}

	if alert := d.checkMultiDevice(event); alert != nil {
		return alert
	}

	return nil
}

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

func (d *AnomalyDetector) checkMultiDevice(event *RequestEvent) *model.AnomalyAlert {
	d.userDevicesLock.Lock()
	defer d.userDevicesLock.Unlock()

	devices, exists := d.userDevices[event.UserID]
	if !exists {
		devices = make(map[string]time.Time)
		d.userDevices[event.UserID] = devices
	}

	devices[event.UserAgent] = event.Timestamp

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
