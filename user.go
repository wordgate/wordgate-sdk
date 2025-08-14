package wordgate

import (
	"fmt"
	"net/url"
	"time"
)

// User represents a user in the WordGate system
type User struct {
	// ID is the unique identifier of the user
	ID uint64 `json:"id"`
	// UID is the user's unique identifier
	UID string `json:"uid"`
	// Nickname is the user's display name
	Nickname string `json:"nickname"`
	// Avatar is the URL to the user's avatar image
	Avatar string `json:"avatar"`
	// Email is the user's email address
	Email string `json:"email"`
	// HasPassword indicates if the user has set a password
	HasPassword bool `json:"has_password"`
	// Status is the user status (1=active, 0=disabled)
	Status int `json:"status"`
	// LastLogin is the timestamp of the user's last login
	LastLogin time.Time `json:"last_login"`
	// CreatedAt is the user creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// Membership contains the user's current membership information
	Membership *UserMembershipInfo `json:"membership,omitempty"`
}

// UserMembershipInfo represents user membership information
type UserMembershipInfo struct {
	// TierName is the membership tier name
	TierName string `json:"tier_name"`
	// Status is the membership status (active/expired/canceled)
	Status string `json:"status"`
	// EndDate is the membership expiration date
	EndDate string `json:"end_date,omitempty"`
}

// UserListRequest represents a request to list users
type UserListRequest struct {
	// Page is the page number (starting from 1)
	Page int `json:"page,omitempty"`
	// Limit is the number of items per page (max 100)
	Limit int `json:"limit,omitempty"`
	// Email filters users by email address
	Email string `json:"email,omitempty"`
	// Nickname filters users by nickname
	Nickname string `json:"nickname,omitempty"`
	// Status filters users by status (1=active, 0=disabled)
	Status *int `json:"status,omitempty"`
	// StartAt filters users created after this date (YYYY-MM-DD)
	StartAt string `json:"start_at,omitempty"`
	// EndAt filters users created before this date (YYYY-MM-DD)
	EndAt string `json:"end_at,omitempty"`
	// MembershipTier filters users by membership tier code
	MembershipTier string `json:"membership_tier,omitempty"`
	// SortBy specifies the sort field (created_at/last_login)
	SortBy string `json:"sort_by,omitempty"`
	// SortDesc indicates whether to sort in descending order
	SortDesc bool `json:"sort_desc,omitempty"`
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	// Items is the list of users
	Items []User `json:"items"`
	// Pagination contains pagination information
	Pagination PaginationInfo `json:"pagination"`
}

// UserDetail represents detailed user information
type UserDetail struct {
	// User contains basic user information
	User UserDetailInfo `json:"user"`
	// Membership contains user membership details
	Membership UserMembershipDetail `json:"membership"`
	// Addresses contains user addresses
	Addresses []UserAddress `json:"addresses"`
	// Orders contains user orders with pagination
	Orders UserOrderList `json:"orders"`
}

// UserDetailInfo represents basic user information in detail view
type UserDetailInfo struct {
	// ID is the unique identifier of the user
	ID uint64 `json:"id"`
	// UID is the user's unique identifier
	UID string `json:"uid"`
	// Nickname is the user's display name
	Nickname string `json:"nickname"`
	// Avatar is the URL to the user's avatar image
	Avatar string `json:"avatar"`
	// HasPassword indicates if the user has set a password
	HasPassword bool `json:"has_password"`
	// Status is the user status (1=active, 0=disabled)
	Status int `json:"status"`
	// LastLogin is the timestamp of the user's last login
	LastLogin time.Time `json:"last_login"`
	// CreatedAt is the user creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// Identities contains login identities (email, phone, oauth, etc.)
	Identities []UserIdentity `json:"login_identities"`
}

// UserIdentity represents a user's login identity
type UserIdentity struct {
	// Provider is the identity provider (email, phone, wechat, etc.)
	Provider string `json:"provider"`
	// Identity is the identity unique identifier
	Identity string `json:"identity"`
	// Verified indicates if the identity is verified
	Verified bool `json:"verified"`
}

// UserMembershipDetail represents detailed membership information
type UserMembershipDetail struct {
	// Current contains current active membership
	Current *UserMembershipDetailItem `json:"current,omitempty"`
	// History contains membership history
	History []UserMembershipDetailItem `json:"history,omitempty"`
}

// UserMembershipDetailItem represents a membership detail item
type UserMembershipDetailItem struct {
	// ID is the membership record ID
	ID uint64 `json:"id"`
	// TierID is the membership tier ID
	TierID uint64 `json:"tier_id"`
	// TierName is the membership tier name
	TierName string `json:"tier_name"`
	// StartDate is the membership start date
	StartDate time.Time `json:"start_date"`
	// EndDate is the membership end date
	EndDate time.Time `json:"end_date"`
	// IsCanceled indicates if the membership is canceled
	IsCanceled bool `json:"is_canceled"`
	// OrderNo is the associated order number
	OrderNo string `json:"order_no,omitempty"`
}

// UserAddress represents a user's address
type UserAddress struct {
	// ID is the address ID
	ID uint64 `json:"id"`
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
	// IsDefault indicates if this is the default address
	IsDefault bool `json:"is_default"`
	// Label is the address label (home, office, etc.)
	Label string `json:"label"`
	// CreatedAt is the creation timestamp
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the last update timestamp
	UpdatedAt string `json:"updated_at"`
}

// UserOrderList represents paginated user orders
type UserOrderList struct {
	// Items is the list of orders
	Items []UserOrder `json:"items"`
	// Pagination contains pagination information
	Pagination PaginationInfo `json:"pagination"`
}

// UserOrder represents a user's order
type UserOrder struct {
	// OrderNo is the order number
	OrderNo string `json:"order_no"`
	// Amount is the order amount in cents
	Amount int64 `json:"amount"`
	// Currency is the order currency
	Currency string `json:"currency"`
	// IsPaid indicates if the order is paid
	IsPaid bool `json:"is_paid"`
	// PaidAt is the payment timestamp
	PaidAt *time.Time `json:"paid_at"`
	// CreatedAt is the order creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// ItemsSummary is a summary of order items
	ItemsSummary string `json:"items_summary"`
	// ItemsCount is the number of items in the order
	ItemsCount int `json:"items_count"`
}

// UpdateUserStatusRequest represents a request to update user status
type UpdateUserStatusRequest struct {
	// Status is the new user status (1=active, 0=disabled)
	Status int `json:"status"`
}

// SetUserMembershipRequest represents a request to set user membership
type SetUserMembershipRequest struct {
	// TierCode is the membership tier code
	TierCode string `json:"tier_code"`
	// StartDate is the membership start date (YYYY-MM-DD), optional (defaults to current date)
	StartDate string `json:"start_date,omitempty"`
	// EndDate is the membership end date (YYYY-MM-DD), required
	EndDate string `json:"end_date"`
	// OrderNo is the associated order number, optional
	OrderNo string `json:"order_no,omitempty"`
}

// SetUserMembershipResponse represents the response from setting user membership
type SetUserMembershipResponse struct {
	// Message is the response message
	Message string `json:"message"`
	// TierCode is the membership tier code that was set
	TierCode string `json:"tier_code"`
	// TierName is the membership tier name that was set
	TierName string `json:"tier_name"`
	// StartDate is the membership start date
	StartDate string `json:"start_date"`
	// EndDate is the membership end date
	EndDate string `json:"end_date"`
}

// FindOrCreateUserRequest represents a request to find or create a user
type FindOrCreateUserRequest struct {
	// Provider is the identity provider (email, phone, username, google, github)
	Provider string `json:"provider"`
	// Identity is the identity unique identifier
	Identity string `json:"identity"`
	// Nickname is the user's display name (optional)
	Nickname string `json:"nickname,omitempty"`
	// Avatar is the URL to the user's avatar image (optional)
	Avatar string `json:"avatar,omitempty"`
}

// FindOrCreateUserResponse represents the response from finding or creating a user
type FindOrCreateUserResponse struct {
	// User contains the user information
	User User `json:"user"`
	// Created indicates if a new user was created
	Created bool `json:"created"`
}

// ListUsers retrieves a paginated list of users
//
// request: The list request containing filter and pagination parameters
// Returns the user list with pagination information and any error
func (c *Client) ListUsers(request *UserListRequest) (*UserListResponse, error) {
	// Build query parameters
	params := url.Values{}
	
	if request != nil {
		if request.Page > 0 {
			params.Set("page", fmt.Sprintf("%d", request.Page))
		}
		if request.Limit > 0 {
			params.Set("limit", fmt.Sprintf("%d", request.Limit))
		}
		if request.Email != "" {
			params.Set("email", request.Email)
		}
		if request.Nickname != "" {
			params.Set("nickname", request.Nickname)
		}
		if request.Status != nil {
			params.Set("status", fmt.Sprintf("%d", *request.Status))
		}
		if request.StartAt != "" {
			params.Set("start_at", request.StartAt)
		}
		if request.EndAt != "" {
			params.Set("end_at", request.EndAt)
		}
		if request.MembershipTier != "" {
			params.Set("membership_tier", request.MembershipTier)
		}
		if request.SortBy != "" {
			params.Set("sort_by", request.SortBy)
		}
		if request.SortDesc {
			params.Set("sort_desc", "true")
		}
	}

	// Build path with query parameters
	path := "/app/users"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result UserListResponse
	err := c.requestJSON("GET", path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return &result, nil
}

// FindOrCreateUser finds an existing user or creates a new one
//
// request: The find or create user request containing identity information
// Returns the user information and creation status and any error
func (c *Client) FindOrCreateUser(request *FindOrCreateUserRequest) (*FindOrCreateUserResponse, error) {
	var result FindOrCreateUserResponse
	err := c.requestJSON("POST", "/app/users/find-or-create", request, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}
	return &result, nil
}

// GetUser retrieves user details by user UID
//
// userUID: The user UID to retrieve
// Returns the user details and any error
func (c *Client) GetUser(userUID string) (*UserDetail, error) {
	path := fmt.Sprintf("/app/users/%s", userUID)
	var result UserDetail
	err := c.requestJSON("GET", path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &result, nil
}

// UpdateUserStatus updates a user's status (active/disabled)
//
// userUID: The user UID to update
// status: The new status (1=active, 0=disabled)
// Returns any error encountered during the update
func (c *Client) UpdateUserStatus(userUID string, status int) error {
	path := fmt.Sprintf("/app/users/%s/status", userUID)
	request := UpdateUserStatusRequest{
		Status: status,
	}
	
	var result map[string]interface{}
	err := c.requestJSON("POST", path, request, &result)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	return nil
}

// SetUserMembership sets a user's membership with specified tier and expiration date
//
// userUID: The user UID to set membership for
// request: The membership setting request
// Returns the membership setting response and any error
func (c *Client) SetUserMembership(userUID string, request *SetUserMembershipRequest) (*SetUserMembershipResponse, error) {
	path := fmt.Sprintf("/app/users/%s/membership", userUID)
	
	var result SetUserMembershipResponse
	err := c.requestJSON("POST", path, request, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to set user membership: %w", err)
	}
	return &result, nil
}

// GrantUserMembership is a convenience method to grant membership to a user
//
// userUID: The user UID to grant membership to
// tierCode: The membership tier code to grant
// durationDays: The number of days the membership should last
// Returns the membership setting response and any error
func (c *Client) GrantUserMembership(userUID string, tierCode string, durationDays int) (*SetUserMembershipResponse, error) {
	now := time.Now()
	endDate := now.AddDate(0, 0, durationDays)
	
	request := &SetUserMembershipRequest{
		TierCode: tierCode,
		EndDate:  endDate.Format("2006-01-02"),
	}
	
	return c.SetUserMembership(userUID, request)
}

// GrantUserMembershipUntil is a convenience method to grant membership to a user until a specific date
//
// userUID: The user UID to grant membership to
// tierCode: The membership tier code to grant
// endDate: The date when the membership should expire
// Returns the membership setting response and any error
func (c *Client) GrantUserMembershipUntil(userUID string, tierCode string, endDate time.Time) (*SetUserMembershipResponse, error) {
	request := &SetUserMembershipRequest{
		TierCode: tierCode,
		EndDate:  endDate.Format("2006-01-02"),
	}
	
	return c.SetUserMembership(userUID, request)
}

// ExtendUserMembership extends a user's current membership by specified days
//
// userUID: The user UID to extend membership for
// tierCode: The membership tier code
// durationDays: The number of days to extend the membership
// Returns the membership setting response and any error
func (c *Client) ExtendUserMembership(userUID string, tierCode string, durationDays int) (*SetUserMembershipResponse, error) {
	// Get current user details to find existing membership end date
	userDetail, err := c.GetUser(userUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}
	
	var startDate time.Time
	if userDetail.Membership.Current != nil && userDetail.Membership.Current.EndDate.After(time.Now()) {
		// Extend from current end date if membership is still active
		startDate = userDetail.Membership.Current.EndDate
	} else {
		// Start from current time if no active membership
		startDate = time.Now()
	}
	
	endDate := startDate.AddDate(0, 0, durationDays)
	
	request := &SetUserMembershipRequest{
		TierCode:  tierCode,
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
	}
	
	return c.SetUserMembership(userUID, request)
}