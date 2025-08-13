package wordgate

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// MembershipTierStatus represents the status of a membership tier
type MembershipTierStatus string

const (
	// MembershipTierStatusActive indicates the tier is active and available
	MembershipTierStatusActive MembershipTierStatus = "active"
	// MembershipTierStatusInactive indicates the tier is inactive and unavailable
	MembershipTierStatusInactive MembershipTierStatus = "inactive"
)

// MembershipPeriodType represents the membership period type
type MembershipPeriodType string

const (
	// PeriodTypeMonth represents monthly membership
	PeriodTypeMonth MembershipPeriodType = "month"
	// PeriodTypeQuarter represents quarterly membership (3 months)
	PeriodTypeQuarter MembershipPeriodType = "quarter"
	// PeriodTypeHalfYear represents half-yearly membership (6 months)
	PeriodTypeHalfYear MembershipPeriodType = "halfyear"
	// PeriodTypeYear represents yearly membership
	PeriodTypeYear MembershipPeriodType = "year"
	// PeriodTypeTwoYear represents two-year membership
	PeriodTypeTwoYear MembershipPeriodType = "twoyear"
	// PeriodTypeThreeYear represents three-year membership
	PeriodTypeThreeYear MembershipPeriodType = "threeyear"
	// PeriodTypeFiveYear represents five-year membership
	PeriodTypeFiveYear MembershipPeriodType = "fiveyear"
)

// MembershipPrice represents a pricing option for a membership tier
type MembershipPrice struct {
	// ID is the unique identifier of the price
	ID uint64 `json:"id"`
	// AppID is the application ID this price belongs to
	AppID uint64 `json:"app_id"`
	// TierID is the membership tier ID this price belongs to
	TierID uint64 `json:"tier_id"`
	// PeriodType is the membership period type
	PeriodType MembershipPeriodType `json:"period_type"`
	// Price is the discounted price in cents
	Price int64 `json:"price"`
	// OriginalPrice is the original price in cents
	OriginalPrice int64 `json:"original_price"`
	// Months is the number of months for this period
	Months int `json:"months"`
	// Status is the price status
	Status string `json:"status"`
	// CreatedAt is the creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the last update timestamp
	UpdatedAt time.Time `json:"updated_at"`
}

// MembershipTier represents a membership tier in the WordGate system
type MembershipTier struct {
	// ID is the unique identifier of the membership tier
	ID uint64 `json:"id"`
	// AppID is the application ID this tier belongs to
	AppID uint64 `json:"app_id"`
	// Code is the unique tier code
	Code string `json:"code"`
	// Name is the tier name
	Name string `json:"name"`
	// Level is the tier level (higher number = higher tier)
	Level int `json:"level"`
	// Status is the tier status (active/inactive)
	Status MembershipTierStatus `json:"status"`
	// IsDefault indicates whether this is the default tier
	IsDefault bool `json:"is_default"`
	// Version is the version number for optimistic locking
	Version int `json:"version"`
	// Prices is the list of pricing options for this tier
	Prices []MembershipPrice `json:"prices,omitempty"`
	// CreatedAt is the creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the last update timestamp
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt is the deletion timestamp (nil if not deleted)
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// MembershipPriceRequest represents a price request for membership tier operations
type MembershipPriceRequest struct {
	// PeriodType is the membership period type
	PeriodType MembershipPeriodType `json:"period_type"`
	// Price is the discounted price in cents
	Price int64 `json:"price"`
	// OriginalPrice is the original price in cents
	OriginalPrice int64 `json:"original_price"`
}

// CreateMembershipTierRequest represents a request to create a membership tier
type CreateMembershipTierRequest struct {
	// Code is the unique tier code
	Code string `json:"code"`
	// Name is the tier name
	Name string `json:"name"`
	// Level is the tier level (higher number = higher tier)
	Level int `json:"level"`
	// IsDefault indicates whether this is the default tier
	IsDefault bool `json:"is_default"`
	// Prices is the list of pricing options for this tier
	Prices []MembershipPriceRequest `json:"prices"`
}

// UpdateMembershipTierRequest represents a request to update a membership tier
type UpdateMembershipTierRequest struct {
	// Name is the tier name
	Name string `json:"name"`
	// Level is the tier level (higher number = higher tier)
	Level int `json:"level"`
	// IsDefault indicates whether this is the default tier
	IsDefault bool `json:"is_default"`
	// Prices is the list of pricing options for this tier
	Prices []MembershipPriceRequest `json:"prices"`
}

// ListMembershipTiersRequest represents a request to list membership tiers
type ListMembershipTiersRequest struct {
	// Status filters tiers by status (optional)
	Status MembershipTierStatus `json:"status,omitempty"`
	// ShowDeleted indicates whether to show deleted tiers
	ShowDeleted bool `json:"show_deleted,omitempty"`
	// Page is the page number (starting from 1)
	Page int `json:"page,omitempty"`
	// Limit is the number of items per page
	Limit int `json:"limit,omitempty"`
}

// MembershipTierListResponse represents a paginated list of membership tiers
type MembershipTierListResponse struct {
	// Data is the list of membership tiers
	Data []MembershipTier `json:"data"`
	// Pagination contains pagination information
	Pagination PaginationInfo `json:"pagination"`
}

// CreateMembershipTier creates a new membership tier
//
// request: The tier creation request containing tier details and pricing
// Returns the created tier information and any error
func (c *Client) CreateMembershipTier(request *CreateMembershipTierRequest) (*MembershipTier, error) {
	var result MembershipTier
	err := c.requestJSON("POST", "/app/membership/tiers", request, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create membership tier: %w", err)
	}
	return &result, nil
}

// GetMembershipTier retrieves membership tier details by tier code
//
// code: The tier code to retrieve
// Returns the tier details and any error
func (c *Client) GetMembershipTier(code string) (*MembershipTier, error) {
	var result MembershipTier
	path := fmt.Sprintf("/app/membership/tiers/%s", url.PathEscape(code))
	err := c.requestJSON("GET", path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get membership tier: %w", err)
	}
	return &result, nil
}

// UpdateMembershipTier updates an existing membership tier
//
// code: The tier code to update
// request: The tier update request containing new tier details and pricing
// Returns the updated tier information and any error
func (c *Client) UpdateMembershipTier(code string, request *UpdateMembershipTierRequest) (*MembershipTier, error) {
	var result MembershipTier
	path := fmt.Sprintf("/app/membership/tiers/%s", url.PathEscape(code))
	err := c.requestJSON("PUT", path, request, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to update membership tier: %w", err)
	}
	return &result, nil
}

// DeleteMembershipTier deletes a membership tier by code
//
// code: The tier code to delete
// Returns any error encountered during deletion
func (c *Client) DeleteMembershipTier(code string) error {
	var result map[string]interface{}
	path := fmt.Sprintf("/app/membership/tiers/%s", url.PathEscape(code))
	err := c.requestJSON("DELETE", path, nil, &result)
	if err != nil {
		return fmt.Errorf("failed to delete membership tier: %w", err)
	}
	return nil
}

// RestoreMembershipTier restores a previously deleted membership tier
//
// code: The tier code to restore
// Returns the restored tier information and any error
func (c *Client) RestoreMembershipTier(code string) (*MembershipTier, error) {
	var result MembershipTier
	path := fmt.Sprintf("/app/membership/tiers/%s/restore", url.PathEscape(code))
	err := c.requestJSON("POST", path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to restore membership tier: %w", err)
	}
	return &result, nil
}

// ListMembershipTiers retrieves a paginated list of membership tiers
//
// request: The list request containing filter and pagination parameters
// Returns the tier list with pagination information and any error
func (c *Client) ListMembershipTiers(request *ListMembershipTiersRequest) (*MembershipTierListResponse, error) {
	// Build query parameters
	params := url.Values{}
	
	if request != nil {
		if request.Status != "" {
			params.Set("status", string(request.Status))
		}
		if request.ShowDeleted {
			params.Set("show_deleted", "true")
		}
		if request.Page > 0 {
			params.Set("page", strconv.Itoa(request.Page))
		}
		if request.Limit > 0 {
			params.Set("limit", strconv.Itoa(request.Limit))
		}
	}

	// Build path with query parameters
	path := "/app/membership/tiers"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result MembershipTierListResponse
	err := c.requestJSON("GET", path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to list membership tiers: %w", err)
	}
	return &result, nil
}

// GetMonthsByPeriodType returns the number of months for a given period type
//
// periodType: The membership period type
// Returns the number of months for the period
func GetMonthsByPeriodType(periodType MembershipPeriodType) int {
	switch periodType {
	case PeriodTypeMonth:
		return 1
	case PeriodTypeQuarter:
		return 3
	case PeriodTypeHalfYear:
		return 6
	case PeriodTypeYear:
		return 12
	case PeriodTypeTwoYear:
		return 24
	case PeriodTypeThreeYear:
		return 36
	case PeriodTypeFiveYear:
		return 60
	default:
		return 0
	}
}

// GetPeriodTypeName returns the display name for a given period type
//
// periodType: The membership period type
// Returns the human-readable name for the period
func GetPeriodTypeName(periodType MembershipPeriodType) string {
	switch periodType {
	case PeriodTypeMonth:
		return "Monthly"
	case PeriodTypeQuarter:
		return "Quarterly"
	case PeriodTypeHalfYear:
		return "Half-yearly"
	case PeriodTypeYear:
		return "Yearly"
	case PeriodTypeTwoYear:
		return "Two-year"
	case PeriodTypeThreeYear:
		return "Three-year"
	case PeriodTypeFiveYear:
		return "Five-year"
	default:
		return "Unknown Period"
	}
}