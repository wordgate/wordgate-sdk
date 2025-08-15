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
    
    // 创建商品订单
    productOrder, err := client.CreateAppProductOrder(&wordgate.CreateAppProductOrderRequest{
        Items: []struct {
            ItemCode string `json:"item_code"`
            Quantity int    `json:"quantity"`
        }{
            {
                ItemCode: "PRODUCT001",
                Quantity: 2,
            },
        },
        UserUID:     "user123",
        AddressID:   1,
        RedirectURL: "https://yoursite.com/payment/success",
    })
    if err != nil {
        log.Fatalf("创建商品订单失败: %v", err)
    }
    
    fmt.Printf("商品订单已创建: %s\n", productOrder.OrderNo)
    fmt.Printf("支付链接: %s\n", productOrder.PayURL)
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

### 📦 商品管理

#### 创建商品
```go
product, err := client.CreateProduct(&wordgate.CreateProductRequest{
    Code:           "PREMIUM_PLAN",
    Name:           "高级会员套餐",
    Price:          9900,  // 99.00 元（以分为单位）
    RequireAddress: false, // 数字商品不需要地址
})
```

#### 获取商品详情
```go
product, err := client.GetProduct("PREMIUM_PLAN")
```

#### 更新商品
```go
product, err := client.UpdateProduct("PREMIUM_PLAN", &wordgate.UpdateProductRequest{
    Name:           "超级高级会员套餐",
    Price:          14900, // 149.00 元
    RequireAddress: false,
})
```

#### 商品列表
```go
products, err := client.ListProducts(&wordgate.ListProductsRequest{
    Status:      "active",
    ShowDeleted: false,
    Page:        1,
    Limit:       20,
})
```

#### 删除/恢复商品
```go
// 软删除
err := client.DeleteProduct("PREMIUM_PLAN")

// 恢复已删除的商品
product, err := client.RestoreProduct("PREMIUM_PLAN")
```

### 💎 会员等级管理

#### 创建会员等级
```go
tier, err := client.CreateMembershipTier(&wordgate.CreateMembershipTierRequest{
    Code:      "VIP",
    Name:      "VIP 会员",
    Level:     2,  // 等级数值，越高等级越高
    IsDefault: false,
    Prices: []struct {
        PeriodType    wordgate.MembershipPeriodType `json:"period_type"`
        Price         int64                         `json:"price"`
        OriginalPrice int64                         `json:"original_price"`
    }{
        {
            PeriodType:    wordgate.PeriodTypeMonth,
            Price:         1900,  // 19.00 元/月
            OriginalPrice: 2900,  // 原价 29.00 元/月
        },
        {
            PeriodType:    wordgate.PeriodTypeYear,
            Price:         19900, // 199.00 元/年
            OriginalPrice: 29900, // 原价 299.00 元/年
        },
    },
})
```

#### 会员等级 CRUD 操作
```go
// 获取等级详情
tier, err := client.GetMembershipTier("VIP")

// 更新等级
tier, err := client.UpdateMembershipTier("VIP", &wordgate.UpdateMembershipTierRequest{
    Name:      "至尊 VIP 会员",
    Level:     3,
    IsDefault: false,
    // ... 价格配置
})

// 列出所有等级
tiers, err := client.ListMembershipTiers(&wordgate.ListMembershipTiersRequest{
    Status:      "active",
    ShowDeleted: false,
    Page:        1,
    Limit:       10,
})

// 删除等级
err := client.DeleteMembershipTier("VIP")

// 恢复等级
tier, err := client.RestoreMembershipTier("VIP")
```

#### 会员周期类型
```go
// 可用的周期类型
wordgate.PeriodTypeMonth     // 月付
wordgate.PeriodTypeQuarter   // 季付（3个月）
wordgate.PeriodTypeHalfYear  // 半年付（6个月）
wordgate.PeriodTypeYear      // 年付
wordgate.PeriodTypeTwoYear   // 两年付
wordgate.PeriodTypeThreeYear // 三年付
wordgate.PeriodTypeFiveYear  // 五年付
```

### 📝 订单管理

#### 创建商品订单
```go
productOrder, err := client.CreateAppProductOrder(&wordgate.CreateAppProductOrderRequest{
    Items: []struct {
        ItemCode string `json:"item_code"`
        Quantity int    `json:"quantity"`
    }{
        {
            ItemCode: "PREMIUM_PLAN",
            Quantity: 1,
        },
        {
            ItemCode: "ADDON_SERVICE",
            Quantity: 2,
        },
    },
    UserUID:     "user123",              // 用户唯一标识
    AddressID:   1,                      // 收货地址ID（可选）
    CouponCode:  "DISCOUNT10",           // 优惠券代码（可选）
    ClientIP:    "192.168.1.100",       // 客户端IP（可选）
    RedirectURL: "https://yoursite.com/payment/success", // 支付完成重定向URL
})
```

#### 创建会员订单
```go
membershipOrder, err := client.CreateAppMembershipOrder(&wordgate.CreateAppMembershipOrderRequest{
    TierID:      1,                      // 会员等级ID
    PeriodType:  "month",                // 周期类型
    UserUID:     "user123",              // 用户唯一标识
    CouponCode:  "VIP_DISCOUNT",         // 优惠券（可选）
    RedirectURL: "https://yoursite.com/membership/success",
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

#### 会员管理
```go
// 设置用户会员资格
response, err := client.SetUserMembership("user123", &wordgate.SetUserMembershipRequest{
    TierCode:  "VIP",
    StartDate: "2024-01-01",  // 可选，默认当前日期
    EndDate:   "2024-12-31",  // 到期日期
    OrderNo:   "ORDER123",    // 关联订单号（可选）
})

// 授予指定天数的会员资格
response, err := client.GrantUserMembership("user123", "VIP", 30) // 30天

// 授予会员资格至指定日期
endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
response, err := client.GrantUserMembershipUntil("user123", "VIP", endDate)

// 延长现有会员资格
response, err := client.ExtendUserMembership("user123", "VIP", 30) // 延长30天
```

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
// 处理商品代码重复
product, err := client.CreateProduct(request)
if err != nil {
    if apiErr, ok := err.(wordgate.APIError); ok && apiErr.Code == 409 {
        fmt.Println("商品代码已存在，请使用其他代码")
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
   - 已配置支付提供商（Stripe、Antom 等）
   - 正确的 Webhook 端点设置

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