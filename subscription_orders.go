package woocommerce

import (
	"fmt"
)

const (
	subscriptionOrdersBasePath = "subscriptions/%v/orders"
)

// SubscriptionOrderService is an interface for interfacing with the subscriptionorders endpoints of woocommerce API
// https://woocommerce.github.io/woocommerce-rest-api-docs/#subscriptionorders
type SubscriptionOrderService interface {
	Create(subscriptionId int64, order Order) (*Order, error)
	Get(subscriptionId int64, orderId int64, options interface{}) (*Order, error)
	List(subscriptionId int64, options SubscriptionOrderListOptions) ([]Order, error)
	Update(subscriptionId int64, order *Order) (*Order, error)
	Delete(subscriptionId int64, subscriptioNorderID int64, options interface{}) (*Order, error)
	Batch(subscriptionId int64, option SubscriptionOrderBatchOption) (*SubscriptionOrderBatchResource, error)
}

// SubscriptionOrderServiceOp handles communication with the order related methods of WooCommerce'API
type SubscriptionOrderServiceOp struct {
	client *Client
}

// SubscriptionOrderListOption list all thee order list option request params
// refrence url:
// https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-subscriptionorders
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
// order	string	Order sort attribute ascending or descending. Options: asc and desc. Default is desc.
// orderby	string	Sort collection by object attribute. Options: date, id, include, title and slug. Default is date.
// parent	array	Limit result set to those of particular parent IDs.
// parent_exclude	array	Limit result set to all items except those of a particular parent ID.
// status	array	Limit result set to subscriptionorders assigned a specific status. Options: any, pending, processing, on-hold, completed, cancelled, refunded, failed and trash. Default is any.
// customer	integer	Limit result set to subscriptionorders assigned a specific customer.
// product	integer	Limit result set to subscriptionorders assigned a specific product.
// dp	integer	Number of decimal points to use in each resource. Default is 2.
type SubscriptionOrderListOptions struct {
	ListOptions
	Parent        []int64  `url:"parent,omitempty"`
	ParentExclude []int64  `url:"parent_exclude,omitempty"`
	Status        []string `url:"status,omitempty"`
	Customer      int64    `url:"customer,omitempty"`
	Product       int64    `url:"product,omitempty"`
	Dp            int      `url:"id,omitempty"`
}

// SubscriptionOrderBatchOption setting  operate for order in batch way
// https://woocommerce.github.io/woocommerce-rest-api-docs/#batch-update-subscriptionorders
type SubscriptionOrderBatchOption struct {
	Create []Order `json:"create,omitempty"`
	Update []Order `json:"update,omitempty"`
	Delete []int64 `json:"delete,omitempty"`
}

// SubscriptionOrderBatchResource conservation the response struct for SubscriptionOrderBatchOption request
type SubscriptionOrderBatchResource struct {
	Create []*Order `json:"create,omitempty"`
	Update []*Order `json:"update,omitempty"`
	Delete []*Order `json:"delete,omitempty"`
}

func (o *SubscriptionOrderServiceOp) List(subscriptionId int64, options SubscriptionOrderListOptions) ([]Order, error) {
	subscriptionorders, _, err := o.ListWithPagination(subscriptionId, options)
	return subscriptionorders, err
}

// ListWithPagination lists products and return pagination to retrieve next/previous results.
func (o *SubscriptionOrderServiceOp) ListWithPagination(subscriptionId int64, options interface{}) ([]Order, *Pagination, error) {
	resource := make([]Order, 0)
	basePath := fmt.Sprintf(subscriptionOrdersBasePath, subscriptionId)
	headers, err := o.client.createAndDoGetHeaders("GET", basePath, nil, options, &resource)
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

func (o *SubscriptionOrderServiceOp) Create(subscriptionId int64, order Order) (*Order, error) {
	basePath := fmt.Sprintf(subscriptionOrdersBasePath, subscriptionId)
	resource := new(Order)

	err := o.client.Post(basePath, &order, resource)
	return resource, err
}

// Get individual order
func (o *SubscriptionOrderServiceOp) Get(subscriptionId int64, orderID int64, options interface{}) (*Order, error) {
	path := fmt.Sprintf("%s/%d", fmt.Sprintf(subscriptionOrdersBasePath, subscriptionId), orderID)
	resource := new(Order)
	err := o.client.Get(path, resource, options)
	return resource, err
}

func (o *SubscriptionOrderServiceOp) Update(subscriptionId int64, order *Order) (*Order, error) {
	path := fmt.Sprintf("%s/%d", fmt.Sprintf(subscriptionOrdersBasePath, subscriptionId), order.ID)
	resource := new(Order)
	err := o.client.Put(path, order, &resource)
	return resource, err
}

func (o *SubscriptionOrderServiceOp) Delete(subscriptionId int64, subscriptionorderID int64, options interface{}) (*Order, error) {
	path := fmt.Sprintf("%s/%d", fmt.Sprintf(subscriptionOrdersBasePath, subscriptionId), subscriptionorderID)
	resource := new(Order)
	err := o.client.Delete(path, options, &resource)
	return resource, err
}

func (o *SubscriptionOrderServiceOp) Batch(subscriptionId int64, data SubscriptionOrderBatchOption) (*SubscriptionOrderBatchResource, error) {
	path := fmt.Sprintf("%s/batch", fmt.Sprintf(subscriptionOrdersBasePath, subscriptionId))
	resource := new(SubscriptionOrderBatchResource)
	err := o.client.Post(path, data, &resource)
	return resource, err
}
