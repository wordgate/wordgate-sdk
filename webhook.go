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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

// WebhookSignature webhook签名相关结构体
type WebhookSignature struct {
	Timestamp int64  `json:"timestamp"` // 时间戳
	Signature string `json:"signature"` // HMAC-SHA256签名
}

// GenerateSignature 生成webhook签名
// timestamp: Unix时间戳(秒)
// body: webhook请求体原文
// secret: 签名密钥
func GenerateSignature(timestamp int64, body []byte, secret string) string {
	// 构造待签名消息: <timestamp>.<body>
	message := strconv.FormatInt(timestamp, 10) + "." + string(body)
	
	// 使用HMAC-SHA256生成签名
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	signature := hex.EncodeToString(h.Sum(nil))
	
	return signature
}

// GenerateSignatureHeader 生成X-Webhook-Signature header值
// timestamp: Unix时间戳(秒)
// body: webhook请求体原文
// secret: 签名密钥
func GenerateSignatureHeader(timestamp int64, body []byte, secret string) string {
	signature := GenerateSignature(timestamp, body, secret)
	return fmt.Sprintf("t=%d,sha256=%s", timestamp, signature)
}

// ParseSignatureHeader 解析X-Webhook-Signature header
// headerValue: X-Webhook-Signature header的值，格式为 "t=<timestamp>,sha256=<signature>"
func ParseSignatureHeader(headerValue string) (*WebhookSignature, error) {
	parts := strings.Split(headerValue, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid signature header format")
	}
	
	var timestamp int64
	var signature string
	
	for _, part := range parts {
		if strings.HasPrefix(part, "t=") {
			var err error
			timestamp, err = strconv.ParseInt(strings.TrimPrefix(part, "t="), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid timestamp: %w", err)
			}
		} else if strings.HasPrefix(part, "sha256=") {
			signature = strings.TrimPrefix(part, "sha256=")
		}
	}
	
	if timestamp == 0 || signature == "" {
		return nil, fmt.Errorf("missing timestamp or signature")
	}
	
	return &WebhookSignature{
		Timestamp: timestamp,
		Signature: signature,
	}, nil
}

// VerifySignature 验证webhook签名
// headerValue: X-Webhook-Signature header的值
// body: webhook请求体原文
// secret: 签名密钥
// maxTimeDiff: 最大时间差(秒)，用于防重放攻击，建议300秒
func VerifySignature(headerValue string, body []byte, secret string, maxTimeDiff int64) error {
	// 解析签名header
	webhookSig, err := ParseSignatureHeader(headerValue)
	if err != nil {
		return fmt.Errorf("parse signature header failed: %w", err)
	}
	
	// 检查时间戳，防重放攻击
	now := time.Now().Unix()
	if now-webhookSig.Timestamp > maxTimeDiff {
		return fmt.Errorf("timestamp too old: %d seconds ago", now-webhookSig.Timestamp)
	}
	
	// 生成预期签名
	expectedSignature := GenerateSignature(webhookSig.Timestamp, body, secret)
	
	// 使用常量时间比较，防时序攻击
	if !hmac.Equal([]byte(webhookSig.Signature), []byte(expectedSignature)) {
		return fmt.Errorf("signature verification failed")
	}
	
	return nil
}
