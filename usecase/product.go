package usecase

import (
	"context"

	"github.com/egaevan/merchant/model"
	"github.com/egaevan/merchant/repository"
	log "github.com/sirupsen/logrus"
)

type Product struct {
	ProductRepo repository.ProductRepository
}

func NewProduct(productRepo repository.ProductRepository) ProductUsecae {
	return &Product{
		ProductRepo: productRepo,
	}
}

func (p *Product) GetOneProduct(ctx context.Context, productID int) (*model.ProductDetail, error) {

	prod, err := p.ProductRepo.FindOne(ctx, productID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return prod, nil
}

func (p *Product) GetProduct(ctx context.Context) ([]model.ProductDetail, error) {

	prod, err := p.ProductRepo.Fetch(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	productList := []model.ProductDetail{}

	for _, v := range prod {
		var product model.ProductDetail

		product.Product.ID = v.Product.ID
		product.Product.Name = v.Product.Name
		product.Product.Sku = v.Product.Sku
		product.Product.Path = v.Product.Path
		product.Price = v.Price
		product.Stock = v.Stock

		productList = append(productList, product)
	}

	return productList, nil
}

func (p *Product) SendProduct(ctx context.Context, product model.ProductDetail) (*model.ProductDetail, error) {

	err := p.ProductRepo.Store(ctx, product)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &product, nil
}

func (p *Product) UpdateProduct(ctx context.Context, product model.ProductDetail, productID int) (*model.ProductDetail, error) {

	err := p.ProductRepo.Update(ctx, product, productID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &product, nil
}

func (p *Product) DeleteProduct(ctx context.Context, productID int) error {

	err := p.ProductRepo.Delete(ctx, productID)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
