package repository

import (
	"context"
	"net/http"

	"github.com/egaevan/merchant/model"
)

type UserRepository interface {
	Login(context.Context, http.ResponseWriter, *http.Request)
	FindOne(context.Context, string, string) map[string]interface{}
	// Fetch(context.Context)
	// Store(context.Context)
	// Delete(context.Context)
	// Aa
}

type ProductRepository interface {
	FindOne(context.Context, int) (*model.Product, error)
	Fetch(context.Context) ([]model.Product, error)
	Store(context.Context, model.Product) error
	Update(context.Context, model.ProductUpdate, int) error
	Delete(context.Context, int) error
}
