package wordgate

import (
	"fmt"
	"time"
)

// OrderItem represents an item in an order
type OrderItem struct {
	// ItemCode is the product/item code
	ItemCode string `json:"item_code"`
	// Quantity is the number of items
	Quantity int `json:"quantity"`
}

// OrderCustomer represents customer information for an order
type OrderCustomer struct {
	// UserUID is the user's unique identifier
	UserUID string `json:"user_uid"`
}


// CreateMembershipOrderRequest represents a request to create a membership order
type CreateMembershipOrderRequest struct {
	// TierID is the membership tier ID
	TierID uint64 `json:"tier_id"`
	// PeriodType is the membership period type (month, quarter, year, etc.)
	PeriodType string `json:"period_type"`
	// CouponCode is an optional coupon code
	CouponCode string `json:"coupon_code,omitempty"`
	// ClientIP is the client's IP address (optional)
	ClientIP string `json:"client_ip,omitempty"`
	// AddressID is the shipping address ID (optional for membership orders)
	AddressID uint64 `json:"address_id,omitempty"`
}

// OrderResponse represents the response when creating an order
type OrderResponse struct {
	// OrderNo is the unique order number
	OrderNo string `json:"order_no"`
	// Amount is the total amount in cents
	Amount int64 `json:"amount"`
	// Currency is the currency code (e.g., "CNY", "USD")
	Currency string `json:"currency"`
	// IsPaid indicates whether the order is paid
	IsPaid bool `json:"is_paid"`
	// PaidAt is the payment timestamp (nil if not paid)
	PaidAt *time.Time `json:"paid_at"`
	// PayURL is the direct payment URL
	PayURL string `json:"pay_url"`
	// RedirectURL is the payment completion redirect URL
	RedirectURL string `json:"redirect_url"`
}

// OrderSummaryResponse represents the response when creating an app order
type OrderSummaryResponse struct {
	// OrderNo is the unique order number
	OrderNo string `json:"order_no"`
	// Amount is the total amount in cents
	Amount int64 `json:"amount"`
	// Currency is the currency code (e.g., "CNY", "USD")
	Currency string `json:"currency"`
	// IsPaid indicates whether the order is paid
	IsPaid bool `json:"is_paid"`
	// PaidAt is the payment timestamp (nil if not paid)
	PaidAt *time.Time `json:"paid_at"`
	// PayURL is the direct payment URL
	PayURL string `json:"pay_url"`
	// RedirectURL is the payment completion redirect URL (optional)
	RedirectURL string `json:"redirect_url"`
}

// OrderItemInfo represents detailed information about an order item
type OrderItemInfo struct {
	// ItemID is the internal item ID
	ItemID uint `json:"item_id"`
	// ItemName is the item name
	ItemName string `json:"item_name"`
	// Quantity is the number of items
	Quantity int `json:"quantity"`
	// UnitPrice is the unit price in cents
	UnitPrice int64 `json:"unit_price"`
	// Subtotal is the subtotal amount in cents
	Subtotal int64 `json:"subtotal"`
}

// OrderDetailResponse represents detailed order information
type OrderDetailResponse struct {
	// ID is the internal order ID
	ID uint `json:"id"`
	// OrderNo is the unique order number
	OrderNo string `json:"order_no"`
	// UserID is the user ID
	UserID uint `json:"user_id"`
	// Amount is the total amount in cents
	Amount int64 `json:"amount"`
	// Currency is the currency code
	Currency string `json:"currency"`
	// IsPaid indicates whether the order is paid
	IsPaid bool `json:"is_paid"`
	// CreatedAt is the creation timestamp
	CreatedAt string `json:"created_at"`
	// PaidAt is the payment timestamp (nil if not paid)
	PaidAt *string `json:"paid_at"`
	// CouponCode is the applied coupon code
	CouponCode string `json:"coupon_code"`
	// DiscountAmount is the discount amount in cents
	DiscountAmount int64 `json:"discount_amount"`
	// Items is the list of order items
	Items []OrderItemInfo `json:"items"`
	// PayURL is the payment URL
	PayURL string `json:"pay_url"`
	// RedirectURL is the payment completion redirect URL
	RedirectURL string `json:"redirect_url"`
}


// CreateMembershipOrder creates a new membership order
//
// request: The membership order creation request containing membership tier and period info
// Returns the created order information and any error
func (c *Client) CreateMembershipOrder(request *CreateMembershipOrderRequest) (*OrderResponse, error) {
	var result OrderResponse
	err := c.requestJSON("POST", "/api/membership-orders/create", request, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create membership order: %w", err)
	}
	return &result, nil
}

// GetOrder retrieves order details by order number
//
// orderNo: The order number to retrieve
// Returns the order details and any error
func (c *Client) GetOrder(orderNo string) (*OrderDetailResponse, error) {
	var result OrderDetailResponse
	path := fmt.Sprintf("/app/orders/%s", orderNo)
	err := c.requestJSON("GET", path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return &result, nil
}

// CreateAppProductOrderRequest represents a request to create a product order via app admin API
type CreateAppProductOrderRequest struct {
	// Items is the list of product items
	Items []struct {
		ItemCode string `json:"item_code"`
		Quantity int    `json:"quantity"`
	} `json:"items"`
	// CouponCode is an optional coupon code
	CouponCode string `json:"coupon_code,omitempty"`
	// ClientIP is the client's IP address (optional)
	ClientIP string `json:"client_ip,omitempty"`
	// AddressID is the shipping address ID
	AddressID uint64 `json:"address_id"`
	// UserUID is the user's unique identifier
	UserUID string `json:"user_uid"`
	// RedirectURL is the payment completion redirect URL (optional)
	RedirectURL string `json:"redirect_url,omitempty"`
}

// CreateAppMembershipOrderRequest represents a request to create a membership order via app admin API
type CreateAppMembershipOrderRequest struct {
	// TierID is the membership tier ID
	TierID uint64 `json:"tier_id"`
	// PeriodType is the membership period type (month, quarter, year, etc.)
	PeriodType string `json:"period_type"`
	// CouponCode is an optional coupon code
	CouponCode string `json:"coupon_code,omitempty"`
	// ClientIP is the client's IP address (optional)
	ClientIP string `json:"client_ip,omitempty"`
	// AddressID is the shipping address ID (optional for membership orders)
	AddressID uint64 `json:"address_id,omitempty"`
	// UserUID is the user's unique identifier
	UserUID string `json:"user_uid"`
	// RedirectURL is the payment completion redirect URL (optional)
	RedirectURL string `json:"redirect_url,omitempty"`
}

// CreateAppProductOrder creates a new product order using admin API
//
// request: The product order creation request containing items and customer info
// Returns the created order information and any error
func (c *Client) CreateAppProductOrder(request *CreateAppProductOrderRequest) (*OrderSummaryResponse, error) {
	var result OrderSummaryResponse
	err := c.requestJSON("POST", "/app/product-orders/create", request, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create app product order: %w", err)
	}
	return &result, nil
}

// CreateAppMembershipOrder creates a new membership order using admin API
//
// request: The membership order creation request containing tier and period info
// Returns the created order information and any error
func (c *Client) CreateAppMembershipOrder(request *CreateAppMembershipOrderRequest) (*OrderSummaryResponse, error) {
	var result OrderSummaryResponse
	err := c.requestJSON("POST", "/app/membership-orders/create", request, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create app membership order: %w", err)
	}
	return &result, nil
}