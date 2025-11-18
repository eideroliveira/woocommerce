package woocommerce

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	ordersBasePath = "orders"
)

// OrderService is an interface for interfacing with the orders endpoints of woocommerce API
// https://woocommerce.github.io/woocommerce-rest-api-docs/#orders
type OrderService interface {
	Create(order Order) (*Order, error)
	Get(orderId int64, options interface{}) (*Order, error)
	List(options interface{}) ([]Order, error)
	Update(order *Order) (*Order, error)
	Delete(orderID int64, options interface{}) (*Order, error)
	Batch(option OrderBatchOption) (*OrderBatchResource, error)
	ListWithPagination(options interface{}) ([]Order, *Pagination, error)
}

// OrderServiceOp handles communication with the order related methods of WooCommerce'API
type OrderServiceOp struct {
	client *Client
}

// OrderListOption list all thee order list option request params
// refrence url:
// https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-orders
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
// status	array	Limit result set to orders assigned a specific status. Options: any, pending, processing, on-hold, completed, cancelled, refunded, failed and trash. Default is any.
// customer	integer	Limit result set to orders assigned a specific customer.
// product	integer	Limit result set to orders assigned a specific product.
// dp	integer	Number of decimal points to use in each resource. Default is 2.
type OrderListOption struct {
	ListOptions
	Parent        []int64  `url:"parent,omitempty"`
	ParentExclude []int64  `url:"parent_exclude,omitempty"`
	Status        []string `url:"status,omitempty"`
	Customer      int64    `url:"customer,omitempty"`
	Product       int64    `url:"product,omitempty"`
	Dp            int      `url:"id,omitempty"`
}

// OrderBatchOption setting  operate for order in batch way
// https://woocommerce.github.io/woocommerce-rest-api-docs/#batch-update-orders
type OrderBatchOption struct {
	Create []Order `json:"create,omitempty"`
	Update []Order `json:"update,omitempty"`
	Delete []int64 `json:"delete,omitempty"`
}

// OrderBatchResource conservation the response struct for OrderBatchOption request
type OrderBatchResource struct {
	Create []*Order `json:"create,omitempty"`
	Update []*Order `json:"update,omitempty"`
	Delete []*Order `json:"delete,omitempty"`
}

// Order represents a WooCommerce Order
// https://woocommerce.github.io/woocommerce-rest-api-docs/#order-properties
type Order struct {
	ID                 int64             `json:"id,omitempty"`
	ParentId           int64             `json:"parent_id,omitempty"`
	Number             string            `json:"number,omitempty"`
	OrderKey           string            `json:"order_key,omitempty"`
	CreatedVia         string            `json:"created_via,omitempty"`
	Version            string            `json:"version,omitempty"`
	Status             string            `json:"status,omitempty"`
	Currency           string            `json:"currency,omitempty"`
	DateCreated        StringTime        `json:"date_created,omitempty"`
	DateCreatedGmt     StringTime        `json:"date_created_gmt,omitempty"`
	DateModified       StringTime        `json:"date_modified,omitempty"`
	DateModifiedGmt    StringTime        `json:"date_modified_gmt,omitempty"`
	DiscountsTotal     StringFloat       `json:"discount_total,omitempty"`
	DiscountsTax       StringFloat       `json:"discount_tax,omitempty"`
	ShippingTotal      StringFloat       `json:"shipping_total,omitempty"`
	ShippingTax        StringFloat       `json:"shipping_tax,omitempty"`
	CartTax            StringFloat       `json:"cart_tax,omitempty"`
	Total              StringFloat       `json:"total,omitempty"`
	TotalTax           StringFloat       `json:"total_tax,omitempty"`
	PricesIncludeTax   bool              `json:"prices_include_tax,omitempty"`
	CustomerId         int64             `json:"customer_id,omitempty"`
	CustomerIpAddress  string            `json:"customer_ip_address,omitempty"`
	CustomerUserAgent  string            `json:"customer_user_agent,omitempty"`
	CustomerNote       string            `json:"customer_note,omitempty"`
	Billing            *Billing          `json:"billing,omitempty"`
	Shipping           *Shipping         `json:"shipping,omitempty"`
	PaymentMethod      string            `json:"payment_method,omitempty"`
	PaymentMethodTitle string            `json:"payment_method_title,omitempty"`
	TransactionId      string            `json:"transaction_id,omitempty"`
	DatePaid           StringTime        `json:"date_paid,omitempty"`
	DatePaidGmt        StringTime        `json:"date_paid_gmt,omitempty"`
	DateCompleted      StringTime        `json:"date_completed,omitempty"`
	DateCompletedGmt   StringTime        `json:"date_completed_gmt,omitempty"`
	CartHash           string            `json:"cart_hash,omitempty"`
	MetaData           []MetaData        `json:"meta,omitempty"`
	Renewal            StringInt         `json:"renewal,omitempty"`
	LineItems          []LineItem        `json:"line_items,omitempty"`
	TaxLines           []TaxLine         `json:"tax_lines,omitempty"`
	ShippingLines      []ShippingLines   `json:"shipping_lines,omitempty"`
	FeeLines           []FeeLine         `json:"fee_lines,omitempty"`
	CouponLines        []CouponLine      `json:"coupon_lines,omitempty"`
	Refunds            []Refund          `json:"refunds,omitempty"`
	PaymentUrl         string            `json:"payment_url,omitempty"`
	CurrencySymbol     string            `json:"currency_symbol,omitempty"`
	Links              Links             `json:"_links"`
	SetPaid            bool              `json:"set_paid,omitempty"`
	IsEditable         bool              `json:"is_editable,omitempty"`
	NeedsPayment       bool              `json:"needs_payment,omitempty"`
	NeedsProcessing    bool              `json:"needs_processing,omitempty"`
	TrackingCode       string            `json:"correios_tracking_code,omitempty"`
	OrderType          string            `json:"order_type,omitempty"`
	Paghiper           *WC_Paghiper_Data `json:"wc_paghiper_data,omitempty"`
	NFE                []NFE             `json:"nfe,omitempty"`
	Paypal             *PaypalData       `json:"paypal,omitempty"`
}

type WC_Paghiper_Data struct {
	OrderTransactionDueDate   StringTime  `json:"order_transaction_due_date,omitempty"`
	TransactionType           string      `json:"transaction_type,omitempty"`
	TransactionID             string      `json:"transaction_id,omitempty"`
	ValueCents                StringFloat `json:"value_cents,omitempty"`
	Status                    string      `json:"status,omitempty"`
	OrderID                   StringInt   `json:"order_id,omitempty"`
	CurrentTransactionDueDate StringTime  `json:"current_transaction_due_date,omitempty"`
	QrcodeBase64              string      `json:"qrcode_base64,omitempty"`
	QrcodeImageURL            string      `json:"qrcode_image_url,omitempty"`
	EMV                       string      `json:"emv,omitempty"`
	BacenURL                  string      `json:"bacen_url,omitempty"`
	PixURL                    string      `json:"pix_url,omitempty"`
	DigitableLine             string      `json:"digitable_line,omitempty"`
	URLSlip                   string      `json:"url_slip,omitempty"`
	URLSlipPDF                string      `json:"url_slip_pdf,omitempty"`
	Barcode                   string      `json:"barcode,omitempty"`
}

type NFE struct {
	UUID                 string     `json:"uuid,omitempty"`
	Status               string     `json:"status,omitempty"`
	Modelo               string     `json:"modelo,omitempty"`
	ChaveAcesso          string     `json:"chave_acesso,omitempty"`
	NRecibo              StringInt  `json:"n_recibo,omitempty"`
	NNFE                 StringInt  `json:"n_nfe,omitempty"`
	NSerie               StringOrInt  `json:"n_serie,omitempty"`
	NFEDoc               string     `json:"nfe_doc,omitempty"`
	PDF                  string     `json:"pdf,omitempty"`
	URLPDF               string     `json:"url_pdf,omitempty"`
	URLXML               string     `json:"url_xml,omitempty"`
	URLDanfe             string     `json:"url_danfe,omitempty"`
	URLDanfeSimplificada string     `json:"url_danfe_simplificada,omitempty"`
	URLDanfeEtiqueta     string     `json:"url_danfe_etiqueta,omitempty"`
	PDFRPS               string     `json:"pdf_rps,omitempty"`
	Data                 StringTime `json:"data,omitempty"`
}

type PaypalData struct {
	TransactionID  string     `json:"transaction_id,omitempty"`
	PaidDate       StringTime `json:"paid_date,omitempty"`
	CompletedDate  StringTime `json:"completed_date,omitempty"`
	SubscriptionID string     `json:"paypal_subscription_id,omitempty"`
	Status         string     `json:"paypal_status,omitempty"`
	IPNTrackingIDs []string   `json:"ipn_tracking_ids,omitempty"`
}

type Links struct {
	Self []struct {
		Href string `json:"href"`
	} `json:"self"`
	Collection []struct {
		Href string `json:"href"`
	} `json:"collection"`
	Customer []struct {
		Href string `json:"href"`
	} `json:"customer"`
	Up []struct {
		Href string `json:"href"`
	} `json:"up"`
}

type PersonType uint

const (
	CustomerTypeUnknown PersonType = iota
	CustomerPessoaFisica
	CustomerPessoaJuridica
)

type Billing struct {
	FirstName      string     `json:"first_name,omitempty"`
	LastName       string     `json:"last_name,omitempty"`
	Company        string     `json:"company,omitempty"`
	Address1       string     `json:"address_1,omitempty"`
	Address2       string     `json:"address_2,omitempty"`
	City           string     `json:"city,omitempty"`
	State          string     `json:"state,omitempty"`
	PostCode       string     `json:"postcode,omitempty"`
	Country        string     `json:"country,omitempty"`
	Email          string     `json:"email,omitempty"`
	BillingEmail   string     `json:"billing_email,omitempty"`
	BillingCompany string     `json:"billing_company,omitempty"`
	Phone          string     `json:"phone,omitempty"`
	CPF            string     `json:"cpf,omitempty"`
	RG             string     `json:"rg,omitempty"`
	CNPJ           string     `json:"cnpj,omitempty"`
	IE             string     `json:"ie,omitempty"`
	Number         string     `json:"number,omitempty"`
	Neighborhood   string     `json:"neighborhood,omitempty"`
	PersonType     PersonType `json:"persontype,omitempty"`
	BirthDate      string     `json:"birthdate,omitempty"`
	CellPhone      string     `json:"cellphone,omitempty"`
	Sex            string     `json:"gender,omitempty"`
	ChurchEmail    string     `json:"church_email,omitempty"`
	ChurchSize     StringInt  `json:"church_size,omitempty"`
	PayerName      string     `json:"payer_name,omitempty"`
	PayerEmail     string     `json:"payer_email,omitempty"`
	PayerPhone     string     `json:"payer_phone,omitempty"`
	Church         string     `json:"church,omitempty"`
}

func (c *Billing) String() string {
	res := []string{}
	for i := 0; i < reflect.TypeOf(*c).NumField(); i++ {
		l := reflect.ValueOf(*c).Field(i)
		if !l.IsZero() {
			f := reflect.TypeOf(*c).Field(i)
			res = append(res, fmt.Sprintf("%v: %v", f.Name, l))
		}
	}
	return strings.Join(res, "\n")
}

type Shipping struct {
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Company      string `json:"company,omitempty"`
	Address1     string `json:"address_1,omitempty"`
	Address2     string `json:"address_2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostCode     string `json:"postcode,omitempty"`
	Country      string `json:"country,omitempty"`
	Number       string `json:"number,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Email        string `json:"email,omitempty"`
}

type LineItem struct {
	ID          int64       `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	ProductID   int64       `json:"product_id,omitempty"`
	VariantID   int64       `json:"variation_id,omitempty"`
	Quantity    int         `json:"quantity,omitempty"`
	TaxClass    string      `json:"tax_class,omitempty"`
	SubTotal    string      `json:"subtotal,omitempty"`
	SubtotalTax string      `json:"subtotal_tax,omitempty"`
	Total       string      `json:"total,omitempty"`
	TotalTax    string      `json:"total_tax,omitempty"`
	Taxes       []TaxLine   `json:"taxes,omitempty"`
	MetaData    []MetaData  `json:"meta,omitempty"`
	SKU         string      `json:"sku,omitempty"`
	Price       StringFloat `json:"price,omitempty"`
	Image       Image       `json:"image,omitempty"`
	ParentName  string      `json:"parent_name,omitempty"`
}

func (p *PersonType) UnmarshalJSON(id []byte) error {
	s := string(id)
	if s == "" || s == "\"\"" {
		*p = CustomerTypeUnknown
		return nil
	}
	switch s {
	case "F":
		*p = CustomerPessoaFisica
	case "J":
		*p = CustomerPessoaJuridica
	default:
		*p = CustomerTypeUnknown
	}
	return nil
}

type StringInt int64

func (i *StringInt) UnmarshalJSON(id []byte) error {
	s := string(id)
	if s == "" || s == "\"\"" {
		*i = StringInt(0)
		return nil
	}
	i_, err := strconv.Atoi(strings.Trim(strings.ReplaceAll(string(id), `"`, ""), " "))
	*i = StringInt(i_)
	if err != nil {
		fmt.Printf("error parsing stringint %s: %v, setting to Zero", id, err)
	}
	return nil
}

func (i *StringInt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", *i)), nil
}

func (i StringInt) String() string {
	return fmt.Sprintf("%d", i)
}

func (i StringInt) Int64() int64 {
	return int64(i)
}

type StringOrInt string

func (i *StringOrInt) UnmarshalJSON(id []byte) error {
	s := string(id)
	if s == "" || s == "\"\"" {
		*i = StringOrInt("")
		return nil
	}
	i_, err := strconv.Atoi(strings.Trim(strings.ReplaceAll(s, `"`, ""), " "))
	if err == nil {
		*i = StringOrInt(fmt.Sprintf("%d", i_))
	} else {
		*i = StringOrInt(s)
	}
	return nil
}

func (i StringOrInt) MarshalJSON() ([]byte, error) {
	return []byte(i), nil
}

func (i StringOrInt) String() string {
	return string(i)
}

func (i StringOrInt) Int64() int64 {
	i_, _ := strconv.Atoi(string(i))
	return int64(i_)
}

type Image struct {
	ID  StringInt `json:"id,omitempty"`
	Src string    `json:"src,omitempty"`
}

type TaxLine struct {
	ID               int64      `json:"id,omitempty"`
	RateCode         string     `json:"rate_code,omitempty"`
	RateId           string     `json:"rate_id,omitempty"`
	Label            string     `json:"label,omitempty"`
	Compound         bool       `json:"compound,omitempty"`
	TaxTotal         string     `json:"tax_total"`
	ShippingTaxTotal string     `json:"shipping_tax_total,omitempty"`
	MetaData         []MetaData `json:"meta,omitempty"`
}

type MetaData struct {
	ID           int64       `json:"id,omitempty"`
	Key          string      `json:"key,omitempty"`
	Value        interface{} `json:"value,omitempty"`
	Label        string      `json:"label,omitempty"`
	DisplayKey   string      `json:"display_key,omitempty"`
	DisplayValue interface{} `json:"display_value,omitempty"`
}

type FeeLine struct {
	ID        int64      `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	TaxClass  string     `json:"tax_class,omitempty"`
	TaxStatus string     `json:"tax_status,omitempty"`
	Amount    string     `json:"amount,omitempty"`
	Total     string     `json:"total,omitempty"`
	TotalTax  string     `json:"total_tax,omitempty"`
	Taxes     []TaxLine  `json:"taxes,omitempty"`
	MetaData  []MetaData `json:"meta,omitempty"`
}

type Refund struct {
	ID     int64  `json:"id,omitempty"`
	Reason string `json:"refund,omitempty"`
	Total  string `json:"total,omitempty"`
}

type ShippingLines struct {
	ID          int64      `json:"id,omitempty"`
	MethodTitle string     `json:"method_title,omitempty"`
	MethodID    string     `json:"method_id,omitempty"`
	Total       string     `json:"total,omitempty"`
	TotalTax    string     `json:"total_tax,omitempty"`
	Taxes       []TaxLine  `json:"taxes,omitempty"`
	MetaData    []MetaData `json:"meta,omitempty"`
}

type CouponLine struct {
	ID          int64      `json:"id,omitempty"`
	Code        string     `json:"code,omitempty"`
	Discount    string     `json:"discount,omitempty"`
	DiscountTax string     `json:"discount_tax,omitempty"`
	MetaData    []MetaData `json:"meta,omitempty"`
}

func (o *OrderServiceOp) List(options interface{}) ([]Order, error) {
	orders, _, err := o.ListWithPagination(options)
	return orders, err
}

// ListWithPagination lists products and return pagination to retrieve next/previous results.
func (o *OrderServiceOp) ListWithPagination(options interface{}) ([]Order, *Pagination, error) {
	path := ordersBasePath
	resource := make([]Order, 0)
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

func (o *OrderServiceOp) Create(order Order) (*Order, error) {
	path := ordersBasePath
	resource := new(Order)

	err := o.client.Post(path, order, &resource)
	return resource, err
}

// Get individual order
func (o *OrderServiceOp) Get(orderID int64, options interface{}) (*Order, error) {
	path := fmt.Sprintf("%s/%d", ordersBasePath, orderID)
	resource := new(Order)
	err := o.client.Get(path, resource, options)
	return resource, err
}

func (o *OrderServiceOp) Update(order *Order) (*Order, error) {
	path := fmt.Sprintf("%s/%d", ordersBasePath, order.ID)
	resource := new(Order)
	err := o.client.Put(path, order, &resource)
	return resource, err
}

func (o *OrderServiceOp) Delete(orderID int64, options interface{}) (*Order, error) {
	path := fmt.Sprintf("%s/%d", ordersBasePath, orderID)
	resource := new(Order)
	err := o.client.Delete(path, options, &resource)
	return resource, err
}

func (o *OrderServiceOp) Batch(data OrderBatchOption) (*OrderBatchResource, error) {
	path := fmt.Sprintf("%s/batch", ordersBasePath)
	resource := new(OrderBatchResource)
	err := o.client.Post(path, data, &resource)
	return resource, err
}
