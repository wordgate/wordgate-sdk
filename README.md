# WordGate Go SDK

A Go SDK for interacting with the WordGate API, providing convenient methods for order management, product management, and membership tier management.

## Installation

```bash
go get github.com/wordgate/wordgate-sdk
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/wordgate/wordgate-sdk"
)

func main() {
    // Create a new client
    client := wordgate.NewClient(
        "your-app-code",
        "your-app-secret", 
        "https://api.wordgate.example.com",
    )
    
    // Create an order with notification and redirect URLs
    order, err := client.CreateOrder(&wordgate.CreateOrderRequest{
        Items: []wordgate.OrderItem{
            {
                ItemCode: "PRODUCT001",
                Quantity: 1,
                ItemType: "product",
            },
        },
        Customer: wordgate.OrderCustomer{
            Provider: "email",
            UID:      "user@example.com",
        },
        RedirectURL: "https://yoursite.com/payment/result",
    })
    if err != nil {
        log.Fatalf("Failed to create order: %v", err)
    }
    
    fmt.Printf("Order created: %s\n", order.OrderNo)
    fmt.Printf("Payment URL: %s\n", order.PayURL)
}
```

## Features

### Order Management

- **Create Orders**: Create new orders with items, customer information, and callback URLs
- **Get Order Details**: Retrieve order information by order number
- **Payment Integration**: Support for notification and redirect URLs

### Product Management

- **Create Products**: Create new products with pricing and shipping requirements
- **Update Products**: Modify product details and pricing
- **Delete/Restore Products**: Soft delete and restore products
- **List Products**: Query products with filtering and pagination
- **Product Status Management**: Active/inactive status control

### Membership Management

- **Create Membership Tiers**: Define membership levels with multiple pricing options
- **Update Tiers**: Modify tier details and pricing structures
- **Delete/Restore Tiers**: Manage tier lifecycle
- **List Tiers**: Query tiers with filtering and pagination
- **Flexible Pricing**: Support for multiple billing periods (monthly, quarterly, yearly, etc.)
- **Default Tier Management**: Designate default membership levels

### User Management

- **List Users**: Query users with filtering and pagination
- **Get User Details**: Retrieve comprehensive user information including membership and order history  
- **Update User Status**: Activate/deactivate user accounts
- **Set User Membership**: Assign membership tiers with custom expiration dates
- **Grant Membership**: Convenient methods to grant memberships for specific durations
- **Extend Membership**: Extend existing user memberships

## API Reference

### Client

```go
// Create a new WordGate API client
client := wordgate.NewClient(appCode, appSecret, baseURL)
```

### Orders

```go
// Create an order
order, err := client.CreateOrder(&wordgate.CreateOrderRequest{
    Items: []wordgate.OrderItem{
        {
            ItemCode: "PRODUCT001",
            Quantity: 1,
            ItemType: "product",
        },
    },
    Customer: wordgate.OrderCustomer{
        Provider: "email",
        UID:      "user@example.com",
    },
    RedirectURL: "https://yoursite.com/payment/result",   // Post-payment redirect URL
})

// Get order details
order, err := client.GetOrder("ORDER123456")
```

### Products

```go
// Create a product
product, err := client.CreateProduct(&wordgate.CreateProductRequest{
    Code:           "PRODUCT001",
    Name:           "Premium Package",
    Price:          9900,  // $99.00 in cents
    RequireAddress: false, // No shipping required
})

// Get product details
product, err := client.GetProduct("PRODUCT001")

// Update a product
product, err := client.UpdateProduct("PRODUCT001", &wordgate.UpdateProductRequest{
    Name:           "Premium Package Plus",
    Price:          14900, // $149.00 in cents
    RequireAddress: false,
})

// List products with filtering
products, err := client.ListProducts(&wordgate.ListProductsRequest{
    Status:      wordgate.ProductStatusActive,
    ShowDeleted: false,
    Page:        1,
    Limit:       20,
})

// Delete a product (soft delete)
err := client.DeleteProduct("PRODUCT001")

// Restore a deleted product
product, err := client.RestoreProduct("PRODUCT001")
```

### Membership Tiers

```go
// Create a membership tier with multiple pricing options
tier, err := client.CreateMembershipTier(&wordgate.CreateMembershipTierRequest{
    Code:      "PREMIUM",
    Name:      "Premium Membership",
    Level:     2,  // Higher level = higher tier
    IsDefault: false,
    Prices: []wordgate.MembershipPriceRequest{
        {
            PeriodType:    wordgate.PeriodTypeMonth,
            Price:         1900,  // $19.00 in cents
            OriginalPrice: 2900,  // $29.00 in cents
        },
        {
            PeriodType:    wordgate.PeriodTypeYear,
            Price:         19900, // $199.00 in cents
            OriginalPrice: 29900, // $299.00 in cents
        },
    },
})

// Get membership tier details
tier, err := client.GetMembershipTier("PREMIUM")

// Update a membership tier
tier, err := client.UpdateMembershipTier("PREMIUM", &wordgate.UpdateMembershipTierRequest{
    Name:      "Premium Plus Membership",
    Level:     3,
    IsDefault: false,
    Prices: []wordgate.MembershipPriceRequest{
        {
            PeriodType:    wordgate.PeriodTypeMonth,
            Price:         2900,
            OriginalPrice: 3900,
        },
    },
})

// List membership tiers
tiers, err := client.ListMembershipTiers(&wordgate.ListMembershipTiersRequest{
    Status:      wordgate.MembershipTierStatusActive,
    ShowDeleted: false,
    Page:        1,
    Limit:       10,
})

// Delete a membership tier
err := client.DeleteMembershipTier("PREMIUM")

// Restore a deleted tier
tier, err := client.RestoreMembershipTier("PREMIUM")
```

### User Management

```go
// List users with filtering and pagination
users, err := client.ListUsers(&wordgate.UserListRequest{
    Page:           1,
    Limit:          20,
    Email:          "user@example.com",  // Filter by email
    Status:         &[]int{1}[0],        // Filter by status (1=active, 0=disabled)
    MembershipTier: "PREMIUM",           // Filter by membership tier
    SortBy:         "created_at",        // Sort by field
    SortDesc:       true,                // Descending order
})

// Get detailed user information
userDetail, err := client.GetUser(12345)
// Returns user info, membership details, addresses, and recent orders

// Update user status (activate/deactivate)
err := client.UpdateUserStatus(12345, 1) // 1=active, 0=disabled

// Set user membership with custom expiration
response, err := client.SetUserMembership(12345, &wordgate.SetUserMembershipRequest{
    TierCode:  "PREMIUM",
    StartDate: "2024-01-01",  // Optional, defaults to current date
    EndDate:   "2024-12-31",  // Required
    OrderNo:   "ORDER123",    // Optional, for tracking
})

// Grant membership for specific duration (convenience method)
response, err := client.GrantUserMembership(12345, "PREMIUM", 30) // 30 days

// Grant membership until specific date
endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
response, err := client.GrantUserMembershipUntil(12345, "PREMIUM", endDate)

// Extend existing membership by specific duration
response, err := client.ExtendUserMembership(12345, "PREMIUM", 30) // Extend by 30 days
```

### Membership Period Types

The SDK supports various billing periods:

```go
wordgate.PeriodTypeMonth     // Monthly billing
wordgate.PeriodTypeQuarter   // Quarterly (3 months)
wordgate.PeriodTypeHalfYear  // Semi-annual (6 months)
wordgate.PeriodTypeYear      // Annual billing
wordgate.PeriodTypeTwoYear   // Biennial (2 years)
wordgate.PeriodTypeThreeYear // Triennial (3 years)
wordgate.PeriodTypeFiveYear  // 5-year billing

// Helper functions
months := wordgate.GetMonthsByPeriodType(wordgate.PeriodTypeYear) // Returns 12
name := wordgate.GetPeriodTypeName(wordgate.PeriodTypeMonth)      // Returns "Monthly"
```

## URL Configuration

### Webhook Notifications

Webhook notifications for payment events are configured at the application level in your WordGate dashboard, not per-order. This provides better security and centralized management.

Configure your webhook endpoint URL in your application settings to receive notifications for:
- Order payment success
- Order payment failure  
- Subscription events

### Redirect URL (RedirectURL)

The redirect URL is where users are redirected after completing payment. This should be a complete URL including protocol and domain:

```go
RedirectURL: "https://yoursite.com/payment/success"
```

## Error Handling

The SDK provides structured error handling:

```go
order, err := client.CreateOrder(request)
if err != nil {
    if apiErr, ok := err.(wordgate.APIError); ok {
        fmt.Printf("API Error (code %d): %s\n", apiErr.Code, apiErr.Message)
    } else {
        fmt.Printf("Network/Other Error: %v\n", err)
    }
    return
}
```

### Common Error Scenarios

```go
// Handle product creation conflicts
product, err := client.CreateProduct(request)
if err != nil {
    if apiErr, ok := err.(wordgate.APIError); ok {
        switch apiErr.Code {
        case 409:  // Conflict - product code already exists
            fmt.Println("Product code already exists")
        case 400:  // Bad request - validation error
            fmt.Println("Invalid product data:", apiErr.Message)
        case 401:  // Unauthorized
            fmt.Println("Authentication failed")
        default:
            fmt.Printf("API error: %s\n", apiErr.Message)
        }
    }
}

// Handle membership tier not found
tier, err := client.GetMembershipTier("NONEXISTENT")
if err != nil {
    fmt.Println("Membership tier not found:", err)
}
```

## Configuration

Ensure your WordGate application is properly configured with:

- Valid app code and app secret
- Correct API base URL
- Proper webhook endpoint setup for notification URLs
- SSL/TLS certificates for HTTPS URLs

## License

This SDK is part of the WordGate project and follows the same licensing terms.