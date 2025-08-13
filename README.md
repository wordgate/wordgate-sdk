# WordGate Go SDK

A Go SDK for interacting with the WordGate API, providing convenient methods for order management.

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
        NotifyURL:   "https://yoursite.com/webhook/payment",
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
    NotifyURL:   "https://yoursite.com/webhook/payment",   // Payment notification URL
    RedirectURL: "https://yoursite.com/payment/result",   // Post-payment redirect URL
})

// Get order details
order, err := client.GetOrder("ORDER123456")
```

## URL Configuration

### Notification URL (NotifyURL)

The notification URL is used for server-to-server webhook notifications when payment status changes. This should be a complete URL including protocol and domain:

```go
NotifyURL: "https://yoursite.com/api/payment/webhook"
```

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

## Configuration

Ensure your WordGate application is properly configured with:

- Valid app code and app secret
- Correct API base URL
- Proper webhook endpoint setup for notification URLs
- SSL/TLS certificates for HTTPS URLs

## License

This SDK is part of the WordGate project and follows the same licensing terms.