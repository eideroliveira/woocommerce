package woocommerce

// ProductVariationService is an interface for interfacing with the product variations endpoints of woocommerce API
// https://woocommerce.github.io/woocommerce-rest-api-docs/#product-variations
type ProductVariationService interface {
	Create()
	Get()
	Delete()
	List()
	Update()
}

type ProductVariationServiceOp struct {
	client *Client
}

func (p *ProductVariationServiceOp) Create() {

}

func (p *ProductVariationServiceOp) Get() {

}

func (p *ProductVariationServiceOp) Delete() {

}

func (p *ProductVariationServiceOp) List() {

}

func (p *ProductVariationServiceOp) Update() {

}

func (p *ProductVariationServiceOp) BatchUpdate() {

}
