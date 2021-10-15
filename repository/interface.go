package repository

import (
	"context"

	"github.com/egaevan/merchant/model"
)

type UserRepository interface {
	FindOne(context.Context, string, string) (model.User, error)
	Fetch(context.Context) error
	Store(context.Context, model.User) error
	Update(context.Context, int, model.User) error
	Delete(context.Context, int) error
}

type ProductRepository interface {
	FindOne(context.Context, int) (*model.ProductDetail, error)
	Fetch(context.Context) ([]model.ProductDetail, error)
	Store(context.Context, model.ProductDetail) error
	Update(context.Context, model.ProductDetail, int) error
	Delete(context.Context, int) error
}
