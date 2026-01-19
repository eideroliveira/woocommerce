package woocommerce

import (
	"fmt"
)

const (
	couponsBasePath = "coupons"
)

// CouponService is an interface for interfacing with the coupons endpoints of woocommerce API
// https://woocommerce.github.io/woocommerce-rest-api-docs/#coupons
type CouponService interface {
	Create(coupon Coupon) (*Coupon, error)
	Get(couponID int64, options interface{}) (*Coupon, error)
	List(options interface{}) ([]Coupon, error)
	Update(coupon *Coupon) (*Coupon, error)
	Delete(couponID int64, options interface{}) (*Coupon, error)
	Batch(option CouponBatchOption) (*CouponBatchResource, error)
	ListWithPagination(options interface{}) ([]Coupon, *Pagination, error)
}

// CouponServiceOp handles communication with the coupon related methods of WooCommerce'API
type CouponServiceOp struct {
	client *Client
}

// CouponListOption list all the coupon list option request params
// refrence url:
// https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-coupons
type CouponListOption struct {
	ListOptions
}

// Coupon represents a WooCommerce Coupon
// https://woocommerce.github.io/woocommerce-rest-api-docs/#coupon-properties
type Coupon struct {
	ID                        int64      `json:"id,omitempty"`
	Code                      string     `json:"code,omitempty"`
	Slug                      string     `json:"slug"`
	Amount                    string     `json:"amount,omitempty"`
	DateCreated               StringTime `json:"date_created,omitempty"`
	DateCreatedGmt            StringTime `json:"date_created_gmt,omitempty"`
	DateModified              StringTime `json:"date_modified,omitempty"`
	DateModifiedGmt           StringTime `json:"date_modified_gmt,omitempty"`
	DiscountType              string     `json:"discount_type,omitempty"`
	Description               string     `json:"description,omitempty"`
	DateExpires               StringTime `json:"date_expires,omitempty"`
	DateExpiresGmt            StringTime `json:"date_expires_gmt,omitempty"`
	UsageCount                int        `json:"usage_count,omitempty"`
	IndividualUse             bool       `json:"individual_use,omitempty"`
	ProductIDs                []int      `json:"product_ids,omitempty"`
	ExcludedProductIDs        []int      `json:"excluded_product_ids,omitempty"`
	UsageLimit                int        `json:"usage_limit,omitempty"`
	UsageLimitPerUser         int        `json:"usage_limit_per_user,omitempty"`
	LimitUsageToXItems        int        `json:"limit_usage_to_x_items,omitempty"`
	FreeShipping              bool       `json:"free_shipping,omitempty"`
	ProductCategories         []int      `json:"product_categories,omitempty"`
	ExcludedProductCategories []int      `json:"excluded_product_categories,omitempty"`
	ExcludeSaleItems          bool       `json:"exclude_sale_items,omitempty"`
	MinimumAmount             string     `json:"minimum_amount,omitempty"`
	MaximumAmount             string     `json:"maximum_amount,omitempty"`
	NominalAmount             float64    `json:"nominal_amount,omitempty"`
	EmailRestrictions         []string   `json:"email_restrictions,omitempty"`
	UsedBy                    []string   `json:"used_by,omitempty"`
	MetaData                  []MetaData `json:"meta_data,omitempty"`
	Status                    string     `json:"status,omitempty"`
	Links                     Links      `json:"_links"`
}

type CouponBatchOption struct {
	Create []Coupon `json:"create,omitempty"`
	Update []Coupon `json:"update,omitempty"`
	Delete []int64  `json:"delete,omitempty"`
}

type CouponBatchResource struct {
	Create []*Coupon `json:"create,omitempty"`
	Update []*Coupon `json:"update,omitempty"`
	Delete []*Coupon `json:"delete,omitempty"`
}

func (c *CouponServiceOp) List(options interface{}) ([]Coupon, error) {
	coupons, _, err := c.ListWithPagination(options)
	return coupons, err
}

func (c *CouponServiceOp) ListWithPagination(options interface{}) ([]Coupon, *Pagination, error) {
	path := couponsBasePath
	resource := make([]Coupon, 0)
	headers, err := c.client.createAndDoGetHeaders("GET", path, nil, options, &resource)
	if err != nil {
		return nil, nil, err
	}
	pagination, err := extractPagination(headers)
	if err != nil {
		return nil, nil, err
	}
	return resource, pagination, err
}

func (c *CouponServiceOp) Create(coupon Coupon) (*Coupon, error) {
	path := couponsBasePath
	resource := new(Coupon)
	err := c.client.Post(path, coupon, &resource)
	return resource, err
}

func (c *CouponServiceOp) Get(couponID int64, options interface{}) (*Coupon, error) {
	path := fmt.Sprintf("%s/%d", couponsBasePath, couponID)
	resource := new(Coupon)
	err := c.client.Get(path, resource, options)
	return resource, err
}

func (c *CouponServiceOp) Update(coupon *Coupon) (*Coupon, error) {
	path := fmt.Sprintf("%s/%d", couponsBasePath, coupon.ID)
	resource := new(Coupon)
	err := c.client.Put(path, coupon, &resource)
	return resource, err
}

func (c *CouponServiceOp) Delete(couponID int64, options interface{}) (*Coupon, error) {
	path := fmt.Sprintf("%s/%d", couponsBasePath, couponID)
	resource := new(Coupon)
	err := c.client.Delete(path, options, &resource)
	return resource, err
}

func (c *CouponServiceOp) Batch(data CouponBatchOption) (*CouponBatchResource, error) {
	path := fmt.Sprintf("%s/batch", couponsBasePath)
	resource := new(CouponBatchResource)
	err := c.client.Post(path, data, &resource)
	return resource, err
}
