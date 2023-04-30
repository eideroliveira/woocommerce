package woocommerce

import (
	"fmt"
	"strings"
	"time"
)

const (
	subscriptionsBasePath = "subscriptions"
)

// SubscriptionService is an interface for interfacing with the subscriptions endpoints of woocommerce API
// https://woocommerce.github.io/woocommerce-rest-api-docs/#subscriptions
type SubscriptionService interface {
	Create(subscription Subscription) (*Subscription, error)
	Get(subscriptionId int64, options interface{}) (*Subscription, error)
	List(options interface{}) ([]Subscription, error)
	Update(subscription *Subscription) (*Subscription, error)
	Delete(subscriptionID int64, options interface{}) (*Subscription, error)
	Batch(option SubscriptionBatchOption) (*SubscriptionBatchResource, error)
}

// SubscriptionServiceOp handles communication with the subscription related methods of WooCommerce'API
type SubscriptionServiceOp struct {
	client *Client
}

// SubscriptionListOption list all thee subscription list option request params
// refrence url:
// https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-subscriptions
// parameters:
// context	string	Scope under which the request is made; determines fields present in response. Options: view and edit. Default is view.
// page	integer	Current page of the collection. Default is 1.
// per_page	integer	Maximum number of items to be returned in result set. Default is 10.
// search	string	Limit results to those matching a string.
// after	string	Limit response to resources published after a given ISO8601 compliant date.
// before	string	Limit response to resources published before a given ISO8601 compliant date.
// exclude	array	Ensure result set excludes specific IDs.
// include	array	Limit result set to specific ids.
// offset	integer	Offset the result set by a specific number of items.
// order	string	Subscription sort attribute ascending or descending. Options: asc and desc. Default is desc.
// orderby	string	Sort collection by object attribute. Options: date, id, include, title and slug. Default is date.
// parent	array	Limit result set to those of particular parent IDs.
// parent_exclude	array	Limit result set to all items except those of a particular parent ID.
// status	array	Limit result set to subscriptions assigned a specific status. Options: any, pending, processing, on-hold, completed, cancelled, refunded, failed and trash. Default is any.
// customer	integer	Limit result set to subscriptions assigned a specific customer.
// product	integer	Limit result set to subscriptions assigned a specific product.
// dp	integer	Number of decimal points to use in each resource. Default is 2.
type SubscriptionListOptions struct {
	ListOptions
	Parent        []int64  `url:"parent,omitempty"`
	ParentExclude []int64  `url:"parent_exclude,omitempty"`
	Status        []string `url:"status,omitempty"`
	Customer      int64    `url:"customer,omitempty"`
	Product       int64    `url:"product,omitempty"`
	Dp            int      `url:"id,omitempty"`
}

// SubscriptionBatchOption setting  operate for subscription in batch way
// https://woocommerce.github.io/woocommerce-rest-api-docs/#batch-update-subscriptions
type SubscriptionBatchOption struct {
	Create []Subscription `json:"create,omitempty"`
	Update []Subscription `json:"update,omitempty"`
	Delete []int64        `json:"delete,omitempty"`
}

// SubscriptionBatchResource conservation the response struct for SubscriptionBatchOption request
type SubscriptionBatchResource struct {
	Create []*Subscription `json:"create,omitempty"`
	Update []*Subscription `json:"update,omitempty"`
	Delete []*Subscription `json:"delete,omitempty"`
}

// Subscription represents a WooCommerce Subscription
// https://woocommerce.github.io/woocommerce-rest-api-docs/#subscription-properties
type Subscription struct {
	ID                       int64           `json:"id,omitempty"`
	ParentId                 int64           `json:"parent_id,omitempty"`
	Status                   string          `json:"status,omitempty"`
	Currency                 string          `json:"currency,omitempty"`
	Version                  string          `json:"version,omitempty"`
	PricesIncludeTax         bool            `json:"prices_include_tax,omitempty"`
	DateCreated              CustomTime      `json:"date_created,omitempty"`
	DateCreatedGmt           string          `json:"date_created_gmt,omitempty"`
	DateModified             string          `json:"date_modified,omitempty"`
	DateModifiedGmt          CustomTime      `json:"date_modified_gmt,omitempty"`
	DateCompleted            string          `json:"date_completed,omitempty"`
	DateCompletedGmt         CustomTime      `json:"date_completed_gmt,omitempty"`
	DatePaid                 string          `json:"date_paid,omitempty"`
	DatePaidGmt              CustomTime      `json:"date_paid_gmt,omitempty"`
	StartDate                CustomTime      `json:"start_date,omitempty"`
	StartDateGmt             CustomTime      `json:"start_date_gmt,omitempty"`
	TrialEnd                 string          `json:"trial_end_date,omitempty"`
	TrialEndGmt              string          `json:"trial_end_date_gmt,omitempty"`
	NextPaymentDate          string          `json:"next_payment_date,omitempty"`
	NextPaymentDateGmt       CustomTime      `json:"next_payment_date_gmt,omitempty"`
	LastPaymentDate          string          `json:"last_payment_date,omitempty"`
	LastPaymentDateGmt       CustomTime      `json:"last_payment_date_gmt,omitempty"`
	PaymentRetryDate         CustomTime      `json:"payment_retry_date,omitempty"`
	PaymentRetryDateGmt      CustomTime      `json:"payment_retry_date_gmt,omitempty"`
	CancelledDate            CustomTime      `json:"cancelled_date,omitempty"`
	CancelledDateGmt         CustomTime      `json:"cancelled_date_gmt,omitempty"`
	EndDate                  string          `json:"end_date,omitempty"`
	EndDateGmt               CustomTime      `json:"end_date_gmt,omitempty"`
	DiscountsTotal           string          `json:"discount_total,omitempty"`
	DiscountsTax             string          `json:"discount_tax,omitempty"`
	ShippingTotal            string          `json:"shipping_total,omitempty"`
	ShippingTax              string          `json:"shipping_tax,omitempty"`
	CartTax                  string          `json:"cart_tax,omitempty"`
	Total                    string          `json:"total,omitempty"`
	TotalTax                 string          `json:"total_tax,omitempty"`
	CustomerId               int64           `json:"customer_id,omitempty"`
	OrderKey                 string          `json:"order_key,omitempty"`
	Billing                  *Billing        `json:"billing,omitempty"`
	Shipping                 *Shipping       `json:"shipping,omitempty"`
	PaymentMethod            string          `json:"payment_method,omitempty"`
	PaymentMethodTitle       string          `json:"payment_method_title,omitempty"`
	CustomerIpAddress        string          `json:"customer_ip_address,omitempty"`
	CustomerUserAgent        string          `json:"customer_user_agent,omitempty"`
	CreatedVia               string          `json:"created_via,omitempty"`
	CustomerNote             string          `json:"customer_note,omitempty"`
	Number                   string          `json:"number,omitempty"`
	MetaData                 []MetaData      `json:"meta_data,omitempty"`
	LineItems                []LineItem      `json:"line_items,omitempty"`
	TaxLines                 []TaxLine       `json:"tax_lines,omitempty"`
	ShippingLines            []ShippingLines `json:"shipping_lines,omitempty"`
	FeeLines                 []FeeLine       `json:"fee_lines,omitempty"`
	CouponLines              []CouponLine    `json:"coupon_lines,omitempty"`
	BillingPeriod            string          `json:"billing_period,omitempty"`
	BillingInterval          Stringint       `json:"billing_interval,omitempty"`
	ResubscribedFrom         string          `json:"resubscribed_from,omitempty"`
	ResubscribedSubscription string          `json:"resubscribed_subscription,omitempty"`
	RemovedLineItems         []LineItem      `json:"removed_line_items,omitempty"`
	PaymentDetails           PaymentDetails  `json:"payment_details,omitempty"`
	PaymentUrl               string          `json:"payment_url,omitempty"`
	TransitionStatus         string          `json:"transition_status,omitempty"`
	TrackingCode             string          `json:"correios_tracking_code,omitempty"`
	NeedsPayment             bool            `json:"needs_payment,omitempty"`
	NeedsProcessing          bool            `json:"needs_processing,omitempty"`
	IsEditable               bool            `json:"is_editable,omitempty"`
	Links                    Links           `json:"_links"`
}

type CustomTime struct {
	time.Time
}

const expiryDateLayout = "2006-01-02T15:04:05"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(expiryDateLayout, s)
	return
}

type PaymentDetails struct {
	PostMeta []MetaData `json:"post_meta,omitempty"`
	UserMeta []MetaData `json:"user_meta,omitempty"`
}

func (o *SubscriptionServiceOp) List(options interface{}) ([]Subscription, error) {
	subscriptions, _, err := o.ListWithPagination(options)
	return subscriptions, err
}

// ListWithPagination lists products and return pagination to retrieve next/previous results.
func (o *SubscriptionServiceOp) ListWithPagination(options interface{}) ([]Subscription, *Pagination, error) {
	resource := make([]Subscription, 0)
	headers, err := o.client.createAndDoGetHeaders("GET", subscriptionsBasePath, nil, options, &resource)
	if err != nil {
		return nil, nil, err
	}
	// Extract pagination info from header
	linkHeader := headers.Get("Link")
	pagination, err := extractPagination(linkHeader)
	if err != nil {
		return nil, nil, err
	}

	return resource, pagination, err
}

func (o *SubscriptionServiceOp) Create(subscription Subscription) (*Subscription, error) {
	resource := new(Subscription)

	err := o.client.Post(subscriptionsBasePath, subscription, &resource)
	return resource, err
}

// Get individual subscription
func (o *SubscriptionServiceOp) Get(subscriptionID int64, options interface{}) (*Subscription, error) {
	path := fmt.Sprintf("%s/%d", subscriptionsBasePath, subscriptionID)
	resource := new(Subscription)
	err := o.client.Get(path, resource, options)
	return resource, err
}

func (o *SubscriptionServiceOp) Update(subscription *Subscription) (*Subscription, error) {
	path := fmt.Sprintf("%s/%d", subscriptionsBasePath, subscription.ID)
	resource := new(Subscription)
	err := o.client.Put(path, subscription, &resource)
	return resource, err
}

func (o *SubscriptionServiceOp) Delete(subscriptionID int64, options interface{}) (*Subscription, error) {
	path := fmt.Sprintf("%s/%d", subscriptionsBasePath, subscriptionID)
	resource := new(Subscription)
	err := o.client.Delete(path, options, &resource)
	return resource, err
}

func (o *SubscriptionServiceOp) Batch(data SubscriptionBatchOption) (*SubscriptionBatchResource, error) {
	path := fmt.Sprintf("%s/batch", subscriptionsBasePath)
	resource := new(SubscriptionBatchResource)
	err := o.client.Post(path, data, &resource)
	return resource, err
}
