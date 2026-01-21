// 訂單模型單元測試
//
// 測試覆蓋：
// - Order.CanPay: 訂單是否可付款（僅待付款狀態可付款）
// - Order.CanCancel: 訂單是否可取消（僅待付款狀態可取消）
// - Order.CanRefund: 訂單是否可退款（僅已付款狀態可退款）
// - OrderStatus.String: 狀態枚舉字串轉換
package model

import (
	"testing"
)

func TestOrder_CanPay(t *testing.T) {
	tests := []struct {
		name  string
		order *Order
		want  bool
	}{
		{
			name:  "pending order can pay",
			order: &Order{Status: OrderStatusPending},
			want:  true,
		},
		{
			name:  "paid order cannot pay again",
			order: &Order{Status: OrderStatusPaid},
			want:  false,
		},
		{
			name:  "cancelled order cannot pay",
			order: &Order{Status: OrderStatusCancelled},
			want:  false,
		},
		{
			name:  "refunded order cannot pay",
			order: &Order{Status: OrderStatusRefunded},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.order.CanPay(); got != tt.want {
				t.Errorf("Order.CanPay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_CanCancel(t *testing.T) {
	tests := []struct {
		name  string
		order *Order
		want  bool
	}{
		{
			name:  "pending order can cancel",
			order: &Order{Status: OrderStatusPending},
			want:  true,
		},
		{
			name:  "paid order cannot cancel directly",
			order: &Order{Status: OrderStatusPaid},
			want:  false,
		},
		{
			name:  "already cancelled order",
			order: &Order{Status: OrderStatusCancelled},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.order.CanCancel(); got != tt.want {
				t.Errorf("Order.CanCancel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_CanRefund(t *testing.T) {
	tests := []struct {
		name  string
		order *Order
		want  bool
	}{
		{
			name:  "paid order can refund",
			order: &Order{Status: OrderStatusPaid},
			want:  true,
		},
		{
			name:  "pending order cannot refund",
			order: &Order{Status: OrderStatusPending},
			want:  false,
		},
		{
			name:  "cancelled order cannot refund",
			order: &Order{Status: OrderStatusCancelled},
			want:  false,
		},
		{
			name:  "already refunded",
			order: &Order{Status: OrderStatusRefunded},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.order.CanRefund(); got != tt.want {
				t.Errorf("Order.CanRefund() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderStatus_String(t *testing.T) {
	tests := []struct {
		status OrderStatus
		want   string
	}{
		{OrderStatusPending, "pending"},
		{OrderStatusPaid, "paid"},
		{OrderStatusCancelled, "cancelled"},
		{OrderStatusRefunded, "refunded"},
		{OrderStatus(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("OrderStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
