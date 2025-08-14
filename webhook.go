/*
Webhook handling utilities for WordGate API.

This file provides structures and utilities for handling webhook events from WordGate,
including order payments, cancellations, and subscription events.

Usage example:

	// Parse webhook events in your webhook handler
	func handleWebhook(w http.ResponseWriter, r *http.Request) {
		var webhookEvent wordgate.WebhookEventData
		if err := json.NewDecoder(r.Body).Decode(&webhookEvent); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		switch webhookEvent.EventType {
		case WebhookEventOrderPaid:
			var orderData WebhookOrderPaidData
			if err := webhookEvent.Parse(&orderData); err != nil {
				http.Error(w, "Failed to parse order data", http.StatusBadRequest)
				return
			}

			// Handle order paid event
			log.Printf("Order %s paid: %d %s", orderData.WordgateOrderNo, orderData.Amount, orderData.Currency)

		case WebhookEventOrderCancelled:
			var cancelData WebhookOrderCancelledData
			if err := webhookEvent.Parse(&cancelData); err != nil {
				http.Error(w, "Failed to parse cancel data", http.StatusBadRequest)
				return
			}

			// Handle order cancelled event
			log.Printf("Order %s cancelled: %s", cancelData.WordgateOrderNo, cancelData.Reason)

		}

		w.WriteHeader(http.StatusOK)
	}
*/
package wordgate

import (
	"encoding/json"
	"fmt"
	"time"
)

// WebhookEventData webhook事件数据格式
type WebhookEventData struct {
	EventType WebhookEventType `json:"event_type"` // 事件类型，如: "order.paid", "order.cancelled" 等
	AppID     uint64           `json:"app_id"`     // 应用ID
	Data      any              `json:"data"`       // 事件数据，具体内容取决于事件类型
	Timestamp int64            `json:"timestamp"`  // 事件时间戳
}

// Parse 解析事件数据为指定类型
func (w *WebhookEventData) Parse(target any) error {
	jsonBytes, err := json.Marshal(w.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	return json.Unmarshal(jsonBytes, target)
}

// WebhookOrderPaidData 订单支付成功事件的数据结构
type WebhookOrderPaidData struct {
	WordgateOrderNo string     `json:"wordgate_order_no"` // 订单号
	Amount          int64      `json:"amount"`            // 订单金额
	Currency        string     `json:"currency"`          // 货币类型
	IsPaid          bool       `json:"is_paid"`           // 是否已支付
	PaidAt          *time.Time `json:"paid_at"`           // 支付时间
	AppID           uint64     `json:"app_id"`            // 应用ID
}

// WebhookOrderCancelledData 订单取消事件的数据结构
type WebhookOrderCancelledData struct {
	WordgateOrderNo string     `json:"wordgate_order_no"` // 订单号
	Amount          int64      `json:"amount"`            // 订单金额
	Currency        string     `json:"currency"`          // 货币类型
	CancelledAt     *time.Time `json:"cancelled_at"`      // 取消时间
	AppID           uint64     `json:"app_id"`            // 应用ID
	Reason          string     `json:"reason"`            // 取消原因
}

// WebhookMembershipActivatedData 会员变动事件的数据结构
type WebhookMembershipActivatedData struct {
	UserID    uint64 `json:"user_id"`    // 用户ID
	TierCode  string `json:"tier_code"`  // 会员等级代码
	ExpiresAt string `json:"expires_at"` // 到期时间 (ISO格式)
	AppID     uint64 `json:"app_id"`     // 应用ID
}

// WebhookEventType 定义支持的webhook事件类型常量
type WebhookEventType string

const (
	WebhookEventOrderPaid           WebhookEventType = "order.paid"           // 订单支付成功
	WebhookEventOrderCancelled      WebhookEventType = "order.cancelled"      // 订单取消
	WebhookEventMembershipActivated WebhookEventType = "membership.activated" // 会员变动
)
