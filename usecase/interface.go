package usecase

import (
	"context"

	"github.com/egaevan/merchant/model"
)

type ProductUsecae interface {
	GetOneProduct(context.Context, int) (*model.ProductDetail, error)
	GetProduct(context.Context) ([]model.ProductDetail, error)
	SendProduct(context.Context, model.ProductDetail) (*model.ProductDetail, error)
	UpdateProduct(context.Context, model.ProductDetail, int) (*model.ProductDetail, error)
	DeleteProduct(context.Context, int) error
}

type UserUsecae interface {
	Login(context.Context, model.User) (model.User, error)
	CreateUser(context.Context, model.User) error
	UpdateUser(context.Context, int, model.User) error
	DeleteUser(context.Context, int) error
}
