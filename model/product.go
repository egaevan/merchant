package model

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Sku  string `json:"sku"`
	Path string `json:"path"`
}

type ProductUpdate struct {
	Name string `json:"name"`
	Sku  string `json:"sku"`
	Path string `json:"path"`
}

type ProductDetail struct {
	Product Product `json:"product"`
	Price   int     `json:"price"`
	Stock   int     `json:"stock"`
}
