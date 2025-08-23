# WordGate Go SDK

Go 语言的 WordGate 管理 API 客户端 SDK，专门用于后台管理系统和服务端集成。

## 📦 安装

```bash
go get github.com/wordgate/wordgate-sdk
```

## 🚀 快速开始

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/wordgate/wordgate-sdk"
)

func main() {
    // 创建客户端实例（使用应用管理凭据）
    client := wordgate.NewClient(
        "your-app-code",    // 应用代码
        "your-app-secret",  // 应用密钥
        "https://api.wordgate.example.com", // API 基础URL
    )
    
    // 创建自定义订单（灵活定价）
    customOrder, err := client.CreateAppCustomOrder(&wordgate.CreateAppCustomOrderRequest{
        UserUID:     "user123",
        Subject:     "自定义套餐订单",
        Description: "包含多项服务的自定义套餐",
        Amount:      19900, // 总金额：199.00元（分为单位）
        Currency:    "CNY",
        Items: []wordgate.CustomOrderItem{
            {
                ItemCode:    "CUSTOM_SERVICE_A",
                ItemName:    "自定义服务 A",
                Quantity:    1,
                UnitPrice:   9900,  // 99.00元
                RequireAddress: false,
            },
            {
                ItemCode:    "CUSTOM_SERVICE_B", 
                ItemName:    "自定义服务 B",
                Quantity:    1,
                UnitPrice:   10000, // 100.00元
                RequireAddress: false,
            },
        },
        RedirectURL: "https://yoursite.com/payment/success",
    })
    if err != nil {
        log.Fatalf("创建自定义订单失败: %v", err)
    }
    
    fmt.Printf("自定义订单已创建: %s\n", customOrder.OrderNo)
    fmt.Printf("支付链接: %s\n", customOrder.PayURL)
}
```

## 🏗️ 架构说明

### 管理 API 专用
本 SDK **仅调用管理接口**（`/app/*` 路径），专门设计用于：
- 后台管理系统
- 服务端到服务端的集成
- 自动化脚本和工具
- 管理员操作

### 认证方式
使用应用级认证：
- **App Code**: 应用标识符
- **App Secret**: 应用密钥
- 通过 `X-App-Code` 和 `X-App-Secret` HTTP 头部传递

## 📚 功能模块

### 📝 订单管理

#### 创建自定义订单（推荐）
自定义订单支持灵活定价，无需预先创建商品或会员等级，适用于各种业务场景：

```go
customOrder, err := client.CreateAppCustomOrder(&wordgate.CreateAppCustomOrderRequest{
    UserUID:     "user123",              // 用户唯一标识（必填）
    Subject:     "专业服务套餐",           // 订单标题（必填）
    Description: "包含咨询、开发、部署的完整服务", // 订单描述（可选）
    Amount:      29900,                  // 总金额：299.00元（必填，单位：分）
    Currency:    "CNY",                  // 货币类型（可选，默认应用货币）
    Items: []wordgate.CustomOrderItem{   // 订单项列表（必填）
        {
            ItemCode:       "CONSULTING",
            ItemName:       "专业咨询服务",
            Quantity:       10,         // 10小时
            UnitPrice:      1000,       // 10.00元/小时
            RequireAddress: false,      // 数字服务不需要地址
        },
        {
            ItemCode:       "DEVELOPMENT",
            ItemName:       "系统开发",
            Quantity:       1,
            UnitPrice:      19900,      // 199.00元
            RequireAddress: false,
        },
    },
    CouponCode:      "FIRST_ORDER",       // 优惠券代码（可选）
    ClientIP:        "192.168.1.100",    // 客户端IP（可选）
    AddressID:       0,                   // 收货地址ID（可选，数字商品可忽略）
    RedirectURL:     "https://yoursite.com/payment/success", // 支付完成重定向（可选）
    NotifyURL:       "https://yoursite.com/webhook/payment", // 支付通知URL（可选）
    RequireAddress:  false,               // 是否需要收货地址
})
```

#### 手动标记订单为已付款
适用于线下支付、银行转账等场景：

```go
err := client.MarkOrderAsPaid(&wordgate.ManualPaymentRequest{
    OrderNo:     "WG202401010001",       // 订单号（必填）
    PaymentNote: "银行转账，转账单号：ABC123456", // 付款说明（必填）
    Amount:      nil,                     // 付款金额（可选，默认为订单金额）
})
```

#### 查询订单详情
```go
orderDetail, err := client.GetAppOrder("WG202401010001")
```

#### 订单列表查询
```go
orders, err := client.ListAppOrders(&wordgate.ListOrdersQuery{
    Page:     1,
    Limit:    20,
    Status:   "paid",                    // 筛选已支付订单
    UserUID:  "user123",                 // 按用户筛选
    Email:    "user@example.com",        // 按邮箱筛选
    StartAt:  "2024-01-01",             // 开始日期
    EndAt:    "2024-01-31",             // 结束日期
    OrderNo:  "WG2024",                 // 订单号模糊匹配
    SortBy:   "created_at",             // 排序字段
    SortDesc: true,                     // 降序排列
})
```

### 👥 用户管理

#### 用户列表查询
```go
users, err := client.ListUsers(&wordgate.UserListRequest{
    Page:           1,
    Limit:          20,
    Email:          "user@example.com",  // 按邮箱筛选
    Status:         1,                   // 按状态筛选（1=激活，0=禁用）
    MembershipTier: "VIP",               // 按会员等级筛选
    SortBy:         "created_at",        // 排序字段
    SortDesc:       true,                // 降序排列
})
```

#### 查找或创建用户
```go
user, err := client.FindOrCreateUser(&wordgate.FindOrCreateUserRequest{
    Provider:  "email",
    UID:       "user@example.com",
    Email:     "user@example.com",
    Name:      "张三",
    AvatarURL: "https://example.com/avatar.jpg",
})
```

#### 获取用户详情
```go
// 包含用户信息、会员状态、地址、订单历史等
userDetail, err := client.GetUser("user123")
```

#### 用户状态管理
```go
// 激活用户
err := client.UpdateUserStatus("user123", 1)

// 禁用用户  
err := client.UpdateUserStatus("user123", 0)
```

## 🎯 核心特性

### ✨ 灵活的自定义订单
- **无依赖创建**：无需预先创建商品或会员等级
- **灵活定价**：支持任意金额和商品项组合
- **完整验证**：自动验证金额一致性，防止数据错误
- **支付兼容**：完全兼容 Stripe 支付提供商

### 🔒 安全的手动付款
- **管理员专用**：使用 App Secret 认证的管理接口
- **详细记录**：必须提供付款说明和相关信息
- **金额校验**：可选择验证付款金额是否与订单一致

### 📊 完整的订单管理
- **详细查询**：获取订单完整信息，包含用户、商品项、支付记录
- **灵活筛选**：支持多维度筛选和排序
- **实时状态**：准确反映订单和支付状态

## 🛠️ 错误处理

### 结构化错误处理
```go
product, err := client.CreateProduct(request)
if err != nil {
    if apiErr, ok := err.(wordgate.APIError); ok {
        fmt.Printf("API 错误 (代码 %d): %s\n", apiErr.Code, apiErr.Message)
        
        switch apiErr.Code {
        case 400:
            fmt.Println("请求参数无效")
        case 401:
            fmt.Println("认证失败，请检查 App Code 和 Secret")
        case 403:
            fmt.Println("权限不足")
        case 404:
            fmt.Println("资源不存在")
        case 409:
            fmt.Println("资源冲突，可能已存在")
        case 500:
            fmt.Println("服务器内部错误")
        }
    } else {
        fmt.Printf("网络或其他错误: %v\n", err)
    }
    return
}
```

### 常见错误场景
```go
// 处理订单金额不一致
customOrder, err := client.CreateAppCustomOrder(request)
if err != nil {
    if apiErr, ok := err.(wordgate.APIError); ok && apiErr.Code == 400 {
        fmt.Println("订单金额与商品项总额不一致，请检查金额计算")
        return
    }
}

// 处理用户不存在
userDetail, err := client.GetUser("nonexistent")
if err != nil {
    if apiErr, ok := err.(wordgate.APIError); ok && apiErr.Code == 404 {
        fmt.Println("用户不存在")
        return
    }
}

// 处理订单已支付
err := client.MarkOrderAsPaid(request)
if err != nil {
    if apiErr, ok := err.(wordgate.APIError); ok && apiErr.Code == 409 {
        fmt.Println("订单已经支付，无需重复标记")
        return
    }
}
```

## ⚙️ 配置要求

### 应用配置
确保您的 WordGate 应用已正确配置：

1. **应用凭据**
   - 有效的 App Code 和 App Secret
   - 在 WordGate 管理后台生成

2. **API 访问**
   - 正确的 API 基础 URL
   - HTTPS 连接（生产环境必需）

3. **支付配置**
   - 已配置 Stripe 支付提供商
   - 正确的 Webhook 端点设置
   - 自定义订单无需预配置商品或价格

### 安全注意事项
```go
// ❌ 不要在客户端代码中暴露
const AppSecret = "your-app-secret" // 危险！

// ✅ 从环境变量或配置文件读取
appSecret := os.Getenv("WORDGATE_APP_SECRET")
if appSecret == "" {
    log.Fatal("WORDGATE_APP_SECRET 环境变量未设置")
}

client := wordgate.NewClient(appCode, appSecret, baseURL)
```

## 🔗 相关链接

- [WordGate 主项目](https://github.com/wordgate/wordgate)
- [API 文档](https://docs.wordgate.example.com)
- [管理后台](https://dashboard.wordgate.example.com)

## 📄 许可证

本 SDK 采用与 WordGate 主项目相同的许可证条款。

## 🆘 技术支持

如果遇到问题：

1. 检查 [Issues](https://github.com/wordgate/wordgate/issues)
2. 提交新的 Issue
3. 查阅 WordGate 文档
4. 联系技术支持团队

---

**WordGate SDK** - 让后台管理更简单 🛠️