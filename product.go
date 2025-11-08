package woocommerce

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var linkRegex = regexp.MustCompile(`^ *<([^>]+)>; rel="(prev|next|first|last)" *$`)

// ProductService allows you to create, view, update, and delete individual, or a batch, of products
// https://woocommerce.github.io/woocommerce-rest-api-docs/#products
type ProductService interface {
	Create(product Product) (*Product, error)
	Get(productID int64, options interface{}) (*Product, error)
	List(options interface{}) ([]Product, error)
	ListWithPagination(options interface{}) ([]Product, *Pagination, error)
	Update(product *Product) (*Product, error)
	Delete(productID int64, options interface{}) (*Product, error)
	Batch(option ProductBatchOption) (*ProductBatchResource, error)
}

// Product represent WooCommerce Product
// https://woocommerce.github.io/woocommerce-rest-api-docs/#product-properties
type Product struct {
	ID                int64              `json:"id"`
	Name              string             `json:"name"`
	Slug              string             `json:"slug"`
	Permalink         string             `json:"permalink"`
	DateCreated       StringTime         `json:"date_created"`
	DateModified      StringTime         `json:"date_modified"`
	Type              string             `json:"type"`
	Status            string             `json:"status"`
	Featured          bool               `json:"featured"`
	CatalogVisibility string             `json:"catalog_visibility"`
	Description       string             `json:"description"`
	ShortDescription  string             `json:"short_description"`
	SKU               string             `json:"sku"`
	Price             *StringFloat       `json:"price"`
	RegularPrice      *StringFloat       `json:"regular_price"`
	SalePrice         *StringFloat       `json:"sale_price"`
	DateOnSaleFrom    StringTime         `json:"date_on_sale_from"`
	DateOnSaleTo      StringTime         `json:"date_on_sale_to"`
	PriceHTML         string             `json:"price_html"`
	OnSale            bool               `json:"on_sale"`
	Purchasable       bool               `json:"purchasable"`
	TotalSales        StringFloat                `json:"total_sales"`
	Virtual           bool               `json:"virtual"`
	Visible           bool               `json:"visible"`
	Downloadable      bool               `json:"downloadable"`
	Downloads         []Download         `json:"downloads"`
	DownloadLimit     int                `json:"download_limit"`
	DownloadExpiry    int                `json:"download_expiry"`
	ExternalURL       string             `json:"external_url"`
	ButtonText        string             `json:"button_text"`
	TaxStatus         string             `json:"tax_status"`
	TaxClass          string             `json:"tax_class"`
	ManageStock       bool               `json:"manage_stock"`
	StockQuantity     *int               `json:"stock_quantity"`
	InStock           bool               `json:"in_stock"`
	Backorders        string             `json:"backorders"`
	BackordersAllowed bool               `json:"backorders_allowed"`
	Backordered       bool               `json:"backordered"`
	SoldIndividually  bool               `json:"sold_individually"`
	Weight            string             `json:"weight"`
	Dimensions        Dimensions         `json:"dimensions"`
	ShippingRequired  bool               `json:"shipping_required"`
	ShippingTaxable   bool               `json:"shipping_taxable"`
	ShippingClass     string             `json:"shipping_class"`
	ShippingClassId   int                `json:"shipping_class_id"`
	ReviewsAllowed    bool               `json:"reviews_allowed"`
	AverageRating     string             `json:"average_rating"`
	RatingCount       int                `json:"rating_count"`
	RelatedIds        []int              `json:"related_ids"`
	UpsellIds         []int              `json:"upsell_ids"`
	CrossSellIds      []int              `json:"cross_sell_ids"`
	ParentId          int                `json:"parent_id"`
	PurchaseNote      string             `json:"purchase_note"`
	Categories        []ProductCategory  `json:"categories"`
	Tags              []ProductTag       `json:"tags"`
	Image             []ProductImage     `json:"image"`
	Images            []ProductImage     `json:"images"`
	Attributes        []ProductAttribute `json:"attributes"`
	DefaultAttributes []ProductAttribute `json:"default_attributes"`
	Variations        []Product          `json:"variations"`
	GroupedProducts   []int              `json:"grouped_products"`
	MenuOrder         int                `json:"menu_order"`
	MetaData          []MetaDatum        `json:"meta"`
	DownloadType      string             `json:"download_type"`
	Links             Links              `json:"_links"`
}


type StringTime time.Time

func (i *StringTime) UnmarshalJSON(t []byte) (err error) {
	s := strings.Trim(string(t), "\"")
	if len(s) == 0 || s == "null" {
		return nil
	}

	for _, format := range []string{"2006-01-02T15:04:05", "2006-01-02", "02/01/2006"} {
		dateTime, err := time.Parse(format, s)
		if err == nil {
			*i = StringTime(dateTime)
			return nil
		}
	}
	return err
}

func (i *StringTime) MarshalJSON() ([]byte, error) {
	t := time.Time(*i)
	return t.MarshalJSON()
}

func (i *StringTime) Time() *time.Time {
	t := time.Time(*i)
	return &t
}

type StringFloat float32

func (i *StringFloat) UnmarshalJSON(t []byte) error {
	s := strings.Trim(string(t), "\"")
	if len(s) == 0 || s == "null" {
		return nil
	}

	f, err := strconv.ParseFloat(s, 32)
	*i = StringFloat(f)
	return err
}

func (i *StringFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(float32(*i))
}

func (i *StringFloat) Float32() float32 {
	return float32(*i)
}

func (i *StringFloat) Float64() float64 {
	return float64(*i)
}

func ParseStringFloat(s string) StringFloat {
	f, _ := strconv.ParseFloat(s, 32)
	return StringFloat(f)
}

type Dimensions struct {
	Length string `json:"length"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type Download struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	File string `json:"file"`
}

type ProductCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ProductTag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ProductImage struct {
	ID           int        `json:"id"`
	DateCreated  StringTime `json:"date_created"`
	DateModified StringTime `json:"date_modified"`
	Src          string     `json:"src"`
	Name         string     `json:"name"`
	Alt          string     `json:"alt"`
	Position     uint       `json:"position"`
}

type ProductAttribute struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Visible   bool     `json:"visible"`
	Variation bool     `json:"variation"`
	Option    string   `json:"option"`
	Options   []string `json:"options"`
}

type MetaDatum struct {
	ID    int         `json:"id"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Pagination of results
type Pagination struct {
	Total uint64
	TotalPages uint64
	NextPageOptions     *ListOptions
	PreviousPageOptions *ListOptions
	FirstPageOptions    *ListOptions
	LastPageOptions     *ListOptions
}

// ProductListOptions specifies the optional parameters to the List method.
type ProductListOptions struct {
	ListOptions
	// WooCommerce API doesn't seem to have specific filter options like parent,
	// status, or customer for products. You might need to use generic filters
	// like `search` or filter by attributes if needed.
	Search string `url:"search,omitempty"`
}

// ProductBatchOption sets options for batch operations on products.
type ProductBatchOption struct {
	Create []Product `json:"create,omitempty"`
	Update []Product `json:"update,omitempty"`
	Delete []int64   `json:"delete,omitempty"`
}

// ProductBatchResource holds the response struct for ProductBatchOption requests.
type ProductBatchResource struct {
	Create []*Product `json:"create,omitempty"`
	Update []*Product `json:"update,omitempty"`
	Delete []*Product `json:"delete,omitempty"`
}

type ProductServiceOp struct {
	client *Client
}

const productsBasePath = "products"

// List products.
func (o *ProductServiceOp) List(options interface{}) ([]Product, error) {
	products, _, err := o.ListWithPagination(options)
	return products, err
}

// ListWithPagination lists products and returns pagination to retrieve next/previous results.
func (o *ProductServiceOp) ListWithPagination(options interface{}) ([]Product, *Pagination, error) {
	resource := make([]Product, 0)
	headers, err := o.client.createAndDoGetHeaders("GET", productsBasePath, nil, options, &resource)
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

func (o *ProductServiceOp) Create(product Product) (*Product, error) {
	resource := new(Product)
	err := o.client.Post(productsBasePath, product, &resource)
	return resource, err
}

// Get individual product.
func (o *ProductServiceOp) Get(productID int64, options interface{}) (*Product, error) {
	path := fmt.Sprintf("%s/%d", productsBasePath, productID)
	resource := new(Product)
	err := o.client.Get(path, resource, options)
	return resource, err
}


func (o *ProductServiceOp) Update(product *Product) (*Product, error) {
	path := fmt.Sprintf("%s/%d", productsBasePath, product.ID)
	resource := new(Product)
	err := o.client.Put(path, product, &resource)
	return resource, err
}

func (o *ProductServiceOp) Delete(productID int64, options interface{}) (*Product, error) {
	path := fmt.Sprintf("%s/%d", productsBasePath, productID)
	resource := new(Product)
	err := o.client.Delete(path, options, &resource)
	return resource, err
}

func (o *ProductServiceOp) Batch(data ProductBatchOption) (*ProductBatchResource, error) {
	path := fmt.Sprintf("%s/batch", productsBasePath)
	resource := new(ProductBatchResource)
	err := o.client.Post(path, data, &resource)
	return resource, err
}

var log = &LeveledLogger{Level: LevelDebug}

// extractPagination extracts pagination info from linkHeader.
// Details on the format are here:
// https://woocommerce.github.io/woocommerce-rest-api-docs/#pagination
// Link: <https://www.example.com/wp-json/wc/v3/products?page=2>; rel="next",
// <https://www.example.com/wp-json/wc/v3/products?page=3>; rel="last"`
func extractPagination(headers http.Header) (*Pagination, error) {
	pagination := new(Pagination)
	linkHeader := headers.Get("Link")
		var err error
	if pagination.Total, err = strconv.ParseUint(headers.Get("X-Wp-Total"), 10, 0); err != nil {
		pagination.Total = 1
	}
	if pagination.TotalPages, err = strconv.ParseUint(headers.Get("X-Wp-Totalpages"), 10, 0); err != nil {
		pagination.TotalPages = 1
	}
	if linkHeader == "" {
		return pagination, nil
	}

	for _, link := range strings.Split(linkHeader, ",") {
		match := linkRegex.FindStringSubmatch(link)
		// Make sure the link is not empty or invalid
		if len(match) != 3 {
			// println("mm", len(match), " ", link, match)
			// We expect 3 values:
			// match[0] = full match
			// match[1] is the URL and match[2] is either 'previous' or 'next', 'first', 'last'
			// err := ResponseDecodingError{
			// 	Message: "could not extract pagination link header from " + linkHeader,
			// }
			return nil, nil
		}

		rel, err := url.Parse(match[1])
		if err != nil {
			err = ResponseDecodingError{
				Message: "pagination does not contain a valid URL",
			}
			return nil, err
		}

		params, err := url.ParseQuery(rel.RawQuery)
		if err != nil {
			return nil, err
		}

		paginationListOptions := ListOptions{}

		page := params.Get("page")
		if page != "" {
			paginationListOptions.Page, err = strconv.Atoi(params.Get("page"))
			if err != nil {
				return nil, err
			}
		}

		switch match[2] {
		case "next":
			pagination.NextPageOptions = &paginationListOptions
		case "prev":
			pagination.PreviousPageOptions = &paginationListOptions
		case "first":
			pagination.FirstPageOptions = &paginationListOptions
		case "last":
			pagination.LastPageOptions = &paginationListOptions
		}

	}

	return pagination, nil
}
