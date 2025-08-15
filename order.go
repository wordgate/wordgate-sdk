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

// ListOrdersQuery represents query parameters for listing orders
type ListOrdersQuery struct {
	// Page is the page number (minimum 1)
	Page int `form:"page" binding:"min=1"`
	// Limit is the number of items per page (1-100)
	Limit int `form:"limit" binding:"min=1,max=100"`
	// Status filters orders by payment status (paid/unpaid)
	Status string `form:"status"`
	// UserUID filters orders by user UID
	UserUID string `form:"user_uid"`
	// Email filters orders by user email
	Email string `form:"email"`
	// StartAt filters orders created after this date (YYYY-MM-DD)
	StartAt string `form:"start_at"`
	// EndAt filters orders created before this date (YYYY-MM-DD)
	EndAt string `form:"end_at"`
	// OrderNo filters orders by order number (partial match)
	OrderNo string `form:"order_no"`
	// SortBy specifies the field to sort by (created_at/amount)
	SortBy string `form:"sort_by"`
	// SortDesc specifies whether to sort in descending order
	SortDesc bool `form:"sort_desc"`
}

// PaymentIntentInfo represents payment intent information
type PaymentIntentInfo struct {
	// ID is the payment intent database ID
	ID uint64 `json:"id"`
	// IntentID is the payment platform generated intent ID
	IntentID string `json:"intent_id"`
	// Provider is the payment provider name
	Provider string `json:"provider"`
	// Status is the payment status
	Status string `json:"status"`
	// Amount is the payment amount in cents
	Amount int64 `json:"amount"`
	// Currency is the currency code
	Currency string `json:"currency"`
	// CreatedAt is the creation timestamp
	CreatedAt string `json:"created_at"`
	// PaidAt is the payment completion timestamp (nil if not paid)
	PaidAt *string `json:"paid_at"`
}

// OrderItemInfo represents order item information
type OrderItemInfo struct {
	// ItemID is the product ID
	ItemID uint64 `json:"item_id"`
	// ItemName is the product name
	ItemName string `json:"item_name"`
	// Quantity is the number of items
	Quantity int `json:"quantity"`
	// UnitPrice is the unit price in cents
	UnitPrice int64 `json:"unit_price"`
	// Subtotal is the total price for this item (UnitPrice * Quantity)
	Subtotal int64 `json:"subtotal"`
	// RequireAddress indicates if this item requires shipping address
	RequireAddress bool `json:"require_address"`
}

// AddressInfo represents address information
type AddressInfo struct {
	// Name is the recipient name
	Name string `json:"name"`
	// Phone is the recipient phone number
	Phone string `json:"phone"`
	// Province is the province/state
	Province string `json:"province"`
	// City is the city
	City string `json:"city"`
	// District is the district/county
	District string `json:"district"`
	// Street is the street address
	Street string `json:"street"`
	// PostalCode is the postal/zip code
	PostalCode string `json:"postal_code"`
	// Label is the address label (e.g., "Home", "Office")
	Label string `json:"label"`
}

// OrderDetailResponse represents detailed order information
type OrderDetailResponse struct {
	// ID is the order database ID
	ID uint64 `json:"id"`
	// OrderNo is the unique order number
	OrderNo string `json:"order_no"`
	// UserID is the user database ID
	UserID uint64 `json:"user_id"`
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
	// Address is the shipping address information (nil if not required)
	Address *AddressInfo `json:"address"`
	// PayURL is the direct payment URL
	PayURL string `json:"pay_url"`
	// PaymentIntents is the list of payment intents
	PaymentIntents []PaymentIntentInfo `json:"payment_intents"`
	// User is the user information
	User interface{} `json:"user,omitempty"`
}

// OrderListItem represents an order item in the list
type OrderListItem struct {
	// ID is the order database ID
	ID uint64 `json:"id"`
	// OrderNo is the unique order number
	OrderNo string `json:"order_no"`
	// UserID is the user database ID
	UserID uint64 `json:"user_id"`
	// Amount is the total amount in cents
	Amount int64 `json:"amount"`
	// Currency is the currency code
	Currency string `json:"currency"`
	// IsPaid indicates whether the order is paid
	IsPaid bool `json:"is_paid"`
	// PaidAt is the payment timestamp (nil if not paid)
	PaidAt *time.Time `json:"paid_at"`
	// CreatedAt is the creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// PaymentIntents is the list of payment intents
	PaymentIntents []PaymentIntentInfo `json:"payment_intents"`
	// ItemsSummary is a summary of order items
	ItemsSummary string `json:"items_summary"`
	// ItemsCount is the number of items in the order
	ItemsCount int `json:"items_count"`
	// UserInfo is the user information
	UserInfo interface{} `json:"user_info,omitempty"`
	// RequireAddress indicates if the order requires shipping address
	RequireAddress bool `json:"require_address"`
}

// ListResult represents a paginated list result
type ListResult struct {
	// Data is the list of items
	Data interface{} `json:"data"`
	// Pagination contains pagination information
	Pagination *Pagination `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	// Page is the current page number
	Page int `json:"page"`
	// Limit is the number of items per page
	Limit int `json:"limit"`
	// Total is the total number of items
	Total int64 `json:"total"`
	// TotalPages is the total number of pages
	TotalPages int `json:"total_pages"`
}

// ManualPaymentRequest represents a request to manually mark order as paid
type ManualPaymentRequest struct {
	// OrderNo is the order number
	OrderNo string `json:"order_no"`
	// PaymentNote is the payment note (required)
	PaymentNote string `json:"payment_note"`
	// Amount is the payment amount in cents (optional, defaults to order amount)
	Amount *int64 `json:"amount,omitempty"`
}

// GetAppOrder retrieves detailed order information by order number
//
// orderNo: The order number to retrieve
// Returns the detailed order information and any error
func (c *Client) GetAppOrder(orderNo string) (*OrderDetailResponse, error) {
	var result OrderDetailResponse
	err := c.requestJSON("GET", "/app/orders/"+orderNo, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get app order: %w", err)
	}
	return &result, nil
}

// ListAppOrders retrieves a paginated list of orders with optional filtering
//
// query: The query parameters for filtering and pagination
// Returns the order list result and any error
func (c *Client) ListAppOrders(query *ListOrdersQuery) (*ListResult, error) {
	var result ListResult
	err := c.requestJSON("GET", "/app/orders", query, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to list app orders: %w", err)
	}
	return &result, nil
}

// MarkOrderAsPaid manually marks an order as paid
//
// request: The manual payment request containing order number and payment note
// Returns any error
func (c *Client) MarkOrderAsPaid(request *ManualPaymentRequest) error {
	var result interface{}
	err := c.requestJSON("POST", "/app/orders/mark_as_paid", request, &result)
	if err != nil {
		return fmt.Errorf("failed to mark order as paid: %w", err)
	}
	return nil
}
