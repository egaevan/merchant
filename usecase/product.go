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

func (p *Product) GetOneProduct(ctx context.Context, productID int) (*model.Product, error) {

	prod, err := p.ProductRepo.FindOne(ctx, productID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return prod, nil
}

func (p *Product) GetProduct(ctx context.Context) ([]model.Product, error) {

	prod, err := p.ProductRepo.Fetch(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	productList := []model.Product{}

	for _, v := range prod {
		var product model.Product

		product.ID = v.ID
		product.Name = v.Name
		product.Sku = v.Sku
		product.Path = v.Path

		productList = append(productList, product)
	}

	return productList, nil
}

func (p *Product) SendProduct(ctx context.Context, product model.Product) (*model.Product, error) {

	err := p.ProductRepo.Store(ctx, product)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &product, nil
}

func (p *Product) UpdateProduct(ctx context.Context, product model.ProductUpdate, productID int) (*model.ProductUpdate, error) {

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
