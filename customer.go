package woocommerce

import (
	"fmt"
)

const (
	customersBasePath = "customers"
)

// CustomerService is an interface for interfacing with the customers endpoints of woocommerce API
// https://woocommerce.github.io/woocommerce-rest-api-docs/#customers
type CustomerService interface {
	Create(customer Customer) (*Customer, error)
	Get(customerId int64, options interface{}) (*Customer, error)
	List(options interface{}) ([]Customer, error)
	ListWithPagination(options interface{}) ([]Customer, *Pagination, error)
	Update(customer *Customer) (*Customer, error)
	Delete(customerID int64, options interface{}) (*Customer, error)
	Batch(option CustomerBatchOption) (*CustomerBatchResource, error)
}

// CustomerServiceOp handles communication with the customer related methods of WooCommerce'API
type CustomerServiceOp struct {
	client *Client
}

// CustomerListOption list all thee customer list option request params
// refrence url:
// https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-customers
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
// customer	string	Customer sort attribute ascending or descending. Options: asc and desc. Default is desc.
// customerby	string	Sort collection by object attribute. Options: date, id, include, title and slug. Default is date.
// parent	array	Limit result set to those of particular parent IDs.
// parent_exclude	array	Limit result set to all items except those of a particular parent ID.
// status	array	Limit result set to customers assigned a specific status. Options: any, pending, processing, on-hold, completed, cancelled, refunded, failed and trash. Default is any.
// customer	integer	Limit result set to customers assigned a specific customer.
// product	integer	Limit result set to customers assigned a specific product.
// dp	integer	Number of decimal points to use in each resource. Default is 2.
type CustomerListOption struct {
	ListOptions
	Parent        []int64  `url:"parent,omitempty"`
	ParentExclude []int64  `url:"parent_exclude,omitempty"`
	Status        []string `url:"status,omitempty"`
	Dp            int      `url:"id,omitempty"`
}

// CustomerBatchOption setting  operate for customer in batch way
// https://woocommerce.github.io/woocommerce-rest-api-docs/#batch-update-customers
type CustomerBatchOption struct {
	Create []Customer `json:"create,omitempty"`
	Update []Customer `json:"update,omitempty"`
	Delete []int64    `json:"delete,omitempty"`
}

// CustomerBatchResource conservation the response struct for CustomerBatchOption request
type CustomerBatchResource struct {
	Create []*Customer `json:"create,omitempty"`
	Update []*Customer `json:"update,omitempty"`
	Delete []*Customer `json:"delete,omitempty"`
}

type CustomerLastOrder struct {
	ID   int64  `json:"id,omitempty"`
	Date string `json:"date,omitempty"`
}

// Customer represents a WooCommerce Customer
// https://woocommerce.github.io/woocommerce-rest-api-docs/#customer-properties
type Customer struct {
	ID                int64                  `json:"id,omitempty"`
	AvatarURL         string                 `json:"avatar_url,omitempty"`
	Capabilities      map[string]interface{} `json:"capabilities,omitempty"`
	DateCreated       string                 `json:"date_created,omitempty"`
	DateCreatedGmt    string                 `json:"date_created_gmt,omitempty"`
	DateModified      string                 `json:"date_modified,omitempty"`
	DateModifiedGmt   string                 `json:"date_modified_gmt,omitempty"`
	LastOrder         CustomerLastOrder      `json:"last_order,omitempty"`
	OrdersCount       uint64                 `json:"orders_count,omitempty"`
	TotalSpent        string                 `json:"total_spent,omitempty"`
	Description       string                 `json:"description,omitempty"`
	Email             string                 `json:"email,omitempty"`
	ExtraCapabilities map[string]interface{} `json:"extra_capabilities,omitempty"`
	FirstName         string                 `json:"first_name,omitempty"`
	IsPayingCustomer  bool                   `json:"is_paying_customer,omitempty"`
	LastName          string                 `json:"last_name,omitempty"`
	Link              string                 `json:"link,omitempty"`
	Name              string                 `json:"name,omitempty"`
	Nickname          string                 `json:"nickname,omitempty"`
	RegisteredDate    string                 `json:"registered_date,omitempty"`
	Role              string                 `json:"role,omitempty"`
	Slug              string                 `json:"slug,omitempty"`
	URL               string                 `json:"url,omitempty"`
	Username          string                 `json:"username,omitempty"`
	Password          string                 `json:"password,omitempty"`
	Billing           *Billing               `json:"billing,omitempty"`
	Shipping          *Shipping              `json:"shipping,omitempty"`
	CartHash          string                 `json:"cart_hash,omitempty"`
	MetaData          []MetaData             `json:"meta_data,omitempty"`
	Links             Links                  `json:"_links"`
}

func (o *CustomerServiceOp) List(options interface{}) ([]Customer, error) {
	customers, _, err := o.ListWithPagination(options)
	return customers, err
}

// ListWithPagination lists products and return pagination to retrieve next/previous results.
func (o *CustomerServiceOp) ListWithPagination(options interface{}) ([]Customer, *Pagination, error) {
	path := customersBasePath
	resource := make([]Customer, 0)
	// headers := http.Header{}
	headers, err := o.client.createAndDoGetHeaders("GET", path, nil, options, &resource)
	if err != nil {
		return nil, nil, err
	}
	// Extract pagination info from header
	pagination, err := extractPagination(headers)
	if err != nil {
		return nil, nil, err
	}

	return resource, pagination, err
}

func (o *CustomerServiceOp) Create(customer Customer) (*Customer, error) {
	path := customersBasePath
	resource := new(Customer)

	err := o.client.Post(path, customer, &resource)
	return resource, err
}

// Get individual customer
func (o *CustomerServiceOp) Get(customerID int64, options interface{}) (*Customer, error) {
	path := fmt.Sprintf("%s/%d", customersBasePath, customerID)
	resource := new(Customer)
	err := o.client.Get(path, resource, options)
	return resource, err
}

func (o *CustomerServiceOp) Update(customer *Customer) (*Customer, error) {
	path := fmt.Sprintf("%s/%d", customersBasePath, customer.ID)
	resource := new(Customer)
	err := o.client.Put(path, customer, &resource)
	return resource, err
}

func (o *CustomerServiceOp) Delete(customerID int64, options interface{}) (*Customer, error) {
	path := fmt.Sprintf("%s/%d", customersBasePath, customerID)
	resource := new(Customer)
	err := o.client.Delete(path, options, &resource)
	return resource, err
}

func (o *CustomerServiceOp) Batch(data CustomerBatchOption) (*CustomerBatchResource, error) {
	path := fmt.Sprintf("%s/batch", customersBasePath)
	resource := new(CustomerBatchResource)
	err := o.client.Post(path, data, &resource)
	return resource, err
}
