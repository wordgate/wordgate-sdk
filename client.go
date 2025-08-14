/*
Package wordgate provides a Go SDK for interacting with the WordGate API.

This SDK provides a simple and convenient way to integrate with WordGate services,
including order management, product management, membership management, and payment processing.

Basic usage examples:

	// Create a new client
	client := wordgate.NewClient("your-app-code", "your-app-secret", "https://api.wordgate.example.com")

	// Create a product order
	order, err := client.CreateProductOrder(&wordgate.CreateProductOrderRequest{
		Items: []wordgate.OrderItem{
			{
				ItemCode: "PRODUCT001",
				Quantity: 1,
			},
		},
		AddressID: 1,
	})
	if err != nil {
		log.Fatalf("Failed to create product order: %v", err)
	}

	fmt.Printf("Product order created: %s\n", order.OrderNo)

	// Create a membership order
	membershipOrder, err := client.CreateMembershipOrder(&wordgate.CreateMembershipOrderRequest{
		TierID:     1,
		PeriodType: "month",
	})
	if err != nil {
		log.Fatalf("Failed to create membership order: %v", err)
	}

	fmt.Printf("Membership order created: %s\n", membershipOrder.OrderNo)

	// Create a product
	product, err := client.CreateProduct(&wordgate.CreateProductRequest{
		Code:           "PRODUCT001",
		Name:           "Premium Package",
		Price:          9900, // $99.00 in cents
		RequireAddress: false,
	})
	if err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}

	fmt.Printf("Product created: %s\n", product.Code)

	// Create a membership tier
	tier, err := client.CreateMembershipTier(&wordgate.CreateMembershipTierRequest{
		Code:      "PREMIUM",
		Name:      "Premium Membership",
		Level:     2,
		IsDefault: false,
		Prices: []wordgate.MembershipPriceRequest{
			{
				PeriodType:    wordgate.PeriodTypeMonth,
				Price:         1900, // $19.00 in cents
				OriginalPrice: 2900, // $29.00 in cents
			},
			{
				PeriodType:    wordgate.PeriodTypeYear,
				Price:         19900, // $199.00 in cents
				OriginalPrice: 29900, // $299.00 in cents
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create membership tier: %v", err)
	}

	fmt.Printf("Membership tier created: %s\n", tier.Code)

	// List users
	users, err := client.ListUsers(&wordgate.UserListRequest{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}

	fmt.Printf("Found %d users\n", len(users.Items))

	// Get user details
	if len(users.Items) > 0 {
		userDetail, err := client.GetUser(users.Items[0].ID)
		if err != nil {
			log.Fatalf("Failed to get user details: %v", err)
		}

		fmt.Printf("User: %s (UID: %s)\n", userDetail.User.Nickname, userDetail.User.UID)

		// Set user membership
		membershipResponse, err := client.SetUserMembership(users.Items[0].ID, &wordgate.SetUserMembershipRequest{
			TierCode: "PREMIUM",
			EndDate:  "2024-12-31",
		})
		if err != nil {
			log.Fatalf("Failed to set user membership: %v", err)
		}

		fmt.Printf("User membership set: %s until %s\n", membershipResponse.TierName, membershipResponse.EndDate)

		// Grant membership for 30 days (convenience method)
		grantResponse, err := client.GrantUserMembership(users.Items[0].ID, "PREMIUM", 30)
		if err != nil {
			log.Fatalf("Failed to grant user membership: %v", err)
		}

		fmt.Printf("User membership granted: %s until %s\n", grantResponse.TierName, grantResponse.EndDate)
	}
*/
package wordgate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a WordGate API client
type Client struct {
	// AppCode is the application code for authentication
	AppCode string
	// AppSecret is the application secret for authentication
	AppSecret string
	// BaseURL is the base URL of the WordGate API
	BaseURL string
	// HTTPClient is the HTTP client used for requests
	HTTPClient *http.Client
}

// APIResponse represents a standard API response wrapper
type APIResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg,omitempty"`
}

// APIError represents an API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}


// Error implements the error interface for APIError
func (e APIError) Error() string {
	return fmt.Sprintf("API error (code %d): %s", e.Code, e.Message)
}

// NewClient creates a new WordGate API client
//
// appCode: The application code for authentication
// appSecret: The application secret for authentication
// baseURL: The base URL of the WordGate API (e.g., "https://api.wordgate.example.com")
func NewClient(appCode, appSecret, baseURL string) *Client {
	return &Client{
		AppCode:   appCode,
		AppSecret: appSecret,
		BaseURL:   baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// request performs an HTTP request to the API
//
// method: HTTP method (GET, POST, etc.)
// path: API endpoint path
// body: Request body (will be JSON encoded if not nil)
func (c *Client) request(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader

	// Encode request body as JSON if provided
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// Build full URL
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	// Create HTTP request
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("X-App-Code", c.AppCode)
	req.Header.Set("X-App-Secret", c.AppSecret)

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	return resp, nil
}

// requestJSON performs an HTTP request and unmarshals the JSON response
//
// method: HTTP method (GET, POST, etc.)
// path: API endpoint path
// body: Request body (will be JSON encoded if not nil)
// result: Pointer to the result structure
func (c *Client) requestJSON(method, path string, body interface{}, result interface{}) error {
	resp, err := c.request(method, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		// Try to parse as APIError
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err == nil && apiErr.Message != "" {
			return apiErr
		}
		// Fallback to HTTP error
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse API response wrapper
	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	// Check API response code
	if apiResp.Code != 0 {
		return APIError{
			Code:    apiResp.Code,
			Message: apiResp.Msg,
		}
	}

	// Marshal and unmarshal data field to target structure
	if result != nil && apiResp.Data != nil {
		dataBytes, err := json.Marshal(apiResp.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal API data: %w", err)
		}

		if err := json.Unmarshal(dataBytes, result); err != nil {
			return fmt.Errorf("failed to unmarshal API data: %w", err)
		}
	}

	return nil
}

