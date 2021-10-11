package usecase

import (
	"context"

	"github.com/egaevan/merchant/model"
)

type ProductUsecae interface {
	GetOneProduct(context.Context, int) (*model.Product, error)
	GetProduct(context.Context) ([]model.Product, error)
	SendProduct(context.Context, model.Product) (*model.Product, error)
	UpdateProduct(context.Context, model.ProductUpdate, int) (*model.ProductUpdate, error)
	DeleteProduct(context.Context, int) error
}
