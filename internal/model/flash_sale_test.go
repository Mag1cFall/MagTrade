package model

import (
	"testing"
	"time"
)

func TestFlashSale_IsActive(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		flashSale *FlashSale
		want      bool
	}{
		{
			name: "active flash sale within time range",
			flashSale: &FlashSale{
				Status:    FlashSaleStatusActive,
				StartTime: now.Add(-1 * time.Hour),
				EndTime:   now.Add(1 * time.Hour),
			},
			want: true,
		},
		{
			name: "pending status should be inactive",
			flashSale: &FlashSale{
				Status:    FlashSaleStatusPending,
				StartTime: now.Add(-1 * time.Hour),
				EndTime:   now.Add(1 * time.Hour),
			},
			want: false,
		},
		{
			name: "not started yet",
			flashSale: &FlashSale{
				Status:    FlashSaleStatusActive,
				StartTime: now.Add(1 * time.Hour),
				EndTime:   now.Add(2 * time.Hour),
			},
			want: false,
		},
		{
			name: "already ended",
			flashSale: &FlashSale{
				Status:    FlashSaleStatusActive,
				StartTime: now.Add(-2 * time.Hour),
				EndTime:   now.Add(-1 * time.Hour),
			},
			want: false,
		},
		{
			name: "finished status",
			flashSale: &FlashSale{
				Status:    FlashSaleStatusFinished,
				StartTime: now.Add(-2 * time.Hour),
				EndTime:   now.Add(-1 * time.Hour),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flashSale.IsActive(); got != tt.want {
				t.Errorf("FlashSale.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlashSale_IsPending(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		flashSale *FlashSale
		want      bool
	}{
		{
			name: "pending and not started",
			flashSale: &FlashSale{
				Status:    FlashSaleStatusPending,
				StartTime: now.Add(1 * time.Hour),
			},
			want: true,
		},
		{
			name: "pending but already started should be false",
			flashSale: &FlashSale{
				Status:    FlashSaleStatusPending,
				StartTime: now.Add(-1 * time.Hour),
			},
			want: false,
		},
		{
			name: "active status",
			flashSale: &FlashSale{
				Status:    FlashSaleStatusActive,
				StartTime: now.Add(-1 * time.Hour),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flashSale.IsPending(); got != tt.want {
				t.Errorf("FlashSale.IsPending() = %v, want %v", got, tt.want)
			}
		})
	}
}
