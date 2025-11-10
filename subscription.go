package wordgate

// SubscriptionPlan 订阅套餐
type SubscriptionPlan struct {
	ID           uint64 `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Price        int64  `json:"price"`
	Currency     string `json:"currency"`
	BillingType  string `json:"billing_type"`  // recurring/one_time
	BillingCycle string `json:"billing_cycle"` // monthly/quarterly/yearly
	Duration     int    `json:"duration"`      // 一次性订阅的天数
	TrialDays    int    `json:"trial_days"`    // 试用期天数
	Status       string `json:"status"`        // active/inactive
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateSubscriptionPlanRequest 创建订阅套餐请求
type CreateSubscriptionPlanRequest struct {
	Code         string            `json:"code" binding:"required,max=50"`
	Name         string            `json:"name" binding:"required,max=100"`
	Description  string            `json:"description" binding:"max=1000"`
	Price        int64             `json:"price" binding:"required,min=0"`
	Currency     string            `json:"currency" binding:"required,max=10"`
	BillingType  string            `json:"billing_type" binding:"required,oneof=recurring one_time"`
	BillingCycle string            `json:"billing_cycle" binding:"max=20"`
	Duration     int               `json:"duration" binding:"min=0"`
	TrialDays    int               `json:"trial_days" binding:"min=0"`
	Status       string            `json:"status" binding:"oneof=active inactive"`
	Metadata     map[string]string `json:"metadata"`
}

// UpdateSubscriptionPlanRequest 更新订阅套餐请求
type UpdateSubscriptionPlanRequest struct {
	Name         *string            `json:"name,omitempty" binding:"omitempty,max=100"`
	Description  *string            `json:"description,omitempty" binding:"omitempty,max=1000"`
	Price        *int64             `json:"price,omitempty" binding:"omitempty,min=0"`
	Currency     *string            `json:"currency,omitempty" binding:"omitempty,max=10"`
	BillingType  *string            `json:"billing_type,omitempty" binding:"omitempty,oneof=recurring one_time"`
	BillingCycle *string            `json:"billing_cycle,omitempty" binding:"omitempty,max=20"`
	Duration     *int               `json:"duration,omitempty" binding:"omitempty,min=0"`
	TrialDays    *int               `json:"trial_days,omitempty" binding:"omitempty,min=0"`
	Status       *string            `json:"status,omitempty" binding:"omitempty,oneof=active inactive"`
	Metadata     map[string]string  `json:"metadata,omitempty"`
}

// SyncSubscriptionPlansRequest 批量同步订阅套餐请求
type SyncSubscriptionPlansRequest struct {
	Plans []struct {
		Code         string `json:"code" binding:"required,max=50"`
		Name         string `json:"name" binding:"required,max=100"`
		Description  string `json:"description" binding:"max=1000"`
		Price        int64  `json:"price" binding:"required,min=0"`
		Currency     string `json:"currency" binding:"required,max=10"`
		BillingType  string `json:"billing_type" binding:"required,oneof=recurring one_time"`
		BillingCycle string `json:"billing_cycle" binding:"max=20"`
		Duration     int    `json:"duration" binding:"min=0"`
		TrialDays    int    `json:"trial_days" binding:"min=0"`
		Status       string `json:"status" binding:"oneof=active inactive"`
	} `json:"plans" binding:"required,min=1"`
}

// CreateSubscriptionRequest 创建订阅请求
type CreateSubscriptionRequest struct {
	PlanCode      string            `json:"plan_code" binding:"required,max=50"`
	PaymentMethod string            `json:"payment_method" binding:"required,max=50"`
	RedirectURL   string            `json:"redirect_url,omitempty" binding:"omitempty,max=500"`
	TrialDays     int               `json:"trial_days,omitempty" binding:"min=0"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// SubscriptionResponse 订阅响应
type SubscriptionResponse struct {
	ID                   uint64            `json:"id"`
	SubscriptionID       string            `json:"subscription_id"`
	OrderNo              string            `json:"order_no"`
	Plan                 SubscriptionPlan  `json:"plan"`
	Amount               int64             `json:"amount"`
	Currency             string            `json:"currency"`
	Status               string            `json:"status"`
	TrialEndDate         *string           `json:"trial_end_date"`
	CurrentPeriodStart   string            `json:"current_period_start"`
	CurrentPeriodEnd     string            `json:"current_period_end"`
	NextBillingDate      *string           `json:"next_billing_date"`
	AutoRenew            bool              `json:"auto_renew"`
	PayURL               string            `json:"pay_url"`
	CreatedAt            string            `json:"created_at"`
	UpdatedAt            string            `json:"updated_at"`
	Metadata             map[string]string `json:"metadata,omitempty"`
}

// SubscriptionListResponse 订阅列表响应
type SubscriptionListResponse struct {
	Items      []SubscriptionResponse `json:"items"`
	Pagination Pagination             `json:"pagination"`
}

// CancelSubscriptionRequest 取消订阅请求
type CancelSubscriptionRequest struct {
	Immediate bool   `json:"immediate"`          // 是否立即取消
	Reason    string `json:"reason,omitempty"`   // 取消原因
}

// UpdateAutoRenewRequest 更新自动续费请求
type UpdateAutoRenewRequest struct {
	AutoRenew bool `json:"auto_renew"`
}

// CreateUserSubscriptionRequest 管理员为用户创建订阅请求
type CreateUserSubscriptionRequest struct {
	PlanCode    string            `json:"plan_code" binding:"required,max=50"`
	StartDate   string            `json:"start_date,omitempty"`
	TrialDays   int               `json:"trial_days,omitempty" binding:"min=0"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// SubscriptionWebhookEventData 订阅 Webhook 事件数据
type SubscriptionWebhookEventData struct {
	AppCode              string            `json:"app_code"`
	SubscriptionID       string            `json:"subscription_id"`
	UserUID              string            `json:"user_uid"`
	Plan                 SubscriptionPlan  `json:"plan"`
	Status               string            `json:"status"`
	CurrentPeriodStart   string            `json:"current_period_start"`
	CurrentPeriodEnd     string            `json:"current_period_end"`
	NextBillingDate      *string           `json:"next_billing_date"`
	AutoRenew            bool              `json:"auto_renew"`
	TrialEndDate         *string           `json:"trial_end_date"`
	Metadata             map[string]string `json:"metadata,omitempty"`
}

// Webhook 事件常量
const (
	WebhookEventSubscriptionCreated   = "subscription.created"
	WebhookEventSubscriptionUpdated   = "subscription.updated" 
	WebhookEventSubscriptionRenewed   = "subscription.renewed"
	WebhookEventSubscriptionCancelled = "subscription.cancelled"
	WebhookEventSubscriptionExpired   = "subscription.expired"
)