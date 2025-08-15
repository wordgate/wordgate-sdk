package wordgate

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// ProductStatus represents the status of a product
type ProductStatus string

const (
	// ProductStatusActive indicates the product is active and available
	ProductStatusActive ProductStatus = "active"
	// ProductStatusInactive indicates the product is inactive and unavailable
	ProductStatusInactive ProductStatus = "inactive"
)

// Product represents a product in the WordGate system
type Product struct {
	// ID is the unique identifier of the product
	ID uint64 `json:"id"`
	// AppID is the application ID this product belongs to
	AppID uint64 `json:"app_id"`
	// Code is the unique product code
	Code string `json:"code"`
	// Name is the product name
	Name string `json:"name"`
	// Price is the product price in cents
	Price int64 `json:"price"`
	// Status is the product status (active/inactive)
	Status ProductStatus `json:"status"`
	// RequireAddress indicates whether this product requires shipping address
	RequireAddress bool `json:"require_address"`
	// Version is the version number for optimistic locking
	Version int `json:"version"`
	// CreatedAt is the creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the last update timestamp
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt is the deletion timestamp (nil if not deleted)
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// CreateProductRequest represents a request to create a product
type CreateProductRequest struct {
	// Code is the unique product code
	Code string `json:"code" binding:"required,max=50"`
	// Name is the product name
	Name string `json:"name" binding:"required,max=100"`
	// Price is the product price in cents
	Price int64 `json:"price" binding:"required,min=0"`
	// RequireAddress indicates whether this product requires shipping address
	RequireAddress bool `json:"require_address"`
}

// UpdateProductRequest represents a request to update a product
type UpdateProductRequest struct {
	// Name is the product name
	Name string `json:"name" binding:"required,max=100"`
	// Price is the product price in cents
	Price int64 `json:"price" binding:"required,min=0"`
	// RequireAddress indicates whether this product requires shipping address
	RequireAddress bool `json:"require_address"`
}

// ListProductsRequest represents a request to list products
type ListProductsRequest struct {
	// Status filters products by status (optional)
	Status ProductStatus `json:"status,omitempty"`
	// ShowDeleted indicates whether to show deleted products
	ShowDeleted bool `json:"show_deleted,omitempty"`
	// Page is the page number (starting from 1)
	Page int `json:"page,omitempty"`
	// Limit is the number of items per page
	Limit int `json:"limit,omitempty"`
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	// CurrentPage is the current page number
	CurrentPage int `json:"current_page"`
	// PerPage is the number of items per page
	PerPage int `json:"per_page"`
	// Total is the total number of items
	Total int64 `json:"total"`
	// TotalPages is the total number of pages
	TotalPages int `json:"total_pages"`
}

// ProductListResponse represents a paginated list of products
type ProductListResponse struct {
	// Data is the list of products
	Data []Product `json:"data"`
	// Pagination contains pagination information
	Pagination PaginationInfo `json:"pagination"`
}

// CreateProduct creates a new product
//
// request: The product creation request containing product details
// Returns the created product information and any error
func (c *Client) CreateProduct(request *CreateProductRequest) (*Product, error) {
	var result Product
	err := c.requestJSON("POST", "/app/products", request, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return &result, nil
}

// GetProduct retrieves product details by product code
//
// code: The product code to retrieve
// Returns the product details and any error
func (c *Client) GetProduct(code string) (*Product, error) {
	var result Product
	path := fmt.Sprintf("/app/products/%s", url.PathEscape(code))
	err := c.requestJSON("GET", path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return &result, nil
}

// UpdateProduct updates an existing product
//
// code: The product code to update
// request: The product update request containing new product details
// Returns the updated product information and any error
func (c *Client) UpdateProduct(code string, request *UpdateProductRequest) (*Product, error) {
	var result Product
	path := fmt.Sprintf("/app/products/%s", url.PathEscape(code))
	err := c.requestJSON("PUT", path, request, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	return &result, nil
}

// DeleteProduct deletes a product by code
//
// code: The product code to delete
// Returns any error encountered during deletion
func (c *Client) DeleteProduct(code string) error {
	var result map[string]interface{}
	path := fmt.Sprintf("/app/products/%s", url.PathEscape(code))
	err := c.requestJSON("DELETE", path, nil, &result)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

// RestoreProduct restores a previously deleted product
//
// code: The product code to restore
// Returns the restored product information and any error
func (c *Client) RestoreProduct(code string) (*Product, error) {
	var result Product
	path := fmt.Sprintf("/app/products/%s/restore", url.PathEscape(code))
	err := c.requestJSON("POST", path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to restore product: %w", err)
	}
	return &result, nil
}

// ListProducts retrieves a paginated list of products
//
// request: The list request containing filter and pagination parameters
// Returns the product list with pagination information and any error
func (c *Client) ListProducts(request *ListProductsRequest) (*ProductListResponse, error) {
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
	path := "/app/products"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result ProductListResponse
	err := c.requestJSON("GET", path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return &result, nil
}