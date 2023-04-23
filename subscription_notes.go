package woocommerce

import (
	"fmt"
)

const (
	subscriptionNotesBasePath = "subscriptions/%v/notes"
)

// SubscriptionNoteService is an interface for interfacing with the subscriptionnotes endpoints of woocommerce API
// https://woocommerce.github.io/woocommerce-rest-api-docs/#subscriptionnotes
type SubscriptionNoteService interface {
	Create(subscriptionId int64, subscriptionNote SubscriptionNote) (*SubscriptionNote, error)
	Get(subscriptionId int64, subscriptionNoteId int64, options interface{}) (*SubscriptionNote, error)
	List(subscriptionId int64, options interface{}) ([]SubscriptionNote, error)
	Update(subscriptionId int64, subscriptioNnote *SubscriptionNote) (*SubscriptionNote, error)
	Delete(subscriptionId int64, subscriptioNnoteID int64, options interface{}) (*SubscriptionNote, error)
	Batch(subscriptionId int64, option SubscriptionNoteBatchOption) (*SubscriptionNoteBatchResource, error)
}

// SubscriptionNoteServiceOp handles communication with the subscriptionnote related methods of WooCommerce'API
type SubscriptionNoteServiceOp struct {
	client *Client
}

// SubscriptionNoteListOption list all thee subscriptionnote list option request params
// refrence url:
// https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-subscriptionnotes
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
// order	string	SubscriptionNote sort attribute ascending or descending. Options: asc and desc. Default is desc.
// orderby	string	Sort collection by object attribute. Options: date, id, include, title and slug. Default is date.
// parent	array	Limit result set to those of particular parent IDs.
// parent_exclude	array	Limit result set to all items except those of a particular parent ID.
// status	array	Limit result set to subscriptionnotes assigned a specific status. Options: any, pending, processing, on-hold, completed, cancelled, refunded, failed and trash. Default is any.
// customer	integer	Limit result set to subscriptionnotes assigned a specific customer.
// product	integer	Limit result set to subscriptionnotes assigned a specific product.
// dp	integer	Number of decimal points to use in each resource. Default is 2.
type SubscriptionNoteListOptions struct {
	ListOptions
	Parent        []int64 `url:"parent,omitempty"`
	ParentExclude []int64 `url:"parent_exclude,omitempty"`
}

// SubscriptionNoteBatchOption setting  operate for subscriptionnote in batch way
// https://woocommerce.github.io/woocommerce-rest-api-docs/#batch-update-subscriptionnotes
type SubscriptionNoteBatchOption struct {
	Create []SubscriptionNote `json:"create,omitempty"`
	Update []SubscriptionNote `json:"update,omitempty"`
	Delete []int64            `json:"delete,omitempty"`
}

// SubscriptionNoteBatchResource conservation the response struct for SubscriptionNoteBatchOption request
type SubscriptionNoteBatchResource struct {
	Create []*SubscriptionNote `json:"create,omitempty"`
	Update []*SubscriptionNote `json:"update,omitempty"`
	Delete []*SubscriptionNote `json:"delete,omitempty"`
}

// SubscriptionNote represents a WooCommerce SubscriptionNote
// https://woocommerce.github.io/woocommerce-rest-api-docs/#subscriptionnote-properties
type SubscriptionNote struct {
	ID             int64  `json:"id,omitempty"`
	DateCreated    string `json:"date_created,omitempty"`
	DateCreatedGmt string `json:"date_created_gmt,omitempty"`
	Note           string `json:"note,omitempty"`
	CustomerNote   bool   `json:"customer_note,omitempty"`
	AddedByUser    bool   `json:"added_by_user,omitempty"`
	Author         string `json:"author,omitempty"`
	Links          Links  `json:"_links"`
}

func (o *SubscriptionNoteServiceOp) List(subscriptionId int64, options interface{}) ([]SubscriptionNote, error) {
	subscriptionnotes, _, err := o.ListWithPagination(subscriptionId, options)
	return subscriptionnotes, err
}

// ListWithPagination lists products and return pagination to retrieve next/previous results.
func (o *SubscriptionNoteServiceOp) ListWithPagination(subscriptionId int64, options interface{}) ([]SubscriptionNote, *Pagination, error) {
	resource := make([]SubscriptionNote, 0)
	basePath := fmt.Sprintf(subscriptionNotesBasePath, subscriptionId)
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

func (o *SubscriptionNoteServiceOp) Create(subscriptionId int64, subscriptionNote SubscriptionNote) (*SubscriptionNote, error) {
	basePath := fmt.Sprintf(subscriptionNotesBasePath, subscriptionId)
	resource := new(SubscriptionNote)

	err := o.client.Post(basePath, &subscriptionNote, resource)
	return resource, err
}

// Get individual subscriptionnote
func (o *SubscriptionNoteServiceOp) Get(subscriptionId int64, subscriptionNoteID int64, options interface{}) (*SubscriptionNote, error) {
	path := fmt.Sprintf("%s/%d", fmt.Sprintf(subscriptionNotesBasePath, subscriptionId), subscriptionNoteID)
	resource := new(SubscriptionNote)
	err := o.client.Get(path, resource, options)
	return resource, err
}

func (o *SubscriptionNoteServiceOp) Update(subscriptionId int64, subscriptionnote *SubscriptionNote) (*SubscriptionNote, error) {
	path := fmt.Sprintf("%s/%d", fmt.Sprintf(subscriptionNotesBasePath, subscriptionId), subscriptionnote.ID)
	resource := new(SubscriptionNote)
	err := o.client.Put(path, subscriptionnote, &resource)
	return resource, err
}

func (o *SubscriptionNoteServiceOp) Delete(subscriptionId int64, subscriptionnoteID int64, options interface{}) (*SubscriptionNote, error) {
	path := fmt.Sprintf("%s/%d", fmt.Sprintf(subscriptionNotesBasePath, subscriptionId), subscriptionnoteID)
	resource := new(SubscriptionNote)
	err := o.client.Delete(path, options, &resource)
	return resource, err
}

func (o *SubscriptionNoteServiceOp) Batch(subscriptionId int64, data SubscriptionNoteBatchOption) (*SubscriptionNoteBatchResource, error) {
	path := fmt.Sprintf("%s/batch", fmt.Sprintf(subscriptionNotesBasePath, subscriptionId))
	resource := new(SubscriptionNoteBatchResource)
	err := o.client.Post(path, data, &resource)
	return resource, err
}
