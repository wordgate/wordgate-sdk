# WordGate Webhook 签名验证规范

## 概述

WordGate 在发送 webhook 通知时，会使用 HMAC-SHA256 算法对请求进行签名，确保请求的完整性和来源真实性。接收方需要验证签名以确保请求来自 WordGate 且未被篡改。

## 签名机制

### 1. HTTP Header 设计

**Header 名称：**
```
X-Webhook-Signature
```

**Header 格式：**
```
t=<unix_timestamp>,sha256=<hmac_signature>
```

**示例：**
```
X-Webhook-Signature: t=1734315480,sha256=8b5a8f9d1c12f8a0d3e9a4a5b5c8e7e8b2f3f4c5a6b7d8e9f0a1b2c3d4e5f6
```

### 2. 签名生成规则（发送方 - WordGate）

1. 获取当前时间戳（Unix 秒级）
2. 将时间戳与 webhook 请求体原文拼接：
   ```
   message = "<timestamp>.<raw_body>"
   ```
3. 使用服务端配置的 **secret** 作为 key，用 **HMAC-SHA256** 算法对 `message` 进行签名
4. 将时间戳和签名按上面格式放入 `X-Webhook-Signature` header

**伪代码示例：**
```go
import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "strconv"
    "time"
)

func generateSignature(body []byte, secret string) string {
    timestamp := time.Now().Unix()
    message := strconv.FormatInt(timestamp, 10) + "." + string(body)
    
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(message))
    signature := hex.EncodeToString(h.Sum(nil))
    
    return fmt.Sprintf("t=%d,sha256=%s", timestamp, signature)
}
```

### 3. 签名验证规则（接收方 - 你的服务器）

1. 从 `X-Webhook-Signature` 中解析出 `t`（时间戳） 和 `sha256`（签名值）
2. 检查时间戳与当前时间差是否超过**5分钟**（防重放攻击）
3. 重新用本地保存的 **secret** 按发送方的算法生成签名
4. 比较生成的签名与 header 中的签名是否一致（建议用 `hmac.Equal` 做常量时间比较，防止时序攻击）

**伪代码示例：**
```go
func verifySignature(headerValue string, body []byte, secret string) error {
    // 1. 解析header
    parts := strings.Split(headerValue, ",")
    if len(parts) != 2 {
        return errors.New("invalid signature header format")
    }
    
    var timestamp int64
    var signature string
    
    for _, part := range parts {
        if strings.HasPrefix(part, "t=") {
            timestamp, _ = strconv.ParseInt(strings.TrimPrefix(part, "t="), 10, 64)
        } else if strings.HasPrefix(part, "sha256=") {
            signature = strings.TrimPrefix(part, "sha256=")
        }
    }
    
    // 2. 检查时间戳（防重放攻击）
    now := time.Now().Unix()
    if now-timestamp > 300 { // 5分钟
        return errors.New("timestamp too old")
    }
    
    // 3. 生成预期签名
    message := strconv.FormatInt(timestamp, 10) + "." + string(body)
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(message))
    expectedSignature := hex.EncodeToString(h.Sum(nil))
    
    // 4. 使用常量时间比较（防时序攻击）
    if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
        return errors.New("signature verification failed")
    }
    
    return nil
}
```

## 安全注意事项

### 1. 时间戳验证
- **必须验证时间戳**：检查请求时间戳与当前时间的差值不超过 5 分钟
- **防重放攻击**：避免恶意者重复发送旧的 webhook 请求

### 2. 签名比较
- **使用常量时间比较**：使用 `hmac.Equal()` 或类似函数进行签名比较
- **防时序攻击**：避免通过比较时间推断出正确的签名

### 3. Secret 管理
- **安全存储**：Webhook secret 应安全存储，不要硬编码在代码中
- **定期轮换**：建议定期更换 webhook secret
- **访问控制**：限制对 secret 的访问权限

### 4. 错误处理
- **记录失败**：记录签名验证失败的请求，用于安全监控
- **返回状态码**：
  - `200` - 验证成功并处理完成
  - `400` - 签名格式错误
  - `401` - 签名验证失败
  - `408` - 时间戳过期

## SDK 支持

WordGate 提供了 SDK 中的便捷方法来简化签名验证：

```go
import "github.com/wordgate/wordgate-sdk"

func handleWebhook(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadAll(r.Body)
    signature := r.Header.Get("X-Webhook-Signature")
    secret := "your_webhook_secret"
    
    // 使用 SDK 验证签名
    err := wordgate.VerifySignature(signature, body, secret, 300) // 5分钟超时
    if err != nil {
        http.Error(w, "签名验证失败", http.StatusUnauthorized)
        return
    }
    
    // 处理 webhook 事件
    var event wordgate.WebhookEventData
    json.Unmarshal(body, &event)
    
    // 处理业务逻辑...
    
    w.WriteHeader(http.StatusOK)
}
```

## 测试建议

### 1. 本地测试
- 使用 ngrok 等工具暴露本地服务
- 在 WordGate 后台配置测试 webhook URL
- 创建测试订单验证 webhook 接收

### 2. 签名测试
- 手动构造签名验证逻辑正确性
- 测试时间戳过期的情况
- 测试错误签名的拒绝逻辑

### 3. 监控告警
- 监控签名验证失败率
- 设置异常请求告警
- 记录可疑的验证失败

## 常见问题

### Q: 签名验证总是失败怎么办？
A: 请检查：
1. 请求体是否完整读取（注意不要预先解析）
2. Secret 是否与后台配置一致
3. 时间戳解析是否正确
4. 消息构造格式是否为 `timestamp.body`

### Q: 如何在开发环境调试签名？
A: 可以：
1. 打印生成的消息字符串和预期签名
2. 使用 WordGate SDK 的验证方法
3. 对比发送方和接收方的签名生成逻辑

### Q: 是否需要验证请求来源 IP？
A: 不建议依赖 IP 验证，因为：
1. WordGate 可能使用多个出口 IP
2. 签名验证已经足够保证安全性
3. IP 可能会变更

### Q: 重放攻击如何防范？
A: 通过时间戳验证：
1. 设置合理的时间窗口（推荐 5 分钟）
2. 可选择记录已处理的请求 ID 避免重复处理
3. 监控异常的时间戳模式